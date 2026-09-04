// Package cmd — fixgit_unmerged.go: detection and resolution of stalled merge conflicts.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func remediateUnmergedConflict(repoRoot string, opts FixGitOptions) ([]FixGitIssue, error) {
	var issues []FixGitIssue

	hasConflict, conflictFiles := detectUnmergedConflict(repoRoot)
	if !hasConflict {
		return issues, nil
	}

	issue := FixGitIssue{
		Category:    "Unmerged",
		Description: fmt.Sprintf("Unresolved merge conflict in %d file(s)", len(conflictFiles)),
	}

	if opts.IsDryRun {
		issue.Remedy = "Would run 'git merge --abort' to clear stalled conflict state"
		issues = append(issues, issue)

		return issues, nil
	}

	abortErr := executeMergeAbort(repoRoot)
	if abortErr != nil {
		issue.ErrorDetail = abortErr.Error()
		issue.Remedy = "Failed to abort stalled merge"
		issues = append(issues, issue)

		return issues, abortErr
	}

	issue.IsFixed = true
	issue.Remedy = "Aborted stalled merge conflict; restored clean working tree"
	issues = append(issues, issue)

	return issues, nil
}

func detectUnmergedConflict(repoRoot string) (bool, []string) {
	mergeHeadPath := filepath.Join(repoRoot, ".git", "MERGE_HEAD")
	if _, err := os.Stat(mergeHeadPath); err == nil {
		return true, []string{"MERGE_HEAD"}
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()

	if err != nil {
		return false, nil
	}

	files := parseUnmergedFiles(string(out))

	return len(files) > 0, files
}

func parseUnmergedFiles(statusOutput string) []string {
	var unmerged []string
	lines := strings.Split(statusOutput, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		file, hasFile := extractUnmergedFilePath(trimmed)
		if hasFile {
			unmerged = append(unmerged, file)
		}
	}

	return unmerged
}

func extractUnmergedFilePath(line string) (string, bool) {
	if !isUnmergedStatusCode(line) {
		return "", false
	}

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return "", false
	}

	return parts[1], true
}

func isUnmergedStatusCode(line string) bool {
	if len(line) < 2 {
		return false
	}

	prefix := line[:2]

	return prefix == "UU" || prefix == "AA" || prefix == "UD" || prefix == "DU"
}

func executeMergeAbort(repoRoot string) error {
	cmd := exec.Command("git", "merge", "--abort")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()

	if err == nil {
		return nil
	}

	if strings.Contains(string(out), "no merge to abort") {
		resetCmd := exec.Command("git", "reset", "--merge")
		resetCmd.Dir = repoRoot

		return resetCmd.Run()
	}

	return err
}
