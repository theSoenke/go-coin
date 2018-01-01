package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/go-coin"
)

var cmdWallet = &cobra.Command{
	Use:   "address",
	Short: "Generate a new address",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, _ := coin.NewWallets(nodeID)
		address, err := wallets.CreateWallet()
		printErr(err)

		err = wallets.SaveToFile(nodeID)
		printErr(err)
		fmt.Printf("Your new address: %s\n", address)
	},
}

func init() {
	RootCmd.AddCommand(cmdWallet)
}
