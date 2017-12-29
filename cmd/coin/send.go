package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var sendFrom string
var sendTo string
var sendAmount int
var cmdSend = &cobra.Command{
	Use:   "send",
	Short: "Send a transaction to an address",
	Run: func(cmd *cobra.Command, args []string) {
		bc, err := coin.NewBlockchain()
		printErr(err)
		defer bc.DB.Close()

		err = validateSendInput(sendFrom, sendTo, sendAmount)
		printErr(err)

		tx, err := coin.NewUTXOTransaction(sendFrom, sendTo, sendAmount, bc)
		printErr(err)

		err = bc.MineBlock([]*coin.Transaction{tx})
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

func validateSendInput(from string, to string, amount int) error {
	// TODO better error messages
	if from == "" || to == "" || amount <= 0 {
		return fmt.Errorf("Please check your input")
	}

	return nil
}
