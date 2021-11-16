package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insApt struct{}

func (ins *insApt) Name() string {
	return "apt"
}

func (ins *insApt) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insApt) Install() error {
	// autoconf to build sass compiler
	// build-essential contains make
	// libnss3-tools is for mkcert, to install certs to Chrome
	script := `
    sudo apt update
    sudo apt install -y wget tar curl autoconf jq git build-essential libnss3-tools
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insApt) BashRC() string {
	return ""
}
