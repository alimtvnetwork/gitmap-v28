// Structured error reporting for cliexit. Builds on the canonical
// single-line text format (see cliexit.go) by adding a richer
// `Context` payload (path, args, mode, free-form extras) and a
// machine-readable JSON output mode.
//
// Why a second entry-point instead of changing Reportf?
//
//   - The existing `gitmap <command>: <op> on <subject> failed: <err>`
//     line is wired into wrapper scripts, CI annotations, and the
//     bare-stderr lint check. Changing its shape would be a breaking
//     observable contract.
//   - Most call sites have only (command, op, subject, err) and want
//     the terse line. Forcing them through a struct API would be
//     ceremony for no gain.
//   - The minority of sites that do have rich context (clone retry
//     loops, scan walkers with a current path + flag mode + invocation
//     args) get a typed entry-point that captures all of it without
//     each site re-inventing its own concatenation.
//
// Output modes:
//
//   - OutputHuman (default): the canonical text line, plus one
//     "  <key>=<value>" indented line per non-empty context field.
//     Stable enough for humans, ignorable enough for greppers that
//     only care about the leading line.
//   - OutputJSON: a single-line JSON object on stderr. Suitable for
//     log shippers / CI harnesses that want structured failures.
//     One object per line so streaming consumers can parse without
//     buffering.
//
// Mode is selected per call (callers thread their `--output` flag
// through). No global state, no env-var sniffing — keeps the helper
// trivially testable and avoids surprising behavior shifts between
// processes.

package cliexit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// OutputMode selects the rendering for Report / FailWith. Zero value
// (OutputHuman) keeps the existing terminal-friendly shape so callers
// that don't care about JSON pay no migration cost.
type OutputMode int

const (
	// OutputHuman renders the canonical text line + indented
	// context fields. This is the default.
	OutputHuman OutputMode = iota
	// OutputJSON renders one single-line JSON object per failure.
	OutputJSON
)

// Context carries the structured fields a rich error site can
// supply. Every field is optional — empty / zero values are elided
// from the output so call sites only fill what they actually have.
//
// Field meanings:
//
//   - Command  : canonical CLI ID (same as Reportf's `command`).
//   - Op       : verb-led operation tag (same as Reportf's `op`).
//   - Path     : repo path / manifest file / destination dir — the
//     "where" of the failure. Used as `subject` in the human line.
//   - Args     : the user-supplied positional args / flags relevant
//     to the failure. Joined with spaces in human mode, kept as a
//     slice in JSON mode.
//   - Mode     : the operating mode the command was in (e.g.
//     "execute", "dry-run", "audit"). Surfaced because the same
//     command often behaves differently per mode and "which mode
//     was active" is the first thing a user asks.
//   - Extras   : free-form key/value pairs for site-specific context
//     (e.g. {"row":"3","url":"https://…"}). Rendered alphabetically
//     so output is deterministic.
//   - Err      : the underlying error. Required; nil triggers the
//     same BUG guard as Reportf.
type Context struct {
	Command string
	Op      string
	Path    string
	Args    []string
	Mode    string
	Extras  map[string]string
	Err     error
}

// Report writes the structured failure to os.Stderr in the requested
// mode. Use FailWith when you also need the os.Exit transition.
func Report(ctx Context, mode OutputMode) {
	writeStructured(os.Stderr, ctx, mode)
}

// FailWith reports then exits with `code`. Atomic (message, exit-code)
// pairing — same rationale as the existing Fail helper.
func FailWith(ctx Context, mode OutputMode, code int) {
	Report(ctx, mode)
	os.Exit(code)
}

// writeStructured is the rendering core, extracted so tests can drive
// it through a bytes.Buffer without intercepting os.Stderr.
func writeStructured(w io.Writer, ctx Context, mode OutputMode) {
	if ctx.Err == nil {
		fmt.Fprintf(w,
			"gitmap %s: BUG: cliexit.Report called with nil err (op=%s path=%s)\n",
			ctx.Command, ctx.Op, ctx.Path)

		return
	}
	if mode == OutputJSON {
		writeJSON(w, ctx)

		return
	}
	writeHuman(w, ctx)
}

// writeHuman emits the canonical text line followed by indented
// "  key=value" lines for every populated context field.
func writeHuman(w io.Writer, ctx Context) {
	fmt.Fprintln(w, formatLine(ctx.Command, ctx.Op, ctx.Path, ctx.Err))
	for _, line := range humanContextLines(ctx) {
		fmt.Fprintln(w, "  "+line)
	}
}

// humanContextLines builds the indented "key=value" tail in a stable
// order: mode, then args, then alphabetized extras. Path is already
// in the leading line so it isn't repeated here.
func humanContextLines(ctx Context) []string {
	out := make([]string, 0, 2+len(ctx.Extras))
	if ctx.Mode != "" {
		out = append(out, "mode="+ctx.Mode)
	}
	if len(ctx.Args) > 0 {
		out = append(out, "args="+strings.Join(ctx.Args, " "))
	}

	return append(out, sortedExtraLines(ctx.Extras)...)
}

// sortedExtraLines renders Extras alphabetically by key so output is
// deterministic across runs / Go versions.
func sortedExtraLines(extras map[string]string) []string {
	if len(extras) == 0 {
		return nil
	}
	keys := make([]string, 0, len(extras))
	for k := range extras {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+"="+extras[k])
	}

	return out
}

// writeJSON emits a single-line JSON object. Trailing newline so
// streaming consumers can split on \n without a sentinel.
func writeJSON(w io.Writer, ctx Context) {
	payload := map[string]any{
		"command": ctx.Command,
		"op":      ctx.Op,
		"error":   ctx.Err.Error(),
	}
	addNonEmptyString(payload, "path", ctx.Path)
	addNonEmptyString(payload, "mode", ctx.Mode)
	if len(ctx.Args) > 0 {
		payload["args"] = ctx.Args
	}
	if len(ctx.Extras) > 0 {
		payload["extras"] = ctx.Extras
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		// Fall back to the human line — losing JSON shape is
		// strictly better than swallowing the error entirely.
		writeHuman(w, ctx)

		return
	}
	fmt.Fprintln(w, string(encoded))
}

// addNonEmptyString keeps writeJSON readable by hiding the
// "skip if empty" guard for optional string fields.
func addNonEmptyString(m map[string]any, key, value string) {
	if value == "" {
		return
	}
	m[key] = value
}
