package main

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/atlassian/go-sentry-api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

type deployFlags struct {
	Project        string
	Memory         string
	ServiceAccount string
	Sentry         string
	VolumeSecret   []string
	Tag            string
}

var (
	flagDeploy deployFlags
)

func init() {
	cmdRoot.AddCommand(cmdDeploy)
	cmdDeploy.PersistentFlags().StringVar(&flagDeploy.Project, "project", "", "Google Cloud project where the container will be stored. Defaults to the GOOGLE_PROJECT environment variable.")
	cmdDeploy.PersistentFlags().StringVar(&flagDeploy.Memory, "memory", "256Mi", "Memory available inside the Cloud Run application.")
	cmdDeploy.PersistentFlags().StringVar(&flagDeploy.ServiceAccount, "service-account", "", "Service account. Defaults to one with the name of the application.")
	cmdDeploy.PersistentFlags().StringVar(&flagDeploy.Sentry, "sentry", "", "Sentry project to configure.")
	cmdDeploy.PersistentFlags().StringSliceVar(&flagDeploy.VolumeSecret, "volume-secret", nil, "Secrets to mount as volumes.")
	cmdDeploy.PersistentFlags().StringVar(&flagDeploy.Tag, "tag", "", "Name of the revision included in the URL. Defaults to the Gerrit change and patchset.")
}

var cmdDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a container to Cloud Run.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		if flagDeploy.Project == "" {
			flagDeploy.Project = os.Getenv("GOOGLE_PROJECT")
		}
		if flagDeploy.Sentry == "" {
			return errors.Errorf("--sentry flag is required")
		}

		client, err := sentry.NewClient(os.Getenv("SENTRY_AUTH_TOKEN"), nil, nil)
		if err != nil {
			return errors.Trace(err)
		}

		org := sentry.Organization{
			Slug: apiString("altipla-consulting"),
		}
		keys, err := client.GetClientKeys(org, sentry.Project{Slug: apiString(flagDeploy.Sentry)})
		if err != nil {
			return errors.Trace(err)
		}

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")
		if os.Getenv("BUILD_CAUSE") == "SCMTRIGGER" {
			version += ".preview"
			if flagDeploy.Tag == "" {
				flagDeploy.Tag = "preview-" + os.Getenv("GERRIT_CHANGE_NUMBER") + "-" + os.Getenv("GERRIT_PATCHSET_NUMBER")
			}
		}

		for _, app := range args {
			serviceAccount := flagDeploy.ServiceAccount
			if serviceAccount == "" {
				serviceAccount = app
			}

			log.WithFields(log.Fields{
				"name":            app,
				"version":         version,
				"memory":          flagDeploy.Memory,
				"service-account": serviceAccount,
			}).Info("Deploy app")

			gcloud := []string{
				"beta", "run", "deploy",
				app,
				"--image", "eu.gcr.io/" + flagDeploy.Project + "/" + app + ":" + version,
				"--region", "europe-west1",
				"--platform", "managed",
				"--concurrency", "50",
				"--timeout", "60s",
				"--service-account", serviceAccount + "@" + flagDeploy.Project + ".iam.gserviceaccount.com",
				"--memory", flagDeploy.Memory,
				"--set-env-vars", "SENTRY_DSN=" + keys[0].DSN.Public,
				"--labels", "app=" + app,
			}
			if len(flagDeploy.VolumeSecret) > 0 {
				var secrets []string
				for _, secret := range flagDeploy.VolumeSecret {
					secrets = append(secrets, "/etc/secrets/"+secret+"="+secret+":latest")
				}
				gcloud = append(gcloud, "--set-secrets", strings.Join(secrets, ","))
			}
			if flagDeploy.Tag != "" {
				gcloud = append(gcloud, "--no-traffic")
				gcloud = append(gcloud, "--max-instances", "1")
				gcloud = append(gcloud, "--tag", flagDeploy.Tag)
			} else {
				gcloud = append(gcloud, "--max-instances", "20")
			}

			log.Debug(strings.Join(append([]string{"gcloud"}, gcloud...), " "))

			build := exec.Command("gcloud", gcloud...)
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
