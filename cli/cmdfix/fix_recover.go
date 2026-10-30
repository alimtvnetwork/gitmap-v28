package cmdfix

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func isRepoAlreadyClean(item *cmdremediation.RemediationItem) bool {
	if item == nil || !gitutil.IsValidGitRepository(item.RepoPath) {
		return false
	}

	if gitutil.InspectDirtyState(item.RepoPath).IsDirty {
		return false
	}

	return !hasMergeHead(item.RepoPath) && !hasGitLockFile(item.RepoPath)
}

func sanitizeRecipeForRepo(repoPath string, recipe gitutil.RemediationRecipe) gitutil.RemediationRecipe {
	if repoPath == "" {
		return recipe
	}
	if r, ok := conflictStateReplacement(repoPath, recipe); ok {
		return r
	}
	if r, ok := noInitialCommitReplacement(repoPath, recipe); ok {
		return r
	}
	if r, ok := localOnlyReplacement(repoPath, recipe); ok {
		return r
	}
	return recipe
}

func noInitialCommitReplacement(repoPath string, recipe gitutil.RemediationRecipe) (gitutil.RemediationRecipe, bool) {
	if gitutil.HasInitialCommit(repoPath) || !isStashOrResetRecipe(recipe) {
		return recipe, false
	}
	return firstRecipe(gitutil.GenerateNoInitialCommitRecipes(repoPath), recipe)
}

func localOnlyReplacement(repoPath string, recipe gitutil.RemediationRecipe) (gitutil.RemediationRecipe, bool) {
	if gitutil.HasRemote(repoPath) || !isPullRecipe(recipe) {
		return recipe, false
	}
	return firstRecipe(gitutil.GenerateLocalOnlyRecipes(repoPath), recipe)
}

func firstRecipe(recipes []gitutil.RemediationRecipe, fallback gitutil.RemediationRecipe) (gitutil.RemediationRecipe, bool) {
	if len(recipes) == 0 {
		return fallback, false
	}
	return recipes[0], true
}

func isStashOrResetRecipe(recipe gitutil.RemediationRecipe) bool {
	cmdLower := strings.ToLower(recipe.Command)
	titleLower := strings.ToLower(recipe.Title)

	return strings.Contains(cmdLower, "stash") || strings.Contains(titleLower, "stash") ||
		strings.Contains(cmdLower, "reset --hard") || strings.Contains(titleLower, "discard")
}

func isPullRecipe(recipe gitutil.RemediationRecipe) bool {
	cmdLower := strings.ToLower(recipe.Command)

	return strings.Contains(cmdLower, "pull")
}

func extractRepoPath(args []string) string {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-C" {
			return args[i+1]
		}
	}

	return ""
}

func isNoInitialCommitError(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "you do not have the initial commit yet") ||
		strings.Contains(lower, "does not have any commits yet")
}

func recoverNoInitialCommit(step gitutil.RemediationStep, output string) (bool, string) {
	if !isNoInitialCommitError(output) {
		return false, ""
	}

	repoPath := extractRepoPath(step.Args)
	if repoPath == "" {
		return false, ""
	}

	diag := gitutil.InspectDirtyState(repoPath)
	if !diag.IsDirty {
		return true, "clean (no commits yet, working tree clean)"
	}

	cmdAdd := exec.Command("git", "-C", repoPath, "add", "-A")
	if err := cmdAdd.Run(); err != nil {
		return false, ""
	}

	commitArgs, hasIdentity := buildCommitArgsWithIdentity(repoPath, "feat: initial commit")
	if !hasIdentity {
		return false, "failed: git identity not configured"
	}
	cmdCommit := exec.Command("git", commitArgs...)
	if err := cmdCommit.Run(); err != nil {
		return false, ""
	}

	return true, "recovered (created initial commit)"
}

func isBenignPullNoTracking(step gitutil.RemediationStep, output string) bool {
	hasPull := false
	for _, arg := range step.Args {
		if arg == "pull" {
			hasPull = true

			break
		}
	}
	if !hasPull {
		return false
	}

	lower := strings.ToLower(output)

	return strings.Contains(lower, "no tracking information") ||
		strings.Contains(lower, "no remote repository specified") ||
		strings.Contains(lower, "no configured push destination")
}
