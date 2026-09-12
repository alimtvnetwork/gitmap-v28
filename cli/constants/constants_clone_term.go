package constants

import "time"

// Shared constants for the `--output terminal` adapter used by every
// clone-related command (clone, clone-next, clone-now, clone-pick,
// clone-from). Centralized here so a contract test can assert that
// every clone command surfaces the same flag name + description.

const (
	// CloneTermDetectTimeout caps how long `git ls-remote --symref
	// <url> HEAD` is allowed to run when previewing a clone in
	// `--output terminal` mode. The clone itself has no timeout —
	// this only protects the per-repo PREVIEW from a hung remote.
	CloneTermDetectTimeout = 4 * time.Second

	// BranchSourceRemoteHEAD is the canonical RepoTermBlock
	// BranchSource label used when the branch was discovered via
	// `git ls-remote --symref <url> HEAD` (i.e. the remote default
	// branch). Mirrors "HEAD" / "manifest" / "default HEAD" used by
	// the other clone commands so the rendered output is grep-able.
	BranchSourceRemoteHEAD = "remote HEAD"

	// FlagCloneTermOutput is the shared --output flag name used by
	// every clone-related command. clone-next and clone-from already
	// use the same string ("output") but kept as a constant so a
	// future rename only happens in one place.
	FlagCloneTermOutput = "output"

	// FlagDescCloneTermOutput is the shared description for clone,
	// clone-now, and clone-pick. clone-next and clone-from keep
	// their existing descriptions verbatim to avoid breaking help
	// snapshots; new commands use this one. Wording calls out the
	// stream split explicitly so users know where to redirect:
	// blocks go to stdout (machine-pipeable), git progress + the
	// human summary go to stderr.
	FlagDescCloneTermOutput = "Per-repo summary format: '' (legacy) or " +
		"'terminal' (standardized branch/from/to/command block on " +
		"stdout, streamed before each clone; git progress stays on stderr)"

	// FlagCloneVerifyCmdFaithful is the shared boolean flag name used
	// by every clone-related command to enable the dry-run verifier
	// that compares the rendered `cmd:` line against the executor's
	// real argv. Single constant so a future rename happens in one
	// place and `gitmap <cmd> --help` stays consistent across surfaces.
	FlagCloneVerifyCmdFaithful = "verify-cmd-faithful"

	// FlagDescCloneVerifyCmdFaithful explains the flag in --help.
	// Wording calls out that it's a SAFETY check (no execution
	// difference, no output unless something is wrong) so users feel
	// safe leaving it on in CI.
	FlagDescCloneVerifyCmdFaithful = "Verify the rendered cmd: line " +
		"matches the executor's real git argv. Prints a structured " +
		"mismatch report to stderr on divergence; silent on match. " +
		"Pure check — does not change clone behavior."

	// FlagCloneVerifyCmdFaithfulExitOnMismatch is the shared boolean
	// flag name used by every clone-related command to upgrade the
	// informational --verify-cmd-faithful checker into a hard failure:
	// when ANY mismatch is detected during the run, the process exits
	// with CloneVerifyCmdFaithfulExitCode after the work completes.
	// Single constant so a future rename happens in one place and
	// `gitmap <cmd> --help` stays consistent across surfaces.
	//
	// Implies --verify-cmd-faithful (the verifier must run to detect
	// drift). The flag is independent so users can opt into hard-fail
	// without re-typing both flags in CI.
	FlagCloneVerifyCmdFaithfulExitOnMismatch = "verify-cmd-faithful-exit-on-mismatch"

	// FlagDescCloneVerifyCmdFaithfulExitOnMismatch explains the flag in
	// --help. Wording is explicit about the exit code AND the timing
	// (after the run completes, not at the first mismatch) so CI logs
	// show the FULL list of divergences before the process dies.
	FlagDescCloneVerifyCmdFaithfulExitOnMismatch = "Exit non-zero (code " +
		"3) at the end of the run when --verify-cmd-faithful detects " +
		"any mismatch. Implies --verify-cmd-faithful. Mismatch reports " +
		"still print to stderr; the executor still finishes the work " +
		"so the full divergence list is logged before the non-zero exit."

	// CloneVerifyCmdFaithfulExitCode is the process exit code used
	// when --verify-cmd-faithful-exit-on-mismatch trips. Distinct from
	// 1 (runtime failure) and 2 (bad CLI usage) so CI scripts can
	// branch on "the code is wrong" vs "git failed" vs "I invoked it
	// wrong" without parsing stderr.
	CloneVerifyCmdFaithfulExitCode = 3

	// FlagClonePrintArgv is the shared boolean flag name used by every
	// clone-related command to dump the exact argv tokens that would
	// be passed to exec.Command. Companion to --verify-cmd-faithful:
	// the verifier proves displayed==executed, this flag SHOWS the
	// executed form so users can paste it elsewhere or grep one token
	// at a time without re-quoting the cmd: string.
	FlagClonePrintArgv = "print-clone-argv"

	// FlagDescClonePrintArgv explains the flag in --help. Wording
	// calls out the format (one token per line, prefixed `argv[i]=`)
	// and the stream (stderr, not stdout) so users redirecting the
	// terminal block to a file aren't surprised by extra noise.
	FlagDescClonePrintArgv = "Print the exact git-clone argv tokens " +
		"to stderr (one per line, `argv[i]=<token>`) right after each " +
		"terminal block. Audit-only — does not change clone behavior."

	// MsgCloneVerifyCmdFaithfulExit is the one-line stderr summary
	// printed immediately before the process exits with
	// CloneVerifyCmdFaithfulExitCode. Phrased as a single sentence so
	// CI UIs that surface the LAST stderr line as the failure reason
	// produce a self-explanatory headline.
	MsgCloneVerifyCmdFaithfulExit = "verify-cmd-faithful: FAIL - " +
		"one or more cmd: lines did not match the executor's argv " +
		"(see per-row reports above); exiting with code 3 because " +
		"--verify-cmd-faithful-exit-on-mismatch was set"

	// CmdFaithfulReportSeverityTag is the severity prefix used on every
	// line of a mismatch report. Picked to be unambiguous in CI logs:
	// "[FAIL]" reads as a real failure even when the surrounding
	// context is missing (e.g. log truncation, single-line scrapers).
	// Paired with CmdFaithfulReportTestPrefix so test-emitted reports
	// can be visually separated from production-emitted ones.
	CmdFaithfulReportSeverityTag = "[FAIL]" // gitmap:cmd skip

	// CmdFaithfulReportHeaderFmt is the multi-arg format string for
	// the report banner. Order: severity tag, repo name, divergence
	// count, exit-flag hint. The trailing hint tells the reader which
	// CLI flag would turn the report into a hard failure — without it,
	// users routinely ask "did this fail the build?" when the answer
	// depends on whether --verify-cmd-faithful-exit-on-mismatch was
	// passed. Verb tense is past ("did not match") to match the
	// printed report's role as a post-hoc diagnostic.
	CmdFaithfulReportHeaderFmt = "%s verify-cmd-faithful: %s — " +
		"%d divergence(s); displayed cmd: did not match executor argv " +
		"(non-fatal unless --verify-cmd-faithful-exit-on-mismatch is set)\n"

	// CmdFaithfulReportTestPrefix wraps a mismatch report when the
	// caller is a Go test that INTENTIONALLY drives a divergent input
	// to exercise the printer. Without this banner, a successful
	// `go test ./...` run prints lines that look identical to a real
	// regression — burying genuine failures and training engineers to
	// ignore "[FAIL]" in test logs. The matching close-tag below makes
	// the bracketed region trivially grep-able and easy for humans
	// scanning a wall of test output.
	CmdFaithfulReportTestPrefix = "--- expected mismatch (test fixture; " +
		"not a real failure) ---\n"

	// CmdFaithfulReportTestSuffix closes the test-fixture banner. Kept
	// as a separate constant (vs. embedding both in one block) so
	// callers can wrap multi-line content without string surgery.
	CmdFaithfulReportTestSuffix = "--- end expected mismatch ---\n" // gitmap:cmd skip
)
