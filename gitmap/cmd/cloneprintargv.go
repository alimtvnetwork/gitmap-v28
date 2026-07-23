package cmd

// cloneprintargv.go — implements --print-clone-argv, the audit-only
// debug flag that dumps the exact argv tokens the executor will hand
// to exec.Command. Sits next to the --verify-cmd-faithful machinery
// because they share the same plumbing model (request-scoped toggle
// flipped once at the dispatcher, read by every per-row print
// helper).
//
// Why a separate flag from --verify-cmd-faithful:
//
//   - --verify-cmd-faithful is a SAFETY check (silent on match,
//     report on drift). Its job is to catch regressions.
//
//   - --print-clone-argv is an AUDIT tool (always prints when set).
//     Its job is to let users see the literal argv without doing
//     mental shell-tokenization on the cmd: line — handy for
//     copy-pasting into a `strace`/`ps` filter or for diffing two
//     runs of the same command.
//
// Both flags can be combined; they read independent state.
//
// Output format (deliberately machine-friendly):
//
//	  argv[0]=git
//	  argv[1]=clone
//	  argv[2]=-b
//	  argv[3]=main
//	  argv[4]=https://x/r.git
//	  argv[5]=r
//
// Two-space leading indent matches the terminal block's column so
// the audit dump visually sits "under" the block it describes.
// Tokens are NOT quoted — git argv tokens never contain newlines or
// the literal string `argv[`, so a downstream parser can split on
// `=` after the closing `]` without ambiguity.
//
// Stream choice: stderr (not stdout). The terminal block stream
// (stdout) MUST stay clean for tools that pipe it into another
// parser; debug dumps belong on stderr alongside git progress.

import (
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// cmdPrintArgv is the request-scoped flag flipped by each
// dispatcher when --print-clone-argv is parsed. See
// clonetermverifystate.go for the same-pattern rationale (CLI is
// single-threaded at dispatch, so no synchronization is needed).
var cmdPrintArgv bool

// setCmdPrintArgv enables (or disables) the argv dump for the
// remainder of the current process. Last write wins.
func setCmdPrintArgv(on bool) { cmdPrintArgv = on }

// cmdPrintArgvEnabled is the predicate the per-row helpers consult.
// Predicate (vs. exposing the var) so a future move to atomic.Bool
// or a context-bound state is a one-line refactor.
func cmdPrintArgvEnabled() bool { return cmdPrintArgv }

// printCloneArgv writes the labeled argv dump to w. The "git"
// prefix is prepended so argv[0] is the real binary (matching what
// exec.Command would put at slot 0 if the caller used "git" as the
// command name). No-op when executorArgv is empty.
//
// Returns the first write error so a closed stderr surfaces (zero-
// swallow policy).
func printCloneArgv(w io.Writer, executorArgv []string) error {
	if len(executorArgv) == 0 {
		return nil
	}
	full := append([]string{constants.GitBin}, executorArgv...)
	for i, tok := range full {
		if _, err := fmt.Fprintf(w, "  argv[%d]=%s\n", i, tok); err != nil {
			return err
		}
	}

	return nil
}

// runCmdPrintArgv is the single integration point used by every
// per-row print helper. No-op when the flag is off so callers can
// invoke it unconditionally on the hot path.
//
// Errors are surfaced (per zero-swallow policy) but do NOT abort
// the clone — the dump is purely informational. Same contract as
// runCmdFaithfulCheck so the two integrations behave identically.
func runCmdPrintArgv(executorArgv []string) {
	if !cmdPrintArgvEnabled() {
		return
	}
	if err := printCloneArgv(os.Stderr, executorArgv); err != nil {
		_, _ = os.Stderr.WriteString(
			"  Warning: --print-clone-argv: failed to write dump: " +
				err.Error() + "\n")
	}
}
