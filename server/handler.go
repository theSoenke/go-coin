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

func handleGetData(request []byte, bc *coin.Blockchain) error {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		fmt.Println("invalid payload")
		return nil
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return err
		}

		err = sendBlock(payload.AddrFrom, &block)
		if err != nil {
			return err
		}
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		err = sendTx(payload.AddrFrom, &tx)
		if err != nil {
			fmt.Printf("Failed sending tx to %s\n", payload.AddrFrom)
			return nil
		}
		// delete(mempool, txID)
	}

	return nil
}

func handleTx(request []byte, bc *coin.Blockchain) error {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		return err
	}

	txData := payload.Transaction
	tx, err := coin.DeserializeTransaction(txData)
	if err != nil {
		return err
	}

	mempool[hex.EncodeToString(tx.ID)] = tx
	// Is the central node
	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				err = sendInv(node, "tx", [][]byte{tx.ID})
				if err != nil {
					fmt.Printf("Could not reach node %s\n", node)
				}
			}
		}
	} else {
		for len(mempool) >= transactionsInBlock {
			err = mineBlock(bc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mineBlock(bc *coin.Blockchain) error {
	var txs []*coin.Transaction

	for id := range mempool {
		tx := mempool[id]
		if bc.VerifyTransaction(&tx) {
			txs = append(txs, &tx)
		}

		if len(txs) == 0 {
			fmt.Println("all transactions are invalid")
			return nil
		}
	}

	cbTx := coin.NewCoinbaseTX(miningAddress, "")
	txs = append(txs, cbTx)
	newBlock, err := bc.MineBlock(txs)
	if err != nil {
		return err
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
			err = sendInv(node, "block", [][]byte{newBlock.Hash})
			fmt.Printf("Could not reach node %s\n", node)
		}
	}

	return nil
}

func handleVersion(request []byte, bc *coin.Blockchain) error {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		return err
	}

	myBestHeight, err := bc.GetBestHeight()
	if err != nil {
		return err
	}

	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		err = sendVersion(payload.AddrFrom, bc)
		if err != nil {
			return err
		}
	}

	// sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		fmt.Printf("New node %s connected\n", payload.AddrFrom)
		knownNodes = append(knownNodes, payload.AddrFrom)
	}

	return nil
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
