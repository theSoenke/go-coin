package main

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var genesisRewardAddress string
var cmdCreate = &cobra.Command{
	Use:   "create-chain",
	Short: "Create a new Blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		bc, err := coin.CreateBlockchain(genesisRewardAddress)
		printErr(err)
		defer bc.DB.Close()

		UTXOSet := coin.UTXOSet{Blockchain: bc}
		UTXOSet.Reindex()
	},
}

func init() {
	cmdCreate.PersistentFlags().StringVar(&genesisRewardAddress, "address", "", "Create a Blockchain and send genesis block reward to specified address")
	RootCmd.AddCommand(cmdCreate)
}
