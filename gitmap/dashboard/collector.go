// Package dashboard collects Git repository data for the HTML dashboard.
package dashboard

import (
	"fmt"
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
	Recent   bool
}

// Collect gathers all repository data into a single DashboardData struct.
func Collect(opts CollectOptions) (model.DashboardData, error) {
	opts = normalizeOptions(opts)
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

// normalizeOptions applies dynamic boundaries like Recent window.
func normalizeOptions(opts CollectOptions) CollectOptions {
	if opts.Recent && len(opts.Since) == 0 {
		opts.Since = recentSinceDate()
	}
	return opts
}

// recentSinceDate returns the date boundary for 7 days ago in UTC formatted as YYYY-MM-DD.
func recentSinceDate() string {
	return time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02")
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
func buildMeta(
	opts CollectOptions,
	totalCommits,
	totalBranches,
	totalTags int,
) model.DashboardMeta {
	return model.DashboardMeta{
		RepoName:      queryRepoName(opts.RepoPath),
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Branch:        queryCurrentBranch(opts.RepoPath),
		RemoteURL:     queryRemoteURL(opts.RepoPath),
		TotalCommits:  totalCommits,
		TotalBranches: totalBranches,
		TotalTags:     totalTags,
		Limit:         opts.Limit,
		Since:         opts.Since,
		Recent:        opts.Recent,
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

// collectBranches parses branch query output into BranchInfo slices.
func collectBranches(repoPath string) []model.BranchInfo {
	lines, err := queryBranches(repoPath)
	if err != nil {
		return nil
	}
	return parseBranchLines(lines)
}

// collectTags parses tag query output into TagInfo slices.
func collectTags(repoPath string) []model.TagInfo {
	lines, err := queryTags(repoPath)
	if err != nil {
		return nil
	}
	return parseTagLines(repoPath, lines)
}
