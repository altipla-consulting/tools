package version

import (
	"bufio"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ebikt/go-debian/control"
	"libs.altipla.consulting/errors"
)

type Info struct {
	CurrentlyInstalled string
	Latest             string
}

func FetchInfo() (*Info, error) {
	currentlyInstalled, err := getCurrentlyInstalled()
	if err != nil {
		return nil, errors.Trace(err)
	}
	latest, err := getLatest()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Info{
		CurrentlyInstalled: currentlyInstalled,
		Latest:             latest,
	}, nil
}

func getCurrentlyInstalled() (string, error) {
	content, err := ioutil.ReadFile("/etc/configure-dev-machine/installed-version")
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", errors.Trace(err)
	}
	return string(content), nil
}

func getLatest() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://europe-west1-apt.pkg.dev/projects/altipla-tools/dists/acpublic/main/binary-amd64/Packages", nil)
	if err != nil {
		return "", errors.Trace(err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer resp.Body.Close()

	indexes, err := control.ParseBinaryIndex(bufio.NewReader(resp.Body))
	if err != nil {
		return "", errors.Trace(err)
	}
	for _, index := range indexes {
		if index.Package == "configure-dev-machine" {
			return index.Version.String(), nil
		}
	}

	return "~~unknown", nil
}
