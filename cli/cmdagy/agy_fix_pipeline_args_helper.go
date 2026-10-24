package cmdagy

import "strings"

func parsePromptOrFileFlag(args []string, idx *int, opts *AgyFixOptions, arg string) bool {
	if *idx+1 >= len(args) {
		return false
	}
	if arg == "--prompt" || arg == "-p" {
		opts.CustomPrompt = args[*idx+1]
		*idx++
		return true
	}
	if arg == "--file" {
		opts.OutputFile = args[*idx+1]
		*idx++
		return true
	}

	return false
}

func parseParamFlag(args []string, idx *int, opts *AgyFixOptions, arg string) bool {
	return parseBatchParamFlag(args, idx, opts, arg) || parsePromptOrFileFlag(args, idx, opts, arg)
}

func parseSingleArg(args []string, idx *int, opts *AgyFixOptions) {
	arg := args[*idx]
	if parseToggleArg(arg, opts) {
		return
	}

	parseArgWithParam(args, idx, opts, arg)
}

func parseArgWithParam(args []string, idx *int, opts *AgyFixOptions, arg string) {
	if parseParamFlag(args, idx, opts, arg) || strings.HasPrefix(arg, "-") {
		return
	}

	if !isSubcommandKeyword(strings.ToLower(arg)) && len(opts.Repo) == 0 {
		opts.Repo = arg
	}
}

func isSubcommandKeyword(word string) bool {
	switch word {
	case "fix", "errors", "error", "err", "agy", "aef",
		"pipeline", "pipeline-fix", "fix-agy", "agy-errors-fix",
		"fix-pipeline", "fixpipeline", "fp":
		return true
	}

	return false
}

func isAgyForceFlag(arg string) bool {
	return arg == "--force" || arg == "-f"
}

func isAgyDetailedFlag(arg string) bool {
	return arg == "--detailed" || arg == "-v"
}

func isAgyDryRunFlag(arg string) bool {
	return arg == "--dry-run" || arg == "-d"
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" || arg == "help" {
			return true
		}
	}

	return false
}
