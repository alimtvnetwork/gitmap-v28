package dbengine

import (
	"context"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type whereType int

const (
	whereOp whereType = iota
	whereLocate
)

type whereClause struct {
	clauseType whereType
	field      string
	op         string
	val        any
}

type joinClause struct {
	joinType string
	table    string
	on       string
}

// QueryBuilder provides a fluent, type-safe SQL query interface.
type QueryBuilder[T any, F ~string] struct {
	repo         *Repository[T, F]
	wheres       []whereClause
	joins        []joinClause
	cteName      string
	cteSql       string
	orderByField string
	orderDir     string
	limit        int
	offset       int
}

// NewQueryBuilder initializes a fluent QueryBuilder for a repository.
func NewQueryBuilder[T any, F ~string](repo *Repository[T, F]) *QueryBuilder[T, F] {
	return &QueryBuilder[T, F]{
		repo:   repo,
		limit:  -1,
		offset: -1,
	}
}

// Where adds a comparison condition (e.g. op "=", ">", "<", "LIKE").
func (b *QueryBuilder[T, F]) Where(field F, op string, val any) *QueryBuilder[T, F] {
	b.wheres = append(b.wheres, whereClause{
		clauseType: whereOp,
		field:      string(field),
		op:         op,
		val:        val,
	})
	return b
}

// WhereEq adds an equality condition (field = val).
func (b *QueryBuilder[T, F]) WhereEq(field F, val any) *QueryBuilder[T, F] {
	return b.Where(field, "=", val)
}

// Locate adds a substring containment filter using the database locate function (e.g. INSTR).
func (b *QueryBuilder[T, F]) Locate(field F, substring string) *QueryBuilder[T, F] {
	b.wheres = append(b.wheres, whereClause{
		clauseType: whereLocate,
		field:      string(field),
		val:        substring,
	})
	return b
}

// Join adds an INNER JOIN clause.
func (b *QueryBuilder[T, F]) Join(table string, on string) *QueryBuilder[T, F] {
	b.joins = append(b.joins, joinClause{joinType: "INNER JOIN", table: table, on: on})
	return b
}

// LeftJoin adds a LEFT JOIN clause.
func (b *QueryBuilder[T, F]) LeftJoin(table string, on string) *QueryBuilder[T, F] {
	b.joins = append(b.joins, joinClause{joinType: "LEFT JOIN", table: table, on: on})
	return b
}

// WithView defines an ad-hoc Common Table Expression (CTE) view: WITH viewName AS (subQuery).
func (b *QueryBuilder[T, F]) WithView(viewName string, subQuery string) *QueryBuilder[T, F] {
	b.cteName = viewName
	b.cteSql = subQuery
	return b
}

// OrderBy sets ascending or descending order for a field.
func (b *QueryBuilder[T, F]) OrderBy(field F, dir string) *QueryBuilder[T, F] {
	b.orderByField = string(field)
	b.orderDir = strings.ToUpper(strings.TrimSpace(dir))
	return b
}

// OrderByDesc sets descending order for a field.
func (b *QueryBuilder[T, F]) OrderByDesc(field F) *QueryBuilder[T, F] {
	return b.OrderBy(field, "DESC")
}

// Limit sets the maximum number of records to return.
func (b *QueryBuilder[T, F]) Limit(limit int) *QueryBuilder[T, F] {
	b.limit = limit
	return b
}

// Offset sets the record offset for pagination.
func (b *QueryBuilder[T, F]) Offset(offset int) *QueryBuilder[T, F] {
	b.offset = offset
	return b
}

// BuildSelect compiles the current query builder state into a SQL string and arguments slice.
func (b *QueryBuilder[T, F]) BuildSelect() (string, []any) {
	compiler := b.repo.db.Compiler()
	quotedTable := compiler.QuoteIdentifier(b.repo.tableName)

	var sqlParts []string
	var args []any

	ctePrefix := b.buildCtePrefix(compiler)
	sqlParts = append(sqlParts, fmt.Sprintf("%sSELECT * FROM %s", ctePrefix, quotedTable))

	joinSql := b.buildJoins(compiler)
	if len(joinSql) > 0 {
		sqlParts = append(sqlParts, joinSql)
	}

	whereSql, whereArgs := b.buildWheres(compiler)
	if len(whereSql) > 0 {
		sqlParts = append(sqlParts, "WHERE "+whereSql)
		args = append(args, whereArgs...)
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

// BuildCount compiles the current query builder state into a SELECT COUNT(*) statement.
func (b *QueryBuilder[T, F]) BuildCount() (string, []any) {
	compiler := b.repo.db.Compiler()
	quotedTable := compiler.QuoteIdentifier(b.repo.tableName)

	var sqlParts []string
	var args []any

	ctePrefix := b.buildCtePrefix(compiler)
	sqlParts = append(sqlParts, fmt.Sprintf("%sSELECT COUNT(*) FROM %s", ctePrefix, quotedTable))

	joinSql := b.buildJoins(compiler)
	if len(joinSql) > 0 {
		sqlParts = append(sqlParts, joinSql)
	}

	whereSql, whereArgs := b.buildWheres(compiler)
	if len(whereSql) > 0 {
		sqlParts = append(sqlParts, "WHERE "+whereSql)
		args = append(args, whereArgs...)
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

	sqlParts = append(sqlParts, fmt.Sprintf("DELETE FROM %s", quotedTable))

	whereSql, whereArgs := b.buildWheres(compiler)
	if len(whereSql) > 0 {
		sqlParts = append(sqlParts, "WHERE "+whereSql)
		args = append(args, whereArgs...)
	}

	fullSql := strings.Join(sqlParts, " ") + ";"
	return fullSql, args
}

func (b *QueryBuilder[T, F]) buildCtePrefix(compiler DialectCompiler) string {
	if len(b.cteName) == 0 || len(b.cteSql) == 0 {
		return ""
	}
	cleanSub := strings.TrimRight(strings.TrimSpace(b.cteSql), ";")
	return fmt.Sprintf("WITH %s AS (%s) ", compiler.QuoteIdentifier(b.cteName), cleanSub)
}

func (b *QueryBuilder[T, F]) buildJoins(compiler DialectCompiler) string {
	if len(b.joins) == 0 {
		return ""
	}
	parts := make([]string, 0, len(b.joins))
	for _, j := range b.joins {
		quotedTarget := compiler.QuoteIdentifier(j.table)
		parts = append(parts, fmt.Sprintf("%s %s ON %s", j.joinType, quotedTarget, j.on))
	}
	return strings.Join(parts, " ")
}

func (b *QueryBuilder[T, F]) buildWheres(compiler DialectCompiler) (string, []any) {
	if len(b.wheres) == 0 {
		return "", nil
	}

	clauses := make([]string, 0, len(b.wheres))
	args := make([]any, 0, len(b.wheres))

	for i, w := range b.wheres {
		quotedField := compiler.QuoteIdentifier(w.field)
		placeholder := compiler.Placeholder(i + 1)

		if w.clauseType == whereLocate {
			clauses = append(clauses, fmt.Sprintf("INSTR(%s, %s) > 0", quotedField, placeholder))
			args = append(args, w.val)
			continue
		}

		clauses = append(clauses, fmt.Sprintf("%s %s %s", quotedField, w.op, placeholder))
		args = append(args, w.val)
	}

	return strings.Join(clauses, " AND "), args
}

func (b *QueryBuilder[T, F]) buildOrderBy(compiler DialectCompiler) string {
	if len(b.orderByField) == 0 {
		return ""
	}
	quotedField := compiler.QuoteIdentifier(b.orderByField)
	dir := b.orderDir
	if dir != "DESC" {
		dir = "ASC"
	}
	return fmt.Sprintf("ORDER BY %s %s", quotedField, dir)
}

// First executes the query and returns the first matching record in an EntityResult envelope.
func (b *QueryBuilder[T, F]) First(ctx context.Context) EntityResult[T] {
	b.limit = 1
	sqlStr, args := b.BuildSelect()

	row, appErr := b.repo.db.QueryRow(ctx, sqlStr, args...)
	if appErr != nil {
		return FailureEntity[T](appErr)
	}

	item, scanErr := b.repo.scanner(row)
	if scanErr != nil {
		return FailureEntity[T](apperror.WrapSimple(scanErr, "scan first "+b.repo.tableName))
	}

	return SuccessEntity(item)
}

// FindAll executes the query and returns all matching records in a ListResult envelope.
func (b *QueryBuilder[T, F]) FindAll(ctx context.Context) ListResult[T] {
	sqlStr, args := b.BuildSelect()

	rows, appErr := b.repo.db.Query(ctx, sqlStr, args...)
	if appErr != nil {
		return FailureList[T](appErr)
	}
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
	sqlStr, args := b.BuildCount()

	row, appErr := b.repo.db.QueryRow(ctx, sqlStr, args...)
	if appErr != nil {
		return FailureInt64(appErr)
	}

	var count int64
	scanErr := row.Scan(&count)
	if scanErr != nil {
		return FailureInt64(apperror.WrapSimple(scanErr, "count "+b.repo.tableName))
	}

	return SuccessInt64(count)
}

// Delete executes a DELETE query matching the builder conditions and returns a RowsAffectedResult.
func (b *QueryBuilder[T, F]) Delete(ctx context.Context) RowsAffectedResult {
	sqlStr, args := b.BuildDelete()
	res, appErr := b.repo.db.Exec(ctx, sqlStr, args...)
	if appErr != nil {
		return FailureRowsAffected(appErr)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return FailureRowsAffected(apperror.WrapSimple(err, "get rows affected for delete "+b.repo.tableName))
	}

	return SuccessRowsAffected(affected)
}
