package main

import (
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/box"
)

func init() {
	CmdRoot.AddCommand(CmdPrintUpdate)
}

var CmdPrintUpdate = &cobra.Command{
	Use:   "print-update",
	Short: "Print update command.",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := box.Box{}
		o.AddLine("Run the following command to update:")
		o.AddLine(aurora.Blue("sudo apt update && sudo apt install configure-dev-machine && configure-dev-machine install"))
		o.Render()

		return nil
	},
}
