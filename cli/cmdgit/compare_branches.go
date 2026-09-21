package cmdgit

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// RunCompare orchestrates comparing 2 or more branches via external diff tool.
func RunCompare(args []string) error {
	branches, err := parseCompareArgs(args)
	if err != nil {
		return err
	}

	repoRoot, err := findGitRoot()
	if err != nil {
		return err
	}

	return executeCompare(repoRoot, branches)
}

func parseCompareArgs(args []string) ([]string, error) {
	if len(args) < 2 {
		return nil, apperror.NewValidationError("Usage: gitmap compare (cmp) <branch1> <branch2> [branch3...]")
	}

	return args, nil
}

func findGitRoot() (string, error) {
	root, err := runGitOutput(".", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", apperror.WrapSimple(err, "not a git repository")
	}

	return root, nil
}

func executeCompare(repoRoot string, branches []string) error {
	if err := validateAllBranches(repoRoot, branches); err != nil {
		return err
	}

	var worktrees []string
	defer func() { cleanupWorktrees(repoRoot, worktrees) }()

	dirs, wts, err := prepareBranchDirectories(repoRoot, branches)
	if err != nil {
		return err
	}
	worktrees = wts

	return launchComparison(dirs)
}

func validateAllBranches(repoRoot string, branches []string) error {
	for _, branch := range branches {
		if err := verifyBranchRef(repoRoot, branch); err != nil {
			return err
		}
	}

	return nil
}

func verifyBranchRef(repoRoot, branch string) error {
	isValid := hasValidBranchRef(repoRoot, branch)
	if isValid {
		return nil
	}

	return apperror.NewValidationError(fmt.Sprintf("invalid branch or git ref: %s", branch))
}

func hasValidBranchRef(repoRoot, branch string) bool {
	_, err := runGitOutput(repoRoot, "rev-parse", "--verify", branch+"^{commit}")

	return err == nil
}

func prepareBranchDirectories(repoRoot string, branches []string) ([]string, []string, error) {
	activeBranch := gitutil.GetActiveBranch(repoRoot)
	dirs := make([]string, 0, len(branches))
	var createdWorktrees []string

	for _, branch := range branches {
		dir, isWorktree, err := resolveBranchPath(repoRoot, branch, activeBranch)
		if err != nil {
			return nil, createdWorktrees, err
		}
		dirs = append(dirs, dir)
		if isWorktree {
			createdWorktrees = append(createdWorktrees, dir)
		}
	}

	return dirs, createdWorktrees, nil
}

func resolveBranchPath(repoRoot, branch, activeBranch string) (string, bool, error) {
	if branch == activeBranch {
		return repoRoot, false, nil
	}

	wtPath, err := createCompareWorktree(repoRoot, branch)
	if err != nil {
		return "", false, err
	}

	return wtPath, true, nil
}

func launchComparison(dirs []string) error {
	tool, err := findCompareTool()
	if err != nil {
		return err
	}

	return launchToolWithPaths(tool, dirs)
}
