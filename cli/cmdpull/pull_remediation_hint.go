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
	if strings.Contains(msg, "Permission denied (publickey)") {
		return "Permission denied (publickey) - SSH key or token rejected"
	}
	if strings.Contains(msg, "Authentication failed") {
		return "Authentication failed - credentials or personal access token rejected"
	}
	if strings.Contains(msg, "CONFLICT") || strings.Contains(msg, "conflict") {
		return "Merge conflict detected during pull"
	}
	if strings.Contains(msg, "Not possible to fast-forward") || strings.Contains(msg, "diverged") {
		return "Cannot fast-forward - local and remote branches have diverged"
	}
	if strings.Contains(msg, "missing repository directory") {
		return "Repository directory does not exist on disk"
	}
	return extractFirstMeaningfulErrorLine(msg)
}

func extractFirstMeaningfulErrorLine(msg string) string {
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "fatal:") || strings.HasPrefix(trimmed, "error:") {
			return trimmed
		}
	}
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
	if strings.Contains(msg, "Can not use the 'cache' credential store on Windows") || strings.Contains(msg, "lack of UNIX socket support") {
		return "Run 'gitmap fix-credential' (alias: fc) to repair Windows git credential store"
	}
	return fmt.Sprintf("gitmap status %s or gitmap fix %s", repoName, repoName)
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
