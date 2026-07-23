package probe

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// recordingSink captures every (record, result) pair under its own
// mutex so tests can assert on the persisted set without racing the
// workers.
type recordingSink struct {
	mu   sync.Mutex
	rows []sinkRow
}

type sinkRow struct {
	record model.ScanRecord
	result Result
}

func (s *recordingSink) sink(rec model.ScanRecord, res Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = append(s.rows, sinkRow{record: rec, result: res})
}

// TestBackgroundRunner_NilWhenWorkersZero ensures the constructor
// returns nil for non-positive worker counts so callers can use a
// single nil-check for the "disabled" path.
func TestBackgroundRunner_NilWhenWorkersZero(t *testing.T) {
	if r := NewBackgroundRunner(0, 5, nil, nil); r != nil {
		t.Fatalf("expected nil runner for workers=0, got %p", r)
	}
	if r := NewBackgroundRunner(-1, 5, nil, nil); r != nil {
		t.Fatalf("expected nil runner for workers=-1, got %p", r)
	}
}

// TestBackgroundRunner_NilSafe verifies every public method tolerates
// a nil receiver. This lets the scan command write
//
//	runner := probe.NewBackgroundRunner(...)  // may be nil
//	runner.Start(rec); runner.Wait()
//
// without an explicit nil-check at every site.
func TestBackgroundRunner_NilSafe(t *testing.T) {
	var r *BackgroundRunner
	r.Start(model.ScanRecord{})
	if got := r.Stats(); (got != Stats{}) {
		t.Fatalf("nil Stats should be zero, got %+v", got)
	}
	if got := r.Remaining(); got != 0 {
		t.Fatalf("nil Remaining should be 0, got %d", got)
	}
	if got := r.Wait(); (got != Stats{}) {
		t.Fatalf("nil Wait should return zero stats, got %+v", got)
	}
}

// TestBackgroundRunner_DrainsAllJobs starts a small batch with a
// stub url-picker that synthesizes Results without touching the
// network, then asserts every record reached the sink.
func TestBackgroundRunner_DrainsAllJobs(t *testing.T) {
	const total = 12
	sink := &recordingSink{}

	// Force the empty-url branch so RunOne is never called (no git
	// process spawned in the unit test).
	r := NewBackgroundRunner(3, total,
		func(model.ScanRecord) string { return "" },
		sink.sink)

	for i := 0; i < total; i++ {
		r.Start(model.ScanRecord{ID: int64(i + 1)})
	}
	stats := r.Wait()

	sink.mu.Lock()
	defer sink.mu.Unlock()
	if len(sink.rows) != total {
		t.Fatalf("sink got %d rows, want %d", len(sink.rows), total)
	}
	if stats.Queued != total || stats.Failed != total {
		t.Fatalf("expected all %d to fail (empty url), got %+v", total, stats)
	}
}

// TestBackgroundRunner_HonorsWorkerCap installs a counter that
// tracks max concurrent invocations of the url-picker. With
// workers=2 the high-water mark must never exceed 2 even when
// many jobs are queued.
func TestBackgroundRunner_HonorsWorkerCap(t *testing.T) {
	var inflight, peak int64
	sink := &recordingSink{}
	pick := func(model.ScanRecord) string {
		now := atomic.AddInt64(&inflight, 1)
		for {
			old := atomic.LoadInt64(&peak)
			if now <= old || atomic.CompareAndSwapInt64(&peak, old, now) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt64(&inflight, -1)

		return "" // empty → RunOne not invoked
	}

	r := NewBackgroundRunner(2, 20, pick, sink.sink)
	for i := 0; i < 20; i++ {
		r.Start(model.ScanRecord{ID: int64(i + 1)})
	}
	r.Wait()

	if peak > 2 {
		t.Fatalf("worker cap violated: peak inflight=%d, want <=2", peak)
	}
}

// TestBackgroundRunner_WaitIdempotent ensures a second Wait is a
// safe no-op (the scan command may Wait both in a defer and in an
// explicit drain line; doing so must not panic with "close of
// closed channel").
func TestBackgroundRunner_WaitIdempotent(t *testing.T) {
	r := NewBackgroundRunner(1, 1,
		func(model.ScanRecord) string { return "" },
		func(model.ScanRecord, Result) {})
	r.Start(model.ScanRecord{ID: 1})
	first := r.Wait()
	second := r.Wait()
	if first != second {
		t.Fatalf("Wait not idempotent: first=%+v second=%+v", first, second)
	}
}

// TestBackgroundRunner_StartAfterWaitSilent verifies that calling
// Start after Wait does not panic — it silently drops the job. This
// is the safety net for any caller that mis-orders the lifecycle.
func TestBackgroundRunner_StartAfterWaitSilent(t *testing.T) {
	r := NewBackgroundRunner(1, 1,
		func(model.ScanRecord) string { return "" },
		func(model.ScanRecord, Result) {})
	r.Wait()

	defer func() {
		if v := recover(); v != nil {
			t.Fatalf("Start after Wait panicked: %v", v)
		}
	}()
	r.Start(model.ScanRecord{ID: 99}) // must NOT panic
}
