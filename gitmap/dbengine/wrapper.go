package dbengine

import (
	"context"
	"database/sql"
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
