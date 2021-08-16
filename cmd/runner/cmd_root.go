package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flagDebug bool

func init() {
	cmdRoot.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "Enable debug logging for this tool")
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
