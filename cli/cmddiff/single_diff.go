package cmddiff

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func runSinglePathDiff(target string) error {
	repoRoot, isGitRepo := detectTargetGitRepo(target)
	if !isGitRepo {
		fmt.Fprintf(os.Stderr, "  ✗ '%s' is not a Git repository.\n", target)
		fmt.Fprintln(os.Stderr, "  💡 Tip: Specify comparison paths (gitmap diff <folder1> <folder2>) or run inside a Git repository.")

		return nil
	}

	compareTool, err := cmdgit.FindCompareTool()
	if err == nil && compareTool.Path != "" {
		return launchExternalToolDiff(compareTool, repoRoot, target)
	}

	return renderTerminalGitDiff(repoRoot, target)
}

func detectTargetGitRepo(target string) (string, bool) {
	evalDir := resolveTargetDir(target)
	res := gitutil.RunGitInDir(evalDir, "rev-parse", "--show-toplevel")
	if res.IsFailure() {
		return "", false
	}

	repoRoot := strings.TrimSpace(string(res.Data))

	return repoRoot, repoRoot != ""
}

func resolveTargetDir(target string) string {
	if target == "" || target == "." {
		return "."
	}

	info, err := os.Stat(target)
	if err == nil && !info.IsDir() {
		return filepath.Dir(target)
	}

	return target
}

func launchExternalToolDiff(tool cmdgit.CompareTool, repoRoot, target string) error {
	diffTarget := target
	if diffTarget == "" {
		diffTarget = repoRoot
	}

	if isWorktreeClean(repoRoot, diffTarget) {
		fmt.Printf("No changes in %s\n", diffTarget)

		return nil
	}

	dirs := []string{diffTarget}

	return cmdgit.LaunchToolWithPathsAsync(tool, dirs)
}

func isWorktreeClean(repoRoot, target string) bool {
	if isDirectoryTarget(target) {
		return false
	}

	res := gitutil.RunGitInDir(repoRoot, "status", "--porcelain", "--", target)
	if res.IsFailure() {
		return false
	}

	return strings.TrimSpace(string(res.Data)) == ""
}

func isDirectoryTarget(target string) bool {
	info, err := os.Stat(target)
	if err != nil {
		return true
	}

	return info.IsDir()
}

func renderTerminalGitDiff(repoRoot, target string) error {
	args := []string{"diff", "--color=always", "HEAD"}
	if target != "" && target != "." {
		args = append(args, "--", target)
	}

	fmt.Printf("Displaying git diff against HEAD for %s:\n", target)

	return cmdgit.ExecGitInspection(repoRoot, args...)
}
