package dbengine

import (
	"fmt"
	"strings"
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
