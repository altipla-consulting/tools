package deploy

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

type cmdFlags struct {
	Project        string
	Memory         string
	ServiceAccount string
	Sentry         string
	VolumeSecret   []string
	EnvSecret      []string
	Tag            string
	Name           string
}

var (
	flags cmdFlags
)

func init() {
	Cmd.PersistentFlags().StringVar(&flags.Project, "project", "", "Google Cloud project where the container will be stored. Defaults to the GOOGLE_PROJECT environment variable.")
	Cmd.PersistentFlags().StringVar(&flags.Memory, "memory", "256Mi", "Memory available inside the Cloud Run application.")
	Cmd.PersistentFlags().StringVar(&flags.ServiceAccount, "service-account", "", "Service account. Defaults to one with the name of the application.")
	Cmd.PersistentFlags().StringVar(&flags.Sentry, "sentry", "", "Sentry project to configure.")
	Cmd.PersistentFlags().StringSliceVar(&flags.VolumeSecret, "volume-secret", nil, "Secrets to mount as volumes.")
	Cmd.PersistentFlags().StringSliceVar(&flags.EnvSecret, "env-secret", nil, "Secrets to mount as environment variables.")
	Cmd.PersistentFlags().StringVar(&flags.Tag, "tag", "", "Name of the revision included in the URL. Defaults to the Gerrit change and patchset.")
	Cmd.PersistentFlags().StringVar(&flags.Name, "name", "", "Name of the application that will be deployed. Defaults to the folder name.")
}

var Cmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploy a container to Cloud Run.",
	Example: "wave deploy foo",
	Args:    cobra.ExactArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		app := args[0]

		if flags.Project == "" {
			flags.Project = os.Getenv("GOOGLE_PROJECT")
		}
		if flags.Sentry == "" {
			return errors.Errorf("--sentry flag is required")
		}
		if flags.ServiceAccount == "" {
			flags.ServiceAccount = app
		}

		client, err := sentry.NewClient(os.Getenv("SENTRY_AUTH_TOKEN"), nil, nil)
		if err != nil {
			return errors.Trace(err)
		}

		org := sentry.Organization{
			Slug: apiString("altipla-consulting"),
		}
		keys, err := client.GetClientKeys(org, sentry.Project{Slug: apiString(flags.Sentry)})
		if err != nil {
			return errors.Trace(err)
		}

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")
		if os.Getenv("BUILD_CAUSE") == "SCMTRIGGER" {
			version += ".preview"
			if flags.Tag == "" {
				flags.Tag = "preview-" + os.Getenv("GERRIT_CHANGE_NUMBER") + "-" + os.Getenv("GERRIT_PATCHSET_NUMBER")
			}
		}

		log.WithFields(log.Fields{
			"name":            app,
			"version":         version,
			"memory":          flags.Memory,
			"service-account": flags.ServiceAccount,
		}).Info("Deploy app")

		gcloud := []string{
			"beta", "run", "deploy",
			app,
			"--image", "eu.gcr.io/" + flags.Project + "/" + app + ":" + version,
			"--region", "europe-west1",
			"--platform", "managed",
			"--concurrency", "50",
			"--timeout", "60s",
			"--service-account", flags.ServiceAccount + "@" + flags.Project + ".iam.gserviceaccount.com",
			"--memory", flags.Memory,
			"--set-env-vars", "SENTRY_DSN=" + keys[0].DSN.Public,
			"--labels", "app=" + app,
		}
		if len(flags.VolumeSecret) > 0 {
			var secrets []string
			for _, secret := range flags.VolumeSecret {
				secrets = append(secrets, "/etc/secrets/"+secret+"="+secret+":latest")
			}
			for _, secret := range flags.EnvSecret {
				varname := strings.Replace(strings.ToUpper(secret), "-", "_", -1)
				secrets = append(secrets, varname+"="+secret+":latest")
			}
			gcloud = append(gcloud, "--set-secrets", strings.Join(secrets, ","))
		}
		if flags.Tag != "" {
			if os.Getenv("BUILD_CAUSE") == "SCMTRIGGER" {
				gcloud = append(gcloud, "--no-traffic")
				gcloud = append(gcloud, "--max-instances", "1")
			}
			gcloud = append(gcloud, "--tag", flags.Tag)
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

		if os.Getenv("BUILD_CAUSE") != "SCMTRIGGER" {
			log.WithFields(log.Fields{
				"name":    app,
				"version": version,
			}).Info("Enable traffic to latest version of the app")

			traffic := exec.Command(
				"gcloud",
				"run", "services", "update-traffic",
				app,
				"--project", flags.Project,
				"--region", "europe-west1",
				"--to-latest",
			)
			traffic.Stdout = os.Stdout
			traffic.Stderr = os.Stderr
			if err := traffic.Run(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}

func apiString(s string) *string {
	return &s
}
