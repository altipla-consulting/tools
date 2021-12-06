package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insCLI struct{}

func (ins *insCLI) Name() string {
	return "cli"
}

func (ins *insCLI) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insCLI) Install() error {
	script := `
		curl -sL https://tools.altipla.consulting/install/cli | bash
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insCLI) BashRC() string {
	return ""
}
