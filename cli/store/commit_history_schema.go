package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const sqlCreateCommitHistoryTables = `
CREATE TABLE IF NOT EXISTS CommitHistoryRecord (
    CommitHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommitSha TEXT NOT NULL,
    ShortSha TEXT NOT NULL,
    RepoName TEXT NOT NULL DEFAULT '',
    RepoSlug TEXT NOT NULL DEFAULT '',
    RepoRoot TEXT NOT NULL DEFAULT '',
    BranchName TEXT NOT NULL DEFAULT '',
    AuthorName TEXT NOT NULL DEFAULT '',
    AuthorEmail TEXT NOT NULL DEFAULT '',
    CommitMessage TEXT NOT NULL DEFAULT '',
    CommittedAtUnix INTEGER NOT NULL DEFAULT 0,
    AiModelName TEXT NOT NULL DEFAULT '',
    AiToolName TEXT NOT NULL DEFAULT '',
    AgentId TEXT NOT NULL DEFAULT '',
    HeatmapSessionId TEXT NOT NULL DEFAULT '',
    HeatmapActionId INTEGER NOT NULL DEFAULT 0,
    TaskSlug TEXT NOT NULL DEFAULT '',
    SpecPath TEXT NOT NULL DEFAULT '',
    PlansCount INTEGER NOT NULL DEFAULT 0,
    SubtasksCount INTEGER NOT NULL DEFAULT 0,
    IsAiCommitted INTEGER NOT NULL DEFAULT 0,
    IsCurrentMachineAi INTEGER NOT NULL DEFAULT 0,
    Notes TEXT NOT NULL DEFAULT '',
    FilesChanged INTEGER NOT NULL DEFAULT 0,
    LinesAdded INTEGER NOT NULL DEFAULT 0,
    LinesDeleted INTEGER NOT NULL DEFAULT 0,
    HasPushed INTEGER NOT NULL DEFAULT 0,
    CreatedAtUnix INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT Uq_CommitHistory_Sha UNIQUE (CommitSha)
);

CREATE TABLE IF NOT EXISTS CommitFileTouch (
    CommitFileTouchId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommitSha TEXT NOT NULL,
    FilePath TEXT NOT NULL DEFAULT '',
    ChangeType TEXT NOT NULL DEFAULT '',
    LinesAdded INTEGER NOT NULL DEFAULT 0,
    LinesDeleted INTEGER NOT NULL DEFAULT 0,
    CreatedAtUnix INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (CommitSha) REFERENCES CommitHistoryRecord(CommitSha) ON DELETE CASCADE
);
`

const sqlCreateCommitHistoryIndexes = `
CREATE INDEX IF NOT EXISTS IdxCommitHistory_Covering 
    ON CommitHistoryRecord(RepoSlug, CommittedAtUnix DESC, CommitSha);

CREATE INDEX IF NOT EXISTS IdxCommitHistory_ShortSha 
    ON CommitHistoryRecord(ShortSha);

CREATE INDEX IF NOT EXISTS IdxCommitHistory_TaskSlug 
    ON CommitHistoryRecord(TaskSlug);

CREATE INDEX IF NOT EXISTS IdxCommitHistory_AgentId 
    ON CommitHistoryRecord(AgentId);

CREATE INDEX IF NOT EXISTS IdxCommitFileTouch_CommitSha 
    ON CommitFileTouch(CommitSha);

CREATE INDEX IF NOT EXISTS IdxCommitFileTouch_FilePath 
    ON CommitFileTouch(FilePath);
`

func initCommitHistorySchema(conn *sql.DB) error {
	if _, err := conn.Exec(sqlCreateCommitHistoryTables); err != nil {
		return apperror.WrapSimple(err, "store.initCommitHistorySchema.tables")
	}

	if _, err := conn.Exec(sqlCreateCommitHistoryIndexes); err != nil {
		return apperror.WrapSimple(err, "store.initCommitHistorySchema.indexes")
	}

	return nil
}
