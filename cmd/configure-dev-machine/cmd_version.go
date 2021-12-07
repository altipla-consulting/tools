package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/version"
)

func init() {
	CmdRoot.AddCommand(CmdVersion)
}

var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print current installed version and latest available one.",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := version.FetchInfo()
		if err != nil {
			return errors.Trace(err)
		}
		fmt.Println("Currently Installed:", info.CurrentlyInstalled)
		fmt.Println("Latest:", info.Latest)
		return nil
	},
}
