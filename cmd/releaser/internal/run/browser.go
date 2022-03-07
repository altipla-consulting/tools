package run

import (
	"os/exec"

	"libs.altipla.consulting/errors"
)

var ErrCannotOpenBrowser = errors.New("run: cannot open browser")

func OpenBrowser(url string) error {
	if _, err := exec.LookPath("xdg-open"); err != nil {
		var ee *exec.Error
		if errors.As(err, &ee) && errors.Is(ee.Err, exec.ErrNotFound) {
			return errors.Trace(ErrCannotOpenBrowser)
		}

		return errors.Trace(err)
	}
	return errors.Trace(exec.Command("xdg-open", url).Start())
}
