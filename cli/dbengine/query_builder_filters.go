package dbengine

import (
	"strings"
)

// Where adds a comparison condition with an operator string.
func (b *QueryBuilder[T, F]) Where(field F, op string, val any) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

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
// Accepts field enums or strings for secondTableField.
func (b *QueryBuilder[T, F]) InnerWhere(firstTableField F, op SqlOperator, secondTableField any) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.wheres = append(b.wheres, whereClause{
		clauseType: whereColumnOp,
		field:      string(firstTableField),
		op:         op.String(),
		targetCol:  toColumnName(secondTableField),
	})

	return b
}

// GroupBy adds fields to the GROUP BY clause.
func (b *QueryBuilder[T, F]) GroupBy(fields ...F) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	for _, f := range fields {
		b.groupByFields = append(b.groupByFields, string(f))
	}

	return b
}

// GroupByRaw adds raw or cross-table column expressions to the GROUP BY clause.
func (b *QueryBuilder[T, F]) GroupByRaw(fields ...string) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.groupByFields = append(b.groupByFields, fields...)

	return b
}

// Having adds a condition to the HAVING clause.
func (b *QueryBuilder[T, F]) Having(field F, op SqlOperator, val any) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.havings = append(b.havings, havingClause{
		field: string(field),
		op:    op,
		val:   val,
	})

	return b
}

// HavingRaw adds a raw SQL condition to the HAVING clause.
func (b *QueryBuilder[T, F]) HavingRaw(condition string) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.havings = append(b.havings, havingClause{
		isRaw:   true,
		rawCond: condition,
	})

	return b
}

// HavingCount adds a COUNT(*) condition to the HAVING clause.
func (b *QueryBuilder[T, F]) HavingCount(op SqlOperator, count int64) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.havings = append(b.havings, havingClause{
		isCount: true,
		op:      op,
		val:     count,
	})

	return b
}

// WithView defines an ad-hoc Common Table Expression (CTE) view: WITH viewName AS (subQuery).
func (b *QueryBuilder[T, F]) WithView(viewName string, subQuery string) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.cteName = viewName
	b.cteSql = subQuery

	return b
}

// OrderBy sets ascending or descending order for a field.
func (b *QueryBuilder[T, F]) OrderBy(field F, dir string) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

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
	if b.err != nil {
		return b
	}

	b.limit = limit

	return b
}

// Offset sets the record offset for pagination.
func (b *QueryBuilder[T, F]) Offset(offset int) *QueryBuilder[T, F] {
	if b.err != nil {
		return b
	}

	b.offset = offset

	return b
}
