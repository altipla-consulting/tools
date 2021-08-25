package preview

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

type cmdFlags struct {
	Project string
	Tag     string
}

var (
	flags cmdFlags
)

func init() {
	Cmd.PersistentFlags().StringVar(&flags.Project, "project", "", "Google Cloud project where the container will be stored. Defaults to the GOOGLE_PROJECT environment variable.")
	Cmd.PersistentFlags().StringVar(&flags.Tag, "tag", "", "Name of the revision included in the URL. Defaults to the Gerrit change and patchset.")
}

var Cmd = &cobra.Command{
	Use:   "preview",
	Short: "Send preview URLs as a comment to Gerrit.",
	Args:  cobra.ExactArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		app := args[0]

		if flags.Project == "" {
			flags.Project = os.Getenv("GOOGLE_PROJECT")
		}

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")
		if os.Getenv("BUILD_CAUSE") == "SCMTRIGGER" {
			version += ".preview"
			if flags.Tag == "" {
				flags.Tag = "preview-" + os.Getenv("GERRIT_CHANGE_NUMBER") + "-" + os.Getenv("GERRIT_PATCHSET_NUMBER")
			}
		}

		suffixcmd := exec.Command(
			"gcloud",
			"run", "services", "describe",
			app,
			"--format", "value(status.url)",
			"--region", "europe-west1",
			"--project", flags.Project,
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
		previews = append(previews, app+" :: https://"+flags.Tag+"---"+app+"-"+suffix+"-ew.a.run.app/")
		log.WithField("previews", previews).Debug("Send comment to Gerrit with the previews")

		log.Info("Send preview URLs as a Gerrit comment")
		gerrit := readGerritInfo()
		ssh := []string{
			"ssh",
			"-p", gerrit.Port,
			fmt.Sprintf("%s@%s", gerrit.BotUsername, gerrit.Host),
			"gerrit", "review", fmt.Sprintf("%v,%v", gerrit.ChangeNumber, gerrit.PatchSetNumber),
			"--message", `"Previews deployed at:` + "\n" + strings.Join(previews, "\n") + `"`,
		}
		log.Debug(strings.Join(ssh, " "))
		comment := exec.Command(ssh[0], ssh[1:]...)
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
