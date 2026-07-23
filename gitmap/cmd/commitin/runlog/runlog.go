package runlog

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// StartRun inserts a new CommitInRun row in the `Running` state and
// returns its primary key. Maps to spec §3.1 stage 04 follow-up (the
// row is created right after MigrateDb so every later stage can carry
// the runId for FK references).
func StartRun(db *sql.DB, sourceRepoPath string, sourceURL *string, wasFreshlyInit bool, profileID *int64, startedAt time.Time) (int64, error) {
	statusID, err := lookupEnumID(db, constants.TableCommitInRunStatus, "RunStatusId", constants.CommitInRunStatusRunning)
	if err != nil {
		return 0, fmt.Errorf("runlog: lookup RunStatus: %w", err)
	}
	res, err := db.Exec(sqlInsertRun, sourceRepoPath, sourceURL, boolToInt(wasFreshlyInit), startedAt.Format(time.RFC3339), statusID, profileID)
	if err != nil {
		return 0, fmt.Errorf("runlog: insert CommitInRun: %w", err)
	}
	return res.LastInsertId()
}

// FinishRun stamps FinishedAt + final status. Status MUST be one of
// the constants.CommitInRunStatus* literals.
func FinishRun(db *sql.DB, runID int64, status string, finishedAt time.Time) error {
	statusID, err := lookupEnumID(db, constants.TableCommitInRunStatus, "RunStatusId", status)
	if err != nil {
		return fmt.Errorf("runlog: lookup status %q: %w", status, err)
	}
	if _, err := db.Exec(sqlUpdateRunFinish, finishedAt.Format(time.RFC3339), statusID, runID); err != nil {
		return fmt.Errorf("runlog: finish run %d: %w", runID, err)
	}
	return nil
}

// boolToInt converts Go bool → SQLite-friendly 0/1.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

const (
	sqlInsertRun = `INSERT INTO CommitInRun
		(SourceRepoPath, SourceRepoUrl, WasSourceFreshlyInit, StartedAt, RunStatusId, ProfileId)
		VALUES (?, ?, ?, ?, ?, ?)`

	sqlUpdateRunFinish = `UPDATE CommitInRun
		SET FinishedAt = ?, RunStatusId = ?
		WHERE CommitInRunId = ?`
)
