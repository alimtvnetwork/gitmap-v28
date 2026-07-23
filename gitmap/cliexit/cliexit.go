// Package cliexit centralizes user-facing CLI failure formatting and
// the os.Exit transition. Every error printed to the user from a
// `gitmap` subcommand should flow through this package so the
// vocabulary, ordering, and operational tags stay uniform across
// 80+ command files.
//
// Format contract (locked):
//
//	gitmap <command>: <op> on <subject> failed: <err>
//
// Where:
//
//   - <command> is the canonical CLI ID (e.g. "scan", "clone-from").
//     This is the string the user typed; it lets a wrapper script
//     grep stderr for `gitmap clone-from:` and route accordingly.
//   - <op>      is a short verb-led operation tag ("parse",
//     "read", "checkout", "persist", …). Mirrors the existing
//     "(operation: ...)" suffix used by ErrConfigLoad / ErrScanFailed
//     so log-grep tooling stays compatible.
//   - <subject> is the most actionable noun the caller has — usually
//     a repo path, manifest file, or URL. Empty subject is allowed
//     (it's elided from the message) but discouraged: the whole
//     point of this helper is per-call attribution.
//   - <err>     is the underlying error's Error() text. Never elided.
//
// Why a helper instead of more constants?
//
//   - Constants encode static prefixes; this helper encodes a
//     *shape*. Every call site supplies the four ingredients and
//     gets a consistent line back, so we don't have to mint a new
//     Err* constant for every (command × op × subject) triple.
//   - It collapses the very common bare `fmt.Fprintln(os.Stderr, err)`
//     anti-pattern (which leaks the underlying error with NO context
//     about which command produced it or what it was doing) into a
//     single typed call: cliexit.Reportf(cmd, op, subject, err).
//   - A future structured-logging migration only has to change this
//     one file.
package cliexit

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// flushers holds best-effort drainers run before os.Exit in Fail.
// Registered at startup by cmd.Run (one per pipe-installing package
// such as glyphs and theme). Without these, bytes written to a
// pipe-wrapped os.Stderr just before os.Exit can be lost on Windows
// because the forwarding goroutine never gets scheduled to drain
// the pipe — the documented root cause of the cliexit subprocess-
// test flakiness on Windows CI.
var (
	flushMu  sync.Mutex
	flushers []func()
)

// RegisterFlusher records a drain function to invoke before os.Exit
// in Fail. Order of registration is preserved; each flusher runs in
// the calling goroutine (Fail's). Idempotent flushers are encouraged.
func RegisterFlusher(f func()) {
	if f == nil {
		return
	}
	flushMu.Lock()
	flushers = append(flushers, f)
	flushMu.Unlock()
}

// runFlushers invokes every registered flusher in order. Recovers
// per-flusher so a panicking drainer can't swallow the exit code —
// we still want Fail to reach os.Exit with the documented code.
func runFlushers() {
	flushMu.Lock()
	fs := append([]func(){}, flushers...)
	flushMu.Unlock()
	for _, f := range fs {
		safeFlush(f)
	}
}

// safeFlush runs f under a recover() so a buggy drainer never
// prevents the exit-code transition the caller asked for.
func safeFlush(f func()) {
	defer func() { _ = recover() }()
	f()
}

// Reportf writes a uniformly-formatted failure line to os.Stderr.
// Returns nothing — callers that need the exit-code transition use
// Fail (which calls Reportf then os.Exit). Splitting the two lets
// non-fatal collectors (e.g. per-row clone loops) reuse the same
// formatter without forcing process exit.
//
//	cliexit.Reportf("clone-from", "parse", manifestPath, err)
//	// → gitmap clone-from: parse on /path/to/manifest.json failed: <err>
//
// `subject` may be empty; the "on <subject>" segment is then elided.
// `err` must NOT be nil — passing nil indicates a logic bug at the
// call site (you wouldn't be in an error path) and we surface it
// loudly instead of silently printing a half-formed message.
func Reportf(command, op, subject string, err error) {
	writeReport(os.Stderr, command, op, subject, err)
}

// Fail prints the standardized failure line, flushes any registered
// output pipes, and exits with the given code. Use this at every cmd
// entry-point error path so the (message, exit-code) pair stays
// atomic and impossible to forget to pair correctly.
func Fail(command, op, subject string, err error, code int) {
	Reportf(command, op, subject, err)
	runFlushers()
	os.Exit(code)
}

// Exit flushes any registered output pipes and exits with the given
// code. Use at non-error os.Exit sites that still need pipe-wrapped
// stdout/stderr drained before process teardown — without this, the
// final lines printed via fmt.Print can be lost on Windows when the
// glyphs/theme forwarding goroutines never get scheduled.
func Exit(code int) {
	runFlushers()
	os.Exit(code)
}

// writeReport is the format core. Extracted so the test suite can
// drive it through a bytes.Buffer without intercepting os.Stderr.
func writeReport(w io.Writer, command, op, subject string, err error) {
	if err == nil {
		// Logic bug at the call site — surface loudly so it gets
		// caught in CI / dev rather than producing a confusing
		// "<no error>" line in production stderr.
		fmt.Fprintf(w,
			"gitmap %s: BUG: cliexit.Reportf called with nil err (op=%s subject=%s)\n",
			command, op, subject)

		return
	}
	fmt.Fprintln(w, formatLine(command, op, subject, err))
}

// formatLine assembles the canonical line. Kept side-effect-free so
// the unit test can assert byte-exact output.
func formatLine(command, op, subject string, err error) string {
	if subject == "" {
		return fmt.Sprintf("gitmap %s: %s failed: %v", command, op, err)
	}

	return fmt.Sprintf("gitmap %s: %s on %s failed: %v", command, op, subject, err)
}
