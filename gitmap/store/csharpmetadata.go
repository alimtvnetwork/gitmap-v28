// Package store — csharpmetadata.go manages Csharp metadata + files.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// UpsertCsharpMetadata inserts or updates C# metadata for a detected project.
func (db *DB) UpsertCsharpMetadata(m model.CsharpProjectMetadata) error {
	_, err := db.conn.Exec(constants.SQLUpsertCsharpMetadata,
		m.DetectedProjectID, m.SlnPath, m.SlnName,
		m.GlobalJsonPath, m.SdkVersion)

	return err
}

// UpsertCsharpProjectFile inserts or updates a C# project file record.
func (db *DB) UpsertCsharpProjectFile(f model.CsharpProjectFile) error {
	_, err := db.conn.Exec(constants.SQLUpsertCsharpProjectFile,
		f.CsharpMetadataID, f.FilePath, f.RelativePath,
		f.FileName, f.ProjectName, f.TargetFramework, f.OutputType, f.Sdk)

	return err
}

// UpsertCsharpKeyFile inserts or updates a C# key file record.
func (db *DB) UpsertCsharpKeyFile(f model.CsharpKeyFile) error {
	_, err := db.conn.Exec(constants.SQLUpsertCsharpKeyFile,
		f.CsharpMetadataID, f.FileType, f.FilePath, f.RelativePath)

	return err
}

// SelectCsharpMetadata returns C# metadata for a detected project.
func (db *DB) SelectCsharpMetadata(detectedProjectID int64) (*model.CsharpProjectMetadata, error) {
	var m model.CsharpProjectMetadata
	err := db.conn.QueryRow(constants.SQLSelectCsharpMetadata, detectedProjectID).Scan(
		&m.ID, &m.DetectedProjectID, &m.SlnPath, &m.SlnName,
		&m.GlobalJsonPath, &m.SdkVersion)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// SelectCsharpProjectFiles returns all .csproj files for a metadata ID.
func (db *DB) SelectCsharpProjectFiles(metadataID int64) ([]model.CsharpProjectFile, error) {
	rows, err := db.conn.Query(constants.SQLSelectCsharpProjectFiles, metadataID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrProjectQuery, err)
	}
	defer rows.Close()

	return scanCsharpFileRows(rows)
}

// SelectCsharpKeyFiles returns all key files for a metadata ID.
func (db *DB) SelectCsharpKeyFiles(metadataID int64) ([]model.CsharpKeyFile, error) {
	rows, err := db.conn.Query(constants.SQLSelectCsharpKeyFiles, metadataID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrProjectQuery, err)
	}
	defer rows.Close()

	return scanCsharpKeyFileRows(rows)
}

// DeleteStaleCsharpFiles removes project files not in the keep list.
func (db *DB) DeleteStaleCsharpFiles(metadataID int64, keepIDs []int64) error {
	if len(keepIDs) == 0 {
		return nil
	}
	placeholders := buildPlaceholders(len(keepIDs))
	query := fmt.Sprintf(constants.SQLDeleteStaleCsharpFiles, placeholders)
	args := buildStaleArgsInt64(metadataID, keepIDs)
	_, err := db.conn.Exec(query, args...)

	return err
}

// DeleteStaleCsharpKeyFiles removes key files not in the keep list.
func (db *DB) DeleteStaleCsharpKeyFiles(metadataID int64, keepIDs []int64) error {
	if len(keepIDs) == 0 {
		return nil
	}
	placeholders := buildPlaceholders(len(keepIDs))
	query := fmt.Sprintf(constants.SQLDeleteStaleCsharpKeyFiles, placeholders)
	args := buildStaleArgsInt64(metadataID, keepIDs)
	_, err := db.conn.Exec(query, args...)

	return err
}

// scanCsharpFileRows scans rows into CsharpProjectFile slices.
func scanCsharpFileRows(rows interface {
	Next() bool
	Scan(...interface{}) error
}) ([]model.CsharpProjectFile, error) {
	var files []model.CsharpProjectFile
	for rows.Next() {
		var f model.CsharpProjectFile
		err := rows.Scan(&f.ID, &f.CsharpMetadataID, &f.FilePath,
			&f.RelativePath, &f.FileName, &f.ProjectName,
			&f.TargetFramework, &f.OutputType, &f.Sdk)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

// scanCsharpKeyFileRows scans rows into CsharpKeyFile slices.
func scanCsharpKeyFileRows(rows interface {
	Next() bool
	Scan(...interface{}) error
}) ([]model.CsharpKeyFile, error) {
	var files []model.CsharpKeyFile
	for rows.Next() {
		var f model.CsharpKeyFile
		err := rows.Scan(&f.ID, &f.CsharpMetadataID, &f.FileType,
			&f.FilePath, &f.RelativePath)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}
