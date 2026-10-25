package cmdlogviewer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// QueryCommits fetches commit list for targetPath.
func QueryCommits(repoRoot, targetPath string, limit int) ([]CommitItem, error) {
	args := buildLogArgs(limit, targetPath)
	out, err := runGitInDir(repoRoot, args...)
	if err != nil {
		return nil, err
	}

	return parseCommitsOutput(out), nil
}

func buildLogArgs(limit int, targetPath string) []string {
	safeLimit := normalizeLogLimit(limit)
	args := []string{"log", fmt.Sprintf("-n%d", safeLimit), "--date=short", "--pretty=format:%H|%h|%an|%ad|%s"}
	if targetPath != "" && targetPath != "." {
		args = append(args, "--", targetPath)
	}

	return args
}

func normalizeLogLimit(limit int) int {
	if limit <= 0 {
		return 100
	}

	return limit
}

func parseCommitsOutput(raw string) []CommitItem {
	lines := strings.Split(raw, "\n")
	items := make([]CommitItem, 0, len(lines))
	for _, line := range lines {
		item, isOk := parseCommitLine(line)
		if isOk {
			items = append(items, item)
		}
	}

	return items
}

func parseCommitLine(line string) (CommitItem, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return CommitItem{}, false
	}

	parts := strings.SplitN(trimmed, "|", 5)
	if len(parts) < 5 {
		return CommitItem{}, false
	}

	item := CommitItem{
		Hash:      parts[0],
		ShortHash: parts[1],
		Author:    parts[2],
		Date:      parts[3],
		Subject:   parts[4],
	}

	return item, true
}

// QueryCommitDetail retrieves metadata and modified files for a commit.
func QueryCommitDetail(repoRoot, commitHash string) (*CommitDetail, error) {
	detail, err := queryCommitMetadata(repoRoot, commitHash)
	if err != nil {
		return nil, err
	}

	files, err := queryCommitFiles(repoRoot, commitHash)
	if err != nil {
		return nil, err
	}
	detail.Files = files

	return detail, nil
}

func queryCommitMetadata(repoRoot, commitHash string) (*CommitDetail, error) {
	format := "%H\x1f%h\x1f%an\x1f%ae\x1f%ad\x1f%s\x1f%b\x1f%P"
	out, err := runGitInDir(repoRoot, "show", "-s", "--format="+format, commitHash)
	if err != nil {
		return nil, err
	}

	return parseMetadataOutput(out), nil
}

func parseMetadataOutput(raw string) *CommitDetail {
	parts := strings.Split(strings.TrimRight(raw, "\r\n"), "\x1f")
	detail := &CommitDetail{}
	if len(parts) > 0 {
		detail.Hash = parts[0]
	}
	if len(parts) > 1 {
		detail.ShortHash = parts[1]
	}
	if len(parts) > 2 {
		detail.Author = parts[2]
	}
	fillExtendedMetadata(parts, detail)

	return detail
}

func fillExtendedMetadata(parts []string, detail *CommitDetail) {
	if len(parts) > 3 {
		detail.AuthorEmail = parts[3]
	}
	if len(parts) > 4 {
		detail.Date = parts[4]
	}
	if len(parts) > 5 {
		detail.Subject = parts[5]
	}
	if len(parts) > 6 {
		detail.Body = parts[6]
	}
	if len(parts) > 7 && strings.TrimSpace(parts[7]) != "" {
		detail.Parents = strings.Fields(parts[7])
	}
}

func queryCommitFiles(repoRoot, commitHash string) ([]FileChange, error) {
	nameOut, err := runGitInDir(repoRoot, "show", "--format=", "--name-status", commitHash)
	if err != nil {
		return nil, err
	}

	numOut, err := runGitInDir(repoRoot, "show", "--format=", "--numstat", commitHash)
	if err != nil {
		return nil, err
	}

	numMap := parseNumstatMap(numOut)

	return parseNameStatusLines(nameOut, numMap), nil
}

func parseNameStatusLines(raw string, numMap map[string]string) []FileChange {
	lines := strings.Split(raw, "\n")
	files := make([]FileChange, 0, len(lines))
	for _, line := range lines {
		fc, isOk := parseSingleFileChange(line, numMap)
		if isOk {
			files = append(files, fc)
		}
	}

	return files
}

func parseSingleFileChange(line string, numMap map[string]string) (FileChange, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return FileChange{}, false
	}

	parts := strings.Split(trimmed, "\t")
	if len(parts) < 2 {
		return FileChange{}, false
	}

	fc := buildFileChangeFromParts(parts, numMap)

	return fc, true
}

func buildFileChangeFromParts(parts []string, numMap map[string]string) FileChange {
	actionCode := parts[0]
	filePath := parts[1]
	fc := FileChange{
		Action: formatAction(actionCode),
		Path:   filePath,
		Stats:  numMap[filePath],
	}
	if strings.HasPrefix(actionCode, "R") && len(parts) > 2 {
		fc.OldPath = parts[1]
		fc.Path = parts[2]
		fc.Stats = numMap[parts[2]]
	}

	return fc
}

func parseNumstatMap(raw string) map[string]string {
	numMap := make(map[string]string)
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		parts := strings.Split(trimmed, "\t")
		if len(parts) >= 3 {
			stat := fmt.Sprintf("+%s, -%s", parts[0], parts[1])
			numMap[parts[2]] = stat
		}
	}

	return numMap
}

func formatAction(code string) string {
	clean := strings.ToUpper(strings.TrimSpace(code))
	switch {
	case strings.HasPrefix(clean, "A"):
		return "Added"
	case strings.HasPrefix(clean, "M"):
		return "Modified"
	case strings.HasPrefix(clean, "D"):
		return "Deleted"
	case strings.HasPrefix(clean, "R"):
		return "Renamed"
	case strings.HasPrefix(clean, "C"):
		return "Copied"
	default:
		return "Modified"
	}
}

func runGitInDir(dir string, args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	if dir != "" {
		cmd.Dir = dir
	}

	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "git "+strings.Join(args, " "))
	}

	return string(out), nil
}
