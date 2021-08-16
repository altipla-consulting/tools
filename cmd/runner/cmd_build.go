package main

import (
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var flagBuildProject string

func init() {
	cmdRoot.AddCommand(cmdBuild)
	cmdBuild.PersistentFlags().StringVarP(&flagBuildProject, "project", "p", "", "Google Cloud project where the container will be stored.")
}

var cmdBuild = &cobra.Command{
	Use:   "build",
	Short: "Build a container from a predefined folder structure.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		if flagBuildProject == "" {
			return errors.Errorf("--project flag is required")
		}

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")

		for _, app := range args {
			logger := log.WithFields(log.Fields{
				"name":    app,
				"version": version,
			})

			logger.Info("Build app")
			build := exec.Command(
				"docker",
				"build",
				"--cache-from", "eu.gcr.io/"+flagBuildProject+"/"+app+":latest",
				"-f", app+"/Dockerfile",
				"-t", "eu.gcr.io/"+flagBuildProject+"/"+app+":latest",
				"-t", "eu.gcr.io/"+flagBuildProject+"/"+app+":"+version,
				".",
			)
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr
			if err := build.Run(); err != nil {
				return errors.Trace(err)
			}

			logger.Info("Push to Container Registry")
			push := exec.Command("docker", "push", "-t", "eu.gcr.io/"+flagBuildProject+"/"+app+":latest")
			push.Stdout = os.Stdout
			push.Stderr = os.Stderr
			if err := push.Run(); err != nil {
				return errors.Trace(err)
			}

			push = exec.Command("docker", "push", "-t", "eu.gcr.io/"+flagBuildProject+"/"+app+":"+version)
			push.Stdout = os.Stdout
			push.Stderr = os.Stderr
			if err := push.Run(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
