package main

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/coin"
)

var cmdAddBlock = &cobra.Command{
	Use:   "add",
	Short: "Add a block to the Blockchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		bc, err := coin.NewBlockchain()
		if err != nil {
			return err
		}

		defer bc.DB.Close()
		data := args[0]
		bc.AddBlock(data)
		return nil
	},
}

func init() {
	cmdAddBlock.Args = cobra.ExactArgs(1)
	RootCmd.AddCommand(cmdAddBlock)
}
