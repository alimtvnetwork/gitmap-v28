package cmdmacro

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

// ParseExecOptions extracts execution configuration from CLI arguments.
func ParseExecOptions(flagArgs []string) macro.ExecOptions {
	opts := macro.ExecOptions{}
	for i := 0; i < len(flagArgs); i++ {
		arg := flagArgs[i]
		checkExecutionModeFlags(arg, &opts)
		checkOutputSuppressionFlags(arg, &opts)
		checkTreeSummaryFlags(arg, &opts)
		checkSerializationFlags(arg, &opts)
		checkFilePathFlags(flagArgs, &i, &opts)
		checkLogFilePathFlags(flagArgs, &i, &opts)
		checkInlineFileArg(arg, &opts)
	}

	return opts
}

func checkExecutionModeFlags(arg string, opts *macro.ExecOptions) {
	switch strings.ToLower(arg) {
	case "--dry-run":
		opts.DryRun = true
	case "--run-until", "--run-until-end", "--keep-going", "--continue-on-error":
		opts.IsRunUntil = true
	case "--verbose", "-v":
		opts.Verbose = true
	}
}

func checkOutputSuppressionFlags(arg string, opts *macro.ExecOptions) {
	switch strings.ToLower(arg) {
	case "--no-terminal", "--quiet", "-q", "--silent":
		opts.IsTerminalSuppressed = true
	}
}

func checkTreeSummaryFlags(arg string, opts *macro.ExecOptions) {
	switch strings.ToLower(arg) {
	case "--summary":
		opts.IsSummaryOnly = true
	case "--no-tree":
		opts.IsTreeSuppressed = true
	}
}

func checkSerializationFlags(arg string, opts *macro.ExecOptions) {
	switch strings.ToLower(arg) {
	case "--json":
		opts.JSON = true
	case "--yaml", "--yml", "-y":
		opts.YAML = true
	}
}

func checkFilePathFlags(flagArgs []string, i *int, opts *macro.ExecOptions) {
	if isFileFlagWithArg(flagArgs[*i]) && *i+1 < len(flagArgs) {
		opts.FilePath = flagArgs[*i+1]
		*i++
	}
}

func checkLogFilePathFlags(flagArgs []string, i *int, opts *macro.ExecOptions) {
	low := strings.ToLower(flagArgs[*i])
	hasLogFlag := low == "--log" || low == "--log-file" || low == "--logfile"
	if hasLogFlag && *i+1 < len(flagArgs) {
		opts.LogFilePath = flagArgs[*i+1]
		*i++
	}
}
