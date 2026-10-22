// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
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
	_, err := ExecWrapper(db.conn, constants.SQLInsertAlias, alias, repoID).Destruct()
	if err != nil {
		return model.Alias{}, apperror.WrapSimple(err, "failed to create alias")
	}

	return db.FindAliasByName(alias)
}

func (db *DB) ensureAliasColumns() {
	if !db.columnExists(constants.TableAlias, "IsPrimary") || !db.columnExists(constants.TableAlias, "Source") {
		db.migrateAliasColumns()
	}
}

// CreateAliasWithDetails inserts a new alias with primary flag and source.
func (db *DB) CreateAliasWithDetails(alias string, repoID int64, isPrimary bool, source string) (model.Alias, error) {
	db.ensureAliasColumns()
	primaryVal := 0
	if isPrimary {
		primaryVal = 1
	}
	_, err := ExecWrapper(db.conn, constants.SQLInsertAliasFull, alias, repoID, primaryVal, source).Destruct()
	if err != nil {
		return model.Alias{}, apperror.WrapSimple(err, "failed to create alias")
	}

	return db.FindAliasByName(alias)
}

// UpdateAlias reassigns an existing alias to a different repo.
func (db *DB) UpdateAlias(alias string, repoID int64) error {
	_, err := ExecWrapper(db.conn, constants.SQLUpdateAlias, repoID, alias).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to create alias")
	}

	return nil
}

// FindAliasByName retrieves a single alias by its name.
func (db *DB) FindAliasByName(alias string) (model.Alias, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectAliasByName, alias)

	return scanOneAlias(row)
}

// FindAliasByRepoID retrieves the alias for a specific repo.
func (db *DB) FindAliasByRepoID(repoID int64) (model.Alias, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectAliasByRepoID, repoID)

	return scanOneAlias(row)
}

// FindAliasesByRepoID retrieves all aliases for a specific repo.
func (db *DB) FindAliasesByRepoID(repoID int64) ([]model.Alias, error) {
	res := db.FindAliasesByRepoIDWrap(repoID)
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

func (db *DB) FindAliasesByRepoIDWrap(repoID int64) AliasSliceWrap {
	db.ensureAliasColumns()
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAliasesByRepoID, repoID).Destruct()
	if err != nil {
		return result.FailSlice[model.Alias](apperror.WrapSimple(err, "FindAliasesByRepoID"))
	}

	defer rows.Close()

	items, errScan := scanAliasRowsFull(rows)
	if errScan != nil {
		return result.FailSlice[model.Alias](apperror.WrapSimple(errScan, "FindAliasesByRepoID"))
	}

	return result.OkSlice(items)
}

// ListAliases returns all aliases ordered by name.
func (db *DB) ListAliases() ([]model.Alias, error) {
	res := db.ListAliasesWrap()
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

func (db *DB) ListAliasesWrap() AliasSliceWrap {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllAliases).Destruct()
	if err != nil {
		return result.FailSlice[model.Alias](apperror.WrapSimple(err, "ListAliases"))
	}

	defer rows.Close()

	items, errScan := scanAliasRows(rows)
	if errScan != nil {
		return result.FailSlice[model.Alias](apperror.WrapSimple(errScan, "ListAliases"))
	}

	return result.OkSlice(items)
}

// ResolveAlias retrieves an alias with its repo path and slug.
func (db *DB) ResolveAlias(alias string) (AliasWithRepo, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectAliasWithRepo, alias)

	var a AliasWithRepo

	err := row.Scan(&a.ID, &a.Alias.Alias, &a.RepoID, &a.CreatedAt, &a.AbsolutePath, &a.Slug)
	if err != nil {
		return AliasWithRepo{}, apperror.WrapNotFound(err, fmt.Sprintf("no alias found: %s", alias))
	}

	return a, nil
}

// ListAliasesWithRepo returns all aliases with resolved repo details.
func (db *DB) ListAliasesWithRepo() ([]AliasWithRepo, error) {
	res := db.ListAliasesWithRepoWrap()
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

func (db *DB) ListAliasesWithRepoWrap() AliasWithRepoSliceWrap {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllAliasesWithRepo).Destruct()
	if err != nil {
		return result.FailSlice[AliasWithRepo](apperror.WrapSimple(err, "ListAliasesWithRepo"))
	}

	defer rows.Close()

	var results []AliasWithRepo

	for rows.Next() {
		var a AliasWithRepo

		scanErr := rows.Scan(&a.ID, &a.Alias.Alias, &a.RepoID, &a.CreatedAt, &a.AbsolutePath, &a.Slug)
		if scanErr != nil {
			continue
		}

		results = append(results, a)
	}

	return result.OkSlice(results)
}

// DeleteAlias removes an alias by name.
func (db *DB) DeleteAlias(alias string) error {
	result, err := ExecWrapper(db.conn, constants.SQLDeleteAlias, alias).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to delete alias")
	}

	return checkDeleted(result, alias)
}

// ListUnaliasedRepos returns repos that have no alias assigned.
func (db *DB) ListUnaliasedRepos() ([]UnaliasedRepo, error) {
	res := db.ListUnaliasedReposWrap()
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

func (db *DB) ListUnaliasedReposWrap() UnaliasedRepoSliceWrap {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectUnaliasedRepos).Destruct()
	if err != nil {
		return result.FailSlice[UnaliasedRepo](apperror.WrapSimple(err, "ListUnaliasedRepos"))
	}

	defer rows.Close()

	var repos []UnaliasedRepo

	for rows.Next() {
		var r UnaliasedRepo

		scanErr := rows.Scan(&r.ID, &r.Slug, &r.RepoName)
		if scanErr != nil {
			continue
		}

		repos = append(repos, r)
	}

	return result.OkSlice(repos)
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

// scanAliasRowsFull scans multiple alias rows including IsPrimary and Source.
func scanAliasRowsFull(rows *sql.Rows) ([]model.Alias, error) { //nolint:unparam // error kept for interface consistency
	var aliases []model.Alias

	for rows.Next() {
		if a := scanSingleAliasFull(rows); a != nil {
			aliases = append(aliases, *a)
		}
	}

	return aliases, nil
}

func scanSingleAliasFull(rows *sql.Rows) *model.Alias {
	var a model.Alias
	var isPrimaryInt int
	if err := rows.Scan(&a.ID, &a.Alias, &a.RepoID, &isPrimaryInt, &a.Source, &a.CreatedAt); err != nil {
		return nil
	}
	a.IsPrimary = isPrimaryInt == 1

	return &a
}
