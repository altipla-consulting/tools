package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/modfile"
	"libs.altipla.consulting/errors"
	"tools.altipla.consulting/cmd/releaser/internal/git"
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
						Message: "Check main branch",
						Handler: releaseCheckMainBranch,
					},
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

	_ = release

	return nil
}

func releaseCheckMainBranch() error {
	branch, err := git.CurrentBranch()
	if err != nil {
		return errors.Trace(err)
	}

	if branch != "main" && branch != "master" {
		return errors.Errorf("Branch master or main expected to make a new release. Current branch: " + branch)
	}

	return nil
}
