package store

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

const sqlSelectReposByTransport = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE IdentifiedTransport = ? COLLATE NOCASE ORDER BY Slug ASC"

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
