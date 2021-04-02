package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/ci/internal/commands"
)

func main() {
	if err := commands.CmdRoot.Execute(); err != nil {
		log.Error(err.Error())
		log.Debug(errors.Stack(err))
		os.Exit(1)
	}
}
