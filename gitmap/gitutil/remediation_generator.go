package gitutil

import (
	"path/filepath"
	"strings"
)

type RemediationStep struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type RemediationRecipe struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Command     string            `json:"command"`
	Steps       []RemediationStep `json:"steps"`
}

func GenerateRemediationRecipes(repoPath string, d DirtyDiagnosis) []RemediationRecipe {
	if !d.IsDirty {
		return nil
	}

	return []RemediationRecipe{
		GenerateStashRecipe(repoPath),
		GenerateCommitRecipe(repoPath),
		GenerateDiscardRecipe(repoPath),
	}
}

// CleanRepoPath normalizes slashes for cross-platform shell execution.
func CleanRepoPath(repoPath string) string {
	p := filepath.ToSlash(filepath.Clean(repoPath))
	if strings.Contains(p, " ") {
		return `"` + p + `"`
	}
	return p
}

// CleanRepoPathRaw normalizes slashes without quotes for direct exec.Command.
func CleanRepoPathRaw(repoPath string) string {
	return filepath.ToSlash(filepath.Clean(repoPath))
}
