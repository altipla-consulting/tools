package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flagDebug bool

func init() {
	CmdRoot.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "Activa el logging de depuración")

	CmdRoot.AddCommand(CmdCheckout)
	CmdRoot.AddCommand(CmdCheckoutShort)
	CmdRoot.AddCommand(CmdLogin)
	CmdRoot.AddCommand(CmdPush)
	CmdRoot.AddCommand(CmdUpdate)
}

var CmdRoot = &cobra.Command{
	Use:           "ci",
	Short:         "Git related helper.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetFormatter(&log.TextFormatter{})
		if flagDebug {
			log.SetLevel(log.DebugLevel)
			log.Debug("DEBUG log level activated")
		}

		return nil
	},
}
