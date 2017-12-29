package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var balanceAddress string
var cmdBalance = &cobra.Command{
	Use:   "balance",
	Short: "Get balance of address",
	Run: func(cmd *cobra.Command, args []string) {
		bc, err := coin.NewBlockchain()
		printErr(err)

		balance := 0
		UTXOs := bc.FindUTXO(balanceAddress)
		for _, out := range UTXOs {
			balance += out.Value
		}

		fmt.Printf("Balance of '%s': %d\n", balanceAddress, balance)
	},
}

func init() {
	cmdBalance.PersistentFlags().StringVar(&balanceAddress, "address", "", "Address to calculate balance for")
	RootCmd.AddCommand(cmdBalance)
}
