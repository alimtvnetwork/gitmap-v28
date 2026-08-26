// Package gitutil — remediation_generator.go produces copy-pasteable remediation recipes.
package gitutil

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
