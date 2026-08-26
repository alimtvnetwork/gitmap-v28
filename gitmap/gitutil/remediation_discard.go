// Package gitutil — remediation_discard.go generates discard/clean workflow commands.
package gitutil

import "fmt"

func GenerateDiscardRecipe(repoPath string) RemediationRecipe {
	return RemediationRecipe{
		Title:       "Option 3 (Discard Local Changes)",
		Description: "Permanently discard local modifications and untracked files to match remote",
		Command:     fmt.Sprintf("git -C %s reset --hard HEAD && git -C %s clean -fd && git -C %s pull", repoPath, repoPath, repoPath),
	}
}
