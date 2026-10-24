package dbengine

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type whereType int

const (
	whereOp whereType = iota
	whereLocate
	whereColumnOp
)

type whereClause struct {
	clauseType whereType
	field      string
	op         string
	val        any
	targetCol  string
}

type joinClause struct {
	joinType        string
	table           string
	on              string
	projectedFields []string
	extraConditions []string
	extraArgs       []any
}

type havingClause struct {
	field   string
	op      SqlOperator
	val     any
	isCount bool
	isRaw   bool
	rawCond string
}

// First executes the query and returns the first matching record in an EntityResult envelope.
func (b *QueryBuilder[T, F]) First(ctx context.Context) EntityResult[T] {
	if b.err != nil {
		return FailureEntity[T](b.err)
	}

	b.limit = 1
	sqlStr, args := b.BuildSelect()

	rowRes := b.repo.db.QueryRow(ctx, sqlStr, args...)
	if rowRes.IsFailure() {
		return FailureEntity[T](rowRes.AppError())
	}

	item, scanErr := b.repo.scanner(rowRes.Value)
	if scanErr != nil {
		return FailureEntity[T](apperror.WrapSimple(scanErr, "scan first "+b.repo.tableName))
	}

	return SuccessEntity(item)
}

// FindAll executes the query and returns all matching records in a ListResult envelope.
func (b *QueryBuilder[T, F]) FindAll(ctx context.Context) ListResult[T] {
	if b.err != nil {
		return FailureList[T](b.err)
	}

	sqlStr, args := b.BuildSelect()

	rowsRes := b.repo.db.Query(ctx, sqlStr, args...)
	if rowsRes.IsFailure() {
		return FailureList[T](rowsRes.AppError())
	}

	rows := rowsRes.Value
	defer rows.Close()

	var items []T
	for rows.Next() {
		item, scanErr := b.repo.scanner(rows)
		if scanErr != nil {
			return FailureList[T](apperror.WrapSimple(scanErr, "scan row "+b.repo.tableName))
		}

		items = append(items, *item)
	}

	return SuccessList(items)
}

// Count executes the query as a count aggregation and returns an Int64Result envelope.
func (b *QueryBuilder[T, F]) Count(ctx context.Context) Int64Result {
	if b.err != nil {
		return FailureInt64(b.err)
	}

	sqlStr, args := b.BuildCount()

	rowRes := b.repo.db.QueryRow(ctx, sqlStr, args...)
	if rowRes.IsFailure() {
		return FailureInt64(rowRes.AppError())
	}

	var count int64
	scanErr := rowRes.Value.Scan(&count)
	if scanErr != nil {
		return FailureInt64(apperror.WrapSimple(scanErr, "count "+b.repo.tableName))
	}

	return SuccessInt64(count)
}

// Delete executes a DELETE query matching the builder conditions and returns a RowsAffectedResult.
func (b *QueryBuilder[T, F]) Delete(ctx context.Context) RowsAffectedResult {
	if b.err != nil {
		return FailureRowsAffected(b.err)
	}

	sqlStr, args := b.BuildDelete()
	execRes := b.repo.db.Exec(ctx, sqlStr, args...)
	if execRes.IsFailure() {
		return FailureRowsAffected(execRes.AppError())
	}

	affected, err := execRes.Value.RowsAffected()
	if err != nil {
		return FailureRowsAffected(apperror.WrapSimple(err, "get rows affected for delete "+b.repo.tableName))
	}

	return SuccessRowsAffected(affected)
}

// SelectTable creates a dynamic QueryBuilder starting with a table name and projected fields.
func SelectTable(db SqlExecutor, tableName string, fields ...string) *QueryBuilder[map[string]any, string] {
	repo := NewRepository[map[string]any, string](db, tableName, func(row RowScanner) (*map[string]any, error) {
		res := make(map[string]any)

		return &res, nil
	})
	qb := NewQueryBuilder(repo)
	qb.SelectRaw(fields...)

	return qb
}
