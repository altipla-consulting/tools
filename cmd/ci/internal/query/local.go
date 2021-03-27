package query

import (
	"os/exec"
	"strings"

	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/ci/internal/run"
)

var org, repo string

func IsGerrit() (bool, error) {
	remote, err := run.GitCaptureOutput("remote", "get-url", "origin")
	if err != nil {
		return false, errors.Trace(err)
	}
	return strings.Contains(remote, "gerrit.altipla.consulting"), nil
}

func MainBranch() (string, error) {
	branch, err := run.GitCaptureOutput("branch", "-a")
	if err != nil {
		return "", errors.Trace(err)
	}
	mainBranch := "master"
	if strings.Contains(branch, "remotes/origin/main") {
		mainBranch = "main"
	}
	return mainBranch, nil
}

func extractGitHub() error {
	if org != "" {
		return nil
	}

	remote, err := run.GitCaptureOutput("remote", "get-url", "origin")
	if err != nil {
		return errors.Trace(err)
	}

	parts := strings.Split(remote, "/")
	org = parts[0][len("git@github.com:"):]
	repo = parts[1][:len(parts[1])-len(".git")]
	return nil
}

func CurrentOrg() (string, error) {
	if err := extractGitHub(); err != nil {
		return "", errors.Trace(err)
	}
	return org, nil
}

func CurrentRepo() (string, error) {
	if err := extractGitHub(); err != nil {
		return "", errors.Trace(err)
	}
	return repo, nil
}

func CurrentBranch() (string, error) {
	branch, err := run.GitCaptureOutput("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", errors.Trace(err)
	}
	return branch, nil
}

func LastCommitMessage() (string, error) {
	msg, err := run.GitCaptureOutput("log", "-1", "--pretty=%B")
	if err != nil {
		return "", errors.Trace(err)
	}
	return msg, nil
}

func BranchExists(name string) (bool, error) {
	if err := run.Git("show-ref", "--verify", "--quiet", "refs/heads/"+name); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			if exit.ProcessState.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, errors.Trace(err)
	}
	return true, nil
}
