package runlog

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// RecordRewritten persists one RewrittenCommit row + (when outcome is
// `Created`) the matching ShaMap entry so the next run dedupes it.
// Maps to spec §3.1 stage 15 (`RecordResult`).
func RecordRewritten(db *sql.DB, runID, sourceCommitID int64, r RewrittenRow) (int64, error) {
	outcomeID, err := lookupEnumID(db, constants.TableCommitInOutcome, "CommitOutcomeId", r.Outcome)
	if err != nil {
		return 0, fmt.Errorf("runlog: lookup outcome %q: %w", r.Outcome, err)
	}
	rewrittenID, err := insertRewritten(db, runID, sourceCommitID, outcomeID, r)
	if err != nil {
		return 0, err
	}
	if r.Outcome == constants.CommitInOutcomeCreated && r.NewSha != "" {
		if err := insertShaMap(db, r.SourceSha, rewrittenID); err != nil {
			return rewrittenID, err
		}
	}
	return rewrittenID, nil
}

// RecordSkip persists one SkipLog entry. Reason MUST be a
// constants.CommitInSkipReason* literal.
func RecordSkip(db *sql.DB, runID, sourceCommitID int64, reason string, previousRewrittenID *int64) error {
	reasonID, err := lookupEnumID(db, constants.TableCommitInSkipReason, "SkipReasonId", reason)
	if err != nil {
		return fmt.Errorf("runlog: lookup skip reason %q: %w", reason, err)
	}
	if _, err := db.Exec(sqlInsertSkipLog, runID, sourceCommitID, reasonID, previousRewrittenID); err != nil {
		return fmt.Errorf("runlog: insert SkipLog: %w", err)
	}
	return nil
}

// insertRewritten is a small helper kept separate so RecordRewritten
// stays under the 15-line cap.
func insertRewritten(db *sql.DB, runID, sourceCommitID, outcomeID int64, r RewrittenRow) (int64, error) {
	var newSha any
	if r.NewSha != "" {
		newSha = r.NewSha
	}
	res, err := db.Exec(sqlInsertRewritten,
		runID, sourceCommitID, newSha, r.FinalMessage,
		r.AuthorName, r.AuthorEmail,
		r.AuthorDateRFC3339, r.CommitterDateRFC3339, outcomeID,
	)
	if err != nil {
		return 0, fmt.Errorf("runlog: insert RewrittenCommit: %w", err)
	}
	return res.LastInsertId()
}

// insertShaMap upserts the cross-run dedupe row.
func insertShaMap(db *sql.DB, sourceSha string, rewrittenID int64) error {
	if _, err := db.Exec(sqlInsertShaMap, sourceSha, rewrittenID); err != nil {
		return fmt.Errorf("runlog: insert ShaMap %s: %w", sourceSha, err)
	}
	return nil
}

const (
	sqlInsertRewritten = `INSERT INTO RewrittenCommit
		(CommitInRunId, SourceCommitId, NewSha, FinalMessage,
		 AppliedAuthorName, AppliedAuthorEmail,
		 AppliedAuthorDate, AppliedCommitterDate, CommitOutcomeId)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlInsertSkipLog = `INSERT INTO SkipLog
		(CommitInRunId, SourceCommitId, SkipReasonId, PreviousRewrittenCommitId)
		VALUES (?, ?, ?, ?)`

	sqlInsertShaMap = `INSERT INTO ShaMap (SourceSha, RewrittenCommitId)
		VALUES (?, ?)`
)
