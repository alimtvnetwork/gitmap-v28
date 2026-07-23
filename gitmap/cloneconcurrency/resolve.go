// Package cloneconcurrency centralizes the worker-count resolver
// shared by every clone-family command (`clone`, `clone-next`,
// `clone-now` / `relclone`, `clone-from`).
//
// The resolver is its own leaf package so the various clone packages
// can depend on it without importing each other (which would create
// cycles: cloner already pulls in model, while clonenow/clonefrom
// must stay independent of cloner so their dry-run renderers can be
// unit-tested without git on the PATH).
//
// Contract — Resolve(n):
//
//   - n  < 0 → returns 0 + ok=false. Callers MUST treat this as a
//     hard CLI usage error (exit 1, print ErrCloneMaxConcurrencyInvalid).
//   - n == 0 → "auto": returns max(1, runtime.NumCPU()) + ok=true.
//     This is the documented default when --max-concurrency is omitted.
//   - n  > 0 → returns n + ok=true verbatim. Capping is the user's
//     responsibility; we do NOT silently clamp because doing so
//     would hide a typo like `--max-concurrency 1000`.
//
// runtime.NumCPU is read once per call (cheap; no caching needed —
// callers invoke this exactly once at command startup).
package cloneconcurrency

import "runtime"

// Resolve translates the user-supplied --max-concurrency value into
// the effective worker count for the bounded pool. See package doc
// for the full contract.
func Resolve(n int) (int, bool) {
	if n < 0 {
		return 0, false
	}
	if n == 0 {
		w := runtime.NumCPU()
		if w < 1 {
			// runtime.NumCPU returning <1 is theoretically impossible
			// per the Go spec but cheap to defend against — guarantees
			// the caller never sees a zero-worker pool.
			w = 1
		}

		return w, true
	}

	return n, true
}
