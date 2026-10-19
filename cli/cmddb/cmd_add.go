package cmddb

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const dbAddHelp = `gitmap db-add — Insert records into a SQLite database table.

Usage:
  gitmap db-add <alias-or-path> <table>
  gitmap db-add <alias-or-path> <table> csv <val1>, <val2>, ...
  gitmap db-add <alias-or-path> <table> col1=val1, col2=val2, ...

Arguments:
  <alias-or-path>   Target database ('master', 'pipeline', 'tasks', 'heatmap', or file path)
  <table>           Target table name
  csv ...           Comma-separated row values to insert

Examples:
  gitmap db-add tasks Subtask csv "Fix bug", "IN_PROGRESS", 1
  gitmap db-add tasks Subtask title="Fix bug", status="IN_PROGRESS"
`

// RunDbAdd executes the database record insertion command.
func RunDbAdd(args []string) *appfault.AppError {
	if len(args) == 0 || (len(args) == 1 && isDbHelpToken(args[0])) {
		fmt.Print(dbAddHelp)
		return nil
	}

	dbArg, tableArg, rawValues, isInteractive, argErr := parseAddArgs(args)
	if argErr != nil {
		return argErr
	}

	conn, dbPath, openErr := OpenResolvedDb(dbArg)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "db.add.open")
	}
	defer conn.Close()

	cols, colErr := fetchColumnInfos(conn, tableArg)
	if colErr != nil {
		return apperror.WrapSimple(colErr, "db.add.columns")
	}

	if len(cols) == 0 {
		return appfault.NewNotFoundError(fmt.Sprintf("table '%s' not found or has no columns in %s", tableArg, dbPath))
	}

	colNames, valArgs, prepErr := prepareInsertPayload(cols, tableArg, rawValues, isInteractive)
	if prepErr != nil {
		return prepErr
	}

	lastID, execErr := executeInsert(conn, tableArg, colNames, valArgs)
	if execErr != nil {
		return apperror.WrapSimple(execErr, "db.add.insert")
	}

	fmt.Printf("Inserted 1 row into '%s' (rowid: %d).\n", tableArg, lastID)
	return nil
}

func isDbHelpToken(token string) bool {
	return token == "--help" || token == "-h" || token == "help"
}

func parseAddArgs(args []string) (string, string, []string, bool, *appfault.AppError) {
	if len(args) < 2 {
		return "", "", nil, false, appfault.NewValidationError("usage: gitmap db-add <alias-or-path> <table> [csv <col1>, <col2>, ...]")
	}

	dbArg := args[0]
	tableArg := args[1]
	if len(args) == 2 {
		return dbArg, tableArg, nil, true, nil
	}

	tokens := extractCsvTokens(args[2:])
	return dbArg, tableArg, tokens, false, nil
}

func extractCsvTokens(remaining []string) []string {
	words := remaining
	if len(remaining) > 0 && strings.EqualFold(remaining[0], "csv") {
		words = remaining[1:]
	}
	return parseTokensFromWords(words)
}

func prepareInsertPayload(cols []ColumnInfo, table string, rawVals []string, isInteractive bool) ([]string, []any, *appfault.AppError) {
	if isInteractive {
		return promptInteractivePayload(cols, table)
	}
	return determineInsertCols(cols, table, rawVals)
}

func promptInteractivePayload(cols []ColumnInfo, table string) ([]string, []any, *appfault.AppError) {
	if !IsInteractiveStdin() {
		return nil, nil, appfault.NewValidationError(fmt.Sprintf("interactive input requires a terminal. Usage: gitmap db-add <alias-or-path> %s csv <col1>, <col2>, ...", table))
	}

	fmt.Printf("Entering new row for table '%s':\n", table)
	scanner := bufio.NewScanner(os.Stdin)
	var colNames []string
	var valArgs []any

	for _, col := range cols {
		val, shouldInclude := promptSingleColumn(scanner, col)
		if shouldInclude {
			colNames = append(colNames, col.Name)
			valArgs = append(valArgs, val)
		}
	}

	return colNames, valArgs, nil
}

func promptSingleColumn(scanner *bufio.Scanner, col ColumnInfo) (any, bool) {
	fmt.Printf("  %s (%s)%s: ", col.Name, col.Type, formatColPromptHint(col))
	if !scanner.Scan() {
		return nil, false
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return nil, !col.IsPrimaryKey && col.DefaultValue == nil && !col.NotNull
	}
	return input, true
}

func formatColPromptHint(col ColumnInfo) string {
	if col.IsPrimaryKey {
		return " [PK - leave blank for auto]"
	}
	if col.DefaultValue != nil {
		return fmt.Sprintf(" [default: %s]", *col.DefaultValue)
	}
	return ""
}

func determineInsertCols(cols []ColumnInfo, table string, tokens []string) ([]string, []any, *appfault.AppError) {
	if checkHasKeyVal(tokens) {
		return parseKeyValInsert(cols, table, tokens)
	}
	return parsePositionalInsert(cols, table, tokens)
}

func parseKeyValInsert(cols []ColumnInfo, table string, tokens []string) ([]string, []any, *appfault.AppError) {
	colMap := indexColumns(cols)
	var colNames []string
	var valArgs []any

	for _, t := range tokens {
		name, val, ok := splitKeyValueToken(t)
		if !ok || !colMap[strings.ToLower(name)] {
			return nil, nil, appfault.NewValidationError(fmt.Sprintf("unknown column or invalid pair '%s' for table '%s'", t, table))
		}
		colNames = append(colNames, name)
		valArgs = append(valArgs, val)
	}

	return colNames, valArgs, nil
}

func parsePositionalInsert(cols []ColumnInfo, table string, tokens []string) ([]string, []any, *appfault.AppError) {
	pkCol := findPrimaryKeyCol(cols)
	targetCols := cols
	if pkCol != nil && len(tokens) == len(cols)-1 {
		targetCols = filterNonPkCols(cols, pkCol.Name)
	}

	if len(tokens) != len(targetCols) {
		return nil, nil, appfault.NewValidationError(fmt.Sprintf("expected %d values for table '%s', got %d", len(targetCols), table, len(tokens)))
	}

	var colNames []string
	var valArgs []any
	for i, c := range targetCols {
		colNames = append(colNames, c.Name)
		valArgs = append(valArgs, tokens[i])
	}
	return colNames, valArgs, nil
}

func executeInsert(conn *sql.DB, table string, colNames []string, valArgs []any) (int64, error) {
	if len(colNames) == 0 {
		return 0, fmt.Errorf("no columns specified for insert")
	}

	quotedCols := make([]string, len(colNames))
	placeholders := make([]string, len(colNames))
	for i, c := range colNames {
		quotedCols[i] = fmt.Sprintf(`"%s"`, strings.ReplaceAll(c, `"`, `""`))
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`,
		strings.ReplaceAll(table, `"`, `""`),
		strings.Join(quotedCols, ", "),
		strings.Join(placeholders, ", "),
	)

	res, err := conn.Exec(query, valArgs...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
