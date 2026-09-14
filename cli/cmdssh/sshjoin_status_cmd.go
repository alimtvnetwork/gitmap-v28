package cmdssh

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	statusPortFlag    int
	statusTimeoutFlag time.Duration
)

// SJStatusCmd represents the 'gitmap ssh-join status' subcommand.
var SJStatusCmd = &cobra.Command{
	Use:     "status [alias|ip] [flags]",
	Aliases: []string{"ping", "health", "check"},
	Short:   "Check connectivity and health of registered SSH machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJStatus(cmd, args, cmd.Context())
	},
}

func extractTargetArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func buildHealthOptions(args []string) SSHHealthOptions {
	return SSHHealthOptions{
		Target:  extractTargetArg(args),
		Port:    statusPortFlag,
		Timeout: statusTimeoutFlag,
	}
}

//nolint:revive
func runSJStatus(cmd *cobra.Command, args []string, ctx context.Context) error {
	opts := buildHealthOptions(args)
	_, appErr := ExecuteHealthCheck(ctx, os.Stdout, opts)
	if appErr != nil {
		return appErr
	}
	return nil
}

// RunSJStatus executes ssh-join status subcommand programmatically.
func RunSJStatus(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSJStatus(cmd, args, ctx)
}

func init() {
	SJStatusCmd.Flags().IntVarP(&statusPortFlag, "port", "p", 22, "Target SSH port to probe")
	SJStatusCmd.Flags().DurationVarP(&statusTimeoutFlag, "timeout", "t", 1500*time.Millisecond, "Probe timeout duration")
	SSHJoinCmd.AddCommand(SJStatusCmd)
}
