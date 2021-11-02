package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var cmdPatch = &cobra.Command{
	Use:     "patch",
	Short:   "Release a new patch version",
	Example: "releaser patch",
	RunE: func(command *cobra.Command, args []string) error {
		return errors.Trace(Release("patch"))
	},
}
