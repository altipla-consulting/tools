package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insGcloud struct{}

func (ins *insGcloud) Name() string {
	return "gcloud"
}

func (ins *insGcloud) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insGcloud) Install() error {
	script := `
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
    sudo apt update
    sudo apt install -y google-cloud-sdk kubectl
  `
	if err := run.Shell(script); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (ins *insGcloud) BashRC() string {
	return `
alias compute='gcloud compute'

export KUBE_EDITOR=nano
alias k='kubectl'
alias kls='kubectl config get-contexts'
alias kuse='kubectl config use-context'
alias kpods='kubectl get pods --field-selector=status.phase!=Succeeded -o wide'
alias knodes='kubectl get nodes -o wide'

source <(kubectl completion bash | sed 's/kubectl/k/g')
`
}
