// Package store — ai_instruction_db.go: isolated SQLite database for AI command execution history.
package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	sqlCreateAiCommandCategory = `CREATE TABLE IF NOT EXISTS AiCommandCategory (
    CategoryId   INTEGER PRIMARY KEY AUTOINCREMENT,
    CategoryCode TEXT NOT NULL UNIQUE,
    Name         TEXT NOT NULL,
    Description  TEXT NOT NULL DEFAULT '',
    IsActive     INTEGER NOT NULL DEFAULT 1,
    CreatedAt    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxAiCommandCategory_Code ON AiCommandCategory(CategoryCode);`

	sqlSeedAiCommandCategory = `INSERT OR IGNORE INTO AiCommandCategory (CategoryCode, Name, Description) VALUES
('general', 'General AI', 'General AI execution'),
('code_gen', 'Code Generation', 'Automated code generation'),
('audit', 'AI Audit', 'AI-assisted code and repo audit'),
('search', 'AI Search', 'AI-assisted search and discovery'),
('instruction', 'Instruction Execution', 'Execution of structured AI instructions');`

	sqlCreateAiExecutionHistory = `CREATE TABLE IF NOT EXISTS AiExecutionHistory (
    ExecutionId   INTEGER PRIMARY KEY AUTOINCREMENT,
    CategoryCode  TEXT NOT NULL DEFAULT 'general',
    CommandText   TEXT NOT NULL,
    ArgsJson      TEXT NOT NULL DEFAULT '',
    WorkDir       TEXT NOT NULL DEFAULT '',
    CallerIp      TEXT NOT NULL DEFAULT '',
    DurationMs    INTEGER NOT NULL DEFAULT 0,
    ExitCode      INTEGER NOT NULL DEFAULT 0,
    OutputSummary TEXT NOT NULL DEFAULT '',
    ErrMsg        TEXT NOT NULL DEFAULT '',
    IsSuccess     INTEGER NOT NULL DEFAULT 1,
    ExecutedAt    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CreatedAt     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxAiExecutionHistory_Category ON AiExecutionHistory(CategoryCode);
CREATE INDEX IF NOT EXISTS IdxAiExecutionHistory_IsSuccess ON AiExecutionHistory(IsSuccess);
CREATE INDEX IF NOT EXISTS IdxAiExecutionHistory_ExecutedAt ON AiExecutionHistory(ExecutedAt);`

	sqlCreateAiExecutionView = `CREATE VIEW IF NOT EXISTS AiExecutionView AS
SELECT
    h.ExecutionId,
    h.CategoryCode,
    COALESCE(c.Name, h.CategoryCode) AS CategoryName,
    h.CommandText,
    h.ArgsJson,
    h.WorkDir,
    h.CallerIp,
    h.DurationMs,
    h.ExitCode,
    h.OutputSummary,
    h.ErrMsg,
    h.IsSuccess,
    h.ExecutedAt,
    h.CreatedAt
FROM AiExecutionHistory h
LEFT JOIN AiCommandCategory c ON h.CategoryCode = c.CategoryCode;`
)

// AiInstructionSplitDB wraps an isolated SQLite database connection for AI instruction history.
type AiInstructionSplitDB struct {
	conn *sql.DB
	Path string
}

// AiInstructionDbPath returns the full path to the AI instruction SQLite DB.
func AiInstructionDbPath() string {
	return ResolveAiInstructionDbPath("")
}

// OpenAiInstructionSplitDB opens or initializes the AI instruction split database.
func OpenAiInstructionSplitDB() (*AiInstructionSplitDB, error) {
	return OpenAiInstructionSplitDBAt(AiInstructionDbPath())
}

// OpenAiInstructionSplitDBAt opens or creates an AI instruction split database at a specific path.
func OpenAiInstructionSplitDBAt(dbPath string) (*AiInstructionSplitDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "ai_instruction.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_instruction.open")
	}

	return initAiInstructionConn(conn, dbPath)
}

func initAiInstructionConn(conn *sql.DB, dbPath string) (*AiInstructionSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "ai_instruction.config")
	}

	if err := initAiInstructionSchema(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	db := &AiInstructionSplitDB{conn: conn, Path: dbPath}
	db.registerWithRegistry()
	return db, nil
}

func initAiInstructionSchema(conn *sql.DB) error {
	statements := []string{
		sqlCreateAiCommandCategory,
		sqlSeedAiCommandCategory,
		sqlCreateAiExecutionHistory,
		sqlCreateAiExecutionView,
	}
	for _, stmt := range statements {
		res := ExecWrapper(conn, stmt)
		if res.IsFailure {
			return apperror.WrapSimple(res.Error, "ai_instruction.initSchema")
		}
	}
	return nil
}

func (db *AiInstructionSplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}
	defer master.Close()
	_ = master.RegisterSplitDB(db.buildRegistryEntry())
}

func (db *AiInstructionSplitDB) buildRegistryEntry() SplitDatabaseEntry {
	return SplitDatabaseEntry{
		DatabaseType:  "ai-instruction",
		DatabaseKey:   "ai_instruction_master",
		DatabasePath:  db.Path,
		Status:        "active",
		IsActive:      true,
		Description:   "AI instruction execution history and category split database",
		SchemaVersion: 1,
	}
}

// Close terminates the database connection.
func (db *AiInstructionSplitDB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}
	return db.conn.Close()
}

// Conn returns the raw database connection.
func (db *AiInstructionSplitDB) Conn() *sql.DB {
	return db.conn
}

// RecordAiExecution logs an AI command execution into the database instance.
func (db *AiInstructionSplitDB) RecordAiExecution(
	category, commandText, argsJson, workDir, callerIp string,
	durationMs, exitCode int,
	outputSummary, errMsg string,
	isSuccess bool,
) error {
	cat := normalizeAiCategory(category)
	successVal := boolToInt(isSuccess)
	q := `INSERT INTO AiExecutionHistory (
		CategoryCode, CommandText, ArgsJson, WorkDir, CallerIp,
		DurationMs, ExitCode, OutputSummary, ErrMsg, IsSuccess
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res := ExecWrapper(db.conn, q, cat, commandText, argsJson, workDir, callerIp, durationMs, exitCode, outputSummary, errMsg, successVal)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "ai_instruction.record_execution")
	}
	return nil
}

// RecordAiExecution opens the default AI instruction database and logs an execution.
func RecordAiExecution(
	category, commandText, argsJson, workDir, callerIp string,
	durationMs, exitCode int,
	outputSummary, errMsg string,
	isSuccess bool,
) error {
	db, err := OpenAiInstructionSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.RecordAiExecution(category, commandText, argsJson, workDir, callerIp, durationMs, exitCode, outputSummary, errMsg, isSuccess)
}

func normalizeAiCategory(category string) string {
	if category != "" {
		return category
	}
	return "general"
}
