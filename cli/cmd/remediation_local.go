package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// InspectLocalRepoItem inspects a local filesystem path for git status.
func InspectLocalRepoItem(query string) (*RemediationItem, string) {
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil, ""
	}

	absPath, err := filepath.Abs(cleanQuery)
	if err != nil {
		return nil, ""
	}

	gitDir := filepath.Join(absPath, ".git")
	if _, statErr := os.Stat(gitDir); statErr != nil {
		return nil, ""
	}

	return buildLocalRemediationItem(absPath)
}

func buildLocalRemediationItem(absPath string) (*RemediationItem, string) {
	diag := gitutil.InspectDirtyState(absPath)
	repoName := filepath.Base(absPath)
	if !diag.IsDirty {
		msg := fmt.Sprintf("%s Repository %q is clean. No remediation needed.", constants.ColorGreen+"✓"+constants.ColorReset, repoName)

		return nil, msg
	}

	recipes := gitutil.GenerateRemediationRecipes(absPath, diag)
	if len(recipes) == 0 {
		return nil, ""
	}

	item := RemediationItem{
		RepoPath:      absPath,
		RepoName:      repoName,
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
		Files:         diag.AllFiles,
	}

	return &item, ""
}
