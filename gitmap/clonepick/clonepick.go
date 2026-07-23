// Package clonepick implements `gitmap clone-pick <repo-url> <paths>`:
// a partial / sparse-checkout clone of a single git repository,
// fetching only the requested repo-relative paths into the current
// working directory (or --dest).
//
// Why a separate package (not clonefrom / clonenow / cloner)?
//
//   - clonefrom is plan-driven from a manifest (rows of url+dest+branch).
//   - clonenow is round-trip from `gitmap scan` output.
//   - cloner is the in-memory pipeline cloner used during scan.
//
// clone-pick sits orthogonal to all three: input is one URL + a path
// list, output is a sparse-checkout. Sharing code with the others
// would force them to grow path-filter awareness they don't need.
//
// Persistence: every successful run writes one row to the
// CloneInteractiveSelection table so the same selection can be
// re-applied later via --replay <id|name>.
//
// See spec/01-app/100-clone-pick.md.
package clonepick

// Plan is the validated, in-memory representation of one clone-pick
// invocation. Built by ParseArgs (or LoadFromDB for --replay) and
// consumed by Render (dry-run) and Execute.
//
// Field order mirrors the CloneInteractiveSelection schema so the
// persist layer can build/destructure rows without a translation map.
type Plan struct {
	// Name is an optional human label persisted with the row. Empty
	// is fine -- the row still gets a SelectionId.
	Name string
	// RepoCanonicalId is the host/owner/repo form returned by
	// gitutil.CanonicalRepoID(RepoUrl). Stored separately so
	// --replay can match HTTPS↔SSH variants of the same repo.
	RepoCanonicalId string
	// RepoUrl is the canonical URL passed to `git clone`. Either the
	// user-supplied URL verbatim (when it was full) or the expanded
	// form built from the owner/repo shorthand + Mode.
	RepoUrl string
	// Mode is "https" | "ssh" -- only meaningful when the input was
	// shorthand. Persisted so --replay reproduces the same URL form.
	Mode string
	// Branch is the optional --branch argument forwarded to git
	// clone. Empty -> git uses remote HEAD.
	Branch string
	// Depth is the shallow-clone depth. 0 = full history, 1 = the
	// classic shallow checkout used by most pickers.
	Depth int
	// Cone is true when sparse-checkout cone mode applies (folder-
	// only patterns). Auto-flipped to false by ParseArgs when any
	// path looks like a glob or carries a file extension after a /.
	Cone bool
	// KeepGit controls whether .git is preserved after checkout.
	// false = files-only mode (rm -rf .git after sparse-checkout).
	KeepGit bool
	// DestDir is the destination directory relative to cwd. "."
	// means clone into the current directory (no shell handoff).
	DestDir string
	// Paths is the deduplicated, sorted list of repo-relative paths
	// to materialize. Already validated by ParseArgs (no empty,
	// no absolute, no traversal, no oversized entries).
	Paths []string
	// UsedAsk records whether --ask was passed. Persisted so
	// --replay can decide whether to re-launch the picker (for now
	// it doesn't, but keeping the bit avoids a schema migration
	// later).
	UsedAsk bool
	// DryRun, Quiet, Force are runtime-only flags (not persisted).
	// Captured on the Plan so render + execute share one source of
	// truth and don't need to thread the cmd-side cfg struct down.
	DryRun bool
	Quiet  bool
	Force  bool
	// PreClonedSrc, when non-empty, is a temp directory holding a
	// metadata-only clone (filter=blob:none, no-checkout) that the
	// executor will promote into DestDir instead of running a fresh
	// `git clone`. Set by the --ask picker so we don't pay the
	// clone cost twice. Always absolute. The executor takes
	// ownership: on success the dir is moved away, on failure the
	// caller is responsible for cleanup.
	PreClonedSrc string
}

// Result describes what Execute did with one Plan. Returned by
// Execute so the cmd layer can pick an exit code without
// re-implementing success/failure heuristics.
type Result struct {
	// Status is one of "ok" | "failed" | "canceled" (the wire value
	// kept on the wire for backwards compat with existing JSON
	// consumers). Numeric exit code mapping lives in the cmd layer
	// (cmd/clonepick.go) so the constants stay co-located with the
	// user-facing strings.
	Status string
	// SelectionId is the DB row id assigned by the persist layer.
	// 0 when DryRun was set (no row written) or when the run failed
	// before reaching the persist step.
	SelectionId int64
	// Detail is a single-line user-facing message that's already
	// printed to stderr by Execute -- duplicated on the Result so
	// callers can include it in JSON summaries without re-deriving.
	Detail string
}

// Status enum values mirror clonenow's vocabulary so downstream
// pipelines that grep "ok"/"failed" keep working.
const (
	StatusOK        = "ok"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled" //nolint:misspell // stable JSON status value, do not rename.
)
