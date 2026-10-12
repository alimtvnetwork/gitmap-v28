package gitutil

import (
	"fmt"
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

	return GenerateDefaultRecipes(repoPath)
}

// GenerateDefaultRecipes generates recipes tailored to repository commit and remote state.
func GenerateDefaultRecipes(repoPath string) []RemediationRecipe {
	if !HasInitialCommit(repoPath) {
		return GenerateNoInitialCommitRecipes(repoPath)
	}

	if !HasRemote(repoPath) {
		return GenerateLocalOnlyRecipes(repoPath)
	}

	return []RemediationRecipe{
		GenerateStashRecipe(repoPath),
		GenerateCommitRecipe(repoPath),
		GenerateDiscardRecipe(repoPath),
	}
}

// GenerateNoInitialCommitRecipes produces safe remediation recipes for repos with no initial commit.
func GenerateNoInitialCommitRecipes(repoPath string) []RemediationRecipe {
	p := CleanRepoPath(repoPath)
	rawPath := CleanRepoPathRaw(repoPath)
	hasRemote := HasRemote(repoPath)

	recipes := []RemediationRecipe{
		{
			Title:       "Option 1 (Initial Commit)",
			Description: "Stage and commit all local files as the initial commit",
			Command:     fmt.Sprintf("git -C %s add -A && git -C %s commit -m \"feat: initial commit\"", p, p),
			Steps: []RemediationStep{
				{Name: "git", Args: []string{"-C", rawPath, "add", "-A"}},
				{Name: "git", Args: []string{"-C", rawPath, "commit", "-m", "feat: initial commit"}},
			},
		},
		{
			Title:       "Option 2 (Discard Untracked Files)",
			Description: "Permanently discard untracked files from new repository",
			Command:     fmt.Sprintf("git -C %s clean -fd", p),
			Steps: []RemediationStep{
				{Name: "git", Args: []string{"-C", rawPath, "clean", "-fd"}},
			},
		},
	}

	if hasRemote {
		recipes = append(recipes, RemediationRecipe{
			Title:       "Option 3 (Fetch Remote)",
			Description: "Fetch remote commits to establish branch history",
			Command:     fmt.Sprintf("git -C %s fetch origin", p),
			Steps: []RemediationStep{
				{Name: "git", Args: []string{"-C", rawPath, "fetch", "origin"}},
			},
		})
	} else {
		recipes = append(recipes, RemediationRecipe{
			Title:       "Option 3 (Inspect Status)",
			Description: "Inspect working tree status of new repository",
			Command:     fmt.Sprintf("git -C %s status", p),
			Steps: []RemediationStep{
				{Name: "git", Args: []string{"-C", rawPath, "status"}},
			},
		})
	}

	return recipes
}

// GenerateLocalOnlyRecipes produces recipes for repos that have commits but no remote configured.
func GenerateLocalOnlyRecipes(repoPath string) []RemediationRecipe {
	p := CleanRepoPath(repoPath)
	rawPath := CleanRepoPathRaw(repoPath)

	return []RemediationRecipe{
		{
			Title:       "Option 1 (Commit WIP)",
			Description: "Commit all modified and untracked files locally",
			Command:     fmt.Sprintf("git -C %s add -A && git -C %s commit -m \"wip: local changes\"", p, p),
			Steps: []RemediationStep{
				{Name: "git", Args: []string{"-C", rawPath, "add", "-A"}},
				{Name: "git", Args: []string{"-C", rawPath, "commit", "-m", "wip: local changes"}},
			},
		},
		{
			Title:       "Option 2 (Stash Changes)",
			Description: "Temporarily save local changes to stash",
			Command:     fmt.Sprintf("git -C %s stash -u", p),
			Steps: []RemediationStep{
				{Name: "git", Args: []string{"-C", rawPath, "stash", "-u"}},
			},
		},
		{
			Title:       "Option 3 (Discard Local Changes)",
			Description: "Permanently discard local modifications and untracked files",
			Command:     fmt.Sprintf("git -C %s reset --hard HEAD && git -C %s clean -fd", p, p),
			Steps: []RemediationStep{
				{Name: "git", Args: []string{"-C", rawPath, "reset", "--hard", "HEAD"}},
				{Name: "git", Args: []string{"-C", rawPath, "clean", "-fd"}},
			},
		},
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
