package cmdpull

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ResolveDualPullRemediationHints produces dual remediation options for an error condition.
func ResolveDualPullRemediationHints(err any, repoDir, repoName string) (string, string, string, string) {
	msg := extractErrorString(err)
	if isMissingRepoFailure(msg) || isMissingRepoDir(repoDir) {
		return resolveMissingRepoDualHints(repoDir, repoName)
	}
	if isNonGitRepoFailure(msg) || (isExistingDir(repoDir) && !isGitRepoDir(repoDir)) {
		return resolveNonGitFolderDualHints(repoDir, repoName)
	}
	if isConflictFailure(msg) {
		return resolveConflictDualHints(repoDir, repoName)
	}
	if isDivergedFailure(msg) {
		return resolveDivergedDualHints(repoDir, repoName)
	}
	if strings.Contains(strings.ToLower(msg), "untracked") {
		return resolveUntrackedDualHints(repoDir)
	}
	return resolveStateOrAuthDualHints(msg, repoDir, repoName)
}

func resolveStateOrAuthDualHints(msg, repoDir, repoName string) (string, string, string, string) {
	if isMissingRepoFailure(msg) || isMissingRepoDir(repoDir) {
		return resolveMissingRepoDualHints(repoDir, repoName)
	}
	if isNonGitRepoFailure(msg) || (isExistingDir(repoDir) && !isGitRepoDir(repoDir)) {
		return resolveNonGitFolderDualHints(repoDir, repoName)
	}
	if isConflictFailure(msg) {
		return resolveConflictDualHints(repoDir, repoName)
	}
	if isDirtyTreeError(msg) {
		return resolveDirtyTreeDualHints(repoDir)
	}
	if isAuthFailure(msg) {
		return resolveAuthDualHints()
	}
	return resolveFallbackDualHints(repoDir, repoName)
}

// ResolveMissingRepoDualHints produces dual remediation options for missing repository errors.
func ResolveMissingRepoDualHints(repoDir, repoName string) (string, string, string, string) {
	return resolveMissingRepoDualHints(repoDir, repoName)
}

func resolveMissingRepoDualHints(repoDir, repoName string) (string, string, string, string) {
	name := repoName
	if name == "" && repoDir != "" {
		name = filepath.Base(repoDir)
	}
	cloneCmd := "gitmap clone " + name
	removeCmd := "gitmap rm --db-only " + name
	return "Clone from Upstream", cloneCmd, "Remove from Registry", removeCmd
}

func resolveNonGitFolderDualHints(repoDir, repoName string) (string, string, string, string) {
	name := resolveEffectiveRepoName(repoName, repoDir)
	cloneCmd := "gitmap clone " + name
	target := name
	if repoDir != "" && !filepath.IsAbs(repoDir) {
		target = filepath.ToSlash(filepath.Clean(repoDir))
	}
	initCmd := fmt.Sprintf("cd %s && git init", target)
	if strings.Contains(target, " ") {
		initCmd = fmt.Sprintf("cd \"%s\" && git init", target)
	}
	return "Clone from Upstream", cloneCmd, "Initialize Repository", initCmd
}

func isDirtyTreeError(msg string) bool {
	low := strings.ToLower(msg)
	return strings.Contains(low, "dirty") || strings.Contains(low, "uncommitted") || strings.Contains(low, "local changes")
}

func resolveConflictDualHints(repoDir, repoName string) (string, string, string, string) {
	name := resolveEffectiveRepoName(repoName, repoDir)
	stashCmd := fmt.Sprintf("gitmap fix %s stash", name)
	discardCmd := fmt.Sprintf("gitmap fix %s discard", name)
	if repoDir != "" && name == "repo" {
		stashCmd = formatRepoGitCmd(repoDir, "stash")
		discardCmd = formatRepoGitCmd(repoDir, "merge --abort")
	}

	return "Stash & Re-pull", stashCmd, "Discard & Abort", discardCmd
}

func resolveDivergedDualHints(repoDir, repoName string) (string, string, string, string) {
	rebaseCmd := formatRepoGitCmd(repoDir, "pull --rebase")
	if repoDir == "" && repoName != "" {
		rebaseCmd = fmt.Sprintf("gitmap pull --rebase %s", repoName)
	}
	resetCmd := formatRepoGitCmd(repoDir, "reset --hard @{u}")
	return "Preserve Local / Rebase", rebaseCmd, "Discard Local / Hard Reset", resetCmd
}

func resolveUntrackedDualHints(repoDir string) (string, string, string, string) {
	stageCmd := formatRepoGitCmd(repoDir, "add .")
	cleanCmd := formatRepoGitCmd(repoDir, "clean -fd")
	return "Track / Stage", stageCmd, "Clean untracked", cleanCmd
}

func cleanRepoPathSlash(p string) string {
	if p == "" {
		return ""
	}
	return strings.ReplaceAll(filepath.ToSlash(p), "\\", "/")
}

func resolveDirtyTreeDualHints(repoDir string) (string, string, string, string) {
	commitCmd := `gitmap cpar "wip: save changes"`
	stashCmd := "gitmap stash"
	if repoDir != "" {
		cleanDir := cleanRepoPathSlash(repoDir)
		commitCmd = fmt.Sprintf("git -C \"%s\" add -A && git -C \"%s\" commit -m \"wip: local changes\" && git -C \"%s\" pull --rebase", cleanDir, cleanDir, cleanDir)
		stashCmd = fmt.Sprintf("git -C \"%s\" stash -u && git -C \"%s\" pull && git -C \"%s\" stash pop", cleanDir, cleanDir, cleanDir)
	}
	return "Commit WIP", commitCmd, "Stash Changes", stashCmd
}

func resolveAuthDualHints() (string, string, string, string) {
	return "Fix Windows Credentials", "gitmap fix-credential", "Deploy SSH Keys", "gitmap ssh deploy-keys"
}

func resolveFallbackDualHints(repoDir, repoName string) (string, string, string, string) {
	name := resolveEffectiveRepoName(repoName, repoDir)
	statusCmd := fmt.Sprintf("gitmap status %s", name)
	if name == "repo" && repoDir != "" {
		statusCmd = formatRepoGitCmd(repoDir, "status")
	}
	fixCmd := fmt.Sprintf("gitmap fix %s", name)
	return "Auto-Fix", fixCmd, "Inspect Status", statusCmd
}

func formatRepoGitCmd(repoDir, gitArgs string) string {
	if repoDir == "" {
		return "git " + gitArgs
	}
	cleanDir := cleanRepoPathSlash(repoDir)
	return fmt.Sprintf("git -C \"%s\" %s", cleanDir, gitArgs)
}

func extractErrorString(err any) string {
	if err == nil {
		return ""
	}
	switch v := err.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func buildRemediationOptions(l1, c1, l2, c2 string) []RemediationOption {
	var opts []RemediationOption
	if l1 != "" && c1 != "" {
		opts = append(opts, RemediationOption{OptionNumber: 1, Title: l1, Command: c1})
	}
	if l2 != "" && c2 != "" {
		opts = append(opts, RemediationOption{OptionNumber: 2, Title: l2, Command: c2})
	}
	return opts
}

func isMissingRepoDir(repoDir string) bool {
	if repoDir == "" {
		return false
	}
	_, err := os.Stat(repoDir)
	return os.IsNotExist(err)
}

func resolveEffectiveRepoName(repoName, repoPath string) string {
	if repoName != "" {
		return repoName
	}
	if repoPath != "" {
		return filepath.Base(repoPath)
	}
	return "repo"
}

func isGitRepoDir(dir string) bool {
	if dir == "" {
		return false
	}
	fi, err := os.Stat(filepath.Join(dir, constants.ExtGit))
	return err == nil && fi != nil
}

func isExistingDir(dir string) bool {
	if dir == "" {
		return false
	}
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}
