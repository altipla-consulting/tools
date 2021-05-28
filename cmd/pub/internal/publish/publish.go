package publish

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/pub/internal/run"
)

type tracker struct {
	step, total int

	NextVersion string
	PackageJSON PackageJSONSpec
}

func (tr *tracker) Announce(msg string) {
	tr.step++
	log.WithField("step", fmt.Sprintf("%d/%d", tr.step, tr.total)).Info(msg)
}

func Run() error {
	nextVersion := filepath.Base(os.Getenv("GERRIT_REFNAME"))

	actions := []func(tr *tracker) error{
		installDependencies,
		runTests,
		incrementVersion,
		publishPackage,
	}

	pkg, err := readPackageJSON()
	if err != nil {
		return errors.Trace(err)
	}
	log.WithFields(log.Fields{
		"current": "v" + pkg.Version,
		"next":    nextVersion,
		"package": pkg.Name,
	}).Info("Release new version for NPM package")

	tr := &tracker{
		total:       len(actions),
		NextVersion: nextVersion,
		PackageJSON: pkg,
	}
	for _, action := range actions {
		if err := action(tr); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

type PackageJSONSpec struct {
	Name    string
	Version string
	Scripts struct {
		Test string
	}
}

func readPackageJSON() (PackageJSONSpec, error) {
	f, err := os.Open("package.json")
	if err != nil {
		return PackageJSONSpec{}, errors.Trace(err)
	}
	defer f.Close()

	var pkg PackageJSONSpec
	if err := json.NewDecoder(f).Decode(&pkg); err != nil {
		return PackageJSONSpec{}, errors.Trace(err)
	}

	return pkg, nil
}

func installDependencies(tr *tracker) error {
	tr.Announce("Install NPM dependencies from scratch")
	return errors.Trace(run.CheckSuccess("npm", "ci", "--engine-strict"))
}

func runTests(tr *tracker) error {
	tr.Announce("Run tests")

	if tr.PackageJSON.Scripts.Test == "" {
		log.Warning("There are no tests defined in the `npm test` script. Skipping step")
		return nil
	}

	return errors.Trace(run.CheckSuccess("npm", "test"))
}

func incrementVersion(tr *tracker) error {
	tr.Announce("Increment package.json version")

	content, err := ioutil.ReadFile(".npmrc")
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Trace(err)
		}

		defer os.Remove(".npmrc")
	} else {
		defer ioutil.WriteFile(".npmrc", content, 0600)
	}

	result := append(content, []byte("\ngit-tag-version=false\n")...)
	if err := ioutil.WriteFile(".npmrc", result, 0600); err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(run.CheckSuccess("npm", "version", tr.NextVersion))
}

func publishPackage(tr *tracker) error {
	tr.Announce("Publish package to NPM")

	content, err := ioutil.ReadFile(".npmrc")
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Trace(err)
		}

		defer os.Remove(".npmrc")
	} else {
		defer ioutil.WriteFile(".npmrc", content, 0600)
	}

	result := append(content, []byte("\nregistry=https://registry.npmjs.org/\n//registry.npmjs.org/:_authToken="+os.Getenv("NPM_TOKEN")+"\n")...)
	if err := ioutil.WriteFile(".npmrc", result, 0600); err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(run.CheckSuccess("npm", "publish", "--access", "public"))
}
