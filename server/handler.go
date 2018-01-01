package server

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	coin "github.com/thesoenke/go-coin"
)

func handleConnection(conn net.Conn, bc *coin.Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Printf("Unknown command '%s'\n", command)
	}

	conn.Close()
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There %d known nodes\n", len(knownNodes))
	requestBlocks()
}

func handleBlock(request []byte, bc *coin.Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := coin.DeserializeBlock(blockData)

	fmt.Printf("Received block %x with height %d\n", block.Hash, block.Height)
	err = bc.AddBlock(block)
	if err != nil {
		log.Fatal(err)
	}

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := coin.UTXOSet{Blockchain: bc}
		UTXOSet.Reindex()
	}
}

func handleInv(request []byte, bc *coin.Blockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *coin.Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(request []byte, bc *coin.Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(request []byte, bc *coin.Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx, err := coin.DeserializeTransaction(txData)
	if err != nil {
		log.Fatal(err)
	}

	mempool[hex.EncodeToString(tx.ID)] = tx
	// Is the central node
	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		for len(mempool) >= 2 {
			mineBlock(bc)
		}
	}
}

func mineBlock(bc *coin.Blockchain) {
	var txs []*coin.Transaction

	for id := range mempool {
		tx := mempool[id]
		if bc.VerifyTransaction(&tx) {
			txs = append(txs, &tx)
		}

		if len(txs) == 0 {
			fmt.Println("all transactions are invalid")
			return
		}
	}

	cbTx := coin.NewCoinbaseTX(miningAddress, "")
	txs = append(txs, cbTx)
	newBlock, err := bc.MineBlock(txs)
	if err != nil {
		log.Fatal(err)
	}

	UTXOSet := coin.UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()

	fmt.Printf("Mined new block with %d transactions\n", len(txs))

	for _, tx := range txs {
		txID := hex.EncodeToString(tx.ID)
		delete(mempool, txID)
	}

	for _, node := range knownNodes {
		if node != nodeAddress {
			sendInv(node, "block", [][]byte{newBlock.Hash})
		}
	}
}

func handleVersion(request []byte, bc *coin.Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	// sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		fmt.Printf("New node %s connected\n", payload.AddrFrom)
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}
