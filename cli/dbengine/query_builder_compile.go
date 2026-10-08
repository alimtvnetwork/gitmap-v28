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
