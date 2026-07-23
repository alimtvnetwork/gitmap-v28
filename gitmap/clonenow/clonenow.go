// Package clonenow implements the `gitmap clone-now <file>` workflow:
// re-run `git clone` against the JSON / CSV / text artifacts produced
// by `gitmap scan`, honoring the recorded folder structure and a
// user-selected SSH/HTTPS mode.
//
// Why a separate package (not clonefrom or cloner)?
//
//   - clonefrom is plan-driven (user-authored row schema).
//   - cloner is the in-memory scan-pipeline cloner that runs as part
//     of the scan command itself.
//   - clonenow is a round-trip cloner: the input is gitmap's own
//     scan output, so the schema (RepoName, HTTPSUrl, SSHUrl, Branch,
//     RelativePath, ...) is fixed and we honor RelativePath verbatim
//     so the destination tree is byte-identical to the original.
//
// Splitting keeps each cloner's contract tight and avoids forcing
// clone-from to grow scan-record-aware fields.
package clonenow

// Plan is the validated, in-memory representation of one input file.
// Built by ParseFile from JSON, CSV, or plain-text scan output and
// consumed by Render (dry-run) and Execute.
type Plan struct {
	// Source is the absolute path the plan was read from. Echoed
	// verbatim in the dry-run header.
	Source string
	// Format is "json" | "csv" | "text" -- used by the dry-run
	// header so the user can confirm we parsed the file the way
	// they expected.
	Format string
	// Mode is "https" | "ssh" -- the URL column the executor will
	// prefer for each row. Captured on the Plan (not just on the
	// CLI cfg) so the dry-run renderer and the executor see the
	// same value, with no risk of drift.
	Mode string
	// OnExists is "skip" | "update" | "force" -- the policy applied
	// when the destination already contains a git repository.
	// Captured on the Plan for the same reason as Mode: a single
	// source of truth shared by render + execute prevents the
	// dry-run preview from advertising a behavior different from
	// what the executor would actually do.
	OnExists string
	// Rows is the deduplicated, validated list of clones to perform.
	// Order matches the on-disk order so dry-run output is stable
	// across runs of the same file.
	Rows []Row
}

// Row is one scan-record-shaped clone target. Unlike clonefrom.Row
// (which carries raw URL + optional branch/depth) this carries the
// full pair of HTTPS / SSH URLs + the recorded relative folder so
// the executor can pick a URL based on the mode flag without re-
// parsing or re-deriving anything.
type Row struct {
	// RepoName is shown in progress + summary lines. Derived from
	// the scan record's RepoName when present; falls back to the
	// last URL segment for text-format inputs.
	RepoName string
	// HTTPSUrl is the recorded https-style clone URL (may be empty
	// for text-format rows when the source line was an ssh URL).
	HTTPSUrl string
	// SSHUrl is the recorded ssh-style clone URL (may be empty for
	// text-format rows when the source line was an https URL).
	SSHUrl string
	// Branch optionally pins the initial branch with --branch.
	// Empty -> git uses the remote's HEAD. Honored for JSON / CSV
	// rows; text-format rows always come through with empty Branch
	// because the plain `git clone <url> <dir>` line carries no
	// branch information.
	Branch string
	// RelativePath is the destination directory relative to cwd at
	// execute time. Always non-empty after parsing -- ParseFile
	// fills it in from the URL basename when the source row didn't
	// supply one (text-format with no explicit destination arg).
	RelativePath string
}

// PickURL returns the URL to clone with for the given mode, falling
// back to the other mode if the preferred URL is missing on this row.
// Returns "" only when both URL fields are empty -- a row that
// reaches Execute with empty PickURL is reported as a failure with
// MsgCloneNowNoURL rather than silently skipped.
//
// Centralized here (rather than duplicated in execute / render) so
// the dry-run preview and the actual git invocation always agree on
// which URL would be used.
func (r Row) PickURL(mode string) string {
	if mode == "ssh" {
		if len(r.SSHUrl) > 0 {
			return r.SSHUrl
		}

		return r.HTTPSUrl
	}
	if len(r.HTTPSUrl) > 0 {
		return r.HTTPSUrl
	}

	return r.SSHUrl
}
