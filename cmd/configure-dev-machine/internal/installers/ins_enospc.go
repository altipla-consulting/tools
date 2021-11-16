package installers

import (
	"io/ioutil"
	"strings"

	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

// insEnospc increases a Linux kernel quota on the allowed watchers. They are too low
// by default on Ubuntu.
//
// Fix is extracted from:
//   https://stackoverflow.com/questions/22475849/node-js-what-is-enospc-error-and-how-to-solve
type insEnospc struct{}

func (ins *insEnospc) Name() string {
	return "enospc"
}

func (ins *insEnospc) Check() (*CheckResult, error) {
	content, err := ioutil.ReadFile("/etc/sysctl.conf")
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "fs.inotify.max_user_watches") {
			return nil, nil
		}
	}
	return NeedsInstall, nil
}

func (ins *insEnospc) Install() error {
	script := `
		echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
		sudo sysctl --system
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insEnospc) BashRC() string {
	return ""
}
