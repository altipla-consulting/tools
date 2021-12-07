package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/box"
	"libs.altipla.consulting/env"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/config"
	"tools.altipla.consulting/cmd/configure-dev-machine/internal/version"
)

func init() {
	CmdRoot.AddCommand(CmdCheckUpdates)
}

var CmdCheckUpdates = &cobra.Command{
	Use:   "check-updates",
	Short: "Check if there are any updates to the tool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if env.IsJenkins() {
			return nil
		}

		filename, err := config.UserFile("last-update-check.txt")
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
			info, err := version.FetchInfo()
			if err != nil {
				return errors.Trace(err)
			}
			if info.CurrentlyInstalled != info.Latest {
				o := box.Box{}
				o.AddLine("Update available ", aurora.Gray(18, info.CurrentlyInstalled), " → ", aurora.BrightGreen(info.Latest))
				o.AddLine()
				o.AddLine("Run the following command to update:")
				o.AddLine(aurora.Blue("apt update && apt install configure-dev-machine && configure-dev-machine install"))
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
