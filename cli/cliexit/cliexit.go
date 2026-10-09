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
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/output"
)

// flushers holds best-effort drainers run before os.Exit in Fail.
//
// Deprecated: the pipe-wrapped os.Stderr they used to drain is gone —
// program 267 replaced it with the synchronous FilterWriter carried
// in the dispatch context (output.UIErr), which cannot lose bytes
// before os.Exit. The registry is kept as a no-op hook so
// migration-window callers keep compiling; nothing registers anymore.
var (
	flushMu  sync.Mutex
	flushers []func()
)

// RegisterFlusher records a drain function to invoke before os.Exit
// in Fail. Order of registration is preserved; each flusher runs in
// the calling goroutine (Fail's). Idempotent flushers are encouraged.
//
// Deprecated: kept for the migration window; prefer writing through
// output.UIErr(), which needs no flush.
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

// Reportf writes a uniformly-formatted failure line to the
// dispatch-context stderr writer (output.UIErr) — the synchronous
// FilterWriter, never the old pipe-wrapped os.Stderr.
// Returns nothing — callers that need the exit-code transition use
// Fail (which calls Reportf then os.Exit). Splitting the two lets
// non-fatal collectors (e.g. per-row clone loops) reuse the same
// formatter without forcing process exit.
//
//	cliexit.Reportf("clone-from", "parse", manifestPath, err)
//	// → gitmap clone-from: parse on /path/to/manifest.json failed: <err>
//
// `subject` may be empty; the "on <subject>" segment is then elided.
// ReportParams encapsulates parameters for generating a CLI exit report.
type ReportParams struct {
	Command string
	Op      string
	Subject string
	Err     error
}

// Reportf formats and prints a standardized CLI failure line to os.Stderr.
func Reportf(command, op, subject string, err error) {
	params := ReportParams{
		Command: command,
		Op:      op,
		Subject: subject,
		Err:     err,
	}

	writeReport(output.UIErr(), params)
}

// Fail prints the standardized failure line and exits with the given
// code. The line goes through the synchronous dispatch-context writer
// (output.UIErr), so it is fully written before os.Exit — no flush
// step needed. runFlushers is retained as a no-op migration hook.
// Use this at every cmd entry-point error path so the (message,
// exit-code) pair stays atomic and impossible to forget to pair
// correctly.
func Fail(command, op, subject string, err error, code int) {
	Reportf(command, op, subject, err)
	runFlushers()
	exitFunc(code)
}

// Exit runs any registered flushers (migration-window no-op) and
// exits with the given code. Use at non-error os.Exit sites.
func Exit(code int) {
	runFlushers()
	exitFunc(code)
}

// writeReport is the format core. Extracted so the test suite can
// drive it through a bytes.Buffer without intercepting os.Stderr.
func writeReport(w io.Writer, params ReportParams) {
	if params.Err == nil {
		// Logic bug at the call site — surface loudly so it gets
		// caught in CI / dev rather than producing a confusing
		// "<no error>" line in production stderr.
		fmt.Fprintf(w,
			"gitmap %s: BUG: cliexit.Reportf called with nil err (op=%s subject=%s)\n",
			params.Command, params.Op, params.Subject)

		return
	}

	fmt.Fprintln(w, formatLine(params))
}

// formatLine assembles the canonical line. Kept side-effect-free so
// the unit test can assert byte-exact output.
func formatLine(params ReportParams) string {
	if params.Op == "not found" {
		return formatNotFoundLine(params)
	}

	if params.Subject == "" {
		return fmt.Sprintf("gitmap %s: %s failed: %v", params.Command, params.Op, params.Err)
	}

	return fmt.Sprintf("gitmap %s: %s on %s failed: %v", params.Command, params.Op, params.Subject, params.Err)
}

func formatNotFoundLine(params ReportParams) string {
	if params.Subject == "" {
		return fmt.Sprintf("gitmap %s: %v", params.Command, params.Err)
	}

	return fmt.Sprintf("gitmap %s: %s: %v", params.Command, params.Subject, params.Err)
}
