package main

import (
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/box"
)

func init() {
	CmdRoot.AddCommand(CmdUpdate)
}

var CmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Imprime el comando de actualización de la herramienta.",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := box.Box{}
		o.AddLine("Run the following command to update:")
		o.AddLine(aurora.Blue("curl https://tools.altipla.consulting/install/configure-dev-machine | bash"))
		o.Render()

		return nil
	},
}
