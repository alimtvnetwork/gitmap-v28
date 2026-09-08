// Package gitutil — remediation_stash.go generates stash/pop workflow commands.
package gitutil

import "fmt"

func GenerateStashRecipe(repoPath string) RemediationRecipe {
	p := CleanRepoPath(repoPath)
	rawPath := CleanRepoPathRaw(repoPath)
	return RemediationRecipe{
		Title:       "Option 1 (Stash & Re-apply)",
		Description: "Temporarily save local changes (including untracked), pull latest remote commits, then re-apply",
		Command:     fmt.Sprintf("git -C %s stash -u && git -C %s pull && git -C %s stash pop", p, p, p),
		Steps: []RemediationStep{
			{Name: "git", Args: []string{"-C", rawPath, "stash", "-u"}},
			{Name: "git", Args: []string{"-C", rawPath, "pull"}},
			{Name: "git", Args: []string{"-C", rawPath, "stash", "pop"}},
		},
	}
}
