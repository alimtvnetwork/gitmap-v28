package cmd

// Helpers shared by the cn batch concurrency E2E tests
// (clonenextbatchconcurrent_e2e_test.go and
// clonenextbatchconcurrent_e2e_csv_test.go). Split out so each test
// file stays under the project's 200-line per-file budget.

import (
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// stubProcessor returns a deterministic batchRowResult per repo path
// with a small randomized sleep so workers genuinely interleave
// completion order. Status cycles ok/failed/skipped by trailing
// digit so test counts are exact, not statistical.
//
// `concurrentSeen` tracks the maximum simultaneous in-flight count
// so tests can assert the pool actually went parallel — a regression
// guard against accidental serialization.
func stubProcessor(concurrentSeen *int64) func(string) batchRowResult {
	var inflight int64
	return func(path string) batchRowResult {
		now := atomic.AddInt64(&inflight, 1)
		for {
			peak := atomic.LoadInt64(concurrentSeen)
			if now <= peak || atomic.CompareAndSwapInt64(concurrentSeen, peak, now) {
				break
			}
		}
		base := filepath.Base(path)
		last := base[len(base)-1] - '0'
		time.Sleep(time.Duration(last%5+1) * time.Millisecond)
		atomic.AddInt64(&inflight, -1)

		return batchRowResult{
			RepoPath:    path,
			FromVersion: "v1",
			ToVersion:   fmt.Sprintf("v%d", last+1),
			Status:      pickStubStatus(last),
		}
	}
}

// pickStubStatus distributes 50 inputs across the three buckets in
// a known ratio (digits 0-3 → ok, 4-7 → failed, 8-9 → skipped) so
// the count assertions are exact, not statistical.
func pickStubStatus(last byte) string {
	switch {
	case last <= 3:
		return constants.BatchStatusOK
	case last <= 7:
		return constants.BatchStatusFailed
	default:
		return constants.BatchStatusSkipped
	}
}

// installStubProcessor swaps processOneBatchRepoFn for the test and
// restores the original in t.Cleanup. Returns the peak-inflight
// counter so callers can assert the pool actually parallelized.
func installStubProcessor(t *testing.T) *int64 {
	t.Helper()
	original := processOneBatchRepoFn
	var peak int64
	processOneBatchRepoFn = stubProcessor(&peak)
	t.Cleanup(func() {
		processOneBatchRepoFn = original
	})
	return &peak
}

// makeRepoPaths returns n synthetic repo paths whose trailing digit
// drives the stubProcessor's status assignment.
func makeRepoPaths(n int) []string {
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = fmt.Sprintf("/tmp/repo-%d", i)
	}
	return out
}
