package constants

// CSV column names + branch validation error for clone-from.
// Split out of constants_clonefrom.go to keep that file under the
// 200-line per-file style budget. The column names are the SINGLE
// source of truth shared by:
//
//   - CSV header parsing (clonefrom/parsecsv.go indexCSVHeader)
//   - Per-row error wrapping (clonefrom/parsecsv.go wrapCSVRowErr)
//   - Per-field validation (clonefrom/validate.go validateRowWithColumn)
//
// Centralizing them ensures a header rename cannot drift the error
// messages out of sync with the actual accepted column names.

const (
	CSVColumnURL      = "url"
	CSVColumnDest     = "dest"
	CSVColumnBranch   = "branch"
	CSVColumnDepth    = "depth"
	CSVColumnCheckout = "checkout"

	// ErrCloneFromBadBranch fires at parse time for branch values
	// that are unambiguously invalid (leading dash → would be parsed
	// as a git flag; whitespace or control characters → cannot be a
	// real ref). Catching these at parse time surfaces the failure
	// with row + column context instead of as an opaque `git
	// checkout` error mid-clone. Stricter ref-name rules (no '..',
	// no '@{', etc.) are deliberately left to git itself.
	// %s = bad branch value.
	ErrCloneFromBadBranch = "branch %q is not a valid git ref name " +
		"(must not be empty after trim, contain whitespace, " +
		"start with '-', or contain control characters)"
)
