package clonefrom

// Validation + dedup helpers shared by both parsers. Kept in a
// dedicated file so format-specific parsing (parse.go) stays
// focused and the validation rules have one obvious home for
// future tightening (e.g., adding a hostname allowlist).

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// validateRow enforces per-row invariants:
//
//   - URL non-empty after trim
//   - URL has a recognizable scheme (https://, http://, git://,
//     ssh://) OR matches the scp-style `user@host:path` shape that
//     git accepts but net/url rejects
//   - Depth non-negative (negative depths are rejected by git but
//     the user-facing error from us is clearer than git's)
//
// Dest is NOT validated for path safety here — we let git fail
// loudly at execute time with its own (well-tested) checks. Adding
// a duplicate validator here would just create a second source of
// truth that could drift from git's actual behavior.
func validateRow(r Row) error {
	_, err := validateRowWithColumn(r)
	return err
}

// validateRowWithColumn returns the offending CSV column name
// alongside the error so callers (the CSV parser) can name the bad
// cell in their wrapped error. Non-CSV callers use validateRow and
// discard the column name.
func validateRowWithColumn(r Row) (string, error) {
	if len(r.URL) == 0 {
		return constants.CSVColumnURL, fmt.Errorf(constants.ErrCloneFromEmptyURL)
	}
	if !looksLikeGitURL(r.URL) {
		return constants.CSVColumnURL, fmt.Errorf(constants.ErrCloneFromBadURL, r.URL)
	}
	if r.Depth < 0 {
		return constants.CSVColumnDepth, fmt.Errorf(constants.ErrCloneFromNegDepth, r.Depth)
	}
	if len(r.Branch) > 0 && !isValidBranchName(r.Branch) {
		return constants.CSVColumnBranch, fmt.Errorf(constants.ErrCloneFromBadBranch, r.Branch)
	}
	if !isValidCheckout(r.Checkout) {
		return constants.CSVColumnCheckout, fmt.Errorf(constants.ErrCloneFromBadCheckout, r.Checkout)
	}

	return "", nil
}

// isValidBranchName rejects values that would either be silently
// reinterpreted by git (a leading '-' becomes a flag) or that we
// know cannot be a real ref (whitespace / control chars). This is
// intentionally permissive — git's own ref-name rules are stricter
// (no '..', no '@{', etc.); we catch only the high-confidence
// failures so the error surfaces with row/column context instead
// of as an opaque `git checkout` failure later.
func isValidBranchName(s string) bool {
	if strings.HasPrefix(s, "-") {
		return false
	}
	for _, r := range s {
		if r <= 0x20 || r == 0x7f {
			return false
		}
	}

	return true
}

// isValidCheckout accepts the empty string (means "inherit global
// default") plus the three explicit modes. Centralized so both the
// row-level validator and the CLI flag validator share one truth.
func isValidCheckout(v string) bool {
	switch v {
	case "",
		constants.CloneFromCheckoutAuto,
		constants.CloneFromCheckoutSkip,
		constants.CloneFromCheckoutForce:

		return true
	}

	return false
}

// looksLikeGitURL is a permissive shape check. We deliberately do
// NOT use net/url.Parse because:
//
//   - scp-style `git@github.com:owner/repo.git` is valid for
//     `git clone` but parses as a relative path in net/url.
//   - we don't want to reject exotic but legal forms (file://,
//     git protocol over SSH tunnels, etc.) just because our
//     validator has a smaller world view than git itself.
//
// The check exists only to catch obvious typos like a bare
// "owner/repo" or a copy-pasted markdown link.
func looksLikeGitURL(s string) bool {
	if hasGitScheme(s) {

		return true
	}

	return looksLikeSCP(s)
}

// hasGitScheme returns true for any URL beginning with one of the
// schemes git supports natively.
func hasGitScheme(s string) bool {
	prefixes := []string{"https://", "http://", "ssh://", "git://", "file://"}
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {

			return true
		}
	}

	return false
}

// looksLikeSCP matches the `[user@]host:path` form. Cheap test:
// has a colon, no slashes before the colon, has at least one slash
// after the colon. Lets `git@github.com:owner/repo.git` through
// while rejecting `owner/repo.git` (no colon) and `://path` (slash
// before colon).
func looksLikeSCP(s string) bool {
	colon := strings.Index(s, ":")
	if colon <= 0 {

		return false
	}
	host := s[:colon]
	path := s[colon+1:]
	if strings.ContainsAny(host, "/\\") {

		return false
	}

	return strings.Contains(path, "/")
}

// dedupRows collapses rows with identical URL+Dest. Later rows
// overwrite earlier ones for branch/depth so users can re-list a
// URL further down the file to override a default — common
// spreadsheet workflow ("global defaults at top, exceptions at
// bottom"). Order of first occurrence is preserved so the dry-run
// preview lists rows in the user's original sequence.
func dedupRows(rows []Row) []Row {
	seen := make(map[string]int, len(rows))
	out := make([]Row, 0, len(rows))
	for _, r := range rows {
		key := r.URL + "\x00" + r.Dest
		if i, ok := seen[key]; ok {
			out[i] = mergeRows(out[i], r)
			continue
		}
		seen[key] = len(out)
		out = append(out, r)
	}

	return out
}

// mergeRows overlays the later row's optional fields onto the
// earlier row. Branch/depth/checkout take the later value if the
// later value is non-empty/non-zero.
func mergeRows(first, later Row) Row {
	out := first
	if len(later.Branch) > 0 {
		out.Branch = later.Branch
	}
	if later.Depth > 0 {
		out.Depth = later.Depth
	}
	if len(later.Checkout) > 0 {
		out.Checkout = later.Checkout
	}

	return out
}
