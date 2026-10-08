package cmdpull

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// ResolvePullErrorDetails extracts a clean, human-readable error summary from pull state.
func ResolvePullErrorDetails(s *PullRepoState) string {
	if s == nil {
		return ""
	}
	if isMissingRepoDir(s.RepoPath) {
		return "Repository directory does not exist on disk"
	}
	if s.IsDirty || s.Changes == "dirty" {
		return "working tree has uncommitted changes"
	}
	if s.ErrorMsg == "" {
		return ""
	}
	return classifyPullErrorString(s.ErrorMsg)
}

func classifyPullErrorString(msg string) string {
	if isWincredmanFailure(msg) {
		return "Windows Credential Manager (wincredman) failed to persist credentials"
	}
	if summary := classifyStandardError(msg); summary != "" {
		return summary
	}
	return extractFirstMeaningfulErrorLine(msg)
}

func classifyStandardError(msg string) string {
	if isAuthFailure(msg) {
		return classifyAuthError(msg)
	}
	if isConflictFailure(msg) {
		return "Merge conflict detected during pull"
	}
	if isDivergedFailure(msg) {
		return "Cannot fast-forward - local and remote branches have diverged"
	}
	if isMissingRepoFailure(msg) {
		return "Repository directory does not exist on disk"
	}
	return ""
}

func classifyAuthError(msg string) string {
	if strings.Contains(msg, "Permission denied (publickey)") {
		return "Permission denied (publickey) - SSH key or token rejected"
	}
	return "Authentication failed - credentials or personal access token rejected"
}

func extractFirstMeaningfulErrorLine(msg string) string {
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "fatal:") || strings.HasPrefix(trimmed, "error:") {
			return trimmed
		}
	}
	return fallbackErrorLine(lines)
}

func fallbackErrorLine(lines []string) string {
	if len(lines) == 0 {
		return "pull execution failed"
	}
	first := strings.TrimSpace(lines[0])
	if first == "" {
		return "pull execution failed"
	}
	if len(first) > 80 {
		return first[:77] + "..."
	}
	return first
}

// ResolvePullRemediationHint generates an actionable next-step command for pull failures.
func ResolvePullRemediationHint(s *PullRepoState) string {
	if s == nil {
		return ""
	}
	if isMissingRepoDir(s.RepoPath) {
		name := resolveEffectiveRepoName(s.RepoName, s.RepoPath)
		return fmt.Sprintf("gitmap clone %s", name)
	}
	if s.IsDirty || s.Changes == "dirty" {
		return fmt.Sprintf("gitmap fix %s, gitmap cpar, or gitmap stash", s.RepoName)
	}
	if s.ErrorMsg == "" {
		return ""
	}
	return buildActionableRemediation(s.RepoName, s.RepoPath, s.ErrorMsg)
}

func buildActionableRemediation(repoName, repoPath, msg string) string {
	effectiveName := resolveEffectiveRepoName(repoName, repoPath)
	if isMissingRepoFailure(msg) || isMissingRepoDir(repoPath) {
		return fmt.Sprintf("gitmap clone %s", effectiveName)
	}
	if isWincredmanFailure(msg) {
		return "Run 'gitmap fix-credential' (alias: fc) or check Windows Credential Manager service"
	}
	if hint := buildSpecificRemediation(effectiveName, msg); hint != "" {
		return hint
	}
	return fmt.Sprintf("gitmap status %s or gitmap fix %s", effectiveName, effectiveName)
}

func buildSpecificRemediation(repoName, msg string) string {
	if isMissingRepoFailure(msg) {
		return fmt.Sprintf("gitmap clone %s", repoName)
	}
	if isAuthFailure(msg) {
		return fmt.Sprintf("gitmap status %s or gitmap fix %s", repoName, repoName)
	}
	if isConflictFailure(msg) {
		return fmt.Sprintf("gitmap fix %s or gitmap stash", repoName)
	}
	if isDivergedFailure(msg) {
		return fmt.Sprintf("gitmap pull %s", repoName)
	}
	if isCacheFailure(msg) {
		return "Run 'gitmap fix-credential' (alias: fc) to repair Windows git credential store"
	}
	return ""
}

func isWincredmanFailure(msg string) bool {
	return strings.Contains(msg, "wincredman") ||
		strings.Contains(msg, "Unable to persist credentials with the 'wincredman' credential store") ||
		strings.Contains(msg, "Unable to persist credentials")
}

func isCacheFailure(msg string) bool {
	return strings.Contains(msg, "Can not use the 'cache' credential store on Windows") ||
		strings.Contains(msg, "lack of UNIX socket support")
}

func isAuthFailure(msg string) bool {
	return strings.Contains(msg, "Permission denied (publickey)") || strings.Contains(msg, "Authentication failed")
}

func isConflictFailure(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "conflict") ||
		strings.Contains(lower, "automatic merge failed") ||
		strings.Contains(lower, "unmerged files")
}

func isDivergedFailure(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "not possible to fast-forward") ||
		strings.Contains(lower, "diverged") ||
		strings.Contains(lower, "cannot fast-forward") ||
		strings.Contains(lower, "reconcile divergent") ||
		strings.Contains(lower, "non-fast-forward")
}

func isMissingRepoFailure(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "missing repository directory") ||
		strings.Contains(lower, "directory does not exist") ||
		strings.Contains(lower, "no such file or directory") ||
		strings.Contains(lower, "cannot change to") ||
		strings.Contains(lower, "does not exist")
}

// RemediationOption represents an actionable alternative fix command.
type RemediationOption struct {
	OptionNumber int    `json:"option_number"`
	Title        string `json:"title"`
	Command      string `json:"command"`
}

// StructuredRemediation captures error reason and dual-option actionable fixes.
type StructuredRemediation struct {
	Reason  string              `json:"reason"`
	Options []RemediationOption `json:"options"`
}

// ResolveStructuredRemediation resolves structured dual options for a repository state.
func ResolveStructuredRemediation(s *PullRepoState) StructuredRemediation {
	if s == nil {
		return StructuredRemediation{}
	}
	reason := ResolvePullErrorDetails(s)
	if reason == "" {
		reason = "pull execution failed"
	}
	l1, c1, l2, c2 := resolveStateDualHints(s)
	return StructuredRemediation{
		Reason:  reason,
		Options: buildRemediationOptions(l1, c1, l2, c2),
	}
}

func resolveStateDualHints(s *PullRepoState) (string, string, string, string) {
	if s == nil {
		return "", "", "", ""
	}
	if isMissingRepoDir(s.RepoPath) {
		return resolveMissingRepoDualHints(s.RepoPath, s.RepoName)
	}
	if s.IsDirty || s.Changes == "dirty" {
		return resolveDirtyStateDualHints(s)
	}
	return ResolveDualPullRemediationHints(s.ErrorMsg, s.RepoPath, s.RepoName)
}

func resolveDirtyStateDualHints(s *PullRepoState) (string, string, string, string) {
	if s.RepoPath == "" {
		return resolveDirtyTreeDualHints(s.RepoPath)
	}

	diag := gitutil.InspectDirtyState(s.RepoPath)
	hasOnlyUntracked := diag.UntrackedCount > 0 && diag.ModifiedCount == 0 && diag.StagedCount == 0 && diag.DeletedCount == 0
	if hasOnlyUntracked {
		return resolveUntrackedDualHints(s.RepoPath)
	}

	return resolveDirtyTreeDualHints(s.RepoPath)
}

// ResolveDualPullRemediationHints produces dual remediation options for an error condition.
func ResolveDualPullRemediationHints(err any, repoDir, repoName string) (string, string, string, string) {
	msg := extractErrorString(err)
	if isMissingRepoFailure(msg) || isMissingRepoDir(repoDir) {
		return resolveMissingRepoDualHints(repoDir, repoName)
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
	return "Clone from Remote", cloneCmd, "Remove from Registry", removeCmd
}

func isDirtyTreeError(msg string) bool {
	low := strings.ToLower(msg)
	return strings.Contains(low, "dirty") || strings.Contains(low, "uncommitted") || strings.Contains(low, "local changes")
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

func resolveDirtyTreeDualHints(repoDir string) (string, string, string, string) {
	commitCmd := `gitmap cpar "wip: save changes"`
	stashCmd := formatRepoGitCmd(repoDir, "stash")
	if repoDir == "" {
		stashCmd = "gitmap stash"
	}
	return "Commit WIP", commitCmd, "Stash Changes", stashCmd
}

func resolveAuthDualHints() (string, string, string, string) {
	return "Fix Windows Credentials", "gitmap fix-credential", "Deploy SSH Keys", "gitmap ssh deploy-keys"
}

func resolveFallbackDualHints(repoDir, repoName string) (string, string, string, string) {
	statusCmd := fmt.Sprintf("gitmap status %s", repoName)
	if repoName == "" {
		statusCmd = formatRepoGitCmd(repoDir, "status")
	}
	pullCmd := fmt.Sprintf("gitmap pull %s", repoName)
	if repoName == "" {
		pullCmd = formatRepoGitCmd(repoDir, "pull")
	}
	return "Inspect Status", statusCmd, "Re-pull Repo", pullCmd
}

func formatRepoGitCmd(repoDir, gitArgs string) string {
	if repoDir == "" {
		return "git " + gitArgs
	}
	return fmt.Sprintf("git -C %q %s", repoDir, gitArgs)
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
