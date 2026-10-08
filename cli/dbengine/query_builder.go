package dbengine

import (
	"fmt"
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// BuildSelect compiles the current query builder state into a SQL string and arguments slice.
func (b *QueryBuilder[T, F]) BuildSelect() (string, []any) {
	compiler := b.repo.db.Compiler()
	quotedTable := compiler.QuoteIdentifier(b.repo.tableName)

	var sqlParts []string
	var args []any
	paramIdx := 1

	ctePrefix := b.buildCtePrefix(compiler)
	projectionList := b.buildProjectionList(compiler)
	sqlParts = append(sqlParts, fmt.Sprintf("%sSELECT %s FROM %s", ctePrefix, projectionList, quotedTable))

	joinSql, joinArgs := b.buildJoins(compiler, &paramIdx)
	if len(joinSql) > 0 {
		sqlParts = append(sqlParts, joinSql)
		args = append(args, joinArgs...)
	}

	whereSql, whereArgs := b.buildWheres(compiler, &paramIdx)
	if len(whereSql) > 0 {
		sqlParts = append(sqlParts, "WHERE "+whereSql)
		args = append(args, whereArgs...)
	}

	groupBySql := b.buildGroupBy(compiler)
	if len(groupBySql) > 0 {
		sqlParts = append(sqlParts, groupBySql)
	}

	havingSql, havingArgs := b.buildHavings(compiler, &paramIdx)
	if len(havingSql) > 0 {
		sqlParts = append(sqlParts, "HAVING "+havingSql)
		args = append(args, havingArgs...)
	}

	orderSql := b.buildOrderBy(compiler)
	if len(orderSql) > 0 {
		sqlParts = append(sqlParts, orderSql)
	}

	paginationSql := compiler.CompilePagination(b.limit, b.offset)
	if len(paginationSql) > 0 {
		sqlParts = append(sqlParts, paginationSql)
	}

	fullSql := strings.Join(sqlParts, " ") + ";"

	return fullSql, args
}

// BuildSelectForView compiles the query into a static SELECT SQL statement suitable for a CREATE VIEW statement.
func (b *QueryBuilder[T, F]) BuildSelectForView() string {
	compiler := b.repo.db.Compiler()
	quotedTable := compiler.QuoteIdentifier(b.repo.tableName)

	var sqlParts []string

	ctePrefix := b.buildCtePrefix(compiler)
	projectionList := b.buildProjectionList(compiler)
	sqlParts = append(sqlParts, fmt.Sprintf("%sSELECT %s FROM %s", ctePrefix, projectionList, quotedTable))

	joinSql := b.buildJoinsForView(compiler)
	if len(joinSql) > 0 {
		sqlParts = append(sqlParts, joinSql)
	}

	whereSql := b.buildWheresForView(compiler)
	if len(whereSql) > 0 {
		sqlParts = append(sqlParts, "WHERE "+whereSql)
	}

	groupBySql := b.buildGroupBy(compiler)
	if len(groupBySql) > 0 {
		sqlParts = append(sqlParts, groupBySql)
	}

	havingSql := b.buildHavingsForView(compiler)
	if len(havingSql) > 0 {
		sqlParts = append(sqlParts, "HAVING "+havingSql)
	}

	orderSql := b.buildOrderBy(compiler)
	if len(orderSql) > 0 {
		sqlParts = append(sqlParts, orderSql)
	}

	return strings.Join(sqlParts, " ")
}

// BuildCount compiles the current query builder state into a SELECT COUNT(*) statement.
func (b *QueryBuilder[T, F]) BuildCount() (string, []any) {
	compiler := b.repo.db.Compiler()
	quotedTable := compiler.QuoteIdentifier(b.repo.tableName)

	var sqlParts []string
	var args []any
	paramIdx := 1

	ctePrefix := b.buildCtePrefix(compiler)
	sqlParts = append(sqlParts, fmt.Sprintf("%sSELECT COUNT(*) FROM %s", ctePrefix, quotedTable))

	joinSql, joinArgs := b.buildJoins(compiler, &paramIdx)
	if len(joinSql) > 0 {
		sqlParts = append(sqlParts, joinSql)
		args = append(args, joinArgs...)
	}

	whereSql, whereArgs := b.buildWheres(compiler, &paramIdx)
	if len(whereSql) > 0 {
		sqlParts = append(sqlParts, "WHERE "+whereSql)
		args = append(args, whereArgs...)
	}

	groupBySql := b.buildGroupBy(compiler)
	if len(groupBySql) > 0 {
		sqlParts = append(sqlParts, groupBySql)
	}

	havingSql, havingArgs := b.buildHavings(compiler, &paramIdx)
	if len(havingSql) > 0 {
		sqlParts = append(sqlParts, "HAVING "+havingSql)
		args = append(args, havingArgs...)
	}

	fullSql := strings.Join(sqlParts, " ") + ";"

	return fullSql, args
}

// BuildDelete compiles the current query builder state into a DELETE statement.
func (b *QueryBuilder[T, F]) BuildDelete() (string, []any) {
	compiler := b.repo.db.Compiler()
	quotedTable := compiler.QuoteIdentifier(b.repo.tableName)

	var sqlParts []string
	var args []any
	paramIdx := 1

	sqlParts = append(sqlParts, fmt.Sprintf("DELETE FROM %s", quotedTable))

	whereSql, whereArgs := b.buildWheres(compiler, &paramIdx)
	if len(whereSql) > 0 {
		sqlParts = append(sqlParts, "WHERE "+whereSql)
		args = append(args, whereArgs...)
	}

	fullSql := strings.Join(sqlParts, " ") + ";"

	return fullSql, args
}

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
