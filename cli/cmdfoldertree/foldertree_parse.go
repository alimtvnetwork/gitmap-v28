package cmdfoldertree

import (
	"strconv"
	"strings"
)

func parseFolderTreeArgs(args []string) (FolderTreeOptions, []string) {
	opts := defaultFolderTreeOptions()
	var remaining []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if isHelpArg(a) {
			opts.Subcommand = "help"
			return opts, nil
		}
		if handled, nextI := handleFlagArg(a, args, i, &opts); handled {
			i = nextI
			continue
		}
		remaining = appendOrSetSubcommand(a, &opts, remaining)
	}
	return opts, remaining
}

func defaultFolderTreeOptions() FolderTreeOptions {
	return FolderTreeOptions{
		MaxDepth:    3,
		ShowFiles:   true,
		ShowNumbers: true,
		Format:      "json",
		IsTree:      true,
	}
}

func isHelpArg(a string) bool {
	return a == "-h" || a == "--help" || a == "help"
}

func appendOrSetSubcommand(a string, opts *FolderTreeOptions, remaining []string) []string {
	if opts.Subcommand == "" && isKnownSubcommand(a) {
		opts.Subcommand = strings.ToLower(a)
		return remaining
	}
	return append(remaining, a)
}

func isKnownSubcommand(a string) bool {
	switch strings.ToLower(a) {
	case "ls", "view", "show", "tree", "export", "exp", "import", "imp":
		return true
	default:
		return false
	}
}

func handleFlagArg(a string, args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if handleBooleanFlags(a, opts) {
		return true, i
	}
	return handleValueFlags(a, args, i, opts)
}

func handleBooleanFlags(a string, opts *FolderTreeOptions) bool {
	switch strings.ToLower(a) {
	case "--preview", "-p", "--gap":
		opts.IsPreview = true
		opts.IsTree = false
		return true
	case "--tree":
		opts.IsTree = true
		opts.IsPreview = false
		return true
	case "--dirs-only", "-d", "--folders-only":
		opts.DirsOnly = true
		opts.ShowFiles = false
		return true
	case "--files", "-f":
		opts.ShowFiles = true
		return true
	case "--no-numbers":
		opts.ShowNumbers = false
		return true
	case "--git-only":
		opts.GitOnly = true
		return true
	case "--all", "-a":
		opts.IncludeHidden = true
		return true
	case "--dry-run", "-n":
		opts.DryRun = true
		return true
	default:
		return false
	}
}

func handleValueFlags(a string, args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if a == "--depth" || a == "-L" {
		return extractDepthFlag(args, i, opts)
	}
	if a == "--format" || a == "-fmt" {
		return extractFormatFlag(args, i, opts)
	}
	if a == "--output" || a == "-o" || a == "--file" {
		return extractOutputFlag(args, i, opts)
	}
	return false, i
}

func extractDepthFlag(args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if i+1 >= len(args) {
		return false, i
	}
	if d, err := strconv.Atoi(args[i+1]); err == nil {
		opts.MaxDepth = d
	}
	return true, i + 1
}

func extractFormatFlag(args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if i+1 < len(args) {
		opts.Format = args[i+1]
		return true, i + 1
	}
	return false, i
}

func extractOutputFlag(args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if i+1 < len(args) {
		opts.OutputFile = args[i+1]
		return true, i + 1
	}
	return false, i
}
