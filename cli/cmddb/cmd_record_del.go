package cmddb

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const dbRecordDelHelp = `gitmap db-record-del (dbrd) — Delete a record from a SQLite database table by primary key.

Usage:
  gitmap db-record-del <db> -table <table> -id <id> [-y] [--json]
  gitmap db-record-del <db> <table> <id> [-y] [--json]

Flags:
  -t, --table, -table <table>   Target table name
  --id, -id <id>                Primary key value of record to delete
  -y, --yes, --force            Bypass confirmation prompt
  -j, --json                    Output results in JSON format
  -h, --help                    Show help text

Examples:
  gitmap db-record-del tasks -table Subtask -id 42 -y
  gitmap dbrd master Repositories 101 -y
`

type recordDelParams struct {
	dbArg       string
	table       string
	id          string
	isConfirmed bool
	isJson      bool
}

// RunDbRecordDel executes primary key record deletion from CLI.
func RunDbRecordDel(args []string) *appfault.AppError {
	if isDbHelpRequested(args) {
		fmt.Print(dbRecordDelHelp)

		return nil
	}

	params, err := parseRecordDelArgs(args)
	if err != nil {
		return err
	}

	return openAndExecuteRecordDel(params)
}

func openAndExecuteRecordDel(params recordDelParams) *appfault.AppError {
	conn, dbPath, openErr := OpenResolvedDb(params.dbArg)
	if openErr != nil {
		return appfault.NewNotFoundError(fmt.Sprintf("open database '%s': %v", params.dbArg, openErr))
	}
	defer conn.Close()

	return executeRecordDeletion(conn, dbPath, params)
}

func parseRecordDelArgs(args []string) (recordDelParams, *appfault.AppError) {
	var p recordDelParams
	i := 0
	for i < len(args) {
		consumed, isFlag := parseRecordDelFlag(args, i, &p)
		if isFlag {
			i += consumed
			continue
		}
		assignRecordDelPositional(args[i], &p)
		i++
	}

	return validateRecordDelParams(p)
}

func validateRecordDelParams(p recordDelParams) (recordDelParams, *appfault.AppError) {
	if p.dbArg == "" || p.table == "" || p.id == "" {
		return p, appfault.NewValidationError("usage: gitmap db-record-del <db> -table <table> -id <id> [-y]")
	}

	return p, nil
}

func parseRecordDelFlag(args []string, idx int, p *recordDelParams) (int, bool) {
	arg := args[idx]
	if isConfirmFlag(arg) {
		p.isConfirmed = true
		return 1, true
	}

	if slices.Contains([]string{"-j", "--json"}, arg) {
		p.isJson = true
		return 1, true
	}

	return parseTargetRecordDelFlag(args, idx, p)
}

func isConfirmFlag(arg string) bool {
	return slices.Contains([]string{"-y", "--yes", "--force"}, arg)
}

func parseTargetRecordDelFlag(args []string, idx int, p *recordDelParams) (int, bool) {
	arg := args[idx]
	if slices.Contains([]string{"-table", "--table", "-t"}, arg) && idx+1 < len(args) {
		p.table = args[idx+1]
		return 2, true
	}

	if slices.Contains([]string{"-id", "--id"}, arg) && idx+1 < len(args) {
		p.id = args[idx+1]
		return 2, true
	}

	return 0, false
}

func assignRecordDelPositional(arg string, p *recordDelParams) {
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
	}
}

func executeRecordDeletion(conn *sql.DB, dbPath string, p recordDelParams) *appfault.AppError {
	cols, err := fetchColumnInfos(conn, p.table)
	if err != nil || len(cols) == 0 {
		return appfault.NewNotFoundError(fmt.Sprintf("table '%s' not found in database '%s'", p.table, dbPath))
	}

	pkCol := resolveTablePkName(cols)
	if !p.isConfirmed && !askDeletionConfirmation(p.table, pkCol, p.id, p.dbArg) {
		fmt.Println("[ABORTED] Record deletion canceled.")

		return nil
	}

	return finishRecordDeletion(conn, dbPath, p, pkCol)
}

func finishRecordDeletion(conn *sql.DB, dbPath string, p recordDelParams, pkCol string) *appfault.AppError {
	rowsAffected, txErr := runDeleteTx(conn, p.table, pkCol, p.id)
	if txErr != nil {
		return appfault.NewValidationError(fmt.Sprintf("delete record failed: %v", txErr))
	}

	if rowsAffected == 0 {
		return appfault.NewNotFoundError(fmt.Sprintf("no record found in '%s' where %s = %s", p.table, pkCol, p.id))
	}

	InvalidateTableMetadata(dbPath)
	outputRecordDelResult(p, pkCol, rowsAffected, dbPath)

	return nil
}

func askDeletionConfirmation(table, pkCol, id, dbArg string) bool {
	fmt.Printf("[CONFIRM] Delete record (%s=%s) from '%s' in '%s'? [y/N]: ", pkCol, id, table, dbArg)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}

	trimmed := strings.ToLower(strings.TrimSpace(response))

	return trimmed == "y" || trimmed == "yes"
}

func runDeleteTx(conn *sql.DB, table, pkCol, id string) (int64, error) {
	tx, err := conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, pkCol)
	res, execErr := tx.Exec(query, id)
	if execErr != nil {
		return 0, execErr
	}

	return commitDeleteTx(tx, res)
}

func commitDeleteTx(tx *sql.Tx, res sql.Result) (int64, error) {
	rowsAffected, affErr := res.RowsAffected()
	if affErr != nil {
		return 0, affErr
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func outputRecordDelResult(p recordDelParams, pkCol string, rowsAffected int64, dbPath string) {
	if p.isJson {
		outputRecordDelJson(p, pkCol, rowsAffected, dbPath)

		return
	}

	fmt.Printf("[SUCCESS] Deleted %d record (%s=%s) from table '%s' in database '%s'.\n", rowsAffected, pkCol, p.id, p.table, p.dbArg)
}

func outputRecordDelJson(p recordDelParams, pkCol string, rowsAffected int64, dbPath string) {
	data := map[string]any{
		"Status": map[string]any{
			"IsSuccess": true,
			"Code":      200,
			"Message":   "Record deleted successfully",
		},
		"Attributes": map[string]any{
			"Database":     dbPath,
			"Table":        p.table,
			"PrimaryKey":   pkCol,
			"Id":           p.id,
			"RowsAffected": rowsAffected,
		},
		"Results": map[string]any{
			"DeletedId": p.id,
		},
	}
	_ = outputJson(os.Stdout, data)
}
