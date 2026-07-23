// Package visibility — backoff.go: step-38 exponential backoff
// helper for provider CLI calls (`gh repo edit`, `gh repo view`).
//
// Pure: no I/O, no time.Sleep — the caller injects the sleep fn so
// tests run instantly. GitHub's secondary rate limit returns within
// a 60s window; default schedule (1s, 2s, 4s, 8s, 16s, 32s) stays
// under that ceiling at 63s total across 6 attempts.
//
// Wiring into the real provider CLI invocation site is item 45
// (provider mock harness) — this file ships the policy + tests now
// so the contract is locked before the integration lands.
package visibility

import (
	"errors"
	"time"
)

// ErrRateLimited is the sentinel callers wrap when the provider
// returns a 403 secondary rate-limit response. Use errors.Is in the
// retry predicate to distinguish from non-retryable failures
// (auth, 404, schema errors).
var ErrRateLimited = errors.New("provider rate-limited")

// backoffSchedule returns the standard 6-attempt exponential delay
// curve (1s, 2s, 4s, 8s, 16s, 32s = 63s total). Exported as a slice
// so tests can substitute a zero-delay schedule.
func backoffSchedule() []time.Duration {
	return []time.Duration{
		1 * time.Second, 2 * time.Second, 4 * time.Second,
		8 * time.Second, 16 * time.Second, 32 * time.Second,
	}
}

// RetryRateLimited runs `op` until it returns nil, a non-rate-limit
// error, or the schedule is exhausted. `sleep` is injected so tests
// can pass a no-op; production passes time.Sleep. Returns the final
// error (the last op call's error) plus the attempt count actually
// made (1-indexed; 1 = succeeded first try).
func RetryRateLimited(op func() error, schedule []time.Duration, sleep func(time.Duration)) (int, error) {
	var lastErr error
	for attempt := 0; attempt <= len(schedule); attempt++ {
		err := op()
		if err == nil {
			return attempt + 1, nil
		}
		lastErr = err
		if !errors.Is(err, ErrRateLimited) {
			return attempt + 1, err
		}
		if attempt < len(schedule) {
			sleep(schedule[attempt])
		}
	}

	return len(schedule) + 1, lastErr
}
