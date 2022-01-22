package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insBuf struct{}

func (ins *insBuf) Name() string {
	return "buf"
}

func (ins *insBuf) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insBuf) Install() error {
	script := `
		curl -L https://github.com/bufbuild/buf/releases/download/v1.0.0-rc11/buf-Linux-x86_64 -o /tmp/buf
		sudo mv /tmp/buf /usr/local/bin/buf
		chmod +x /usr/local/bin/buf
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insBuf) BashRC() string {
	return ""
}
