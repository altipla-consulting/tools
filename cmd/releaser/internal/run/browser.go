package run

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func OpenBrowser(url string) error {
	if _, err := exec.LookPath("xdg-open"); err != nil {
		var ee *exec.Error
		if errors.As(err, &ee) && errors.Is(ee.Err, exec.ErrNotFound) {
			log.Println("Open the following URL in your browser:")
			log.Println("\t" + url)
			return nil
		}

		return errors.Trace(err)
	}
	return errors.Trace(exec.Command("xdg-open", url).Start())
}
