package constants

// commit-in tag-replay constants. See spec/03-commit-in/09-commit-in-replay-map.md
// and spec/03-commit-in/08-tag-mirroring-and-release-branches.md.
//
// Per Core memory: no magic strings. The five `TagReplayOutcome`
// member names are centralized here and referenced from SQL seeds,
// store helpers, runlog inserters, and acceptance tests. Adding a
// new outcome requires touching exactly one Go file and one SQL seed.

// ---- Table-name constants -----------------------------------------
const (
	TableCommitInReplayMap  = "CommitInReplayMap"
	TableCommitInTagOutcome = "TagReplayOutcome"
)

// ---- Enum: TagReplayOutcome (spec §9.4) ---------------------------
const (
	TagReplayOutcomeCreated       = "Created"
	TagReplayOutcomeCreatedDryRun = "CreatedDryRun"
	TagReplayOutcomeSkipped       = "Skipped"
	TagReplayOutcomeAlreadyExists = "AlreadyExists"
	TagReplayOutcomeFailed        = "Failed"
)

// TagReplayOutcomeAll is the canonical, alphabetised list of every
// outcome name. The migration seed and the parity test in
// gitmap/store/migrate_commitin_test.go MUST agree with this slice.
var TagReplayOutcomeAll = []string{
	TagReplayOutcomeAlreadyExists,
	TagReplayOutcomeCreated,
	TagReplayOutcomeCreatedDryRun,
	TagReplayOutcomeFailed,
	TagReplayOutcomeSkipped,
}

// VersionTagPattern is the canonical regex source matching annotated
// "version" tags per spec §08 / §09.4 (`IsVersionTag`). Accepts an
// optional leading `v` and a SemVer 2.0.0 core (MAJOR.MINOR.PATCH)
// with an optional pre-release / build-metadata suffix. Anchored.
//
// Examples that match:   v1.2.3, 1.2.3, v1.0.0-rc.1, v2.0.0+build.7
// Examples that DON'T:   v1.2, nightly, release-1.0, 1, v1.2.3.4
//
// Centralized so the detector, the §08 release-branch step, and any
// future tooling all agree on what a "version tag" is. Compiled lazily
// by callers via regexp.MustCompile so this stays a pure string const.
const VersionTagPattern = `^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)` +
	`(?:-((?:0|[1-9]\d*|\d*[A-Za-z-][0-9A-Za-z-]*)(?:\.(?:0|[1-9]\d*|\d*[A-Za-z-][0-9A-Za-z-]*))*))?` +
	`(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`
