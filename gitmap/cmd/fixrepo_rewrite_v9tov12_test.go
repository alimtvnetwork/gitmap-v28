package cmd

// Regression: end-to-end fixture bump from v9 to v12 using the
// fix-repo token-rewrite engine. Locks the width-crossing boundary
// (1-digit -> 2-digit) and cross-validates that pairsForTarget +
// remoteSlugRe still agree with the rewritten bytes.
//
// Background: historically the rewriter, the pair builder, and the
// remote-slug regex drifted independently when the project's own
// version bumped past v9. This test wires all three into a single
// assertion so any future desync fails loudly with a layered diff.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fixtureversion"
)

// fixRepoV9ToV12FixtureBody is the on-disk fixture: every realistic
// shape we have seen in third-party Go repos that depend on a
// versioned module — bare slug, dash form, slash form, and a digit-
// adjacent token (`gitmap-v28`) that MUST NOT match `gitmap-v28`.
// We use `-v10` (a real, plausible neighbor version) rather than the
// nonsensical `-v90` to keep the fixture readable while still locking
// the negative-lookahead guard against `-v9` matching inside `-v10`.
//
// The first line is a fixtureversion stamp so a stale fixture (one
// whose generation lags the test's MinGeneration) fails with an
// actionable "regenerate via ..." message instead of a confusing
// rewrite-count mismatch.
//
// IMPORTANT: this body MUST contain `gitmap-v28` tokens (the rewrite
// target) and a `gitmap-v28` guarded neighbor. A previous global
// rename collapsed every `-v9` token into `-v16` and silently broke
// the test (rewriter found 0 tokens to bump). See
// .lovable/memory/issues/2026-05-01-fixrepo-digit-capture-desync.md.
// IMPORTANT: base MUST be a synthetic name (e.g. `acme`) — using the
// project's own module suffix (`gitmap`) collapses to `gitmap-v28` /
// `gitmap-v110` under fix-repo runs against this very repo, silently
// erasing the v9 tokens. See mem://FIX-REPO DIGIT-CAPTURE GAP and
// .lovable/memory/issues/2026-05-01-fixrepo-digit-capture-desync.md.
const fixRepoV9ToV12FixtureBody = `// fixture-stamp: name=fixrepo-v9-to-v12 generation=3 min-current=12 for=v9->v12-width-cross sha=916cb2573c38
module example.com/consumer

require (
	github.com/example/acme-v9 v0.0.0
)

import gm "github.com/example/acme-v9/pkg"

// repo URL: https://github.com/example/acme-v9.git
// guarded:  acme-v10 must NOT be rewritten by target=9 (v9 is a
//           prefix of v10 — the negative-lookahead guard skips it)
`

// TestFixRepoRewriteV9ToV12Fixture is the end-to-end regression test.
// It writes a fixture, invokes the rewrite engine for target=9 with
// current=12, then asserts the rewritten bytes against pairsForTarget
// and feeds the new slug back through remoteSlugRe.
func TestFixRepoRewriteV9ToV12Fixture(t *testing.T) {
	const (
		base    = "acme"
		target  = 9
		current = 12
	)
	// Fail fast if the embedded fixture is stale relative to what
	// this test expects. Under `make fixtures-bump` the marker is
	// rewritten in this very source file automatically; otherwise
	// this t.Fatals with an actionable regenerate recipe.
	fixtureversion.MustValidateBodyWithAutobump(t, fixRepoV9ToV12FixtureBody,
		"fixrepo_rewrite_v9tov12_test.go",
		fixtureversion.Expectation{
			MinGeneration:    1,
			CurrentVersion:   current,
			RegenerateRecipe: "run `make fixtures-bump RUN=TestFixRepoRewriteV9ToV12Fixture` (or hand-edit the // fixture-stamp: marker)",
		})
	path := writeV9Fixture(t)

	count, err := rewriteFixRepoFile(path, base, current, []int{target}, false)
	if err != nil {
		t.Fatalf("rewriteFixRepoFile: %v", err)
	}
	got := readFile(t, path)

	assertDashFormBumped(t, got, base, target, current, count)
	assertGuardedNeighborPreserved(t, got, base)
	assertPairsForTargetAgrees(t, base, target, current, got)
	assertRemoteSlugRegexAgrees(t, base, current)
}

// writeV9Fixture materializes the fixture body in a temp file and
// returns its path. t.TempDir auto-cleans on test exit.
func writeV9Fixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "go.mod.fixture")
	if err := os.WriteFile(path, []byte(fixRepoV9ToV12FixtureBody), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	return path
}

// readFile is a tiny test helper that fatals on read error so
// individual assertions can stay focused on content checks.
func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}

	return string(b)
}

// assertDashFormBumped verifies every `{base}-v{target}` token was
// rewritten to `{base}-v{current}` and the replacement count matches
// the number of dash-form occurrences in the fixture.
func assertDashFormBumped(t *testing.T, got, base string, target, current, count int) {
	t.Helper()
	oldTok := fmt.Sprintf("%s-v%d", base, target)
	newTok := fmt.Sprintf("%s-v%d", base, current)
	if countUnguardedHits(got, oldTok) > 0 {
		t.Errorf("found stale unguarded %q after bump\n%s",
			oldTok, renderFixRepoFailureDiff(got, oldTok, newTok, target, current))
	}
	if !strings.Contains(got, newTok) {
		t.Errorf("missing bumped %q after rewrite\n%s",
			newTok, renderFixRepoFailureDiff(got, oldTok, newTok, target, current))
	}
	wantCount := countUnguardedHits(fixRepoV9ToV12FixtureBody, oldTok)
	if count != wantCount {
		t.Errorf("replacement count = %d, want %d (unguarded hits in fixture)\n%s",
			count, wantCount,
			renderFixRepoFailureDiff(got, oldTok, newTok, target, current))
	}
}

// countUnguardedHits is a thin wrapper around the production
// scanner so the test cannot drift from the rewriter's actual
// negative-lookahead guard. Any future tweak to the guard
// predicate (e.g. supporting non-ASCII digits) automatically
// propagates here. See fixrepo_rewrite_scan.go for the contract.
func countUnguardedHits(body, token string) int {
	return CountUnguardedTokenHits(body, token)
}

// assertGuardedNeighborPreserved locks the negative-lookahead guard:
// `gitmap-v28` must survive untouched when bumping target=9, because
// `-v9` is a prefix of `-v10` and the rewriter's negative-lookahead
// must skip digit-adjacent matches.
func assertGuardedNeighborPreserved(t *testing.T, got, base string) {
	t.Helper()
	guarded := base + "-v10"
	if !strings.Contains(got, guarded) {
		oldTok := base + "-v9"
		newTok := base + "-v12"
		t.Errorf("guarded neighbor %q was incorrectly rewritten\n%s",
			guarded, renderFixRepoFailureDiff(got, oldTok, newTok, 9, 12))
	}
}

// assertPairsForTargetAgrees feeds the same (base, target, current)
// triple through pairsForTarget and asserts the dash-form `new`
// string is exactly what landed in the rewritten file. This is the
// cross-layer guard against the v9->v12 width-crossing desync.
func assertPairsForTargetAgrees(t *testing.T, base string, target, current int, got string) {
	t.Helper()
	pairs := pairsForTarget(base, target, current)
	if len(pairs) < 1 {
		t.Fatalf("pairsForTarget returned %d pairs, want >=1", len(pairs))
	}
	if !strings.Contains(got, pairs[0].new) {
		oldTok := fmt.Sprintf("%s-v%d", base, target)
		t.Errorf("rewriter output missing pairsForTarget dash.new=%q\n%s",
			pairs[0].new,
			renderFixRepoFailureDiff(got, oldTok, pairs[0].new, target, current))
	}
}

// assertRemoteSlugRegexAgrees feeds the bumped slug back through
// remoteSlugRe and confirms the captured base/num match the values
// the rewriter just produced. Closes the loop: rewriter -> regex.
func assertRemoteSlugRegexAgrees(t *testing.T, base string, current int) {
	t.Helper()
	bumpedSlug := fmt.Sprintf("%s-v%d", base, current)
	m := remoteSlugRe.FindStringSubmatch(bumpedSlug)
	if m == nil {
		t.Fatalf("remoteSlugRe did not match bumped slug %q", bumpedSlug)
	}
	wantNum := fmt.Sprintf("%d", current)
	if m[1] != base || m[2] != wantNum {
		t.Errorf("remoteSlugRe(%q) = base=%q num=%q, want base=%q num=%q",
			bumpedSlug, m[1], m[2], base, wantNum)
	}
}

// renderFixRepoFailureDiff produces a single multi-section diagnostic
// block for any fix-repo rewrite assertion failure. CI logs only show
// what is in the t.Errorf payload, so we pre-bake the most useful
// signals: target/current versions, hit counts on both the fixture
// and the rewritten output, every unguarded `-vN` match with a line
// number + surrounding context, and the full rewritten file.
func renderFixRepoFailureDiff(got, oldTok, newTok string, target, current int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "  -- fix-repo rewrite diff (target=v%d -> current=v%d)\n",
		target, current)
	fmt.Fprintf(&b, "     old token        = %q\n", oldTok)
	fmt.Fprintf(&b, "     new token        = %q\n", newTok)
	fmt.Fprintf(&b, "     fixture unguarded %s = %d\n",
		oldTok, countUnguardedHits(fixRepoV9ToV12FixtureBody, oldTok))
	fmt.Fprintf(&b, "     output unguarded  %s = %d (want 0)\n",
		oldTok, countUnguardedHits(got, oldTok))
	fmt.Fprintf(&b, "     output occurrences of %s = %d\n",
		newTok, strings.Count(got, newTok))
	b.WriteString(renderUnguardedHitContext(got, oldTok))
	b.WriteString("  -- rewritten file --\n")
	b.WriteString(indentLines(got, "    "))

	return b.String()
}

// renderUnguardedHitContext walks every unguarded occurrence of token
// in body and emits a `line N: <line>` block per hit so the CI log
// pinpoints exactly where the rewriter missed.
func renderUnguardedHitContext(body, token string) string {
	var b strings.Builder
	hits := unguardedHitOffsets(body, token)
	if len(hits) == 0 {
		return ""
	}
	b.WriteString("  -- unguarded stale matches --\n")
	for _, off := range hits {
		line, col, text := lineAtOffset(body, off)
		fmt.Fprintf(&b, "    line %d col %d: %s\n", line, col, text)
	}

	return b.String()
}

// unguardedHitOffsets returns every byte offset where token appears
// without a digit immediately after it. Delegates to the production
// scanner so the test's "where did the rewriter miss" diagnostic and
// the rewriter's actual substitution decision use byte-identical
// predicates. See fixrepo_rewrite_scan.go.
func unguardedHitOffsets(body, token string) []int {
	return ScanUnguardedTokenHits(body, token)
}

// lineAtOffset returns the 1-based line number, 1-based column, and
// the line text containing byte offset off in body. Used to label
// each stale-match site in the failure diff.
func lineAtOffset(body string, off int) (int, int, string) {
	if off < 0 || off > len(body) {
		return 0, 0, ""
	}
	line := 1 + strings.Count(body[:off], "\n")
	lineStart := strings.LastIndexByte(body[:off], '\n') + 1
	lineEnd := lineStart + strings.IndexByte(body[lineStart:], '\n')
	if lineEnd < lineStart {
		lineEnd = len(body)
	}

	return line, off - lineStart + 1, body[lineStart:lineEnd]
}

// indentLines prefixes every line in s with prefix. Keeps the
// rewritten-file dump visually distinct from the surrounding test
// log so CI scrollback stays scannable.
func indentLines(s, prefix string) string {
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}
