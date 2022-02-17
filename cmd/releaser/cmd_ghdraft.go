package main

import (
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/releaser/internal/git"
	"tools.altipla.consulting/cmd/releaser/internal/run"
)

var cmdGHDraft = &cobra.Command{
	Use:     "ghdraft",
	Short:   "Open a GitHub release draft with the info filled from the latest release.",
	Example: "releaser ghdraft",
	RunE: func(command *cobra.Command, args []string) error {
		latest, err := git.LatestRemoteTag()
		if err != nil {
			return errors.Trace(err)
		}
		if latest == "" {
			return errors.Errorf("There is no latest release to open a GitHub draft.")
		}

		remote, err := git.RemoteURL("origin")
		if err != nil {
			return errors.Trace(err)
		}
		repo := "https://github.com/" + strings.TrimSuffix(strings.Split(remote, ":")[1], ".git")
		u, err := url.Parse(repo + "/releases/new")
		if err != nil {
			return errors.Trace(err)
		}
		qs := make(url.Values)
		qs.Set("tag", latest)
		notes, err := releaseNotes(repo, latest)
		if err != nil {
			return errors.Trace(err)
		}
		qs.Set("body", notes)
		u.RawQuery = qs.Encode()
		if err := run.OpenBrowser(u.String()); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}
