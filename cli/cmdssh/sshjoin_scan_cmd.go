package cmdssh

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	scanPortFlag    int
	scanTimeoutFlag time.Duration
	scanWorkersFlag int
)

// SJScanCmd represents the 'gitmap ssh-join scan' subcommand.
var SJScanCmd = &cobra.Command{
	Use:     "scan [subnet] [flags]",
	Aliases: []string{"find", "discover", "probe"},
	Short:   "Scan local subnet or CIDR for active SSH machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJScan(cmd, args, cmd.Context())
	},
}

func extractSubnetArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func buildScanOptions(args []string) SSHScanOptions {
	return SSHScanOptions{
		Subnet:  extractSubnetArg(args),
		Port:    scanPortFlag,
		Timeout: scanTimeoutFlag,
		Workers: scanWorkersFlag,
	}
}

//nolint:revive
func runSJScan(cmd *cobra.Command, args []string, ctx context.Context) error {
	opts := buildScanOptions(args)
	_, appErr := ExecuteSubnetScan(ctx, os.Stdout, opts)
	if appErr != nil {
		return appErr
	}
	return nil
}

// RunSJScan executes ssh-join scan subcommand programmatically.
func RunSJScan(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSJScan(cmd, args, ctx)
}

func init() {
	SJScanCmd.Flags().IntVarP(&scanPortFlag, "port", "p", 22, "Target SSH port to probe")
	SJScanCmd.Flags().DurationVarP(&scanTimeoutFlag, "timeout", "t", 800*time.Millisecond, "Per-host probe timeout")
	SJScanCmd.Flags().IntVarP(&scanWorkersFlag, "workers", "w", 40, "Concurrent worker pool size")
	SSHJoinCmd.AddCommand(SJScanCmd)
}
