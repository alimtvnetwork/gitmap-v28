package dbengine

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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
