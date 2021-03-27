package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/ci/internal/pr"
	"tools.altipla.consulting/cmd/ci/internal/query"
	"tools.altipla.consulting/cmd/ci/internal/run"
)

var CmdCheckout = &cobra.Command{
	Use:   "checkout",
	Short: "Establece el código a la versión exacta de un Pull Request en GitHub",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.Errorf("Especifica como primer argumento el ID del PR que quieres descargar: %v: %v", args[0], err)
		}

		branch, err := pr.Branch(ctx, id)
		if err != nil {
			return errors.Trace(err)
		}

		exists, err := query.BranchExists(branch)
		if err != nil {
			return errors.Trace(err)
		}
		if exists {
			if err := run.Git("branch", "-D", branch); err != nil {
				return errors.Trace(err)
			}
		}
		if err := run.Git("fetch", "origin", fmt.Sprintf("pull/%d/head:%s", id, branch)); err != nil {
			return errors.Trace(err)
		}

		if err := run.Git("checkout", branch); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}

var CmdCheckoutShort = &cobra.Command{
	Use:   "co",
	Short: CmdCheckout.Short,
	Args:  CmdCheckout.Args,
	RunE:  CmdCheckout.RunE,
}
