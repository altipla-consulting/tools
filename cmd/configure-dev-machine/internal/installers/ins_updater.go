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
	return NeedsInstall, nil
}

func (ins *insUpdater) Install() error {
	script := `
		sudo mkdir -p /etc/configure-dev-machine
		echo REPLACE_VERSION | sudo tee /etc/configure-dev-machine/installed-version
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insUpdater) BashRC() string {
	return "/usr/local/bin/configure-dev-machine check-updates"
}
