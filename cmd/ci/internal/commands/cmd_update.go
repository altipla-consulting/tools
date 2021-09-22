package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/ci/internal/prompt"
	"tools.altipla.consulting/cmd/ci/internal/query"
	"tools.altipla.consulting/cmd/ci/internal/run"
)

var flagForce bool

func init() {
	cmdUpdate.Flags().BoolVarP(&flagForce, "force", "f", false, "Fuerza la actualización aunque haya cambios pendientes. WARNING: Es una operación destructiva.")
}

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Actualiza a la última versión de master borrando todo lo que haya en local",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run.Git("fetch", "origin"); err != nil {
			return errors.Trace(err)
		}

		mainBranch, err := query.MainBranch()
		if err != nil {
			return errors.Trace(err)
		}

		status, err := run.GitCaptureOutput("status", "-s")
		if err != nil {
			return errors.Trace(err)
		}
		if len(status) > 0 && !flagForce {
			keep, err := prompt.Confirm(fmt.Sprintf("El proyecto tiene cambios. ¿Estás seguro de que deseas borrar todo y pasar a %s?", mainBranch))
			if err != nil {
				return errors.Trace(err)
			}
			if !keep {
				return nil
			}
		}

		if err := run.Git("checkout", "--", "."); err != nil {
			return errors.Trace(err)
		}
		if err := run.Git("checkout", mainBranch); err != nil {
			return errors.Trace(err)
		}
		if err := run.Git("reset", "--hard", "origin/"+mainBranch); err != nil {
			return errors.Trace(err)
		}

		// Remove any untracked file remaining in the working directory
		status, err = run.GitCaptureOutput("status", "-s")
		if err != nil {
			return errors.Trace(err)
		}
		if len(status) > 0 {
			if err := run.Git("stash", "-q", "--include-untracked"); err != nil {
				return errors.Trace(err)
			}
			if err := run.Git("stash", "drop", "-q"); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
