package cmd

// Pre-flight safety prompt for `gitmap reclone --execute`.
//
// Goal: refuse to run `git clone` against destinations that already
// exist on disk unless the user explicitly confirms (interactive
// prompt) or pre-confirms (--yes). The per-row --on-exists policy
// (skip / update / force) still decides what actually happens to
// each existing directory; this gate is a single high-level
// "are you sure?" that stops the most common foot-gun -- accidentally
// running `force` (or any execute) over a populated tree -- before
// any side effect occurs.
//
// Resolution rules (priority order):
//
//  1. --yes              => skip the prompt, proceed.
//  2. No existing dests  => no prompt needed, proceed.
//  3. Non-TTY stdin      => refuse (exit 2). CI must pass --yes.
//  4. Interactive prompt => "y" proceeds, anything else aborts (exit 2).

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// confirmCloneNowExistingDestsOrExit gates the executor on a user
// confirmation when any planned destination already exists. Exits
// CloneNowExitConfirmAborted (2) on refusal / non-TTY without --yes
// so wrappers can distinguish "user said no" from per-row clone
// failures (exit 1).
func confirmCloneNowExistingDestsOrExit(plan clonenow.Plan, cfg cloneNowFlags) {
	if cfg.assumeYes {

		return
	}
	existing := collectExistingDests(plan, cfg.cwd)
	if len(existing) == 0 {

		return
	}
	printExistingDestsPreview(existing, plan.OnExists)
	if !isStdinInteractive() {
		fmt.Fprint(os.Stderr, constants.MsgCloneNowConfirmNonTTY)
		os.Exit(constants.CloneNowExitConfirmAborted)
	}
	if !readUserConfirmation() {
		fmt.Fprint(os.Stderr, constants.MsgCloneNowConfirmAborted)
		os.Exit(constants.CloneNowExitConfirmAborted)
	}
}

// collectExistingDests returns the relative paths whose resolved
// destination directory already exists on disk. Resolution joins
// the row's RelativePath onto cwd (empty cwd = process cwd). Stat
// errors other than "not exist" are treated as "exists" -- safer
// to over-prompt than to miss a permission-denied dir that still
// blocks `git clone`.
func collectExistingDests(plan clonenow.Plan, cwd string) []string {
	base := cwd
	if base == "" {
		base = "."
	}
	out := make([]string, 0, len(plan.Rows))
	for _, r := range plan.Rows {
		if r.RelativePath == "" {

			continue
		}
		cleaned := model.CleanRelativePath(r.RelativePath)
		if destPathExists(filepath.Join(base, cleaned)) {
			out = append(out, cleaned)
		}
	}

	return out
}

// destPathExists treats any non-IsNotExist error as "exists" so a
// permission-denied parent can't sneak past the gate.
func destPathExists(p string) bool {
	_, err := os.Stat(p)
	if err == nil {

		return true
	}

	return !os.IsNotExist(err)
}

// printExistingDestsPreview renders the header + bullet list to
// stderr, capping the visible bullets at CloneNowExistingPreviewLimit
// while always reporting the full count.
func printExistingDestsPreview(existing []string, onExists string) {
	fmt.Fprintf(os.Stderr, constants.MsgCloneNowConfirmHeader,
		len(existing), onExists)
	limit := constants.CloneNowExistingPreviewLimit
	shown := len(existing)
	if shown > limit {
		shown = limit
	}
	for i := 0; i < shown; i++ {
		fmt.Fprintf(os.Stderr, constants.MsgCloneNowConfirmBullet, existing[i])
	}
	if len(existing) > shown {
		fmt.Fprintf(os.Stderr, constants.MsgCloneNowConfirmTruncated,
			len(existing)-shown)
	}
}

// isStdinInteractive returns true when stdin is attached to a
// character device (terminal). Required so CI / piped runs don't
// block forever on a Read that will never receive input. Uses the
// FileMode character-device bit instead of pulling in
// golang.org/x/term as a direct dependency.
func isStdinInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {

		return false
	}

	return (info.Mode() & os.ModeCharDevice) != 0
}

// readUserConfirmation prints the prompt to stderr and reads one
// trimmed, lower-cased line from stdin. Only "y" proceeds; anything
// else (including the empty default and EOF) aborts. The strict
// match also makes `yes | gitmap reclone --execute` work the same
// way it would for any other "[y/N]" prompt.
func readUserConfirmation() bool {
	fmt.Fprint(os.Stderr, constants.MsgCloneNowConfirmPrompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {

		return false
	}
	answer := strings.ToLower(strings.TrimSpace(line))

	return answer == constants.CloneNowConfirmYes
}
