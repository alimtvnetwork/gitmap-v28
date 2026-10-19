package cmddb

import (
	"database/sql"
	"fmt"
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
	if len(args) == 0 || (len(args) == 1 && isDbHelpToken(args[0])) {
		fmt.Print(dbEditHelp)
		return nil
	}

	params, argErr := parseEditArgs(args)
	if argErr != nil {
		return argErr
	}

	conn, dbPath, openErr := OpenResolvedDb(params.dbArg)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "db.edit.open")
	}
	defer conn.Close()

	cols, colErr := fetchColumnInfos(conn, params.table)
	if colErr != nil {
		return apperror.WrapSimple(colErr, "db.edit.columns")
	}

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

	if rowsAffected == 0 {
		return appfault.NewNotFoundError(fmt.Sprintf("no row found in '%s' where %s = %s", params.table, pkCol, params.id))
	}

	fmt.Printf("Updated %d row(s) in '%s' (%s = %s).\n", rowsAffected, params.table, pkCol, params.id)
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

	if p.dbArg == "" || p.table == "" || p.id == "" || len(p.csvWords) == 0 {
		return p, appfault.NewValidationError("usage: gitmap db-edit <alias-or-path> -table <table> -id <id> csv <col1>, <col2>, ... (or \"col=val, ...\")")
	}

	return p, nil
}

func parseEditFlagToken(args []string, idx int, p *editParams) (int, bool) {
	arg := args[idx]
	if (arg == "-table" || arg == "--table" || arg == "-t") && idx+1 < len(args) {
		p.table = args[idx+1]
		return 2, true
	}
	if strings.HasPrefix(arg, "-table=") || strings.HasPrefix(arg, "--table=") {
		p.table = strings.TrimPrefix(strings.TrimPrefix(arg, "--table="), "-table=")
		return 1, true
	}
	if (arg == "-id" || arg == "--id") && idx+1 < len(args) {
		p.id = args[idx+1]
		return 2, true
	}
	if strings.HasPrefix(arg, "-id=") || strings.HasPrefix(arg, "--id=") {
		p.id = strings.TrimPrefix(strings.TrimPrefix(arg, "--id="), "-id=")
		return 1, true
	}
	if strings.EqualFold(arg, "csv") {
		p.csvWords = append(p.csvWords, args[idx+1:]...)
		return len(args) - idx, true
	}
	return 0, false
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
	raw := strings.Join(words, " ")
	tokens := parseCsvTokens(raw)
	if len(tokens) == 0 {
		return nil, nil, appfault.NewValidationError("no update values provided")
	}

	if checkHasKeyVal(tokens) {
		return parseKeyValInsert(cols, table, tokens)
	}

	return parseEditPositional(cols, table, tokens)
}

func parseEditPositional(cols []ColumnInfo, table string, tokens []string) ([]string, []any, *appfault.AppError) {
	pkCol := findPrimaryKeyCol(cols)
	targetCols := cols
	if pkCol != nil && len(tokens) == len(cols)-1 {
		targetCols = filterNonPkCols(cols, pkCol.Name)
	}

	if len(tokens) != len(targetCols) {
		return nil, nil, appfault.NewValidationError(fmt.Sprintf("expected %d values to update in '%s', got %d", len(targetCols), table, len(tokens)))
	}

	var colNames []string
	var valArgs []any
	for i, c := range targetCols {
		colNames = append(colNames, c.Name)
		valArgs = append(valArgs, tokens[i])
	}

	return colNames, valArgs, nil
}

func executeUpdate(conn *sql.DB, table, pkCol, idVal string, colNames []string, valArgs []any) (int64, error) {
	setClauses := make([]string, len(colNames))
	for i, c := range colNames {
		setClauses[i] = fmt.Sprintf(`"%s" = ?`, strings.ReplaceAll(c, `"`, `""`))
	}

	query := fmt.Sprintf(`UPDATE "%s" SET %s WHERE "%s" = ?`,
		strings.ReplaceAll(table, `"`, `""`),
		strings.Join(setClauses, ", "),
		strings.ReplaceAll(pkCol, `"`, `""`),
	)

	execArgs := make([]any, 0, len(valArgs)+1)
	execArgs = append(execArgs, valArgs...)
	execArgs = append(execArgs, idVal)

	res, err := conn.Exec(query, execArgs...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
