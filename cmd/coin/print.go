package main

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var cmdPrint = &cobra.Command{
	Use:   "print",
	Short: "Print the blockchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		bc := coin.NewBlockchain()
		defer bc.DB.Close()
		bc.Print()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(cmdPrint)
}
