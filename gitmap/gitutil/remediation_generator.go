package gitutil

import (
	"path/filepath"
	"strings"
)

type RemediationRecipe struct {
	Title       string
	Description string
	Command     string
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
