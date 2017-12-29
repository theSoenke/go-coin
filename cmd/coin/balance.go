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
		balance := getBalance(balanceAddress)
		fmt.Printf("Balance of '%s': %d\n", balanceAddress, balance)
	},
}

func init() {
	cmdBalance.PersistentFlags().StringVar(&balanceAddress, "address", "", "Address to calculate balance for")
	RootCmd.AddCommand(cmdBalance)
}

func getBalance(address string) int {

	if !coin.ValidateAddress(address) {
		err := fmt.Errorf("address '%s' is not valid", address)
		printErr(err)
	}

	bc, err := coin.NewBlockchain("1")
	printErr(err)
	defer bc.DB.Close()

	balance := 0
	pubKeyHash := coin.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOSet := coin.UTXOSet{Blockchain: bc}
	UTXOs, err := UTXOSet.FindUTXO(pubKeyHash)
	printErr(err)

	for _, out := range UTXOs {
		balance += out.Value
	}

	return balance
}
