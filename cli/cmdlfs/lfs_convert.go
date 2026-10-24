package cmdlfs

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

type lfsFlags struct {
	dryRun      bool
	autoConfirm bool
	force       bool
	sizeRaw     string
	threshold   int64
}

type lfsCandidate struct {
	RelPath       string
	Size          int64
	Pattern       string
	AlreadyInAttr bool
}

// RunLFSConvertCmd executes the 'gitmap lfs convert' command.
func RunLFSConvertCmd(args []string) error {
	cleanArgs := normalizeLFSArgs(args)
	if hasHelpFlag(cleanArgs) {
		helpdoc.Print("lfs-convert")
		return nil
	}

	flags, err := parseLFSFlags(cleanArgs)
	if err != nil {
		return err
	}

	if !isInsideGitRepo() {
		return apperror.NewValidation("lfs_convert", "E9001", "not inside a Git repository")
	}

	return processLFSConvert(flags)
}

func normalizeLFSArgs(args []string) []string {
	if len(args) >= 1 && (args[0] == "convert" || args[0] == "conv") {
		return args[1:]
	}
	return args
}

var lfsHelpFlagTokens = []string{"-h", "--help", "help"}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if slices.Contains(lfsHelpFlagTokens, arg) {
			return true
		}
	}
	return false
}

func parseLFSFlags(args []string) (lfsFlags, error) {
	var lf lfsFlags
	fs := flag.NewFlagSet("lfs-convert", flag.ContinueOnError)
	fs.BoolVar(&lf.dryRun, "dry-run", false, "Preview without modifying files or git")
	fs.BoolVar(&lf.autoConfirm, "y", false, "Auto-confirm execution prompt")
	fs.BoolVar(&lf.autoConfirm, "yes", false, "Auto-confirm execution prompt")
	fs.BoolVar(&lf.force, "force", false, "Force convert execution")
	fs.StringVar(&lf.sizeRaw, "size", "5MB", "Size threshold for LFS conversion")
	fs.StringVar(&lf.sizeRaw, "s", "5MB", "Size threshold for LFS conversion")
	if err := fs.Parse(args); err != nil {
		return lf, apperror.WrapSimple(err, "parse_lfs_flags")
	}
	threshold, err := parseByteSize(lf.sizeRaw)
	lf.threshold = threshold
	return lf, err
}

func parseByteSize(val string) (int64, error) {
	trimmed := strings.ToUpper(strings.TrimSpace(val))
	if trimmed == "" {
		return 5 * 1024 * 1024, nil
	}

	multiplier, numStr := resolveSizeMultiplier(trimmed)
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, apperror.WrapSimple(err, "invalid_size_format")
	}

	return int64(num * float64(multiplier)), nil
}

func resolveSizeMultiplier(s string) (int64, string) {
	for _, u := range []struct {
		s string
		m int64
	}{
		{"GB", 1024 * 1024 * 1024}, {"G", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024}, {"M", 1024 * 1024},
		{"KB", 1024}, {"K", 1024}, {"B", 1},
	} {
		if strings.HasSuffix(s, u.s) {
			return u.m, strings.TrimSuffix(s, u.s)
		}
	}
	return 1, s
}

func formatByteSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(1024), 0
	for n := bytes / 1024; n >= 1024; n /= 1024 {
		div *= 1024
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

func isLFSAvailable() bool {
	cmd := exec.Command("git", "lfs", "version")
	return cmd.Run() == nil
}
