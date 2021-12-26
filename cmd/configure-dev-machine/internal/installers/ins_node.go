package installers

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

const wantedNode = "16"

type insNode struct{}

func (ins *insNode) Name() string {
	return "node"
}

func (ins *insNode) Check() (*CheckResult, error) {
	if _, err := exec.LookPath("node"); err != nil {
		log.Info("not found")
		return NeedsInstall, nil
	}

	output, err := run.InteractiveCaptureOutput("node", "-v")
	if err != nil {
		return nil, errors.Trace(err)
	}
	version := strings.Split(output, ".")[0][1:]

	if version != wantedNode {
		log.WithFields(log.Fields{
			"wanted": wantedNode,
			"found":  version,
		}).Info("update node")

		return NeedsInstall, nil
	}
	return nil, nil
}

func (ins *insNode) Install() error {
	script := `
    curl -sL https://deb.nodesource.com/setup_$VERSION.x | sudo -E bash -
    sudo apt install -y nodejs
    node -v
  `
	vars := map[string]string{"VERSION": wantedNode}
	return errors.Trace(run.Shell(script, vars))
}

func (ins *insNode) BashRC() string {
	return ""
}
