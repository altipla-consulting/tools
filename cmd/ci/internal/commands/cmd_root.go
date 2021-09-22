package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flagDebug bool

func init() {
	CmdRoot.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "Activa el logging de depuración")

	CmdRoot.AddCommand(cmdCheckout)
	CmdRoot.AddCommand(cmdCheckoutShort)
	CmdRoot.AddCommand(cmdLogin)
	CmdRoot.AddCommand(cmdPush)
	CmdRoot.AddCommand(cmdUpdate)
	CmdRoot.AddCommand(cmdPR)
}

var CmdRoot = &cobra.Command{
	Use:           "ci",
	Short:         "Git related helper.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetFormatter(new(log.TextFormatter))
		if flagDebug {
			log.SetLevel(log.DebugLevel)
			log.Debug("DEBUG log level activated")
		}

		return nil
	},
}
