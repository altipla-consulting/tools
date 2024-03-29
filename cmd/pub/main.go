package main

import (
	"encoding/json"
	"fmt"
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

	gerrit := readGerritInfo()
	nextVersion := filepath.Base(gerrit.RefName)

	pkg, err := readPackageJSON()
	if err != nil {
		return errors.Trace(err)
	}
	log.WithFields(log.Fields{
		"current": "v" + pkg.Version,
		"next":    nextVersion,
		"package": pkg.Name,
	}).Info("Release new version for NPM package")

	log.Info("Increment package.json version")
	if err := runCommand("npm", "version", nextVersion); err != nil {
		return errors.Trace(err)
	}

	log.Info("Publish package to NPM")
	if err := runCommand("npm", "publish", "--access", "public"); err != nil {
		return errors.Trace(err)
	}

	log.Info("Push commit updating version to Gerrit")
	if err := runCommand("scp", "-p", "-P", gerrit.Port, fmt.Sprintf("%s@%s:hooks/commit-msg", gerrit.BotUsername, gerrit.Host), ".git/hooks/"); err != nil {
		return errors.Trace(err)
	}
	if err := runCommand("git", "add", "package.json", "package-lock.json"); err != nil {
		return errors.Trace(err)
	}
	if err := runCommand("git", "commit", "-m", "Release "+nextVersion); err != nil {
		return errors.Trace(err)
	}
	if err := runCommand("ci", "push"); err != nil {
		return errors.Trace(err)
	}

	log.Info("Assign reviewers to the commit")
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Trace(err)
	}
	commit := strings.TrimSpace(string(output))
	if err := runCommand("ssh", fmt.Sprintf("%s@%s", gerrit.BotUsername, gerrit.Host), "-p", gerrit.Port, "gerrit", "set-reviewers", commit, "-a", gerrit.ReviewersGroup); err != nil {
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

type gerritInfo struct {
	BotUsername    string
	RefName        string
	Host           string
	Port           string
	ReviewersGroup string
}

func readGerritInfo() gerritInfo {
	return gerritInfo{
		BotUsername:    os.Getenv("GERRIT_BOT_USERNAME"),
		RefName:        os.Getenv("GERRIT_REFNAME"),
		Host:           os.Getenv("GERRIT_HOST"),
		Port:           os.Getenv("GERRIT_PORT"),
		ReviewersGroup: os.Getenv("GERRIT_REVIEWERS_GROUP"),
	}
}
