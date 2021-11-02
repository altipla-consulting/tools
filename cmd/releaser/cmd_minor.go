package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var cmdMinor = &cobra.Command{
	Use:     "minor",
	Short:   "Release a new minor version",
	Example: "releaser minor",
	RunE: func(command *cobra.Command, args []string) error {
		return errors.Trace(Release("minor"))
	},
}
