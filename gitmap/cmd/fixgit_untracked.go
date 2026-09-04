// Package cmd — fixgit_untracked.go: backup and remediation of untracked files blocking pull.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func remediateUntrackedOverwrites(repoRoot string, opts FixGitOptions) ([]FixGitIssue, error) {
	var issues []FixGitIssue

	hasCollision, files := probeUntrackedCollisions(repoRoot)
	if !hasCollision || len(files) == 0 {
		return issues, nil
	}

	issue := FixGitIssue{
		Category:    "Untracked",
		Description: fmt.Sprintf("%d untracked file(s) blocking incoming merge", len(files)),
	}

	if opts.IsDryRun {
		issue.Remedy = fmt.Sprintf("Would back up %d file(s) to .gitmap/backup/ and pull", len(files))
		issues = append(issues, issue)

		return issues, nil
	}

	backupDir, pullErr := backupAndPullRepo(repoRoot, files)
	if pullErr != nil {
		issue.ErrorDetail = pullErr.Error()
		issue.Remedy = "Failed to complete pull after backup"
		issues = append(issues, issue)

		return issues, pullErr
	}

	relBackup, _ := filepath.Rel(repoRoot, backupDir)
	issue.IsFixed = true
	issue.Remedy = fmt.Sprintf("Backed up %d file(s) to %s and completed git pull", len(files), relBackup)
	issues = append(issues, issue)

	return issues, nil
}

func probeUntrackedCollisions(repoRoot string) (bool, []string) {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()

	if err == nil {
		return false, nil
	}

	outputStr := string(out)
	if !strings.Contains(outputStr, "untracked working tree files would be overwritten") {
		return false, nil
	}

	files := parseUntrackedConflictFiles(outputStr)

	return len(files) > 0, files
}

func parseUntrackedConflictFiles(output string) []string {
	var files []string
	isCollecting := false
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "untracked working tree files would be overwritten") {
			isCollecting = true
			continue
		}

		if isCollecting && strings.HasPrefix(trimmed, "Please move or remove") {
			break
		}

		if isCollecting && trimmed != "" && !strings.HasPrefix(trimmed, "Aborting") {
			files = append(files, trimmed)
		}
	}

	return files
}

func backupAndPullRepo(repoRoot string, files []string) (string, error) {
	backupDir := filepath.Join(repoRoot, ".gitmap", "backup", "untracked-overwrites", fmt.Sprintf("%d", time.Now().Unix()))

	backupErr := backupConflictFiles(repoRoot, backupDir, files)
	if backupErr != nil {
		return "", backupErr
	}

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = repoRoot
	out, pullErr := cmd.CombinedOutput()

	if pullErr != nil {
		return backupDir, fmt.Errorf("git pull failed: %s: %w", strings.TrimSpace(string(out)), pullErr)
	}

	return backupDir, nil
}

func backupConflictFiles(repoRoot, backupDir string, files []string) error {
	for _, relFile := range files {
		src := filepath.Join(repoRoot, filepath.FromSlash(relFile))
		dst := filepath.Join(backupDir, filepath.FromSlash(relFile))

		copyErr := copyAndRemoveFile(src, dst)
		if copyErr != nil {
			return copyErr
		}
	}

	return nil
}

func copyAndRemoveFile(src, dst string) error {
	data, readErr := os.ReadFile(src)
	if readErr != nil {
		return nil
	}

	dirErr := os.MkdirAll(filepath.Dir(dst), 0755)
	if dirErr != nil {
		return dirErr
	}

	writeErr := os.WriteFile(dst, data, 0644)
	if writeErr != nil {
		return writeErr
	}

	return os.Remove(src)
}
