package macro

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const macroTableDDL = `
CREATE TABLE IF NOT EXISTS macros (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    total_steps INTEGER DEFAULT 0,
    tags TEXT DEFAULT ''
);
CREATE TABLE IF NOT EXISTS macro_steps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    macro_name TEXT NOT NULL,
    step_num INTEGER NOT NULL,
    command_line TEXT NOT NULL,
    working_dir TEXT DEFAULT '',
    continue_on_error INTEGER DEFAULT 0,
    timeout_seconds INTEGER DEFAULT 0,
    FOREIGN KEY (macro_name) REFERENCES macros(name) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_macro_steps_name ON macro_steps(macro_name, step_num);
`

// ExportMacrosToSQLite exports a slice of macros into a standalone SQLite database.
func ExportMacrosToSQLite(macros []Macro, dbPath string) error {
	prepareSQLiteExportPath(dbPath)
	db, err := openExportSQLiteDB(dbPath)
	if err != nil {
		return err
	}

	defer db.Close()
	if err := initMacroSQLiteSchema(db); err != nil {
		return err
	}

	return writeMacrosTransaction(db, macros)
}

func prepareSQLiteExportPath(dbPath string) {
	_ = os.Remove(dbPath)
	_ = os.Remove(dbPath + "-wal")
	_ = os.Remove(dbPath + "-shm")
	if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
		_ = os.MkdirAll(dir, constants.DirPermission)
	}
}

func openExportSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open sqlite export db")
	}

	db.SetMaxOpenConns(1)
	if err := applyMacroSQLitePragmas(db); err != nil {
		_ = db.Close()

		return nil, err
	}

	return db, nil
}

func applyMacroSQLitePragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return apperror.WrapSimple(err, "apply sqlite pragma")
		}
	}

	return nil
}

func initMacroSQLiteSchema(db *sql.DB) error {
	if _, err := db.Exec(macroTableDDL); err != nil {
		return apperror.WrapSimple(err, "initialize macro sqlite schema")
	}

	return nil
}

func writeMacrosTransaction(db *sql.DB, macros []Macro) error {
	tx, err := db.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "begin export transaction")
	}

	for _, m := range macros {
		if err := insertMacroRow(tx, m); err != nil {
			_ = tx.Rollback()

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return apperror.WrapSimple(err, "commit export transaction")
	}

	return nil
}

func insertMacroRow(tx *sql.Tx, m Macro) error {
	query := `INSERT INTO macros (name, description, created_at, updated_at, total_steps, tags) VALUES (?, ?, ?, ?, ?, ?)`
	createdStr := m.CreatedAt.Format(time.RFC3339)
	updatedStr := m.UpdatedAt.Format(time.RFC3339)
	if _, err := tx.Exec(query, m.Name, m.Description, createdStr, updatedStr, len(m.Steps), m.Tags); err != nil {
		return apperror.WrapSimple(err, "insert macro record")
	}

	return insertMacroStepRows(tx, m.Name, m.Steps)
}

func insertMacroStepRows(tx *sql.Tx, macroName string, steps []MacroStep) error {
	query := `INSERT INTO macro_steps (macro_name, step_num, command_line, working_dir, continue_on_error, timeout_seconds) VALUES (?, ?, ?, ?, ?, ?)`
	for i, step := range steps {
		stepNum := step.StepNum
		if stepNum <= 0 {
			stepNum = i + 1
		}

		isContinue := 0
		if step.ContinueOnError {
			isContinue = 1
		}

		if _, err := tx.Exec(query, macroName, stepNum, step.CommandLine, step.WorkingDir, isContinue, step.TimeoutSeconds); err != nil {
			return apperror.WrapSimple(err, "insert macro step record")
		}
	}

	return nil
}
