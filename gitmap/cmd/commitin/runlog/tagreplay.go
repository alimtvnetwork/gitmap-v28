package runlog

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TagReplayFacts is the persistence-layer projection of one mirrored
// annotated tag. Mirrors spec §9.4 columns. Caller resolves the
// `RewrittenCommit` row id BEFORE invoking RecordTagReplay so this
// stays a pure write helper (no FK guessing in the runlog layer).
//
// `Outcome` MUST be a constants.TagReplayOutcome* literal — the enum
// FK is resolved here via the standard mirror-table lookup pattern.
//
// `IsAnnotated` records whether the SOURCE tag is annotated (true) or
// lightweight (false). Lightweight tags can never be classified as
// "version tags" regardless of their name — see ClassifyVersionTag
// and the recorder gate below. This field is the single source of
// truth for the strict-semver contract: the mapping layer refuses to
// promote a lightweight `v1.2.3` to `IsVersionTag=1`.
type TagReplayFacts struct {
	SourceTagName         string
	SourceTagSha          string
	SourceCommitSha       string
	DestTagSha            string // empty for DryRun / Failed / Skipped
	DestCommitSha         string // empty for DryRun
	MirroredReleaseBranch string // empty when no branch was mirrored
	IsAnnotated           bool
	IsVersionTag          bool
	Outcome               string
}

// TagReplayLookup is the cross-run idempotency view returned by
// LookupTagReplay (spec §9.5). Every field carries the destination
// state recorded by the previous successful mirror; an empty
// DestCommitSha means the prior row was Created/AlreadyExists under
// a flow that did NOT capture the dest commit (defensive — should
// not occur for those two outcomes per §9.4).
type TagReplayLookup struct {
	DestTagSha            string
	DestCommitSha         string
	MirroredReleaseBranch string
}

// ErrTagReplayMiss is returned by LookupTagReplay when no row matches
// the (SourceTagName, SourceTagSha) pair under a Created /
// AlreadyExists outcome. Callers MUST use errors.Is to detect this
// (Core memory: zero-swallow, errors.Is everywhere).
var ErrTagReplayMiss = errors.New("runlog: tag replay lookup miss")

// ErrLightweightVersionTag is returned by RecordTagReplay when the
// caller tries to persist a row with `IsVersionTag=true` for a
// lightweight (non-annotated) tag. Strict-semver contract: only
// annotated tags whose name matches `constants.VersionTagPattern` may
// be classified as version tags — see ClassifyVersionTag for the
// canonical decision. Callers MUST use errors.Is to detect this
// (Core memory: zero-swallow, errors.Is everywhere).
var ErrLightweightVersionTag = errors.New("runlog: lightweight tag cannot be a version tag")

// RecordTagReplay persists one CommitInReplayMap row. Empty string
// fields are stored as SQL NULL where the column is nullable per
// spec §9.4 (DestTagSha, DestCommitSha, MirroredReleaseBranch).
//
// Strict-semver gate: if the caller asserts `IsVersionTag=true` for a
// lightweight tag (`IsAnnotated=false`), the insert is rejected with
// ErrLightweightVersionTag. This is the SINGLE choke point that
// guarantees `CommitInReplayMap.IsVersionTag=1` rows are always
// annotated semver tags — no upstream filter can quietly drift.
func RecordTagReplay(db *sql.DB, runID, rewrittenID int64, f TagReplayFacts) (int64, error) {
	if f.IsVersionTag && !f.IsAnnotated {
		return 0, fmt.Errorf("runlog: tag %q: %w", f.SourceTagName, ErrLightweightVersionTag)
	}
	outcomeID, err := lookupEnumID(db, constants.TableCommitInTagOutcome,
		"TagReplayOutcomeId", f.Outcome)
	if err != nil {
		return 0, fmt.Errorf("runlog: lookup tag outcome %q: %w", f.Outcome, err)
	}
	res, err := db.Exec(constants.SQLInsertCommitInReplayMap,
		runID, rewrittenID,
		f.SourceTagName, f.SourceTagSha, f.SourceCommitSha,
		nullIfEmpty(f.DestTagSha), nullIfEmpty(f.DestCommitSha),
		nullIfEmpty(f.MirroredReleaseBranch),
		boolToInt(f.IsVersionTag), outcomeID,
	)
	if err != nil {
		return 0, fmt.Errorf("runlog: insert CommitInReplayMap: %w", err)
	}
	return res.LastInsertId()
}

// LookupTagReplay implements the §9.5 cross-run short-circuit query.
// Returns ErrTagReplayMiss when no prior Created / AlreadyExists row
// matches. Pure read; safe to call before §08 attempts `git tag`.
func LookupTagReplay(db *sql.DB, sourceTagName, sourceTagSha string) (TagReplayLookup, error) {
	var got TagReplayLookup
	var dt, dc, mb sql.NullString
	err := db.QueryRow(constants.SQLSelectCommitInReplayLookup,
		sourceTagName, sourceTagSha).Scan(&dt, &dc, &mb)
	if errors.Is(err, sql.ErrNoRows) {
		return got, ErrTagReplayMiss
	}
	if err != nil {
		return got, fmt.Errorf("runlog: lookup tag replay %q: %w", sourceTagName, err)
	}
	got.DestTagSha = dt.String
	got.DestCommitSha = dc.String
	got.MirroredReleaseBranch = mb.String
	return got, nil
}

// IsAnnotatedSemverVersionTag returns true iff `tagName` matches the
// canonical version-tag pattern (spec §08, §9.4 `IsVersionTag`). The
// regex compile is cached once via sync.Once because the pattern is a
// package constant and never mutates.
//
// NAME-ONLY: this helper performs the regex match in isolation. It is
// the building block for ClassifyVersionTag, which is the canonical
// strict gate (annotated AND semver). Direct callers exist only at
// the §08 walker boundary, where the annotated-vs-lightweight bit is
// derived from the git object kind. Persistence-layer callers MUST
// use ClassifyVersionTag instead.
func IsAnnotatedSemverVersionTag(tagName string) bool {
	return versionTagRegex().MatchString(tagName)
}

// ClassifyVersionTag is the canonical strict-semver classifier used by
// the mapping layer. A tag is a "version tag" iff:
//
//  1. The tag object is annotated (`isAnnotated=true`), AND
//  2. The tag name matches `constants.VersionTagPattern`.
//
// Lightweight tags are NEVER version tags, even when their name looks
// like `v1.2.3`. This matches the spec §08 §3 contract that release
// branches and the §09 `IsVersionTag` flag both reflect deliberate,
// signed-or-annotated release intent — not an accidental ref name.
func ClassifyVersionTag(tagName string, isAnnotated bool) bool {
	if !isAnnotated {
		return false
	}
	return IsAnnotatedSemverVersionTag(tagName)
}

var (
	versionTagRegexOnce sync.Once
	versionTagRegexVal  *regexp.Regexp
)

func versionTagRegex() *regexp.Regexp {
	versionTagRegexOnce.Do(func() {
		versionTagRegexVal = regexp.MustCompile(constants.VersionTagPattern)
	})
	return versionTagRegexVal
}

// nullIfEmpty maps "" → SQL NULL so spec §9.4's NULL-on-dry-run /
// NULL-on-failed contract is honored without the caller juggling
// `any` typed parameters.
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
