package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List addresses stored in wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, err := coin.NewWallets()
		if err != nil {
			printErr(err)
		}

		addresses := wallets.GetAddresses()
		for _, address := range addresses {
			balance := getBalance(address)
			fmt.Printf("Address: %s Balance: %d\n", address, balance)
		}
	},
}

func init() {
	RootCmd.AddCommand(cmdList)
}
