// Package cmd — vscodepmsofterror.go: shared soft-fail reporter for
// every code path that touches the alefragnani.project-manager
// projects.json file.
//
// Soft-fail policy: a failed VS Code Project Manager interaction
// (missing user-data root, extension not installed, parse error on
// a hand-edited projects.json, transient write failure) MUST NEVER
// turn a successful gitmap operation into a non-zero exit code.
// Every caller in gitmap/cmd routes through this single function so
// the wording, formatting, and stderr destination stay consistent —
// and so future changes to the policy (e.g. demoting the line to a
// debug-only trace) happen in exactly one place.
//
// The reporter intentionally writes to os.Stderr (not os.Stdout) so
// scripts that pipe gitmap output through `jq` / `tee` keep clean
// streams. Callers do not need to print anything else after invoking
// it — the produced line already includes the "vscode:" namespace
// prefix used by every other helper in gitmap/vscodepm.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// reportVSCodePMSoftError prints a single namespaced line describing
// a non-fatal VS Code Project Manager failure. nil errors are a
// no-op so callers can pass through error values unconditionally.
//
// Special-cased errors get a friendlier message that tells the user
// exactly which install step is missing (VS Code itself vs. the
// project-manager extension) instead of the bare error chain.
func reportVSCodePMSoftError(err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, vscodepm.ErrUserDataMissing):
		fmt.Fprintln(os.Stderr, constants.MsgVSCodePMPathRootMissing)
	case errors.Is(err, vscodepm.ErrExtensionMissing):
		fmt.Fprintf(os.Stderr, constants.MsgVSCodePMPathExtMissing, err.Error())
	default:
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
	}
}
