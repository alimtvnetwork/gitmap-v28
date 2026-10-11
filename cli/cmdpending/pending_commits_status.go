package cmdpending

import (
	"strings"
)

func queryRepoStatus(dir string) (int, int, int, []string) {
	porcelain, _ := currentPendingCommitsGitExecutor(dir, "status", "--porcelain")
	return ParsePorcelainStatusLines(porcelain)
}

func ParsePorcelainStatusLines(output string) (int, int, int, []string) {
	var untracked, modified, staged int
	var pendingFiles []string

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if len(line) == 0 {
			continue
		}
		pendingFiles = append(pendingFiles, strings.TrimSpace(line))
		if strings.HasPrefix(line, "??") {
			untracked++
			continue
		}
		stagedDelta, modifiedDelta := classifyStatusChars(line)
		staged += stagedDelta
		modified += modifiedDelta
	}

	return untracked, modified, staged, pendingFiles
}

func classifyStatusChars(line string) (int, int) {
	if len(line) < 2 {
		return 0, 0
	}
	stagedDelta := 0
	modifiedDelta := 0
	idxChar := line[0]
	wtChar := line[1]
	if idxChar != ' ' && idxChar != '?' {
		stagedDelta = 1
	}
	if wtChar != ' ' && wtChar != '?' {
		modifiedDelta = 1
	}
	return stagedDelta, modifiedDelta
}

func parsePendingCommitLines(output string) []string {
	var shas []string
	lines := strings.Split(output, "\n")
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t != "" {
			shas = append(shas, t)
		}
	}
	return shas
}
