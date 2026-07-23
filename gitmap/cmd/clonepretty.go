// Package cmd — clonepretty.go: shared colorful runner for every
// `git clone` invocation triggered by gitmap (clone, cfr, cfrp,
// clone-replace temp swap). Unifies the log formatting requested in
// the v6.49.0 spec: cyan headers, green ✓ on success with elapsed
// time, red panel + retry hints on failure, and a `--dry-run` short
// circuit that prints the exact command without executing it.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Package-level dry-run + spinner state. Set by the top-level
// command handlers (runClone, runCloneFixRepoPipeline) before any
// `git clone` runs. Defaults to off so existing callers/tests keep
// their byte-for-byte legacy behavior unless the flag is opted in.
var (
	cloneDryRunFlag  bool
	cloneSpinnerOff  bool
	isCloneAssumeYes bool
)

// SetCloneDryRun toggles the dry-run short circuit for every
// subsequent runCloneCommand call in this process. Exported lowercase
// (package-private) intentionally — only the cmd package wires it.
func SetCloneDryRun(on bool) { cloneDryRunFlag = on }

// IsCloneDryRun reports the current dry-run flag state. Callers
// outside runCloneCommand (e.g. cfr's chained fix-repo step) can
// branch on it to suppress destructive follow-ups in dry-run mode.
func IsCloneDryRun() bool { return cloneDryRunFlag }

// SetCloneAssumeYes toggles auto-accept-new-host-key behavior for SSH
// clone commands when the user passes -y / --yes.
func SetCloneAssumeYes(on bool) { isCloneAssumeYes = on }

// SetCloneSpinnerOff disables the inline spinner. Useful in tests
// or CI where carriage-return updates clutter captured output.
func SetCloneSpinnerOff(off bool) { cloneSpinnerOff = off }

// runCloneCommandPretty replaces the bare `exec.Command("git",
// "clone", url, dest).Run()` pattern with a colorful, timed, and
// dry-run-aware runner. The exact git argv is printed up-front so
// the user always sees what would (or did) run.
func runCloneCommandPretty(url, dest string) error {
	printClonePrettyHeader(url, dest)
	if cloneDryRunFlag {
		fmt.Println(constants.MsgCloneDryRunNoop)
		return nil
	}
	stopSpinner := startCloneSpinnerForURL(url)
	start := time.Now()
	cmd := newCloneCommand(url, dest)
	runErr := cmd.Run()
	stopSpinner()
	elapsed := time.Since(start).Truncate(time.Millisecond)
	if runErr != nil {
		printClonePrettyFailure(url, dest, runErr, elapsed)
		return runErr
	}
	fmt.Printf(constants.MsgClonePrettyOK, dest, elapsed)
	return nil
}

func startCloneSpinnerForURL(url string) func() {
	if isSSHCloneURL(url) {
		return func() {}
	}

	return startCloneSpinner(constants.MsgCloneSpinnerLabel)
}

func newCloneCommand(url, dest string) *exec.Cmd {
	cmd := exec.Command(constants.GitBin, constants.GitClone, url, dest)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if isCloneAssumeYes && isSSHCloneURL(url) {
		cmd.Env = cloneEnvWithSSHAcceptNew()
	}

	return cmd
}

// printClonePrettyHeader emits the cyan banner shared by every
// clone-class command. Matches the visual language already used by
// the chrome-profile-copy and release pipelines.
func printClonePrettyHeader(url, dest string) {
	fmt.Printf(constants.MsgClonePrettyHeader, url, dest,
		constants.GitBin, constants.GitClone, url, dest)
}

// printClonePrettyFailure renders the red failure panel: exact
// argv that ran, exit code (when available), and a list of concrete
// retry hints. Hint set was tuned for the failure modes we see most
// often in support tickets (auth, transport, half-cloned trees).
func printClonePrettyFailure(url, dest string, runErr error, elapsed time.Duration) {
	code := -1
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		code = exitErr.ExitCode()
	}
	cmdline := fmt.Sprintf("%s %s %s %s", constants.GitBin, constants.GitClone, url, dest)
	fmt.Fprintf(os.Stderr, constants.MsgClonePrettyFail,
		cmdline, code, runErr, elapsed, buildClonePrettyHints(url, dest))
}

// buildClonePrettyHints picks retry suggestions based on URL shape.
// Returns a single newline-joined string so the format-string in
// constants_clone_pretty.go stays a single %s slot.
func buildClonePrettyHints(url, dest string) string {
	hints := []string{
		fmt.Sprintf("  • clean up: "+constants.ColorCyan+"rm -rf %q"+constants.ColorReset, dest),
		"  • retry without replace:    " + constants.ColorCyan + "gitmap clone " + url + " --no-replace" + constants.ColorReset,
	}
	if strings.HasPrefix(strings.ToLower(url), "https://") {
		hints = append(hints, "  • switch transport (SSH):   "+constants.ColorCyan+"gitmap clone "+url+" --ssh"+constants.ColorReset)
	}
	if strings.HasPrefix(strings.ToLower(url), "git@") ||
		strings.HasPrefix(strings.ToLower(url), "ssh://") {
		hints = append(hints, "  • switch transport (HTTPS): "+constants.ColorCyan+"gitmap clone "+url+" --https"+constants.ColorReset)
		hints = append(hints, "  • accept new SSH host key: "+constants.ColorCyan+"gitmap clone "+url+" -y"+constants.ColorReset)
	}
	hints = append(hints,
		"  • preview without cloning:  "+constants.ColorCyan+"gitmap clone "+url+" --dry-run"+constants.ColorReset)
	return strings.Join(hints, "\n")
}
