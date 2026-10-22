// Package gitutil provides Git operations and metadata extraction.
package gitutil

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// EnableGitRerere enables git rerere.enabled and rerere.autoupdate.
func EnableGitRerere(repoDir string, isGlobal bool) error {
	argsEnabled := []string{constants.GitConfigCmd, "rerere.enabled", "true"}
	argsAutoUpdate := []string{constants.GitConfigCmd, "rerere.autoupdate", "true"}

	if isGlobal {
		argsEnabled = []string{constants.GitConfigCmd, "--global", "rerere.enabled", "true"}
		argsAutoUpdate = []string{constants.GitConfigCmd, "--global", "rerere.autoupdate", "true"}
	}

	if _, err := runGit(repoDir, argsEnabled...); err != nil {
		return apperror.WrapSimple(err, "failed to enable rerere")
	}

	if _, err := runGit(repoDir, argsAutoUpdate...); err != nil {
		return apperror.WrapSimple(err, "failed to enable rerere.autoupdate")
	}

	return nil
}

// CheckGitRerereEnabled reports whether git rerere is enabled.
func CheckGitRerereEnabled(repoDir string) bool {
	out, err := runGit(repoDir, constants.GitConfigCmd, constants.GitGetFlag, "rerere.enabled")
	if err != nil {
		return false
	}

	return strings.EqualFold(strings.TrimSpace(out), "true")
}

// CheckoutConflictSide checks out either the "theirs" or "ours" side for conflicted files and stages them.
func CheckoutConflictSide(repoDir string, side string, files []string) error {
	if len(files) == 0 {
		return nil
	}

	normalizedSide := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(side)), "--")
	if normalizedSide != "theirs" && normalizedSide != "ours" {
		return apperror.NewValidationError(fmt.Sprintf("invalid conflict side %q: expected 'theirs' or 'ours'", side))
	}

	checkoutArgs := append([]string{constants.GitCheckout, "--" + normalizedSide, "--"}, files...)
	if _, err := runGit(repoDir, checkoutArgs...); err != nil {
		return apperror.WrapSimple(err, "git checkout --"+normalizedSide+" failed")
	}

	addArgs := append([]string{constants.GitAdd, "--"}, files...)
	if _, err := runGit(repoDir, addArgs...); err != nil {
		return apperror.WrapSimple(err, "git add failed")
	}

	return nil
}

// SafeAbortMergeOrRebase safely aborts any active merge or rebase in repoDir.
func SafeAbortMergeOrRebase(repoDir string) error {
	if _, err := runGit(repoDir, "merge", "--abort"); err == nil {
		return nil
	}

	if _, err := runGit(repoDir, "rebase", "--abort"); err == nil {
		return nil
	}

	if _, err := runGit(repoDir, "reset", "--merge"); err == nil {
		return nil
	}

	return apperror.NewExecutionError(fmt.Sprintf("failed to abort merge or rebase in %s", repoDir))
}
