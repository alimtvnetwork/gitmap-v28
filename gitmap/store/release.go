package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// UpsertRelease inserts or updates a release record in the database.
// v17: requires RepoID FK pointing at an existing Repo row.
func (db *DB) UpsertRelease(r model.ReleaseRecord) error {
	if r.RepoID == 0 {
		return fmt.Errorf(constants.ErrReleaseNoRepo, "<unset>")
	}

	return db.runUpsertReleaseTx(r)
}

func (db *DB) runUpsertReleaseTx(r model.ReleaseRecord) error {
	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	appErr = wrap.WithTransaction(context.Background(), func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeUpsertReleaseTx(tx.Tx(), r)
	})
	if appErr != nil {
		return fmt.Errorf(constants.ErrDBReleaseUpsert, appErr)
	}

	return nil
}

func executeUpsertReleaseTx(tx *sql.Tx, r model.ReleaseRecord) *apperror.AppError {
	if err := upsertReleaseTx(tx, r); err != nil {
		return apperror.WrapSimple(err, "executeUpsertReleaseTx")
	}

	return nil
}

func syncLatestReleaseFlag(runner sqlExecutor, r model.ReleaseRecord) error {
	if !r.IsLatest {
		return nil
	}

	return clearLatestRunner(runner, r.RepoID)
}

func upsertReleaseTx(runner sqlExecutor, r model.ReleaseRecord) error {
	if err := syncLatestReleaseFlag(runner, r); err != nil {
		return err
	}

	return insertOrUpdateRelease(runner, r)
}

func insertOrUpdateRelease(runner sqlExecutor, r model.ReleaseRecord) error {
	isDraft := boolToInt(r.IsDraft)
	isPreRelease := boolToInt(r.IsPreRelease)
	isLatest := boolToInt(r.IsLatest)

	_, err := ExecWrapper(runner, constants.SQLUpsertRelease,
		r.RepoID, r.Version, r.Tag, r.Branch, r.SourceBranch,
		r.CommitSha, r.Changelog, r.Notes, isDraft, isPreRelease, isLatest, r.Source, r.CreatedAt,
	).Destruct()

	return err
}

func clearLatestRunner(runner sqlExecutor, repoID int64) error {
	_, err := ExecWrapper(runner, constants.SQLClearLatestRelease, repoID).Destruct()
	if err != nil {
		return fmt.Errorf(constants.ErrDBReleaseUpsert, err)
	}

	return nil
}

// ListReleases returns all releases ordered by creation date descending.
func (db *DB) ListReleases() ([]model.ReleaseRecord, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllReleases).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBReleaseQuery, err)
	}

	defer rows.Close()

	return scanReleaseRows(rows)
}

// FindReleaseByTag returns a release matching the given tag.
func (db *DB) FindReleaseByTag(tag string) (model.ReleaseRecord, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectReleaseByTag, tag)

	return scanOneRelease(row)
}

// scanReleaseRows reads ReleaseRecord values from query result rows.
func scanReleaseRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.ReleaseRecord, error) {
	var results []model.ReleaseRecord

	for rows.Next() {
		r, err := scanOneReleaseRow(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}

// scanOneReleaseRow reads a single ReleaseRecord from a row scanner.
// v17: RepoId is now part of the projected column list.
func scanOneReleaseRow(row interface{ Scan(dest ...any) error }) (model.ReleaseRecord, error) {
	var r model.ReleaseRecord
	var isDraft, isPreRelease, isLatest int

	err := row.Scan(&r.ID, &r.RepoID, &r.Version, &r.Tag, &r.Branch, &r.SourceBranch,
		&r.CommitSha, &r.Changelog, &r.Notes, &isDraft, &isPreRelease, &isLatest, &r.Source, &r.CreatedAt)
	if err != nil {
		return model.ReleaseRecord{}, err
	}

	r.IsDraft = isDraft == 1
	r.IsPreRelease = isPreRelease == 1
	r.IsLatest = isLatest == 1

	return r, nil
}

// scanOneRelease reads a single ReleaseRecord from a QueryRow result.
func scanOneRelease(row interface{ Scan(dest ...any) error }) (model.ReleaseRecord, error) {
	return scanOneReleaseRow(row)
}

// JoinChangelog joins changelog notes into a newline-separated string.
func JoinChangelog(notes []string) string {
	if len(notes) == 0 {
		return ""
	}

	return strings.Join(notes, "\n")
}

// boolToInt converts a bool to SQLite integer (0/1).
func boolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

// ResolveCurrentRepoID returns the RepoId for the repo at absPath. Returns 0
// and an error when the repo has not been scanned (caller should advise the
// user to run `gitmap scan` first).
func (db *DB) ResolveCurrentRepoID(absPath string) (int64, error) {
	repos, err := db.FindByPath(absPath)
	if err != nil {
		return 0, err
	}

	if len(repos) == 0 {
		return 0, fmt.Errorf(constants.ErrReleaseNoRepo, absPath)
	}

	return repos[0].ID, nil
}
