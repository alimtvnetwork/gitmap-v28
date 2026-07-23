package formatter

// Pre-write validation for ScanRecord slices destined for JSON/CSV output.
//
// Policy (warn-and-write, decided in v3.43.0):
//
//   - ValidateRecords NEVER fails the write. It returns a slice of issues
//     so the caller can surface them to stderr while still flushing the
//     file to disk. Partial output is preferred over silent corruption,
//     and full silence is preferred over no output at all.
//   - "Required" means: RepoName, RelativePath, AND at least one of
//     HTTPSUrl / SSHUrl. A repo with no URL cannot be re-cloned, so the
//     downstream `gitmap clone` round-trip would fail anyway.
//   - "Consistent" means: when both Slug and RepoName are populated,
//     Slug must equal strings.ToLower(RepoName). The DB upsert path
//     keys on Slug; drift here causes duplicate rows on re-scan.

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// ValidationIssue describes one problem discovered in one record. Rendered
// straight to stderr by the writer; carries enough context for users to
// fix the source data without re-running the scan.
type ValidationIssue struct {
	RowIndex int    // 0-based position in the input slice
	RepoName string // best-effort identifier (may be "" when RepoName is the missing field)
	Field    string // the field that triggered the issue
	Reason   string // human-readable explanation
}

// String renders an issue as a single line suitable for stderr output.
func (v ValidationIssue) String() string {
	name := v.RepoName
	if len(name) == 0 {
		name = "<unnamed>"
	}

	return fmt.Sprintf("row %d (%s): %s — %s", v.RowIndex, name, v.Field, v.Reason)
}

// ValidateRecords inspects every record and returns the union of all
// completeness + consistency issues found. An empty return slice means
// the input is safe to encode.
func ValidateRecords(records []model.ScanRecord) []ValidationIssue {
	var issues []ValidationIssue
	for i, rec := range records {
		issues = append(issues, validateOne(i, rec)...)
	}

	return issues
}

// validateOne runs every check against a single record and returns its
// per-record issue list (possibly empty).
func validateOne(idx int, rec model.ScanRecord) []ValidationIssue {
	var out []ValidationIssue
	out = appendIfMissing(out, idx, rec, rec.RepoName, "RepoName", "required field is empty")
	out = appendIfMissing(out, idx, rec, rec.RelativePath, "RelativePath", "required field is empty")
	out = checkURLPresence(out, idx, rec)
	out = checkSlugConsistency(out, idx, rec)

	return out
}

// appendIfMissing records an issue when value is empty after trimming.
func appendIfMissing(issues []ValidationIssue, idx int, rec model.ScanRecord, value, field, reason string) []ValidationIssue {
	if len(strings.TrimSpace(value)) > 0 {
		return issues
	}

	return append(issues, ValidationIssue{
		RowIndex: idx,
		RepoName: rec.RepoName,
		Field:    field,
		Reason:   reason,
	})
}

// checkURLPresence flags a record that has neither HTTPSUrl nor SSHUrl —
// such a record cannot be re-cloned by `gitmap clone`.
func checkURLPresence(issues []ValidationIssue, idx int, rec model.ScanRecord) []ValidationIssue {
	if hasURL(rec) {
		return issues
	}

	return append(issues, ValidationIssue{
		RowIndex: idx,
		RepoName: rec.RepoName,
		Field:    "HTTPSUrl|SSHUrl",
		Reason:   "record has no clone URL — downstream `gitmap clone` will skip it",
	})
}

// hasURL reports whether at least one clone URL is populated.
func hasURL(rec model.ScanRecord) bool {
	return len(strings.TrimSpace(rec.HTTPSUrl)) > 0 ||
		len(strings.TrimSpace(rec.SSHUrl)) > 0
}

// checkSlugConsistency verifies Slug matches the lowercased RepoName when
// both are populated. Mismatches break the DB dedupe key.
func checkSlugConsistency(issues []ValidationIssue, idx int, rec model.ScanRecord) []ValidationIssue {
	if len(rec.Slug) == 0 || len(rec.RepoName) == 0 {
		return issues
	}
	want := strings.ToLower(rec.RepoName)
	if rec.Slug == want {
		return issues
	}

	return append(issues, ValidationIssue{
		RowIndex: idx,
		RepoName: rec.RepoName,
		Field:    "Slug",
		Reason:   fmt.Sprintf("slug %q does not match lowercased RepoName %q", rec.Slug, want),
	})
}
