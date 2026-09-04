// Package cmd — fixgit_index.go: index corruption detection and rebuild.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func remediateGitIndex(repoRoot, gitDir string, opts FixGitOptions) ([]FixGitIssue, error) {
	var issues []FixGitIssue

	hasCorrupt, reason := inspectIndexCorrupt(repoRoot, gitDir)
	if !hasCorrupt {
		return issues, nil
	}

	issue := FixGitIssue{
		Category:    "Index",
		Description: fmt.Sprintf("Corrupt or invalid index detected: %s", reason),
	}

	if opts.IsDryRun {
		issue.Remedy = "Would back up broken index and rebuild via 'git reset'"
		issues = append(issues, issue)

		return issues, nil
	}

	rebuildErr := executeIndexRebuild(repoRoot, gitDir)
	if rebuildErr != nil {
		issue.ErrorDetail = rebuildErr.Error()
		issue.Remedy = "Index rebuild failed"
		issues = append(issues, issue)

		return issues, rebuildErr
	}

	issue.IsFixed = true
	issue.Remedy = "Backed up damaged index and restored healthy index from HEAD"
	issues = append(issues, issue)

	return issues, nil
}

func inspectIndexCorrupt(repoRoot, gitDir string) (bool, string) {
	indexPath := filepath.Join(gitDir, "index")

	info, err := os.Stat(indexPath)
	if err == nil && info.Size() == 0 {
		return true, "index file is 0 bytes (truncated)"
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoRoot
	out, cmdErr := cmd.CombinedOutput()

	if cmdErr != nil && isIndexCorruptOutput(string(out)) {
		return true, strings.TrimSpace(string(out))
	}

	return false, ""
}

func isIndexCorruptOutput(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "index file corrupt") ||
		strings.Contains(lower, "bad signature") ||
		strings.Contains(lower, "could not write new index file") ||
		strings.Contains(lower, "smaller than minimum")
}

func executeIndexRebuild(repoRoot, gitDir string) error {
	indexPath := filepath.Join(gitDir, "index")

	backupErr := backupCorruptIndex(gitDir, indexPath)
	if backupErr != nil {
		return backupErr
	}

	_ = os.Remove(indexPath)

	cmd := exec.Command("git", "reset", "HEAD")
	cmd.Dir = repoRoot

	return cmd.Run()
}

func backupCorruptIndex(gitDir, indexPath string) error {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
	}

	backupName := fmt.Sprintf("index.corrupt.%d", time.Now().Unix())
	backupPath := filepath.Join(gitDir, backupName)

	data, readErr := os.ReadFile(indexPath)
	if readErr != nil {
		return readErr
	}

	return os.WriteFile(backupPath, data, 0644)
}
