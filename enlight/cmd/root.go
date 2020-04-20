package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the hook for all of the other commands in the enlight binary.
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "enlight",
	Short:         "Build Enlight applications with ease",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
