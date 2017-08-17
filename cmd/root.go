package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().Bool("help", false, "help for "+RootCmd.Name())
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// RootCmd is executed when no arguments are provided.
var RootCmd = &cobra.Command{
	Use:          "hoggle",
	Short:        "hoggle - standalone custom transfer agent for Git LFS",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}
