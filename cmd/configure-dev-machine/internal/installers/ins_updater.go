package installers

import (
	"libs.altipla.consulting/errors"
	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insUpdater struct{}

func (ins *insUpdater) Name() string {
	return "updater"
}

func (ins *insUpdater) Check() (*CheckResult, error) {
	return nil, nil
}

func (ins *insUpdater) Install() error {
	script := `
		echo REPLACE_VERSION | sudo tee /etc/configure-dev-machine/installed-version
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insUpdater) BashRC() string {
	return "configure-dev-machine check-updates"
}
