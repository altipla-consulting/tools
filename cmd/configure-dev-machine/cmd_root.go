package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CmdRoot = &cobra.Command{
	Use:   "configure-dev-machine",
	Short: "Instalador y configurador para los ordenadores de desarrolladores",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if Version == "dev" {
			log.Warning("Running development version. To download a production version run: curl https://tools.altipla.consulting/install/configure-dev-machine | bash")
		}

		return nil
	},
}
