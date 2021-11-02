package git

import (
	"sort"
	"strings"

	"golang.org/x/mod/semver"
	"libs.altipla.consulting/errors"
	"tools.altipla.consulting/cmd/releaser/internal/run"
)

func CurrentBranch() (string, error) {
	branch, err := run.GitCaptureOutput("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", errors.Trace(err)
	}
	return branch, nil
}

func LatestRemoteTag() (string, error) {
	lines, err := run.GitCaptureOutput("ls-remote", "-q", "--tags", "--refs")
	if err != nil {
		return "", errors.Trace(err)
	}

	tags := []string{"v0.0.0"}
	for _, line := range strings.Split(lines, "\n") {
		if line == "" {
			continue
		}

		tag := strings.Fields(line)[1]
		if strings.HasPrefix(tag, "refs/tags/") {
			tag = tag[len("refs/tags/"):]
		}
		tags = append(tags, tag)
	}

	sort.Sort(sort.Reverse(semver.ByVersion(tags)))
	return tags[0], nil
}
