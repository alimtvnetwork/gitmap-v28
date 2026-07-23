package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// FindNext returns every repo whose latest VersionProbe row reports an
// available update. When scanFolderID > 0, results are scoped to that
// scan folder; pass 0 to query the whole database.
func (db *DB) FindNext(scanFolderID int64) ([]model.FindNextRow, error) {
	rows, err := queryFindNext(db.conn, scanFolderID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFindNextQuery, err)
	}
	defer rows.Close()

	return scanFindNextRows(rows)
}

// queryFindNext picks the right SQL based on the scan-folder filter.
func queryFindNext(conn *sql.DB, scanFolderID int64) (*sql.Rows, error) {
	if scanFolderID > 0 {
		return conn.Query(constants.SQLSelectFindNextByScanFolder, scanFolderID)
	}

	return conn.Query(constants.SQLSelectFindNext)
}

// scanFindNextRows materializes every result row into model.FindNextRow.
func scanFindNextRows(rows *sql.Rows) ([]model.FindNextRow, error) {
	var results []model.FindNextRow
	for rows.Next() {
		row, err := scanOneFindNextRow(rows)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrFindNextScanRow, err)
		}
		results = append(results, row)
	}

	return results, nil
}

// scanOneFindNextRow reads one joined Repo + VersionProbe row.
func scanOneFindNextRow(rows *sql.Rows) (model.FindNextRow, error) {
	var r model.FindNextRow
	err := rows.Scan(
		&r.Repo.ID, &r.Repo.Slug, &r.Repo.RepoName, &r.Repo.HTTPSUrl, &r.Repo.SSHUrl,
		&r.Repo.Branch, &r.Repo.RelativePath, &r.Repo.AbsolutePath,
		&r.Repo.CloneInstruction, &r.Repo.Notes,
		&r.NextVersionTag, &r.NextVersionNum, &r.Method, &r.ProbedAt,
	)

	return r, err
}
