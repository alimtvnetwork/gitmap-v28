package cmdrun

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const sqlCreateRunErrorsTable = `CREATE TABLE IF NOT EXISTS run_errors (
    ErrorId       TEXT PRIMARY KEY,
    FilePath      TEXT NOT NULL,
    FileExtension TEXT NOT NULL,
    Interpreter   TEXT NOT NULL,
    ExitCode      INTEGER NOT NULL,
    DurationMs    INTEGER NOT NULL,
    ErrorMessage  TEXT NOT NULL,
    StdoutSnippet TEXT NOT NULL DEFAULT '',
    StderrSnippet TEXT NOT NULL DEFAULT '',
    ExecutedArgs  TEXT NOT NULL DEFAULT '',
    CreatedAt     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_run_errors_created_at ON run_errors(CreatedAt DESC);`

// EnsureRunErrorsTable creates the run_errors table and indexes if they do not exist.
func EnsureRunErrorsTable(conn *sql.DB) error {
	res := store.ExecWrapper(conn, sqlCreateRunErrorsTable)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdrun.ensure_run_errors_table")
	}

	return nil
}

// RecordRunError persists an execution failure into the run_errors table and internal error log.
func RecordRunError(rec RunErrorRecord) error {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "cmdrun.record_run_error: open tasks db")
	}
	defer tasksDB.Close()

	if tableErr := EnsureRunErrorsTable(tasksDB.Conn()); tableErr != nil {
		return tableErr
	}

	insertErr := insertRunErrorRow(tasksDB.Conn(), rec)
	if insertErr != nil {
		return insertErr
	}

	logInternalRunError(rec)

	return nil
}

func insertRunErrorRow(conn *sql.DB, rec RunErrorRecord) error {
	if rec.ErrorID == "" {
		rec.ErrorID = fmt.Sprintf("runerr-%d", time.Now().UnixNano())
	}

	query := `INSERT INTO run_errors (ErrorId, FilePath, FileExtension, Interpreter, ExitCode, DurationMs, ErrorMessage, StdoutSnippet, StderrSnippet, ExecutedArgs, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	res := store.ExecWrapper(conn, query, rec.ErrorID, rec.FilePath, rec.FileExtension, rec.Interpreter, rec.ExitCode, rec.DurationMs, rec.ErrorMessage, rec.StdoutSnippet, rec.StderrSnippet, rec.ExecutedArgs)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdrun.insert_run_error")
	}

	return nil
}

func logInternalRunError(rec RunErrorRecord) {
	internalRec := store.InternalErrorRecord{
		ErrorCode:  fmt.Sprintf("E_RUN_%d", rec.ExitCode),
		ErrorType:  "run-error",
		Command:    "gitmap run",
		Message:    rec.ErrorMessage,
		Details:    fmt.Sprintf("file=%s interpreter=%s args=%s", rec.FilePath, rec.Interpreter, rec.ExecutedArgs),
		SourceFile: rec.FilePath,
	}

	store.LogInternalErrorRecord(internalRec)
}

// QueryRecentRunErrors returns the most recent execution failure records.
func QueryRecentRunErrors(limit int) ([]RunErrorRecord, error) {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdrun.query_run_errors: open tasks db")
	}
	defer tasksDB.Close()

	if tableErr := EnsureRunErrorsTable(tasksDB.Conn()); tableErr != nil {
		return nil, tableErr
	}

	boundedLimit := sanitizeLimit(limit)
	query := `SELECT ErrorId, FilePath, FileExtension, Interpreter, ExitCode, DurationMs, ErrorMessage, StdoutSnippet, StderrSnippet, ExecutedArgs, CreatedAt
FROM run_errors ORDER BY CreatedAt DESC LIMIT ?`

	rows, queryErr := tasksDB.Conn().Query(query, boundedLimit)
	if queryErr != nil {
		return nil, apperror.WrapSimple(queryErr, "cmdrun.query_run_errors: select")
	}
	defer rows.Close()

	return scanRunErrorRows(rows)
}

func sanitizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}

	return limit
}

func scanRunErrorRows(rows *sql.Rows) ([]RunErrorRecord, error) {
	var records []RunErrorRecord
	for rows.Next() {
		var r RunErrorRecord
		var createdAtStr string
		scanErr := rows.Scan(
			&r.ErrorID, &r.FilePath, &r.FileExtension, &r.Interpreter,
			&r.ExitCode, &r.DurationMs, &r.ErrorMessage,
			&r.StdoutSnippet, &r.StderrSnippet, &r.ExecutedArgs, &createdAtStr,
		)
		if scanErr != nil {
			continue
		}
		t, _ := time.Parse(time.RFC3339, createdAtStr)
		r.CreatedAt = t
		records = append(records, r)
	}

	return records, nil
}

// ClearRunErrors purges all records from the run_errors table.
func ClearRunErrors() error {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "cmdrun.clear_run_errors: open tasks db")
	}
	defer tasksDB.Close()

	if tableErr := EnsureRunErrorsTable(tasksDB.Conn()); tableErr != nil {
		return tableErr
	}

	res := store.ExecWrapper(tasksDB.Conn(), `DELETE FROM run_errors`)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdrun.clear_run_errors: delete")
	}

	return nil
}
