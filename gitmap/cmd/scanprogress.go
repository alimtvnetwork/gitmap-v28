package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/scanner"
)

// scanProgressRenderer renders a single-line, CR-overwritten live status
// for a running gitmap scan. It is wired in via scanner.ScanOptions and
// invoked from the scanner's emitter goroutine — see
// gitmap/scanner/progress.go for the cadence and Final-snapshot contract.
//
// Design choices:
//
//   - All output goes to stderr. Stdout is reserved for machine-readable
//     artifacts (CSV/JSON paths) and the banner the rest of executeScan
//     prints; mixing a CR-overwritten line into stdout would corrupt
//     pipes and tee'd captures.
//   - When stderr is NOT a terminal (CI logs, redirected output, file
//     piping) we suppress the live frames entirely and emit only the
//     final summary line — no \r tricks, just one clean line. This
//     keeps logs grep-able without sacrificing the human UX.
//   - When --quiet is set we suppress everything. The --quiet flag is
//     the project-wide signal "I want minimal output," and this
//     renderer respects it without needing its own flag.
type scanProgressRenderer struct {
	enabled bool // emit live frames + final summary
	mu      sync.Mutex
	last    scanner.ScanProgress
	dirty   bool // a non-final frame is currently on the line
}

// newScanProgressRenderer constructs a renderer honoring --quiet and the
// stderr-isatty check. When disabled, Callback returns nil so the
// scanner never spawns its emitter goroutine at all.
func newScanProgressRenderer(quiet bool) *scanProgressRenderer {
	return &scanProgressRenderer{
		enabled: !quiet && stderrIsTerminal(),
	}
}

// Callback returns the function passed into scanner.ScanOptions. nil is
// returned when the renderer is disabled — the scanner treats that as
// "no progress hook" and skips the goroutine entirely.
func (r *scanProgressRenderer) Callback() func(scanner.ScanProgress) {
	if !r.enabled {
		return nil
	}

	return r.handle
}

// handle is invoked from the scanner's single emitter goroutine. We
// still take a lock so the eventual Done() call (which may run on the
// caller's goroutine after the scanner returns) cannot interleave with
// the final snapshot.
func (r *scanProgressRenderer) handle(p scanner.ScanProgress) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.last = p
	if p.Final {
		r.renderFinalLocked()

		return
	}

	r.renderFrameLocked(p)
}

// renderFrameLocked writes one live frame. Called only with r.mu held.
func (r *scanProgressRenderer) renderFrameLocked(p scanner.ScanProgress) {
	fmt.Fprintf(os.Stderr,
		constants.ScanProgressLineFmt,
		constants.ColorCyan, constants.ScanProgressPrefix, constants.ColorReset,
		constants.ColorWhite, p.DirsWalked, constants.ColorReset,
		constants.ColorGreen, p.ReposFound, constants.ColorReset,
	)
	r.dirty = true
}

// renderFinalLocked clears the in-flight frame (if any) and writes the
// terminating summary line. Called only with r.mu held.
func (r *scanProgressRenderer) renderFinalLocked() {
	if r.dirty {
		fmt.Fprint(os.Stderr, constants.ScanProgressClearLine)
		r.dirty = false
	}
	fmt.Fprintf(os.Stderr,
		constants.ScanProgressDoneFmt,
		constants.ColorGreen, r.last.DirsWalked, r.last.ReposFound, constants.ColorReset,
	)
}

// stderrIsTerminal reports whether stderr is connected to a real TTY.
// We avoid pulling in golang.org/x/term here — the os.Stat ModeCharDevice
// check is dependency-free and matches what isatty does on POSIX and
// Windows for the common cases (regular consoles, redirected files,
// pipes). Edge cases like Cygwin pseudo-terminals fall back to "not a
// terminal", which is the safe default for this renderer.
func stderrIsTerminal() bool {
	info, err := os.Stderr.Stat()
	if err != nil {
		return false
	}

	return (info.Mode() & os.ModeCharDevice) != 0
}
