// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"fmt"
	"strings"
)

// RunBackupBranch implements `gitmap space backup-branch "<task>"`:
// it creates backup/<slug> from current HEAD and pushes it unless
// --no-push is passed. The working tree must be clean (spec 243.2).
func RunBackupBranch(args []string) error {
	if hasHelpFlag(args) {
		printBackupBranchHelp()
		return nil
	}

	task, noPush, force := parseBackupBranchArgs(args)
	if task == "" {
		return fmt.Errorf(`missing task string: usage: gitmap space backup-branch "<task string>" [--no-push] [--force]`)
	}

	slug, slugErr := makeBackupSlug(task)
	if slugErr != nil {
		return slugErr
	}
	branch := "backup/" + slug

	if cleanErr := requireCleanTree(); cleanErr != nil {
		return cleanErr
	}

	if existsErr := ensureBranchAbsence(branch, force); existsErr != nil {
		return existsErr
	}

	if createErr := createBackupBranch(branch, force); createErr != nil {
		return createErr
	}

	sha, shaErr := runGitCapture("rev-parse", "--short", "HEAD")
	if shaErr != nil {
		return shaErr
	}
	fmt.Printf("%s @ %s\n", branch, sha)

	if !noPush {
		if pushErr := runGitInherit("push", "-u", "origin", branch); pushErr != nil {
			return pushErr
		}
	}

	return nil
}

// hasHelpFlag reports whether any arg requests help text.
func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}

	return false
}

// parseBackupBranchArgs splits args into the task string (positional args
// joined with a space) plus the --no-push and --force flags, position
// agnostic, matching the space common flag style.
func parseBackupBranchArgs(args []string) (task string, noPush, force bool) {
	var positional []string

	for _, a := range args {
		switch a {
		case "--no-push":
			noPush = true
		case "--force", "-f":
			force = true
		default:
			if !strings.HasPrefix(a, "-") {
				positional = append(positional, a)
			}
		}
	}

	return strings.Join(positional, " "), noPush, force
}
