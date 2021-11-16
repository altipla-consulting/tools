package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insGoprivate struct{}

func (ins *insGoprivate) Name() string {
	return "goprivate"
}

func (ins *insGoprivate) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insGoprivate) Install() error {
	script := `
		git config --global url."ssh://git@github.com:".insteadOf "https://github.com"
		/usr/local/go/bin/go env -w GOPRIVATE=github.com/lavozdealmeria,github.com/altipla-consulting
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insGoprivate) BashRC() string {
	return ""
}
