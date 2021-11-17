package main

import (
	"github.com/spf13/cobra"
)

var CmdRoot = &cobra.Command{
	Use:   "configure-dev-machine",
	Short: "Configure Altipla development machines with an automated script.",
}
