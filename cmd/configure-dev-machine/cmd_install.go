package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/collections"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/config"
)

func init() {
	CmdRoot.AddCommand(CmdInstall)
}

var CmdInstall = &cobra.Command{
	Use:          "install",
	Short:        "Run all the installers.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		fmt.Println()
		log.Info("======================================================================")
		log.Info(">>> configure ~/.bashrc file")
		log.Info("======================================================================")
		filename, err := config.UserFile("bashrc.sh")
		if err != nil {
			return errors.Trace(err)
		}
		if _, err := os.Stat(filename); err != nil && !os.IsNotExist(err) {
			return errors.Trace(err)
		} else if err == nil {
			if err := os.Remove(filename); err != nil {
				return errors.Trace(err)
			}
		}

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
			"source ~/.config/configure-dev-machine/bashrc.sh",
		}

		var result []string
		for _, line := range strings.Split(string(bashrc), "\n") {
			if collections.HasString(remove, line) {
				log.WithField("line", line).Warning("Removing old configuration line from .bashrc")
			} else {
				result = append(result, line)
			}
		}
		if err := ioutil.WriteFile(filepath.Join(usr, ".bashrc"), []byte(strings.Join(result, "\n")), 0744); err != nil {
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
	},
}
