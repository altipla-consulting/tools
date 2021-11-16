package installers

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

const wantedDC = "1.25.4"

type insDC struct{}

func (ins *insDC) Name() string {
	return "docker-compose"
}

func (ins *insDC) Check() (*CheckResult, error) {
	if _, err := exec.LookPath("docker-compose"); err != nil {
		log.Info("not found")
		return NeedsInstall, nil
	}

	version, err := run.InteractiveCaptureOutput("docker-compose", "version", "--short")
	if err != nil {
		return nil, errors.Trace(err)
	}

	if version != wantedDC {
		log.WithFields(log.Fields{
			"wanted": wantedDC,
			"found":  version,
		}).Info("update docker-compose")

		return NeedsInstall, nil
	}
	return nil, nil
}

func (ins *insDC) Install() error {
	script := `
    sudo curl -L "https://github.com/docker/compose/releases/download/$VERSION/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    docker-compose --version
  `
	vars := map[string]string{"VERSION": wantedDC}
	if err := run.Shell(script, vars); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (ins *insDC) BashRC() string {
	return `
export USR_ID=$(id -u)
export GRP_ID=$(id -g)

alias dc='docker-compose'
alias dcrun='docker-compose run --rm'
alias dps='docker ps --format="table {{.ID}}\t{{.Names}}\t{{.Ports}}\t{{.Status}}"'
	`
}
