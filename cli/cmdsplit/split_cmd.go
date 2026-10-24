package cmdsplit

import (
	"flag"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

type splitFlags struct {
	dryRun      bool
	autoConfirm bool
	force       bool
	sizeRaw     string
	threshold   int64
}

type largeFileInfo struct {
	RelPath           string
	Size              int64
	Method            uint16
	MethodName        string
	ProjectedChunks   int
	ProjectedPartList []string
}

// RunSplitLargeFilesCmd executes the split-large-files command.
func RunSplitLargeFilesCmd(args []string) error {
	cleanArgs := normalizeSplitArgs(args)
	if hasHelpFlag(cleanArgs) {
		helpdoc.Print("split-large-files")
		return nil
	}

	flags, err := parseSplitFlags(cleanArgs)
	if err != nil {
		return err
	}

	if !isInsideGitRepo() {
		return apperror.NewValidation("split_large_files", "E9001", "not inside a Git repository")
	}

	return processSplitOperation(flags)
}

var splitHelpFlagTokens = []string{"-h", "--help", "help"}

func normalizeSplitArgs(args []string) []string {
	if len(args) >= 2 && args[0] == "large" && args[1] == "files" {
		return args[2:]
	}

	if len(args) >= 1 && slices.Contains(splitLargeFilesTokens, args[0]) {
		return args[1:]
	}

	return args
}

var splitLargeFilesTokens = []string{"large-files", "large_files"}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if slices.Contains(splitHelpFlagTokens, arg) {
			return true
		}
	}

	return false
}

func parseSplitFlags(args []string) (splitFlags, error) {
	var sf splitFlags
	fs := flag.NewFlagSet("split-large-files", flag.ContinueOnError)

	fs.BoolVar(&sf.dryRun, "dry-run", false, "Preview without modifying files or git")
	fs.BoolVar(&sf.autoConfirm, "y", false, "Auto-confirm execution prompt")
	fs.BoolVar(&sf.autoConfirm, "yes", false, "Auto-confirm execution prompt")
	fs.BoolVar(&sf.force, "force", false, "Force split execution")
	fs.StringVar(&sf.sizeRaw, "size", "10MB", "Size threshold for splitting")
	fs.StringVar(&sf.sizeRaw, "s", "10MB", "Size threshold for splitting")

	if err := fs.Parse(args); err != nil {
		return sf, apperror.WrapSimple(err, "parse_split_flags")
	}

	threshold, err := parseByteSize(sf.sizeRaw)
	sf.threshold = threshold

	return sf, err
}

func parseByteSize(val string) (int64, error) {
	trimmed := strings.ToUpper(strings.TrimSpace(val))
	if trimmed == "" {
		return 10 * 1024 * 1024, nil
	}

	multiplier, numStr := resolveSizeMultiplier(trimmed)
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, apperror.WrapSimple(err, "invalid_size_format")
	}

	return int64(num * float64(multiplier)), nil
}

func resolveSizeMultiplier(s string) (int64, string) {
	if strings.HasSuffix(s, "GB") {
		return 1024 * 1024 * 1024, strings.TrimSuffix(s, "GB")
	}
	if strings.HasSuffix(s, "G") {
		return 1024 * 1024 * 1024, strings.TrimSuffix(s, "G")
	}
	if strings.HasSuffix(s, "MB") {
		return 1024 * 1024 * 1024, strings.TrimSuffix(s, "MB")
	}
	if strings.HasSuffix(s, "M") {
		return 1024 * 1024 * 1024, strings.TrimSuffix(s, "M")
	}
	if strings.HasSuffix(s, "KB") {
		return 1024, strings.TrimSuffix(s, "KB")
	}
	if strings.HasSuffix(s, "K") {
		return 1024, strings.TrimSuffix(s, "K")
	}

	return 1, strings.TrimSuffix(s, "B")
}

func formatByteSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func isInsideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}
