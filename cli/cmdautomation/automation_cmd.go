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
		Aliases: []string{"auto", "aum", "py-auto", "scripts", "ai"},
		Short:   "High-performance native Go automation, search, newlines, and benchmarks",
		RunE:    runDefaultAutomationCmd,
	}

	searchCmd = &cobra.Command{
		Use:     "search <pattern> [dir]",
		Aliases: []string{"grep", "find-text", "find", "text", "scan"},
		Short:   "Multi-core streaming search with lazy regex and literal fast path",
		Example: `  gitmap aum search "letterly" --path "02-spec"
  gitmap aum search --root "cli" --query "cmdpipeline"
  gitmap aum search --query "func Run" --dir cli/store --ext .go
  gitmap aum search "Candidate Response|CANDIDATE RESPONSE" cli -i
  gitmap aum search "timer|Clock|8:47|border.*timer" src/runner.tsx -i
  gitmap aum search -r "(BEGIN [A-Z ]*PRIVATE KEY|AKIA[0-9A-Z]{16})" .`,
		RunE: runSearchCmd,
	}

	benchCmd = &cobra.Command{
		Use:     "benchmark [target]",
		Aliases: []string{"bench", "compare"},
		Short:   "Run side-by-side Go vs Python execution benchmarks",
		RunE:    runBenchmarkCmd,
	}

	searchPathFlag    string
	searchDirFlag     string
	searchRootFlag    string
	searchQueryFlag   string
	searchPatternFlag string

	searchOpts SearchOptions
)

func runDefaultAutomationCmd(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func runSearchCmd(cmd *cobra.Command, args []string) error {
	defer resetSearchFlags()
	pattern, dir := resolveSearchArgs(args)
	if pattern == "" {
		return apperror.NewValidationError("search pattern is required")
	}

	searchOpts.Pattern = pattern
	searchOpts.Dir = dir

	return executeAndRenderSearch(searchOpts)
}

func executeAndRenderSearch(opts SearchOptions) error {
	res, err := RunSearch(opts)
	if err != nil {
		return err
	}

	renderSearchResults(res)

	return nil
}

func resolveSearchArgs(args []string) (string, string) {
	flagPattern, flagDir := resolveExplicitSearchFlags()
	hasExplicitFlag := flagPattern != "" || flagDir != ""
	if hasExplicitFlag {
		return resolveHybridSearchArgs(flagPattern, flagDir, args)
	}

	return parseSearchPositionalArgs(args)
}

func resolveExplicitSearchFlags() (string, string) {
	flagPattern := resolveSearchPatternFlag()
	flagDir := resolveSearchDirFlag()

	return flagPattern, flagDir
}

func resolveSearchPatternFlag() string {
	flagPattern := strings.TrimSpace(searchQueryFlag)
	if flagPattern != "" {
		return flagPattern
	}

	return strings.TrimSpace(searchPatternFlag)
}

func resolveSearchDirFlag() string {
	flagDir := strings.TrimSpace(searchRootFlag)
	if flagDir != "" {
		return flagDir
	}
	if strings.TrimSpace(searchDirFlag) != "" {
		return strings.TrimSpace(searchDirFlag)
	}

	return strings.TrimSpace(searchPathFlag)
}

func resolveHybridSearchArgs(flagPattern, flagDir string, args []string) (string, string) {
	pattern := cleanSearchPattern(flagPattern)
	dir := cleanSearchPattern(flagDir)
	if pattern == "" && len(args) > 0 {
		pattern = cleanSearchPattern(args[0])
	}
	if dir == "" && len(args) > 0 && flagPattern != "" {
		dir = cleanSearchPattern(args[0])
	}
	if dir == "" && len(args) > 1 && flagPattern == "" {
		dir = cleanSearchPattern(args[1])
	}

	return pattern, dir
}

func resetSearchFlags() {
	searchPathFlag = ""
	searchDirFlag = ""
	searchRootFlag = ""
	searchQueryFlag = ""
	searchPatternFlag = ""
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
	return token == "automation" || token == "auto" || token == "aum" ||
		token == "py-auto" || token == "scripts" || token == "ai"
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
	initSearchPathFlags()
	initSearchPatternFlags()
	initSearchOptionFlags()
}

func initSearchPathFlags() {
	searchCmd.Flags().StringVarP(&searchPathFlag, "path", "p", "", "Target directory or file path to search")
	searchCmd.Flags().StringVarP(&searchDirFlag, "dir", "d", "", "Target directory or file path to search")
	searchCmd.Flags().StringVar(&searchRootFlag, "root", "", "Target root directory or file path to search")
}

func initSearchPatternFlags() {
	searchCmd.Flags().StringVarP(&searchQueryFlag, "query", "q", "", "Search query or keyword to find")
	searchCmd.Flags().StringVar(&searchPatternFlag, "pattern", "", "Search pattern or keyword to find")
}

func initSearchOptionFlags() {
	searchCmd.Flags().BoolVarP(&searchOpts.IsRegex, "regex", "r", false, "Use regular expression matching")
	searchCmd.Flags().BoolVarP(&searchOpts.IsCaseInsensitive, "ignore-case", "i", false, "Perform case-insensitive matching")
	searchCmd.Flags().StringSliceVarP(&searchOpts.Extensions, "ext", "e", nil, "Filter by file extensions (e.g. .go, .ts)")
	searchCmd.Flags().IntVarP(&searchOpts.Workers, "workers", "w", 0, "Number of worker threads (default: CPU count)")
	searchCmd.Flags().IntVar(&searchOpts.MaxJsonKb, "max-json-kb", 500, "Maximum JSON size in KB before auto-exclusion")
	searchCmd.Flags().BoolVar(&searchOpts.IncludeBinaries, "include-binaries", false, "Include binary files in search")
	searchCmd.Flags().BoolVar(&searchOpts.IncludeLargeJson, "include-large-json", false, "Include oversized JSON files in search")
}
