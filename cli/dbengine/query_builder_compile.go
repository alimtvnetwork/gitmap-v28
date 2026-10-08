package dbengine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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

// QueryHash computes a deterministic SHA-256 hex string for the query structure.
func (b *QueryBuilder[T, F]) QueryHash() string {
	viewSql := b.BuildSelectForView()

	return ComputeSqlHash(viewSql)
}

// Compile compiles the query into SQL and arguments, wrapped in a CompiledQueryResult.
// Resulting SQL is cached in GlobalQueryCache for instant reuse on subsequent executions.
func (b *QueryBuilder[T, F]) Compile() CompiledQueryResult {
	if b.err != nil {
		return FailureCompiledQuery(b.err)
	}

	cacheKey := b.Signature()
	if cachedSql, found := GlobalQueryCache.Get(cacheKey); found {
		_, args := b.BuildSelect()
		hash := ComputeSqlHash(cachedSql)

		return SuccessCompiledQuery(CompiledQuery{
			SQL:       cachedSql,
			Args:      args,
			QueryHash: hash,
		})
	}

	sqlStr, args := b.BuildSelect()
	GlobalQueryCache.Put(cacheKey, sqlStr)
	hash := ComputeSqlHash(sqlStr)

	return SuccessCompiledQuery(CompiledQuery{
		SQL:       sqlStr,
		Args:      args,
		QueryHash: hash,
	})
}

// CompileRaw returns raw SQL string and arguments without result wrapper.
func (b *QueryBuilder[T, F]) CompileRaw() (string, []any) {
	res := b.Compile()
	if res.IsFailed() {
		return "", nil
	}

	cq := res.Value

	return cq.SQL, cq.Args
}

// ViewCreator defines methods for creating or reusing database views.
type ViewCreator interface {
	CreateViewOrUseViewWithHash(ctx context.Context, name string, selectSql string, queryHash string) BoolResult
	CreateViewOrUseView(ctx context.Context, name string, selectSql string, requiredColumns ...string) BoolResult
}

// CreateViewOrUseView checks if a view exists and contains the required columns (or matches query hash if omitted).
// If valid, it reuses the view. If missing or schema differs, it validates SQL, drops old view, creates updated view, and saves metadata.
func (b *QueryBuilder[T, F]) CreateViewOrUseView(ctx context.Context, viewName string, requiredColumns ...string) BoolResult {
	if b.err != nil {
		return FailureBool(b.err)
	}

	vc, isViewCreator := b.repo.db.(ViewCreator)
	if !isViewCreator {
		return FailureBool(apperror.WrapSimple(errors.New("executor does not support view creation"), "create view"))
	}

	return b.dispatchCreateView(ctx, vc, viewName, requiredColumns)
}

func (b *QueryBuilder[T, F]) dispatchCreateView(
	ctx context.Context,
	vc ViewCreator,
	viewName string,
	requiredColumns []string,
) BoolResult {
	viewSql := b.BuildSelectForView()
	if len(requiredColumns) == 0 {
		return vc.CreateViewOrUseViewWithHash(ctx, viewName, viewSql, b.QueryHash())
	}

	return vc.CreateViewOrUseView(ctx, viewName, viewSql, requiredColumns...)
}

func toColumnName(col any) string {
	if s, ok := col.(string); ok {
		return s
	}

	if str, ok := col.(fmt.Stringer); ok {
		return str.String()
	}

	return fmt.Sprintf("%v", col)
}

// JoinBuilder provides scoped methods for configuring a table join before returning to QueryBuilder.
type JoinBuilder[T any, F ~string] struct {
	parent          *QueryBuilder[T, F]
	joinType        string
	targetTable     string
	projectedFields []string
	extraConditions []string
	extraArgs       []any
	err             *apperror.AppError
}

func newJoinBuilder[T any, F ~string](parent *QueryBuilder[T, F], joinType, table string) *JoinBuilder[T, F] {
	jb := &JoinBuilder[T, F]{
		parent:      parent,
		joinType:    joinType,
		targetTable: table,
	}

	if parent.err != nil {
		jb.err = parent.err
	}

	return jb
}

// Select specifies projected fields from the joined table. Accepts strings or typed field enums.
func (j *JoinBuilder[T, F]) Select(fields ...any) *JoinBuilder[T, F] {
	if j.err != nil {
		return j
	}

	for _, f := range fields {
		j.projectedFields = append(j.projectedFields, toColumnName(f))
	}

	return j
}

// And adds an extra filter condition to the ON clause. Accepts strings or typed field enums.
func (j *JoinBuilder[T, F]) And(column any, op SqlOperator, val any) *JoinBuilder[T, F] {
	if j.err != nil {
		return j
	}

	colName := toColumnName(column)
	compiler := j.parent.repo.db.Compiler()
	quotedCol := j.parent.qualifyColumn(compiler, j.targetTable, colName)
	j.extraConditions = append(j.extraConditions, fmt.Sprintf("%s %s ?", quotedCol, op.String()))
	j.extraArgs = append(j.extraArgs, val)

	return j
}

// AndRaw adds a raw SQL condition to the ON clause.
func (j *JoinBuilder[T, F]) AndRaw(condition string) *JoinBuilder[T, F] {
	if j.err != nil {
		return j
	}

	j.extraConditions = append(j.extraConditions, condition)

	return j
}

// On specifies the raw ON condition and transitions back to QueryBuilder.
func (j *JoinBuilder[T, F]) On(condition string) *QueryBuilder[T, F] {
	if j.err != nil {
		j.parent.err = j.err

		return j.parent
	}

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
// Accepts field enums or strings for both root table column and joined table column.
func (j *JoinBuilder[T, F]) OnField(mainCol any, op SqlOperator, joinCol any) *QueryBuilder[T, F] {
	if j.err != nil {
		j.parent.err = j.err

		return j.parent
	}

	firstCol := toColumnName(mainCol)
	secondCol := toColumnName(joinCol)
	compiler := j.parent.repo.db.Compiler()
	quotedFirst := j.parent.qualifyColumn(compiler, j.parent.repo.tableName, firstCol)
	quotedSecond := j.parent.qualifyColumn(compiler, j.targetTable, secondCol)
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
	err            *apperror.AppError
}
