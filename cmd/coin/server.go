package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/go-coin"
	"github.com/thesoenke/go-coin/server"
)

var nodeID string
var minerAddress string
var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Start a new node",
	Run: func(cmd *cobra.Command, args []string) {
		if !coin.ValidateAddress(minerAddress) {
			err := fmt.Errorf("miner address is not valid")
			printErr(err)
		}

		fmt.Printf("Started mining. Address to receive rewards: %s\n", minerAddress)
		server.StartServer(nodeID, minerAddress)
	},
}

func init() {
	cmdServer.PersistentFlags().StringVar(&nodeID, "node", "", "ID of the node to identify on a single machine")
	cmdServer.PersistentFlags().StringVar(&minerAddress, "address", "", "Address of the miner for rewards")
	RootCmd.AddCommand(cmdServer)
}
