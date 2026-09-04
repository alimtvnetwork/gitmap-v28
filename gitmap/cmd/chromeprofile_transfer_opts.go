package cmd

import (
	"strconv"
	"strings"
)

type chromeTransferOptions struct {
	Format     string
	Profile    string
	Email      string
	Limit      int
	Except     []string
	Positional []string
}

func parseChromeTransferOptions(args []string) chromeTransferOptions {
	opts := chromeTransferOptions{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if parseOptionFormat(arg, args, &i, &opts) {
			continue
		}
		if parseOptionLimit(arg, args, &i, &opts) {
			continue
		}
		if parseOptionProfile(arg, args, &i, &opts) {
			continue
		}
		if parseOptionEmail(arg, args, &i, &opts) {
			continue
		}
		if parseOptionExcept(arg, args, &i, &opts) {
			continue
		}
		if !strings.HasPrefix(arg, "-") {
			checkAndAssignPositionalEmail(arg, &opts)
			opts.Positional = append(opts.Positional, arg)
		}
	}
	return opts
}

func checkAndAssignPositionalEmail(arg string, opts *chromeTransferOptions) {
	if opts.Email == "" && strings.Contains(arg, "@") && !isSnapshotFileExtension(arg) {
		opts.Email = arg
	}
}

func isSnapshotFileExtension(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".json") ||
		strings.HasSuffix(lower, ".zip") ||
		strings.HasSuffix(lower, ".sqlite") ||
		strings.HasSuffix(lower, ".db") ||
		strings.HasSuffix(lower, ".yaml") ||
		strings.HasSuffix(lower, ".yml") ||
		strings.HasSuffix(lower, ".csv")
}

func parseOptionFormat(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--format=") {
		opts.Format = strings.TrimPrefix(arg, "--format=")
		return true
	}
	if arg == "--format" && *i+1 < len(args) {
		*i++
		opts.Format = args[*i]
		return true
	}
	return false
}

func parseOptionLimit(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--limit=") {
		opts.Limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "--limit="))
		return true
	}
	if strings.HasPrefix(arg, "-n=") {
		opts.Limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "-n="))
		return true
	}
	if strings.HasPrefix(arg, "-l=") {
		opts.Limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "-l="))
		return true
	}
	isLimitFlag := arg == "--limit" || arg == "-n" || arg == "-l" || arg == "-limit"
	if isLimitFlag && *i+1 < len(args) {
		*i++
		opts.Limit, _ = strconv.Atoi(args[*i])
		return true
	}
	return false
}

func parseOptionProfile(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--profile=") {
		opts.Profile = strings.TrimPrefix(arg, "--profile=")
		return true
	}
	if strings.HasPrefix(arg, "-p=") {
		opts.Profile = strings.TrimPrefix(arg, "-p=")
		return true
	}
	isProfileFlag := arg == "--profile" || arg == "-p" || arg == "-profile"
	if isProfileFlag && *i+1 < len(args) {
		*i++
		opts.Profile = args[*i]
		return true
	}
	return false
}

func parseOptionEmail(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--email=") {
		opts.Email = strings.TrimPrefix(arg, "--email=")
		return true
	}
	if strings.HasPrefix(arg, "-email=") {
		opts.Email = strings.TrimPrefix(arg, "-email=")
		return true
	}
	isFlag := arg == "--email" || arg == "-email"
	if isFlag && *i+1 < len(args) {
		*i++
		opts.Email = args[*i]
		return true
	}
	return false
}

func parseOptionExcept(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--except=") {
		appendExceptValues(opts, strings.TrimPrefix(arg, "--except="))
		return true
	}
	if strings.HasPrefix(arg, "--exclude=") {
		appendExceptValues(opts, strings.TrimPrefix(arg, "--exclude="))
		return true
	}
	if strings.HasPrefix(arg, "--skip=") {
		appendExceptValues(opts, strings.TrimPrefix(arg, "--skip="))
		return true
	}
	if strings.HasPrefix(arg, "-e=") {
		appendExceptValues(opts, strings.TrimPrefix(arg, "-e="))
		return true
	}
	isFlag := arg == "--except" || arg == "--exclude" || arg == "--skip" || arg == "-e"
	if isFlag && *i+1 < len(args) {
		*i++
		appendExceptValues(opts, args[*i])
		return true
	}
	return false
}

func appendExceptValues(opts *chromeTransferOptions, val string) {
	parts := strings.Split(val, ",")
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			opts.Except = append(opts.Except, trimmed)
		}
	}
}
