package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/go-coin"
)

var sendFrom string
var sendTo string
var sendAmount int
var cmdSend = &cobra.Command{
	Use:   "send",
	Short: "Send a transaction to an address",
	Run: func(cmd *cobra.Command, args []string) {
		if !coin.ValidateAddress(sendFrom) {
			err := fmt.Errorf("sender address '%s' is not valid", sendFrom)
			printErr(err)
		}
		if !coin.ValidateAddress(sendTo) {
			err := fmt.Errorf("receiver address '%s' is not valid", sendTo)
			printErr(err)
		}

		if sendAmount < 0 {
			err := fmt.Errorf("amount needs to be > 0")
			printErr(err)

		}

		bc, err := coin.NewBlockchain("1")
		printErr(err)
		defer bc.DB.Close()

		UTXOSet := coin.UTXOSet{Blockchain: bc}

		tx, err := coin.NewUTXOTransaction(sendFrom, sendTo, sendAmount, &UTXOSet)
		printErr(err)
		cbTx := coin.NewCoinbaseTX(sendFrom, "")
		txs := []*coin.Transaction{cbTx, tx}

		newBlock, err := bc.MineBlock(txs)
		printErr(err)

		err = UTXOSet.Update(newBlock)
		printErr(err)
		fmt.Println("Success!")
	},
}

func init() {
	cmdSend.PersistentFlags().StringVar(&sendFrom, "from", "", "Sender of the transaction")
	cmdSend.PersistentFlags().StringVar(&sendTo, "to", "", "Receiver of the transaction")
	cmdSend.PersistentFlags().IntVar(&sendAmount, "amount", 0, "Amount that will be send")
	RootCmd.AddCommand(cmdSend)
}
