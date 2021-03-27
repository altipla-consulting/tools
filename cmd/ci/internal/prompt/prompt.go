package prompt

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"libs.altipla.consulting/errors"
)

func Confirm(msg string) (bool, error) {
	var reply bool
	prompt := &survey.Confirm{
		Message: msg,
	}
	if err := survey.AskOne(prompt, &reply); err != nil {
		return false, errors.Trace(err)
	}
	return reply, nil
}

func TextDefault(msg, def string) (string, error) {
	var text string
	prompt := &survey.Input{
		Message: msg,
		Default: def,
	}
	if err := survey.AskOne(prompt, &text); err != nil {
		return "", errors.Trace(err)
	}
	return strings.TrimSpace(text), nil
}
