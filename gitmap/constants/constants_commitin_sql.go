package constants

// commit-in SQL DDL — see spec/03-commit-in/04-database-schema.md.
//
// Table-name conventions (Core memory rules):
//   - PascalCase, singular.
//   - Every PK is `<TableName>Id INTEGER PRIMARY KEY AUTOINCREMENT`.
//   - Every classifier FK is `<EnumName>Id INTEGER NOT NULL` with a
//     mirror table `(<EnumName>Id, Name TEXT NOT NULL UNIQUE)`.
//   - Every table has `CreatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`.
//   - All FKs are `ON DELETE RESTRICT` (purge requires explicit
//     `commit-in --reset`, out of scope v1).
//
// All seed inserts use `INSERT OR IGNORE` so re-running Migrate() is a
// no-op on existing rows.

// ---- Table-name constants -----------------------------------------
const (
	TableCommitInRun        = "CommitInRun"
	TableCommitInRunStatus  = "RunStatus"
	TableCommitInInputRepo  = "InputRepo"
	TableCommitInInputKind  = "InputKind"
	TableCommitInSrcCommit  = "SourceCommit"
	TableCommitInSrcFile    = "SourceCommitFile"
	TableCommitInRewritten  = "RewrittenCommit"
	TableCommitInOutcome    = "CommitOutcome"
	TableCommitInSkipLog    = "SkipLog"
	TableCommitInSkipReason = "SkipReason"
	TableCommitInShaMap     = "ShaMap"
	TableCommitInProfile    = "Profile"
	TableCommitInProfileExc = "ProfileExclusion"
	TableCommitInExcKind    = "ExclusionKind"
	TableCommitInProfileMsg = "ProfileMessageRule"
	TableCommitInMsgKind    = "MessageRuleKind"
	TableCommitInLanguage   = "FunctionIntelLanguage"
	TableCommitInConflict   = "ConflictMode"
)

// ---- Enum-mirror tables -------------------------------------------
const SQLCreateCommitInRunStatus = `CREATE TABLE IF NOT EXISTS RunStatus (
	RunStatusId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name        TEXT    NOT NULL UNIQUE,
	CreatedAt   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInInputKind = `CREATE TABLE IF NOT EXISTS InputKind (
	InputKindId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name        TEXT    NOT NULL UNIQUE,
	CreatedAt   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInOutcome = `CREATE TABLE IF NOT EXISTS CommitOutcome (
	CommitOutcomeId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name            TEXT    NOT NULL UNIQUE,
	CreatedAt       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInSkipReason = `CREATE TABLE IF NOT EXISTS SkipReason (
	SkipReasonId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name         TEXT    NOT NULL UNIQUE,
	CreatedAt    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInExclusionKind = `CREATE TABLE IF NOT EXISTS ExclusionKind (
	ExclusionKindId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name            TEXT    NOT NULL UNIQUE,
	CreatedAt       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInMessageRuleKind = `CREATE TABLE IF NOT EXISTS MessageRuleKind (
	MessageRuleKindId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name              TEXT    NOT NULL UNIQUE,
	CreatedAt         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInLanguage = `CREATE TABLE IF NOT EXISTS FunctionIntelLanguage (
	FunctionIntelLanguageId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name                    TEXT    NOT NULL UNIQUE,
	CreatedAt               DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInConflictMode = `CREATE TABLE IF NOT EXISTS ConflictMode (
	ConflictModeId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name           TEXT    NOT NULL UNIQUE,
	CreatedAt      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

// ---- Profile (defined before CommitInRun so the FK can resolve) ---
const SQLCreateCommitInProfile = `CREATE TABLE IF NOT EXISTS Profile (
	ProfileId      INTEGER PRIMARY KEY AUTOINCREMENT,
	Name           TEXT    NOT NULL UNIQUE,
	SourceRepoPath TEXT    NOT NULL,
	IsDefault      INTEGER NOT NULL DEFAULT 0,
	PayloadJson    TEXT    NOT NULL,
	CreatedAt      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

// One default profile per source repo (partial unique index).
const SQLCreateCommitInProfileDefaultIdx = `CREATE UNIQUE INDEX IF NOT EXISTS
	UX_Profile_DefaultPerSource ON Profile (SourceRepoPath) WHERE IsDefault = 1`

const SQLCreateCommitInProfileExclusion = `CREATE TABLE IF NOT EXISTS ProfileExclusion (
	ProfileExclusionId INTEGER PRIMARY KEY AUTOINCREMENT,
	ProfileId          INTEGER NOT NULL REFERENCES Profile(ProfileId) ON DELETE RESTRICT,
	Value              TEXT    NOT NULL,
	ExclusionKindId    INTEGER NOT NULL REFERENCES ExclusionKind(ExclusionKindId) ON DELETE RESTRICT,
	CreatedAt          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInProfileMessageRule = `CREATE TABLE IF NOT EXISTS ProfileMessageRule (
	ProfileMessageRuleId INTEGER PRIMARY KEY AUTOINCREMENT,
	ProfileId            INTEGER NOT NULL REFERENCES Profile(ProfileId) ON DELETE RESTRICT,
	MessageRuleKindId    INTEGER NOT NULL REFERENCES MessageRuleKind(MessageRuleKindId) ON DELETE RESTRICT,
	Value                TEXT    NOT NULL,
	CreatedAt            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

// ---- Core run/commit tables ---------------------------------------
const SQLCreateCommitInRun = `CREATE TABLE IF NOT EXISTS CommitInRun (
	CommitInRunId       INTEGER PRIMARY KEY AUTOINCREMENT,
	SourceRepoPath      TEXT    NOT NULL,
	SourceRepoUrl       TEXT,
	WasSourceFreshlyInit INTEGER NOT NULL DEFAULT 0,
	StartedAt           DATETIME NOT NULL,
	FinishedAt          DATETIME,
	RunStatusId         INTEGER NOT NULL REFERENCES RunStatus(RunStatusId) ON DELETE RESTRICT,
	ProfileId           INTEGER REFERENCES Profile(ProfileId) ON DELETE RESTRICT,
	CreatedAt           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInRunSourceIdx = `CREATE INDEX IF NOT EXISTS
	IX_CommitInRun_SourceRepoPath ON CommitInRun (SourceRepoPath)`

const SQLCreateCommitInInputRepo = `CREATE TABLE IF NOT EXISTS InputRepo (
	InputRepoId   INTEGER PRIMARY KEY AUTOINCREMENT,
	CommitInRunId INTEGER NOT NULL REFERENCES CommitInRun(CommitInRunId) ON DELETE RESTRICT,
	OrderIndex    INTEGER NOT NULL,
	OriginalRef   TEXT    NOT NULL,
	ResolvedPath  TEXT    NOT NULL,
	InputKindId   INTEGER NOT NULL REFERENCES InputKind(InputKindId) ON DELETE RESTRICT,
	CreatedAt     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (CommitInRunId, OrderIndex)
)`

const SQLCreateCommitInSourceCommit = `CREATE TABLE IF NOT EXISTS SourceCommit (
	SourceCommitId  INTEGER PRIMARY KEY AUTOINCREMENT,
	InputRepoId     INTEGER NOT NULL REFERENCES InputRepo(InputRepoId) ON DELETE RESTRICT,
	SourceSha       TEXT    NOT NULL,
	AuthorName      TEXT    NOT NULL,
	AuthorEmail     TEXT    NOT NULL,
	AuthorDate      DATETIME NOT NULL,
	CommitterDate   DATETIME NOT NULL,
	OriginalMessage TEXT    NOT NULL,
	OrderIndex      INTEGER NOT NULL,
	CreatedAt       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (InputRepoId, SourceSha)
)`

const SQLCreateCommitInSourceCommitShaIdx = `CREATE INDEX IF NOT EXISTS
	IX_SourceCommit_SourceSha ON SourceCommit (SourceSha)`

const SQLCreateCommitInSourceCommitFile = `CREATE TABLE IF NOT EXISTS SourceCommitFile (
	SourceCommitFileId INTEGER PRIMARY KEY AUTOINCREMENT,
	SourceCommitId     INTEGER NOT NULL REFERENCES SourceCommit(SourceCommitId) ON DELETE RESTRICT,
	RelativePath       TEXT    NOT NULL,
	CreatedAt          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (SourceCommitId, RelativePath)
)`

const SQLCreateCommitInRewritten = `CREATE TABLE IF NOT EXISTS RewrittenCommit (
	RewrittenCommitId    INTEGER PRIMARY KEY AUTOINCREMENT,
	CommitInRunId        INTEGER NOT NULL REFERENCES CommitInRun(CommitInRunId) ON DELETE RESTRICT,
	SourceCommitId       INTEGER NOT NULL REFERENCES SourceCommit(SourceCommitId) ON DELETE RESTRICT,
	NewSha               TEXT,
	FinalMessage         TEXT    NOT NULL,
	AppliedAuthorName    TEXT    NOT NULL,
	AppliedAuthorEmail   TEXT    NOT NULL,
	AppliedAuthorDate    DATETIME NOT NULL,
	AppliedCommitterDate DATETIME NOT NULL,
	CommitOutcomeId      INTEGER NOT NULL REFERENCES CommitOutcome(CommitOutcomeId) ON DELETE RESTRICT,
	CreatedAt            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (CommitInRunId, SourceCommitId)
)`

const SQLCreateCommitInRewrittenShaIdx = `CREATE INDEX IF NOT EXISTS
	IX_RewrittenCommit_NewSha ON RewrittenCommit (NewSha)`

const SQLCreateCommitInSkipLog = `CREATE TABLE IF NOT EXISTS SkipLog (
	SkipLogId                 INTEGER PRIMARY KEY AUTOINCREMENT,
	CommitInRunId             INTEGER NOT NULL REFERENCES CommitInRun(CommitInRunId) ON DELETE RESTRICT,
	SourceCommitId            INTEGER NOT NULL REFERENCES SourceCommit(SourceCommitId) ON DELETE RESTRICT,
	SkipReasonId              INTEGER NOT NULL REFERENCES SkipReason(SkipReasonId) ON DELETE RESTRICT,
	PreviousRewrittenCommitId INTEGER REFERENCES RewrittenCommit(RewrittenCommitId) ON DELETE RESTRICT,
	CreatedAt                 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInShaMap = `CREATE TABLE IF NOT EXISTS ShaMap (
	ShaMapId          INTEGER PRIMARY KEY AUTOINCREMENT,
	SourceSha         TEXT    NOT NULL UNIQUE,
	RewrittenCommitId INTEGER NOT NULL REFERENCES RewrittenCommit(RewrittenCommitId) ON DELETE RESTRICT,
	CreatedAt         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const SQLCreateCommitInShaMapIdx = `CREATE INDEX IF NOT EXISTS
	IX_ShaMap_SourceSha ON ShaMap (SourceSha)`

// ---- Enum-mirror seed inserts -------------------------------------
// Seed Names MUST match the Go enum String() exactly (parity test in
// gitmap/cmd/commitin/enums_test.go locks the Go side; another test
// in gitmap/store/migrate_commitin_test.go locks the SQL side).
const (
	SQLSeedCommitInRunStatus = `INSERT OR IGNORE INTO RunStatus (Name) VALUES
		('Pending'), ('Running'), ('Completed'), ('Failed'), ('PartiallyFailed')`

	SQLSeedCommitInInputKind = `INSERT OR IGNORE INTO InputKind (Name) VALUES
		('LocalFolder'), ('GitUrl'), ('VersionedSibling')`

	SQLSeedCommitInOutcome = `INSERT OR IGNORE INTO CommitOutcome (Name) VALUES
		('Created'), ('Skipped'), ('Failed')`

	SQLSeedCommitInSkipReason = `INSERT OR IGNORE INTO SkipReason (Name) VALUES
		('DuplicateSourceSha'), ('ExcludedAllFiles'),
		('EmptyAfterMessageRules'), ('DryRun')`

	SQLSeedCommitInExclusionKind = `INSERT OR IGNORE INTO ExclusionKind (Name) VALUES
		('PathFolder'), ('PathFile')`

	SQLSeedCommitInMessageRuleKind = `INSERT OR IGNORE INTO MessageRuleKind (Name) VALUES
		('StartsWith'), ('EndsWith'), ('Contains')`

	SQLSeedCommitInLanguage = `INSERT OR IGNORE INTO FunctionIntelLanguage (Name) VALUES
		('Go'), ('JavaScript'), ('TypeScript'), ('Rust'),
		('Python'), ('Php'), ('Java'), ('CSharp')`

	SQLSeedCommitInConflictMode = `INSERT OR IGNORE INTO ConflictMode (Name) VALUES
		('ForceMerge'), ('Prompt')`
)
