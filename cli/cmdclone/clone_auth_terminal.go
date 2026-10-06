// Package cmdclone — clone_auth_terminal.go detects interactive terminals so
// authentication prompts fail fast instead of hanging forever when gitmap
// runs without a TTY (AI agents, CI pipelines, background jobs).
package cmdclone

import (
	"errors"
	"os"

	"golang.org/x/term"
)

// ErrAuthNonInteractive is returned when repository authentication is needed
// but stdin is not a terminal, so no interactive prompt can be shown.
// The message tells the caller exactly how to supply a token instead.
var ErrAuthNonInteractive = errors.New(
	"authentication required but stdin is not a terminal; " +
		"supply a token non-interactively via `gitmap clone <url> --token <PAT>`, " +
		"the GITHUB_TOKEN or GH_TOKEN environment variable, or `gh auth login`",
)

// stdinIsTerminal reports whether stdin is attached to an interactive terminal.
func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
