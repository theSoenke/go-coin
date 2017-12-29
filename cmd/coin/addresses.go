package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var cmdAddresses = &cobra.Command{
	Use:   "addresses",
	Short: "List addresses stored in wallet file",
	Run: func(cmd *cobra.Command, args []string) {
		wallets, err := coin.NewWallets()
		if err != nil {
			printErr(err)
		}

		addresses := wallets.GetAddresses()
		for _, address := range addresses {
			fmt.Println(address)
		}
	},
}

func init() {
	RootCmd.AddCommand(cmdAddresses)
}
