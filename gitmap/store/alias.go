// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// AliasWithRepo extends Alias with resolved repo details.
type AliasWithRepo struct {
	model.Alias
	AbsolutePath string
	Slug         string
}

// UnaliasedRepo holds a repo that has no alias assigned.
type UnaliasedRepo struct {
	ID       int64
	Slug     string
	RepoName string
}

// CreateAlias inserts a new alias for the given repo ID.
func (db *DB) CreateAlias(alias string, repoID int64) (model.Alias, error) {
	_, err := db.conn.Exec(constants.SQLInsertAlias, alias, repoID)
	if err != nil {
		return model.Alias{}, fmt.Errorf(constants.ErrAliasCreate, err)
	}

	return db.FindAliasByName(alias)
}

// UpdateAlias reassigns an existing alias to a different repo.
func (db *DB) UpdateAlias(alias string, repoID int64) error {
	_, err := db.conn.Exec(constants.SQLUpdateAlias, repoID, alias)
	if err != nil {
		return fmt.Errorf(constants.ErrAliasCreate, err)
	}

	return nil
}

// FindAliasByName retrieves a single alias by its name.
func (db *DB) FindAliasByName(alias string) (model.Alias, error) {
	row := db.conn.QueryRow(constants.SQLSelectAliasByName, alias)

	return scanOneAlias(row)
}

// FindAliasByRepoID retrieves the alias for a specific repo.
func (db *DB) FindAliasByRepoID(repoID int64) (model.Alias, error) {
	row := db.conn.QueryRow(constants.SQLSelectAliasByRepoID, repoID)

	return scanOneAlias(row)
}

// ListAliases returns all aliases ordered by name.
func (db *DB) ListAliases() ([]model.Alias, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllAliases)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrAliasQuery, err)
	}
	defer rows.Close()

	return scanAliasRows(rows)
}

// ResolveAlias retrieves an alias with its repo path and slug.
func (db *DB) ResolveAlias(alias string) (AliasWithRepo, error) {
	row := db.conn.QueryRow(constants.SQLSelectAliasWithRepo, alias)

	var a AliasWithRepo

	err := row.Scan(&a.ID, &a.Alias.Alias, &a.RepoID, &a.CreatedAt, &a.AbsolutePath, &a.Slug)
	if err != nil {
		return AliasWithRepo{}, fmt.Errorf(constants.ErrAliasNotFound, alias)
	}

	return a, nil
}

// ListAliasesWithRepo returns all aliases with resolved repo details.
func (db *DB) ListAliasesWithRepo() ([]AliasWithRepo, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllAliasesWithRepo)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrAliasQuery, err)
	}
	defer rows.Close()

	var results []AliasWithRepo

	for rows.Next() {
		var a AliasWithRepo

		err := rows.Scan(&a.ID, &a.Alias.Alias, &a.RepoID, &a.CreatedAt, &a.AbsolutePath, &a.Slug)
		if err != nil {
			continue
		}

		results = append(results, a)
	}

	return results, nil
}

// DeleteAlias removes an alias by name.
func (db *DB) DeleteAlias(alias string) error {
	result, err := db.conn.Exec(constants.SQLDeleteAlias, alias)
	if err != nil {
		return fmt.Errorf(constants.ErrAliasDelete, err)
	}

	return checkDeleted(result, alias)
}

// ListUnaliasedRepos returns repos that have no alias assigned.
func (db *DB) ListUnaliasedRepos() ([]UnaliasedRepo, error) {
	rows, err := db.conn.Query(constants.SQLSelectUnaliasedRepos)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrAliasQuery, err)
	}
	defer rows.Close()

	var repos []UnaliasedRepo

	for rows.Next() {
		var r UnaliasedRepo

		err := rows.Scan(&r.ID, &r.Slug, &r.RepoName)
		if err != nil {
			continue
		}

		repos = append(repos, r)
	}

	return repos, nil
}

// AliasExists returns true if an alias with the given name exists.
func (db *DB) AliasExists(alias string) bool {
	_, err := db.FindAliasByName(alias)

	return err == nil
}

// scanOneAlias scans a single alias row.
func scanOneAlias(row *sql.Row) (model.Alias, error) {
	var a model.Alias

	err := row.Scan(&a.ID, &a.Alias, &a.RepoID, &a.CreatedAt)
	if err != nil {
		return model.Alias{}, err
	}

	return a, nil
}

// scanAliasRows scans multiple alias rows.
func scanAliasRows(rows *sql.Rows) ([]model.Alias, error) { //nolint:unparam // error kept for interface consistency
	var aliases []model.Alias

	for rows.Next() {
		var a model.Alias

		err := rows.Scan(&a.ID, &a.Alias, &a.RepoID, &a.CreatedAt)
		if err != nil {
			continue
		}

		aliases = append(aliases, a)
	}

	return aliases, nil
}
