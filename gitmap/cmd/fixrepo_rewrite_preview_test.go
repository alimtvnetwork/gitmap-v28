package cmd

// v5.40.0: dry-run preview breakdown coverage. Asserts that the
// preview engine reports per-rule hit counts (including the bare-base
// sentinel) without mutating its input and matches the rewrite
// engine's total count byte-for-byte.

import "testing"

func TestPreviewAllTargets_NumberedRulesOnly(t *testing.T) {
	body := "acme-v1 and acme-v2 and acme-v2 and acme-v4"
	total, hits := previewAllTargets(body, "acme", 4, []int{1, 2, 3}, false)
	if total != 3 {
		t.Fatalf("total: got %d want 3", total)
	}
	want := map[int]int{1: 1, 2: 2}
	for _, h := range hits {
		if want[h.n] != h.count {
			t.Errorf("hit %+v not in want %v", h, want)
		}
		delete(want, h.n)
	}
	if len(want) != 0 {
		t.Fatalf("missing hits: %v", want)
	}
}

func TestPreviewAllTargets_BareBaseSweepAtV2(t *testing.T) {
	body := "acme bare and acme-v1 numbered"
	total, hits := previewAllTargets(body, "acme", 2, []int{1}, false)
	if total != 2 {
		t.Fatalf("total: got %d want 2", total)
	}
	// Expect one numbered hit (v1×1) + one bare-base hit (bare×1).
	sawBare := false
	sawV1 := false
	for _, h := range hits {
		if h.n == fixRepoBareBaseSentinel && h.count == 1 {
			sawBare = true
		}
		if h.n == 1 && h.count == 1 {
			sawV1 = true
		}
	}
	if !sawBare || !sawV1 {
		t.Fatalf("missing hit categories: hits=%v", hits)
	}
}

func TestPreviewAllTargets_RestrictSuppressesBareBase(t *testing.T) {
	body := "acme bare and acme-v1 numbered"
	total, hits := previewAllTargets(body, "acme", 2, []int{1}, true)
	if total != 1 {
		t.Fatalf("total: got %d want 1 (bare suppressed)", total)
	}
	for _, h := range hits {
		if h.n == fixRepoBareBaseSentinel {
			t.Fatalf("bare sentinel leaked despite restrictNoVersion: %v", hits)
		}
	}
}

func TestPreviewAllTargets_TotalMatchesRewriteEngine(t *testing.T) {
	// The preview engine MUST agree with applyAllTargetsR on the
	// total count for every (current, targets) shape — otherwise the
	// dry-run summary lies about what a real run would do.
	body := "acme and acme-v1 plus acme-v2 plus acme-v3"
	cases := []struct {
		current int
		targets []int
		restr   bool
	}{
		{2, []int{1}, false},
		{2, []int{1}, true},
		{3, []int{1, 2}, false},
		{4, []int{1, 2, 3}, false},
	}
	for _, c := range cases {
		_, rw := applyAllTargetsR(body, "acme", c.current, c.targets, c.restr)
		prev, _ := previewAllTargets(body, "acme", c.current, c.targets, c.restr)
		if rw != prev {
			t.Errorf("count drift (current=%d restr=%v): rewrite=%d preview=%d",
				c.current, c.restr, rw, prev)
		}
	}
}
