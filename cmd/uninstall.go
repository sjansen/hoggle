package cmd

import (
	"github.com/sjansen/hoggle/pkg/engine"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Update Git LFS to stop using hoggle as a custom transfer agent",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return engine.Uninstall()
	},
}
