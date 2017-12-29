package main

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var cmdLog = &cobra.Command{
	Use:   "log",
	Short: "Print the Blockchain log",
	RunE: func(cmd *cobra.Command, args []string) error {
		bc, err := coin.NewBlockchain()
		if err != nil {
			return err
		}

		defer bc.DB.Close()
		bc.Print()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(cmdLog)
}
