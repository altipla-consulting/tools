package installers

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

const wantedGo = "1.17.3"

type insGo struct{}

func (ins *insGo) Name() string {
	return "go"
}

func (ins *insGo) Check() (*CheckResult, error) {
	if _, err := exec.LookPath("go"); err != nil {
		log.Info("not found")
		return NeedsInstall, nil
	}

	output, err := run.InteractiveCaptureOutput("go", "version")
	if err != nil {
		return nil, errors.Trace(err)
	}
	version := strings.Split(output, " ")[2]

	if version != "go"+wantedGo {
		log.WithFields(log.Fields{
			"wanted": "go" + wantedGo,
			"found":  version,
		}).Info("update go")

		return NeedsInstall, nil
	}
	return nil, nil
}

func (ins *insGo) Install() error {
	script := `
    sudo rm -rf /usr/local/go
    wget -q -O /tmp/go.tar.gz "https://dl.google.com/go/go$VERSION.linux-amd64.tar.gz"
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    /usr/local/go/bin/go version
  `
	vars := map[string]string{"VERSION": wantedGo}
	if err := run.Shell(script, vars); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (ins *insGo) BashRC() string {
	return `
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$HOME/go/bin
`
}
