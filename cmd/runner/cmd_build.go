package main

import (
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

type buildFlags struct {
	Project string
}

var (
	flagBuild buildFlags
)

func init() {
	cmdRoot.AddCommand(cmdBuild)
	cmdBuild.PersistentFlags().StringVar(&flagBuild.Project, "project", "", "Google Cloud project where the container will be stored. Defaults to the GOOGLE_PROJECT environment variable.")
}

var cmdBuild = &cobra.Command{
	Use:   "build",
	Short: "Build a container from a predefined folder structure.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		if flagBuild.Project == "" {
			flagBuild.Project = os.Getenv("GOOGLE_PROJECT")
		}

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")
		if os.Getenv("BUILD_CAUSE") == "SCMTRIGGER" {
			version += ".preview"
		}

		for _, app := range args {
			logger := log.WithFields(log.Fields{
				"name":    app,
				"version": version,
			})

			logger.Info("Build app")
			build := exec.Command(
				"docker",
				"build",
				"--cache-from", "eu.gcr.io/"+flagBuild.Project+"/"+app+":latest",
				"-f", app+"/Dockerfile",
				"-t", "eu.gcr.io/"+flagBuild.Project+"/"+app+":latest",
				"-t", "eu.gcr.io/"+flagBuild.Project+"/"+app+":"+version,
				".",
			)
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr
			if err := build.Run(); err != nil {
				return errors.Trace(err)
			}

			logger.Info("Push to Container Registry")
			push := exec.Command("docker", "push", "eu.gcr.io/"+flagBuild.Project+"/"+app+":latest")
			push.Stdout = os.Stdout
			push.Stderr = os.Stderr
			if err := push.Run(); err != nil {
				return errors.Trace(err)
			}

			push = exec.Command("docker", "push", "eu.gcr.io/"+flagBuild.Project+"/"+app+":"+version)
			push.Stdout = os.Stdout
			push.Stderr = os.Stderr
			if err := push.Run(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
