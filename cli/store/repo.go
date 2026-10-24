package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// NormalizeStoragePath produces a canonical, clean absolute path for database persistence.
// Storage paths use forward slashes for cross-platform portability.
// On Windows, drive letters are capitalized to ensure consistent visual presentation.
func NormalizeStoragePath(pathStr string) string {
	// Convert Windows separators explicitly before cleaning:
	// filepath.ToSlash is a no-op on Unix, where '\' is a legal filename
	// character, so a Windows-style path would otherwise never match the
	// forward-slash form the record was stored under.
	slashPath := strings.ReplaceAll(strings.TrimSpace(pathStr), "\\", "/")
	clean := filepath.ToSlash(filepath.Clean(slashPath))
	vol := filepath.VolumeName(clean)
	if len(vol) >= 2 && vol[1] == ':' {
		clean = strings.ToUpper(string(vol[0])) + clean[1:]
	}

	return clean
}

// OnRepoAddedHook is an optional callback triggered when repositories are upserted into the catalog.
var OnRepoAddedHook func(records []model.ScanRecord)

// UpsertRepos inserts or updates all records by absolute_path.
func (db *DB) UpsertRepos(records []model.ScanRecord) error {
	wrap, appErr := dbengine.WrapDB(db.conn, dbengine.DBSQLite)
	if appErr != nil {
		return appErr
	}

	appErr = wrap.WithTransaction(context.Background(), func(tx *dbengine.TxWrapper) *apperror.AppError {
		return db.upsertReposTx(tx.Tx(), records)
	})
	if appErr != nil {
		return apperror.WrapSimple(appErr, "failed to upsert repo")
	}

	if OnRepoAddedHook != nil && len(records) > 0 {
		OnRepoAddedHook(records)
	}

	return nil
}

func (db *DB) upsertReposTx(tx *sql.Tx, records []model.ScanRecord) *apperror.AppError {
	for _, r := range records {
		if err := upsertOneRepo(tx, r); err != nil {
			return apperror.WrapSimple(err, "upsertOneRepo")
		}
	}

	return nil
}

func upsertOneRepo(runner sqlExecutor, r model.ScanRecord) error {
	normalizedPath := NormalizeStoragePath(r.AbsolutePath)
	_, err := ExecWrapper(runner, constants.SQLUpsertRepoByPath,
		r.Slug, r.RepoName, r.HTTPSUrl, r.SSHUrl,
		r.Branch, r.RelativePath, normalizedPath,
		r.CloneInstruction, r.Notes, r.IdentifiedTransport,
	).Destruct()

	return err
}

// DeleteByPathWrap removes repo rows matching absolute path encapsulated in RowsAffectedWrap.
func (db *DB) DeleteByPathWrap(absPath string) RowsAffectedWrap {
	cleanPath := NormalizeStoragePath(absPath)

	res, err := ExecWrapper(db.conn, constants.SQLDeleteRepoByPath, cleanPath).Destruct()
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "store.DeleteByPath"))
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return result.Fail[int64](apperror.WrapSimple(rowsErr, "store.DeleteByPath.rows"))
	}

	return result.Ok(rows)
}

// DeleteByPath removes the repo row whose AbsolutePath matches.
// Returns the number of rows deleted (0 when no match).
func (db *DB) DeleteByPath(absPath string) (int64, error) {
	res := db.DeleteByPathWrap(absPath)
	if res.IsFailure() {
		return 0, res.AppError()
	}

	val, err := res.Unwrap()
	if err != nil {
		return 0, err
	}

	return val, nil
}

// DeleteBySlugWrap removes repo rows matching slug encapsulated in RowsAffectedWrap.
func (db *DB) DeleteBySlugWrap(slug string) RowsAffectedWrap {
	res, err := ExecWrapper(db.conn, constants.SQLDeleteRepoBySlug, slug).Destruct()
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "store.DeleteBySlug"))
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return result.Fail[int64](apperror.WrapSimple(rowsErr, "store.DeleteBySlug.rows"))
	}

	return result.Ok(rows)
}

// DeleteBySlug removes every repo row whose Slug matches.
// Returns the number of rows deleted.
func (db *DB) DeleteBySlug(slug string) (int64, error) {
	res := db.DeleteBySlugWrap(slug)
	if res.IsFailure() {
		return 0, res.AppError()
	}

	return res.Unwrap()
}

// scanRows reads ScanRecord values from query result rows.
func scanRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.ScanRecord, error) {
	var results []model.ScanRecord

	for rows.Next() {
		r, err := scanOneRow(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}

// scanOneRow reads a single ScanRecord from the current row.
func scanOneRow(row interface{ Scan(dest ...any) error }) (model.ScanRecord, error) {
	var r model.ScanRecord
	err := row.Scan(
		&r.ID, &r.Slug, &r.RepoName, &r.HTTPSUrl, &r.SSHUrl,
		&r.Branch, &r.RelativePath, &r.AbsolutePath,
		&r.CloneInstruction, &r.Notes, &r.IdentifiedTransport,
	)

	return r, err
}

// GetRepoSuggestions returns up to 10 repo slugs matching the partial string.
func (db *DB) GetRepoSuggestions(partial string) ([]string, error) {
	likePattern := "%" + partial + "%"
	rows, err := QueryWrapper(db.conn, constants.SQLSuggestRepoBySlug, likePattern, likePattern).Destruct()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var matches []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}

		matches = append(matches, slug)
	}

	return matches, nil
}

// UpdateRepoAbsolutePath modifies the tracked AbsolutePath of a repository record.
// Both oldPath and newPath are canonicalized using NormalizeStoragePath.
func (db *DB) UpdateRepoAbsolutePath(oldPath, newPath string) error {
	cleanOld := NormalizeStoragePath(oldPath)
	cleanNew := NormalizeStoragePath(newPath)
	if cleanOld == "" || cleanNew == "" || cleanOld == cleanNew {
		return nil
	}

	_, err := ExecWrapper(db.conn, constants.SQLUpdateRepoAbsolutePath, cleanNew, cleanOld).Destruct()

	return err
}

// GetRepoPrimaryURLByPath retrieves the primary remote URL and repository ID for a path.
func (db *DB) GetRepoPrimaryURLByPath(repoPath string) (int64, string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return 0, "", nil
	}

	rec, err := db.findFirstRepoRecord(repoPath)
	if err != nil {
		return 0, "", err
	}

	url := resolvePrimaryURLFromFields(rec.HTTPSUrl, rec.SSHUrl, rec.IdentifiedTransport)

	return rec.ID, url, nil
}

func (db *DB) findFirstRepoRecord(repoPath string) (model.ScanRecord, error) {
	records, err := db.FindByPath(repoPath)
	if err != nil {
		return model.ScanRecord{}, err
	}

	if len(records) == 0 {
		return model.ScanRecord{}, sql.ErrNoRows
	}

	return records[0], nil
}

func resolvePrimaryURLFromFields(httpsURL, sshURL, transport string) string {
	isSSH := strings.EqualFold(transport, constants.TransportSSH)
	if isSSH && sshURL != "" {
		return sshURL
	}

	if httpsURL != "" {
		return httpsURL
	}

	return sshURL
}
