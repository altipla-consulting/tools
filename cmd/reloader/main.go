package main

import (
	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func main() {
	if err := CmdRoot.Execute(); err != nil {
		log.Fatal(errors.Stack(err))
	}
}
