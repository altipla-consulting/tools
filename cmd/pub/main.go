package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/pub/internal/publish"
)

func main() {
	var flagDebug bool
	flag.BoolVarP(&flagDebug, "debug", "d", false, "Debug logging")
	flag.Parse()

	log.SetFormatter(new(log.TextFormatter))
	if flagDebug {
		log.SetLevel(log.DebugLevel)
		log.Debug("DEBUG log level activated")
	}

	if err := publish.Run(); err != nil {
		log.Error(err.Error())
		log.Debug(errors.Stack(err))
		os.Exit(1)
	}
}
