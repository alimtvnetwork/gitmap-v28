package dbengine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// DbWrapper encapsulates a SQL database connection with dialect query compilation.
type DbWrapper struct {
	conn     *sql.DB
	dialect  DatabaseDialectType
	compiler DialectCompiler
}

// TxWrapper encapsulates an active SQL transaction.
type TxWrapper struct {
	tx       *sql.Tx
	compiler DialectCompiler
}

// OpenDb opens a database connection for the specified dialect.
func OpenDb(dialect DatabaseDialectType, dsn string) (*DbWrapper, *apperror.AppError) {
	conn, err := sql.Open(string(dialect), dsn)
	if err != nil {
		return nil, apperror.WrapSimple(err, fmt.Sprintf("open database for dialect %s", dialect))
	}
	return WrapDb(conn, dialect)
}

// WrapDb wraps an existing sql.DB connection with dialect compilation.
func WrapDb(conn *sql.DB, dialect DatabaseDialectType) (*DbWrapper, *apperror.AppError) {
	compiler, appErr := ResolveCompiler(dialect)
	if appErr != nil {
		return nil, appErr
	}
	return &DbWrapper{
		conn:     conn,
		dialect:  dialect,
		compiler: compiler,
	}, nil
}

// Close closes the underlying connection.
func (w *DbWrapper) Close() *apperror.AppError {
	if w.conn == nil {
		return nil
	}
	err := w.conn.Close()
	if err != nil {
		return apperror.WrapSimple(err, "close database connection")
	}
	return nil
}

// QueryRow executes a query that is expected to return at most one row.
func (w *DbWrapper) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, *apperror.AppError) {
	row := w.conn.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, apperror.WrapSimple(row.Err(), "execute query row: "+query)
	}
	return row, nil
}

// Query executes a query that returns rows.
func (w *DbWrapper) Query(ctx context.Context, query string, args ...any) (*sql.Rows, *apperror.AppError) {
	rows, err := w.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "execute query: "+query)
	}
	return rows, nil
}

// Exec executes a query without returning rows.
func (w *DbWrapper) Exec(ctx context.Context, query string, args ...any) (sql.Result, *apperror.AppError) {
	result, err := w.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "execute exec: "+query)
	}
	return result, nil
}

// WithTransaction runs a function within a database transaction.
func (w *DbWrapper) WithTransaction(ctx context.Context, fn func(tx *TxWrapper) *apperror.AppError) *apperror.AppError {
	tx, err := w.conn.BeginTx(ctx, nil)
	if err != nil {
		return apperror.WrapSimple(err, "begin transaction")
	}

	txWrap := &TxWrapper{tx: tx, compiler: w.compiler}
	txErr := fn(txWrap)
	if txErr != nil {
		_ = tx.Rollback()
		return txErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return apperror.WrapSimple(commitErr, "commit transaction")
	}
	return nil
}

// Compiler returns the active DialectCompiler.
func (w *DbWrapper) Compiler() DialectCompiler {
	return w.compiler
}

// Dialect returns the active DatabaseDialectType.
func (w *DbWrapper) Dialect() DatabaseDialectType {
	return w.dialect
}

// Conn returns the underlying *sql.DB connection.
func (w *DbWrapper) Conn() *sql.DB {
	return w.conn
}

// CreateView creates a database view with the specified SELECT statement.
func (w *DbWrapper) CreateView(ctx context.Context, name string, selectSql string) BoolResult {
	query := w.compiler.CompileCreateView(name, selectSql)
	_, err := w.Exec(ctx, query)
	if err != nil {
		return FailureBool(err)
	}
	return SuccessBool(true)
}

// DropView drops a database view.
func (w *DbWrapper) DropView(ctx context.Context, name string) BoolResult {
	query := w.compiler.CompileDropView(name)
	_, err := w.Exec(ctx, query)
	if err != nil {
		return FailureBool(err)
	}
	return SuccessBool(true)
}

// CallFunction executes a database function and returns the scalar string result.
func (w *DbWrapper) CallFunction(ctx context.Context, name string, args ...any) StringResult {
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
func (w *DbWrapper) ExecRowsAffected(ctx context.Context, query string, args ...any) RowsAffectedResult {
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
func (w *DbWrapper) ViewExists(ctx context.Context, name string) (bool, *apperror.AppError) {
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
func (w *DbWrapper) GetTableColumns(ctx context.Context, name string) ([]string, *apperror.AppError) {
	query := w.compiler.CompileInspectColumns(name)
	if len(query) == 0 {
		return nil, nil
	}

	if w.dialect == DbSQLite {
		return w.scanSqliteColumns(ctx, query)
	}

	return w.scanStandardColumns(ctx, query, name)
}

func (w *DbWrapper) scanSqliteColumns(ctx context.Context, query string) ([]string, *apperror.AppError) {
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

func (w *DbWrapper) scanStandardColumns(ctx context.Context, query string, name string) ([]string, *apperror.AppError) {
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

func (w *DbWrapper) verifyViewColumns(ctx context.Context, name string, requiredColumns []string) (bool, *apperror.AppError) {
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

// CreateViewOrUseView inspects if a view exists and contains the required columns.
// If valid, it reuses the view. If missing or columns differ, it creates or recreates the view.
func (w *DbWrapper) CreateViewOrUseView(ctx context.Context, name string, selectSql string, requiredColumns ...string) BoolResult {
	exists, appErr := w.ViewExists(ctx, name)
	if appErr != nil {
		return FailureBool(appErr)
	}

	if !exists {
		return w.CreateView(ctx, name, selectSql)
	}

	if len(requiredColumns) == 0 {
		return SuccessBool(true)
	}

	hasAll, verifyErr := w.verifyViewColumns(ctx, name, requiredColumns)
	if verifyErr != nil {
		return FailureBool(verifyErr)
	}

	if hasAll {
		return SuccessBool(true)
	}

	dropRes := w.DropView(ctx, name)
	if dropRes.IsFailed() {
		return dropRes
	}

	return w.CreateView(ctx, name, selectSql)
}

