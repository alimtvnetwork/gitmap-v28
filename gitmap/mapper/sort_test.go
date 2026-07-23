package mapper

// Tests for SortRecords -- pins the (RelativePath, HTTPSUrl, SSHUrl,
// AbsolutePath) ordering so renderers downstream can rely on the
// sequence being stable across runs and platforms.

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// TestSortRecords_PathPrimaryKey pins the documented primary key:
// shuffled records come back ordered by RelativePath.
func TestSortRecords_PathPrimaryKey(t *testing.T) {
	in := []model.ScanRecord{
		{RelativePath: "z", HTTPSUrl: "https://x/z.git"},
		{RelativePath: "a", HTTPSUrl: "https://x/a.git"},
		{RelativePath: "m", HTTPSUrl: "https://x/m.git"},
	}
	SortRecords(in)
	want := []string{"a", "m", "z"}
	for i, w := range want {
		if in[i].RelativePath != w {
			t.Errorf("idx %d: got %q want %q", i, in[i].RelativePath, w)
		}
	}
}

// TestSortRecords_HTTPSURLTiebreaker pins the documented secondary
// key: same path, HTTPSUrl decides.
func TestSortRecords_HTTPSURLTiebreaker(t *testing.T) {
	in := []model.ScanRecord{
		{RelativePath: "same", HTTPSUrl: "https://x/b.git"},
		{RelativePath: "same", HTTPSUrl: "https://x/a.git"},
	}
	SortRecords(in)
	if in[0].HTTPSUrl != "https://x/a.git" {
		t.Errorf("HTTPSUrl tiebreaker broken: %+v", in)
	}
}

// TestSortRecords_SSHFallback verifies that when HTTPSUrl is empty
// on both sides, SSHUrl is used as the next tiebreaker.
func TestSortRecords_SSHFallback(t *testing.T) {
	in := []model.ScanRecord{
		{RelativePath: "same", SSHUrl: "git@x:b.git"},
		{RelativePath: "same", SSHUrl: "git@x:a.git"},
	}
	SortRecords(in)
	if in[0].SSHUrl != "git@x:a.git" {
		t.Errorf("SSHUrl fallback broken: %+v", in)
	}
}

// TestSortRecords_AbsPathFinal verifies the final fallback when
// every URL field is empty (e.g. a repo with no remotes).
func TestSortRecords_AbsPathFinal(t *testing.T) {
	in := []model.ScanRecord{
		{RelativePath: "same", AbsolutePath: "/b"},
		{RelativePath: "same", AbsolutePath: "/a"},
	}
	SortRecords(in)
	if in[0].AbsolutePath != "/a" {
		t.Errorf("AbsolutePath final fallback broken: %+v", in)
	}
}
