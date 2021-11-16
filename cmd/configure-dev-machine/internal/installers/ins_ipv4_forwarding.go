package installers

import (
	"io/ioutil"
	"strings"

	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/run"
)

// insIPV4Forwarding fixes an error that disables networking on Docker.
//
// Fix is extracted from:
//   https://stackoverflow.com/questions/41453263/docker-networking-disabled-warning-ipv4-forwarding-is-disabled-networking-wil
type insIPV4Forwarding struct{}

func (ins *insIPV4Forwarding) Name() string {
	return "ipv4-forwarding"
}

func (ins *insIPV4Forwarding) Check() (*CheckResult, error) {
	content, err := ioutil.ReadFile("/etc/sysctl.conf")
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "net.ipv4.ip_forward") {
			return nil, nil
		}
	}
	return NeedsInstall, nil
}

func (ins *insIPV4Forwarding) Install() error {
	script := `
		echo net.ipv4.ip_forward=1 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
		sudo sysctl --system
  `
	return errors.Trace(run.Shell(script))
}

func (ins *insIPV4Forwarding) BashRC() string {
	return ""
}
