package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/modfile"
	"libs.altipla.consulting/errors"
	"tools.altipla.consulting/cmd/releaser/internal/git"
	"tools.altipla.consulting/cmd/releaser/internal/run"
	"tools.altipla.consulting/cmd/releaser/internal/tasks"
)

func Release(update string) error {
	gomod, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return errors.Trace(err)
	}
	modname := modfile.ModulePath(gomod)

	current, err := git.LatestRemoteTag()
	if err != nil {
		return errors.Trace(err)
	}
	if current == "" {
		current = "v0.0.0"
	}
	components := strings.Split(current, ".")
	major, err := strconv.ParseInt(components[0][1:], 10, 64)
	if err != nil {
		return errors.Trace(err)
	}
	minor, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return errors.Trace(err)
	}
	patch, err := strconv.ParseInt(components[2], 10, 64)
	if err != nil {
		return errors.Trace(err)
	}
	var release string
	switch update {
	case "patch":
		release = fmt.Sprintf("v%d.%d.%d", major, minor, patch+1)
	case "minor":
		release = fmt.Sprintf("v%d.%d.0", major, minor+1)
	case "major":
		release = fmt.Sprintf("v%d.0.0", major+1)
	default:
		panic("should not reach here")
	}

	fmt.Println()
	fmt.Println("Publish a new version of", aurora.Magenta(modname).Bold(), aurora.Gray(22, "(current: "+current+")").Faint())
	fmt.Println()
	fmt.Println()

	runner := &tasks.Runner{
		Tasks: []tasks.Runnable{
			&tasks.ParentTask{
				Message: "Prerequisites check",
				Tasks: []*tasks.Task{
					{
						Message: "Check git remote",
						Handler: func(w io.Writer) error {
							_, err := run.Git("ls-remote", "origin", "HEAD")
							return errors.Trace(err)
						},
					},
					{
						Message: "Check main branch",
						Handler: func(w io.Writer) error {
							branch, err := git.CurrentBranch()
							if err != nil {
								return errors.Trace(err)
							}
							if branch != "main" && branch != "master" {
								return errors.Errorf("Branch master or main expected to make a new release. Current branch: " + branch)
							}
							return nil
						},
					},
					{
						Message: "Check local working tree",
						Handler: func(w io.Writer) error {
							dirty, err := git.DirtyWorkingTree()
							if err != nil {
								return errors.Trace(err)
							}
							if dirty {
								return errors.Errorf("Unclean working tree. Commit or stash changes first.")
							}
							return nil
						},
					},
					{
						Message: "Check remote history",
						Handler: func(w io.Writer) error {
							clean, err := git.RemoteHistoryClean()
							if err != nil {
								return errors.Trace(err)
							}
							if !clean {
								return errors.Errorf("Remote history differs. Please pull changes.")
							}
							return nil
						},
					},
				},
			},
			&tasks.ParentTask{
				Message: "Release new version",
				Tasks: []*tasks.Task{
					{
						Message: "Commit new tag",
						Handler: func(w io.Writer) error {
							_, err := run.GitCapture(w, "commit", "--allow-empty", "-m", release[1:])
							return errors.Trace(err)
						},
					},
					{
						Message: "Tag repo",
						Handler: func(w io.Writer) error {
							return errors.Trace(git.Tag(release))
						},
					},
					{
						Message: "Push tags",
						Handler: func(w io.Writer) error {
							_, err := run.GitCapture(w, "push", "--follow-tags")
							return errors.Trace(err)
						},
					},
				},
			},
			&tasks.Task{
				Disable: true,
				Message: "Create release draft on GitHub",
				Handler: func(w io.Writer) error {
					remote, err := git.RemoteURL("origin")
					if err != nil {
						return errors.Trace(err)
					}
					repo := "https://github.com/" + strings.TrimSuffix(strings.Split(remote, ":")[1], ".git")
					u, err := url.Parse(repo + "/releases/new")
					if err != nil {
						return errors.Trace(err)
					}
					qs := make(url.Values)
					qs.Set("tag", release)
					notes, err := releaseNotes(repo, release)
					if err != nil {
						return errors.Trace(err)
					}
					qs.Set("body", notes)
					u.RawQuery = qs.Encode()
					if err := run.OpenBrowser(u.String()); err != nil {
						return errors.Trace(err)
					}

					return nil
				},
			},
		},
	}

	if err := runner.Run(); err != nil {
		fmt.Println()
		fmt.Println(aurora.Red("✖ "), err.Error())
		log.Debug(errors.Stack(err) + "\n")
		return nil
	}

	return nil
}

func releaseNotes(repo, release string) (string, error) {
	prev, err := git.PreviousTag()
	if err != nil {
		return "", errors.Trace(err)
	}
	if prev == "" {
		prev, err = git.FirstCommit()
		if err != nil {
			return "", errors.Trace(err)
		}
	}

	commitlog, err := git.CommitLogFrom(prev)
	if err != nil {
		return "", errors.Trace(err)
	}
	if commitlog == "" {
		return "", nil
	}

	var body []string
	for _, line := range strings.Split(commitlog, "\n") {
		index := strings.LastIndex(line, " ")
		body = append(body, "- "+line[:index]+"  "+line[index+1:])
	}

	body = append(body, "", "")
	body = append(body, repo+"/compare/"+prev+"..."+release)

	return strings.Join(body, "\n"), nil
}
