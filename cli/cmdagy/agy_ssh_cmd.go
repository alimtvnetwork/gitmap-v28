package cmdagy

import (
	"github.com/spf13/cobra"
)

// RunAgySSHFn delegates agy SSH execution to cmdssh package.
var RunAgySSHFn func(args []string) error

var agySshCmd = &cobra.Command{
	Use:                "ssh [target] [flags] <command...>",
	Aliases:            []string{"remote"},
	Short:              "Execute AGY commands across remote SSH fleet nodes in parallel",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if RunAgySSHFn != nil {
			return RunAgySSHFn(args)
		}
		return nil
	},
}

func init() {
	AgyCmd.AddCommand(agySshCmd)
}
