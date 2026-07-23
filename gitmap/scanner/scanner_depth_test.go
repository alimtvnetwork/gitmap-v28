package scanner

import (
	"testing"
)

// TestScanDirDefaultMaxDepthStopsAtFour asserts the default cap of 4
// levels: a repo at depth 4 (root/d1/d2/d3/d4) IS found, but a repo at
// depth 5 is NOT. The cap is inclusive of depth 4.
func TestScanDirDefaultMaxDepthStopsAtFour(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "d1/d2/d3/d4/in-budget")   // depth 5? d1=1 d2=2 d3=3 d4=4 in=5
	makeRepo(t, root, "d1/d2/d3/d4/d5/too-deep") // depth 6 — past cap

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	// With DefaultMaxDepth=4: root(0) → d1(1) → d2(2) → d3(3) → d4(4)
	// is the deepest dir we walk. We can READ d4's entries (find
	// "in-budget" / "d5" subdirs) but we MUST NOT enqueue them.
	// Therefore neither repo is reachable — both live at depth >=5.
	if len(got) != 0 {
		t.Fatalf("default depth cap should hide depth-5+ repos, got %+v", got)
	}
}

// TestScanDirDefaultDepthFindsRepoAtCap verifies a repo whose directory
// IS at exactly the cap depth gets discovered: the walker reads that
// dir and detects its `.git` marker before considering descent.
func TestScanDirDefaultDepthFindsRepoAtCap(t *testing.T) {
	root := t.TempDir()
	// d1/d2/d3/repo => repo dir is at depth 4 == DefaultMaxDepth.
	makeRepo(t, root, "d1/d2/d3/repo")

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("repo at cap depth should be discovered, got %+v", got)
	}
}

// TestScanDirCustomMaxDepthOne pins the cap at 1 — only direct
// children of the scan root may be inspected. A repo at depth 1 is
// found; a repo at depth 2 is not.
func TestScanDirCustomMaxDepthOne(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "shallow")      // depth 1
	makeRepo(t, root, "outer/nested") // depth 2 — should be skipped

	got, err := ScanDirWithOptions(root, ScanOptions{MaxDepth: 1})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("MaxDepth=1 should hide depth-2 repos, got %+v", got)
	}
}

// TestScanDirNegativeMaxDepthIsUnbounded confirms that callers can opt
// out of the cap entirely (legacy behavior) by passing MaxDepth < 0.
// A deeply-nested repo is found regardless of depth.
func TestScanDirNegativeMaxDepthIsUnbounded(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "a/b/c/d/e/f/g/deep") // depth 8

	got, err := ScanDirWithOptions(root, ScanOptions{MaxDepth: -1})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("MaxDepth<0 should be unbounded, got %+v", got)
	}
}

// TestScanDirRepoStopsDescentWithinBudget reasserts that the existing
// "stop at .git" rule still applies within the depth budget: a repo
// found at depth 1 means its depth-2/3/4 children are NOT walked, even
// though they would otherwise be in budget.
func TestScanDirRepoStopsDescentWithinBudget(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "outer")                       // depth 1, has .git
	makeRepo(t, root, "outer/nested/inner/sub/leaf") // would be depth 5 — also past cap, double-covered

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("nested repo under discovered repo must be hidden, got %+v", got)
	}
}
