// Package cmd — codingguidelines_commit.go: auto-commit + push step
// that runs immediately after RunCodingGuidelinesInstall succeeds.
//
// Contract:
//   - Runs `git add -A` in the cloned working tree, then checks
//     porcelain output. If the installer produced no working-tree
//     changes we short-circuit (MsgCGNoChanges) instead of creating
//     an empty commit.
//   - Otherwise creates one commit with CodingGuidelinesCommitMessage
//     and, unless NoPush is set OR no upstream is configured, pushes
//     to the current branch's upstream.
//   - Errors are logged to opts.Stderr in the zero-swallow format
//     before being returned so the caller can decide whether to
//     halt the pipeline (cfr currently exits with the standard
//     chain-failed code).
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// CGCommitOpts controls a single post-install commit/push run.
// Zero-value opts default to real exec + os stdio.
type CGCommitOpts struct {
	WorkingDir string
	NoCommit   bool
	NoPush     bool
	Runner     func(name string, args ...string) *exec.Cmd
	Stdout     io.Writer
	Stderr     io.Writer
}

// CommitCodingGuidelines stages, commits, and (optionally) pushes any
// working-tree changes produced by the v24 installer.
func CommitCodingGuidelines(opts CGCommitOpts) error {
	opts = withCGCommitDefaults(opts)
	if opts.NoCommit {
		emitCGSkipNotes(opts, true, opts.NoPush)
		return nil
	}
	if err := runGitStep(opts, "add", "-A"); err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGCommitFailed, err)
		return err
	}
	dirty, err := hasStagedChanges(opts)
	if err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGCommitFailed, err)
		return err
	}
	if !dirty {
		fmt.Fprint(opts.Stderr, constants.MsgCGNoChanges)
		return nil
	}
	return commitAndMaybePush(opts)
}

func withCGCommitDefaults(opts CGCommitOpts) CGCommitOpts {
	if opts.Runner == nil {
		opts.Runner = exec.Command
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	return opts
}

// runGitStep runs a git subcommand with stdio streamed through the
// caller's writers so users see live output.
func runGitStep(opts CGCommitOpts, args ...string) error {
	cmd := opts.Runner("git", args...)
	cmd.Dir = opts.WorkingDir
	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr
	return cmd.Run()
}

// hasStagedChanges returns true when `git status --porcelain` reports
// at least one staged or unstaged path. We capture stdout separately
// (not streamed) because the parser needs the raw bytes.
func hasStagedChanges(opts CGCommitOpts) (bool, error) {
	cmd := opts.Runner("git", "status", "--porcelain")
	cmd.Dir = opts.WorkingDir
	cmd.Stderr = opts.Stderr
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(out))) > 0, nil
}

func commitAndMaybePush(opts CGCommitOpts) error {
	if err := runGitStep(opts, "commit", "-m", constants.CodingGuidelinesCommitMessage); err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGCommitFailed, err)
		return err
	}
	fmt.Fprintf(opts.Stderr, constants.MsgCGCommitted, constants.CodingGuidelinesCommitMessage)
	if opts.NoPush {
		emitCGSkipNotes(opts, false, true)
		return nil
	}
	upstream, ok := detectUpstream(opts)
	if !ok {
		emitCGSkipNotes(opts, false, true)
		return nil
	}
	if err := runGitStep(opts, "push"); err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGPushFailed, err)
		return err
	}
	fmt.Fprintf(opts.Stderr, constants.MsgCGPushed, upstream)
	return nil
}

// detectUpstream returns the current branch's upstream ref (e.g.
// "origin/main") or false when no upstream is configured. Missing
// upstream is a normal state (fresh branch, detached HEAD) — not an
// error worth halting the pipeline for.
func detectUpstream(opts CGCommitOpts) (string, bool) {
	cmd := opts.Runner("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	cmd.Dir = opts.WorkingDir
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	ref := strings.TrimSpace(string(out))
	if len(ref) == 0 {
		return "", false
	}
	return ref, true
}

// emitCGSkipNotes prints the --no-commit and --no-push skip notices
// in a single atomic write so downstream CI captures (Tee-Object on
// Windows, tee on POSIX) never lose the second message to pipe-drain
// races. The notes are mirrored to Stdout as well because they are
// user-facing status, not errors, and some CI runners buffer stderr
// separately from stdout when interleaving through a merged pipe.
func emitCGSkipNotes(opts CGCommitOpts, noCommit, noPush bool) {
	var b strings.Builder
	if noCommit {
		b.WriteString(constants.MsgCGSkipCommit)
	}
	if noPush {
		b.WriteString(constants.MsgCGSkipPush)
	}
	if b.Len() == 0 {
		return
	}
	msg := b.String()
	fmt.Fprint(opts.Stderr, msg)
	fmt.Fprint(opts.Stdout, msg)
}
