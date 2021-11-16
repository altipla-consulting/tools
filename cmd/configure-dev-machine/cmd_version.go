package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	CmdRoot.AddCommand(CmdVersion)
}

var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Versión actual del configurador",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(Version)
		return nil
	},
}
