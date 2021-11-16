package run

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"libs.altipla.consulting/errors"
)

func Interactive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return errors.Trace(cmd.Run())
}

func Shell(script string, vars ...map[string]string) error {
	script = "#!/bin/bash\nset -eu\n" + script

	for _, replace := range vars {
		for k, v := range replace {
			script = strings.Replace(script, "$"+k, v, -1)
		}
	}

	f, err := ioutil.TempFile("", "cdm")
	if err != nil {
		return errors.Trace(err)
	}
	defer os.Remove(f.Name())
	fmt.Fprint(f, script)
	if err := f.Close(); err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(Interactive("bash", f.Name()))
}

func InteractiveCaptureOutput(name string, args ...string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", errors.Trace(err)
	}

	return strings.TrimSpace(buf.String()), nil
}
