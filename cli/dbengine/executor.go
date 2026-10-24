package dbengine

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SqlExecuter provides a unified interface for executing queries across database connections and transactions.
type SqlExecuter interface {
	QueryRow(ctx context.Context, query string, args ...any) DBRowWrap
	Query(ctx context.Context, query string, args ...any) DBRowsWrap
	Exec(ctx context.Context, query string, args ...any) DBResultWrap
	ExecRowsAffected(ctx context.Context, query string, args ...any) RowsAffectedResult
	Compiler() DialectCompiler
}

// SqlExecutor is an alias to SqlExecuter for conventional spelling.
type SqlExecutor = SqlExecuter

var (
	_ SqlExecutor = (*DBWrapper)(nil)
	_ SqlExecutor = (*TxWrapper)(nil)
)

// Tx returns the underlying *sql.Tx transaction.
func (t *TxWrapper) Tx() *sql.Tx {
	return t.tx
}

// Compiler returns the active DialectCompiler.
func (t *TxWrapper) Compiler() DialectCompiler {
	return t.compiler
}

// QueryRow executes a query within the transaction expecting at most one row.
func (t *TxWrapper) QueryRow(ctx context.Context, query string, args ...any) DBRowWrap {
	row := t.tx.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return FailureDBRow(apperror.WrapSimple(row.Err(), "execute tx query row: "+query))
	}

	return SuccessDBRow(row)
}

// Query executes a query within the transaction returning rows.
func (t *TxWrapper) Query(ctx context.Context, query string, args ...any) DBRowsWrap {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return FailureDBRows(apperror.WrapSimple(err, "execute tx query: "+query))
	}

	return SuccessDBRows(rows)
}

// Exec executes a query within the transaction without returning rows.
func (t *TxWrapper) Exec(ctx context.Context, query string, args ...any) DBResultWrap {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return FailureDBResult(apperror.WrapSimple(err, "execute tx exec: "+query))
	}

	return SuccessDBResult(result)
}

// ExecRowsAffected executes a statement within the transaction and returns rows affected.
func (t *TxWrapper) ExecRowsAffected(ctx context.Context, query string, args ...any) RowsAffectedResult {
	execRes := t.Exec(ctx, query, args...)
	if execRes.IsFailure() {
		return FailureRowsAffected(execRes.AppError())
	}

	affected, err := execRes.Value.RowsAffected()
	if err != nil {
		return FailureRowsAffected(apperror.WrapSimple(err, "get tx rows affected"))
	}

	return SuccessRowsAffected(affected)
}

// Prepare creates a prepared statement within the transaction.
func (t *TxWrapper) Prepare(ctx context.Context, query string) DBStmtWrap {
	stmt, err := t.tx.PrepareContext(ctx, query)
	if err != nil {
		return FailureDBStmt(apperror.WrapSimple(err, "prepare tx statement: "+query))
	}

	return SuccessDBStmt(stmt)
}

// ValidateSql performs a dry-run syntax and reference check using EXPLAIN within the transaction.
func (t *TxWrapper) ValidateSql(ctx context.Context, sqlStr string) *apperror.AppError {
	cleanSql := strings.TrimRight(strings.TrimSpace(sqlStr), ";")
	explainQuery := fmt.Sprintf("EXPLAIN %s", cleanSql)
	rowsRes := t.Query(ctx, explainQuery)
	if rowsRes.IsFailure() {
		return rowsRes.AppError()
	}

	defer rowsRes.Value.Close()

	return nil
}

// CallFunction executes a database function within the transaction and returns the scalar string result.
func (t *TxWrapper) CallFunction(ctx context.Context, name string, args ...any) StringResult {
	query := t.compiler.CompileFunctionCall(name, len(args))
	rowRes := t.QueryRow(ctx, query, args...)
	if rowRes.IsFailure() {
		return FailureString(rowRes.AppError())
	}

	var res string
	if err := rowRes.Value.Scan(&res); err != nil {
		return FailureString(apperror.WrapSimple(err, "scan result of tx function "+name))
	}

	return SuccessString(res)
}
