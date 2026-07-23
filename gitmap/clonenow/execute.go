package clonenow

// Executor: walk Plan.Rows sequentially, shell out to `git clone`
// for each, accumulate per-row Results. Sequential by design so the
// progress lines arrive in stable order; parallelism is a future
// follow-up because Result has no shared state between rows.
//
// Idempotency policy (Plan.OnExists):
//
//   - "skip" (default): a destination that's already a git repo
//     pointing at the same URL on the same branch is reported as
//     `skipped` with detail "already matches". Any other state
//     (URL/branch mismatch, non-repo dir) is also skipped but with
//     a detail line that explains WHY -- so a re-run never silently
//     hides drift. The detail is the user's signal to re-run with
//     `--on-exists update` or `--on-exists force`.
//   - "update": same probe as skip; a matching repo is still a
//     no-op (`skipped: already matches`), a mismatched repo runs
//     `git fetch` + `git checkout <branch>` to align without
//     destroying local commits. Implemented in update_existing.go.
//   - "force": destination is removed (only after we confirm it's
//     a git repo or empty -- never blow away an unrelated dir),
//     then re-cloned from scratch. Implemented in force_reclone.go.
//
// All three policies share a single inspect step (inspectExistingRepo)
// so the "what's actually on disk?" question is asked exactly once
// per row, regardless of which branch fires next.

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Result is one row's outcome. Status is one of "ok" | "skipped"
// | "failed". Detail carries human-readable context: empty for ok,
// MsgCloneNowDestExists for skipped, the trimmed git stderr for
// failed (capped at CloneNowErrTrimLimit chars).
type Result struct {
	Row      Row
	URL      string // the URL actually used (after mode + fallback)
	Dest     string // the destination path actually used (after cwd join)
	Status   string
	Detail   string
	Duration time.Duration
}

// Execute runs every row sequentially and writes per-row progress
// lines to `progress` (typically os.Stderr; pass io.Discard to
// silence). Returns the full result slice -- callers render the
// summary AFTER Execute returns so a write failure on `progress`
// doesn't truncate the result set.
//
// Working directory: each clone runs in `cwd`. Empty `cwd` -> the
// process cwd at call time, captured once up-front so a row that
// shells out into a subprocess (which mutates cwd) doesn't shift
// where subsequent rows land.
func Execute(plan Plan, cwd string, progress io.Writer) []Result {
	if len(cwd) == 0 {
		if wd, err := os.Getwd(); err == nil {
			cwd = wd
		}
	}
	out := make([]Result, 0, len(plan.Rows))
	for i, r := range plan.Rows {
		res := executeRow(r, plan, cwd)
		out = append(out, res)
		writeProgress(progress, i+1, len(plan.Rows), res)
	}

	return out
}

// executeRow handles one row's lifecycle: pick URL, resolve dest,
// inspect what's on disk, dispatch to the on-exists policy, run.
//
// The dispatch is done here (rather than buried inside a sub-helper)
// so the four possible terminal states -- ok, skipped, failed,
// updated -- are visible at a glance when reading the executor.
func executeRow(r Row, plan Plan, cwd string) Result {
	start := time.Now()
	url := r.PickURL(plan.Mode)
	dest := r.RelativePath
	absDest := dest
	if !filepath.IsAbs(absDest) {
		absDest = filepath.Join(cwd, dest)
	}
	base := Result{Row: r, URL: url, Dest: dest}
	if len(url) == 0 {
		base.Status = constants.CloneNowStatusFailed
		base.Detail = constants.MsgCloneNowNoURL
		base.Duration = time.Since(start)

		return base
	}
	state := inspectExistingRepo(absDest)
	res := dispatchOnExists(r, url, absDest, cwd, plan.OnExists, state)
	res.Row = r
	res.URL = url
	res.Dest = dest
	res.Duration = time.Since(start)

	return res
}

// runGitClone shells out to `git clone` with the row's options.
// Returns (detail, ok). On success detail is empty. On failure
// detail is a single-line summary of the trimmed stderr.
//
// Pre-creates the destination's parent directory so nested
// RelativePath values (e.g. `org-a/team/repo-x`) succeed on a
// fresh checkout where the intermediate folders don't yet exist.
// MkdirAll is idempotent: a pre-existing parent is a no-op and
// safe under concurrent sibling clones. Failure is logged in the
// project's Code Red format AND surfaced as the row Detail so the
// per-row line + summary table carry the same diagnosis.
func runGitClone(r Row, url, dest, cwd string) (string, bool) {
	absDest := dest
	if !filepath.IsAbs(absDest) {
		absDest = filepath.Join(cwd, dest)
	}
	parent := filepath.Dir(absDest)
	if err := os.MkdirAll(parent, constants.DirPermission); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNowMkdirParent, parent, err)
		return fmt.Sprintf(constants.MsgCloneNowMkdirParentFailFmt, err), false
	}
	args := buildGitArgs(r, url, dest)
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Dir = cwd
	if isSSHCloneURL(url) {
		return runInteractiveGitClone(cmd)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return trimGitError(string(out), err), false
	}

	return "", true
}

// buildGitArgs assembles the `git clone [...] <url> <dest>` argv.
// Order matters for git's flag parser: --branch must precede the
// positional URL/dest pair.
func buildGitArgs(r Row, url, dest string) []string {
	args := []string{constants.GitClone}
	if len(r.Branch) > 0 {
		args = append(args, constants.GitBranchFlag, r.Branch)
	}
	args = append(args, url, dest)

	return args
}

// trimGitError collapses multi-line git stderr to a single line
// with a length cap so the summary table stays scannable. Full
// stderr is in the user's terminal scrollback already (we use
// CombinedOutput).
func trimGitError(stderr string, err error) string {
	last := strings.TrimSpace(stderr)
	if i := strings.LastIndex(last, "\n"); i >= 0 {
		last = strings.TrimSpace(last[i+1:])
	}
	if len(last) == 0 {
		last = err.Error()
	}
	if len(last) > constants.CloneNowErrTrimLimit {
		last = last[:constants.CloneNowErrTrimLimit] + "..."
	}

	return last
}

// writeProgress emits one line per finished row. nil writer is a
// no-op so callers can pass io.Discard (or nothing at all in tests)
// without sprinkling nil-checks at every call site.
func writeProgress(w io.Writer, n, total int, res Result) {
	if w == nil {
		return
	}
	fmt.Fprintf(w, "  [%d/%d] %-7s %s -> %s\n", n, total, res.Status, res.URL, res.Dest)
}
