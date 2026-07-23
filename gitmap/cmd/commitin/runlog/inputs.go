package runlog

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// InsertInputRepo persists one staged input. Returns the new
// InputRepoId so the walker can FK SourceCommit rows back to it.
func InsertInputRepo(db *sql.DB, runID int64, orderIndex int, originalRef, resolvedPath, kind string) (int64, error) {
	kindID, err := lookupEnumID(db, constants.TableCommitInInputKind, "InputKindId", kind)
	if err != nil {
		return 0, fmt.Errorf("runlog: lookup InputKind %q: %w", kind, err)
	}
	res, err := db.Exec(sqlInsertInputRepo, runID, orderIndex, originalRef, resolvedPath, kindID)
	if err != nil {
		return 0, fmt.Errorf("runlog: insert InputRepo: %w", err)
	}
	return res.LastInsertId()
}

// InsertSourceCommit persists one walked commit + its files. Wraps
// both writes in a single transaction so a partial insert never
// leaves orphaned rows.
func InsertSourceCommit(db *sql.DB, inputRepoID int64, c SourceCommitRow) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("runlog: begin tx: %w", err)
	}
	id, err := insertSourceCommitTx(tx, inputRepoID, c)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if err := insertSourceFilesTx(tx, id, c.Files); err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("runlog: commit tx: %w", err)
	}
	return id, nil
}

// insertSourceCommitTx writes the SourceCommit row only.
func insertSourceCommitTx(tx *sql.Tx, inputRepoID int64, c SourceCommitRow) (int64, error) {
	res, err := tx.Exec(sqlInsertSourceCommit,
		inputRepoID, c.Sha, c.AuthorName, c.AuthorEmail,
		c.AuthorDateRFC3339, c.CommitterDateRFC3339,
		c.OriginalMessage, c.OrderIndex,
	)
	if err != nil {
		return 0, fmt.Errorf("runlog: insert SourceCommit %s: %w", c.Sha, err)
	}
	return res.LastInsertId()
}

// insertSourceFilesTx batch-writes one row per touched file.
func insertSourceFilesTx(tx *sql.Tx, sourceCommitID int64, files []string) error {
	for _, rel := range files {
		if _, err := tx.Exec(sqlInsertSourceFile, sourceCommitID, rel); err != nil {
			return fmt.Errorf("runlog: insert SourceCommitFile %q: %w", rel, err)
		}
	}
	return nil
}

const (
	sqlInsertInputRepo = `INSERT INTO InputRepo
		(CommitInRunId, OrderIndex, OriginalRef, ResolvedPath, InputKindId)
		VALUES (?, ?, ?, ?, ?)`

	sqlInsertSourceCommit = `INSERT INTO SourceCommit
		(InputRepoId, SourceSha, AuthorName, AuthorEmail,
		 AuthorDate, CommitterDate, OriginalMessage, OrderIndex)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	sqlInsertSourceFile = `INSERT INTO SourceCommitFile
		(SourceCommitId, RelativePath) VALUES (?, ?)`
)
