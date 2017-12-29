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
	RunE: func(cmd *cobra.Command, args []string) error {
		bc, err := coin.NewBlockchain()
		if err != nil {
			return err
		}

		balance := 0
		UTXOs := bc.FindUTXO(balanceAddress)
		for _, out := range UTXOs {
			balance += out.Value
		}

		fmt.Printf("Balance of '%s': %d\n", balanceAddress, balance)
		return nil
	},
}

func init() {
	cmdBalance.PersistentFlags().StringVar(&balanceAddress, "address", "", "Address to calculate balance for")
	RootCmd.AddCommand(cmdBalance)
}
