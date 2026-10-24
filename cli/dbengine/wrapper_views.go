package dbengine

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DropView drops a database view.
func (w *DBWrapper) DropView(ctx context.Context, name string) BoolResult {
	query := w.compiler.CompileDropView(name)
	_, err := w.Exec(ctx, query)
	if err != nil {
		return FailureBool(err)
	}

	return SuccessBool(true)
}

// CallFunction executes a database function and returns the scalar string result.
func (w *DBWrapper) CallFunction(ctx context.Context, name string, args ...any) StringResult {
	query := w.compiler.CompileFunctionCall(name, len(args))
	row, appErr := w.QueryRow(ctx, query, args...)
	if appErr != nil {
		return FailureString(appErr)
	}

	var res string
	if err := row.Scan(&res); err != nil {
		return FailureString(apperror.WrapSimple(err, "scan result of function "+name))
	}

	return SuccessString(res)
}

// ExecRowsAffected executes a statement and returns the number of rows affected wrapped in RowsAffectedResult.
func (w *DBWrapper) ExecRowsAffected(ctx context.Context, query string, args ...any) RowsAffectedResult {
	res, appErr := w.Exec(ctx, query, args...)
	if appErr != nil {
		return FailureRowsAffected(appErr)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return FailureRowsAffected(apperror.WrapSimple(err, "get rows affected"))
	}

	return SuccessRowsAffected(affected)
}

// ViewExists checks whether a database view exists.
func (w *DBWrapper) ViewExists(ctx context.Context, name string) (bool, *apperror.AppError) {
	query := w.compiler.CompileInspectViewExists(name)
	if len(query) == 0 {
		return false, nil
	}

	row, appErr := w.QueryRow(ctx, query, name)
	if appErr != nil {
		return false, appErr
	}

	var dummy int
	err := row.Scan(&dummy)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return false, apperror.WrapSimple(err, "check view existence "+name)
}

// GetTableColumns returns column names for a table or view.
func (w *DBWrapper) GetTableColumns(ctx context.Context, name string) ([]string, *apperror.AppError) {
	query := w.compiler.CompileInspectColumns(name)
	if len(query) == 0 {
		return nil, nil
	}

	if w.dialect == DBSQLite {
		return w.scanSqliteColumns(ctx, query)
	}

	return w.scanStandardColumns(ctx, query, name)
}

func (w *DBWrapper) scanSqliteColumns(ctx context.Context, query string) ([]string, *apperror.AppError) {
	rows, appErr := w.Query(ctx, query)
	if appErr != nil {
		return nil, appErr
	}

	defer rows.Close()

	var cols []string
	for rows.Next() {
		var (
			cid     int
			colName string
			ctype   string
			notnull int
			dflt    any
			pk      int
		)
		if err := rows.Scan(&cid, &colName, &ctype, &notnull, &dflt, &pk); err != nil {
			return nil, apperror.WrapSimple(err, "scan sqlite column info")
		}

		cols = append(cols, colName)
	}

	return cols, nil
}

func (w *DBWrapper) scanStandardColumns(ctx context.Context, query string, name string) ([]string, *apperror.AppError) {
	rows, appErr := w.Query(ctx, query, name)
	if appErr != nil {
		return nil, appErr
	}

	defer rows.Close()

	var cols []string
	for rows.Next() {
		var colName string
		if err := rows.Scan(&colName); err != nil {
			return nil, apperror.WrapSimple(err, "scan column info for "+name)
		}

		cols = append(cols, colName)
	}

	return cols, nil
}

func (w *DBWrapper) verifyViewColumns(ctx context.Context, name string, requiredColumns []string) (bool, *apperror.AppError) {
	existingCols, appErr := w.GetTableColumns(ctx, name)
	if appErr != nil {
		return false, appErr
	}

	colMap := make(map[string]bool, len(existingCols))
	for _, col := range existingCols {
		colMap[col] = true
	}

	for _, req := range requiredColumns {
		if !colMap[req] {
			return false, nil
		}
	}

	return true, nil
}

const sqlCreateViewMeta = `
CREATE TABLE IF NOT EXISTS __gitmap_view_meta (
    ViewName TEXT PRIMARY KEY,
    QueryHash TEXT NOT NULL,
    ViewSql TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);`

// ComputeSqlHash computes a deterministic SHA-256 hex string for a given SQL statement.
func ComputeSqlHash(sqlStr string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(sqlStr)))

	return hex.EncodeToString(sum[:])
}

// ValidateSql performs a dry-run syntax and reference check using EXPLAIN on the database connection.
func (w *DBWrapper) ValidateSql(ctx context.Context, sqlStr string) *apperror.AppError {
	cleanSql := strings.TrimRight(strings.TrimSpace(sqlStr), ";")
	explainQuery := fmt.Sprintf("EXPLAIN %s", cleanSql)
	rows, err := w.conn.QueryContext(ctx, explainQuery)
	if err != nil {
		return apperror.WrapSimple(err, "validate sql query syntax")
	}

	defer rows.Close()

	return nil
}

// EnsureViewMetaTable initializes the view metadata table if it does not already exist.
func (w *DBWrapper) EnsureViewMetaTable(ctx context.Context) *apperror.AppError {
	_, err := w.Exec(ctx, sqlCreateViewMeta)
	if err != nil {
		return apperror.WrapSimple(err, "ensure view metadata table")
	}

	return nil
}

// GetViewHash retrieves the recorded query hash for a view from __gitmap_view_meta.
func (w *DBWrapper) GetViewHash(ctx context.Context, name string) (string, *apperror.AppError) {
	query := "SELECT QueryHash FROM __gitmap_view_meta WHERE ViewName = ?"
	row, appErr := w.QueryRow(ctx, query, name)
	if appErr != nil {
		return "", appErr
	}

	var hash string
	err := row.Scan(&hash)
	if err == nil {
		return hash, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	return "", apperror.WrapSimple(err, "scan view hash for "+name)
}

// SaveViewMeta records or updates the view hash and SQL in __gitmap_view_meta.
func (w *DBWrapper) SaveViewMeta(ctx context.Context, name, queryHash, viewSql string) *apperror.AppError {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `
INSERT OR REPLACE INTO __gitmap_view_meta (ViewName, QueryHash, ViewSql, UpdatedAt)
VALUES (?, ?, ?, ?);`
	_, appErr := w.Exec(ctx, query, name, queryHash, viewSql, now)
	if appErr != nil {
		return apperror.WrapSimple(appErr, "save view metadata for "+name)
	}

	return nil
}

func (w *DBWrapper) isViewCurrent(ctx context.Context, name, queryHash string) (bool, *apperror.AppError) {
	existingHash, appErr := w.GetViewHash(ctx, name)
	if appErr != nil {
		return false, appErr
	}

	if len(existingHash) == 0 {
		return false, nil
	}

	return existingHash == queryHash, nil
}
