package cmdagy

import (
	"strconv"
)

func defaultAgyFixOptions() AgyFixOptions {
	return AgyFixOptions{
		IsDetailed:    agyFixDetailed,
		IsNoRelease:   agyFixNoRelease,
		CustomPrompt:  agyFixCustomPrompt,
		IsNoClipboard: agyFixNoClipboard,
		OutputFile:    agyFixOutputFile,
		IsDryRun:      agyFixDryRun,
		IsForce:       agyFixForce,
		IsAll:         agyFixAll,
		ProjectsCount: agyFixProjects,
		Limit:         agyFixLimit,
		IsResetBatch:  agyFixResetBatch,
		IsNoInject:    agyFixNoInject,
	}
}

func parseAgyFixArgs(args []string) AgyFixOptions {
	opts := defaultAgyFixOptions()
	for i := 0; i < len(args); i++ {
		parseSingleArg(args, &i, &opts)
	}

	return opts
}

func parseBatchToggle(arg string, opts *AgyFixOptions) bool {
	switch arg {
	case "--all":
		opts.IsAll = true
		return true
	case "--reset-batch":
		opts.IsResetBatch = true
		return true
	case "--no-inject":
		opts.IsNoInject = true
		return true
	default:
		return false
	}
}

func parseBehaviorToggle(arg string, opts *AgyFixOptions) bool {
	if isAgyForceFlag(arg) {
		opts.IsForce = true
		return true
	}
	if isAgyDetailedFlag(arg) {
		opts.IsDetailed = true
		return true
	}

	return parseBatchToggle(arg, opts)
}

func parseOutputToggle(arg string, opts *AgyFixOptions) bool {
	switch arg {
	case "--no-release":
		opts.IsNoRelease = true
		return true
	case "--no-clipboard":
		opts.IsNoClipboard = true
		return true
	default:
		if isAgyDryRunFlag(arg) {
			opts.IsDryRun = true
			return true
		}
		return false
	}
}

func parseToggleArg(arg string, opts *AgyFixOptions) bool {
	return parseBehaviorToggle(arg, opts) || parseOutputToggle(arg, opts)
}

func parseNumericParam(args []string, idx *int, target *int) bool {
	if *idx+1 >= len(args) {
		return false
	}

	val, err := strconv.Atoi(args[*idx+1])
	if err != nil {
		return false
	}

	*target = val
	*idx++

	return true
}

func parseBatchParamFlag(args []string, idx *int, opts *AgyFixOptions, arg string) bool {
	switch arg {
	case "--projects":
		return parseNumericParam(args, idx, &opts.ProjectsCount)
	case "--limit":
		return parseNumericParam(args, idx, &opts.Limit)
	default:
		return false
	}
}
