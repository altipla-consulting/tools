package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/pub/internal/publish"
)

func main() {
	if err := publish.Run(); err != nil {
		log.Error(err.Error())
		log.Debug(errors.Stack(err))
		os.Exit(1)
	}
}
