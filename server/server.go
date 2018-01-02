package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/thesoenke/go-coin"
)

const (
	protocol            = "tcp"
	nodeVersion         = 1
	commandLength       = 12
	transactionsInBlock = 2
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
func Start(nodeID int, minerAddress string) error {
	nodeAddress = fmt.Sprintf("localhost:%d", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		return err
	}

	defer ln.Close()
	bc, err := coin.NewBlockchain(nodeID)
	if err != nil {
		return err
	}

	if nodeAddress != knownNodes[0] {
		err = sendVersion(knownNodes[0], bc)
		if err != nil {
			fmt.Println("Failed sending version to central node")
			return err
		}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn, bc)
	}
}

// SendTx sends a transaction to another node
func SendTx(tx *coin.Transaction) {
	err := sendTx(knownNodes[0], tx)
	if err != nil {
		fmt.Println("Failed sending transaction to central node")
	}
}

func sendAddr(address string) error {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	err := sendData(address, request)
	return err
}

func sendBlock(addr string, b *coin.Block) error {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	err := sendData(addr, request)
	return err
}

func sendData(addr string, data []byte) error {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		removeNode(addr)
		return err
	}

	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	return err
}

func sendInv(address, kind string, items [][]byte) error {
	inventory := inv{AddrFrom: nodeAddress, Type: kind, Items: items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	err := sendData(address, request)
	return err
}

func sendGetBlocks(address string) error {
	payload := gobEncode(getblocks{AddrFrom: nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	err := sendData(address, request)
	return err
}

func sendGetData(address, kind string, id []byte) error {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	err := sendData(address, request)
	return err
}

func sendTx(addr string, tnx *coin.Transaction) error {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	err := sendData(addr, request)
	return err
}

func sendVersion(addr string, bc *coin.Blockchain) error {
	bestHeight, err := bc.GetBestHeight()
	if err != nil {
		return err
	}

	payload := gobEncode(version{
		Version:    nodeVersion,
		BestHeight: bestHeight,
		AddrFrom:   nodeAddress,
	})

	request := append(commandToBytes("version"), payload...)

	err = sendData(addr, request)
	return err
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
