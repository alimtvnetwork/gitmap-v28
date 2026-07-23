package cmd

// Tests for the regoldens diff-summary helpers. These tests target
// pure-function helpers that don't shell out, so they remain fast
// and hermetic. The git-invoking layer (collectGoldenDiffEntries)
// is exercised end-to-end by the regoldens CLI in CI.

import (
	"strings"
	"testing"
)

func TestNormalizePorcelainStatus_Untracked_ReturnsAdded(t *testing.T) {
	if got := normalizePorcelainStatus("??"); got != "A" {
		t.Fatalf("?? -> %q, want %q", got, "A")
	}
}

func TestNormalizePorcelainStatus_DeletionBeatsModification(t *testing.T) {
	if got := normalizePorcelainStatus("MD"); got != "D" {
		t.Fatalf("MD -> %q, want %q", got, "D")
	}
}

func TestNormalizePorcelainStatus_PlainModified(t *testing.T) {
	if got := normalizePorcelainStatus(" M"); got != "M" {
		t.Fatalf(" M -> %q, want %q", got, "M")
	}
}

func TestSortGoldenDiffEntries_OrdersByStatusThenPath(t *testing.T) {
	in := []goldenDiffEntry{
		{status: "D", path: "z/testdata/x"},
		{status: "A", path: "b/testdata/y"},
		{status: "M", path: "a/testdata/y"},
		{status: "A", path: "a/testdata/x"},
	}
	sortGoldenDiffEntries(in)
	wantOrder := []string{
		"a/testdata/x", "b/testdata/y", "a/testdata/y", "z/testdata/x",
	}
	for i, e := range in {
		if e.path != wantOrder[i] {
			t.Fatalf("position %d: got %q, want %q", i, e.path, wantOrder[i])
		}
	}
}

func TestGoldenDiffTotals_AccumulateMixedStatuses(t *testing.T) {
	totals := goldenDiffTotals{count: 4}
	totals.accumulate(goldenDiffEntry{status: "A", added: 10, deleted: 0})
	totals.accumulate(goldenDiffEntry{status: "M", added: 4, deleted: 7})
	totals.accumulate(goldenDiffEntry{status: "D", added: 0, deleted: 12})
	totals.accumulate(goldenDiffEntry{status: "R", added: 1, deleted: 1})
	assertTotalsEqual(t, totals, 1, 1, 1, 1, 15, 20)
}

// assertTotalsEqual centralizes the multi-field comparison so the
// test bodies above stay one-purpose and well under 15 lines.
func assertTotalsEqual(t *testing.T, totals goldenDiffTotals,
	wantAdded, wantModified, wantRenamed, wantDeleted, wantLinesAdded, wantLinesDeleted int,
) {
	t.Helper()
	mismatches := []string{}
	if totals.added != wantAdded {
		mismatches = append(mismatches, "added")
	}
	if totals.modified != wantModified {
		mismatches = append(mismatches, "modified")
	}
	if totals.renamed != wantRenamed {
		mismatches = append(mismatches, "renamed")
	}
	if totals.deleted != wantDeleted {
		mismatches = append(mismatches, "deleted")
	}
	if totals.linesAdded != wantLinesAdded {
		mismatches = append(mismatches, "linesAdded")
	}
	if totals.linesDeleted != wantLinesDeleted {
		mismatches = append(mismatches, "linesDeleted")
	}
	if len(mismatches) > 0 {
		t.Fatalf("totals mismatch in fields [%s]: got %+v",
			strings.Join(mismatches, ", "), totals)
	}
}

func TestGoldenDiffStatusRank_OrderingContract(t *testing.T) {
	ranks := []int{
		goldenDiffStatusRank("A"),
		goldenDiffStatusRank("M"),
		goldenDiffStatusRank("R"),
		goldenDiffStatusRank("D"),
	}
	for i := 1; i < len(ranks); i++ {
		if ranks[i-1] >= ranks[i] {
			t.Fatalf("rank ordering broken at %d: %v", i, ranks)
		}
	}
}

func TestIsGoldenFixturePath_AcceptsTestdataGoldenBasename(t *testing.T) {
	cases := []string{
		"gitmap/clonefrom/testdata/clone_report.golden.json",
		"pkg/x/testdata/sub/Foo_Golden.txt",
		"a/testdata/snapshot.golden",
	}
	for _, p := range cases {
		if !isGoldenFixturePath(p) {
			t.Fatalf("expected accept, got reject for %q", p)
		}
	}
}

func TestIsGoldenFixturePath_RejectsNonGoldenOrNonTestdata(t *testing.T) {
	cases := []string{
		"gitmap/clonefrom/testdata/corpus.json",      // no "golden"
		"gitmap/clonefrom/testdata/sample_input.txt", // no "golden"
		"src/golden/main.go",                         // not under testdata/
		"docs/golden-rules.md",                       // not under testdata/
		"testdata-not-a-dir/foo.golden.json",         // wrong segment
	}
	for _, p := range cases {
		if isGoldenFixturePath(p) {
			t.Fatalf("expected reject, got accept for %q", p)
		}
	}
}
