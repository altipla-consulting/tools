package run

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func Git(args ...string) error {
	log.Debugf("RUN: git %s", strings.Join(args, " "))

	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return errors.Trace(cmd.Run())
}

func GitCaptureOutput(args ...string) (string, error) {
	log.Debugf("RUN: git %s", strings.Join(args, " "))

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Trace(err)
	}

	return strings.TrimSpace(string(output)), nil
}
