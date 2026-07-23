package release

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// withRetry executes fn up to maxAttempts times with exponential backoff.
func withRetry(label string, maxAttempts int, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			printRetrySuccess(label, attempt)

			return nil
		}

		if verbose.IsEnabled() {
			verbose.Get().Log("retry: %s attempt %d/%d failed: %v", label, attempt, maxAttempts, lastErr)
		}

		if attempt < maxAttempts {
			waitAndLog(label, attempt, maxAttempts)
		}
	}

	return fmt.Errorf("%s after %d attempts: %w", label, maxAttempts, lastErr)
}

// waitAndLog prints the retry message and sleeps with exponential backoff.
func waitAndLog(label string, attempt, maxAttempts int) {
	delay := computeDelay(attempt)
	fmt.Printf(constants.MsgRetryAttempt, attempt, maxAttempts, label, delay)

	if verbose.IsEnabled() {
		verbose.Get().Log("retry: %s sleeping %v before attempt %d", label, delay, attempt+1)
	}

	time.Sleep(delay)
}

// computeDelay returns the backoff duration for the given attempt.
func computeDelay(attempt int) time.Duration {
	base := time.Duration(constants.RetryBaseDelayMs) * time.Millisecond
	factor := time.Duration(1)

	for i := 1; i < attempt; i++ {
		factor *= time.Duration(constants.RetryBackoffFactor)
	}

	return base * factor
}

// printRetrySuccess logs success with attempt number when retries occurred.
func printRetrySuccess(label string, attempt int) {
	if attempt > 1 {
		fmt.Printf(constants.MsgRetrySuccess, label, attempt)
	}
}
