package run

import (
	"os/exec"

	"libs.altipla.consulting/errors"
)

func OpenBrowser(url string) error {
	return errors.Trace(exec.Command("xdg-open", url).Start())
}
