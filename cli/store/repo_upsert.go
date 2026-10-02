package store

import (
	"context"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// NormalizeRepoRecord ensures a record's AbsolutePath is canonicalized before storage.
func NormalizeRepoRecord(r *model.ScanRecord) {
	if r != nil {
		r.AbsolutePath = NormalizeStoragePath(r.AbsolutePath)
	}
}

// IsPathDuplicateCaseInsensitive checks if two paths are equivalent, respecting OS case sensitivity.
func IsPathDuplicateCaseInsensitive(pathA, pathB string) bool {
	normA := NormalizeStoragePath(pathA)
	normB := NormalizeStoragePath(pathB)

	isCaseInsensitive := runtime.GOOS == "windows" || fsutil.IsPathCaseInsensitive(pathA) || fsutil.IsPathCaseInsensitive(pathB)
	if isCaseInsensitive {
		return strings.EqualFold(normA, normB)
	}

	return normA == normB
}

func buildRepoUniqueQuery(normPath string) string {
	query := "SELECT COUNT(*) FROM Repo WHERE AbsolutePath = ?"
	isCaseInsensitive := runtime.GOOS == "windows" || fsutil.IsPathCaseInsensitive(normPath)
	if isCaseInsensitive {
		query = "SELECT COUNT(*) FROM Repo WHERE AbsolutePath = ? COLLATE NOCASE"
	}

	return query
}

func validateDBConn(db *DB) *apperror.AppError {
	if db == nil || db.conn == nil {
		return apperror.New("EnsureRepoUniqueInDB", "E_DB_NIL", map[string]any{
			"error": "database connection is nil",
		})
	}

	return nil
}

// EnsureRepoUniqueInDB queries the Repo table in SQLite by NormalizeStoragePath(absPath)
// (using case folding on Windows, exact matching on Unix) to confirm uniqueness before inserting.
// Returns true if unique (no duplicate found), or false if an existing record already exists.
func EnsureRepoUniqueInDB(db *DB, absPath string) (bool, error) {
	if appErr := validateDBConn(db); appErr != nil {
		return false, appErr
	}

	normPath := NormalizeStoragePath(absPath)
	query := buildRepoUniqueQuery(normPath)
	var count int
	if err := db.conn.QueryRow(query, normPath).Scan(&count); err != nil {
		return false, apperror.WrapSimple(err, "EnsureRepoUniqueInDB")
	}

	isUnique := count == 0

	return isUnique, nil
}

// UpsertOneRepoRecord inserts or updates a single repository record ensuring path normalization.
func (db *DB) UpsertOneRepoRecord(ctx context.Context, r model.ScanRecord) error {
	NormalizeRepoRecord(&r)

	return db.UpsertRepos([]model.ScanRecord{r})
}
