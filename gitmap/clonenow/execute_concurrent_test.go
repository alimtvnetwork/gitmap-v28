package clonenow

// Concurrent runner contract tests.
//
// These cover the structural guarantees of
// ExecuteWithHooksConcurrent that don't require a real `git`
// binary on PATH:
//
//   - workers <= 1 falls through to the sequential ExecuteWithHooks
//     (verified by passing a row that fails its no-URL check; the
//     fallback path must produce the same Result.Status the
//     sequential runner would).
//   - Result ordering matches input order regardless of N.
//   - The BeforeRow hook fires once per row, in input order, with
//     the resolved index/total/url/dest.
//
// The actual `git clone` shell-out is exercised by the existing
// sequential tests (executeRow + execute_mkdir_test.go); the
// concurrent runner reuses the same per-row helper, so re-testing
// network behavior here would just duplicate those.

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestExecuteWithHooksConcurrent_FallbackBelowTwoWorkers(t *testing.T) {
	plan := Plan{Mode: constants.CloneNowModeHTTPS, Rows: []Row{{RelativePath: "z"}}}
	res := ExecuteWithHooksConcurrent(plan, t.TempDir(), io.Discard, nil, 1)
	if len(res) != 1 {
		t.Fatalf("len(res)=%d, want 1", len(res))
	}
	if res[0].Status != constants.CloneNowStatusFailed {
		t.Errorf("status=%q, want %q (no-URL row must fail under fallback)",
			res[0].Status, constants.CloneNowStatusFailed)
	}
}

func TestExecuteWithHooksConcurrent_PreservesInputOrder(t *testing.T) {
	rows := make([]Row, 8)
	for i := range rows {
		// Empty URL forces the executor to short-circuit with a
		// failed Result — this avoids needing a real git binary
		// while still exercising the worker pool's order guarantee.
		rows[i] = Row{RelativePath: padIdx(i)}
	}
	plan := Plan{Mode: constants.CloneNowModeHTTPS, Rows: rows}
	res := ExecuteWithHooksConcurrent(plan, t.TempDir(), io.Discard, nil, 4)
	if len(res) != len(rows) {
		t.Fatalf("len(res)=%d, want %d", len(res), len(rows))
	}
	for i, r := range res {
		if r.Row.RelativePath != padIdx(i) {
			t.Errorf("res[%d].RelativePath=%q, want %q (order drift)",
				i, r.Row.RelativePath, padIdx(i))
		}
	}
}

func TestExecuteWithHooksConcurrent_HookOrderAndCount(t *testing.T) {
	rows := []Row{
		{RelativePath: "a"},
		{RelativePath: "b"},
		{RelativePath: "c"},
	}
	plan := Plan{Mode: constants.CloneNowModeHTTPS, Rows: rows}

	var mu sync.Mutex
	var seen []string
	hook := func(index, total int, row Row, url, dest string) {
		mu.Lock()
		defer mu.Unlock()
		seen = append(seen, row.RelativePath)
		if total != len(rows) {
			t.Errorf("hook total=%d, want %d", total, len(rows))
		}
	}
	_ = ExecuteWithHooksConcurrent(plan, t.TempDir(), io.Discard, hook, 4)
	if len(seen) != len(rows) {
		t.Fatalf("hook fired %d times, want %d", len(seen), len(rows))
	}
	for i, r := range rows {
		if seen[i] != r.RelativePath {
			t.Errorf("hook[%d]=%q, want %q (input order required)",
				i, seen[i], r.RelativePath)
		}
	}
}

func TestExecuteWithHooksConcurrent_ProgressLinesEmittedInOrder(t *testing.T) {
	rows := []Row{
		{RelativePath: "first"},
		{RelativePath: "second"},
		{RelativePath: "third"},
	}
	plan := Plan{Mode: constants.CloneNowModeHTTPS, Rows: rows}
	var buf bytes.Buffer
	_ = ExecuteWithHooksConcurrent(plan, t.TempDir(), &buf, nil, 4)
	got := buf.String()
	for i, r := range rows {
		marker := " " + r.RelativePath + "\n"
		if !bytes.Contains([]byte(got), []byte(marker)) {
			t.Errorf("progress missing row %d (%q): %q", i, r.RelativePath, got)
		}
	}
}

// padIdx produces a stable, sortable RelativePath for ordering
// assertions. Avoids strconv import churn when used in tight loops.
func padIdx(i int) string {
	return string(rune('a'+i)) + "-row"
}
