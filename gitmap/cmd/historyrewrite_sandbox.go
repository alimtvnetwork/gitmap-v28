package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ensureFilterRepoInstalled exits 3 with an OS-appropriate install
// hint when `git filter-repo --version` cannot run.
func ensureFilterRepoInstalled() {
	cmd := exec.Command(constants.HistoryGitBin, "filter-repo", "--version")
	if err := cmd.Run(); err == nil {
		return
	}
	fmt.Fprint(os.Stderr, constants.HistoryErrNoFilterRepo)
	switch runtime.GOOS {
	case "darwin":
		fmt.Fprint(os.Stderr, constants.HistoryMsgInstallHintMac)
	case "windows":
		fmt.Fprint(os.Stderr, constants.HistoryMsgInstallHintWin)
	default:
		fmt.Fprint(os.Stderr, constants.HistoryMsgInstallHintLinux)
	}
	os.Exit(constants.HistoryExitNoFilterRepo)
}

// readOriginURL invokes `git remote get-url origin` in the cwd. Exits
// 2 when not in a repo or no origin is configured.
func readOriginURL() string {
	cmd := exec.Command(constants.HistoryGitBin, "remote", "get-url", constants.HistoryRemoteOrigin)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrNoOrigin, err)
		os.Exit(constants.HistoryExitNotInRepo)
	}
	url := strings.TrimSpace(string(out))
	if url == "" {
		fmt.Fprintf(os.Stderr, constants.HistoryErrNoOrigin, fmt.Errorf("empty origin URL"))
		os.Exit(constants.HistoryExitNotInRepo)
	}
	return url
}

// mirrorClone creates an os.MkdirTemp sandbox and runs
// `git clone --mirror <origin> <sandbox>` into it.
func mirrorClone(originURL string, opts historyOpts) string {
	sandbox, err := os.MkdirTemp("", constants.HistorySandboxPrefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrSandbox, err)
		os.Exit(constants.HistoryExitBadArgs)
	}
	if !opts.quiet {
		fmt.Fprintf(os.Stderr, constants.HistoryMsgPhaseClone, originURL, sandbox)
	}
	cmd := exec.Command(constants.HistoryGitBin, "clone", "--mirror", originURL, sandbox)
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(sandbox)
		fmt.Fprintf(os.Stderr, constants.HistoryErrMirrorClone, err)
		os.Exit(constants.HistoryExitFilterFailed)
	}
	return sandbox
}

// runFilterRepo dispatches to the per-mode runner.
func runFilterRepo(mode historyMode, sandbox string, paths []string,
	pinPayloads map[string][]byte, opts historyOpts,
) {
	if mode == historyModePurge {
		runFilterRepoPurge(sandbox, paths, opts)
		return
	}
	runFilterRepoPin(sandbox, paths, pinPayloads, opts)
}

// runFilterRepoPurge invokes filter-repo with --invert-paths --path P
// for every requested path.
func runFilterRepoPurge(sandbox string, paths []string, opts historyOpts) {
	if !opts.quiet {
		fmt.Fprintf(os.Stderr, constants.HistoryMsgPhaseFilterPurge, len(paths))
	}
	args := []string{"-C", sandbox, "filter-repo", "--force", "--invert-paths"}
	for _, p := range paths {
		args = append(args, "--path", p)
	}
	args = append(args, historyMessageArgs(opts, sandbox, paths)...)
	execFilterRepo(args)
}

// runFilterRepoPin generates a Python --blob-callback that swaps every
// historical blob for the path with the current bytes loaded from the
// working tree.
func runFilterRepoPin(sandbox string, paths []string,
	pinPayloads map[string][]byte, opts historyOpts,
) {
	if !opts.quiet {
		fmt.Fprintf(os.Stderr, constants.HistoryMsgPhaseFilterPin, len(paths))
	}
	manifest, err := writePinManifest(sandbox, paths, pinPayloads)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrManifest, err)
		os.Exit(constants.HistoryExitFilterFailed)
	}
	args := []string{
		"-C", sandbox, "filter-repo", "--force",
		"--blob-callback", buildPinCallbackPython(manifest),
	}
	args = append(args, historyMessageArgs(opts, sandbox, paths)...)
	execFilterRepo(args)
}

// historyMessageArgs returns the filter-repo args needed to rewrite
// commit messages of ONLY commits that touched one of `paths` in the
// PRE-rewrite history, leaving every other commit's message
// untouched. Returns nil when --message is empty.
//
// We pre-compute the set of touched commit SHAs via `git log` in the
// sandbox BEFORE filter-repo runs, then pass that set into the
// commit-callback. This is required for `purge` mode, where
// --invert-paths drops the file_changes for the target paths before
// the callback fires — the callback would otherwise see zero touches
// and rewrite no messages (smoke-test failure mode).
func historyMessageArgs(opts historyOpts, sandbox string, paths []string) []string {
	if opts.message == "" {
		return nil
	}
	touched := touchedCommitSHAs(sandbox, paths)
	return []string{"--commit-callback", buildScopedMessagePython(opts.message, touched)}
}

// touchedCommitSHAs returns the lowercase hex SHAs of every commit
// reachable from any ref in `sandbox` that touched at least one of
// `paths`. Empty slice on git failure (caller treats as "rewrite
// nothing", which is the safe default).
func touchedCommitSHAs(sandbox string, paths []string) []string {
	args := []string{"-C", sandbox, "log", "--all", "--format=%H", "--"}
	args = append(args, paths...)
	out, err := exec.Command(constants.HistoryGitBin, args...).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrFilterRepo, exitCodeOf(err), err.Error())
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	shas := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			shas = append(shas, ln)
		}
	}
	return shas
}

// buildScopedMessagePython renders a Python snippet for filter-repo's
// --commit-callback that only rewrites commit.message when the
// commit's pre-rewrite SHA (commit.original_id) is in the set of
// touched commits computed by `git log -- <paths>`. We rely on the
// pre-computed SHA set instead of inspecting commit.file_changes
// because purge mode (--invert-paths) strips those entries before the
// callback fires.
func buildScopedMessagePython(message string, shas []string) string {
	quoted := make([]string, 0, len(shas))
	for _, s := range shas {
		quoted = append(quoted, fmt.Sprintf("%q", s))
	}
	setLiteral := "{" + strings.Join(quoted, ", ") + "}"
	return fmt.Sprintf(`
_touched_shas = %s
_oid = commit.original_id
if _oid is not None:
    try:
        _oid = _oid.decode("ascii")
    except Exception:
        _oid = ""
    if _oid in _touched_shas:
        commit.message = b%q
`, setLiteral, message)
}

// execFilterRepo runs `git ...` with stdio inherited and exits 5 on
// non-zero. Caller assembles the full arg vector.
func execFilterRepo(args []string) {
	cmd := exec.Command(constants.HistoryGitBin, args...)
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrFilterRepo, exitCodeOf(err), err.Error())
		os.Exit(constants.HistoryExitFilterFailed)
	}
}

// exitCodeOf extracts the process exit code from an exec.ExitError, or
// returns -1 when the error is something else.
func exitCodeOf(err error) int {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}
