package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"libs.altipla.consulting/errors"
)

func main() {
	if err := run(); err != nil {
		log.Error(err.Error())
		log.Debug(errors.Stack(err))
		os.Exit(1)
	}
}

func run() error {
	var flagDebug bool
	flag.BoolVarP(&flagDebug, "debug", "d", false, "Debug logging")
	flag.Parse()

	log.SetFormatter(new(log.TextFormatter))
	if flagDebug {
		log.SetLevel(log.DebugLevel)
		log.Debug("DEBUG log level activated")
	}

	nextVersion := filepath.Base(os.Getenv("GERRIT_REFNAME"))

	pkg, err := readPackageJSON()
	if err != nil {
		return errors.Trace(err)
	}
	log.WithFields(log.Fields{
		"current": "v" + pkg.Version,
		"next":    nextVersion,
		"package": pkg.Name,
	}).Info("Release new version for NPM package")

	log.Info("Install NPM dependencies from scratch")
	if err := runCommand("npm", "ci", "--engine-strict"); err != nil {
		return errors.Trace(err)
	}

	log.Info("Run linter")
	if pkg.Scripts.Lint == "" {
		log.Warning("There is no linter defined in the `npm run lint` script. Skipping step")
	} else {
		if err := runCommand("npm", "run", "lint"); err != nil {
			return errors.Trace(err)
		}
	}

	log.Info("Run tests")
	if pkg.Scripts.Test == "" {
		log.Warning("There are no tests defined in the `npm test` script. Skipping step")
	} else {
		if err := runCommand("npm", "test"); err != nil {
			return errors.Trace(err)
		}
	}

	log.Info("Configure NPM to release a new package")
	content, err := ioutil.ReadFile(".npmrc")
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Trace(err)
		}
		defer os.Remove(".npmrc")
	} else {
		defer ioutil.WriteFile(".npmrc", content, 0600)
	}
	newlines := []string{
		"",
		"git-tag-version=false",
		"registry=https://registry.npmjs.org/",
		"//registry.npmjs.org/:_authToken=" + os.Getenv("NPM_TOKEN"),
		"",
	}
	result := append(content, []byte(strings.Join(newlines, "\n"))...)
	if err := ioutil.WriteFile(".npmrc", result, 0600); err != nil {
		return errors.Trace(err)
	}

	log.Info("Increment package.json version")
	if err := runCommand("npm", "version", nextVersion, "-m", "Release "+nextVersion); err != nil {
		return errors.Trace(err)
	}

	log.Info("Publish package to NPM")
	if err := runCommand("npm", "publish", "--access", "public"); err != nil {
		return errors.Trace(err)
	}

	log.Info("Push commit updating version to Gerrit")
	if err := runCommand("ci", "push"); err != nil {
		return errors.Trace(err)
	}

	return nil
}

type PackageJSONSpec struct {
	Name    string
	Version string
	Scripts struct {
		Lint string
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

func runCommand(command ...string) error {
	log.WithField("command", strings.Join(command, " ")).Debug("Running command")

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return errors.Trace(cmd.Run())
}
