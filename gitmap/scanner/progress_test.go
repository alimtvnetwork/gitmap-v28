package scanner

import (
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestScanProgressFinalSnapshotMatchesTotals verifies that the very last
// snapshot delivered to the caller carries Final=true, the right repo
// count, and a directory count that's at least the number of dirs
// physically created (root + subdirs + repo dirs themselves all get
// read).
func TestScanProgressFinalSnapshotMatchesTotals(t *testing.T) {
	root := t.TempDir()
	want := []string{
		"a",
		"b",
		"deep/nested/c",
		"side/d",
		"side/sub/e",
	}
	for _, r := range want {
		makeRepo(t, root, r)
	}

	var (
		mu         sync.Mutex
		snapshots  []ScanProgress
		finalCount int
	)

	got, err := ScanDirWithOptions(root, ScanOptions{
		Progress: func(p ScanProgress) {
			mu.Lock()
			defer mu.Unlock()
			snapshots = append(snapshots, p)
			if p.Final {
				finalCount++
			}
		},
	})
	if err != nil {
		t.Fatalf("ScanDirWithOptions: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(snapshots) == 0 {
		t.Fatal("expected at least one progress snapshot, got 0")
	}
	if finalCount != 1 {
		t.Fatalf("expected exactly 1 Final snapshot, got %d", finalCount)
	}

	last := snapshots[len(snapshots)-1]
	if !last.Final {
		t.Fatalf("last snapshot must be Final, got %+v", last)
	}
	if last.ReposFound != int64(len(got)) {
		t.Errorf("ReposFound: got %d, want %d", last.ReposFound, len(got))
	}
	if last.DirsWalked < int64(len(want)) {
		t.Errorf("DirsWalked: got %d, want >= %d", last.DirsWalked, len(want))
	}
}

// TestScanProgressMonotonic verifies that DirsWalked / ReposFound never
// decrease across consecutive snapshots — a bug we'd otherwise risk if
// counters were ever re-initialized mid-run or read non-atomically into
// a mutated buffer.
func TestScanProgressMonotonic(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 20; i++ {
		makeRepo(t, root, filepath.Join("g", string(rune('a'+i%5)), "r"+string(rune('0'+i%10))))
	}

	var (
		mu        sync.Mutex
		snapshots []ScanProgress
	)

	if _, err := ScanDirWithOptions(root, ScanOptions{
		Progress: func(p ScanProgress) {
			mu.Lock()
			snapshots = append(snapshots, p)
			mu.Unlock()
		},
	}); err != nil {
		t.Fatalf("ScanDirWithOptions: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	for i := 1; i < len(snapshots); i++ {
		prev, cur := snapshots[i-1], snapshots[i]
		if cur.DirsWalked < prev.DirsWalked {
			t.Errorf("DirsWalked regressed at i=%d: %d -> %d", i, prev.DirsWalked, cur.DirsWalked)
		}
		if cur.ReposFound < prev.ReposFound {
			t.Errorf("ReposFound regressed at i=%d: %d -> %d", i, prev.ReposFound, cur.ReposFound)
		}
	}
}

// TestScanProgressNilCallbackDoesNotCrash guarantees the legacy code
// path — no Progress hook supplied — still works. Catches regressions
// where startProgress is invoked unconditionally.
func TestScanProgressNilCallbackDoesNotCrash(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "solo")

	got, err := ScanDirWithOptions(root, ScanOptions{Progress: nil})
	if err != nil {
		t.Fatalf("ScanDirWithOptions: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(got))
	}
}

// TestScanProgressFinalAlwaysFiresEvenOnFastScans is a regression guard
// for the contract in progress.go: even if the scan finishes inside a
// single tick interval (so the periodic loop never runs) the caller
// must still receive exactly one Final=true snapshot. We give an empty
// directory — the scan completes in microseconds, well under the
// 100ms tick.
func TestScanProgressFinalAlwaysFiresEvenOnFastScans(t *testing.T) {
	root := t.TempDir()

	var (
		gotFinal atomic.Bool
		done     = make(chan struct{})
	)

	if _, err := ScanDirWithOptions(root, ScanOptions{
		Progress: func(p ScanProgress) {
			if p.Final {
				gotFinal.Store(true)
				close(done)
			}
		},
	}); err != nil {
		t.Fatalf("ScanDirWithOptions: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Final snapshot never delivered within 2s")
	}

	if !gotFinal.Load() {
		t.Fatal("expected Final=true snapshot to fire even on empty scans")
	}
}
