// Package cmd — prompt_args_parser.go parses CLI flags and targets for ct install-prompts.
package cmd

import (
	"strings"
)

type promptInstallOptions struct {
	Targets   []string
	Exclude   string
	IsDryRun  bool
	IsAll     bool
	Action    string
}

func parsePromptArgs(args []string) promptInstallOptions {
	var opts promptInstallOptions
	opts.Action = "install"

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--dry-run":
			opts.IsDryRun = true
		case a == "--all" || a == "-A":
			opts.IsAll = true
		case (a == "--exclude" || a == "-e") && i+1 < len(args):
			opts.Exclude = args[i+1]
			i++
		case strings.HasPrefix(a, "--exclude="):
			opts.Exclude = strings.TrimPrefix(a, "--exclude=")
		case a == "install" || a == "install-prompts" || a == "install-prompt":
			opts.Action = "install"
		case a == "update" || a == "update-prompts":
			opts.Action = "update"
		case a == "status" || a == "prompts-status":
			opts.Action = "status"
		case a == "version" || a == "prompts-version":
			opts.Action = "version"
		case !strings.HasPrefix(a, "-"):
			opts.Targets = append(opts.Targets, a)
		}
	}

	return opts
}
