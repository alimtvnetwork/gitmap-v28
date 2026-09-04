// Package cmd — fixgit.go: CLI entry point for gitmap fix-git.
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runFixGit(args []string) error {
	checkHelp(constants.CmdFixGit, args)

	opts := parseFixGitOptions(args)

	repoRoot, gitDir, findErr := resolveGitDirectory(opts.TargetDir)
	if findErr != nil {
		return findErr
	}

	result, execErr := executeFixGitPipelines(repoRoot, gitDir, opts)
	if execErr != nil {
		return execErr
	}

	return renderFixGitOutput(result, opts.IsJSON)
}

func parseFixGitOptions(args []string) FixGitOptions {
	opts := FixGitOptions{TargetDir: "."}

	for _, arg := range args {
		switch {
		case arg == "--dry-run" || arg == "-n":
			opts.IsDryRun = true
		case arg == "--json":
			opts.IsJSON = true
		case arg == "--verbose" || arg == "-v":
			opts.IsVerbose = true
		case arg == "--permissions" || arg == "--perms" || arg == "--permissions-only":
			opts.IsPermissionsOnly = true
		case arg == "--locks" || arg == "--locks-only":
			opts.IsLocksOnly = true
		case arg == "--index" || arg == "--index-only":
			opts.IsIndexOnly = true
		case arg == "--safe-dir":
			opts.IsSafeDirOnly = true
		case !strings.HasPrefix(arg, "-"):
			opts.TargetDir = arg
		}
	}

	return opts
}

func resolveGitDirectory(target string) (string, string, error) {
	absTarget, absErr := filepath.Abs(target)
	if absErr != nil {
		return "", "", apperror.WrapSimple(absErr, "resolve target path")
	}

	gitDir := filepath.Join(absTarget, ".git")

	info, statErr := os.Stat(gitDir)
	if statErr != nil || !info.IsDir() {
		return "", "", apperror.NewValidationError("not a git repository (missing .git directory)")
	}

	return absTarget, gitDir, nil
}

func executeFixGitPipelines(repoRoot, gitDir string, opts FixGitOptions) (FixGitResult, error) {
	var allIssues []FixGitIssue

	allIssues = appendFilteredIssues(allIssues, runPermsPipeline(gitDir, opts))
	allIssues = appendFilteredIssues(allIssues, runLocksPipeline(gitDir, opts))
	allIssues = appendFilteredIssues(allIssues, runIndexPipeline(repoRoot, gitDir, opts))
	allIssues = appendFilteredIssues(allIssues, runSafeDirPipeline(repoRoot, opts))
	allIssues = appendFilteredIssues(allIssues, runUnmergedPipeline(repoRoot, opts))
	allIssues = appendFilteredIssues(allIssues, runUntrackedPipeline(repoRoot, opts))

	res := summarizeFixGitResult(repoRoot, allIssues)

	return res, nil
}

func runPermsPipeline(gitDir string, opts FixGitOptions) []FixGitIssue {
	if opts.IsLocksOnly || opts.IsIndexOnly || opts.IsSafeDirOnly {
		return nil
	}

	issues, _ := remediateGitPermissions(gitDir, opts)

	return issues
}

func runLocksPipeline(gitDir string, opts FixGitOptions) []FixGitIssue {
	if opts.IsPermissionsOnly || opts.IsIndexOnly || opts.IsSafeDirOnly {
		return nil
	}

	issues, _ := remediateGitLocks(gitDir, opts)

	return issues
}

func runIndexPipeline(repoRoot, gitDir string, opts FixGitOptions) []FixGitIssue {
	if opts.IsPermissionsOnly || opts.IsLocksOnly || opts.IsSafeDirOnly {
		return nil
	}

	issues, _ := remediateGitIndex(repoRoot, gitDir, opts)

	return issues
}

func runSafeDirPipeline(repoRoot string, opts FixGitOptions) []FixGitIssue {
	if opts.IsPermissionsOnly || opts.IsLocksOnly || opts.IsIndexOnly {
		return nil
	}

	issues, _ := remediateSafeDir(repoRoot, opts)

	return issues
}

func runUnmergedPipeline(repoRoot string, opts FixGitOptions) []FixGitIssue {
	if opts.IsPermissionsOnly || opts.IsLocksOnly || opts.IsIndexOnly || opts.IsSafeDirOnly {
		return nil
	}

	issues, _ := remediateUnmergedConflict(repoRoot, opts)

	return issues
}

func runUntrackedPipeline(repoRoot string, opts FixGitOptions) []FixGitIssue {
	if opts.IsPermissionsOnly || opts.IsLocksOnly || opts.IsIndexOnly || opts.IsSafeDirOnly {
		return nil
	}

	issues, _ := remediateUntrackedOverwrites(repoRoot, opts)

	return issues
}

func appendFilteredIssues(target []FixGitIssue, add []FixGitIssue) []FixGitIssue {
	if len(add) == 0 {
		return target
	}

	return append(target, add...)
}

func summarizeFixGitResult(repoRoot string, issues []FixGitIssue) FixGitResult {
	fixedCount := 0

	for _, issue := range issues {
		if issue.IsFixed {
			fixedCount++
		}
	}

	isClean := len(issues) == 0 || fixedCount == len(issues)

	return FixGitResult{
		TargetDir:   repoRoot,
		IsClean:     isClean,
		IssuesFound: len(issues),
		IssuesFixed: fixedCount,
		Issues:      issues,
	}
}
