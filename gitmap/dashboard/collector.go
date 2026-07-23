// Package dashboard collects Git repository data for the HTML dashboard.
package dashboard

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// CollectOptions holds user-facing flags for data collection.
type CollectOptions struct {
	RepoPath string
	Limit    int
	Since    string
	NoMerges bool
}

// Collect gathers all repository data into a single DashboardData struct.
func Collect(opts CollectOptions) (model.DashboardData, error) {
	commits, err := collectCommits(opts)
	if err != nil {
		return model.DashboardData{}, fmt.Errorf(constants.ErrDashCollect, err)
	}

	branches := collectBranches(opts.RepoPath)
	tags := collectTags(opts.RepoPath)
	commits = attachTagsToCommits(commits, tags)
	authors := buildAuthors(commits)
	frequency := buildFrequency(commits)
	meta := buildMeta(opts, len(commits), len(branches), len(tags))

	return assembleDashboard(meta, branches, tags, authors, commits, frequency), nil
}

// assembleDashboard constructs the final DashboardData struct.
func assembleDashboard(
	meta model.DashboardMeta,
	branches []model.BranchInfo,
	tags []model.TagInfo,
	authors []model.AuthorInfo,
	commits []model.CommitInfo,
	frequency model.FrequencyData,
) model.DashboardData {

	return model.DashboardData{
		Meta:      meta,
		Branches:  branches,
		Tags:      tags,
		Authors:   authors,
		Commits:   commits,
		Frequency: frequency,
	}
}

// buildMeta constructs the metadata header for the dashboard.
func buildMeta(opts CollectOptions, totalCommits, totalBranches, totalTags int) model.DashboardMeta {
	repoName := queryRepoName(opts.RepoPath)
	remoteURL := queryRemoteURL(opts.RepoPath)
	branch := queryCurrentBranch(opts.RepoPath)

	return model.DashboardMeta{
		RepoName:      repoName,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Branch:        branch,
		RemoteURL:     remoteURL,
		TotalCommits:  totalCommits,
		TotalBranches: totalBranches,
		TotalTags:     totalTags,
		Limit:         opts.Limit,
		Since:         opts.Since,
	}
}

// queryRemoteURL returns the remote origin URL or empty string.
func queryRemoteURL(repoPath string) string {
	out, err := runDashGit(repoPath,
		constants.GitConfigCmd, constants.GitGetFlag, constants.GitRemoteOrigin)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(out)
}

// queryCurrentBranch returns the current HEAD branch name.
func queryCurrentBranch(repoPath string) string {
	out, err := runDashGit(repoPath,
		constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(out)
}

// collectCommits parses the git log output into CommitInfo slices.
func collectCommits(opts CollectOptions) ([]model.CommitInfo, error) {
	raw, err := queryLog(opts.RepoPath, opts.Limit, opts.Since, opts.NoMerges)
	if err != nil {
		return nil, err
	}

	return parseCommitLog(raw), nil
}

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

	return model.CommitInfo{
		SHA:          fields[0],
		ShortSHA:     fields[1],
		Author:       fields[2],
		Email:        fields[3],
		Date:         fields[4],
		Message:      fields[5],
		IsMerge:      isMergeCommit(fields[6]),
		FilesChanged: files,
		Insertions:   ins,
		Deletions:    del,
	}, true
}

// parseNumstat tallies file changes, insertions, and deletions.
func parseNumstat(lines []string) (int, int, int) {
	files, ins, del := 0, 0, 0

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		files++
		added, _ := strconv.Atoi(parts[0])
		removed, _ := strconv.Atoi(parts[1])
		ins += added
		del += removed
	}

	return files, ins, del
}

// collectBranches parses branch query output into BranchInfo slices.
func collectBranches(repoPath string) []model.BranchInfo {
	lines, err := queryBranches(repoPath)
	if err != nil {
		return nil
	}

	return parseBranchLines(lines)
}

// parseBranchLines converts raw branch lines to BranchInfo structs.
func parseBranchLines(lines []string) []model.BranchInfo {
	branches := make([]model.BranchInfo, 0, len(lines))

	for _, line := range lines {
		fields := strings.SplitN(line, "|", 3)
		if len(fields) < 3 {
			continue
		}

		branches = append(branches, model.BranchInfo{
			Name:           fields[0],
			IsRemote:       strings.HasPrefix(fields[0], constants.GitOriginPrefix),
			LastCommitSHA:  fields[1],
			LastCommitDate: fields[2],
		})
	}

	return branches
}

// collectTags parses tag query output into TagInfo slices.
func collectTags(repoPath string) []model.TagInfo {
	lines, err := queryTags(repoPath)
	if err != nil {
		return nil
	}

	return parseTagLines(repoPath, lines)
}

// parseTagLines converts raw tag lines to TagInfo structs with distances.
func parseTagLines(repoPath string, lines []string) []model.TagInfo {
	tags := make([]model.TagInfo, 0, len(lines))

	for i, line := range lines {
		fields := strings.SplitN(line, "|", 3)
		if len(fields) < 3 {
			continue
		}

		count := 0
		if i < len(lines)-1 {
			nextFields := strings.SplitN(lines[i+1], "|", 3)
			if len(nextFields) >= 1 {
				count = queryTagDistance(repoPath, nextFields[0], fields[0])
			}
		}

		tags = append(tags, model.TagInfo{
			Name:        fields[0],
			SHA:         fields[1],
			Date:        fields[2],
			CommitCount: count,
		})
	}

	return tags
}
