package cmdagy

import (
	"fmt"
	"os/exec"
	"strings"
)

func runGitLogStat(repoDir string, count int) (string, bool) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-n%d", count), "--stat", "--no-merges")
	if len(repoDir) > 0 {
		cmd.Dir = repoDir
	}

	out, err := cmd.Output()

	return strings.TrimSpace(string(out)), err == nil && len(out) > 0
}

// ExtractGitLog extracts recent commit history from local git repository.
func ExtractGitLog(repoDir string, count int) string {
	if count <= 0 {
		count = 5
	}

	if out, hasStat := runGitLogStat(repoDir, count); hasStat {
		return out
	}

	return fallbackGitLog(repoDir, count)
}

func fallbackGitLog(repoDir string, count int) string {
	cmd := exec.Command("git", "log", fmt.Sprintf("-n%d", count), "--oneline")
	if len(repoDir) > 0 {
		cmd.Dir = repoDir
	}

	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}

	return "(recent git commit log unavailable)"
}
