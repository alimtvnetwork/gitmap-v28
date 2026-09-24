// Package cmd — lowercasefix_commit.go handles git staging and committing of renames.
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func commitRenames(pairs []RenamePair, customMsg string, isNoPush bool) (string, bool, error) {
	if len(pairs) == 0 {
		return "", false, nil
	}

	if err := stageRenamedFiles(); err != nil {
		return "", false, err
	}

	msg := resolveCommitMsg(pairs, customMsg)
	if err := executeGitCommit(msg); err != nil {
		return "", false, err
	}

	sha := getLatestCommitSHA()
	shortSHA := sha
	if len(sha) > 8 {
		shortSHA = sha[:8]
	}

	fmt.Printf("\n%s✓ Committed %d lowercase file rename(s): %q (%s)%s\n",
		constants.ColorGreen, len(pairs), msg, shortSHA, constants.ColorReset)

	if isNoPush {
		fmt.Println()
		return sha, false, nil
	}

	pushErr := executeGitPush()
	if pushErr != nil {
		fmt.Printf("%s⚠️  Git commit created (%s), but git push failed: %v%s\n\n",
			constants.ColorYellow, shortSHA, pushErr, constants.ColorReset)
		return sha, false, nil
	}

	fmt.Printf("%s✓ Pushed commit (%s) to remote tracking branch%s\n\n",
		constants.ColorGreen, shortSHA, constants.ColorReset)
	return sha, true, nil
}

func executeGitPush() error {
	branch := resolveCurrentBranchName()
	cmd := exec.Command("git", "push", "origin", branch)
	if _, err := cmd.CombinedOutput(); err == nil {
		return nil
	}

	cmdSimple := exec.Command("git", "push")
	if outSimple, errSimple := cmdSimple.CombinedOutput(); errSimple != nil {
		return fmt.Errorf("push failed: %w (%s)", errSimple, strings.TrimSpace(string(outSimple)))
	}

	return nil
}

func stageRenamedFiles() error {
	cmd := exec.Command("git", "add", "-A")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add -A failed: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	return nil
}

func executeGitCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	return nil
}

func getLatestCommitSHA() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func resolveCommitMsg(pairs []RenamePair, customMsg string) string {
	if strings.TrimSpace(customMsg) != "" {
		return customMsg
	}
	if len(pairs) == 1 {
		return fmt.Sprintf("chore: rename %s to lowercase %s", pairs[0].OldBase, pairs[0].NewBase)
	}

	return fmt.Sprintf("chore: rename %d files to lowercase across repository", len(pairs))
}
