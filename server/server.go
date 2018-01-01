package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/thesoenke/go-coin"
)

const (
	protocol      = "tcp"
	nodeVersion   = 1
	commandLength = 12
)

var (
	nodeAddress     string
	miningAddress   string
	knownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	mempool         = make(map[string]coin.Transaction)
)

type addr struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

type getblocks struct {
	AddrFrom string
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

// Start server to run a node
func Start(nodeID int, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%d", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc, err := coin.NewBlockchain(nodeID)
	if err != nil {
		log.Fatal(err)
	}

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

// SendTx sends a transaction to another node
func SendTx(tx *coin.Transaction) {
	sendTx(knownNodes[0], tx)
}

func sendAddr(address string) {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *coin.Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		removeNode(addr)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{AddrFrom: nodeAddress, Type: kind, Items: items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{AddrFrom: nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(addr string, tnx *coin.Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func sendVersion(addr string, bc *coin.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{
		Version:    nodeVersion,
		BestHeight: bestHeight,
		AddrFrom:   nodeAddress,
	})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func removeNode(addr string) {
	var updatedNodes []string
	for _, node := range knownNodes {
		if node != addr {
			updatedNodes = append(updatedNodes, node)
		}
	}

	knownNodes = updatedNodes
}
