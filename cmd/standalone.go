package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(standaloneCmd)
}

var standaloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "Used internally by Git LFS",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}
