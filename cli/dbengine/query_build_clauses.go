package dbengine

import (
	"fmt"
	"strings"
)

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
