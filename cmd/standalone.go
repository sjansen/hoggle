package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/sjansen/hoggle/pkg/engine"
)

func init() {
	RootCmd.AddCommand(standaloneCmd)
}

var standaloneCmd = &cobra.Command{
	Use:   "standalone URL",
	Short: "Run in standalone mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing required argument")
		}
		if len(args) > 1 {
			return errors.New("too many arguments")
		}
		return engine.Standalone(args[0])
	},
}
