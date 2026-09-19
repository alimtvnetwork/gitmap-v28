package cmdautomation

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	runCmd = &cobra.Command{
		Use:     "run <runtime> [code|file]",
		Aliases: []string{"exec", "worker"},
		Short:   "Execute inline code or script files across polyglot worker pool",
		RunE:    runRunCmd,
	}

	workerOpts WorkerRunOptions
)

func runRunCmd(cmd *cobra.Command, args []string) error {
	validErr := validateRunArgs(args)
	if validErr != nil {
		return validErr
	}
	opts := buildRunOptions(args)
	return toError(RunWorkerPool(opts))
}

func validateRunArgs(args []string) *apperror.AppError {
	if len(args) == 0 {
		return apperror.NewValidationError("runtime is required (e.g. py, node, go, rust, ps, bash)")
	}
	if len(args) < 2 {
		return apperror.NewValidationError("inline code string or script file path is required")
	}
	return nil
}

func buildRunOptions(args []string) WorkerRunOptions {
	opts := workerOpts
	opts.Runtime = args[0]
	opts.Target = args[1]
	return opts
}

func initRunFlags() {
	runCmd.Flags().IntVarP(&workerOpts.Workers, "workers", "w", 0, "Number of worker groups (default: CPU cores)")
	runCmd.Flags().IntVar(&workerOpts.Threads, "threads", 1, "Parallel worker threads per worker group")
	runCmd.Flags().StringVar(&workerOpts.Encoding, "encoding", "utf-8", "Stream encoding ('utf-8' or 'utf16')")
	runCmd.Flags().IntVar(&workerOpts.TimeoutSec, "timeout", 45, "Execution timeout in seconds per worker")
	runCmd.Flags().BoolVar(&workerOpts.AsJson, "json", false, "Output worker results as JSON")
	runCmd.Flags().BoolVar(&workerOpts.IsPreRead, "pre-read", false, "Pre-read file contents into context before dispatch")
}

func init() {
	AutomationCmd.AddCommand(runCmd)
	initRunFlags()
}
