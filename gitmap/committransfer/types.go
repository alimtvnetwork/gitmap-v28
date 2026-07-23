// Package committransfer implements the commit-replay family:
// `gitmap commit-left`, `commit-right`, `commit-both`. Each command
// resolves two endpoints (LEFT, RIGHT) and replays one side's commit
// timeline onto the other as a sequence of fresh, cleaned commits via
// the manual-reconstruct mechanism (checkout + file-snapshot + commit).
//
// Spec: spec/01-app/106-commit-left-right-both.md
//
// **Status (v3.102.0):** all three directions live —
// `commit-right` (Phase 1, v3.76.0), `commit-left` (Phase 2), and
// `commit-both` (Phase 3) all run on the same Plan/Replay primitives.
package committransfer

import "time"

// Direction names which side receives the new commits.
type Direction int

const (
	// DirRight replays LEFT → RIGHT (writes commits on RIGHT).
	DirRight Direction = iota
	// DirLeft replays RIGHT → LEFT (writes commits on LEFT).
	DirLeft
	// DirBoth replays in both directions sequentially: LEFT → RIGHT
	// first, then RIGHT → LEFT. Each pass builds its own plan and
	// shares the same Options (with a directional log-prefix suffix).
	DirBoth
)

// MessagePolicy controls the per-commit message normalization pipeline.
// All fields default to the "conservative" values described in spec §6
// when zero-initialized; the CLI layer overrides via flags.
type MessagePolicy struct {
	// DropPatterns are regexes matched against the source subject.
	// A match means the whole commit is skipped (not replayed).
	DropPatterns []string

	// StripPatterns are regexes matched against the source subject and
	// replaced with the empty string. Run before Conventional.
	StripPatterns []string

	// Conventional, when true, normalizes the cleaned subject into a
	// `<type>: <subject>` form (spec §6.3).
	Conventional bool

	// Provenance, when true, appends a `gitmap-replay:` footer block
	// (spec §6.4) so re-runs can detect already-replayed commits.
	Provenance bool

	// SourceDisplayName is what the provenance footer records as the
	// source repo identifier. Caller sets it from Endpoint.DisplayName.
	SourceDisplayName string

	// CommandName is the user-visible command label embedded in the
	// `gitmap-replay-cmd:` footer line. Set by the dispatcher
	// (e.g. "commit-right").
	CommandName string
}

// PreferPolicy mirrors movemerge.PreferPolicy for file-level conflicts
// during the snapshot-copy step. We re-export the values rather than
// importing the type to keep the package decoupled in tests.
type PreferPolicy int

const (
	// PreferNone means use the existing file (no overwrite on conflict).
	PreferNone PreferPolicy = iota
	// PreferSource overwrites target with source on every conflict.
	PreferSource
	// PreferTarget keeps target on every conflict (effectively skip).
	PreferTarget
)

// Options bundles every CLI flag for the commit-transfer family.
// Mirrors movemerge.Options field-for-field where the semantics overlap.
//
// v6.0.0 breaking change: IncludeMerges defaults to true in the CLI layer.
// Direct API users who create Options{} get the zero value (false), which
// preserves legacy strip behaviour; set IncludeMerges = true explicitly
// for the new default.
type Options struct {
	Yes            bool          // skip the confirm prompt
	DryRun         bool          // print the plan; no writes
	NoPush         bool          // skip the final git push
	NoCommit       bool          // copy + stage but do not commit
	IncludeMerges  bool          // default true in CLI; zero value = legacy strip (v6.0.0)
	IncludeVCS     bool          // copy .git/* during snapshot
	IncludeNodeMod bool          // copy node_modules/* during snapshot
	Mirror         bool          // delete target-only files (true mirror)
	ForceReplay    bool          // replay even commits with provenance footer
	Interleave     bool          // commit-both only: author-date interleave
	Limit          int           // 0 = no limit; replay at most N (oldest first)
	Since          string        // override divergence base (sha or date)
	MaxHistoryScan int           // 0 = unbounded target log scan (spec 114 Gap A); >0 caps the idempotence-check log query at N commits
	Prefer         PreferPolicy  // file-conflict policy
	Message        MessagePolicy // §6 pipeline knobs
	CommandName    string        // "commit-right" etc.
	LogPrefix      string        // "[commit-right]" etc.
}

// SourceCommit is one entry in the resolved replay set.
type SourceCommit struct {
	SHA       string    // full source SHA (rev-parse)
	ShortSHA  string    // %h, used in logs and the provenance footer
	Subject   string    // raw source subject line
	Body      string    // raw source body (everything after subject + blank line)
	Author    string    // "Name <email>"
	AuthorAt  time.Time // author date (used by commit-both interleave)
	Cleaned   string    // post-pipeline message; "" means skipped
	SkipCause string    // populated when the planner pre-skips this commit
}

// ReplayPlan is the fully-computed replay set for one direction. It is
// built once before any side effects so the user sees the full preview
// up front (spec §3).
type ReplayPlan struct {
	SourceDir   string         // resolved source working dir
	TargetDir   string         // resolved target working dir
	SourceHEAD  string         // original branch/sha at source (for restore)
	BaseSHA     string         // git merge-base output (or "" if unrelated)
	Commits     []SourceCommit // oldest-first
	SkippedDrop int            // commits dropped by the drop filter

	// MergeExcluded is the number of merge commits stripped from the
	// source range because Options.IncludeMerges was false. Surfaced
	// in PrintPlan / PrintSummary so the user can reconcile the plan
	// count against `git log` (which includes merges by default).
	MergeExcluded int

	// IncludeMerges mirrors the option that produced this plan so
	// PrintPlan can emit the correct notice (v6.0.0+).
	IncludeMerges bool
}

// ReplayResult is the outcome after the replay loop runs.
type ReplayResult struct {
	Replayed        int
	SkippedDrop     int
	SkippedReplayed int
	SkippedEmpty    int
	NewSHAs         []string // target-side SHAs in commit order
	Pushed          bool
}
