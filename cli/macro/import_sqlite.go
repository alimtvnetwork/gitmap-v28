package macro

import (
	"database/sql"
	"os"
	"time"

	_ "modernc.org/sqlite"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ParseImportSQLite reads macros and their steps from a SQLite database file.
func ParseImportSQLite(dbPath string) ([]Macro, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, apperror.WrapSimple(err, "sqlite file not found")
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open sqlite import db")
	}

	defer db.Close()
	db.SetMaxOpenConns(1)

	return queryAllMacrosWithSteps(db)
}

func queryAllMacrosWithSteps(db *sql.DB) ([]Macro, error) {
	stepMap, err := queryAllMacroSteps(db)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, COALESCE(description, ''), created_at, updated_at, total_steps, COALESCE(tags, '') FROM macros ORDER BY id ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query macros table")
	}

	defer rows.Close()

	return scanMacroRows(rows, stepMap)
}

func scanMacroRows(rows *sql.Rows, stepMap map[string][]MacroStep) ([]Macro, error) {
	var list []Macro
	for rows.Next() {
		m, err := scanSingleMacroRow(rows, stepMap)
		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	return list, nil
}

func scanSingleMacroRow(rows *sql.Rows, stepMap map[string][]MacroStep) (Macro, error) {
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

func queryAllMacroSteps(db *sql.DB) (map[string][]MacroStep, error) {
	query := `SELECT id, macro_name, step_num, command_line, COALESCE(working_dir, ''), continue_on_error, timeout_seconds FROM macro_steps ORDER BY step_num ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query macro_steps table")
	}

	defer rows.Close()

	return scanMacroStepsMap(rows)
}

func scanMacroStepsMap(rows *sql.Rows) (map[string][]MacroStep, error) {
	stepMap := make(map[string][]MacroStep)
	for rows.Next() {
		var step MacroStep
		var macroName string
		var continueVal int
		if err := rows.Scan(&step.ID, &macroName, &step.StepNum, &step.CommandLine, &step.WorkingDir, &continueVal, &step.TimeoutSeconds); err != nil {
			return nil, apperror.WrapSimple(err, "scan macro step row")
		}

		step.ContinueOnError = continueVal == 1
		stepMap[macroName] = append(stepMap[macroName], step)
	}

	return stepMap, nil
}
