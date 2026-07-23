package cmd

// Locks the guard-aware scanner contract so the rewriter and any
// downstream consumer (e.g. TestFixRepoRewriteV9ToV12Fixture) cannot
// silently disagree about what counts as a stale `{base}-vN` token.
//
// IMPORTANT: this file MUST use a synthetic base (e.g. `acme-v9`) that is
// NOT the project's own module suffix. The repo has been renamed to
// `gitmap-v27` in the past, which collapses any test using `gitmap-vN`
// into "token == its own guarded neighbor" and silently breaks the
// negative-lookahead logic. See mem://FIX-REPO DIGIT-CAPTURE GAP and
// .lovable/memory/issues/2026-05-01-fixrepo-digit-capture-desync.md.

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestScanUnguardedTokenHits(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		token     string
		wantHits  []int
		wantCount int
	}{
		{
			name:      "single match mid-line",
			body:      "use acme-v9 here",
			token:     "acme-v9",
			wantHits:  []int{4},
			wantCount: 1,
		},
		{
			name:     "guarded by trailing digit (-v9 inside -v10)",
			body:     "import acme-v10 // not v9",
			token:    "acme-v9",
			wantHits: nil, wantCount: 0,
		},
		{
			name:      "EOF-adjacent counts as unguarded",
			body:      "tail acme-v9",
			token:     "acme-v9",
			wantHits:  []int{5},
			wantCount: 1,
		},
		{
			name:      "mixed guarded + unguarded in one body",
			body:      "a acme-v9 b acme-v10 c acme-v9\n",
			token:     "acme-v9",
			wantHits:  []int{2, 23},
			wantCount: 2,
		},
		{
			name:      "non-digit neighbor (letter) is unguarded",
			body:      "acme-v9z",
			token:     "acme-v9",
			wantHits:  []int{0},
			wantCount: 1,
		},
		{
			name:      "empty token returns nothing",
			body:      "anything",
			token:     "",
			wantHits:  nil,
			wantCount: 0,
		},
		{
			name:      "token longer than body returns nothing",
			body:      "x",
			token:     "acme-v9",
			wantHits:  nil,
			wantCount: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotHits := ScanUnguardedTokenHits(tc.body, tc.token)
			if !reflect.DeepEqual(gotHits, tc.wantHits) {
				t.Errorf("hits = %v, want %v", gotHits, tc.wantHits)
			}
			if got := CountUnguardedTokenHits(tc.body, tc.token); got != tc.wantCount {
				t.Errorf("count = %d, want %d", got, tc.wantCount)
			}
		})
	}
}

// TestScannerMatchesRewriter is the cross-layer guard: for a body the
// rewriter would touch, the scanner's unguarded count MUST equal the
// rewriter's substitution count. Locks the invariant that powers
// assertDashFormBumped's `wantCount` derivation.
func TestScannerMatchesRewriter(t *testing.T) {
	// Body mixes 4 unguarded `acme-v9` tokens with one guarded `acme-v10`
	// neighbor; rewriter must touch the 4 and leave the neighbor alone.
	body := "acme-v9 + acme-v9 + acme-v9 (eof)acme-v9 // keep acme-v10"
	const (
		base    = "acme"
		target  = 9
		current = 12
	)
	token := "acme-v9"
	want := CountUnguardedTokenHits(body, token)
	out, count := applyAllTargets(body, base, current, []int{target})
	if count != want {
		t.Errorf("rewriter substituted %d, scanner counted %d (must agree)",
			count, want)
	}
	// Derive the expected rewritten token from `current` rather than
	// hard-coding a sibling literal. See mem://FIX-REPO DIGIT-CAPTURE GAP:
	// any version-bearing expectation must be built from the same int the
	// rewriter received, otherwise width-crossing bumps silently desync.
	wantToken := fmt.Sprintf("%s-v%d", base, current)
	if strings.Count(out, wantToken) != want {
		t.Errorf("output has %d %s tokens, want %d",
			strings.Count(out, wantToken), wantToken, want)
	}
	// guarded neighbor (acme-v10) must survive untouched
	if !strings.Contains(out, "acme-v10") {
		t.Errorf("guarded acme-v10 was rewritten: %q", out)
	}
}
