package gitutil

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// DetectPRStatus checks for tracking branches or remote pull references.
func DetectPRStatus(repoPath string) string {
	if prs, ok := detectGitHubPRs(repoPath); ok && prs != "0 PRs" {
		return prs
	}
	return detectGitBranchTracking(repoPath)
}

func detectGitHubPRs(repoPath string) (string, bool) {
	cmdGh := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number")
	cmdGh.Dir = repoPath
	out, err := cmdGh.Output()
	if err != nil {
		return "", false
	}
	var prs []interface{}
	if err := json.Unmarshal(out, &prs); err != nil {
		return "", false
	}
	return formatPRCount(len(prs)), true
}

func formatPRCount(count int) string {
	if count == 0 {
		return "0 PRs"
	}
	if count == 1 {
		return "1 Open PR"
	}
	return fmt.Sprintf("%d Open PRs", count)
}

func detectGitBranchTracking(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "for-each-ref", "--format=%(upstream:track)", "refs/heads")
	out, err := cmd.Output()
	if err != nil {
		return "—"
	}
	track := strings.TrimSpace(string(out))
	if strings.Contains(track, "ahead") && strings.Contains(track, "behind") {
		return "diverged"
	}
	if strings.Contains(track, "ahead") {
		return "ahead"
	}
	if strings.Contains(track, "behind") {
		return "behind"
	}
	if track == "" {
		return "local"
	}
	return "synced"
}
