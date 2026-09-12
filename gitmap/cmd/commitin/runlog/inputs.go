package runlog

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

// InsertInputRepo persists one staged input. Returns the new
// InputRepoId so the walker can FK SourceCommit rows back to it.
func InsertInputRepo(
	db *sql.DB,
	runID int64,
	orderIndex int,
	originalRef,
	resolvedPath,
	kind string,
) (int64, error) {
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
	wrapper, appErr := dbengine.WrapDb(db, dbengine.DbSQLite)
	if appErr != nil {
		return 0, fmt.Errorf("runlog: wrap db: %w", appErr)
	}

	var id int64
	txErr := wrapper.WithTransaction(context.Background(), func(tx *dbengine.TxWrapper) *apperror.AppError {
		var err *apperror.AppError
		id, err = executeCommitInsertTx(tx, inputRepoID, c)

		return err
	})
	if txErr != nil {
		return 0, txErr
	}

	return id, nil
}

func executeCommitInsertTx(tx *dbengine.TxWrapper, inputRepoID int64, c SourceCommitRow) (int64, *apperror.AppError) {
	id, err := insertSourceCommitTx(tx, inputRepoID, c)
	if err != nil {
		return 0, err
	}
	if err := insertSourceFilesTx(tx, id, c.Files); err != nil {
		return 0, err
	}

	return id, nil
}

// insertSourceCommitTx writes the SourceCommit row only.
func insertSourceCommitTx(tx *dbengine.TxWrapper, inputRepoID int64, c SourceCommitRow) (int64, *apperror.AppError) {
	res, err := tx.Exec(
		context.Background(), sqlInsertSourceCommit,
		inputRepoID, c.Sha, c.AuthorName, c.AuthorEmail,
		c.AuthorDateRFC3339, c.CommitterDateRFC3339,
		c.OriginalMessage, c.OrderIndex,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, fmt.Sprintf("runlog: insert SourceCommit %s", c.Sha))
	}

	rowID, idErr := res.LastInsertId()
	if idErr != nil {
		return 0, apperror.WrapSimple(idErr, "get last insert id for source commit")
	}

	return rowID, nil
}

// insertSourceFilesTx batch-writes one row per touched file.
func insertSourceFilesTx(tx *dbengine.TxWrapper, sourceCommitID int64, files []string) *apperror.AppError {
	ctx := context.Background()
	for _, rel := range files {
		if _, err := tx.Exec(ctx, sqlInsertSourceFile, sourceCommitID, rel); err != nil {
			return apperror.WrapSimple(err, fmt.Sprintf("runlog: insert SourceCommitFile %q", rel))
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
