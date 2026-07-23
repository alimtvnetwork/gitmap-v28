// Package cmd — clonespinner.go: tiny goroutine-driven spinner used
// by runCloneCommandPretty to give the user visible feedback while
// `git clone` is doing its work. Writes to stderr with carriage
// returns; the stop func clears the line on exit so subsequent
// output (success banner, git porcelain) starts at column 0.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/uipref"
)

// startCloneSpinner launches a background ticker that repaints a
// single-line spinner every 100ms. Returns a stop func the caller
// MUST invoke (usually `defer stop()`); calling it twice is safe.
//
// No-ops when stdout isn't a TTY-ish stream OR when
// cloneSpinnerOff is true — in those cases git's own output already
// gives the user enough signal and a CR-spinner would pollute logs.
func startCloneSpinner(label string) func() {
	if cloneSpinnerOff || !isStderrInteractive() || uipref.IsQuiet() || uipref.IsNoColor() {
		return func() {}
	}
	frames := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	start := time.Now()
	go func() {
		defer close(doneCh)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-stopCh:
				fmt.Fprint(os.Stderr, "\r\033[K")
				return
			case <-ticker.C:
				elapsed := time.Since(start).Truncate(time.Second)
				fmt.Fprintf(os.Stderr, "\r%s%c%s %s %s(%s)%s ",
					constants.ColorCyan, frames[i%len(frames)], constants.ColorReset,
					label,
					constants.ColorDim, elapsed, constants.ColorReset)
				i++
			}
		}
	}()
	return func() {
		select {
		case <-stopCh:
			return
		default:
			close(stopCh)
			<-doneCh
		}
	}
}

// isStderrInteractive returns true when stderr looks like a real
// terminal. We probe by checking the device mode — pipes/files give
// no IsCharDevice bit. Conservative on error: returns false so we
// never spam non-TTY logs.
func isStderrInteractive() bool {
	fi, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
