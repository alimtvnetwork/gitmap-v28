package committransfer

import (
	"testing"
	"time"
)

// TestBuildInterleavedStreamSortsByAuthorDate proves the merged stream
// from both directional plans ends up in chronological order regardless
// of which plan it came from. This is the core invariant of --interleave.
func TestBuildInterleavedStreamSortsByAuthorDate(t *testing.T) {
	t1 := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
	t4 := time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC)

	ltr := ReplayPlan{Commits: []SourceCommit{
		{ShortSHA: "L1", AuthorAt: t1},
		{ShortSHA: "L3", AuthorAt: t3},
	}}
	rtl := ReplayPlan{Commits: []SourceCommit{
		{ShortSHA: "R2", AuthorAt: t2},
		{ShortSHA: "R4", AuthorAt: t4},
	}}

	stream := buildInterleavedStream(ltr, rtl)

	wantOrder := []string{"L1", "R2", "L3", "R4"}
	wantDirs := []string{"L→R", "R→L", "L→R", "R→L"}

	if len(stream) != len(wantOrder) {
		t.Fatalf("expected %d steps; got %d", len(wantOrder), len(stream))
	}
	for i, step := range stream {
		if step.Commit.ShortSHA != wantOrder[i] {
			t.Errorf("step %d: want SHA %s; got %s", i, wantOrder[i], step.Commit.ShortSHA)
		}
		if step.Direction != wantDirs[i] {
			t.Errorf("step %d: want dir %s; got %s", i, wantDirs[i], step.Direction)
		}
	}
}

// TestBuildInterleavedStreamStableForTies guarantees within-side order
// is preserved when timestamps tie (sort.SliceStable). Without this,
// `git rev-list --reverse` ordering could be silently shuffled.
func TestBuildInterleavedStreamStableForTies(t *testing.T) {
	tie := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	ltr := ReplayPlan{Commits: []SourceCommit{
		{ShortSHA: "L_a", AuthorAt: tie},
		{ShortSHA: "L_b", AuthorAt: tie},
	}}
	rtl := ReplayPlan{Commits: []SourceCommit{
		{ShortSHA: "R_a", AuthorAt: tie},
	}}

	stream := buildInterleavedStream(ltr, rtl)

	// L→R entries appended first → must remain before R→L on tie.
	want := []string{"L_a", "L_b", "R_a"}
	for i, s := range stream {
		if s.Commit.ShortSHA != want[i] {
			t.Errorf("step %d: want %s; got %s", i, want[i], s.Commit.ShortSHA)
		}
	}
}

// TestBuildInterleavedStreamEmptyPlans returns empty when both sides
// have nothing — RunBothInterleaved short-circuits on this.
func TestBuildInterleavedStreamEmptyPlans(t *testing.T) {
	stream := buildInterleavedStream(ReplayPlan{}, ReplayPlan{})
	if len(stream) != 0 {
		t.Fatalf("expected empty stream; got %d entries", len(stream))
	}
}
