package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"tools.altipla.consulting/cmd/runner/internal/build"
	"tools.altipla.consulting/cmd/runner/internal/deploy"
	"tools.altipla.consulting/cmd/runner/internal/preview"
)

var flagDebug bool

func init() {
	cmdRoot.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "Enable debug logging for this tool")
	cmdRoot.AddCommand(build.Cmd)
	cmdRoot.AddCommand(deploy.Cmd)
	cmdRoot.AddCommand(preview.Cmd)
}

var cmdRoot = &cobra.Command{
	Use:          "runner",
	Short:        "Publish Cloud Run applications.",
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(new(log.TextFormatter))
		if flagDebug {
			log.SetLevel(log.DebugLevel)
		}
	},
}
