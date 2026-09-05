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

// JoinBuilder provides scoped methods for configuring a table join before returning to QueryBuilder.
type JoinBuilder[T any, F ~string] struct {
	parent          *QueryBuilder[T, F]
	joinType        string
	targetTable     string
	projectedFields []string
	extraConditions []string
	extraArgs       []any
}

func newJoinBuilder[T any, F ~string](parent *QueryBuilder[T, F], joinType, table string) *JoinBuilder[T, F] {
	return &JoinBuilder[T, F]{
		parent:      parent,
		joinType:    joinType,
		targetTable: table,
	}
}

// Select specifies projected fields from the joined table.
func (j *JoinBuilder[T, F]) Select(fields ...string) *JoinBuilder[T, F] {
	j.projectedFields = append(j.projectedFields, fields...)
	return j
}

// And adds an extra filter condition to the ON clause.
func (j *JoinBuilder[T, F]) And(column string, op SqlOperator, val any) *JoinBuilder[T, F] {
	compiler := j.parent.repo.db.Compiler()
	quotedCol := j.parent.qualifyColumn(compiler, j.targetTable, column)
	j.extraConditions = append(j.extraConditions, fmt.Sprintf("%s %s ?", quotedCol, op.String()))
	j.extraArgs = append(j.extraArgs, val)
	return j
}

// AndRaw adds a raw SQL condition to the ON clause.
func (j *JoinBuilder[T, F]) AndRaw(condition string) *JoinBuilder[T, F] {
	j.extraConditions = append(j.extraConditions, condition)
	return j
}

// On specifies the raw ON condition and transitions back to QueryBuilder.
func (j *JoinBuilder[T, F]) On(condition string) *QueryBuilder[T, F] {
	j.parent.joins = append(j.parent.joins, joinClause{
		joinType:        j.joinType,
		table:           j.targetTable,
		on:              condition,
		projectedFields: j.projectedFields,
		extraConditions: j.extraConditions,
		extraArgs:       j.extraArgs,
	})
	return j.parent
}

// OnField specifies a type-safe column-to-column condition and transitions back to QueryBuilder.
func (j *JoinBuilder[T, F]) OnField(mainCol F, op SqlOperator, joinCol string) *QueryBuilder[T, F] {
	compiler := j.parent.repo.db.Compiler()
	quotedFirst := j.parent.qualifyColumn(compiler, j.parent.repo.tableName, string(mainCol))
	quotedSecond := j.parent.qualifyColumn(compiler, j.targetTable, joinCol)
	baseOn := fmt.Sprintf("%s %s %s", quotedFirst, op.String(), quotedSecond)
	return j.On(baseOn)
}

// QueryBuilder provides a fluent, type-safe SQL query interface.
type QueryBuilder[T any, F ~string] struct {
	repo           *Repository[T, F]
	selectedFields []string
	wheres         []whereClause
	joins          []joinClause
	groupByFields  []string
	havings        []havingClause
	cteName        string
	cteSql         string
	orderByField   string
	orderDir       string
	limit          int
	offset         int
}

// NewQueryBuilder initializes a fluent QueryBuilder for a repository.
func NewQueryBuilder[T any, F ~string](repo *Repository[T, F]) *QueryBuilder[T, F] {
	return &QueryBuilder[T, F]{
		repo:   repo,
		limit:  -1,
		offset: -1,
	}
}

// Select specifies the projected fields for the root table.
func (b *QueryBuilder[T, F]) Select(fields ...F) *QueryBuilder[T, F] {
	for _, f := range fields {
		b.selectedFields = append(b.selectedFields, string(f))
	}
	return b
}

// SelectRaw specifies raw or cross-table projected fields.
func (b *QueryBuilder[T, F]) SelectRaw(fields ...string) *QueryBuilder[T, F] {
	b.selectedFields = append(b.selectedFields, fields...)
	return b
}

// Join starts an INNER JOIN clause returning a scoped JoinBuilder.
func (b *QueryBuilder[T, F]) Join(table string) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "INNER JOIN", table)
}

// InnerJoin starts an INNER JOIN clause returning a scoped JoinBuilder.
func (b *QueryBuilder[T, F]) InnerJoin(table string) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "INNER JOIN", table)
}

// LeftJoin starts a LEFT OUTER JOIN clause returning a scoped JoinBuilder.
func (b *QueryBuilder[T, F]) LeftJoin(table string) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "LEFT JOIN", table)
}

// RightJoin starts a RIGHT OUTER JOIN clause returning a scoped JoinBuilder.
func (b *QueryBuilder[T, F]) RightJoin(table string) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "RIGHT JOIN", table)
}

// FullOuterJoin starts a FULL OUTER JOIN clause returning a scoped JoinBuilder.
func (b *QueryBuilder[T, F]) FullOuterJoin(table string) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "FULL OUTER JOIN", table)
}

// OuterJoin starts a FULL OUTER JOIN clause returning a scoped JoinBuilder.
func (b *QueryBuilder[T, F]) OuterJoin(table string) *JoinBuilder[T, F] {
	return newJoinBuilder(b, "FULL OUTER JOIN", table)
}

// JoinOn adds an INNER JOIN directly with an ON condition without transitioning to JoinBuilder.
func (b *QueryBuilder[T, F]) JoinOn(table string, on string, projectedFields ...string) *QueryBuilder[T, F] {
	b.joins = append(b.joins, joinClause{
		joinType:        "INNER JOIN",
		table:           table,
		on:              on,
		projectedFields: projectedFields,
	})
	return b
}

// InnerJoinOn adds an INNER JOIN directly with an ON condition without transitioning to JoinBuilder.
func (b *QueryBuilder[T, F]) InnerJoinOn(table string, on string, projectedFields ...string) *QueryBuilder[T, F] {
	return b.JoinOn(table, on, projectedFields...)
}

// LeftJoinOn adds a LEFT JOIN directly with an ON condition without transitioning to JoinBuilder.
func (b *QueryBuilder[T, F]) LeftJoinOn(table string, on string, projectedFields ...string) *QueryBuilder[T, F] {
	b.joins = append(b.joins, joinClause{
		joinType:        "LEFT JOIN",
		table:           table,
		on:              on,
		projectedFields: projectedFields,
	})
	return b
}

// Where adds a comparison condition with an operator string.
func (b *QueryBuilder[T, F]) Where(field F, op string, val any) *QueryBuilder[T, F] {
	b.wheres = append(b.wheres, whereClause{
		clauseType: whereOp,
		field:      string(field),
		op:         op,
		val:        val,
	})
	return b
}

// WhereOp adds a comparison condition using the strongly typed SqlOperator enum.
func (b *QueryBuilder[T, F]) WhereOp(field F, op SqlOperator, val any) *QueryBuilder[T, F] {
	return b.Where(field, op.String(), val)
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

// InnerWhere adds a column-to-column condition (e.g. Table1.Field1 = Table2.Field2).
func (b *QueryBuilder[T, F]) InnerWhere(firstTableField F, op SqlOperator, secondTableField string) *QueryBuilder[T, F] {
	b.wheres = append(b.wheres, whereClause{
		clauseType: whereColumnOp,
		field:      string(firstTableField),
		op:         op.String(),
		targetCol:  secondTableField,
	})
	return b
}

// GroupBy adds fields to the GROUP BY clause.
func (b *QueryBuilder[T, F]) GroupBy(fields ...F) *QueryBuilder[T, F] {
	for _, f := range fields {
		b.groupByFields = append(b.groupByFields, string(f))
	}
	return b
}

// GroupByRaw adds raw or cross-table column expressions to the GROUP BY clause.
func (b *QueryBuilder[T, F]) GroupByRaw(fields ...string) *QueryBuilder[T, F] {
	b.groupByFields = append(b.groupByFields, fields...)
	return b
}

// Having adds a condition to the HAVING clause.
func (b *QueryBuilder[T, F]) Having(field F, op SqlOperator, val any) *QueryBuilder[T, F] {
	b.havings = append(b.havings, havingClause{
		field: string(field),
		op:    op,
		val:   val,
	})
	return b
}

// HavingRaw adds a raw SQL condition to the HAVING clause.
func (b *QueryBuilder[T, F]) HavingRaw(condition string) *QueryBuilder[T, F] {
	b.havings = append(b.havings, havingClause{
		isRaw:   true,
		rawCond: condition,
	})
	return b
}

// HavingCount adds a COUNT(*) condition to the HAVING clause.
func (b *QueryBuilder[T, F]) HavingCount(op SqlOperator, count int64) *QueryBuilder[T, F] {
	b.havings = append(b.havings, havingClause{
		isCount: true,
		op:      op,
		val:     count,
	})
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

// Signature computes a deterministic cache key for the query structure.
func (b *QueryBuilder[T, F]) Signature() string {
	var sb strings.Builder
	sb.WriteString(b.repo.tableName)
	sb.WriteString("|sel:")
	sb.WriteString(strings.Join(b.selectedFields, ","))
	sb.WriteString("|cte:")
	sb.WriteString(b.cteName)
	sb.WriteString(":")
	sb.WriteString(b.cteSql)
	for _, j := range b.joins {
		sb.WriteString("|join:")
		sb.WriteString(j.joinType)
		sb.WriteString(":")
		sb.WriteString(j.table)
		sb.WriteString(":")
		sb.WriteString(j.on)
		sb.WriteString(":")
		sb.WriteString(strings.Join(j.projectedFields, ","))
		sb.WriteString(":")
		sb.WriteString(strings.Join(j.extraConditions, ","))
	}
	for _, w := range b.wheres {
		sb.WriteString("|w:")
		sb.WriteString(fmt.Sprintf("%d:%s:%s:%s", w.clauseType, w.field, w.op, w.targetCol))
	}
	sb.WriteString("|grp:")
	sb.WriteString(strings.Join(b.groupByFields, ","))
	for _, h := range b.havings {
		sb.WriteString("|h:")
		sb.WriteString(fmt.Sprintf("%s:%s:%v", h.field, h.op.String(), h.val))
	}
	sb.WriteString("|ord:")
	sb.WriteString(b.orderByField)
	sb.WriteString(":")
	sb.WriteString(b.orderDir)
	sb.WriteString("|lim:")
	sb.WriteString(fmt.Sprintf("%d:%d", b.limit, b.offset))
	return sb.String()
}

// Compile compiles the query into SQL and extracts arguments.
// Resulting SQL is cached in GlobalQueryCache for instant reuse on subsequent executions.
func (b *QueryBuilder[T, F]) Compile() (string, []any) {
	cacheKey := b.Signature()
	if cachedSql, found := GlobalQueryCache.Get(cacheKey); found {
		_, args := b.BuildSelect()
		return cachedSql, args
	}

	sqlStr, args := b.BuildSelect()
	GlobalQueryCache.Put(cacheKey, sqlStr)
	return sqlStr, args
}

// CreateViewOrUseView checks if a view exists and contains the required columns.
// If valid, it reuses the view. If missing or schema differs, it compiles the query and creates/updates the view.
func (b *QueryBuilder[T, F]) CreateViewOrUseView(ctx context.Context, viewName string, requiredColumns ...string) BoolResult {
	viewSql := b.BuildSelectForView()
	return b.repo.db.CreateViewOrUseView(ctx, viewName, viewSql, requiredColumns...)
}

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

func (b *QueryBuilder[T, F]) buildCtePrefix(compiler DialectCompiler) string {
	if len(b.cteName) == 0 || len(b.cteSql) == 0 {
		return ""
	}
	cleanSub := strings.TrimRight(strings.TrimSpace(b.cteSql), ";")
	return fmt.Sprintf("WITH %s AS (%s) ", compiler.QuoteIdentifier(b.cteName), cleanSub)
}

func (b *QueryBuilder[T, F]) buildProjectionList(compiler DialectCompiler) string {
	var cols []string
	quotedMain := compiler.QuoteIdentifier(b.repo.tableName)

	for _, f := range b.selectedFields {
		if strings.Contains(f, ".") {
			parts := strings.SplitN(f, ".", 2)
			cols = append(cols, compiler.QuoteIdentifier(parts[0])+"."+compiler.QuoteIdentifier(parts[1]))
			continue
		}
		if len(b.joins) > 0 {
			cols = append(cols, quotedMain+"."+compiler.QuoteIdentifier(f))
			continue
		}
		cols = append(cols, compiler.QuoteIdentifier(f))
	}

	for _, j := range b.joins {
		quotedJoin := compiler.QuoteIdentifier(j.table)
		for _, pf := range j.projectedFields {
			if strings.Contains(pf, ".") {
				parts := strings.SplitN(pf, ".", 2)
				cols = append(cols, compiler.QuoteIdentifier(parts[0])+"."+compiler.QuoteIdentifier(parts[1]))
				continue
			}
			cols = append(cols, quotedJoin+"."+compiler.QuoteIdentifier(pf))
		}
	}

	if len(cols) == 0 {
		return "*"
	}
	return strings.Join(cols, ", ")
}

func (b *QueryBuilder[T, F]) buildJoins(compiler DialectCompiler, paramIdx *int) (string, []any) {
	if len(b.joins) == 0 {
		return "", nil
	}

	parts := make([]string, 0, len(b.joins))
	var args []any

	for _, j := range b.joins {
		quotedTarget := compiler.QuoteIdentifier(j.table)
		onCondition := j.on

		if len(j.extraConditions) > 0 {
			var conditions []string
			for i, cond := range j.extraConditions {
				placeholder := compiler.Placeholder(*paramIdx)
				*paramIdx++
				conditions = append(conditions, strings.Replace(cond, "?", placeholder, 1))
				args = append(args, j.extraArgs[i])
			}
			onCondition = fmt.Sprintf("%s AND %s", onCondition, strings.Join(conditions, " AND "))
		}

		parts = append(parts, fmt.Sprintf("%s %s ON %s", j.joinType, quotedTarget, onCondition))
	}

	return strings.Join(parts, " "), args
}

func (b *QueryBuilder[T, F]) buildJoinsForView(compiler DialectCompiler) string {
	if len(b.joins) == 0 {
		return ""
	}

	parts := make([]string, 0, len(b.joins))
	for _, j := range b.joins {
		quotedTarget := compiler.QuoteIdentifier(j.table)
		onCondition := j.on

		if len(j.extraConditions) > 0 {
			var conditions []string
			for i, cond := range j.extraConditions {
				lit := formatSqlLiteral(j.extraArgs[i])
				conditions = append(conditions, strings.Replace(cond, "?", lit, 1))
			}
			onCondition = fmt.Sprintf("%s AND %s", onCondition, strings.Join(conditions, " AND "))
		}

		parts = append(parts, fmt.Sprintf("%s %s ON %s", j.joinType, quotedTarget, onCondition))
	}

	return strings.Join(parts, " ")
}

func (b *QueryBuilder[T, F]) qualifyColumn(compiler DialectCompiler, defaultTable, col string) string {
	if strings.Contains(col, ".") {
		parts := strings.SplitN(col, ".", 2)
		return compiler.QuoteIdentifier(parts[0]) + "." + compiler.QuoteIdentifier(parts[1])
	}
	if len(defaultTable) > 0 {
		return compiler.QuoteIdentifier(defaultTable) + "." + compiler.QuoteIdentifier(col)
	}
	return compiler.QuoteIdentifier(col)
}

func (b *QueryBuilder[T, F]) buildWheres(compiler DialectCompiler, paramIdx *int) (string, []any) {
	if len(b.wheres) == 0 {
		return "", nil
	}

	clauses := make([]string, 0, len(b.wheres))
	args := make([]any, 0, len(b.wheres))

	for _, w := range b.wheres {
		if w.clauseType == whereColumnOp {
			quotedFirst := b.qualifyColumn(compiler, b.repo.tableName, w.field)
			quotedSecond := b.qualifyColumn(compiler, "", w.targetCol)
			clauses = append(clauses, fmt.Sprintf("%s %s %s", quotedFirst, w.op, quotedSecond))
			continue
		}

		quotedField := b.qualifyColumn(compiler, b.repo.tableName, w.field)
		placeholder := compiler.Placeholder(*paramIdx)
		*paramIdx++

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

func (b *QueryBuilder[T, F]) buildWheresForView(compiler DialectCompiler) string {
	if len(b.wheres) == 0 {
		return ""
	}

	clauses := make([]string, 0, len(b.wheres))
	for _, w := range b.wheres {
		if w.clauseType == whereColumnOp {
			quotedFirst := b.qualifyColumn(compiler, b.repo.tableName, w.field)
			quotedSecond := b.qualifyColumn(compiler, "", w.targetCol)
			clauses = append(clauses, fmt.Sprintf("%s %s %s", quotedFirst, w.op, quotedSecond))
			continue
		}

		quotedField := b.qualifyColumn(compiler, b.repo.tableName, w.field)
		if w.clauseType == whereLocate {
			literalStr := formatSqlLiteral(w.val)
			clauses = append(clauses, fmt.Sprintf("INSTR(%s, %s) > 0", quotedField, literalStr))
			continue
		}

		literalVal := formatSqlLiteral(w.val)
		clauses = append(clauses, fmt.Sprintf("%s %s %s", quotedField, w.op, literalVal))
	}

	return strings.Join(clauses, " AND ")
}

func (b *QueryBuilder[T, F]) buildGroupBy(compiler DialectCompiler) string {
	if len(b.groupByFields) == 0 {
		return ""
	}

	cols := make([]string, 0, len(b.groupByFields))
	for _, f := range b.groupByFields {
		quoted := b.qualifyColumn(compiler, b.repo.tableName, f)
		cols = append(cols, quoted)
	}

	return "GROUP BY " + strings.Join(cols, ", ")
}

func (b *QueryBuilder[T, F]) buildHavings(compiler DialectCompiler, paramIdx *int) (string, []any) {
	if len(b.havings) == 0 {
		return "", nil
	}

	clauses := make([]string, 0, len(b.havings))
	args := make([]any, 0, len(b.havings))

	for _, h := range b.havings {
		if h.isCount {
			placeholder := compiler.Placeholder(*paramIdx)
			*paramIdx++
			clauses = append(clauses, fmt.Sprintf("COUNT(*) %s %s", h.op.String(), placeholder))
			args = append(args, h.val)
			continue
		}

		if h.isRaw {
			clauses = append(clauses, h.rawCond)
			continue
		}

		quoted := b.qualifyColumn(compiler, b.repo.tableName, h.field)
		placeholder := compiler.Placeholder(*paramIdx)
		*paramIdx++
		clauses = append(clauses, fmt.Sprintf("%s %s %s", quoted, h.op.String(), placeholder))
		args = append(args, h.val)
	}

	return strings.Join(clauses, " AND "), args
}

func (b *QueryBuilder[T, F]) buildHavingsForView(compiler DialectCompiler) string {
	if len(b.havings) == 0 {
		return ""
	}

	clauses := make([]string, 0, len(b.havings))
	for _, h := range b.havings {
		if h.isCount {
			clauses = append(clauses, fmt.Sprintf("COUNT(*) %s %v", h.op.String(), h.val))
			continue
		}

		if h.isRaw {
			clauses = append(clauses, h.rawCond)
			continue
		}

		quoted := b.qualifyColumn(compiler, b.repo.tableName, h.field)
		literalVal := formatSqlLiteral(h.val)
		clauses = append(clauses, fmt.Sprintf("%s %s %s", quoted, h.op.String(), literalVal))
	}

	return strings.Join(clauses, " AND ")
}

func formatSqlLiteral(val any) string {
	if val == nil {
		return "NULL"
	}
	switch v := val.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case bool:
		if v {
			return "1"
		}
		return "0"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v)
	case fmt.Stringer:
		return "'" + strings.ReplaceAll(v.String(), "'", "''") + "'"
	default:
		return "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
	}
}

func (b *QueryBuilder[T, F]) buildOrderBy(compiler DialectCompiler) string {
	if len(b.orderByField) == 0 {
		return ""
	}
	quotedField := b.qualifyColumn(compiler, b.repo.tableName, b.orderByField)
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

// SelectTable creates a dynamic QueryBuilder starting with a table name and projected fields.
func SelectTable(db *DbWrapper, tableName string, fields ...string) *QueryBuilder[map[string]any, string] {
	repo := NewRepository[map[string]any, string](db, tableName, func(row RowScanner) (*map[string]any, error) {
		res := make(map[string]any)
		return &res, nil
	})
	qb := NewQueryBuilder(repo)
	qb.SelectRaw(fields...)
	return qb
}
