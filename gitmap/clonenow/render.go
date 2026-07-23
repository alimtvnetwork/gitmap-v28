package clonenow

// Dry-run renderer + post-execute summary. Both are intentionally
// monochrome (no ANSI) so the output is grep-friendly and identical
// across TTY / pipe / file targets. Callers that want colored
// per-repo blocks should pass through render.RenderRepoTermBlocks
// at the call site instead -- this package keeps the contract
// minimal so it can be re-used from non-CLI surfaces (tests,
// programmatic invocations) without a renderer dependency.

import (
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Render writes the dry-run preview: a header echoing the source +
// format + mode, then one block per row showing the URL that would
// be used and the destination path it would land at. Returns any
// io.Writer error so the CLI can exit 1 on a broken pipe instead of
// silently truncating.
func Render(w io.Writer, plan Plan) error {
	if _, err := fmt.Fprintf(w, constants.MsgCloneNowDryHeader,
		plan.Source, plan.Format, plan.Mode, len(plan.Rows)); err != nil {
		return err
	}
	for i, r := range plan.Rows {
		if err := renderRow(w, i+1, len(plan.Rows), r, plan.Mode); err != nil {
			return err
		}
	}

	return nil
}

// renderRow writes one row's dry-run block. Kept separate from
// Render so the per-row format can evolve (add branch, depth, ...)
// without churning the header logic.
func renderRow(w io.Writer, idx, total int, r Row, mode string) error {
	url := r.PickURL(mode)
	if len(url) == 0 {
		url = "(" + constants.MsgCloneNowNoURL + ")"
	}
	branch := r.Branch
	if len(branch) == 0 {
		branch = "(default)"
	}
	_, err := fmt.Fprintf(w,
		"  [%d/%d] %s\n        url:    %s\n        dest:   %s\n        branch: %s\n",
		idx, total, r.RepoName, url, r.RelativePath, branch)

	return err
}

// RenderSummary writes the end-of-batch tally + the per-row status
// table. Always prints, even when every row is "ok", so users have a
// definitive "this is what happened" record without hunting through
// the per-row progress lines.
func RenderSummary(w io.Writer, results []Result) error {
	ok, skipped, failed := tally(results)
	if _, err := fmt.Fprintf(w, constants.MsgCloneNowSummaryHeader,
		ok, skipped, failed, len(results)); err != nil {
		return err
	}
	for _, r := range results {
		if _, err := fmt.Fprintf(w, "  %-7s %s -> %s",
			r.Status, r.URL, r.Dest); err != nil {
			return err
		}
		if len(r.Detail) > 0 {
			if _, err := fmt.Fprintf(w, "  (%s)", r.Detail); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	return nil
}

// tally counts results by status. Centralized so Render and the
// CLI exit-code helper share one source of truth -- a "skipped"
// count drift between the summary line and the exit code would be
// confusing to debug.
func tally(results []Result) (ok, skipped, failed int) {
	for _, r := range results {
		switch r.Status {
		case constants.CloneNowStatusOK:
			ok++
		case constants.CloneNowStatusSkipped:
			skipped++
		case constants.CloneNowStatusFailed:
			failed++
		}
	}

	return ok, skipped, failed
}
