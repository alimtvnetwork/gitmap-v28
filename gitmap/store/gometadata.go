// Package store — gometadata.go manages GoProjectMetadata + GoRunnableFiles.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// UpsertGoMetadata inserts or updates Go metadata for a detected project.
func (db *DB) UpsertGoMetadata(m model.GoProjectMetadata) error {
	_, err := ExecWrapper(db.conn, constants.SQLUpsertGoMetadata,
		m.DetectedProjectID, m.GoModPath, m.GoSumPath,
		m.ModuleName, m.GoVersion).Destruct()

	return err
}

// UpsertGoRunnable inserts or updates a Go runnable file record.
func (db *DB) UpsertGoRunnable(r model.GoRunnableFile) error {
	_, err := ExecWrapper(db.conn, constants.SQLUpsertGoRunnable,
		r.GoMetadataID, r.RunnableName, r.FilePath, r.RelativePath).Destruct()

	return err
}

// SelectGoMetadata returns Go metadata for a detected project.
func (db *DB) SelectGoMetadata(detectedProjectID int64) (*model.GoProjectMetadata, error) {
	var m model.GoProjectMetadata
	err := db.conn.QueryRow(constants.SQLSelectGoMetadata, detectedProjectID).Scan(
		&m.ID, &m.DetectedProjectID, &m.GoModPath, &m.GoSumPath,
		&m.ModuleName, &m.GoVersion)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// SelectGoRunnables returns all runnable files for a Go metadata ID.
func (db *DB) SelectGoRunnables(goMetadataID int64) ([]model.GoRunnableFile, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectGoRunnables, goMetadataID).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrProjectQuery, err)
	}
	defer rows.Close()

	return scanGoRunnableRows(rows)
}

// DeleteStaleGoRunnables removes runnables not in the given ID list.
func (db *DB) DeleteStaleGoRunnables(goMetadataID int64, keepIDs []int64) error {
	if len(keepIDs) == 0 {
		return nil
	}
	placeholders := buildPlaceholders(len(keepIDs))
	query := fmt.Sprintf(constants.SQLDeleteStaleGoRunnables, placeholders)
	args := buildStaleArgsInt64(goMetadataID, keepIDs)
	_, err := ExecWrapper(db.conn, query, args...).Destruct()

	return err
}

// scanGoRunnableRows scans rows into GoRunnableFile slices.
func scanGoRunnableRows(rows interface {
	Next() bool
	Scan(...interface{}) error
}) ([]model.GoRunnableFile, error) {
	var runnables []model.GoRunnableFile
	for rows.Next() {
		var r model.GoRunnableFile
		err := rows.Scan(&r.ID, &r.GoMetadataID, &r.RunnableName,
			&r.FilePath, &r.RelativePath)
		if err != nil {
			return nil, err
		}
		runnables = append(runnables, r)
	}

	return runnables, nil
}
