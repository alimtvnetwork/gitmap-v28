// Package gitutil — remediation_stash.go generates stash/pop workflow commands.
package gitutil

import "fmt"

func GenerateStashRecipe(repoPath string) RemediationRecipe {
	return RemediationRecipe{
		Title:       "Option 1 (Stash & Re-apply)",
		Description: "Temporarily save local changes, pull latest remote commits, then re-apply",
		Command:     fmt.Sprintf("git -C %s stash && git -C %s pull && git -C %s stash pop", repoPath, repoPath, repoPath),
	}
}
