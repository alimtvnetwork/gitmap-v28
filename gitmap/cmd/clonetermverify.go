package cmd

// clonetermverify.go — implements --verify-cmd-faithful, the dry-run
// safety net that proves the `cmd:` line printed in --output terminal
// mode matches the argv the executor would actually hand to
// exec.Command.
//
// Why this exists: every clone command builds two parallel
// representations of the same git invocation:
//
//  1. The DISPLAYED form (CloneTermBlockInput → buildCloneCommand →
//     printed as the "command:" line in the terminal block).
//  2. The EXECUTED form (Row/Plan → executor's BuildGitArgs →
//     exec.Command argv).
//
// A regression in either side (e.g. forgetting to add a new flag to
// one of the two paths) silently lies to the user. The verifier
// tokenizes the displayed form, compares it position-by-position to
// the executed argv, and prints a structured mismatch report listing
// every divergence — without running git, so it's safe to pile onto
// every clone command's pre-flight checks.
//
// Behavior contract:
//   - Pass: silent (no output) so the flag is cheap to leave on in CI.
//   - Mismatch: print a multi-line report to stderr (machine-readable
//     "diff" format) and SET the verifier's HadMismatch flag. Caller
//     decides whether to abort or just warn — this file does not call
//     os.Exit so the executor stays in charge of process lifecycle.
//
// Stream choice: stderr (not stdout) so the report doesn't pollute
// the terminal-block stream that may be piped into another tool.

import (
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// CmdFaithfulMismatch describes one position-level divergence between
// the displayed cmd: tokens and the executor's argv. Index is 0-based
// over the joined slice (git, clone, …). Either Displayed or Executed
// may be empty when one slice is shorter than the other.
type CmdFaithfulMismatch struct {
	Index     int
	Displayed string
	Executed  string
	Reason    string // short tag: "differs", "missing-in-displayed", "missing-in-executed"
}

// CmdFaithfulReport bundles the per-row verification result. Empty
// Mismatches means the two forms are byte-identical when joined by
// single spaces — which is the contract --output terminal advertises.
type CmdFaithfulReport struct {
	Repo       string
	Displayed  string // exact `cmd:` string the user would see
	Executed   string // space-joined executor argv (incl. "git")
	Mismatches []CmdFaithfulMismatch
}

// HasMismatch is a convenience predicate so callers can branch
// without poking at the slice length directly.
func (r CmdFaithfulReport) HasMismatch() bool {
	return len(r.Mismatches) > 0
}

// VerifyCmdFaithful computes the displayed cmd: line via
// buildCloneCommand and compares it token-by-token to executorArgv
// (which the caller obtains from clonenow.BuildGitArgs /
// clonefrom.BuildGitArgs / clonepick.BuildGitArgs). The "git" prefix
// is prepended to executorArgv internally so the two forms align —
// the executors return argv WITHOUT the binary, matching exec.Command's
// convention.
//
// Pure function: no I/O, deterministic. Caller decides where (if
// anywhere) to print the resulting report — see PrintCmdFaithfulReport.
func VerifyCmdFaithful(in CloneTermBlockInput, executorArgv []string) CmdFaithfulReport {
	displayed := buildCloneCommand(in)
	fullExecuted := append([]string{constants.GitBin}, executorArgv...)
	executed := strings.Join(fullExecuted, " ")

	report := CmdFaithfulReport{
		Repo:      in.Name,
		Displayed: displayed,
		Executed:  executed,
	}
	if displayed == executed {
		return report
	}
	report.Mismatches = diffArgvTokens(strings.Split(displayed, " "), fullExecuted)

	return report
}

// diffArgvTokens walks the two slices in parallel and records each
// position where they diverge. Length differences surface as
// "missing-in-X" entries for the trailing positions so the report
// surfaces the FULL extent of the drift (not just the first byte).
//
// We compare on the already-split-by-space displayed form rather
// than re-tokenizing the executed slice because buildCloneCommand
// joins on a single space and never emits embedded spaces inside a
// token (URLs, branch names, dest paths don't contain spaces in
// practice; if that ever changes the renderer needs quoting anyway,
// at which point this comparison will fail loudly — by design).
func diffArgvTokens(displayed, executed []string) []CmdFaithfulMismatch {
	var out []CmdFaithfulMismatch
	maxLen := len(displayed)
	if len(executed) > maxLen {
		maxLen = len(executed)
	}
	for i := 0; i < maxLen; i++ {
		d, e := tokenAt(displayed, i), tokenAt(executed, i)
		if d == e {
			continue
		}
		out = append(out, CmdFaithfulMismatch{
			Index: i, Displayed: d, Executed: e, Reason: classifyDiff(d, e),
		})
	}

	return out
}

// tokenAt safely indexes into a token slice, returning "" past the end.
func tokenAt(s []string, i int) string {
	if i >= len(s) {
		return ""
	}

	return s[i]
}

// classifyDiff tags the divergence so the printed report tells the
// reader at a glance whether a token was added, removed, or changed.
func classifyDiff(displayed, executed string) string {
	if len(displayed) == 0 {
		return "missing-in-displayed"
	}
	if len(executed) == 0 {
		return "missing-in-executed"
	}

	return "differs"
}

// PrintCmdFaithfulReport writes the report to w. No-op when the
// report has no mismatches so callers can invoke it unconditionally.
// Returns the first write error so a closed stderr surfaces (zero-
// swallow policy).
//
// Format (v4.13.0+): every line is prefixed with the severity tag
// `[FAIL]` so a CI log scraped one line at a time is unambiguous, and
// the banner ends with a hint about which flag promotes the report
// into a hard exit. The intent is to stop users asking "is this
// actually broken?" — the banner answers it inline.
func PrintCmdFaithfulReport(w io.Writer, r CmdFaithfulReport) error {
	if !r.HasMismatch() {
		return nil
	}
	tag := constants.CmdFaithfulReportSeverityTag
	header := fmt.Sprintf(constants.CmdFaithfulReportHeaderFmt,
		tag, r.Repo, len(r.Mismatches))
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s   displayed: %s\n", tag, r.Displayed); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s   executed:  %s\n", tag, r.Executed); err != nil {
		return err
	}
	for _, m := range r.Mismatches {
		if _, err := fmt.Fprintf(w,
			"  %s   [#%d %s] displayed=%q executed=%q\n",
			tag, m.Index, m.Reason, m.Displayed, m.Executed); err != nil {
			return err
		}
	}

	return nil
}
