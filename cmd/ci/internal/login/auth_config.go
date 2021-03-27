package login

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"libs.altipla.consulting/errors"
)

type AuthConfig struct {
	AccessToken string
	Username    string
}

func ReadAuthConfig() (*AuthConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Trace(err)
	}
	content, err := ioutil.ReadFile(filepath.Join(home, ".ci", "auth.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Trace(err)
	}

	auth := new(AuthConfig)
	if err := json.Unmarshal(content, auth); err != nil {
		return nil, errors.Trace(err)
	}

	return auth, nil
}

func (auth *AuthConfig) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.Trace(err)
	}
	if err := os.MkdirAll(filepath.Join(home, ".ci"), 0700); err != nil {
		return errors.Trace(err)
	}

	f, err := os.Create(filepath.Join(home, ".ci", "auth.json"))
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(auth); err != nil {
		return errors.Trace(err)
	}

	return nil
}
