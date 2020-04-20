package cmd

import (
	"github.com/juliankoehn/enlight/enlight/cmd/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate application components",
	Aliases: []string{"g"},
}

func init() {
	generateCmd.AddCommand(generate.MiddlewareCmd)
}
