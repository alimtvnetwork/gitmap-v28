package cmdagy

import (
	"context"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	agyROPDryRun bool
)

// AGYROPCmd represents the 'gitmap agy reread-optimize-project' command.
var AGYROPCmd = &cobra.Command{
	Use:     "reread-optimize-project [N] [flags]",
	Aliases: []string{"rop", "reread-optimize", "re-read-optimize", "optimize-reread"},
	Short:   "Reread and optimize top N recently communicated projects with Split-DB conversation backup",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAGYROPCLI(args)
	},
}

func init() {
	AGYROPCmd.Flags().BoolVarP(&agyROPDryRun, "dry-run", "d", false, "Preview reread and optimization without modifying conversations")
	AgyCmd.AddCommand(AGYROPCmd)
}

func parseROPArgs(args []string) (int, bool) {
	n := 5
	isDryRun := agyROPDryRun

	for _, a := range args {
		if a == "--dry-run" || a == "-d" {
			isDryRun = true
			continue
		}
		if val, err := strconv.Atoi(a); err == nil && val > 0 {
			n = val
			continue
		}
	}
	return n, isDryRun
}

// RunAGYROPCLI executes reread-optimize-project from commandline args.
func RunAGYROPCLI(args []string) error {
	n, isDryRun := parseROPArgs(args)
	ctx := context.Background()
	_, appErr := ExecuteROP(ctx, os.Stdout, n, isDryRun)
	if appErr != nil {
		return appErr
	}
	return nil
}
