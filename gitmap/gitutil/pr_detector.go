// Package gitutil — pr_detector.go detects remote PR and tracking branch status.
package gitutil

import (
	"os/exec"
	"strings"
)

// DetectPRStatus checks for tracking branches or remote pull references.
func DetectPRStatus(repoPath string) string {
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
