package config

import (
	"io/ioutil"
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

func Version() (string, error) {
	filename, err := Filename("version.txt")
	if err != nil {
		return "", errors.Trace(err)
	}
	version, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return "dev", nil
		}
		return "", errors.Trace(err)
	}
	return string(version), nil
}
