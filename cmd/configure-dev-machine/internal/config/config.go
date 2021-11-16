package config

import (
	"os"
	"path/filepath"

	"libs.altipla.consulting/errors"
)

func Filename(filename string) (string, error) {
	configdir, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Trace(err)
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return "", errors.Trace(err)
	}

	return filepath.Join(configdir, "configure-dev-machine", filename), nil
}
