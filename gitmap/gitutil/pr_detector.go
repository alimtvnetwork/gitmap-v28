package gitutil

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// DetectPRStatus checks for tracking branches or remote pull references,
// and optionally counts Open PRs using the GitHub CLI.
func DetectPRStatus(repoPath string) string {
	// 1. Check for open PRs via GitHub CLI
	cmdGh := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number")
	cmdGh.Dir = repoPath
	if out, err := cmdGh.Output(); err == nil {
		var prs []interface{}
		if json.Unmarshal(out, &prs) == nil {
			if len(prs) == 0 {
				return "0 PRs"
			} else if len(prs) == 1 {
				return "1 Open PR"
			}
			return fmt.Sprintf("%d Open PRs", len(prs))
		}
	}

	// 2. Fallback to basic git branch tracking
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
