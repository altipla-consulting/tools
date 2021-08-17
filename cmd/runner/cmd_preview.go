package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

type previewFlags struct {
	Project string
	Tag     string
}

var (
	flagPreview previewFlags
)

func init() {
	cmdRoot.AddCommand(cmdPreview)
	cmdPreview.PersistentFlags().StringVar(&flagPreview.Project, "project", "", "Google Cloud project where the container will be stored. Defaults to the GOOGLE_PROJECT environment variable.")
	cmdPreview.PersistentFlags().StringVar(&flagDeploy.Tag, "tag", "", "Name of the revision included in the URL. Defaults to the Gerrit change and patchset.")
}

var cmdPreview = &cobra.Command{
	Use:   "preview",
	Short: "Send preview URLs as a comment to Gerrit.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		if flagPreview.Project == "" {
			flagPreview.Project = os.Getenv("GOOGLE_PROJECT")
		}

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")
		if os.Getenv("BUILD_CAUSE") == "SCMTRIGGER" {
			version += ".preview"
			if flagPreview.Tag == "" {
				flagPreview.Tag = "preview-" + os.Getenv("GERRIT_CHANGE_NUMBER") + "-" + os.Getenv("GERRIT_PATCHSET_NUMBER")
			}
		}

		suffixcmd := exec.Command(
			"gcloud",
			"run", "services", "describe",
			args[0],
			"--format", "value(status.url)",
			"--region", "europe-west1",
			"--project", flagPreview.Project,
		)
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
			previews = append(previews, app+" :: https://"+flagPreview.Tag+"---"+app+"-"+suffix+"-ew.a.run.app/")
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

		return nil
	},
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
