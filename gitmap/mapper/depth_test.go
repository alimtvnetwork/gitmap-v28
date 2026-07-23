package mapper

// Verifies that scanner.RepoInfo.Depth survives the mapper layer
// and lands in model.ScanRecord.Depth. Lightweight unit test —
// avoids the filesystem entirely by constructing a synthetic
// RepoInfo. The end-to-end "depth captured by the walker" path is
// covered separately in scanner/scanner_depth_test.go.

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/scanner"
)

// TestBuildRecords_PropagatesDepth confirms that depth values
// flowing in via RepoInfo are surfaced unchanged on every emitted
// ScanRecord. Two repos at different depths catch any accidental
// "all rows get the same depth" regression where a shared variable
// would clobber per-record state.
func TestBuildRecords_PropagatesDepth(t *testing.T) {
	repos := []scanner.RepoInfo{
		{AbsolutePath: "/x/a", RelativePath: "a", Depth: 1},
		{AbsolutePath: "/x/deep/b/c/d", RelativePath: "deep/b/c/d", Depth: 4},
	}
	records := BuildRecords(repos, "https", "")
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Depth != 1 {
		t.Errorf("record[0].Depth: got %d, want 1", records[0].Depth)
	}
	if records[1].Depth != 4 {
		t.Errorf("record[1].Depth: got %d, want 4", records[1].Depth)
	}
}
