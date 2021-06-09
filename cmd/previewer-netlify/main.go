package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	var flagDebug, flagProduction bool
	flag.BoolVarP(&flagDebug, "debug", "d", false, "Debug logging")
	flag.BoolVarP(&flagProduction, "production", "p", false, "Deploy to production instead of a preview")
	flag.Parse()

	log.SetFormatter(new(log.TextFormatter))
	if flagDebug {
		log.SetLevel(log.DebugLevel)
		log.Debug("DEBUG log level activated")
	}

	gerrit := readGerritInfo()

	log.Info("Get last commit message")
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Trace(err)
	}
	var filtered []string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "Change-Id") {
			continue
		}
		filtered = append(filtered, line)
	}
	lastCommit := strings.TrimSpace(strings.Join(filtered, "\n"))

	log.Info("Deploy to Netlify")
	args := []string{
		"netlify",
		"deploy",
		"--dir", "dist",
		"--json",
		"--message", lastCommit,
	}
	if flagProduction {
		args = append(args, "--prod")
	}
	cmd = exec.Command(args[0], args[1:]...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Trace(err)
	}
	var deployment map[string]interface{}
	if err := json.Unmarshal(output, &deployment); err != nil {
		log.Warningf("Cannot parse Netlify output:\n%s", output)
		return errors.Trace(err)
	}
	deployURL := deployment["deploy_url"].(string)

	if !flagProduction {
		log.Info("Send preview URL as a Gerrit comment")
		args = []string{
			"ssh",
			"-p", gerrit.Port,
			fmt.Sprintf("%s@%s", gerrit.BotUsername, gerrit.Host),
			"gerrit", "review", fmt.Sprintf("%v,%v", gerrit.ChangeNumber, gerrit.PatchSetNumber),
			"--message", `"Previsualización desplegada en Netlify: ` + deployURL + `"`,
		}
		cmd = exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

type gerritInfo struct {
	BotUsername    string
	Host           string
	Port           string
	ChangeNumber   string
	PatchSetNumber string
}

func readGerritInfo() gerritInfo {
	return gerritInfo{
		BotUsername:    os.Getenv("GERRIT_BOT_USERNAME"),
		Host:           os.Getenv("GERRIT_HOST"),
		Port:           os.Getenv("GERRIT_PORT"),
		ChangeNumber:   os.Getenv("GERRIT_CHANGE_NUMBER"),
		PatchSetNumber: os.Getenv("GERRIT_PATCHSET_NUMBER"),
	}
}
