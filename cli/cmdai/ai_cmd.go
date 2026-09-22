package cmdai

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	// AiCmd is the root Cobra command for GitMap native AI scripts automation.
	AiCmd = &cobra.Command{
		Use:     "ai",
		Aliases: []string{"scripts"},
		Short:   "Native AI scripts catalog, execution, and autofix engine",
		RunE:    runDefaultAiCmd,
	}

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "catalog"},
		Short:   "List registered AI automation scripts",
		RunE:    runListCmd,
	}

	runCmd = &cobra.Command{
		Use:     "run <token> [args...]",
		Aliases: []string{"exec"},
		Short:   "Execute an AI automation script with streaming output",
		RunE:    runRunCmd,
	}

	fixCmd = &cobra.Command{
		Use:   "fix [target] [args...]",
		Short: "Run standardized repository autofix targets",
		RunE:  runFixCmd,
	}

	pwshCmd = &cobra.Command{
		Use:     "pwsh [command/args...]",
		Aliases: []string{"ps", "powershell"},
		Short:   "Execute PowerShell automation and record execution history",
		RunE:    runPwshCmd,
	}

	historyCmd = &cobra.Command{
		Use:     "history",
		Aliases: []string{"hist", "frequent", "freq"},
		Short:   "Inspect and copy frequently executed AI commands",
		RunE:    runHistoryCmd,
	}

	listOpts ScriptListOptions
)

func runDefaultAiCmd(cmd *cobra.Command, args []string) error {
	return toError(RunAiList(listOpts))
}

func runListCmd(cmd *cobra.Command, args []string) error {
	return toError(RunAiList(listOpts))
}

func runRunCmd(cmd *cobra.Command, args []string) error {
	isEmpty := len(args) == 0
	if isEmpty {
		return apperror.NewValidationError("script token required (number, slug, or alias)")
	}

	return toError(RunAiScript(args[0], args[1:]))
}

func runPwshCmd(cmd *cobra.Command, args []string) error {
	return toError(RunAiPowershell(args))
}

func runHistoryCmd(cmd *cobra.Command, args []string) error {
	return toError(RunAiFrequent(25, listOpts.IsCopyClipboard))
}

func runFixCmd(cmd *cobra.Command, args []string) error {
	target := "all"
	extraArgs := []string{}
	hasTarget := len(args) > 0
	if hasTarget {
		target = args[0]
		extraArgs = args[1:]
	}

	return toError(RunAiFix(target, extraArgs))
}

// DispatchAi routes CLI arguments to native AI scripts commands.
func DispatchAi(args []string) error {
	cleanArgs := stripAiPrefix(args)
	isEmpty := len(cleanArgs) == 0
	if isEmpty {
		return toError(RunAiList(ScriptListOptions{}))
	}

	first := cleanArgs[0]
	if first == "-h" || first == "--help" || first == "help" {
		RenderAiHelp()
		return nil
	}

	isDirect := isDirectScriptDispatch(cleanArgs[0])
	if isDirect {
		return toError(RunAiScript(cleanArgs[0], cleanArgs[1:]))
	}

	AiCmd.SetArgs(cleanArgs)

	return AiCmd.Execute()
}

func stripAiPrefix(args []string) []string {
	isEmpty := len(args) == 0
	if isEmpty {
		return args
	}

	first := strings.ToLower(args[0])
	isPrefix := first == "ai" || first == "scripts"
	if isPrefix {
		return args[1:]
	}

	return args
}

func isDirectScriptDispatch(token string) bool {
	isCommand := isAiSubcommand(token)
	if isCommand {
		return false
	}

	_, findErr := FindScriptByToken(token)
	isFound := findErr == nil

	return isFound
}

func isAiSubcommand(token string) bool {
	clean := strings.ToLower(token)
	isKnown := clean == "list" || clean == "ls" || clean == "catalog" ||
		clean == "run" || clean == "exec" || clean == "fix" ||
		clean == "create" || clean == "new" || clean == "scaffold" || clean == "gen" ||
		clean == "history" || clean == "hist" || clean == "frequent" || clean == "freq" ||
		clean == "pwsh" || clean == "ps" || clean == "powershell" ||
		clean == "help" || clean == "--help" || clean == "-h"

	return isKnown
}

func toError(appErr *apperror.AppError) error {
	hasErr := appErr != nil
	if hasErr {
		return appErr
	}

	return nil
}

func init() {
	AiCmd.AddCommand(listCmd)
	AiCmd.AddCommand(historyCmd)
	AiCmd.AddCommand(runCmd)
	AiCmd.AddCommand(pwshCmd)
	AiCmd.AddCommand(fixCmd)

	initListFlags()
	initCreateCmd()
}

func initListFlags() {
	listCmd.Flags().StringVarP(&listOpts.CategoryFilter, "category", "c", "", "Filter scripts by category")
	listCmd.Flags().StringVarP(&listOpts.SearchQuery, "search", "s", "", "Search scripts by keyword")
	listCmd.Flags().BoolVarP(&listOpts.IsFixOnly, "fix", "f", false, "Show only scripts with fix mode")
	listCmd.Flags().BoolVar(&listOpts.IsJsonFormat, "json", false, "Output results in JSON format")
	listCmd.Flags().BoolVarP(&listOpts.IsVerbose, "verbose", "v", false, "Show verbose script information")
	listCmd.Flags().BoolVar(&listOpts.IsCopyClipboard, "copy", false, "Copy frequent AI commands to clipboard")
	listCmd.Flags().BoolVar(&listOpts.IsFrequentOnly, "frequent", false, "Show only frequent AI commands")

	historyCmd.Flags().BoolVar(&listOpts.IsCopyClipboard, "copy", false, "Copy frequent AI commands to clipboard")
}
