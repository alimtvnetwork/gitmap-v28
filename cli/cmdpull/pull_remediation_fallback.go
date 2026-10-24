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
	if isWincredmanFailure(msg) || isCacheFailure(msg) || isAuthFailure(msg) || isPromptDisabledFailure(msg) {
		return resolveCredentialAndSSHDualHints(repoDir, repoName)
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
	target := quoteRemediationTarget(resolveEffectiveRepoName(repoName, repoDir))
	cloneCmd := fmt.Sprintf("gitmap clone %s", target)
	fixCmd := fmt.Sprintf("gitmap fix %s", target)
	return "Clone from Upstream", cloneCmd, "Initialize Repository", fixCmd
}

func isDirtyTreeError(msg string) bool {
	low := strings.ToLower(msg)
	return strings.Contains(low, "dirty") || strings.Contains(low, "uncommitted") || strings.Contains(low, "local changes")
}

func resolveConflictDualHints(repoDir, repoName string) (string, string, string, string) {
	target := quoteRemediationTarget(resolveEffectiveRepoName(repoName, repoDir))
	stashCmd := fmt.Sprintf("gitmap fix %s stash", target)
	discardCmd := fmt.Sprintf("gitmap fix %s discard", target)
	return "Stash & Re-pull", stashCmd, "Discard & Abort", discardCmd
}

func resolveDivergedDualHints(repoDir, repoName string) (string, string, string, string) {
	target := quoteRemediationTarget(resolveEffectiveRepoName(repoName, repoDir))
	pullCmd := fmt.Sprintf("gitmap pull %s", target)
	fixCmd := fmt.Sprintf("gitmap fix %s", target)
	return "Preserve Local / Rebase", pullCmd, "Discard Local / Auto-Fix", fixCmd
}

func resolveUntrackedDualHints(repoDir string) (string, string, string, string) {
	target := quoteRemediationTarget(resolveRemediationTarget(repoDir))
	wipCmd := fmt.Sprintf("gitmap fix %s wip", target)
	cleanCmd := fmt.Sprintf("gitmap fix %s clean", target)
	return "Track / Stage", wipCmd, "Clean untracked", cleanCmd
}

func cleanRepoPathSlash(p string) string {
	if p == "" {
		return ""
	}
	return strings.ReplaceAll(filepath.ToSlash(p), "\\", "/")
}

func resolveDirtyTreeDualHints(repoDir string) (string, string, string, string) {
	target := quoteRemediationTarget(resolveRemediationTarget(repoDir))
	commitCmd := fmt.Sprintf(`gitmap sends cp %s "wip: save changes"`, target)
	stashCmd := fmt.Sprintf("gitmap stash %s", target)
	return "Commit WIP", commitCmd, "Stash Changes", stashCmd
}

func resolveRemediationTarget(repoDir string) string {
	if repoDir == "" {
		return "repo"
	}
	base := filepath.Base(cleanRepoPathSlash(repoDir))
	if base == "" || base == "." || base == "/" || base == "\\" {
		return "repo"
	}
	return base
}

func quoteRemediationTarget(s string) string {
	if strings.Contains(s, " ") && !strings.HasPrefix(s, "\"") {
		return fmt.Sprintf(`"%s"`, s)
	}
	return s
}

func resolveCredentialAndSSHDualHints(repoDir, repoName string) (string, string, string, string) {
	target := quoteRemediationTarget(resolveEffectiveRepoName(repoName, repoDir))
	fixCmd := fmt.Sprintf("gitmap fix %s", target)
	return "Auto-Fix Auth & Convert to SSH", fixCmd, "Fix Credential Store", "gitmap fix-credential"
}

func resolveFallbackDualHints(repoDir, repoName string) (string, string, string, string) {
	target := quoteRemediationTarget(resolveEffectiveRepoName(repoName, repoDir))
	statusCmd := fmt.Sprintf("gitmap status %s", target)
	fixCmd := fmt.Sprintf("gitmap fix %s", target)
	return "Auto-Fix", fixCmd, "Inspect Status", statusCmd
}

// Deprecated: formatRepoGitCmd is deprecated in favor of native GitMap commands.
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
