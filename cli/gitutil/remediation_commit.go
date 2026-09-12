// Package gitutil — remediation_commit.go generates commit workflow commands.
package gitutil

import "fmt"

func GenerateCommitRecipe(repoPath string) RemediationRecipe {
	p := CleanRepoPath(repoPath)
	rawPath := CleanRepoPathRaw(repoPath)

	return RemediationRecipe{
		Title:       "Option 2 (Commit Work-In-Progress)",
		Description: "Commit all modified and untracked files locally before pulling",
		Command:     fmt.Sprintf("git -C %s add -A && git -C %s commit -m \"wip: local changes\" && git -C %s pull --rebase", p, p, p),
		Steps: []RemediationStep{
			{Name: "git", Args: []string{"-C", rawPath, "add", "-A"}},
			{Name: "git", Args: []string{"-C", rawPath, "commit", "-m", "wip: local changes"}},
			{Name: "git", Args: []string{"-C", rawPath, "pull", "--rebase"}},
		},
	}
}
