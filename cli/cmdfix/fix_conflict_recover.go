package cmdfix

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func isNeedsMergeError(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "needs merge") ||
		strings.Contains(lower, "unmerged") ||
		strings.Contains(lower, "error building trees") ||
		strings.Contains(lower, "cannot save the current worktree state") ||
		strings.Contains(lower, "cannot save the current index state")
}

func isConflictError(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "conflict") ||
		strings.Contains(lower, "failed to merge") ||
		strings.Contains(lower, "automatic merge failed") ||
		strings.Contains(lower, "fix conflicts") ||
		strings.Contains(lower, "could not restore untracked") ||
		strings.Contains(lower, "overwritten by merge") ||
		strings.Contains(lower, "already exists, no checkout") ||
		strings.Contains(lower, "cannot merge") ||
		strings.Contains(lower, "cannot apply") ||
		strings.Contains(lower, "unmerged")
}

func hasMergeHead(repoPath string) bool {
	gitDir := filepath.Join(repoPath, ".git")
	for _, name := range []string{"MERGE_HEAD", "AUTO_MERGE", "REBASE_HEAD", "CHERRY_PICK_HEAD"} {
		if _, err := os.Stat(filepath.Join(gitDir, name)); err == nil {
			return true
		}
	}

	return false
}

func hasActiveConflict(repoPath string) bool {
	if hasMergeHead(repoPath) {
		return true
	}
	files, err := cmdconflict.GetConflictedFiles(repoPath)
	if err != nil {
		return false
	}

	return len(files) > 0
}

func conflictStateReplacement(repoPath string, recipe gitutil.RemediationRecipe) (gitutil.RemediationRecipe, bool) {
	if !hasActiveConflict(repoPath) {
		return recipe, false
	}
	autoResolveRepoConflicts(repoPath)

	return generateConflictResolveRecipe(repoPath), true
}

func generateConflictResolveRecipe(repoPath string) gitutil.RemediationRecipe {
	rawPath := gitutil.CleanRepoPathRaw(repoPath)
	p := gitutil.CleanRepoPath(repoPath)

	return gitutil.RemediationRecipe{
		Title:       "Option 1 (Auto-Resolve Conflicts & Pull)",
		Description: "Resolve merge conflicts, stage files, commit merge, and pull latest",
		Command:     fmt.Sprintf("git -C %s add -A && git -C %s pull", p, p),
		Steps: []gitutil.RemediationStep{
			{Name: "git", Args: []string{"-C", rawPath, "add", "-A"}},
			{Name: "git", Args: []string{"-C", rawPath, "pull"}},
		},
	}
}

func autoResolveRepoConflicts(repoPath string) bool {
	files, _ := cmdconflict.GetConflictedFiles(repoPath)
	for _, f := range files {
		resolveSingleConflictFile(repoPath, f)
	}
	cleanWorkingTreeConflictMarkers(repoPath)
	_ = exec.Command("git", "-C", repoPath, "add", "-A").Run()
	if hasMergeHead(repoPath) {
		concludeMergeCommit(repoPath)
	}

	return true
}

func cleanWorkingTreeConflictMarkers(repoPath string) {
	diag := gitutil.InspectDirtyState(repoPath)
	allCandidateFiles := append(diag.ModifiedFiles, diag.StagedFiles...)
	for _, relFile := range allCandidateFiles {
		absPath := filepath.Join(repoPath, relFile)
		_ = resolveConflictMarkersInFile(absPath)
	}
}

func concludeMergeCommit(repoPath string) bool {
	commitArgs, hasIdentity := buildCommitArgsWithIdentity(repoPath, "fix: auto-resolve merge conflicts")
	if !hasIdentity {
		return false
	}

	return exec.Command("git", commitArgs...).Run() == nil
}

func resolveSingleConflictFile(repoPath, relFile string) {
	absPath := filepath.Join(repoPath, relFile)
	hasResolved := resolveConflictMarkersInFile(absPath)
	if !hasResolved {
		_ = exec.Command("git", "-C", repoPath, "checkout", "--theirs", "--", relFile).Run()
	}
	_ = exec.Command("git", "-C", repoPath, "add", "--", relFile).Run()
}

func resolveConflictMarkersInFile(absPath string) bool {
	b, err := os.ReadFile(absPath)
	if err != nil {
		return false
	}
	content := string(b)
	if !strings.Contains(content, "<<<<<<<") || !strings.Contains(content, ">>>>>>>") {
		return true
	}
	resolved := stripConflictMarkers(content)
	if strings.Contains(resolved, "<<<<<<<") || strings.Contains(resolved, ">>>>>>>") {
		return false
	}

	return os.WriteFile(absPath, []byte(resolved), 0644) == nil
}

func isConflictMarkerLine(trimmed string) bool {
	return strings.HasPrefix(trimmed, "<<<<<<<") ||
		strings.HasPrefix(trimmed, "=======") ||
		strings.HasPrefix(trimmed, ">>>>>>>")
}

func filterConflictLines(lines []string) []string {
	var kept []string
	seen := make(map[string]bool)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if isConflictMarkerLine(trimmed) {
			continue
		}
		if trimmed != "" && strings.HasPrefix(trimmed, "|") && seen[trimmed] {
			continue
		}
		if trimmed != "" && strings.HasPrefix(trimmed, "|") {
			seen[trimmed] = true
		}
		kept = append(kept, line)
	}

	return kept
}

func stripConflictMarkers(content string) string {
	lines := strings.Split(content, "\n")
	kept := filterConflictLines(lines)

	return strings.Join(kept, "\n")
}

func recoverNeedsMerge(step gitutil.RemediationStep, output string) (bool, string) {
	if !isNeedsMergeError(output) {
		return false, ""
	}
	repoPath := extractRepoPath(step.Args)
	if repoPath == "" {
		return false, ""
	}
	isResolved := autoResolveRepoConflicts(repoPath)
	if !isResolved {
		return false, ""
	}
	if exec.Command(step.Name, step.Args...).Run() != nil {
		diag := gitutil.InspectDirtyState(repoPath)

		return !diag.IsDirty, "recovered (conflicts resolved, working tree clean)"
	}

	return true, "recovered (auto-resolved unmerged conflicts and completed step)"
}

func recoverStashPopConflict(step gitutil.RemediationStep, output string) (bool, string) {
	if !checkHasPop(step.Args) || !isConflictError(output) {
		return false, ""
	}
	repoPath := extractRepoPath(step.Args)
	if repoPath == "" {
		return false, ""
	}
	autoResolveRepoConflicts(repoPath)
	_ = exec.Command("git", "-C", repoPath, "stash", "drop").Run()

	return true, "recovered (auto-resolved stash pop conflicts and dropped stash)"
}

func isDivergentPullError(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "divergent") ||
		strings.Contains(lower, "need to specify how to reconcile") ||
		strings.Contains(lower, "not possible to fast-forward")
}

func execDivergentMerge(repoPath string) (bool, string) {
	if exec.Command("git", "-C", repoPath, "pull", "--no-rebase").Run() == nil {
		return true, "recovered (merged divergent remote changes)"
	}

	if hasActiveConflict(repoPath) {
		autoResolveRepoConflicts(repoPath)

		return true, "recovered (auto-resolved divergent conflicts)"
	}

	return false, ""
}

func recoverDivergentPull(step gitutil.RemediationStep, output string) (bool, string) {
	if !isDivergentPullError(output) {
		return false, ""
	}

	repoPath := extractRepoPath(step.Args)
	if repoPath == "" {
		return false, ""
	}

	return execDivergentMerge(repoPath)
}
