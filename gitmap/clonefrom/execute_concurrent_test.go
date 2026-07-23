package clonefrom

// Concurrent runner contract tests for clonefrom.
//
// We exercise the worker pool WITHOUT a real `git` binary by
// pre-creating each row's destination as a non-empty directory.
// executeRow's shouldSkip short-circuits these as `skipped`, so
// the network/exec layer is bypassed but the pool's ordering,
// hook-firing, and progress-emission contracts are still exercised.
//
// The actual `git clone` shell-out is covered by the existing
// sequential tests; the concurrent runner reuses executeRow so
// re-testing network behavior here would just duplicate them.

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// seedSkippableRows pre-creates each row's destination so executeRow
// reports `skipped` instead of trying to shell out to git. Returns
// the rows in input order.
func seedSkippableRows(t *testing.T, dir string, n int) []Row {
	t.Helper()
	rows := make([]Row, n)
	for i := 0; i < n; i++ {
		dest := filepath.Join(dir, padIdx(i))
		if err := os.MkdirAll(dest, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dest, err)
		}
		marker := filepath.Join(dest, ".keep")
		if err := os.WriteFile(marker, []byte("x"), 0o644); err != nil {
			t.Fatalf("write %s: %v", marker, err)
		}
		rows[i] = Row{URL: "https://example.invalid/" + padIdx(i) + ".git", Dest: padIdx(i)}
	}

	return rows
}

func TestExecuteWithHooksConcurrent_FallbackBelowTwoWorkers(t *testing.T) {
	dir := t.TempDir()
	rows := seedSkippableRows(t, dir, 1)
	plan := Plan{Rows: rows}
	res := ExecuteWithHooksConcurrent(plan, dir, io.Discard, nil, 1)
	if len(res) != 1 {
		t.Fatalf("len(res)=%d, want 1", len(res))
	}
	if res[0].Status != constants.CloneFromStatusSkipped {
		t.Errorf("status=%q, want %q (pre-populated dest must skip)",
			res[0].Status, constants.CloneFromStatusSkipped)
	}
}

func TestExecuteWithHooksConcurrent_PreservesInputOrder(t *testing.T) {
	dir := t.TempDir()
	rows := seedSkippableRows(t, dir, 8)
	plan := Plan{Rows: rows}
	res := ExecuteWithHooksConcurrent(plan, dir, io.Discard, nil, 4)
	if len(res) != len(rows) {
		t.Fatalf("len(res)=%d, want %d", len(res), len(rows))
	}
	for i, r := range res {
		if r.Row.URL != rows[i].URL {
			t.Errorf("res[%d].URL=%q, want %q (order drift)",
				i, r.Row.URL, rows[i].URL)
		}
	}
}

func TestExecuteWithHooksConcurrent_HookOrderAndCount(t *testing.T) {
	dir := t.TempDir()
	rows := seedSkippableRows(t, dir, 3)
	plan := Plan{Rows: rows}

	var mu sync.Mutex
	var seen []string
	hook := func(index, total int, row Row, dest string) {
		mu.Lock()
		defer mu.Unlock()
		seen = append(seen, row.URL)
		if total != len(rows) {
			t.Errorf("hook total=%d, want %d", total, len(rows))
		}
	}
	_ = ExecuteWithHooksConcurrent(plan, dir, io.Discard, hook, 4)
	if len(seen) != len(rows) {
		t.Fatalf("hook fired %d times, want %d", len(seen), len(rows))
	}
	for i, r := range rows {
		if seen[i] != r.URL {
			t.Errorf("hook[%d]=%q, want %q (input order required)",
				i, seen[i], r.URL)
		}
	}
}

func TestExecuteWithHooksConcurrent_ProgressLinesEmitted(t *testing.T) {
	dir := t.TempDir()
	rows := seedSkippableRows(t, dir, 3)
	plan := Plan{Rows: rows}
	var buf bytes.Buffer
	_ = ExecuteWithHooksConcurrent(plan, dir, &buf, nil, 4)
	got := buf.Bytes()
	for i, r := range rows {
		if !bytes.Contains(got, []byte(r.URL)) {
			t.Errorf("progress missing row %d URL %q in %q", i, r.URL, got)
		}
	}
}

// padIdx produces a stable, sortable Dest path for ordering
// assertions.
func padIdx(i int) string {
	return string(rune('a'+i)) + "-row"
}
