// Package gitutil — upstream_delta.go calculates ahead and behind commit counts.
package gitutil

import (
	"os/exec"
	"strconv"
	"strings"
)

// GetUpstreamDelta returns (ahead, behind) commit counts relative to upstream.
func GetUpstreamDelta(repoPath string) (int, int) {
	cmd := exec.Command("git", "-C", repoPath, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) != 2 {
		return 0, 0
	}
	ahead, _ := strconv.Atoi(parts[0])
	behind, _ := strconv.Atoi(parts[1])
	return ahead, behind
}
