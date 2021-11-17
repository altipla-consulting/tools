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
		o.AddLine(aurora.Blue("curl -sL https://tools.altipla.consulting/install/configure-dev-machine | bash"))
		o.Render()

		return nil
	},
}
