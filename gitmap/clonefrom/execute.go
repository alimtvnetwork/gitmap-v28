package clonefrom

// Executor: walk Plan.Rows sequentially, shell out to `git clone`
// for each, and accumulate per-row Results. Sequential by design
// — parallel fan-out is a follow-up (see _TODO in
// .lovable/question-and-ambiguity/03-clone-from-scope.md). Adding
// it later is a one-function change because Result has no shared
// state between rows.
//
// Skip rule: if the resolved dest already exists AND is a non-
// empty directory → mark as `skipped` without invoking git. Makes
// re-running the same plan idempotent (common pattern: user fixes
// a typo in row 4, re-runs, doesn't want rows 1-3 re-cloned).
//
// We deliberately do NOT try to detect "dest exists AND points at
// the same URL" — that requires `git remote get-url` which (a)
// adds latency to the skip check and (b) would make the rule
// behave differently on partially-cloned dests. Conservative
// "non-empty dir = skip" is easier to reason about.

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Result is one row's outcome. Status is one of "ok" | "skipped"
// | "failed". Detail is human-readable context: for "ok" the empty
// string; for "skipped" the reason ("dest exists"); for "failed"
// the trimmed git stderr (capped at GitErrorTrimLimit chars to
// keep the summary table readable).
type Result struct {
	Row      Row
	Dest     string // resolved (after DeriveDest fallback)
	Status   string
	Detail   string
	Duration time.Duration
}

// Execute runs every row sequentially and writes per-row progress
// lines to progress (typically os.Stderr; pass io.Discard to
// silence). Returns the full result slice — callers render the
// summary AFTER Execute returns so a write failure on progress
// doesn't truncate the result set.
//
// Working directory: each clone runs in `cwd`. Empty cwd → use
// the current process cwd at call time.
func Execute(plan Plan, cwd string, progress io.Writer) []Result {
	cwd = resolveCwd(cwd)
	out := make([]Result, 0, len(plan.Rows))
	for i, r := range plan.Rows {
		res := executeRow(r, cwd)
		out = append(out, res)
		writeProgress(progress, i+1, len(plan.Rows), res)
	}

	return out
}

// RowLifecycleParams encapsulates parameters for running a clone row lifecycle.
type RowLifecycleParams struct {
	Row     Row
	Dest    string
	AbsDest string
	Cwd     string
	Start   time.Time
}

// PrepareCloneParams encapsulates parameters for preparing parent dir and cloning.
type PrepareCloneParams struct {
	Row     Row
	Dest    string
	AbsDest string
	Cwd     string
}

// FailedResultParams encapsulates parameters for creating a failed clone result.
type FailedResultParams struct {
	Row    Row
	Dest   string
	Detail string
	Start  time.Time
}

// executeRow handles one row's lifecycle: resolve dest, check
// skip rule, ensure parent dir, git clone, checkout, time, return.
func executeRow(r Row, cwd string) Result {
	start := time.Now()
	dest, absDest := resolveDest(r, cwd)
	isSkipped := shouldSkip(absDest)
	if isSkipped {
		return Result{Row: r, Dest: dest, Status: constants.CloneFromStatusSkipped,
			Detail: constants.MsgCloneFromDestExists, Duration: time.Since(start)}
	}

	lifecycleParams := RowLifecycleParams{
		Row:     r,
		Dest:    dest,
		AbsDest: absDest,
		Cwd:     cwd,
		Start:   start,
	}

	return runRowLifecycle(lifecycleParams)
}

func runRowLifecycle(params RowLifecycleParams) Result {
	cloneParams := PrepareCloneParams{
		Row:     params.Row,
		Dest:    params.Dest,
		AbsDest: params.AbsDest,
		Cwd:     params.Cwd,
	}

	detail, ok := prepareAndClone(cloneParams)
	isFailed := !ok
	if isFailed {
		return makeFailedResult(FailedResultParams{Row: params.Row, Dest: params.Dest, Detail: detail, Start: params.Start})
	}
	coDetail, coOK := runPostCloneCheckout(params.Row, params.Dest, params.Cwd)
	isCheckoutFailed := !coOK
	if isCheckoutFailed {
		return makeFailedResult(FailedResultParams{Row: params.Row, Dest: params.Dest, Detail: coDetail, Start: params.Start})
	}
	return Result{Row: params.Row, Dest: params.Dest, Status: constants.CloneFromStatusOK,
		Detail: "", Duration: time.Since(params.Start)}
}

func prepareAndClone(params PrepareCloneParams) (string, bool) {
	detail, ok := prepareDestParent(params.AbsDest)
	isParentFailed := !ok
	if isParentFailed {
		return detail, false
	}
	return runGitClone(params.Row, params.Dest, params.Cwd)
}

func makeFailedResult(params FailedResultParams) Result {
	return Result{
		Row:      params.Row,
		Dest:     params.Dest,
		Status:   constants.CloneFromStatusFailed,
		Detail:   params.Detail,
		Duration: time.Since(params.Start),
	}
}

// runGitClone shells out to `git clone` with the row's options.
// Returns (detail, ok). On success detail is empty. On failure
// detail is a single-line summary of the trimmed stderr.
func runGitClone(r Row, dest, cwd string) (string, bool) {
	out, err := execGitClone(r, dest, cwd)
	isSuccess := err == nil
	if isSuccess {
		return "", true
	}
	return handleGitCloneError(dest, string(out), err)
}

func execGitClone(r Row, dest, cwd string) ([]byte, error) {
	args := buildGitArgs(r, dest)
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Dir = cwd
	return cmd.CombinedOutput()
}

func handleGitCloneError(dest, outputStr string, err error) (string, bool) {
	file, isSmudge := detectLFSSmudgeError(outputStr)
	if isSmudge {
		fixParams := LfsFixParams{
			Dest:        dest,
			File:        file,
			OutputStr:   outputStr,
			OriginalErr: err,
		}

		return tryLfsAutoFix(fixParams)
	}
	return trimGitError(outputStr, err), false
}

// LfsFixParams encapsulates parameters for attempting an LFS auto-fix.
type LfsFixParams struct {
	Dest        string
	File        string
	OutputStr   string
	OriginalErr error
}

func tryLfsAutoFix(params LfsFixParams) (string, bool) {
	fmt.Fprintf(os.Stderr, "\n[Warning] Git clone succeeded but checkout failed due to missing LFS object (404) for file: %s\n", params.File)
	confirmed := confirmYesNo("Do you want to automatically drop this broken LFS pointer to fix the clone and push the fix?")
	isDeclined := !confirmed
	if isDeclined {
		return trimGitError(params.OutputStr, params.OriginalErr), false
	}
	return applyLfsFix(params)
}

func applyLfsFix(params LfsFixParams) (string, bool) {
	fixErr := executeLFSFix(params.Dest, params.File)
	isFixFailed := fixErr != nil
	if isFixFailed {
		return trimGitError(params.OutputStr+"\n[LFS Fix Failed: "+fixErr.Error()+"]", params.OriginalErr), false
	}
	return "", true
}

func resolveCwd(cwd string) string {
	hasCwd := len(cwd) > 0
	if hasCwd {
		return cwd
	}
	wd, err := os.Getwd()
	isSuccess := err == nil
	if isSuccess {
		return wd
	}
	return ""
}

// buildGitArgs translates a Row + resolved dest into the git
// clone argument vector.
func buildGitArgs(r Row, dest string) []string {
	args := []string{constants.GitClone}
	args = appendBranchArg(args, r.Branch)
	args = appendDepthArg(args, r.Depth)
	args = appendCheckoutArg(args, r)
	return append(args, r.URL, dest)
}

func appendBranchArg(args []string, branch string) []string {
	hasBranch := len(branch) > 0
	if hasBranch {
		return append(args, constants.GitBranchFlag, branch)
	}
	return args
}

func appendDepthArg(args []string, depth int) []string {
	hasDepth := depth > 0
	if hasDepth {
		return append(args, fmt.Sprintf(constants.CloneFromDepthFlagFmt, depth))
	}
	return args
}

func appendCheckoutArg(args []string, r Row) []string {
	isSkipCheckout := EffectiveCheckout(r) == constants.CloneFromCheckoutSkip
	if isSkipCheckout {
		return append(args, constants.CloneFromNoCheckoutFlag)
	}
	return args
}

// trimGitError collapses multi-line git stderr to a single line
// with a length cap.
func trimGitError(stderr string, err error) string {
	last := extractLastStderrLine(stderr, err)
	isExceedingLimit := len(last) > constants.CloneFromErrTrimLimit
	if isExceedingLimit {
		return last[:constants.CloneFromErrTrimLimit] + "..."
	}
	return last
}

func extractLastStderrLine(stderr string, err error) string {
	last := stderr
	idx := strings.LastIndex(strings.TrimSpace(stderr), "\n")
	hasNewline := idx >= 0
	if hasNewline {
		last = strings.TrimSpace(stderr)[idx+1:]
	}
	last = strings.TrimSpace(last)
	isEmpty := len(last) == 0
	if isEmpty {
		return err.Error()
	}
	return last
}

// writeProgress emits one line per finished row.
func writeProgress(w io.Writer, n, total int, res Result) {
	isNilWriter := w == nil
	if isNilWriter {
		return
	}
	fmt.Fprintf(w, "  [%d/%d] %-7s %s\n", n, total, res.Status, res.Row.URL)
}
