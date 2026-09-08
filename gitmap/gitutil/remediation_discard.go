// Package gitutil — remediation_discard.go generates discard/clean workflow commands.
package gitutil

import "fmt"

func GenerateDiscardRecipe(repoPath string) RemediationRecipe {
	p := CleanRepoPath(repoPath)
	rawPath := CleanRepoPathRaw(repoPath)
	return RemediationRecipe{
		Title:       "Option 3 (Discard Local Changes)",
		Description: "Permanently discard local modifications and untracked files to match remote",
		Command:     fmt.Sprintf("git -C %s reset --hard HEAD && git -C %s clean -fd && git -C %s pull", p, p, p),
		Steps: []RemediationStep{
			{Name: "git", Args: []string{"-C", rawPath, "reset", "--hard", "HEAD"}},
			{Name: "git", Args: []string{"-C", rawPath, "clean", "-fd"}},
			{Name: "git", Args: []string{"-C", rawPath, "pull"}},
		},
	}
}
