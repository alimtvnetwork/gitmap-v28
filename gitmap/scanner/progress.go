// Live progress emitter for the parallel scanner.
//
// startProgress spawns a single goroutine that, on a fixed cadence,
// reads atomic counters from scanState and forwards them to the
// caller-supplied callback. A stop function is returned; calling it
//
//	(a) signals the goroutine to exit its loop, and
//	(b) blocks until the goroutine has emitted its final snapshot.
//
// The callback is therefore guaranteed to receive AT LEAST one
// invocation (the Final=true terminating snapshot) for every scan that
// passes a non-nil progress hook — even if the scan finishes before the
// first tick fires. This lets renderers print a summary line without
// having to track whether the periodic loop ever ran.
package scanner

import "time"

// progressTickInterval controls how often the emitter samples atomic
// counters and invokes the callback. ~10 Hz is fast enough to feel live
// in a human-facing terminal but slow enough that even a million-dir
// scan only emits a few hundred updates total — well below the cost
// threshold for any reasonable renderer.
const progressTickInterval = 100 * time.Millisecond

// startProgress wires the throttled emitter goroutine to the running
// scan. When `progress` is nil the cost is a single allocation and a
// no-op closure: there is no goroutine and no ticker overhead.
//
// The returned stop function is idempotent for callers but is invoked
// exactly once internally by walkParallel after the worker pool drains.
func startProgress(st *scanState, progress func(ScanProgress)) func() {
	if progress == nil {
		return func() {}
	}

	stop := make(chan struct{})
	done := make(chan struct{})

	go runProgressLoop(st, progress, stop, done)

	return func() {
		close(stop)
		<-done
	}
}

// runProgressLoop is the emitter goroutine body. It ticks at
// progressTickInterval, emitting non-final snapshots, then emits exactly
// one Final=true snapshot when stop is closed. Splitting this out keeps
// startProgress tiny and the cleanup contract obvious.
func runProgressLoop(st *scanState, progress func(ScanProgress), stop, done chan struct{}) {
	defer close(done)

	ticker := time.NewTicker(progressTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			progress(st.snapshot(false))
		case <-stop:
			progress(st.snapshot(true))

			return
		}
	}
}
