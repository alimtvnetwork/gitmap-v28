package store

import (
	"database/sql"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// Spec-number issuance DDL (spec 252). UNIQUE(number) is the race-safety
// anchor: INSERT OR IGNORE either wins the number or reports zero rows
// affected, which the caller reads as "already taken".
const sqlCreateSpecNumbersTable = `CREATE TABLE IF NOT EXISTS spec_numbers (
  number     INTEGER NOT NULL,
  repo_root  TEXT    NOT NULL DEFAULT '',
  task_slug  TEXT    NOT NULL DEFAULT '',
  issued_at  TEXT    NOT NULL,
  UNIQUE(number)
);`

// EnsureSpecNumbersSchema creates the spec_numbers table when missing.
// Idempotent; runs at every DB open.
func EnsureSpecNumbersSchema(conn *sql.DB) *appfault.AppError {
	for _, ddl := range []string{sqlCreateSpecNumbersTable} {
		res := ExecWrapper(conn, ddl)
		if res.IsFailure {
			return appfault.WrapSimple(res.Error, "EnsureSpecNumbersSchema")
		}
	}

	return nil
}

// ClaimSpecNumber inserts (number, repoRoot, taskSlug, now-UTC-RFC3339)
// with INSERT OR IGNORE. Returns (true, nil) when this caller won the
// claim, (false, nil) when the number was already taken, and an error
// only on a genuine DB failure. RowsAffected == 0 is the normal
// "already taken" signal and must never be treated as an error.
func ClaimSpecNumber(conn *sql.DB, number int, repoRoot, taskSlug string) (bool, error) {
	const insertSQL = `INSERT OR IGNORE INTO spec_numbers (number, repo_root, task_slug, issued_at)
  VALUES (?, ?, ?, ?);`

	res := ExecWrapper(conn, insertSQL, number, repoRoot, taskSlug, time.Now().UTC().Format(time.RFC3339))
	if res.IsFailure {
		return false, res.Error
	}

	affected, rowsErr := res.Data.RowsAffected()
	if rowsErr != nil {
		return false, rowsErr
	}

	claimed := affected == 1

	return claimed, nil
}

// MaxIssuedSpecNumber returns the highest spec number ever issued through
// this DB, or 0 when the table is empty (fresh repo or pre-issuance state).
func MaxIssuedSpecNumber(conn *sql.DB) (int, error) {
	const maxSQL = `SELECT COALESCE(MAX(number), 0) FROM spec_numbers;`

	var maxNumber int
	scanErr := QueryRowWrapper(conn, maxSQL).Scan(&maxNumber)
	if scanErr != nil {
		return 0, scanErr
	}

	return maxNumber, nil
}
