package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/go-coin"
)

var genesisRewardAddress string
var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Create a new Blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		if genesisRewardAddress == "" {
			err := fmt.Errorf("genesis reward address cannot be empty")
			printErr(err)
		}

		bc, err := coin.CreateBlockchain(genesisRewardAddress, genesisNodeID)
		printErr(err)
		defer bc.DB.Close()

		UTXOSet := coin.UTXOSet{Blockchain: bc}
		UTXOSet.Reindex()
	},
}

func init() {
	cmdInit.PersistentFlags().StringVar(&genesisRewardAddress, "address", "", "Create a Blockchain and send genesis block reward to specified address")
	RootCmd.AddCommand(cmdInit)
}
