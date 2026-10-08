package dbengine

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// NewQueryBuilder initializes a fluent QueryBuilder for a repository.
func NewQueryBuilder[T any, F ~string](repo *Repository[T, F]) *QueryBuilder[T, F] {
	return &QueryBuilder[T, F]{
		repo:   repo,
		limit:  -1,
		offset: -1,
	}
}

// SetError records an error on the builder, causing terminal methods to return this error.
func (b *QueryBuilder[T, F]) SetError(err *apperror.AppError) *QueryBuilder[T, F] {
	if b.err == nil {
		b.err = err
	}

	return b
}

// Err returns any error accumulated on the builder.
func (b *QueryBuilder[T, F]) Err() *apperror.AppError {
	return b.err
}

// Select specifies the projected fields for the root table.
func (b *QueryBuilder[T, F]) Select(fields ...F) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	for _, f := range fields {
		b.selectedFields = append(b.selectedFields, string(f))
	}

	return b
}

// SelectRaw specifies raw or cross-table projected fields.
func (b *QueryBuilder[T, F]) SelectRaw(fields ...string) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.selectedFields = append(b.selectedFields, fields...)

	return b
}

// Join starts an INNER JOIN clause returning a scoped JoinBuilder. Accepts strings or table constants.
func (b *QueryBuilder[T, F]) Join(table any) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "INNER JOIN", toColumnName(table))
}

// InnerJoin starts an INNER JOIN clause returning a scoped JoinBuilder. Accepts strings or table constants.
func (b *QueryBuilder[T, F]) InnerJoin(table any) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "INNER JOIN", toColumnName(table))
}

// LeftJoin starts a LEFT OUTER JOIN clause returning a scoped JoinBuilder. Accepts strings or table constants.
func (b *QueryBuilder[T, F]) LeftJoin(table any) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "LEFT JOIN", toColumnName(table))
}

// RightJoin starts a RIGHT OUTER JOIN clause returning a scoped JoinBuilder. Accepts strings or table constants.
func (b *QueryBuilder[T, F]) RightJoin(table any) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "RIGHT JOIN", toColumnName(table))
}

// FullOuterJoin starts a FULL OUTER JOIN clause returning a scoped JoinBuilder. Accepts strings or table constants.
func (b *QueryBuilder[T, F]) FullOuterJoin(table any) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "FULL OUTER JOIN", toColumnName(table))
}

// OuterJoin starts a FULL OUTER JOIN clause returning a scoped JoinBuilder. Accepts strings or table constants.
func (b *QueryBuilder[T, F]) OuterJoin(table any) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "FULL OUTER JOIN", toColumnName(table))
}

// JoinOn adds an INNER JOIN directly with an ON condition without transitioning to JoinBuilder.
func (b *QueryBuilder[T, F]) JoinOn(table any, on string, projectedFields ...any) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	proj := make([]string, 0, len(projectedFields))
	for _, pf := range projectedFields {
		proj = append(proj, toColumnName(pf))
	}

	b.joins = append(b.joins, joinClause{
		joinType:        "INNER JOIN",
		table:           toColumnName(table),
		on:              on,
		projectedFields: proj,
	})

	return b
}

// InnerJoinOn adds an INNER JOIN directly with an ON condition without transitioning to JoinBuilder.
func (b *QueryBuilder[T, F]) InnerJoinOn(table any, on string, projectedFields ...any) *QueryBuilder[T, F] {
	return b.JoinOn(table, on, projectedFields...)
}

// LeftJoinOn adds a LEFT JOIN directly with an ON condition without transitioning to JoinBuilder.
func (b *QueryBuilder[T, F]) LeftJoinOn(table any, on string, projectedFields ...any) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	proj := make([]string, 0, len(projectedFields))
	for _, pf := range projectedFields {
		proj = append(proj, toColumnName(pf))
	}

	b.joins = append(b.joins, joinClause{
		joinType:        "LEFT JOIN",
		table:           toColumnName(table),
		on:              on,
		projectedFields: proj,
	})

	return b
}
