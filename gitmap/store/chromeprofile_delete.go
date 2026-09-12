// Package store — chromeprofile_delete.go: removes a ChromeProfile +
// its cascaded ChromeProfileExport rows. Returns the artifact file
// paths that were tracked so the CLI can rm() them on disk.
package store

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

const sqlSelectChromeProfileExports = `
SELECT e.FilePath FROM ChromeProfileExport e
JOIN ChromeProfile p ON p.ChromeProfileId = e.ChromeProfileId
WHERE p.Name = ?`

const sqlDeleteChromeProfile = `DELETE FROM ChromeProfile WHERE Name = ?`

const sqlDeleteChromeProfileExports = `
DELETE FROM ChromeProfileExport
WHERE ChromeProfileId IN (SELECT ChromeProfileId FROM ChromeProfile WHERE Name = ?)`

// DeleteChromeProfile removes the named profile and its artifact rows.
// Returns the list of artifact file paths so the caller can clean disk.
func (db *DB) DeleteChromeProfile(name string) ([]string, error) {
	if err := db.EnsureChromeProfileTables(); err != nil {
		return nil, err
	}

	paths, err := db.collectChromeArtifactPaths(name)
	if err != nil {
		return nil, err
	}

	if delErr := db.deleteChromeProfileTx(name); delErr != nil {
		return nil, delErr
	}

	return paths, nil
}

func (db *DB) deleteChromeProfileTx(name string) error {
	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	ctx := context.Background()
	appErr = wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeDeleteChromeProfileTx(ctx, tx, name)
	})
	if appErr != nil {
		return appErr
	}

	return nil
}

func executeDeleteChromeProfileTx(ctx context.Context, tx *dbengine.TxWrapper, name string) *apperror.AppError {
	if _, appErr := tx.Exec(ctx, sqlDeleteChromeProfileExports, name); appErr != nil {
		return apperror.WrapSimple(appErr, "delete chrome-profile exports")
	}

	if _, appErr := tx.Exec(ctx, sqlDeleteChromeProfile, name); appErr != nil {
		return apperror.WrapSimple(appErr, "delete chrome-profile")
	}

	return nil
}

// ChromeProfileExists reports whether a row with the given name exists.
func (db *DB) ChromeProfileExists(name string) bool {
	if err := db.EnsureChromeProfileTables(); err != nil {
		return false
	}

	var id int64

	return db.conn.QueryRow(sqlSelectChromeProfileId, name).Scan(&id) == nil
}

// collectChromeArtifactPaths reads the FilePath column for every export
// row tied to the named profile.
func (db *DB) collectChromeArtifactPaths(name string) ([]string, error) {
	rows, err := QueryWrapper(db.conn, sqlSelectChromeProfileExports, name).Destruct()
	if err != nil {
		return nil, fmt.Errorf("query chrome-profile artifacts: %w", err)
	}

	defer rows.Close()
	var out []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("scan artifact: %w", err)
		}

		out = append(out, p)
	}

	return out, nil
}
