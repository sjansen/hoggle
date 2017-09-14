package cmd

import (
	"github.com/sjansen/hoggle/pkg/engine"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init URL",
	Short: "Configure Git LFS to use hoggle as a custom transfer agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return engine.Init(args[0])
	},
}
