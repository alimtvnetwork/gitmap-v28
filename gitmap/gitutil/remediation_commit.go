// Package gitutil — remediation_commit.go generates commit workflow commands.
package gitutil

import "fmt"

func GenerateCommitRecipe(repoPath string) RemediationRecipe {
	p := CleanRepoPath(repoPath)
	return RemediationRecipe{
		Title:       "Option 2 (Commit Work-In-Progress)",
		Description: "Commit all modified and untracked files locally before pulling",
		Command:     fmt.Sprintf("git -C %s add -A && git -C %s commit -m \"wip: local changes\" && git -C %s pull --rebase", p, p, p),
	}
}
