package cmd

import (
	"fmt"
	"testing"
)

// TestSlugFromRemote covers every remote-URL shape we care about: HTTPS,
// SSH (git@host:path), bare slug, with and without trailing `.git`. The
// expected result is always the last path segment with `.git` trimmed.
func TestSlugFromRemote(t *testing.T) {
	cases := map[string]string{
		"https://github.com/alimtvnetwork/gitmap-v28.git": "gitmap-v28",
		"https://github.com/alimtvnetwork/gitmap-v28":     "gitmap-v28",
		"git@github.com:alimtvnetwork/gitmap-v28.git":     "gitmap-v28",
		"git@github.com:alimtvnetwork/gitmap-v28":         "gitmap-v28",
		"ssh://git@host.example/foo/bar/gitmap-v28.git":   "gitmap-v28",
		"gitmap-v28":     "gitmap-v28",
		"gitmap-v28.git": "gitmap-v28",
	}

	for in, want := range cases {
		got := slugFromRemote(in)
		if got != want {
			t.Errorf("slugFromRemote(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestRemoteSlugRegex documents what the version-suffix regex accepts
// and rejects. A failed match must return nil so detectVersion can
// emit the spec's clear "expected suffix -vN" error.
//
// IMPORTANT: every {base}-vN literal here pairs with a SIBLING string
// representing the same N. Both halves MUST be derived from the same
// `int` (see fmt.Sprintf below) so a fix-repo bump cannot rewrite one
// without the other. See .lovable/memory/issues/2026-05-02-fixrepo-
// paired-literal-desync.md.
func TestRemoteSlugRegex(t *testing.T) {
	cases := buildSlugRegexCases()
	for in, w := range cases {
		m := remoteSlugRe.FindStringSubmatch(in)
		if (m != nil) != w.matches {
			t.Fatalf("regex match for %q = %v, want %v", in, m != nil, w.matches)
		}
		if !w.matches {
			continue
		}
		if m[1] != w.base || m[2] != w.num {
			t.Errorf("regex %q -> base=%q num=%q, want base=%q num=%q",
				in, m[1], m[2], w.base, w.num)
		}
	}
}

// buildSlugRegexCases derives every paired (slug, expected-num) entry
// from a single int so fix-repo cannot half-rewrite the test.
func buildSlugRegexCases() map[string]struct {
	matches bool
	base    string
	num     string
} {
	type want = struct {
		matches bool
		base    string
		num     string
	}
	out := map[string]want{}
	add := func(base string, n int) {
		slug := fmt.Sprintf("%s-v%d", base, n)
		out[slug] = want{true, base, fmt.Sprintf("%d", n)}
	}
	add("gitmap", 12)
	add("my-tool", 123)
	add("some-app-prefix", 0)
	out["gitmap"] = want{false, "", ""}
	out["gitmap-v"] = want{false, "", ""}
	out["gitmap-vX"] = want{false, "", ""}
	out[fmt.Sprintf("gitmap-v%d-extra", 12)] = want{false, "", ""}

	return out
}

// TestPairsForTarget locks the dual-form contract: every target produces
// both a `-vN` and a `/vN` replacement so Go module import paths and
// repo URLs are bumped in the same pass. Both pairs MUST use the same
// `current` value — historically this test pinned `gitmap-v28` next to
// `gitmap/v9`, a width-crossing desync that drifted whenever the
// project's own version was bumped.
func TestPairsForTarget(t *testing.T) {
	const target, current = 4, 12
	got := pairsForTarget("gitmap", target, current)
	// Diagnostic log: surfaces the actual returned values when CI
	// reports a failure here. Cheap insurance against truncated
	// failure logs that hide the assertion detail.
	t.Logf("pairsForTarget(gitmap, %d, %d) = %+v", target, current, got)
	if len(got) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(got))
	}
	wantDashOld := fmt.Sprintf("gitmap-v%d", target)
	wantDashNew := fmt.Sprintf("gitmap-v%d", current)
	if got[0].old != wantDashOld || got[0].new != wantDashNew {
		t.Errorf("dash form: got %+v, want {old:%q new:%q} (target=%d current=%d)",
			got[0], wantDashOld, wantDashNew, target, current)
	}
	wantSlashOld := fmt.Sprintf("gitmap/v%d", target)
	wantSlashNew := fmt.Sprintf("gitmap/v%d", current)
	if got[1].old != wantSlashOld || got[1].new != wantSlashNew {
		t.Errorf("slash form: got %+v, want {old:%q new:%q} (target=%d current=%d)",
			got[1], wantSlashOld, wantSlashNew, target, current)
	}
}

// TestPairsForTargetWidthCrossing locks the v9 -> v10/v12 boundary
// where the captured digit goes from 1 char to 2. Regression guard
// against the historical desync where test fixtures hard-coded
// `gitmap-v28` next to a `current=9` argument.
func TestPairsForTargetWidthCrossing(t *testing.T) {
	cases := []struct{ target, current int }{
		{9, 10}, {9, 12}, {1, 100}, {99, 100},
	}
	for _, c := range cases {
		got := pairsForTarget("gitmap", c.target, c.current)
		wantDash := fmt.Sprintf("gitmap-v%d", c.current)
		wantSlash := fmt.Sprintf("gitmap/v%d", c.current)
		if got[0].new != wantDash || got[1].new != wantSlash {
			t.Errorf("target=%d current=%d: dash.new=%q slash.new=%q, want %q / %q",
				c.target, c.current, got[0].new, got[1].new, wantDash, wantSlash)
		}
	}
}
