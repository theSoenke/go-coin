package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var cmdReindex = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex unspent transactions (UTXO)",
	Run: func(cmd *cobra.Command, args []string) {
		bc, err := coin.NewBlockchain("")
		printErr(err)

		UTXOSet := coin.UTXOSet{Blockchain: bc}
		err = UTXOSet.Reindex()
		printErr(err)

		count, err := UTXOSet.CountTransactions()
		printErr(err)
		fmt.Printf("Reindex of %d UTXO transactions successful\n", count)
	},
}

func init() {
	RootCmd.AddCommand(cmdReindex)
}
