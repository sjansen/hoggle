package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sjansen/hoggle/cmd"
)

func main() {
	cmd.RootCmd.AddCommand(versionCmd)
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print hoggle's version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("hoggle", version)
		return nil
	},
}
