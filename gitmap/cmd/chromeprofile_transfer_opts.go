package cmd

import (
	"strconv"
	"strings"
)

type chromeTransferOptions struct {
	Format     string
	Profile    string
	Limit      int
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
		if !strings.HasPrefix(arg, "-") {
			opts.Positional = append(opts.Positional, arg)
		}
	}
	return opts
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
	isLimitFlag := arg == "--limit" || arg == "-n" || arg == "-limit"
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
