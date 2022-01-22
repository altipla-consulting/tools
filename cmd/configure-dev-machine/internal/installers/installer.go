package installers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/collections"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/config"
)

var (
	register = []Installer{
		new(insApt),
		new(insGo),
		new(insGoprivate),
		new(insNode),
		new(insDC),
		new(insTools),
		new(insBuf),
		new(insActools),
		new(insGcloud),
		new(insMkcert),
		new(insStern),
		new(insEnospc),
		new(insIPV4Forwarding),
		new(insNpmpackages),
		new(insUpdater),
	}

	NeedsInstall = &CheckResult{Install: true}
)

type Installer interface {
	Name() string
	Check() (*CheckResult, error)
	Install() error
	BashRC() string
}

type CheckResult struct {
	Install bool
}

func Run(filter string) error {
	bashrc := []string{
		"#!/bin/bash",
		"",
	}
	for _, installer := range register {
		if filter != "" && installer.Name() != filter {
			continue
		}

		fmt.Println()
		fmt.Println()
		log.Info("======================================================================")
		log.Info(">>> install ", installer.Name())
		log.Info("======================================================================")

		result, err := installer.Check()
		if err != nil {
			return errors.Trace(err)
		}
		if result != nil && result.Install {
			if err := installer.Install(); err != nil {
				return errors.Trace(err)
			}
		}

		if script := installer.BashRC(); script != "" {
			bashrc = append(bashrc, "# "+installer.Name())
			bashrc = append(bashrc, strings.TrimSpace(script))
			bashrc = append(bashrc, "")
		}
	}

	fmt.Println()
	fmt.Println()
	log.Info("======================================================================")
	log.Info(">>> configure ~/.bashrc file")
	log.Info("======================================================================")
	filename, err := config.UserFile("bashrc.sh")
	if err != nil {
		return errors.Trace(err)
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return errors.Trace(err)
	}
	if err := ioutil.WriteFile(filename, []byte(strings.Join(bashrc, "\n")), 0700); err != nil {
		return errors.Trace(err)
	}

	if err := configureBashRC(); err != nil {
		return errors.Trace(err)
	}

	fmt.Println()
	fmt.Println()
	log.Info("======================================================================")
	log.Info("======================================================================")
	log.Info()
	log.Info("Finished successfully!")
	log.Info()

	return nil
}

// Clean old lines from bashrc that shouldn't be configured that way anymore
// and install the CDM script.
func configureBashRC() error {
	const installLine = "source ~/.config/configure-dev-machine/bashrc.sh"

	usr, err := os.UserHomeDir()
	if err != nil {
		return errors.Trace(err)
	}
	bashrc, err := ioutil.ReadFile(filepath.Join(usr, ".bashrc"))
	if err != nil {
		return errors.Trace(err)
	}

	remove := []string{
		`export CONFIGURE_DEV_MACHINE_UPDATER=true`,
		`configure-dev-machine check-updates`,
		`export USR_ID=$(id -u)`,
		`export GRP_ID=$(id -g)`,
		`alias dc='docker-compose'`,
		`alias dcrun='docker-compose run --rm'`,
		`alias dps='docker ps --format="table {{.ID}}\t{{.Names}}\t{{.Ports}}\t{{.Status}}"'`,
		`export CDM_GCLOUD=1`,
		`alias compute='gcloud compute'`,
		`export KUBE_EDITOR=gedit`,
		`alias k='kubectl'`,
		`alias kls='kubectl config get-contexts'`,
		`alias kuse='kubectl config use-context'`,
		`alias kpods='kubectl get pods --field-selector=status.phase!=Succeeded -o wide'`,
		`alias knodes='kubectl get nodes -o wide'`,
		`source <(kubectl completion bash | sed 's/kubectl/k/g')`,
		`export GOROOT=/usr/local/go`,
		`export PATH=$PATH:$GOROOT/bin:$HOME/go/bin`,
		`source ~/.config/configure-dev-machine/bashrc.shexport CDM_GCLOUD=1`,
	}

	var result []string
	var cdm bool
	for _, line := range strings.Split(string(bashrc), "\n") {
		if collections.HasString(remove, line) {
			log.WithField("line", line).Warning("Removing old configuration line from .bashrc")
		} else {
			result = append(result, line)
		}

		if line == installLine {
			cdm = true
		}
	}
	if !cdm {
		result = append(result, installLine)
	}

	if err := ioutil.WriteFile(filepath.Join(usr, ".bashrc"), []byte(strings.Join(result, "\n")), 0744); err != nil {
		return errors.Trace(err)
	}

	return nil
}
