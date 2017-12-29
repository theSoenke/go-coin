package main

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var address string
var cmdCreate = &cobra.Command{
	Use:   "create-chain",
	Short: "Create a new Blockchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		bc, err := coin.CreateBlockchain(address)
		if err != nil {
			return err
		}

		defer bc.DB.Close()
		bc.Print()
		return nil
	},
}

func init() {
	cmdCreate.PersistentFlags().StringVar(&address, "address", "", "Create a Blockchain and send genesis block reward to specified address")
	RootCmd.AddCommand(cmdCreate)
}
