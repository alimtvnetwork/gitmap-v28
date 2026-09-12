// Package cmd — fixgit_safedir.go: safe.directory registration for ownership errors.
package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func remediateSafeDir(repoRoot string, opts FixGitOptions) ([]FixGitIssue, error) {
	var issues []FixGitIssue

	hasDubious := checkDubiousOwnership(repoRoot)
	if !hasDubious {
		return issues, nil
	}

	absPath, absErr := filepath.Abs(repoRoot)
	if absErr != nil {
		absPath = repoRoot
	}

	issue := FixGitIssue{
		Category:    "SafeDirectory",
		Description: "Dubious ownership detected by Git security policy",
	}

	if opts.IsDryRun {
		issue.Remedy = fmt.Sprintf("Would register %s in git config --global safe.directory", absPath)
		issues = append(issues, issue)

		return issues, nil
	}

	regErr := registerSafeDirectory(absPath)
	if regErr != nil {
		issue.ErrorDetail = regErr.Error()
		issue.Remedy = "Failed to register safe.directory in global config"
		issues = append(issues, issue)

		return issues, regErr
	}

	issue.IsFixed = true
	issue.Remedy = fmt.Sprintf("Registered %s in git config --global safe.directory", absPath)
	issues = append(issues, issue)

	return issues, nil
}

func checkDubiousOwnership(repoRoot string) bool {
	cmd := exec.Command("git", "status")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()

	if err == nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(out)), "detected dubious ownership")
}

func registerSafeDirectory(absPath string) error {
	normalized := filepath.ToSlash(absPath)
	cmd := exec.Command("git", "config", "--global", "--add", "safe.directory", normalized)

	return cmd.Run()
}
