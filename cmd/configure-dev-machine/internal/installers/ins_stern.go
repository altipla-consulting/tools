package installers

import (
	"os/exec"
	"strings"

	"libs.altipla.consulting/errors"
	log "github.com/sirupsen/logrus"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

const wantedStern = "1.11.0"

type insStern struct{}

func (ins *insStern) Name() string {
	return "stern"
}

func (ins *insStern) Check() (*CheckResult, error) {
	if _, err := exec.LookPath("stern"); err != nil {
		log.Info("not found")
		return NeedsInstall, nil
	}

	output, err := run.InteractiveCaptureOutput("stern", "-v")
	if err != nil {
		return nil, errors.Trace(err)
	}
	version := strings.Split(output, " ")[2]

	if version != wantedStern {
		log.WithFields(log.Fields{
			"wanted": wantedStern,
			"found":  version,
		}).Info("update stern")

		return NeedsInstall, nil
	}
	return nil, nil
}

func (ins *insStern) Install() error {
	script := `
		curl -L -o ~/bin/stern https://github.com/wercker/stern/releases/download/$VERSION/stern_linux_amd64
    chmod +x ~/bin/stern
	`
	vars := map[string]string{"VERSION": wantedStern}
	return errors.Trace(run.Shell(script, vars))
}

func (ins *insStern) BashRC() string {
	return ""
}
