// Package cmd — lowercasefix_git.go handles Git detection and safe two-step renaming.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()

	return err == nil && strings.TrimSpace(string(out)) == "true"
}

func isFileTrackedInGit(path string) bool {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", path)

	return cmd.Run() == nil
}

func performSafeGitRename(p RenamePair) error {
	tmpPath := p.OldPath + ".tmp-lcf"
	if p.IsGitTracked {
		return performTwoStepGitMv(p.OldPath, tmpPath, p.NewPath)
	}

	return performTwoStepFSRename(p.OldPath, tmpPath, p.NewPath)
}

func performTwoStepGitMv(src, tmp, dst string) error {
	_ = os.Remove(tmp)
	if err := runGitMvForLcf(src, tmp); err != nil {
		return fallbackRename(src, tmp, dst, err)
	}
	if err := runGitMvForLcf(tmp, dst); err != nil {
		_ = runGitMvForLcf(tmp, src)

		return fmt.Errorf("git mv step 2 (%s -> %s) failed: %w", tmp, dst, err)
	}

	return nil
}

func performTwoStepFSRename(src, tmp, dst string) error {
	_ = os.Remove(tmp)
	if err := os.Rename(src, tmp); err != nil {
		return fmt.Errorf("fs rename step 1 (%s -> %s) failed: %w", src, tmp, err)
	}
	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Rename(tmp, src)

		return fmt.Errorf("fs rename step 2 (%s -> %s) failed: %w", tmp, dst, err)
	}
	if isGitRepository() {
		_ = exec.Command("git", "add", dst).Run()
	}

	return nil
}

func fallbackRename(src, tmp, dst string, gitErr error) error {
	_ = os.Remove(tmp)
	if err := os.Rename(src, tmp); err != nil {
		return fmt.Errorf("git mv failed (%v) and fs rename to tmp failed: %w", gitErr, err)
	}
	if err := os.Rename(tmp, dst); err != nil {
		return fmt.Errorf("fs rename tmp to %s failed: %w", dst, err)
	}
	if isGitRepository() {
		_ = exec.Command("git", "add", dst).Run()
	}

	return nil
}

func runGitMvForLcf(src, dst string) error {
	cmd := exec.Command("git", "mv", src, dst)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git mv %s to %s failed: %w (output: %s)", src, dst, err, strings.TrimSpace(string(out)))
	}

	return nil
}

// GitWorkingTreeStatus records pending dirty or conflicted files.
type GitWorkingTreeStatus struct {
	HasConflicts  bool
	ConflictFiles []string
	DirtyFiles    []string
}

func checkGitWorkingTreeStatus() GitWorkingTreeStatus {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return GitWorkingTreeStatus{}
	}

	var status GitWorkingTreeStatus
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		if len(trimmed) < 3 {
			continue
		}
		code := trimmed[:2]
		filePath := strings.TrimSpace(trimmed[3:])
		if isGitConflictCode(code) {
			status.HasConflicts = true
			status.ConflictFiles = append(status.ConflictFiles, filePath)
		} else {
			status.DirtyFiles = append(status.DirtyFiles, filePath)
		}
	}

	return status
}

func isGitConflictCode(code string) bool {
	return code == "UU" || code == "AA" || code == "DD" ||
		code == "AU" || code == "UA" || code == "UD" || code == "DU"
}

func discardGitWorkingTreeChanges() error {
	cmdReset := exec.Command("git", "reset", "--hard", "HEAD")
	if out, err := cmdReset.CombinedOutput(); err != nil {
		return fmt.Errorf("git reset --hard failed: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	cmdClean := exec.Command("git", "clean", "-fd")
	if out, err := cmdClean.CombinedOutput(); err != nil {
		return fmt.Errorf("git clean -fd failed: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	return nil
}
