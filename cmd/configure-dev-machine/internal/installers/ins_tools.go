package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insTools struct{}

func (ins *insTools) Name() string {
	return "tools"
}

func (ins *insTools) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insTools) Install() error {
	script := `
	  rm -f ~/bin/ci ~/bin/configure-dev-machine ~/bin/gaestage ~/bin/gendc ~/bin/impsort ~/bin/jnet ~/bin/linter ~/bin/previewer-netlify ~/bin/pub ~/bin/releaser ~/bin/reloader ~/bin/wave ~/bin/altipla ~/bin/lavoz
		
		curl https://europe-west1-apt.pkg.dev/doc/repo-signing-key.gpg | sudo apt-key add -
		echo 'deb https://europe-west1-apt.pkg.dev/projects/altipla-tools acpublic main' | sudo tee /etc/apt/sources.list.d/acpublic.list
		sudo apt update
		sudo apt install -y tools/acpublic
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insTools) BashRC() string {
	return ""
}
