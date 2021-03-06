package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/go-coin"
	"github.com/thesoenke/go-coin/server"
)

var minerAddress string
var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Start a new node server",
	Run: func(cmd *cobra.Command, args []string) {
		if !coin.ValidateAddress(minerAddress) {
			err := fmt.Errorf("miner address is not valid")
			printErr(err)
		}

		fmt.Printf("Started mining. Address to receive rewards: %s\n", minerAddress)
		err := server.Start(nodeID, minerAddress)
		printErr(err)
	},
}

func init() {
	cmdServer.PersistentFlags().StringVar(&minerAddress, "address", "", "Address of the miner for rewards")
	RootCmd.AddCommand(cmdServer)
}
