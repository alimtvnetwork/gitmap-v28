package cmdhistory

import (
	"fmt"
	"os/exec"
	"strings"
)

const rangeSeparator = ".."

func expandCommitRange(repo, span string) ([]string, error) {
	if !strings.Contains(span, rangeSeparator) {
		return nil, fmt.Errorf("range %q must look like start..end", span)
	}
	cmd := exec.Command("git", "-C", repo, "rev-list", "--reverse", span)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git rev-list %s: %w", span, err)
	}

	return splitSHAs(string(out)), nil
}

func splitSHAs(raw string) []string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var shas []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			shas = append(shas, trimmed)
		}
	}

	return shas
}
