package generate

import (
	"fmt"

	"github.com/spf13/cobra"
)

const middlewareExample = `$ enlight g middleware auth`

// MiddlewareCmd generates a new middleware/middleware file and a stub test.
var MiddlewareCmd = &cobra.Command{
	Use:     "middleware [name]",
	Example: middlewareExample,
	Aliases: []string{"m"},
	Short:   "Generate a new middleware/middleware file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must supply a name")
		}
		return nil
	},
}
