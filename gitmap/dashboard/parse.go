// Package dashboard collects Git repository data for the HTML dashboard.
package dashboard

import (
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// parseCommitLog splits raw git log output into CommitInfo entries.
func parseCommitLog(raw string) []model.CommitInfo {
	blocks := strings.Split(strings.TrimSpace(raw), "\n\n")
	commits := make([]model.CommitInfo, 0, len(blocks))
	for _, block := range blocks {
		commit, ok := parseOneCommit(block)
		if ok {
			commits = append(commits, commit)
		}
	}
	return commits
}

// parseOneCommit extracts a CommitInfo from a single log block.
func parseOneCommit(block string) (model.CommitInfo, bool) {
	lines := strings.Split(strings.TrimSpace(block), "\n")
	if len(lines) == 0 {
		return model.CommitInfo{}, false
	}
	fields := strings.SplitN(lines[0], "|", 7)
	if len(fields) < 7 {
		return model.CommitInfo{}, false
	}
	files, ins, del := parseNumstat(lines[1:])
	return makeCommitInfo(fields, files, ins, del), true
}

// makeCommitInfo constructs a CommitInfo struct from parsed fields.
func makeCommitInfo(f []string, files, ins, del int) model.CommitInfo {
	return model.CommitInfo{
		SHA:          f[0],
		ShortSHA:     f[1],
		Author:       f[2],
		Email:        f[3],
		Date:         f[4],
		Message:      f[5],
		IsMerge:      isMergeCommit(f[6]),
		FilesChanged: files,
		Insertions:   ins,
		Deletions:    del,
	}
}

// parseNumstat tallies file changes, insertions, and deletions.
func parseNumstat(lines []string) (int, int, int) {
	files, ins, del := 0, 0, 0
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			files++
			added, _ := strconv.Atoi(parts[0])
			removed, _ := strconv.Atoi(parts[1])
			ins += added
			del += removed
		}
	}
	return files, ins, del
}

// parseBranchLines converts raw branch lines to BranchInfo structs.
func parseBranchLines(lines []string) []model.BranchInfo {
	branches := make([]model.BranchInfo, 0, len(lines))
	for _, line := range lines {
		fields := strings.SplitN(line, "|", 3)
		if len(fields) >= 3 {
			branches = append(branches, model.BranchInfo{
				Name:           fields[0],
				IsRemote:       strings.HasPrefix(fields[0], constants.GitOriginPrefix),
				LastCommitSHA:  fields[1],
				LastCommitDate: fields[2],
			})
		}
	}
	return branches
}

// parseTagLines converts raw tag lines to TagInfo structs with distances.
func parseTagLines(repoPath string, lines []string) []model.TagInfo {
	tags := make([]model.TagInfo, 0, len(lines))
	for i, line := range lines {
		fields := strings.SplitN(line, "|", 3)
		if len(fields) >= 3 {
			tags = append(tags, parseTagEntry(repoPath, lines, i, fields))
		}
	}
	return tags
}

// parseTagEntry constructs a single TagInfo with its commit distance.
func parseTagEntry(repoPath string, lines []string, i int, fields []string) model.TagInfo {
	count := 0
	if i < len(lines)-1 {
		nextFields := strings.SplitN(lines[i+1], "|", 3)
		count = resolveTagDistance(repoPath, fields[0], nextFields)
	}
	return model.TagInfo{
		Name:        fields[0],
		SHA:         fields[1],
		Date:        fields[2],
		CommitCount: count,
	}
}

// resolveTagDistance calculates commit distance between adjacent tags.
func resolveTagDistance(repoPath, field string, nextFields []string) int {
	if len(nextFields) >= 1 {
		return queryTagDistance(repoPath, nextFields[0], field)
	}
	return 0
}
