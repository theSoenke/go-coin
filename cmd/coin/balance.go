package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var balanceAddress string
var cmdBalance = &cobra.Command{
	Use:   "balance",
	Short: "Get balance of address",
	Run: func(cmd *cobra.Command, args []string) {
		if !coin.ValidateAddress(address) {
			log.Panic("ERROR: Address is not valid")
		}
		bc, err := coin.NewBlockchain()
		printErr(err)
		defer bc.DB.Close()

		balance := 0
		pubKeyHash := coin.Base58Decode([]byte(address))
		pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
		UTXOs := bc.FindUTXO(pubKeyHash)

		for _, out := range UTXOs {
			balance += out.Value
		}

		fmt.Printf("Balance of '%s': %d\n", address, balance)
	},
}

func init() {
	cmdBalance.PersistentFlags().StringVar(&balanceAddress, "address", "", "Address to calculate balance for")
	RootCmd.AddCommand(cmdBalance)
}
