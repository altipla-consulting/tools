package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
	"tools.altipla.consulting/cmd/configure-dev-machine/internal/config"
)

func init() {
	CmdRoot.AddCommand(CmdVersion)
}

var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Versión actual del configurador",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := config.Version()
		if err != nil {
			return errors.Trace(err)
		}
		fmt.Println(version)
		return nil
	},
}
