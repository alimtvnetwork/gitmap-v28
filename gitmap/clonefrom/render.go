package clonefrom

// Dry-run renderer. Pure function over a Plan; never touches disk
// or the network. Output is human-readable plain text designed to
// be grep-able and diff-able across runs of the same plan.
//
// Format (chosen for clarity over compactness — dry-run output is
// read once by a human, not parsed by tooling):
//
//   gitmap clone-from: dry-run
//   source: /abs/path/to/plan.csv (csv)
//   3 row(s) — pass --execute to actually clone
//
//     1. https://github.com/a/b.git
//        dest:   ./b
//        branch: (default HEAD)
//        depth:  full
//     2. ...
//
// Tooling that wants machine-readable dry-run output can pass
// --format=json to the CLI; that path renders a JSON array of the
// same Plan.Rows and is covered by a separate JSON Schema (see
// spec/08-json-schemas/_TODO.md entry for clone-from).

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
)

// Render writes the legacy human-readable dry-run preview to w.
// Returns the first write error so callers (CLI) can surface it
// instead of silently truncating output to a closed pipe.
func Render(w io.Writer, p Plan) error {
	header := fmt.Sprintf(constants.MsgCloneFromDryHeader, p.Source, p.Format, len(p.Rows))
	if _, err := io.WriteString(w, header); err != nil {

		return err
	}
	for i, r := range p.Rows {
		block := renderRow(i+1, r)
		if _, err := io.WriteString(w, block); err != nil {

			return err
		}
	}

	return nil
}

// RenderTerminal writes the standardized per-repo block (one block
// per row) using render.RenderRepoTermBlocks. The block matches the
// format emitted by scan, clone-next, and probe so users learn one
// shape regardless of which command produced it.
func RenderTerminal(w io.Writer, p Plan) error {
	header := fmt.Sprintf(constants.MsgCloneFromDryHeader, p.Source, p.Format, len(p.Rows))
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	blocks := make([]render.RepoTermBlock, 0, len(p.Rows))
	for i, r := range p.Rows {
		blocks = append(blocks, rowToBlock(i+1, r))
	}

	return render.RenderRepoTermBlocks(w, blocks)
}

// rowToBlock maps a clone-from Row into the standardized block.
// OriginalURL == TargetURL because clone-from does not rewrite URLs
// — the user's manifest specifies the URL verbatim.
func rowToBlock(n int, r Row) render.RepoTermBlock {
	dest := r.Dest
	if len(dest) == 0 {
		dest = DeriveDest(r.URL)
	}

	return render.RepoTermBlock{
		Index:        n,
		Name:         dest,
		Branch:       r.Branch,
		BranchSource: branchSourceForRow(r),
		OriginalURL:  r.URL,
		TargetURL:    r.URL,
		CloneCommand: cloneCommandForRow(r, dest),
	}
}

// branchSourceForRow surfaces "manifest" when the row pinned a
// branch explicitly, or "default HEAD" when git will pick.
func branchSourceForRow(r Row) string {
	if len(strings.TrimSpace(r.Branch)) > 0 {
		return "manifest"
	}

	return "default HEAD"
}

// cloneCommandForRow renders the would-be `git clone` invocation,
// matching buildGitArgs in execute.go (same flag order so the
// preview is faithful to the executed command).
func cloneCommandForRow(r Row, dest string) string {
	parts := []string{"git", "clone"}
	if len(strings.TrimSpace(r.Branch)) > 0 {
		parts = append(parts, "-b", r.Branch)
	}
	if r.Depth > 0 {
		parts = append(parts, fmt.Sprintf(constants.CloneFromDepthFlagFmt, r.Depth))
	}
	if EffectiveCheckout(r) == constants.CloneFromCheckoutSkip {
		parts = append(parts, constants.CloneFromNoCheckoutFlag)
	}
	parts = append(parts, r.URL, dest)

	return strings.Join(parts, " ")
}

// renderRow formats one row's block. The first four lines match the
// pre-checkout-feature output byte-for-byte (locked by render_test.go
// + downstream log scrapers); the optional fifth `checkout:` line is
// appended ONLY when the resolved mode is non-default ("skip" or
// "force") so existing fixtures and grep expectations stay valid for
// the common case.
func renderRow(n int, r Row) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "  %d. %s\n", n, r.URL)
	fmt.Fprintf(&sb, "     dest:   %s\n", displayDest(r))
	fmt.Fprintf(&sb, "     branch: %s\n", displayBranch(r))
	fmt.Fprintf(&sb, "     depth:  %s\n", displayDepth(r))
	if mode := EffectiveCheckout(r); mode != constants.CloneFromCheckoutDefault {
		fmt.Fprintf(&sb, "     checkout: %s\n", mode)
	}

	return sb.String()
}

// displayDest shows the resolved destination, marking implicit
// derivations with `(derived)` so the user knows we're computing
// the dest from the URL basename rather than honoring an explicit
// value. Helps when a typo'd URL produces a surprising dest.
func displayDest(r Row) string {
	if len(r.Dest) > 0 {

		return r.Dest
	}

	return DeriveDest(r.URL) + "  (derived)"
}

// DeriveDest computes the directory name `git clone <url>` would
// pick by default: the last path segment with a trailing `.git`
// stripped. Public so the executor can use the same logic when
// the row's Dest field is empty.
func DeriveDest(url string) string {
	// Strip scp-style user@host: prefix if present so path.Base
	// works on the right portion.
	if i := strings.LastIndex(url, ":"); i > 0 && !strings.HasPrefix(url, "ssh://") &&
		!strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = url[i+1:]
	}
	// A trailing slash means the URL has no real path segment
	// (e.g. "https://example.org/") — fall back to "repo" rather
	// than letting path.Base return the host name.
	trimmed := strings.TrimRight(url, "/")
	if !strings.Contains(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(trimmed, "https://"), "http://"), "ssh://"), "/") {

		return "repo"
	}
	base := path.Base(trimmed)
	base = strings.TrimSuffix(base, ".git")
	if len(base) == 0 || base == "." || base == "/" {

		return "repo"
	}

	return base
}

// displayBranch shows the explicit branch or a "(default HEAD)"
// placeholder. Placeholder uses parens to visually distinguish
// from a literal branch named e.g. `default`.
func displayBranch(r Row) string {
	if len(r.Branch) > 0 {

		return r.Branch
	}

	return "(default HEAD)"
}

// displayDepth shows "full" for zero-depth (full history) or the
// explicit shallow depth. "full" reads better than "0" in user-
// facing output.
func displayDepth(r Row) string {
	if r.Depth > 0 {

		return fmt.Sprintf("%d", r.Depth)
	}

	return "full"
}
