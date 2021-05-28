package run

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func CheckSuccess(command ...string) error {
	log.WithField("command", strings.Join(command, " ")).Debug("Running command checking output code")
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("COMMAND OUTPUT:\n", string(output))
		return errors.Trace(err)
	}
	log.Debug("COMMAND OUTPUT:\n", string(output))
	return nil
}

func CheckFails(command ...string) error {
	log.WithField("command", strings.Join(command, " ")).Debug("Running command checking output code expecting failure")
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			log.WithField("exit-code", exit.ExitCode()).Debug("Command failed as expected")
			return nil
		}
		log.Error("COMMAND OUTPUT:\n", string(output))
		return errors.Trace(err)
	}
	log.Debug("COMMAND OUTPUT:\n", string(output))
	return errors.Errorf("command did not fail")
}

func Output(command ...string) (string, error) {
	log.WithField("command", strings.Join(command, " ")).Debug("Running command to obtain output")
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("COMMAND OUTPUT:\n", string(output))
		return "", errors.Trace(err)
	}
	log.Debug("COMMAND OUTPUT:\n", string(output))
	return strings.TrimSpace(string(output)), nil
}
