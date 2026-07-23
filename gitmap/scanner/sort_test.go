package scanner

// Tests for the deterministic-sort contract. We don't rely on the
// real walker (which is already deterministic at the FS level for
// the toy fixtures the other tests use); instead we feed SortRepos
// hand-crafted slices so the assertions pin EXACTLY the documented
// key order without depending on filesystem layout.

import (
	"testing"
)

// TestSortRepos_PathPrimaryKey verifies that RelativePath is the
// primary sort key: shuffled input must come back lexicographically
// ordered by path.
func TestSortRepos_PathPrimaryKey(t *testing.T) {
	in := []RepoInfo{
		{RelativePath: "z/last", AbsolutePath: "/abs/z/last"},
		{RelativePath: "a/first", AbsolutePath: "/abs/a/first"},
		{RelativePath: "m/middle", AbsolutePath: "/abs/m/middle"},
	}
	SortRepos(in)
	want := []string{"a/first", "m/middle", "z/last"}
	for i, w := range want {
		if in[i].RelativePath != w {
			t.Errorf("idx %d: got %q want %q", i, in[i].RelativePath, w)
		}
	}
}

// TestSortRepos_AbsPathTiebreaker verifies that when two entries
// share a RelativePath the AbsolutePath decides ordering, so the
// sort is total even in the (unusual) duplicate-path case.
func TestSortRepos_AbsPathTiebreaker(t *testing.T) {
	in := []RepoInfo{
		{RelativePath: "same", AbsolutePath: "/b/abs"},
		{RelativePath: "same", AbsolutePath: "/a/abs"},
	}
	SortRepos(in)
	if in[0].AbsolutePath != "/a/abs" {
		t.Errorf("tiebreaker broken: %+v", in)
	}
}

// TestSortRepos_StableForEqualKeys verifies sort stability: two
// entries with identical RelativePath + AbsolutePath keep their
// original order so callers that pre-sorted on a custom key see
// their tiebreaker preserved.
func TestSortRepos_StableForEqualKeys(t *testing.T) {
	in := []RepoInfo{
		{RelativePath: "x", AbsolutePath: "/x", Depth: 1},
		{RelativePath: "x", AbsolutePath: "/x", Depth: 2},
	}
	SortRepos(in)
	if in[0].Depth != 1 || in[1].Depth != 2 {
		t.Errorf("sort not stable: %+v", in)
	}
}
