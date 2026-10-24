package dbengine

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// Exact typed result envelopes wrapping canonical result types.
type (
	Uint64Result        = result.Uint64Result
	Int64Result         = result.Int64Result
	StringResult        = result.StringResult
	BoolResult          = result.BoolResult
	RowsAffectedResult  = result.Int64Result
	EntityResult[T any] = result.Wrap[*T]
	ListResult[T any]   = result.WrapSlice[T]
	CompiledQueryWrap   = result.Wrap[CompiledQuery]
	CompiledQueryResult = CompiledQueryWrap

	// DBWrapperWrap is the canonical single reusable result container for a database wrapper.
	DBWrapperWrap = result.Wrap[*DBWrapper]

	// DBRowWrap is the canonical single reusable result container for a single query row.
	DBRowWrap = result.Wrap[*sql.Row]

	// DBRowsWrap is the canonical single reusable result container for a query rows cursor.
	DBRowsWrap = result.Wrap[*sql.Rows]

	// DBResultWrap is the canonical single reusable result container for an exec result.
	DBResultWrap = result.Wrap[sql.Result]

	// DBStmtWrap is the canonical single reusable result container for a prepared statement.
	DBStmtWrap = result.Wrap[*sql.Stmt]
)

// SuccessUint64 wraps a uint64 value in a successful result envelope.
func SuccessUint64(val uint64) Uint64Result {

	return result.Ok(val)

}

// FailureUint64 wraps an AppError in a failed Uint64Result envelope.
func FailureUint64(err *apperror.AppError) Uint64Result {

	return result.Fail[uint64](err)

}

// SuccessInt64 wraps an int64 value in a successful result envelope.
func SuccessInt64(val int64) Int64Result {

	return result.Ok(val)

}

// FailureInt64 wraps an AppError in a failed Int64Result envelope.
func FailureInt64(err *apperror.AppError) Int64Result {

	return result.Fail[int64](err)

}

// SuccessString wraps a string value in a successful result envelope.
func SuccessString(val string) StringResult {

	return result.Ok(val)

}

// FailureString wraps an AppError in a failed StringResult envelope.
func FailureString(err *apperror.AppError) StringResult {

	return result.Fail[string](err)

}

// SuccessBool wraps a boolean value in a successful result envelope.
func SuccessBool(val bool) BoolResult {

	return result.Ok(val)

}

// FailureBool wraps an AppError in a failed BoolResult envelope.
func FailureBool(err *apperror.AppError) BoolResult {

	return result.Fail[bool](err)

}

// SuccessRowsAffected wraps the number of rows affected in a successful result envelope.
func SuccessRowsAffected(val int64) RowsAffectedResult {

	return result.Ok(val)

}

// FailureRowsAffected wraps an AppError in a failed RowsAffectedResult envelope.
func FailureRowsAffected(err *apperror.AppError) RowsAffectedResult {

	return result.Fail[int64](err)

}

// SuccessEntity wraps an entity pointer in a successful result envelope.
func SuccessEntity[T any](val *T) EntityResult[T] {

	return result.Ok(val)

}

// FailureEntity wraps an AppError in a failed EntityResult envelope.
func FailureEntity[T any](err *apperror.AppError) EntityResult[T] {

	return result.Fail[*T](err)

}

// SuccessList wraps an entity slice in a successful result envelope.
func SuccessList[T any](val []T) ListResult[T] {

	return result.OkSlice(val)

}

// FailureList wraps an AppError in a failed ListResult envelope.
func FailureList[T any](err *apperror.AppError) ListResult[T] {

	return result.FailSlice[T](err)

}

// CompiledQuery encapsulates compiled SQL text, bound parameters, and deterministic hash.
type CompiledQuery struct {
	SQL       string
	Args      []any
	QueryHash string
}

// SuccessCompiledQuery wraps a CompiledQuery in a successful result envelope.
func SuccessCompiledQuery(val CompiledQuery) CompiledQueryResult {

	return result.Ok(val)

}

// FailureCompiledQuery wraps an AppError in a failed CompiledQueryResult envelope.
func FailureCompiledQuery(err *apperror.AppError) CompiledQueryResult {

	return result.Fail[CompiledQuery](err)

}

// SuccessDBWrapper wraps a database wrapper in a successful result envelope.
func SuccessDBWrapper(val *DBWrapper) DBWrapperWrap {
	return result.Ok(val)
}

// FailureDBWrapper wraps an AppError in a failed DBWrapperWrap envelope.
func FailureDBWrapper(err *apperror.AppError) DBWrapperWrap {
	return result.Fail[*DBWrapper](err)
}

// SuccessDBRow wraps a single query row in a successful result envelope.
func SuccessDBRow(val *sql.Row) DBRowWrap {
	return result.Ok(val)
}

// FailureDBRow wraps an AppError in a failed DBRowWrap envelope.
func FailureDBRow(err *apperror.AppError) DBRowWrap {
	return result.Fail[*sql.Row](err)
}

// SuccessDBRows wraps a query rows cursor in a successful result envelope.
func SuccessDBRows(val *sql.Rows) DBRowsWrap {
	return result.Ok(val)
}

// FailureDBRows wraps an AppError in a failed DBRowsWrap envelope.
func FailureDBRows(err *apperror.AppError) DBRowsWrap {
	return result.Fail[*sql.Rows](err)
}

// SuccessDBResult wraps an exec result in a successful result envelope.
func SuccessDBResult(val sql.Result) DBResultWrap {
	return result.Ok(val)
}

// FailureDBResult wraps an AppError in a failed DBResultWrap envelope.
func FailureDBResult(err *apperror.AppError) DBResultWrap {
	return result.Fail[sql.Result](err)
}

// SuccessDBStmt wraps a prepared statement in a successful result envelope.
func SuccessDBStmt(val *sql.Stmt) DBStmtWrap {
	return result.Ok(val)
}

// FailureDBStmt wraps an AppError in a failed DBStmtWrap envelope.
func FailureDBStmt(err *apperror.AppError) DBStmtWrap {
	return result.Fail[*sql.Stmt](err)
}
