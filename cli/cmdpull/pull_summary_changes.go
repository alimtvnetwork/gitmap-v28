package cmdpull

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func renderIncomingPullChanges(state *PullRepoState, cwd string) {
	if !hasIncomingChanges(state) {
		return
	}
	lines := queryGitDiffStatDetailed(cwd, state.OldSHA, state.NewSHA)
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			fmt.Printf("   %s\n", l)
		}
	}
}

func hasIncomingChanges(state *PullRepoState) bool {
	if state == nil || state.OldSHA == "" || state.NewSHA == "" {
		return false
	}

	return state.OldSHA != state.NewSHA && state.OldSHA != "-" && state.NewSHA != "-"
}

func queryGitDiffStatDetailed(cwd, oldSHA, newSHA string) []string {
	args := buildDiffStatArgs(cwd, oldSHA, newSHA)
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimRight(string(out), "\r\n"), "\n")
}

func buildDiffStatArgs(cwd, oldSHA, newSHA string) []string {
	args := []string{"-C", cwd}
	if resolveTerminalMode() {
		args = append(args, "-c", "color.diff=always")
	}

	return append(args, "diff", "--stat", oldSHA+".."+newSHA)
}

func renderGitDirtyStatus(cwd string) {
	lines := queryGitStatusShort(cwd)
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			fmt.Printf("   %s\n", l)
		}
	}
}

func renderInlineDirtyFiles(diag gitutil.DirtyDiagnosis) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	renderGitDirtyStatus(cwd)
}

func queryGitStatusShort(cwd string) []string {
	args := buildStatusShortArgs(cwd)
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return filterGitStatusLines(string(out))
}

func buildStatusShortArgs(cwd string) []string {
	args := []string{"-C", cwd}
	if resolveTerminalMode() {
		args = append(args, "-c", "color.status=always")
	}

	return append(args, "status", "--short")
}

func filterGitStatusLines(raw string) []string {
	lines := strings.Split(strings.TrimRight(raw, "\r\n"), "\n")
	var filtered []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if isIgnoredStatusLine(trimmed) {
			continue
		}
		filtered = append(filtered, l)
	}

	return filtered
}

func isIgnoredStatusLine(trimmed string) bool {
	if trimmed == "" {
		return true
	}

	return strings.Contains(trimmed, ".gitmap")
}
