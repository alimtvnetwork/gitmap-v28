package cmd

import "strings"

// Shared guard-aware token scanner used by the rewriter and by tests
// that need to count or locate stale `{base}-vN` tokens in a body.
//
// Why a shared helper: previously the v9->v12 regression test
// hand-rolled its own `countUnguardedHits` / `unguardedHitOffsets`
// pair. The predicate (next byte is an ASCII digit ⇒ guarded skip)
// was duplicated verbatim from writeOneTokenHit, so any change to
// the rewriter's guard had to be mirrored in two places or the test
// would silently disagree with reality. ScanUnguardedTokenHits
// makes the rewriter and the test consume the exact same predicate
// and the exact same step-after-skip semantics, so a future tweak
// to the guard automatically propagates.
//
// Contract (mirrors writeOneTokenHit + the surrounding loop in
// rewriteToken):
//   - A "hit" is any byte offset where token starts in body.
//   - A hit is "unguarded" iff the byte immediately after the token
//     is past EOF OR is not an ASCII digit. Unguarded hits are what
//     the rewriter would substitute; guarded hits are what it would
//     leave intact (e.g. `-v9` inside `-v10`).
//   - After every hit (guarded or not) the scan advances by len(token)
//     bytes — matching `text = text[idx+tlen:]` in rewriteToken so a
//     digit-adjacent skip never re-matches the same prefix.

// ScanUnguardedTokenHits returns every byte offset in body where
// token appears AND the rewriter's negative-lookahead guard would
// allow a substitution. The returned slice is in ascending order.
func ScanUnguardedTokenHits(body, token string) []int {
	if token == "" || len(token) > len(body) {
		return nil
	}
	var hits []int
	visit := func(start int) { hits = append(hits, start) }
	walkTokenHits(body, token, visit, nil)

	return hits
}

// CountUnguardedTokenHits is a convenience wrapper that returns just
// the hit count. Equivalent to len(ScanUnguardedTokenHits(...)) but
// avoids the slice allocation on hot paths.
func CountUnguardedTokenHits(body, token string) int {
	if token == "" || len(token) > len(body) {
		return 0
	}
	count := 0
	visit := func(int) { count++ }
	walkTokenHits(body, token, visit, nil)

	return count
}

// walkTokenHits is the single source-of-truth scanner. onUnguarded is
// called with the start offset of every unguarded match; onGuarded
// (may be nil) is called for every digit-adjacent skip. Step
// semantics deliberately match rewriteToken in fixrepo_rewrite.go:
// advance len(token) bytes after every hit regardless of guard
// outcome (so a digit-adjacent skip never re-matches the same prefix).
func walkTokenHits(body, token string, onUnguarded, onGuarded func(int)) {
	tlen := len(token)
	pos := 0
	for pos+tlen <= len(body) {
		rel := strings.Index(body[pos:], token)
		if rel < 0 {
			return
		}
		idx := pos + rel
		end := idx + tlen
		if end < len(body) && isASCIIDigit(body[end]) {
			if onGuarded != nil {
				onGuarded(idx)
			}
		} else {
			onUnguarded(idx)
		}
		pos = end
	}
}
