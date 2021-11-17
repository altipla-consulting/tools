package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/box"
	"libs.altipla.consulting/env"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/config"
)

func init() {
	CmdRoot.AddCommand(CmdCheckUpdates)
}

type ghRelease struct {
	TagName string `json:"tag_name"`
}

var CmdCheckUpdates = &cobra.Command{
	Use:   "check-updates",
	Short: "Check if there are any updates to the tool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if env.IsJenkins() {
			return nil
		}

		filename, err := config.Filename("last-update-check.txt")
		if err != nil {
			return errors.Trace(err)
		}
		var lastUpdate time.Time
		if content, err := ioutil.ReadFile(filename); err != nil && !os.IsNotExist(err) {
			return errors.Trace(err)
		} else if err == nil {
			if err := lastUpdate.UnmarshalText(content); err != nil {
				return errors.Trace(err)
			}
		}

		if time.Now().Sub(lastUpdate) > 1*time.Hour {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/altipla-consulting/tools/releases/latest", nil)
			if err != nil {
				return errors.Trace(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return errors.Trace(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return errors.Errorf("unexpected github status: %s", resp.Status)
			}
			var release ghRelease
			if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
				return errors.Trace(err)
			}

			version, err := config.Version()
			if err != nil {
				return errors.Trace(err)
			}
			if release.TagName != version {
				o := box.Box{}
				o.AddLine("Update available ", aurora.Gray(18, version), " → ", aurora.BrightGreen(release.TagName))
				o.AddLine()
				o.AddLine("Run the following command to update:")
				o.AddLine(aurora.Blue("curl https://tools.altipla.consulting/install/configure-dev-machine | bash"))
				o.Render()

				return nil
			}

			check, err := time.Now().MarshalText()
			if err != nil {
				return errors.Trace(err)
			}
			if err := ioutil.WriteFile(filename, check, 0600); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
