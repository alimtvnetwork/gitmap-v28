package scanner

// Deterministic ordering for discovered repos.
//
// The walker runs a parallel worker pool, so the order in which
// `processDir` enqueues RepoInfo onto `scanState.repos` depends on
// goroutine scheduling and is not reproducible across runs. That
// non-determinism leaks all the way out to the terminal table, the
// CSV/JSON exports, and the generated clone scripts, so the same
// scan looks "different" each time it runs.
//
// SortRepos pins a stable, lexicographic order so every downstream
// renderer produces byte-identical output across runs:
//
//   - Primary key: RelativePath (the folder path relative to the
//     scan root). This is what users actually see in the terminal
//     tree and the CSV's `relativePath` column, so sorting on it
//     gives a layout that reads top-to-bottom like a directory tree.
//   - Tiebreaker:  AbsolutePath. Two RepoInfo entries can in
//     principle share a RelativePath (e.g. when two scan roots get
//     stitched together by an upstream caller), so we fall back to
//     the absolute path to keep the order total and deterministic
//     even in that pathological case.
//
// Sorting is done once at the very end of walkParallel, so the
// O(n log n) cost is paid once per scan and not per renderer.

import "sort"

// SortRepos sorts the slice in place by (RelativePath, AbsolutePath)
// using filepath-agnostic byte comparison. Exported so callers that
// build RepoInfo from a non-walker source (tests, custom adapters)
// can apply the same canonical order before handing the slice to a
// renderer.
func SortRepos(repos []RepoInfo) {
	sort.SliceStable(repos, func(i, j int) bool {
		if repos[i].RelativePath != repos[j].RelativePath {
			return repos[i].RelativePath < repos[j].RelativePath
		}

		return repos[i].AbsolutePath < repos[j].AbsolutePath
	})
}
