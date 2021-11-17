package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/installers"
)

var filter string

func init() {
	CmdInstall.PersistentFlags().StringVarP(&filter, "filter", "f", "", "Filter installers to run by name.")
	CmdRoot.AddCommand(CmdInstall)
}

var CmdInstall = &cobra.Command{
	Use:          "install",
	Short:        "Run all the installers.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.Trace(installers.Run(filter))
	},
}
