package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func main() {
	if err := cmdRoot.Execute(); err != nil {
		log.Debug(errors.Stack(err))
		os.Exit(1)
	}
}
