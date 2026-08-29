// Package cmd — workdir_flags.go parses CLI options for work directory commands.
package cmd

import (
	"strings"
)

type workDirOptions struct {
	Action string
	Target string
	Label  string
	Args   []string
}

func parseWorkDirFlags(args []string) workDirOptions {
	var opts workDirOptions
	var cleanArgs []string

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case (a == "--label" || a == "-l" || a == "-label") && i+1 < len(args):
			opts.Label = args[i+1]
			i++
		case strings.HasPrefix(a, "--label="):
			opts.Label = strings.TrimPrefix(a, "--label=")
		default:
			cleanArgs = append(cleanArgs, a)
		}
	}

	if len(cleanArgs) == 0 {
		opts.Action = "ls"
		return opts
	}

	opts.Action = strings.ToLower(cleanArgs[0])
	if len(cleanArgs) > 1 {
		opts.Target = cleanArgs[1]
		opts.Args = cleanArgs[1:]
	}
	return opts
}
