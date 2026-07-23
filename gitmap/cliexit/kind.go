// Centralized error-kind → exit-code mapping for `gitmap`. Pairs
// with cliexit.Report* so a single call classifies the failure,
// emits the standardized context line(s) on stderr, and exits with
// the canonical code for that class.
//
// Why centralize this?
//
//   - Today exit codes are duplicated as ad-hoc constants across the
//     codebase (CloneNowExitConfirmAborted=2, scattered os.Exit(1),
//     etc.). Wrapper scripts that branch on `$?` need a single
//     authoritative table to grep against.
//   - Bundling "what kind of failure this is" with "what code do we
//     exit with" prevents the recurring mistake of printing a "user
//     canceled" message and then exiting 1, or vice versa.
//   - Future additions (timeout, conflict) become a one-line table
//     change instead of a sweep across cmd/.
//
// Code table (locked — wrapper scripts depend on these numbers):
//
//	0 KindSuccess         — reserved for symmetry; never passed to FailKind
//	1 KindExecutionFailed — operational failure (git clone died, IO error, …)
//	2 KindUserCanceled    — user said no, or non-TTY without --yes
//	2 KindInvalidInput    — bad flag, malformed manifest, missing positional
//	3 KindVerifyFailed    — read-only verify/audit found a divergence
//	4 KindPreconditionFailed — environment unmet (missing git, locked DB)
//
// Note KindUserCanceled and KindInvalidInput share code 2 by design:
// both mean "we did nothing destructive; fix your input or your
// answer and rerun". This matches the long-standing convention used
// by `git`, `curl`, and POSIX utilities (exit 2 = usage / refusal).

package cliexit

// Kind classifies a CLI failure into one of the documented buckets.
// Use FailKind(ctx, kind, mode) at the cmd entry-point error path so
// the (message, exit-code) pair stays atomic.
type Kind int

const (
	// KindSuccess is the zero-value sentinel. Never pass it to
	// FailKind — it's only here so `var k Kind` defaults to a
	// meaningful name in debuggers.
	KindSuccess Kind = iota
	// KindExecutionFailed covers any operational failure that
	// happened while doing real work: git clone died, a file
	// couldn't be written, a network call timed out. Exit 1.
	KindExecutionFailed
	// KindUserCanceled is "the user said no" — interactive
	// abort, non-TTY without --yes, Ctrl-C at a confirm prompt.
	// Exit 2. Distinct semantics from KindInvalidInput but
	// shares the code (see file header).
	KindUserCanceled
	// KindInvalidInput is bad CLI input: unknown flag, missing
	// positional, malformed manifest, unparseable URL. Exit 2.
	KindInvalidInput
	// KindVerifyFailed is a read-only check that found a
	// divergence (clone-audit mismatch, verify-cmd-faithful
	// mismatch). Exit 3 — distinguishable from operational
	// failures so CI can branch on "drift vs broken".
	KindVerifyFailed
	// KindPreconditionFailed means the environment isn't ready:
	// `git` not on PATH, DB locked by another process, required
	// secret missing. Exit 4 — separate from input errors so
	// retry logic can wait-and-retry instead of giving up.
	KindPreconditionFailed
)

// kindCodes is the single authoritative table. Wrapper scripts and
// the cliexit_*_test.go suite both grep these numbers.
var kindCodes = map[Kind]int{
	KindSuccess:            0,
	KindExecutionFailed:    1,
	KindUserCanceled:       2,
	KindInvalidInput:       2,
	KindVerifyFailed:       3,
	KindPreconditionFailed: 4,
}

// kindLabels is the short tag rendered in human/JSON output so the
// user sees *why* the process is exiting with that code, not just
// the number.
var kindLabels = map[Kind]string{
	KindSuccess:            "success",
	KindExecutionFailed:    "execution-failed",
	KindUserCanceled:       "user-canceled",
	KindInvalidInput:       "invalid-input",
	KindVerifyFailed:       "verify-failed",
	KindPreconditionFailed: "precondition-failed",
}

// KindCode returns the canonical exit code for a Kind. Unknown
// values fall back to 1 so a future enum addition that forgets to
// update the table still produces a "something went wrong" exit
// rather than a misleading 0.
func KindCode(k Kind) int {
	if code, ok := kindCodes[k]; ok {
		return code
	}

	return 1
}

// KindLabel returns the short human label for a Kind. Used by
// FailKind to tag the stderr context with the failure class so a
// reader doesn't have to memorize the code table.
func KindLabel(k Kind) string {
	if label, ok := kindLabels[k]; ok {
		return label
	}

	return "execution-failed"
}

// FailKind reports the structured failure (with the kind tag added
// to Extras as `kind=<label>`) and exits with KindCode(kind). This
// is the single front door for "I have a classified failure, write
// it out and exit." Use it at every cmd entry-point error path.
//
// The kind is always added to Extras so JSON consumers can branch
// on a stable string field instead of reverse-engineering it from
// the exit code.
func FailKind(ctx Context, kind Kind, mode OutputMode) {
	tagged := withKindExtra(ctx, kind)
	FailWith(tagged, mode, KindCode(kind))
}

// withKindExtra returns a copy of ctx with kind=<label> appended to
// Extras. Defensive copy so callers can reuse their Context value
// without surprise mutations.
func withKindExtra(ctx Context, kind Kind) Context {
	extras := make(map[string]string, len(ctx.Extras)+1)
	for k, v := range ctx.Extras {
		extras[k] = v
	}
	extras["kind"] = KindLabel(kind)
	ctx.Extras = extras

	return ctx
}
