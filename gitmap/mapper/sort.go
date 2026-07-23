package mapper

// Deterministic ordering for ScanRecord slices.
//
// scanner.SortRepos already sorts the raw RepoInfo slice by
// (RelativePath, AbsolutePath) before BuildRecords sees it, but the
// user-facing contract is "sort by folder path then repo URL". URL
// is only known AFTER the mapper has resolved it from
// `git remote get-url`, so we re-sort the records here using the
// URL as the tiebreaker the spec actually documents.
//
// Sort order, in priority:
//
//   1. RelativePath          -- the folder path; primary key so the
//                                terminal/CSV/JSON output reads in
//                                directory order.
//   2. HTTPSUrl              -- documented secondary key; gives
//                                deterministic order when two repos
//                                somehow share a relative path.
//   3. SSHUrl                -- final fallback so a repo with only
//                                an SSH remote still sorts stably
//                                next to its peers.
//   4. AbsolutePath          -- last-resort total order so the sort
//                                is fully reproducible even when
//                                every URL field is empty (e.g. a
//                                repo with no remotes configured).
//
// Sort is stable so equal-key records keep the order they had after
// scanner.SortRepos -- this matters for callers that pre-sorted on a
// custom key and want the mapper to preserve it within a key bucket.

import (
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// SortRecords sorts the slice in place using the documented
// (path, URL) ordering. Exported so non-mapper callers (custom
// pipelines, tests) can apply the same canonical order.
func SortRecords(records []model.ScanRecord) {
	sort.SliceStable(records, func(i, j int) bool {
		return lessRecord(records[i], records[j])
	})
}

// lessRecord implements the four-level comparison documented at the
// top of this file. Split out so SortRecords stays a one-liner and
// each comparison step is easy to read top-to-bottom.
func lessRecord(a, b model.ScanRecord) bool {
	if a.RelativePath != b.RelativePath {
		return a.RelativePath < b.RelativePath
	}
	if a.HTTPSUrl != b.HTTPSUrl {
		return a.HTTPSUrl < b.HTTPSUrl
	}
	if a.SSHUrl != b.SSHUrl {
		return a.SSHUrl < b.SSHUrl
	}

	return a.AbsolutePath < b.AbsolutePath
}
