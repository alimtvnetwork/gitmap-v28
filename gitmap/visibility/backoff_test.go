package visibility

import (
	"errors"
	"testing"
	"time"
)

func noSleep(time.Duration) {}

func TestRetryRateLimitedSucceedsFirstTry(t *testing.T) {
	calls := 0
	op := func() error { calls++; return nil }
	n, err := RetryRateLimited(op, backoffSchedule(), noSleep)
	if err != nil || n != 1 || calls != 1 {
		t.Fatalf("first-try: n=%d calls=%d err=%v", n, calls, err)
	}
}

func TestRetryRateLimitedRecoversMidway(t *testing.T) {
	calls := 0
	op := func() error {
		calls++
		if calls < 3 {
			return ErrRateLimited
		}

		return nil
	}
	n, err := RetryRateLimited(op, backoffSchedule(), noSleep)
	if err != nil || n != 3 || calls != 3 {
		t.Fatalf("recover: n=%d calls=%d err=%v", n, calls, err)
	}
}

func TestRetryRateLimitedNonRetryableExitsImmediately(t *testing.T) {
	boom := errors.New("404 not found")
	calls := 0
	op := func() error { calls++; return boom }
	n, err := RetryRateLimited(op, backoffSchedule(), noSleep)
	if !errors.Is(err, boom) || calls != 1 || n != 1 {
		t.Fatalf("non-retryable must NOT retry: n=%d calls=%d err=%v", n, calls, err)
	}
}

func TestRetryRateLimitedExhaustsSchedule(t *testing.T) {
	calls := 0
	op := func() error { calls++; return ErrRateLimited }
	sched := []time.Duration{0, 0, 0}
	n, err := RetryRateLimited(op, sched, noSleep)
	if !errors.Is(err, ErrRateLimited) || calls != 4 || n != 4 {
		t.Fatalf("exhaust: n=%d calls=%d err=%v (want 4/4/ErrRateLimited)", n, calls, err)
	}
}

func TestBackoffScheduleSumUnderRateLimitWindow(t *testing.T) {
	var total time.Duration
	for _, d := range backoffSchedule() {
		total += d
	}
	if total >= 64*time.Second {
		t.Fatalf("schedule %v sums to %v — must stay under GitHub's 60s window ceiling (63s budget)", backoffSchedule(), total)
	}
}
