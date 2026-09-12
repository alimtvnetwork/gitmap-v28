package macro

import (
	"database/sql"
	"os"
	"time"

	_ "modernc.org/sqlite"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ParseImportSQLite reads macros and their steps from a SQLite database file.
func ParseImportSQLite(dbPath string) result.ResultSlice[Macro] {
	if _, err := os.Stat(dbPath); err != nil {
		return result.FailSlice[Macro](apperror.WrapSimple(err, "sqlite file not found"))
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return result.FailSlice[Macro](apperror.WrapSimple(err, "open sqlite import db"))
	}

	defer db.Close()
	db.SetMaxOpenConns(1)

	return queryAllMacrosWithSteps(db)
}

func queryAllMacrosWithSteps(db *sql.DB) result.ResultSlice[Macro] {
	stepRes := queryAllMacroSteps(db)
	if stepRes.IsFailure() {
		return result.FailSlice[Macro](stepRes.AppError())
	}

	query := `SELECT id, name, COALESCE(description, ''), created_at, updated_at, total_steps, COALESCE(tags, '') FROM macros ORDER BY id ASC`
	rows, err := db.Query(query)
	if err != nil {
		return result.FailSlice[Macro](apperror.WrapSimple(err, "query macros table"))
	}

	defer rows.Close()

	return scanMacroRows(rows, stepRes.Data)
}

func scanMacroRows(rows *sql.Rows, stepMap map[string][]MacroStep) result.ResultSlice[Macro] {
	var list []Macro
	for rows.Next() {
		m, err := scanSingleMacroRow(rows, stepMap)
		if err != nil {
			return result.FailSlice[Macro](err)
		}

		list = append(list, m)
	}

	return result.OkSlice(list)
}

func scanSingleMacroRow(rows *sql.Rows, stepMap map[string][]MacroStep) (Macro, *apperror.AppError) {
	var m Macro
	var createdStr, updatedStr string
	if err := rows.Scan(&m.ID, &m.Name, &m.Description, &createdStr, &updatedStr, &m.TotalSteps, &m.Tags); err != nil {
		return m, apperror.WrapSimple(err, "scan macro row")
	}

	m.CreatedAt = parseSQLiteTime(createdStr)
	m.UpdatedAt = parseSQLiteTime(updatedStr)
	m.Steps = stepMap[m.Name]
	m.TotalSteps = len(m.Steps)

	return m, nil
}

func parseSQLiteTime(raw string) time.Time {
	t, err := time.Parse(time.RFC3339, raw)
	if err == nil {
		return t
	}

	return time.Now()
}

func queryAllMacroSteps(db *sql.DB) result.ResultMap[string, []MacroStep] {
	query := `SELECT id, macro_name, step_num, command_line, COALESCE(working_dir, ''), continue_on_error, timeout_seconds FROM macro_steps ORDER BY step_num ASC`
	rows, err := db.Query(query)
	if err != nil {
		appErr := apperror.WrapSimple(err, "query macro_steps table")

		return result.FailMap[string, []MacroStep](appErr)
	}

	defer rows.Close()

	return scanMacroStepsMap(rows)
}

func scanMacroStepsMap(rows *sql.Rows) result.ResultMap[string, []MacroStep] {
	stepMap := make(map[string][]MacroStep)
	for rows.Next() {
		var s MacroStep
		var name string
		var continueVal int
		if err := rows.Scan(&s.ID, &name, &s.StepNum, &s.CommandLine, &s.WorkingDir, &continueVal, &s.TimeoutSeconds); err != nil {
			appErr := apperror.WrapSimple(err, "scan macro step row")

			return result.FailMap[string, []MacroStep](appErr)
		}

		s.ContinueOnError = continueVal == 1
		stepMap[name] = append(stepMap[name], s)
	}

	return result.OkMap(stepMap)
}
