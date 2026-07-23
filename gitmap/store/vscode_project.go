package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// UpsertVSCodeProject inserts or updates a row keyed by RootPath
// (case-insensitive). Bumps Name + LastSeenAt + UpdatedAt on conflict.
// The Paths column is NOT touched here — use SetVSCodeProjectPaths
// to mutate the DB-side multi-root list.
func (db *DB) UpsertVSCodeProject(rootPath, name string) error {
	if _, err := db.conn.Exec(constants.SQLUpsertVSCodeProject, rootPath, name); err != nil {
		return fmt.Errorf(constants.ErrVSCodePMUpsert, rootPath, err)
	}

	return nil
}

// ListVSCodeProjects returns every row in the VSCodeProject table.
func (db *DB) ListVSCodeProjects() ([]model.VSCodeProject, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllVSCodeProjects)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrVSCodePMList, err)
	}
	defer rows.Close()

	return scanVSCodeProjectRows(rows)
}

// FindVSCodeProjectByPath returns the row matching RootPath (case-insensitive)
// or sql.ErrNoRows when missing.
func (db *DB) FindVSCodeProjectByPath(rootPath string) (model.VSCodeProject, error) {
	row := db.conn.QueryRow(constants.SQLSelectVSCodeProjectByPath, rootPath)

	return scanOneVSCodeProjectRow(row)
}

// FindVSCodeProjectByName returns the first row whose Name matches (case-
// insensitive) or sql.ErrNoRows. Used by `gitmap code paths` to look up
// an alias without requiring the user to supply rootPath.
func (db *DB) FindVSCodeProjectByName(name string) (model.VSCodeProject, error) {
	row := db.conn.QueryRow(constants.SQLSelectVSCodeProjectByName, name)

	return scanOneVSCodeProjectRow(row)
}

// RenameVSCodeProjectByPath updates the Name column for the matching RootPath.
// Returns the number of rows affected so callers can detect "no match".
func (db *DB) RenameVSCodeProjectByPath(rootPath, newName string) (int64, error) {
	res, err := db.conn.Exec(constants.SQLRenameVSCodeProject, newName, rootPath)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrVSCodePMRename, rootPath, err)
	}

	affected, _ := res.RowsAffected()

	return affected, nil
}

// SetVSCodeProjectPaths replaces the JSON-encoded Paths column for a row.
// Caller is responsible for de-duplication.
func (db *DB) SetVSCodeProjectPaths(rootPath string, paths []string) error {
	encoded, err := encodePaths(rootPath, paths)
	if err != nil {
		return err
	}

	if _, err := db.conn.Exec(constants.SQLUpdateVSCodeProjectPaths, encoded, rootPath); err != nil {
		return fmt.Errorf(constants.ErrVSCodePMUpdatePaths, rootPath, err)
	}

	return nil
}

// DeleteVSCodeProjectByPath removes a row by RootPath.
func (db *DB) DeleteVSCodeProjectByPath(rootPath string) error {
	if _, err := db.conn.Exec(constants.SQLDeleteVSCodeProjectByPath, rootPath); err != nil {
		return fmt.Errorf(constants.ErrVSCodePMDelete, rootPath, err)
	}

	return nil
}

// scanOneVSCodeProjectRow narrows sql.ErrNoRows so callers can switch on it.
func scanOneVSCodeProjectRow(row interface{ Scan(dest ...any) error }) (model.VSCodeProject, error) {
	p, err := scanOneVSCodeProject(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.VSCodeProject{}, sql.ErrNoRows
		}

		return model.VSCodeProject{}, fmt.Errorf(constants.ErrVSCodePMList, err)
	}

	return p, nil
}

// scanOneVSCodeProject reads a single VSCodeProject row including the
// JSON-encoded Paths column.
func scanOneVSCodeProject(row interface{ Scan(dest ...any) error }) (model.VSCodeProject, error) {
	var (
		p         model.VSCodeProject
		enabled   int64
		pathsJSON string
	)

	err := row.Scan(&p.ID, &p.RootPath, &p.Name, &pathsJSON, &enabled, &p.Profile,
		&p.LastSeenAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, err
	}

	p.Enabled = enabled != 0

	paths, decodeErr := decodePaths(p.RootPath, pathsJSON)
	if decodeErr != nil {
		return p, decodeErr
	}
	p.Paths = paths

	return p, nil
}

// scanVSCodeProjectRows reads VSCodeProject values from query result rows.
func scanVSCodeProjectRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.VSCodeProject, error) {
	var results []model.VSCodeProject

	for rows.Next() {
		p, err := scanOneVSCodeProject(rows)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrVSCodePMList, err)
		}

		results = append(results, p)
	}

	return results, nil
}

// encodePaths marshals the multi-root list into the JSON shape stored in
// the Paths TEXT column. Empty/nil slices encode to "[]" so the column
// never holds NULL.
func encodePaths(rootPath string, paths []string) (string, error) {
	if paths == nil {
		paths = []string{}
	}

	bytes, err := json.Marshal(paths)
	if err != nil {
		return "", fmt.Errorf(constants.ErrVSCodePMPathsEncode, rootPath, err)
	}

	return string(bytes), nil
}

// decodePaths unmarshals the Paths TEXT column. Empty / "[]" / "null"
// all decode to a non-nil empty slice.
func decodePaths(rootPath, encoded string) ([]string, error) {
	if encoded == "" || encoded == "null" {
		return []string{}, nil
	}

	var paths []string
	if err := json.Unmarshal([]byte(encoded), &paths); err != nil {
		return nil, fmt.Errorf(constants.ErrVSCodePMPathsDecode, rootPath, err)
	}

	if paths == nil {
		paths = []string{}
	}

	return paths, nil
}
