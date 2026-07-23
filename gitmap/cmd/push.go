package cmd

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// transportFlags holds the parsed --ssh / --https selection plus any
// remaining positional args to forward to the underlying git command.
type transportFlags struct {
	useSSH   bool
	useHTTPS bool
	rest     []string
}

// parseTransportFlags parses the shared --ssh/--https flag pair used
// by `gitmap push` and the `gitmap pull` cwd short-circuit.
func parseTransportFlags(cmdName string, args []string) transportFlags {
	fs := flag.NewFlagSet(cmdName, flag.ExitOnError)
	sshFlag := fs.Bool("ssh", false, "Rewrite remote.origin.url to SSH and persist via `git remote set-url`")
	fs.BoolVar(sshFlag, "sh", false, "Short alias for --ssh")
	httpsFlag := fs.Bool("https", false, "Rewrite remote.origin.url to HTTPS and persist via `git remote set-url`")
	fs.BoolVar(httpsFlag, "ht", false, "Short alias for --https")
	fs.BoolVar(httpsFlag, "pub", false, "Alias for --https (public HTTPS clone URL)")

	fs.Parse(reorderFlagsBeforeArgs(args))

	return transportFlags{useSSH: *sshFlag, useHTTPS: *httpsFlag, rest: fs.Args()}
}

// runPush is the entry point for `gitmap push`. Short-circuits to
// `git push` in the cwd, with optional remote rewrite and auto
// pull --rebase + retry on non-fast-forward rejection.
func runPush(args []string) {
	checkHelp(constants.CmdPush, args)
	requireOnline()
	tf := parseTransportFlags(constants.CmdPush, args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getwd: %v\n", err)
		exitWith(1)

		return
	}
	if !isGitRepoCWD() {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap push` inside a repo)")
		exitWith(1)

		return
	}

	if _, _, _, applyErr := ApplyTransportFlag(cwd, tf.useSSH, tf.useHTTPS); applyErr != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", applyErr)
		exitWith(1)

		return
	}

	pushWithAutoRebase(cwd, tf.rest)
}

// pushWithAutoRebase runs `git push`, and on non-fast-forward
// rejection auto-runs `git pull --rebase` then retries once.
func pushWithAutoRebase(cwd string, rest []string) {
	gitArgs := append([]string{"push"}, rest...)
	fmt.Printf("→ Running: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	runErr, stderr := runGitCapturingStderr(gitArgs)
	if runErr == nil {
		return
	}
	if !isNonFastForwardRejection(stderr) {
		handleGitExit("git push", runErr)

		return
	}

	fmt.Fprintln(os.Stderr, "↻ push rejected (non-fast-forward) — auto-running `git pull --rebase` and retrying")
	if pullErr := runGitInherit([]string{"pull", "--rebase"}); pullErr != nil {
		fmt.Fprintln(os.Stderr, "✗ auto pull --rebase failed — resolve conflicts then re-run `gitmap push`")
		handleGitExit("git pull --rebase", pullErr)

		return
	}
	fmt.Printf("→ Retrying: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	if retryErr := runGitInherit(gitArgs); retryErr != nil {
		handleGitExit("git push (retry)", retryErr)
	}
}

// runGitCapturingStderr runs git, streaming stdio to the user while
// tee-ing stderr to a buffer for pattern matching.
func runGitCapturingStderr(gitArgs []string) (error, string) {
	var buf bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)

	return cmd.Run(), buf.String()
}

// runGitInherit runs git with stdio fully inherited.
func runGitInherit(gitArgs []string) error {
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// isNonFastForwardRejection matches git's canonical non-fast-forward
// rejection markers in captured stderr.
func isNonFastForwardRejection(stderr string) bool {
	lower := strings.ToLower(stderr)
	if !strings.Contains(lower, "[rejected]") && !strings.Contains(lower, "failed to push some refs") {
		return false
	}

	return strings.Contains(lower, "fetch first") || strings.Contains(lower, "non-fast-forward")
}

// handleGitExit translates a git ExitError into a clean process exit
// preserving the underlying status code.
func handleGitExit(label string, runErr error) {
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		exitWith(exitErr.ExitCode())

		return
	}
	fmt.Fprintf(os.Stderr, "%s failed: %v\n", label, runErr)
	exitWith(1)
}

// joinForLog renders argv as a space-separated banner for stdout.
func joinForLog(args []string) string {
	out := ""
	for i, a := range args {
		if i > 0 {
			out += " "
		}
		out += a
	}

	return out
}
