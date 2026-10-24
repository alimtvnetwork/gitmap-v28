package dbengine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DbWrapper encapsulates a SQL database connection with dialect query compilation.
type DBWrapper struct {
	conn     *sql.DB
	dialect  DatabaseDialectType
	compiler DialectCompiler
}

// TxWrapper encapsulates an active SQL transaction.
type TxWrapper struct {
	tx       *sql.Tx
	compiler DialectCompiler
}

var sqlitePragmas = []string{
	"PRAGMA busy_timeout = 5000",
	"PRAGMA journal_mode = WAL",
	"PRAGMA synchronous = NORMAL",
	"PRAGMA foreign_keys = ON",
}

func applySQLitePragmas(conn *sql.DB) *apperror.AppError {
	for _, pragma := range sqlitePragmas {
		if _, err := conn.Exec(pragma); err != nil {
			_ = conn.Close()

			return apperror.WrapSimple(err, "configure sqlite pragma: "+pragma)
		}
	}

	return nil
}

func configureConnForDialect(dialect DatabaseDialectType, conn *sql.DB) *apperror.AppError {
	if dialect != DBSQLite {
		return nil
	}

	conn.SetMaxOpenConns(1)

	return applySQLitePragmas(conn)
}

// OpenDb opens a database connection for the specified dialect.
func OpenDB(dialect DatabaseDialectType, dsn string) (*DBWrapper, *apperror.AppError) {
	conn, err := sql.Open(string(dialect), dsn)
	if err != nil {
		return nil, apperror.WrapSimple(err, fmt.Sprintf("open database for dialect %s", dialect))
	}

	if appErr := configureConnForDialect(dialect, conn); appErr != nil {
		return nil, appErr
	}

	return WrapDB(conn, dialect)
}

// WrapDb wraps an existing sql.DB connection with dialect compilation.
func WrapDB(conn *sql.DB, dialect DatabaseDialectType) (*DBWrapper, *apperror.AppError) {
	compiler, appErr := ResolveCompiler(dialect)
	if appErr != nil {
		return nil, appErr
	}

	return &DBWrapper{
		conn:     conn,
		dialect:  dialect,
		compiler: compiler,
	}, nil
}

// Close closes the underlying connection.
func (w *DBWrapper) Close() *apperror.AppError {
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
func (w *DBWrapper) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, *apperror.AppError) {
	row := w.conn.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, apperror.WrapSimple(row.Err(), "execute query row: "+query)
	}

	return row, nil
}

// Query executes a query that returns rows.
func (w *DBWrapper) Query(ctx context.Context, query string, args ...any) (*sql.Rows, *apperror.AppError) {
	rows, err := w.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "execute query: "+query)
	}

	return rows, nil
}

// Exec executes a query without returning rows.
func (w *DBWrapper) Exec(ctx context.Context, query string, args ...any) (sql.Result, *apperror.AppError) {
	result, err := w.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "execute exec: "+query)
	}

	return result, nil
}

// Prepare creates a prepared statement on the database connection.
func (w *DBWrapper) Prepare(ctx context.Context, query string) (*sql.Stmt, *apperror.AppError) {
	stmt, err := w.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "prepare statement: "+query)
	}

	return stmt, nil
}

func handleRollback(tx *sql.Tx, txErr *apperror.AppError) *apperror.AppError {
	rbErr := tx.Rollback()
	if rbErr == nil {
		return txErr
	}

	msg := "rollback failed after: " + txErr.Error()

	return apperror.WrapWithDetails(
		rbErr, "rollback transaction", "E9000", msg,
		"dbengine", apperror.ErrorTypeExecution, apperror.SeverityError, nil,
	)
}

var (
	immediateTxOptions = &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}

	readOnlyTxOptions = &sql.TxOptions{
		ReadOnly: true,
	}

	exclusiveTxOptions = &sql.TxOptions{
		Isolation: sql.LevelLinearizable,
	}
)

func recoverTxPanic(tx *sql.Tx) {
	p := recover()
	if p == nil {
		return
	}

	rbErr := tx.Rollback()
	if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
		panic(fmt.Sprintf("%v (rollback failed: %v)", p, rbErr)) // lint-allow: panic
	}

	panic(p) // lint-allow: panic
}

func (w *DBWrapper) executeTx(
	ctx context.Context,
	opts *sql.TxOptions,
	fn func(tx *TxWrapper) *apperror.AppError,
) *apperror.AppError {
	tx, err := w.conn.BeginTx(ctx, opts)
	if err != nil {
		return apperror.WrapSimple(err, "begin transaction")
	}

	defer recoverTxPanic(tx)

	return w.runTxFunc(tx, fn)
}

func (w *DBWrapper) runTxFunc(tx *sql.Tx, fn func(tx *TxWrapper) *apperror.AppError) *apperror.AppError {
	txWrap := &TxWrapper{tx: tx, compiler: w.compiler}
	if txErr := fn(txWrap); txErr != nil {
		return handleRollback(tx, txErr)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return apperror.WrapSimple(commitErr, "commit transaction")
	}

	return nil
}

// WithTransaction runs a function within a database transaction.
func (w *DBWrapper) WithTransaction(ctx context.Context, fn func(tx *TxWrapper) *apperror.AppError) *apperror.AppError {
	return w.executeTx(ctx, nil, fn)
}

// WithImmediateTransaction runs a function within an immediate (serializable) database transaction.
func (w *DBWrapper) WithImmediateTransaction(ctx context.Context, fn func(tx *TxWrapper) *apperror.AppError) *apperror.AppError {
	return w.executeTx(ctx, immediateTxOptions, fn)
}

// WithTxOptions runs a function within a database transaction using the provided options.
func (w *DBWrapper) WithTxOptions(
	ctx context.Context,
	opts *sql.TxOptions,
	fn func(tx *TxWrapper) *apperror.AppError,
) *apperror.AppError {
	return w.executeTx(ctx, opts, fn)
}

// WithReadOnlyTransaction runs a function within a read-only database transaction.
func (w *DBWrapper) WithReadOnlyTransaction(
	ctx context.Context,
	fn func(tx *TxWrapper) *apperror.AppError,
) *apperror.AppError {
	return w.executeTx(ctx, readOnlyTxOptions, fn)
}

// WithExclusiveTransaction runs a function within an exclusive database transaction.
func (w *DBWrapper) WithExclusiveTransaction(
	ctx context.Context,
	fn func(tx *TxWrapper) *apperror.AppError,
) *apperror.AppError {
	return w.executeTx(ctx, exclusiveTxOptions, fn)
}

// Compiler returns the active DialectCompiler.
func (w *DBWrapper) Compiler() DialectCompiler {
	return w.compiler
}

// Dialect returns the active DatabaseDialectType.
func (w *DBWrapper) Dialect() DatabaseDialectType {
	return w.dialect
}

// Conn returns the underlying *sql.DB connection.
func (w *DBWrapper) Conn() *sql.DB {
	return w.conn
}

// CreateView creates a database view with the specified SELECT statement.
func (w *DBWrapper) CreateView(ctx context.Context, name string, selectSql string) BoolResult {
	query := w.compiler.CompileCreateView(name, selectSql)
	_, err := w.Exec(ctx, query)
	if err != nil {
		return FailureBool(err)
	}

	return SuccessBool(true)
}
