package run

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func OpenBrowser(url string) error {
	if _, err := exec.LookPath("xdg-open"); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Println("Open the following URL in your browser:")
			log.Println("\t" + url)
			return nil
		}

		return errors.Trace(err)
	}
	return errors.Trace(exec.Command("xdg-open", url).Start())
}
