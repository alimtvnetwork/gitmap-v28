package constants

// commit-in tag-replay SQL DDL (migration 007). See
// spec/03-commit-in/09-commit-in-replay-map.md.
//
// Idempotent per Core memory: every statement uses CREATE / INDEX
// IF NOT EXISTS, and the enum-mirror seed uses INSERT OR IGNORE.

// ---- Enum-mirror table --------------------------------------------
const SQLCreateCommitInTagOutcome = `CREATE TABLE IF NOT EXISTS TagReplayOutcome (
	TagReplayOutcomeId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name               TEXT    NOT NULL UNIQUE,
	CreatedAt          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

// ---- CommitInReplayMap (one row per mirrored annotated tag) -------
const SQLCreateCommitInReplayMap = `CREATE TABLE IF NOT EXISTS CommitInReplayMap (
	CommitInReplayMapId    INTEGER PRIMARY KEY AUTOINCREMENT,
	CommitInRunId          INTEGER NOT NULL REFERENCES CommitInRun(CommitInRunId) ON DELETE RESTRICT,
	RewrittenCommitId      INTEGER NOT NULL REFERENCES RewrittenCommit(RewrittenCommitId) ON DELETE RESTRICT,
	SourceTagName          TEXT    NOT NULL,
	SourceTagSha           TEXT    NOT NULL,
	SourceCommitSha        TEXT    NOT NULL,
	DestTagSha             TEXT,
	DestCommitSha          TEXT,
	MirroredReleaseBranch  TEXT,
	IsVersionTag           INTEGER NOT NULL,
	TagReplayOutcomeId     INTEGER NOT NULL REFERENCES TagReplayOutcome(TagReplayOutcomeId) ON DELETE RESTRICT,
	CreatedAt              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (CommitInRunId, SourceTagName),
	UNIQUE (CommitInRunId, RewrittenCommitId, SourceTagName)
)`

const SQLCreateCommitInReplayMapTagNameIdx = `CREATE INDEX IF NOT EXISTS
	IX_CommitInReplayMap_SourceTagName ON CommitInReplayMap (SourceTagName)`

const SQLCreateCommitInReplayMapDestShaIdx = `CREATE INDEX IF NOT EXISTS
	IX_CommitInReplayMap_DestCommitSha ON CommitInReplayMap (DestCommitSha)`

const SQLCreateCommitInReplayMapBranchIdx = `CREATE INDEX IF NOT EXISTS
	IX_CommitInReplayMap_MirroredReleaseBranch ON CommitInReplayMap (MirroredReleaseBranch)`

// ---- Enum-mirror seed (matches constants.TagReplayOutcomeAll) -----
const SQLSeedCommitInTagOutcome = `INSERT OR IGNORE INTO TagReplayOutcome (Name) VALUES
	('AlreadyExists'), ('Created'), ('CreatedDryRun'), ('Failed'), ('Skipped')`

// ---- Cross-run lookup query (spec §9.5) ---------------------------
//
// Bound parameters: ?1 = SourceTagName, ?2 = SourceTagSha. Returns
// (DestTagSha, DestCommitSha, MirroredReleaseBranch) of the most
// recent successful row, or no rows when this tag has never been
// mirrored. The IN-list locks the lookup to outcomes that REPRESENT
// a real destination state (Created / AlreadyExists), per §9.5.
const SQLSelectCommitInReplayLookup = `SELECT
	m.DestTagSha, m.DestCommitSha, m.MirroredReleaseBranch
FROM   CommitInReplayMap m
JOIN   TagReplayOutcome  o ON o.TagReplayOutcomeId = m.TagReplayOutcomeId
WHERE  m.SourceTagName = ?
  AND  m.SourceTagSha  = ?
  AND  o.Name IN ('Created', 'AlreadyExists')
ORDER BY m.CommitInReplayMapId DESC
LIMIT  1`

// SQLInsertCommitInReplayMap inserts one row. Caller resolves the
// outcome FK via a normal enum-mirror lookup (same pattern as
// CommitOutcome / SkipReason).
const SQLInsertCommitInReplayMap = `INSERT INTO CommitInReplayMap (
	CommitInRunId, RewrittenCommitId,
	SourceTagName, SourceTagSha, SourceCommitSha,
	DestTagSha, DestCommitSha, MirroredReleaseBranch,
	IsVersionTag, TagReplayOutcomeId
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
