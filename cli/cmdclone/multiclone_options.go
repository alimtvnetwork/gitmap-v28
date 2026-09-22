package cmdclone

import (
	"io"
	"os"
	"strconv"
	"strings"
)

// ResolveMultiCloneOptions parses flags, target directory, files, and multiline inputs.
func ResolveMultiCloneOptions(args []string) (MultiCloneOptions, error) {
	opts := MultiCloneOptions{Concurrency: 1}
	var textParts []string

	for i := 0; i < len(args); i++ {
		consumed, isFlag := parseMultiCloneFlag(args, i, &opts)
		if isFlag {
			i += consumed
			continue
		}

		handleMultiClonePositional(args[i], &opts, &textParts)
	}

	opts.RawInput = assembleRawInput(opts, textParts)
	return opts, nil
}

func parseMultiCloneFlag(args []string, idx int, opts *MultiCloneOptions) (int, bool) {
	arg := args[idx]
	if arg == "--dry-run" {
		opts.IsDryRun = true
		return 0, true
	}
	if arg == "--no-replace" {
		opts.NoReplace = true
		return 0, true
	}
	if arg == "--no-vscode-sync" {
		opts.NoVSCodeSync = true
		return 0, true
	}
	if arg == "--github-desktop" || arg == "--desktop" {
		opts.GHDesktop = true
		return 0, true
	}
	return parseMultiCloneValueFlag(args, idx, opts)
}

func parseMultiCloneValueFlag(args []string, idx int, opts *MultiCloneOptions) (int, bool) {
	arg := args[idx]
	hasValue := idx+1 < len(args)

	if (arg == "-d" || arg == "--dir" || arg == "--target") && hasValue {
		opts.TargetDir = args[idx+1]
		return 1, true
	}
	if (arg == "-f" || arg == "--file") && hasValue {
		opts.FilePath = args[idx+1]
		return 1, true
	}
	if (arg == "-j" || arg == "-p" || arg == "--concurrency" || arg == "--parallel") && hasValue {
		opts.Concurrency = parseConcurrencyValue(args[idx+1])
		opts.IsConcurrent = opts.Concurrency > 1
		return 1, true
	}
	return 0, false
}

func parseConcurrencyValue(val string) int {
	num, err := strconv.Atoi(val)
	if err == nil && num > 0 {
		return num
	}
	return 1
}

func handleMultiClonePositional(arg string, opts *MultiCloneOptions, textParts *[]string) {
	isExistingFile := checkIsFile(arg)
	if isExistingFile && opts.FilePath == "" {
		opts.FilePath = arg
		return
	}

	isFolderCandidate := opts.TargetDir == "" && len(*textParts) == 0 && isDirCandidate(arg)
	if isFolderCandidate {
		opts.TargetDir = arg
		return
	}

	*textParts = append(*textParts, arg)
}

func checkIsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isDirCandidate(arg string) bool {
	isURL := strings.Contains(arg, "://") || strings.HasPrefix(arg, "git@")
	isFence := strings.HasPrefix(arg, "```") || strings.HasPrefix(arg, "~~~")
	isOwnerRepo := strings.Contains(arg, "/") && !strings.HasSuffix(arg, "/") && !strings.Contains(arg, "\\")

	return !isURL && !isFence && !isOwnerRepo
}

func assembleRawInput(opts MultiCloneOptions, parts []string) string {
	if opts.FilePath != "" {
		content, err := os.ReadFile(opts.FilePath)
		if err == nil {
			return string(content)
		}
	}

	if len(parts) > 0 {
		return strings.Join(parts, "\n")
	}

	return readStdinIfAvailable()
}

func readStdinIfAvailable() string {
	stat, err := os.Stdin.Stat()
	isPiped := err == nil && (stat.Mode()&os.ModeCharDevice) == 0
	if !isPiped {
		return ""
	}

	bytes, readErr := io.ReadAll(os.Stdin)
	if readErr != nil {
		return ""
	}

	return string(bytes)
}
