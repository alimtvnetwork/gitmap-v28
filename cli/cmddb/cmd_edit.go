package cmddb

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const dbEditHelp = `gitmap db-edit — Update a record in a SQLite database table by ID.

Usage:
  gitmap db-edit <alias-or-path> -table <table> -id <id> csv <col1>, <col2>, ...
  gitmap db-edit <alias-or-path> -table <table> -id <id> "col1=val1, col2=val2, ..."

Flags:
  -table, --table, -t <table>   Target table name
  -id, --id <id>               Record identifier to update
  csv ...                      New values to apply (positional or key=value)

Examples:
  gitmap db-edit tasks -table Subtask -id 42 csv "Updated title", "DONE"
  gitmap db-edit tasks -table Subtask -id 42 "Status=DONE, Priority=HIGH"
`

type editParams struct {
	dbArg    string
	table    string
	id       string
	csvWords []string
}

// RunDbEdit executes the database record update command.
func RunDbEdit(args []string) *appfault.AppError {
	if isDbHelpRequested(args) {
		fmt.Print(dbEditHelp)

		return nil
	}

	params, argErr := parseEditArgs(args)
	if argErr != nil {
		return argErr
	}

	return processDbEdit(params)
}

func processDbEdit(params editParams) *appfault.AppError {
	conn, dbPath, openErr := OpenResolvedDb(params.dbArg)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "db.edit.open")
	}
	defer conn.Close()

	cols, colErr := fetchColumnInfos(conn, params.table)
	if colErr != nil {
		return apperror.WrapSimple(colErr, "db.edit.columns")
	}

	return executeValidatedUpdate(conn, dbPath, params, cols)
}

func executeValidatedUpdate(conn *sql.DB, dbPath string, params editParams, cols []ColumnInfo) *appfault.AppError {
	if len(cols) == 0 {
		return appfault.NewNotFoundError(fmt.Sprintf("table '%s' not found or has no columns in %s", params.table, dbPath))
	}

	colNames, valArgs, prepErr := prepareEditPayload(cols, params.table, params.csvWords)
	if prepErr != nil {
		return prepErr
	}

	pkCol := resolveTablePkName(cols)
	rowsAffected, execErr := executeUpdate(conn, params.table, pkCol, params.id, colNames, valArgs)
	if execErr != nil {
		return apperror.WrapSimple(execErr, "db.edit.update")
	}

	return finalizeEditResult(dbPath, params.table, pkCol, params.id, rowsAffected)
}

func finalizeEditResult(dbPath, table, pkCol, id string, rowsAffected int64) *appfault.AppError {
	if rowsAffected == 0 {
		return appfault.NewNotFoundError(fmt.Sprintf("no row found in '%s' where %s = %s", table, pkCol, id))
	}

	InvalidateTableMetadata(dbPath)
	fmt.Printf("Updated %d row(s) in '%s' (%s = %s).\n", rowsAffected, table, pkCol, id)

	return nil
}

func parseEditArgs(args []string) (editParams, *appfault.AppError) {
	var p editParams
	i := 0
	for i < len(args) {
		consumed, isFlag := parseEditFlagToken(args, i, &p)
		if isFlag {
			i += consumed
			continue
		}
		assignEditPositional(args[i], &p)
		i++
	}

	return validateEditParams(p)
}

func validateEditParams(p editParams) (editParams, *appfault.AppError) {
	if p.dbArg == "" || p.table == "" || p.id == "" || len(p.csvWords) == 0 {
		return p, appfault.NewValidationError("usage: gitmap db-edit <alias-or-path> -table <table> -id <id> csv <col1>, <col2>, ... (or \"col=val, ...\")")
	}

	return p, nil
}

func parseEditFlagToken(args []string, idx int, p *editParams) (int, bool) {
	if n, isTable := parseTableEditFlag(args, idx, p); isTable {
		return n, true
	}

	if n, isId := parseIdEditFlag(args, idx, p); isId {
		return n, true
	}

	if strings.EqualFold(args[idx], "csv") {
		p.csvWords = append(p.csvWords, args[idx+1:]...)

		return len(args) - idx, true
	}

	return 0, false
}

func parseTableEditFlag(args []string, idx int, p *editParams) (int, bool) {
	arg := args[idx]
	if slices.Contains([]string{"-table", "--table", "-t"}, arg) && idx+1 < len(args) {
		p.table = strings.Trim(args[idx+1], `"'`)
		return 2, true
	}

	if strings.HasPrefix(arg, "-table") && len(arg) > len("-table") {
		p.table = cleanPrefixedFlagValue(arg, "-table")
		return 1, true
	}

	if strings.HasPrefix(arg, "--table") && len(arg) > len("--table") {
		p.table = cleanPrefixedFlagValue(arg, "--table")
		return 1, true
	}

	return 0, false
}

func parseIdEditFlag(args []string, idx int, p *editParams) (int, bool) {
	arg := args[idx]
	if slices.Contains([]string{"-id", "--id"}, arg) && idx+1 < len(args) {
		p.id = strings.Trim(args[idx+1], `"'`)
		return 2, true
	}

	if strings.HasPrefix(arg, "-id") && len(arg) > len("-id") {
		p.id = cleanPrefixedFlagValue(arg, "-id")
		return 1, true
	}

	if strings.HasPrefix(arg, "--id") && len(arg) > len("--id") {
		p.id = cleanPrefixedFlagValue(arg, "--id")
		return 1, true
	}

	return 0, false
}

func cleanPrefixedFlagValue(raw, prefix string) string {
	val := strings.TrimPrefix(raw, prefix)
	val = strings.TrimPrefix(val, "=")

	return strings.Trim(val, `"'`)
}

func assignEditPositional(arg string, p *editParams) {
	if p.dbArg == "" {
		p.dbArg = arg
		return
	}

	if p.table == "" {
		p.table = arg
		return
	}

	if p.id == "" {
		p.id = arg
		return
	}

	p.csvWords = append(p.csvWords, arg)
}

func resolveTablePkName(cols []ColumnInfo) string {
	for _, c := range cols {
		if c.IsPrimaryKey {
			return c.Name
		}
	}

	for _, c := range cols {
		if strings.EqualFold(c.Name, "id") {
			return c.Name
		}
	}

	return "rowid"
}

func prepareEditPayload(cols []ColumnInfo, table string, words []string) ([]string, []any, *appfault.AppError) {
	tokens := parseTokensFromWords(words)
	if len(tokens) == 0 {
		return nil, nil, appfault.NewValidationError("no update values provided")
	}

	if checkHasKeyVal(tokens) {
		return parseKeyValInsert(cols, table, tokens)
	}

	return parseEditPositional(cols, table, tokens)
}

func parseEditPositional(cols []ColumnInfo, table string, tokens []string) ([]string, []any, *appfault.AppError) {
	targetCols, err := determineEditTargetCols(cols, table, tokens)
	if err != nil {
		return nil, nil, err
	}

	var colNames []string
	var valArgs []any
	for i, c := range targetCols {
		colNames = append(colNames, c.Name)
		valArgs = append(valArgs, tokens[i])
	}

	return colNames, valArgs, nil
}

func determineEditTargetCols(cols []ColumnInfo, table string, tokens []string) ([]ColumnInfo, *appfault.AppError) {
	pkCol := findPrimaryKeyCol(cols)
	targetCols := cols
	if pkCol != nil {
		targetCols = filterNonPkCols(cols, pkCol.Name)
	}

	if len(tokens) != len(targetCols) && len(tokens) != len(cols) {
		return nil, appfault.NewValidationError(fmt.Sprintf("expected %d values to update in '%s', got %d", len(targetCols), table, len(tokens)))
	}

	if len(tokens) == len(cols) {
		return cols, nil
	}

	return targetCols, nil
}

func executeUpdate(conn *sql.DB, table, pkCol, idVal string, colNames []string, valArgs []any) (int64, error) {
	query := buildUpdateQuery(table, pkCol, colNames)
	execArgs := make([]any, 0, len(valArgs)+1)
	execArgs = append(execArgs, valArgs...)
	execArgs = append(execArgs, idVal)

	res, err := conn.Exec(query, execArgs...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func buildUpdateQuery(table, pkCol string, colNames []string) string {
	setClauses := make([]string, len(colNames))
	for i, c := range colNames {
		setClauses[i] = fmt.Sprintf(`"%s" = ?`, strings.ReplaceAll(c, `"`, `""`))
	}

	return fmt.Sprintf(`UPDATE "%s" SET %s WHERE "%s" = ?`,
		strings.ReplaceAll(table, `"`, `""`),
		strings.Join(setClauses, ", "),
		strings.ReplaceAll(pkCol, `"`, `""`),
	)
}
