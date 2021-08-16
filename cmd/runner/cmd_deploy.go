package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/atlassian/go-sentry-api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var (
	flagDeployProject  string
	flagMemory         string
	flagServiceAccount string
	flagSentry         string
)

func init() {
	cmdRoot.AddCommand(cmdDeploy)
	cmdDeploy.PersistentFlags().StringVarP(&flagDeployProject, "project", "p", "", "Google Cloud project where the container will be stored.")
	cmdDeploy.PersistentFlags().StringVarP(&flagMemory, "memory", "m", "256Mi", "Memory available inside the Cloud Run application.")
	cmdDeploy.PersistentFlags().StringVarP(&flagServiceAccount, "service-account", "a", "", "Service account. Defaults to one with the name of the application.")
	cmdDeploy.PersistentFlags().StringVarP(&flagSentry, "sentry", "s", "", "Sentry project to configure.")
}

var cmdDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a container to Cloud Run.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		if flagDeployProject == "" {
			return errors.Errorf("--project flag is required")
		}
		if flagSentry == "" {
			return errors.Errorf("--sentry flag is required")
		}

		client, err := sentry.NewClient(os.Getenv("SENTRY_AUTH_TOKEN"), nil, nil)
		if err != nil {
			return errors.Trace(err)
		}

		org := sentry.Organization{
			Slug: apiString("altipla-consulting"),
		}
		keys, err := client.GetClientKeys(org, sentry.Project{Slug: apiString(flagSentry)})
		if err != nil {
			return errors.Trace(err)
		}

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")

		for _, app := range args {
			logger := log.WithFields(log.Fields{
				"name":    app,
				"version": version,
			})

			logger.Info("Deploy app")
			build := exec.Command(
				"gcloud",
				"run", "deploy",
				app,
				"--image", "eu.gcr.io/"+flagDeployProject+"/"+app+":"+version,
				"--region", "europe-west1",
				"--platform", "managed",
				"--concurrency", "50",
				"--timeout", "60s",
				"--service-account", flagServiceAccount+"@"+flagDeployProject+".iam.gserviceaccount.com",
				"--max-instances", "20",
				"--memory", flagMemory,
				"--set-env-vars", "SENTRY_DSN="+keys[0].DSN.Public,
				"--labels", "app="+app,
			)
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr
			if err := build.Run(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}

func apiString(s string) *string {
	return &s
}
