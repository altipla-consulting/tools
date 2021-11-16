package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/box"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/configure-dev-machine/internal/config"
)

func init() {
	CmdRoot.AddCommand(CmdCheckUpdates)
}

var CmdCheckUpdates = &cobra.Command{
	Use:   "check-updates",
	Short: "Comprueba si hay actualizaciones de esta herramienta.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Jenkins y el entorno de desarrollo no deben comprobar la versión
		if Version == "dev" || os.Getenv("CI") == "true" {
			return nil
		}

		filename, err := config.Filename("last-update-check.txt")
		if err != nil {
			return errors.Trace(err)
		}

		lastUpdate := time.Time{}
		if content, err := ioutil.ReadFile(filename); err != nil && !os.IsNotExist(err) {
			return errors.Trace(err)
		} else if err == nil {
			if err := lastUpdate.UnmarshalText(content); err != nil {
				return errors.Trace(err)
			}
		}

		if time.Now().Sub(lastUpdate) > 1*time.Hour {
			reply, err := http.Get("https://tools.altipla.consulting/version-manifest/configure-dev-machine")
			if err != nil {
				return errors.Trace(err)
			}
			defer reply.Body.Close()
			if reply.StatusCode != http.StatusOK {
				return errors.Errorf("unexpected status: %s", reply.Status)
			}
			body, err := ioutil.ReadAll(reply.Body)
			if err != nil {
				return errors.Trace(err)
			}
			expected := strings.TrimSpace(string(body))

			if expected != Version {
				o := box.Box{}
				o.AddLine("Update available ", aurora.Gray(18, Version), " → ", aurora.BrightGreen(expected))
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
