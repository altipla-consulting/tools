package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var cmdMajor = &cobra.Command{
	Use:     "major",
	Short:   "Release a new major version",
	Example: "releaser major",
	RunE: func(command *cobra.Command, args []string) error {
		return errors.Trace(Release("major"))
	},
}
