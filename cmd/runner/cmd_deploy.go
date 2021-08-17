package main

import (
	"fmt"
	"net/url"
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
	cmdDeploy.PersistentFlags().StringVar(&flagDeploy.Project, "project", "", "Google Cloud project where the container will be stored.")
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
			return errors.Errorf("--project flag is required")
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

			args := []string{
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
				args = append(args, "--set-secrets", strings.Join(secrets, ","))
			}
			if flagDeploy.Tag != "" {
				args = append(args, "--no-traffic")
				args = append(args, "--max-instances", "1")
				args = append(args, "--tag", flagDeploy.Tag)
			} else {
				args = append(args, "--max-instances", "20")
			}

			log.Debug(strings.Join(append([]string{"gcloud"}, args...), " "))

			build := exec.Command("gcloud", args...)
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr
			if err := build.Run(); err != nil {
				return errors.Trace(err)
			}
		}

		if os.Getenv("BUILD_CAUSE") == "SCMTRIGGER" {
			suffixcmd := exec.Command("gcloud", "run", "services", "describe", args[0], "--format", "value(status.url)")
			output, err := suffixcmd.CombinedOutput()
			if err != nil {
				log.Error(string(output))
				return errors.Trace(err)
			}
			u, err := url.Parse(strings.TrimSpace(string(output)))
			if err != nil {
				return errors.Trace(err)
			}
			parts := strings.Split(strings.Split(u.Host, ".")[0], "-")
			suffix := parts[len(parts)-2]
			var previews []string
			for _, app := range args {
				previews = append(previews, "https://"+flagDeploy.Tag+"---"+app+"-"+suffix+"-ew.a.run.app/")
			}
			log.WithField("previews", previews).Debug("Send comment to Gerrit with the previews")

			log.Info("Send preview URLs as a Gerrit comment")
			gerrit := readGerritInfo()
			args = []string{
				"ssh",
				"-p", gerrit.Port,
				fmt.Sprintf("%s@%s", gerrit.BotUsername, gerrit.Host),
				"gerrit", "review", fmt.Sprintf("%v,%v", gerrit.ChangeNumber, gerrit.PatchSetNumber),
				"--message", `"Previews deployed at:` + "\n" + strings.Join(previews, "\n") + `"`,
			}
			log.Debug(strings.Join(args, " "))
			comment := exec.Command(args[0], args[1:]...)
			comment.Stdout = os.Stdout
			comment.Stderr = os.Stderr
			if err := comment.Run(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}

func apiString(s string) *string {
	return &s
}

type gerritInfo struct {
	BotUsername    string
	Host           string
	Port           string
	ChangeNumber   string
	PatchSetNumber string
}

func readGerritInfo() gerritInfo {
	return gerritInfo{
		BotUsername:    os.Getenv("GERRIT_BOT_USERNAME"),
		Host:           os.Getenv("GERRIT_HOST"),
		Port:           os.Getenv("GERRIT_PORT"),
		ChangeNumber:   os.Getenv("GERRIT_CHANGE_NUMBER"),
		PatchSetNumber: os.Getenv("GERRIT_PATCHSET_NUMBER"),
	}
}
