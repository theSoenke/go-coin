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
		if !coin.ValidateAddress(balanceAddress) {
			err := fmt.Errorf("address '%s' is not valid", balanceAddress)
			printErr(err)
		}

		bc, err := coin.NewBlockchain()
		printErr(err)
		defer bc.DB.Close()

		balance := 0
		pubKeyHash := coin.Base58Decode([]byte(balanceAddress))
		pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
		UTXOSet := coin.UTXOSet{Blockchain: bc}
		UTXOs := UTXOSet.FindUTXO(pubKeyHash)

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
