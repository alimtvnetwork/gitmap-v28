// Package clonefrom implements the `gitmap clone-from <file>` workflow:
// read a JSON or CSV plan from disk, validate every row, render a
// dry-run preview by default, and on `--execute` shell out to
// `git clone` for each row with a per-row pass/fail summary.
//
// Why a separate package (not gitmap/cloner or gitmap/clonenext)?
//
//   - gitmap/cloner is the scan-driven cloner: it consumes records
//     emitted by the scan workflow and assumes a DB-backed repo
//     model. `clone-from` is plan-driven: input comes from a user-
//     provided file, no scan, no model.CloneRecord round-trip.
//   - gitmap/clonenext is the version-bumping cloner for existing
//     local repos (`vN+1` of an already-cloned repo). `clone-from`
//     clones brand-new URLs to user-chosen destinations.
//
// Splitting keeps each cloner's contract tight and lets clone-from
// evolve (e.g., adding `--depth`, `--single-branch`, parallel
// fan-out) without touching the more constrained scan/cn paths.
package clonefrom

// Plan is the validated, in-memory representation of one input file.
// Built by ParseFile from either JSON or CSV; consumed by Render
// (dry-run) and Execute. The Source field carries the on-disk path
// so the dry-run header can echo it back to the user without the
// caller having to thread the original argument through.
type Plan struct {
	// Source is the absolute or user-supplied path the plan was
	// read from. Echoed verbatim in the dry-run header.
	Source string
	// Format is "json" or "csv" — used by the dry-run header so the
	// user can confirm we parsed the file the way they expected.
	Format string
	// Rows is the deduplicated, validated list of clones to perform.
	// Order matches the on-disk order so dry-run output is stable
	// across runs of the same file.
	Rows []Row
}

// Row is one git-clone target. Every field except URL is optional;
// zero values translate to "use git's default" (HEAD branch, full
// history, dest derived from URL basename).
type Row struct {
	// URL is the clone source. Required. Validated for non-empty
	// and rough HTTPS/SSH/scp-style shape — we do NOT round-trip
	// through net/url.Parse because git accepts forms (scp-style
	// `user@host:path`) that net/url rejects.
	URL string
	// Dest is the target directory relative to cwd at execute
	// time. Empty → derived from the last URL segment (matches
	// `git clone <url>` default).
	Dest string
	// Branch optionally pins the initial branch with --branch.
	// Empty → git uses the remote's HEAD.
	Branch string
	// Depth optionally enables a shallow clone with --depth=N.
	// Zero → full history.
	Depth int
	// Checkout controls post-clone working-tree behavior. One of
	// "" / "auto" / "skip" / "force" (validated at parse time).
	// Empty string means "inherit the global default" (set via
	// --checkout on the CLI; defaults to "auto"). Per-row value
	// always wins over the global default. See
	// constants.FlagDescCloneFromCheckout for the semantics of each
	// mode.
	Checkout string
}
