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

// Resolve translates the user-supplied --max-concurrency value into
// the effective worker count for the bounded pool. See package doc
// for the full contract.
func Resolve(n int) (int, bool) {
	return ResolveWithPriority(n, IsSSHSession())
}

// ResolveWorkerHands translates worker and hand parameters into effective concurrency.
// When isWWOH is true, enforces workerCount = 1, handCount = 1 (or keeps workerCount > 0 with handCount = 1).
// When workerCount <= 0, resolves via ResolveWithPriority.
// When handCount <= 0, defaults handCount = 1.
func ResolveWorkerHands(workerCount, handCount int, isWWOH bool, isSSH bool) (int, int) {
	if isWWOH {
		return resolveWWOHWorkerHands(workerCount)
	}
	workerCount = normalizeWorkerCount(workerCount, isSSH)
	handCount = normalizeHandCount(handCount)
	return workerCount, handCount
}

func resolveWWOHWorkerHands(workerCount int) (int, int) {
	if workerCount <= 0 {
		return 1, 1
	}
	return workerCount, 1
}

func normalizeWorkerCount(workerCount int, isSSH bool) int {
	if workerCount > 0 {
		return workerCount
	}
	resolved, ok := ResolveWithPriority(workerCount, isSSH)
	if !ok {
		return 1
	}
	if resolved <= 0 {
		return 1
	}
	return resolved
}

func normalizeHandCount(handCount int) int {
	if handCount <= 0 {
		return 1
	}
	return handCount
}
