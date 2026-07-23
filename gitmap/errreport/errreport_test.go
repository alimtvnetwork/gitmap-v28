package errreport

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestNilCollectorIsNoOp confirms that the nil-receiver short-circuit
// is fully wired: Add, Count, and WriteIfAny must all silently
// succeed without panicking when called on a nil *Collector.
func TestNilCollectorIsNoOp(t *testing.T) {
	var c *Collector
	c.Add(PhaseScan, Entry{RepoPath: "/x", Error: "boom"})
	if s, k := c.Count(); s != 0 || k != 0 {
		t.Fatalf("nil Count must be 0/0, got %d/%d", s, k)
	}
	got, err := c.WriteIfAny(t.TempDir())
	if err != nil || got != "" {
		t.Fatalf("nil WriteIfAny must be ('',nil), got (%q,%v)", got, err)
	}
}

// TestWriteIfAnySkipsCleanRun verifies that an empty collector
// produces NO file on disk — clean runs must not litter the reports
// directory with empty JSON arrays.
func TestWriteIfAnySkipsCleanRun(t *testing.T) {
	c := New("test-1.0.0", "scan")
	dir := t.TempDir()
	got, err := c.WriteIfAny(dir)
	if err != nil {
		t.Fatalf("WriteIfAny: %v", err)
	}
	if got != "" {
		t.Fatalf("empty collector must return '', got %q", got)
	}
	entries, _ := os.ReadDir(filepath.Join(dir, reportDirRel))
	if len(entries) != 0 {
		t.Fatalf("clean run must leave reports dir empty, got %d entries", len(entries))
	}
}

// TestWriteIfAnyEmitsGroupedPayload exercises the happy path: add
// failures across both phases, write, then re-parse the file and
// assert the grouped shape (meta + scan[] + clone[]) plus that
// every Entry's required fields round-trip cleanly.
func TestWriteIfAnyEmitsGroupedPayload(t *testing.T) {
	c := New("test-1.2.3", "scan+clone-next")
	c.Add(PhaseScan, Entry{RepoPath: "/r1", Step: "readdir", Error: "denied"})
	c.Add(PhaseScan, Entry{RepoPath: "/r2", Step: "ls-remote", Error: "auth"})
	c.Add(PhaseClone, Entry{RepoPath: "/r3", Step: "clone", Error: "exit 128"})

	dir := t.TempDir()
	path, err := c.WriteIfAny(dir)
	if err != nil {
		t.Fatalf("WriteIfAny: %v", err)
	}
	if path == "" {
		t.Fatal("expected non-empty path for non-empty collector")
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	var out fileShape
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Meta.Version != "test-1.2.3" || out.Meta.Command != "scan+clone-next" {
		t.Fatalf("meta mismatch: %+v", out.Meta)
	}
	if out.Meta.TotalScan != 2 || out.Meta.TotalClone != 1 {
		t.Fatalf("totals wrong: scan=%d clone=%d", out.Meta.TotalScan, out.Meta.TotalClone)
	}
	if len(out.Scan) != 2 || len(out.Clone) != 1 {
		t.Fatalf("entry counts wrong: %+v", out)
	}
	for _, e := range append(out.Scan, out.Clone...) {
		if e.TimestampUnixMS == 0 {
			t.Fatalf("entry missing timestamp: %+v", e)
		}
	}
}

// TestConcurrentAdd asserts the collector survives many goroutines
// piling on simultaneously without dropping or duplicating entries.
// Race-detector under `go test -race` catches any unprotected map
// access; we additionally assert the final length matches.
func TestConcurrentAdd(t *testing.T) {
	c := New("v", "scan")
	const writers, perWriter = 16, 32
	var wg sync.WaitGroup
	wg.Add(writers)
	for w := 0; w < writers; w++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perWriter; i++ {
				c.Add(PhaseScan, Entry{RepoPath: "/r", Error: "x"})
			}
		}(w)
	}
	wg.Wait()
	gotScan, _ := c.Count()
	if gotScan != writers*perWriter {
		t.Fatalf("expected %d scan entries, got %d", writers*perWriter, gotScan)
	}
}
