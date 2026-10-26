package cmddiff

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func runSinglePathDiff(target string) error {
	repoRoot, isGitRepo := detectTargetGitRepo(target)
	if !isGitRepo {
		return apperror.NewExecutionError("not a git repository: " + target)
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

	dirs := []string{diffTarget}

	return cmdgit.LaunchToolWithPathsAsync(tool, dirs)
}

func renderTerminalGitDiff(repoRoot, target string) error {
	args := []string{"diff", "--color=always", "HEAD"}
	if target != "" && target != "." {
		args = append(args, "--", target)
	}

	fmt.Printf("Displaying git diff against HEAD for %s:\n", target)

	return cmdgit.ExecGitInspection(repoRoot, args...)
}
