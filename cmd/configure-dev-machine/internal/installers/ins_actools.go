package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insActools struct{}

func (ins *insActools) Name() string {
	return "actools"
}

func (ins *insActools) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insActools) Install() error {
	script := `
	  sudo curl https://tools.altipla.consulting/bin/actools -o /usr/bin/actools
	  sudo chmod +x /usr/bin/actools

	  actools pull
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insActools) BashRC() string {
	return ""
}
