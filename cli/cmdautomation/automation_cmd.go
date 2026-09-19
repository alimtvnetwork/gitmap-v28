package cmdautomation

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	// AutomationCmd is the root command for native Go repository automation.
	AutomationCmd = &cobra.Command{
		Use:     "automation",
		Aliases: []string{"auto", "py-auto", "scripts", "ai"},
		Short:   "High-performance native Go automation, search, newlines, and benchmarks",
		RunE:    runDefaultAutomationCmd,
	}

	searchCmd = &cobra.Command{
		Use:     "search <pattern> [dir]",
		Aliases: []string{"grep", "find-text"},
		Short:   "Multi-core streaming search with lazy regex and literal fast path",
		RunE:    runSearchCmd,
	}

	benchCmd = &cobra.Command{
		Use:     "benchmark [target]",
		Aliases: []string{"bench", "compare"},
		Short:   "Run side-by-side Go vs Python execution benchmarks",
		RunE:    runBenchmarkCmd,
	}

	searchOpts SearchOptions
)

func runDefaultAutomationCmd(cmd *cobra.Command, args []string) error {
	return AutomationCmd.Help()
}

func runSearchCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("search pattern is required")
	}
	searchOpts.Pattern = args[0]
	if len(args) > 1 {
		searchOpts.Dir = args[1]
	}
	res, err := RunSearch(searchOpts)
	if err != nil {
		return err
	}
	renderSearchResults(res)
	return nil
}

func runBenchmarkCmd(cmd *cobra.Command, args []string) error {
	target := "all"
	if len(args) > 0 {
		target = args[0]
	}
	return toError(RunBenchmark(target))
}

// DispatchAutomation routes CLI arguments to native automation commands.
func DispatchAutomation(args []string) error {
	clean := stripAutomationPrefix(args)
	if len(clean) == 0 {
		return AutomationCmd.Help()
	}
	AutomationCmd.SetArgs(clean)
	return AutomationCmd.Execute()
}

func stripAutomationPrefix(args []string) []string {
	if len(args) == 0 {
		return args
	}
	first := strings.ToLower(args[0])
	if isAutomationTrigger(first) {
		return args[1:]
	}
	return args
}

func isAutomationTrigger(token string) bool {
	return token == "automation" || token == "auto" || token == "py-auto" ||
		token == "scripts" || token == "ai"
}

func toError(appErr *apperror.AppError) error {
	if appErr != nil {
		return appErr
	}
	return nil
}

func init() {
	AutomationCmd.AddCommand(searchCmd)
	AutomationCmd.AddCommand(benchCmd)

	initSearchFlags()
	initSubsystems()
}

func initSearchFlags() {
	searchCmd.Flags().BoolVarP(&searchOpts.IsRegex, "regex", "r", false, "Use regular expression matching")
	searchCmd.Flags().BoolVarP(&searchOpts.IsCaseInsensitive, "ignore-case", "i", false, "Perform case-insensitive matching")
	searchCmd.Flags().StringSliceVarP(&searchOpts.Extensions, "ext", "e", nil, "Filter by file extensions (e.g. .go, .ts)")
	searchCmd.Flags().IntVarP(&searchOpts.Workers, "workers", "w", 0, "Number of worker threads (default: CPU count)")
}
