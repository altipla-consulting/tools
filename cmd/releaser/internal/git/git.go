package git

import (
	"log"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
	"libs.altipla.consulting/errors"
	"tools.altipla.consulting/cmd/releaser/internal/run"
)

func CurrentBranch() (string, error) {
	branch, err := run.Git("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", errors.Trace(err)
	}
	return branch, nil
}

func LatestRemoteTag() (string, error) {
	lines, err := run.Git("ls-remote", "-q", "--tags", "--refs")
	if err != nil {
		return "", errors.Trace(err)
	}

	var tags []string
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

	if len(tags) == 0 {
		return "", nil
	}

	sort.Sort(sort.Reverse(semver.ByVersion(tags)))
	return tags[0], nil
}

func PreviousTag() (string, error) {
	lines, err := run.Git("ls-remote", "-q", "--tags", "--refs")
	if err != nil {
		return "", errors.Trace(err)
	}

	var tags []string
	for _, line := range strings.Split(lines, "\n") {
		if line == "" {
			continue
		}

		tag := strings.Fields(line)[1]
		if strings.HasPrefix(tag, "refs/tags/") {
			tag = tag[len("refs/tags/"):]
			tags = append(tags, tag)
		}
	}

	if len(tags) < 2 {
		return "", nil
	}

	log.Println(tags)

	sort.Sort(sort.Reverse(semver.ByVersion(tags)))
	log.Println(tags)
	return tags[1], nil
}

func DirtyWorkingTree() (bool, error) {
	status, err := run.Git("status", "-s")
	if err != nil {
		return false, errors.Trace(err)
	}
	return len(status) > 0, nil
}

func RemoteHistoryClean() (bool, error) {
	history, err := run.Git("rev-list", "--count", "--left-only", "@{u}...HEAD")
	if err != nil {
		return false, errors.Trace(err)
	}
	return history == "" || history == "0", nil
}

func Tag(tag string) error {
	_, err := run.Git("tag", tag)
	return errors.Trace(err)
}

func RemoteURL(name string) (string, error) {
	remote, err := run.Git("remote", "get-url", name)
	if err != nil {
		return "", errors.Trace(err)
	}
	return remote, nil
}

func FirstCommit() (string, error) {
	commit, err := run.Git("rev-list", "--max-parents=0", "HEAD")
	if err != nil {
		return "", errors.Trace(err)
	}
	return commit, nil
}

func CommitLogFrom(from string) (string, error) {
	commitlog, err := run.Git("log", "--format=%s %h", from+"..HEAD")
	if err != nil {
		return "", errors.Trace(err)
	}
	return commitlog, nil
}

func RemoteTagExists(name string) (bool, error) {
	_, err := run.Git("rev-parse", "--quiet", "--verify", "refs/tags/"+name)
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 {
			return false, nil
		}

		return false, errors.Trace(err)
	}

	return true, nil
}
