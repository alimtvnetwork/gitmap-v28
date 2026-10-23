package cmdpull

import (
	"os/exec"
	"strconv"
	"strings"
)

func resolveCommitMetadata(repoPath string) (string, string) {
	cmd := exec.Command("git", "-C", repoPath, "log", "-1", "--pretty=format:%an%x1f%s")
	out, err := cmd.Output()
	if err != nil {
		return "", ""
	}

	parts := strings.Split(string(out), "\x1f")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), ""
	}

	return "", ""
}

func resolveCommitTrace(repoPath, oldSHA, newSHA string) string {
	diffTrace := queryRangeCommitTrace(repoPath, oldSHA, newSHA)
	if len(diffTrace) > 0 {
		return diffTrace
	}

	return queryLatestCommitTrace(repoPath)
}

func queryRangeCommitTrace(repoPath, oldSHA, newSHA string) string {
	if oldSHA == "" || newSHA == "" || oldSHA == newSHA {
		return ""
	}
	cmd := exec.Command("git", "-C", repoPath, "log", oldSHA+".."+newSHA, "--oneline")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func queryLatestCommitTrace(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "log", "-1", "--oneline")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func parseFilesChanged(changes, repoPath, oldSHA, newSHA string) int {
	count := extractFilesFromChangesString(changes)
	if count > 0 {
		return count
	}

	return queryDiffStatFilesCount(repoPath, oldSHA, newSHA)
}

func extractFilesFromChangesString(changes string) int {
	start := strings.Index(changes, "(")
	end := strings.Index(changes, ")")
	if start == -1 || end <= start+1 {
		return parseFilesKeywordCount(changes)
	}

	num, err := strconv.Atoi(strings.TrimSpace(changes[start+1 : end]))
	if err != nil {
		return parseFilesKeywordCount(changes)
	}

	return num
}

func parseFilesKeywordCount(changes string) int {
	if !strings.Contains(changes, "file") {
		return 0
	}

	fields := strings.Fields(changes)
	for i, f := range fields {
		count, isMatch := tryParseFileField(fields, i, f)
		if isMatch {
			return count
		}
	}

	return 0
}

func tryParseFileField(fields []string, i int, f string) (int, bool) {
	if i <= 0 || !strings.HasPrefix(f, "file") {
		return 0, false
	}

	num, err := strconv.Atoi(fields[i-1])
	if err != nil {
		return 0, false
	}

	return num, true
}

func queryDiffStatFilesCount(repoPath, oldSHA, newSHA string) int {
	if oldSHA == "" || newSHA == "" || oldSHA == newSHA {
		return 0
	}

	cmd := exec.Command("git", "-C", repoPath, "diff", "--shortstat", oldSHA+".."+newSHA)
	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	return parseFilesKeywordCount(string(out))
}
