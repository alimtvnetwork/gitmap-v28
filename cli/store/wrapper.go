package store

import (
	"database/sql"
	"fmt"
	"os"
)

// QueryResult wraps the result of a database query, providing structured error states.
type QueryResult[T any] struct {
	IsSuccess bool
	IsFailure bool
	Data      T
	Error     error
}

type sqlExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
}

type sqlQueryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type sqlRowQueryer interface {
	QueryRow(query string, args ...any) *sql.Row
}

type sqlQueryerExecutor interface {
	sqlQueryer
	sqlExecutor
}

// ExecWrapper wraps db.Exec, explicitly logging failures to os.Stderr.
func ExecWrapper(db sqlExecutor, query string, args ...any) QueryResult[sql.Result] {
	res, err := db.Exec(query, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[QueryWrapper Error]: exec failed: %v\nquery: %s\n", err, query)

		return QueryResult[sql.Result]{IsSuccess: false, IsFailure: true, Error: err}
	}

	return QueryResult[sql.Result]{IsSuccess: true, IsFailure: false, Data: res}
}

// QueryWrapper wraps db.Query, explicitly logging failures to os.Stderr.
func QueryWrapper(db sqlQueryer, query string, args ...any) QueryResult[*sql.Rows] {
	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[QueryWrapper Error]: query failed: %v\nquery: %s\n", err, query)

		return QueryResult[*sql.Rows]{IsSuccess: false, IsFailure: true, Error: err}
	}

	return QueryResult[*sql.Rows]{IsSuccess: true, IsFailure: false, Data: rows}
}

// QueryRowWrapper delegates to QueryRow. The error is deferred until Scan is called.
func QueryRowWrapper(db sqlRowQueryer, query string, args ...any) *sql.Row {
	return db.QueryRow(query, args...)
}

func (q QueryResult[T]) Destruct() (T, error) {
	return q.Data, q.Error
}
