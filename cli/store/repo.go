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
	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
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

// ListReposWrap returns all tracked repositories encapsulated in ScanRecordSliceWrap.
func (db *DB) ListReposWrap() ScanRecordSliceWrap {
	return NewRepoRepository(db.conn).ListAll()
}

// ListRepos returns all tracked repositories ordered by slug.
func (db *DB) ListRepos() ([]model.ScanRecord, error) {
	res := db.ListReposWrap()
	if res.IsFailure() {
		return nil, res.AppError()
	}

	return res.Items(), nil
}

const sqlSelectReposByTransport = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE IdentifiedTransport = ? COLLATE NOCASE ORDER BY Slug ASC"

// ListReposByTransport returns all repositories configured for the specified transport protocol.
func (db *DB) ListReposByTransport(transport string) ([]model.ScanRecord, error) {
	clean := strings.TrimSpace(transport)
	if clean == "" {
		return db.ListRepos()
	}

	rows, err := QueryWrapper(db.conn, sqlSelectReposByTransport, clean).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query repos")
	}

	defer rows.Close()

	return scanRows(rows)
}

// ListReposByBranch returns all repositories configured for the specified branch.
func (db *DB) ListReposByBranch(branch string) ScanRecordSliceWrap {
	return NewRepoRepository(db.conn).ListByBranch(branch)
}

// ListReposDirtyStatus returns repositories flagged as dirty or requiring synchronization.
func (db *DB) ListReposDirtyStatus() ScanRecordSliceWrap {
	return NewRepoRepository(db.conn).ListDirty()
}

// CountReposWrap returns the total number of repositories encapsulated in RepoCountWrap.
func (db *DB) CountReposWrap() RepoCountWrap {
	var total int

	err := db.conn.QueryRow("SELECT COUNT(*) FROM Repo").Scan(&total)
	if err != nil {
		return result.Fail[int](apperror.WrapSimple(err, "store.CountRepos"))
	}

	return result.Ok(total)
}

// CountRepos returns the total number of repositories tracked in the database.
func (db *DB) CountRepos() (int, error) {
	res := db.CountReposWrap()
	if res.IsFailure() {
		return 0, res.AppError()
	}

	return res.Unwrap()
}

// FindByID returns all repos matching the given ID.
func (db *DB) FindByID(id int64) ([]model.ScanRecord, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectRepoByID, id).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query repos")
	}

	defer rows.Close()

	return scanRows(rows)
}

func findBySlugRunner(runner sqlQueryer, slug string) ([]model.ScanRecord, error) {
	rows, err := QueryWrapper(runner, constants.SQLSelectRepoBySlug, slug).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query repos")
	}

	defer rows.Close()

	return scanRows(rows)
}

func isMissingError(err *apperror.AppError) bool {
	if err == nil {
		return false
	}

	return err.Code == "E_NOT_FOUND" || strings.Contains(err.Error(), "not found")
}

// FindBySlugWrap returns repository matching slug encapsulated in ScanRecordWrap.
func (db *DB) FindBySlugWrap(slug string) ScanRecordWrap {
	return NewRepoRepository(db.conn).FindBySlug(slug)
}

// FindBySlug returns all repos matching the given slug.
func (db *DB) FindBySlug(slug string) ([]model.ScanRecord, error) {
	res := db.FindBySlugWrap(slug)
	if res.IsSuccess() {
		rec, _ := res.Unwrap()

		return []model.ScanRecord{rec}, nil
	}

	if isMissingError(res.AppError()) {
		return nil, nil
	}

	return nil, res.AppError()
}

// FindOneBySlugWrap returns the first repository matching slug encapsulated in ScanRecordWrap.
func (db *DB) FindOneBySlugWrap(slug string) ScanRecordWrap {
	return NewRepoRepository(db.conn).FindBySlug(slug)
}

// FindOneBySlug returns the first repository matching the given slug or nil if not found.
func (db *DB) FindOneBySlug(slug string) (*model.ScanRecord, error) {
	res := db.FindOneBySlugWrap(slug)
	if res.IsSuccess() {
		rec, _ := res.Unwrap()

		return &rec, nil
	}

	if isMissingError(res.AppError()) {
		return nil, nil
	}

	return nil, res.AppError()
}

// FindByPathWrap returns repositories matching absolute path encapsulated in ScanRecordSliceWrap.
func (db *DB) FindByPathWrap(absPath string) ScanRecordSliceWrap {
	cleanPath := NormalizeStoragePath(absPath)

	rows, err := QueryWrapper(db.conn, constants.SQLSelectRepoByPath, cleanPath).Destruct()
	if err != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "store.FindByPath"))
	}

	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(scanErr, "store.FindByPath.scan"))
	}

	return result.OkSlice(records)
}

// FindByPath returns the repo at the given absolute path.
func (db *DB) FindByPath(absPath string) ([]model.ScanRecord, error) {
	res := db.FindByPathWrap(absPath)
	if res.IsFailure() {
		return nil, res.AppError()
	}

	return res.Items(), nil
}

// RepoExists returns true if a repository with the given RepoId exists.
func (db *DB) RepoExists(repoID int64) bool {
	if repoID <= 0 {
		return false
	}

	var isPresent bool
	_ = db.conn.QueryRow("SELECT EXISTS(SELECT 1 FROM Repo WHERE RepoId = ?)", repoID).Scan(&isPresent)

	return isPresent
}

// SelectRepoIDByPath returns the persisted RepoId for a repository absolute path.
func (db *DB) SelectRepoIDByPath(absPath string) (int64, error) {
	cleanPath := NormalizeStoragePath(absPath)
	if id := db.queryRepoIDExact(cleanPath); id > 0 {
		return id, nil
	}

	return db.queryRepoIDPrefix(cleanPath)
}

func (db *DB) queryRepoIDExact(cleanPath string) int64 {
	var id int64
	err := db.conn.QueryRow(
		"SELECT RepoId FROM Repo WHERE REPLACE(AbsolutePath, '\\', '/') = ? COLLATE NOCASE",
		cleanPath,
	).Scan(&id)
	if err == nil && id > 0 {
		return id
	}

	return 0
}

func (db *DB) queryRepoIDPrefix(cleanPath string) (int64, error) {
	var id int64
	err := db.conn.QueryRow(
		"SELECT RepoId FROM Repo WHERE ? LIKE (REPLACE(AbsolutePath, '\\', '/') || '%') ORDER BY LENGTH(AbsolutePath) DESC LIMIT 1",
		cleanPath,
	).Scan(&id)
	if err == nil && id > 0 {
		return id, nil
	}

	return 0, err
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
