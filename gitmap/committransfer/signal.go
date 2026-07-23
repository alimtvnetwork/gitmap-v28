package committransfer

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// installInterruptGuard arms a SIGINT/SIGTERM handler that restores the
// source working dir to sourceHead before the process exits. It returns
// a stop func that the caller MUST defer — stop() removes the handler
// and is safe to call after a normal completion.
//
// Why this exists: Replay() already restores the source ref via a
// regular `defer` on the success/failure path, but a deferred call does
// NOT run when the process is killed by a signal. Without this guard a
// Ctrl-C in the middle of replay leaves the source repo on a detached
// HEAD pointing at some intermediate source SHA — confusing and easy
// to mistake for a real working-tree change.
//
// Behavior on signal:
//  1. Print a clear "interrupted, restoring source HEAD" line to stderr.
//  2. Run checkoutRef(sourceDir, sourceHead) best-effort.
//  3. Re-raise the signal with the default disposition so the exit
//     status reflects the interrupt (128+signo), matching git's own
//     interrupt convention.
func installInterruptGuard(sourceDir, sourceHead, logPrefix string) func() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	var once sync.Once
	done := make(chan struct{})

	go func() {
		select {
		case sig := <-ch:
			once.Do(func() {
				restoreOnSignal(sourceDir, sourceHead, logPrefix, sig)
			})
		case <-done:
			return
		}
	}()

	return func() {
		signal.Stop(ch)
		close(done)
	}
}

// restoreOnSignal runs the best-effort checkout and re-raises the
// signal so the caller exits with the conventional 128+signo code.
func restoreOnSignal(sourceDir, sourceHead, logPrefix string, sig os.Signal) {
	fmt.Fprintf(os.Stderr,
		"\n%s interrupted (%s) — restoring source HEAD to %s\n",
		logPrefix, sig, sourceHead)

	if err := checkoutRef(sourceDir, sourceHead); err != nil {
		fmt.Fprintf(os.Stderr,
			"%s WARNING: failed to restore source HEAD %s: %v\n"+
				"  → run manually: git -C %q checkout %s\n",
			logPrefix, sourceHead, err, sourceDir, sourceHead)
	} else {
		fmt.Fprintf(os.Stderr, "%s source HEAD restored.\n", logPrefix)
	}

	// Re-raise with default disposition so exit code = 128+signo.
	signal.Reset(sig)
	if p, err := os.FindProcess(os.Getpid()); err == nil {
		_ = p.Signal(sig)
	}
}
