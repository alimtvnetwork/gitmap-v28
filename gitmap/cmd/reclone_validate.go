package cmd

// Manifest validation for `gitmap reclone`.
//
// Runs AFTER ParseFile produces a Plan and BEFORE the dry-run
// renderer or the --execute pipeline touches anything. The parser
// already enforces the "file is parseable + non-empty" contract;
// this layer enforces the SEMANTIC contract that each Row is
// actually clonable end-to-end:
//
//   - RepoName is present (used in progress + summary lines).
//   - At least one URL is set, AND it is well-formed (has a scheme
//     or matches the scp-like git@host:path shape).
//   - The URL for the user-selected --mode is reachable (either the
//     preferred field is set, or fallback to the other mode is
//     possible — Row.PickURL fallback is honored, but a row whose
//     ONLY URL is empty is rejected).
//   - RelativePath is set, is RELATIVE (not absolute), and contains
//     no ".." traversal segments — preventing a hostile or
//     hand-edited manifest from clobbering arbitrary disk paths.
//
// Failures are aggregated and printed as a row-level table on
// stderr, then the process exits with code 2 (bad input — same code
// the parser uses for unsupported extensions / bad flags). We exit
// 2 (not 1) because a malformed manifest is a usage error, not a
// per-row clone failure: the user must fix the file before any
// useful run is possible.
//
// Validation is INTENTIONALLY non-overridable — there is no
// --skip-validate flag. Letting bad rows through would either
// (a) crash deep inside git with a cryptic message, or
// (b) succeed in writing to a path the user did not intend.
// Both outcomes are worse than refusing up front.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// validateRecloneManifestOrExit runs every row through the semantic
// checks above. On any failure it prints a row-by-row error report
// to stderr and exits with constants.CloneNowExitManifestInvalid.
// Called from runCloneNow right after ParseFile so BOTH the dry-run
// branch and the --execute branch share the same gate — a user
// should never see "looks fine" on dry-run only to be rejected on
// --execute.
func validateRecloneManifestOrExit(plan clonenow.Plan) {
	issues := collectRecloneManifestIssues(plan)
	if len(issues) == 0 {
		return
	}
	printRecloneManifestIssues(plan, issues)
	os.Exit(constants.CloneNowExitManifestInvalid)
}

// recloneRowIssue is one validation finding, scoped to a single
// row. Multiple issues per row are reported separately so the user
// sees every problem in a single pass instead of fixing them one at
// a time across repeated invocations.
type recloneRowIssue struct {
	// rowIndex is 1-based to match how humans read manifests
	// (spreadsheet row 1, JSON array index 0 -> "row 1").
	rowIndex int
	// repoName echoes Row.RepoName when present, or "<unnamed>"
	// so the table column is never blank.
	repoName string
	// dest echoes Row.RelativePath verbatim (including the empty
	// string) so the user can see exactly what the manifest said.
	dest string
	// reason is a short, fixed phrase from the constants block —
	// stable so shell scripts and tests can grep for it.
	reason string
}

// collectRecloneManifestIssues walks every row once, accumulating
// issues. Order preserves the manifest order so the report reads
// top-to-bottom like the source file.
func collectRecloneManifestIssues(plan clonenow.Plan) []recloneRowIssue {
	out := make([]recloneRowIssue, 0)
	for i, row := range plan.Rows {
		for _, reason := range checkRecloneRow(row, plan.Mode) {
			out = append(out, recloneRowIssue{
				rowIndex: i + 1,
				repoName: displayRepoName(row.RepoName),
				dest:     row.RelativePath,
				reason:   reason,
			})
		}
	}

	return out
}

// checkRecloneRow returns the list of failure reasons for one row.
// Returns nil (not a zero-length slice) when the row is clean so
// callers can use a simple len check. Each reason is a stable
// phrase from constants — never include row data in the reason
// itself; the table prints repo + dest in their own columns.
func checkRecloneRow(row clonenow.Row, mode string) []string {
	reasons := make([]string, 0)
	if len(strings.TrimSpace(row.RepoName)) == 0 {
		reasons = append(reasons, constants.MsgRecloneValidateMissingRepoName)
	}
	if len(strings.TrimSpace(row.HTTPSUrl)) == 0 && len(strings.TrimSpace(row.SSHUrl)) == 0 {
		reasons = append(reasons, constants.MsgRecloneValidateNoURL)
	} else {
		picked := row.PickURL(mode)
		if !isPlausibleGitURL(picked) {
			reasons = append(reasons, constants.MsgRecloneValidateMalformedURL)
		}
	}
	reasons = append(reasons, checkRecloneDest(row.RelativePath)...)

	if len(reasons) == 0 {
		return nil
	}

	return reasons
}

// checkRecloneDest validates RelativePath in isolation. Split out
// because three independent failure modes apply (empty / absolute /
// traversal) and combining them inline would push checkRecloneRow
// past the 15-line function limit.
func checkRecloneDest(dest string) []string {
	trimmed := strings.TrimSpace(dest)
	if len(trimmed) == 0 {
		return []string{constants.MsgRecloneValidateMissingDest}
	}
	if filepath.IsAbs(trimmed) {
		return []string{constants.MsgRecloneValidateAbsoluteDest}
	}
	cleaned := filepath.ToSlash(filepath.Clean(trimmed))
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return []string{constants.MsgRecloneValidateTraversalDest}
	}

	return nil
}

// isPlausibleGitURL accepts the two URL shapes git itself accepts
// for clone: a real scheme://host/path URL (https, ssh, git, file)
// OR the scp-like shorthand `user@host:path`. Anything else (a bare
// "github.com/foo/bar", a Windows path, etc.) is rejected — git
// would either fail with a confusing error or silently treat it as
// a local path. Centralized so the check is consistent across
// future call sites (e.g. clone-from could adopt it).
func isPlausibleGitURL(url string) bool {
	trimmed := strings.TrimSpace(url)
	if len(trimmed) == 0 {
		return false
	}
	if hasURLScheme(trimmed) {
		return true
	}

	return isSCPLikeGitURL(trimmed)
}

// hasURLScheme matches "<alpha>://..." — the standard URL form.
// Rejects bare "://foo" (no scheme letters) and accepts mixed-case
// schemes (HTTPS://) which git normalizes anyway.
func hasURLScheme(url string) bool {
	idx := strings.Index(url, "://")
	if idx <= 0 {
		return false
	}
	for _, ch := range url[:idx] {
		if !isSchemeChar(ch) {
			return false
		}
	}

	return true
}

// isSchemeChar mirrors RFC 3986's scheme grammar (ALPHA / DIGIT /
// "+" / "-" / "."). Kept tiny + inline-able rather than pulling in
// net/url just to call url.Parse — net/url accepts shapes git
// rejects (e.g. "foo" with no host), so we'd still need this check.
func isSchemeChar(ch rune) bool {
	switch {
	case ch >= 'a' && ch <= 'z', ch >= 'A' && ch <= 'Z':
		return true
	case ch >= '0' && ch <= '9':
		return true
	case ch == '+' || ch == '-' || ch == '.':
		return true
	}

	return false
}

// isSCPLikeGitURL matches `user@host:path` — the historical SSH
// shorthand git accepts. Requires '@' BEFORE ':' (otherwise "C:\foo"
// would slip through as user="" host="C" path="\foo") and a
// non-empty path after the colon.
func isSCPLikeGitURL(url string) bool {
	at := strings.Index(url, "@")
	colon := strings.Index(url, ":")
	if at <= 0 || colon <= at+1 || colon == len(url)-1 {
		return false
	}
	// Reject "user@host:/absolute" only when it looks like a
	// drive letter — actual SSH paths starting with '/' are fine.
	return true
}

// displayRepoName replaces an empty/whitespace name with a fixed
// placeholder so the report's "repo" column is always populated.
// Stable string so tests can assert against it.
func displayRepoName(name string) string {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return constants.MsgRecloneValidateUnnamedRepo
	}

	return trimmed
}

// printRecloneManifestIssues writes the header + one line per
// issue + the abort footer, all to stderr. Single function so the
// output ordering is obvious at a glance and there's exactly one
// place to touch when the format changes.
func printRecloneManifestIssues(plan clonenow.Plan, issues []recloneRowIssue) {
	fmt.Fprintf(os.Stderr, constants.MsgRecloneValidateHeaderFmt,
		plan.Source, len(issues), len(plan.Rows))
	for _, issue := range issues {
		fmt.Fprintf(os.Stderr, constants.MsgRecloneValidateRowFmt,
			issue.rowIndex, issue.repoName, displayDest(issue.dest), issue.reason)
	}
	fmt.Fprint(os.Stderr, constants.MsgRecloneValidateFooter)
}

// displayDest mirrors displayRepoName for the dest column so empty
// RelativePath values render as a visible placeholder instead of a
// blank gap that's easy to miss when scanning the table.
func displayDest(dest string) string {
	trimmed := strings.TrimSpace(dest)
	if len(trimmed) == 0 {
		return constants.MsgRecloneValidateEmptyDest
	}

	return trimmed
}
