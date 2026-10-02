package cmdpull

import (
	"fmt"
	"strings"
)

// ResolvePullErrorDetails extracts a clean, human-readable error summary from pull state.
func ResolvePullErrorDetails(s *PullRepoState) string {
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
	if s.IsDirty || s.Changes == "dirty" {
		return fmt.Sprintf("gitmap fix %s, gitmap cpar, or gitmap stash", s.RepoName)
	}
	if s.ErrorMsg == "" {
		return ""
	}
	return buildActionableRemediation(s.RepoName, s.RepoPath, s.ErrorMsg)
}

func buildActionableRemediation(repoName, repoPath, msg string) string {
	if isWincredmanFailure(msg) {
		return "Run 'gitmap fix-credential' (alias: fc) or check Windows Credential Manager service"
	}
	if hint := buildSpecificRemediation(repoName, msg); hint != "" {
		return hint
	}
	return fmt.Sprintf("gitmap status %s or gitmap fix %s", repoName, repoName)
}

func buildSpecificRemediation(repoName, msg string) string {
	if isAuthFailure(msg) {
		return fmt.Sprintf("gitmap status %s or gitmap fix %s", repoName, repoName)
	}
	if isConflictFailure(msg) {
		return fmt.Sprintf("gitmap fix %s or gitmap stash", repoName)
	}
	if isDivergedFailure(msg) {
		return fmt.Sprintf("gitmap pull %s", repoName)
	}
	if isMissingRepoFailure(msg) {
		return fmt.Sprintf("gitmap clone %s", repoName)
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
	return strings.Contains(msg, "CONFLICT") || strings.Contains(msg, "conflict")
}

func isDivergedFailure(msg string) bool {
	return strings.Contains(msg, "Not possible to fast-forward") || strings.Contains(msg, "diverged")
}

func isMissingRepoFailure(msg string) bool {
	return strings.Contains(msg, "missing repository directory")
}
