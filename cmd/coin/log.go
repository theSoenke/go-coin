package main

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/go-coin"
)

var cmdLog = &cobra.Command{
	Use:   "log",
	Short: "Print the Blockchain log",
	Run: func(cmd *cobra.Command, args []string) {
		bc, err := coin.NewBlockchain("1")
		printErr(err)

		defer bc.DB.Close()
		bc.Print()
	},
}

func init() {
	RootCmd.AddCommand(cmdLog)
}
