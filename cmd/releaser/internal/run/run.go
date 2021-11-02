package run

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func Git(args ...string) (string, error) {
	return GitCapture(nil, args...)
}

func GitCapture(w io.Writer, args ...string) (string, error) {
	var buf bytes.Buffer
	var out io.Writer = &buf
	if w != nil {
		out = io.MultiWriter(&buf, w)
	}

	log.Debugf("RUN: git %s", strings.Join(args, " "))
	cmd := exec.Command("git", args...)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return "", errors.Trace(err)
	}

	log.Debugf("OUTPUT:\n%s", buf.String())

	return strings.TrimSpace(buf.String()), nil
}
