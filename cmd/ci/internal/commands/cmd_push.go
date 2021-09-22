package commands

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/ci/internal/query"
	"tools.altipla.consulting/cmd/ci/internal/run"
)

var cmdPush = &cobra.Command{
	Use:     "push",
	Short:   "Envía el commit a Gerrit/GitHub.",
	Example: "ci push",
	RunE: func(cmd *cobra.Command, args []string) error {
		gerrit, err := query.IsGerrit()
		if err != nil {
			return errors.Trace(err)
		}
		mainBranch, err := query.MainBranch()
		if err != nil {
			return errors.Trace(err)
		}

		if gerrit {
			if err := run.Git("push", "origin", "HEAD:refs/for/"+mainBranch); err != nil {
				return errors.Trace(err)
			}
			return nil
		}

		return errors.Trace(run.Git("push"))
	},
}
