package installers

import (
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

type insNpmpackages struct{}

func (ins *insNpmpackages) Name() string {
	return "npmpackages"
}

func (ins *insNpmpackages) Check() (*CheckResult, error) {
	return NeedsInstall, nil
}

func (ins *insNpmpackages) Install() error {
	// We need to install NPM in a different batch because any update will make
	// the next packages to miss the files npm itself needs because of the update.
	script := `
		sudo rm -rf /usr/lib/node_modules/yarn /usr/lib/node_modules/netlify-cli
		sudo npm install -g npm@latest
		sudo npm install -g yarn@latest
		sudo npm install -g --unsafe-perm=true netlify-cli@latest
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insNpmpackages) BashRC() string {
	return ""
}
