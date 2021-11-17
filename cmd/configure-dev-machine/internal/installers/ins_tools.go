package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insTools struct{}

func (ins *insTools) Name() string {
	return "tools"
}

func (ins *insTools) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insTools) Install() error {
	script := `
		curl -s https://tools.altipla.consulting/install/all | sudo bash
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insTools) BashRC() string {
	return ""
}
