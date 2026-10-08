// Package store — ai_analysis_split_db.go: isolated Split SQLite database for AI analysis tasks, file tracking, and line diffs.
package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	sqlCreateAiTask = `CREATE TABLE IF NOT EXISTS AiTask (
    AiTaskId    INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskUuid    TEXT UNIQUE,
    Goal        TEXT NOT NULL DEFAULT '',
    Status      TEXT NOT NULL DEFAULT 'pending',
    Category    TEXT NOT NULL DEFAULT 'refactor',
    Reasoning   TEXT NOT NULL DEFAULT '',
    TotalFiles  INTEGER NOT NULL DEFAULT 0,
    TotalLines  INTEGER NOT NULL DEFAULT 0,
    IsActive    INTEGER NOT NULL DEFAULT 1,
    HasFailed   INTEGER NOT NULL DEFAULT 0,
    Description TEXT NOT NULL DEFAULT '',
    Notes       TEXT NOT NULL DEFAULT '',
    Comments    TEXT NOT NULL DEFAULT '',
    StartedAt   TEXT NOT NULL DEFAULT '',
    CompletedAt TEXT NOT NULL DEFAULT '',
    CreatedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS Idx_AiTask_TaskUuid ON AiTask(TaskUuid);
CREATE INDEX IF NOT EXISTS Idx_AiTask_Status ON AiTask(Status);
CREATE INDEX IF NOT EXISTS Idx_AiTask_IsActive ON AiTask(IsActive);`

	sqlCreateAiTaskFile = `CREATE TABLE IF NOT EXISTS AiTaskFile (
    AiTaskFileId INTEGER PRIMARY KEY AUTOINCREMENT,
    AiTaskId     INTEGER NOT NULL,
    RelPath      TEXT NOT NULL,
    AbsPath      TEXT NOT NULL DEFAULT '',
    Action       TEXT NOT NULL DEFAULT '',
    BeforeSha256 TEXT NOT NULL DEFAULT '',
    AfterSha256  TEXT NOT NULL DEFAULT '',
    BackupPath   TEXT NOT NULL DEFAULT '',
    LineCount    INTEGER NOT NULL DEFAULT 0,
    Reasoning    TEXT NOT NULL DEFAULT '',
    IsModified   INTEGER NOT NULL DEFAULT 0,
    IsRemoved    INTEGER NOT NULL DEFAULT 0,
    Notes        TEXT NOT NULL DEFAULT '',
    Comments     TEXT NOT NULL DEFAULT '',
    CreatedAt    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (AiTaskId) REFERENCES AiTask(AiTaskId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS Idx_AiTaskFile_TaskId ON AiTaskFile(AiTaskId);
CREATE INDEX IF NOT EXISTS Idx_AiTaskFile_RelPath ON AiTaskFile(RelPath);
CREATE INDEX IF NOT EXISTS Idx_AiTaskFile_IsRemoved ON AiTaskFile(IsRemoved);`

	sqlCreateAiTaskLine = `CREATE TABLE IF NOT EXISTS AiTaskLine (
    AiTaskLineId    INTEGER PRIMARY KEY AUTOINCREMENT,
    AiTaskFileId    INTEGER NOT NULL,
    LineNumber      INTEGER NOT NULL DEFAULT 0,
    OriginalContent TEXT NOT NULL DEFAULT '',
    ProposedContent TEXT NOT NULL DEFAULT '',
    DiffKind        TEXT NOT NULL DEFAULT '',
    RuleViolation   TEXT NOT NULL DEFAULT '',
    Reasoning       TEXT NOT NULL DEFAULT '',
    IsApplied       INTEGER NOT NULL DEFAULT 0,
    Notes           TEXT NOT NULL DEFAULT '',
    Comments        TEXT NOT NULL DEFAULT '',
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (AiTaskFileId) REFERENCES AiTaskFile(AiTaskFileId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS Idx_AiTaskLine_FileId ON AiTaskLine(AiTaskFileId);
CREATE INDEX IF NOT EXISTS Idx_AiTaskLine_RuleViolation ON AiTaskLine(RuleViolation);
CREATE INDEX IF NOT EXISTS Idx_AiTaskLine_IsApplied ON AiTaskLine(IsApplied);`

	sqlCreateAiAnalysisTask = `CREATE TABLE IF NOT EXISTS AiAnalysisTask (
    task_id             TEXT PRIMARY KEY,
    repo_path           TEXT NOT NULL DEFAULT '',
    task_description    TEXT NOT NULL DEFAULT '',
    model_name          TEXT NOT NULL DEFAULT '',
    status              TEXT NOT NULL DEFAULT 'in_progress',
    started_at          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at        TEXT,
    total_files_read     INTEGER NOT NULL DEFAULT 0,
    total_files_modified INTEGER NOT NULL DEFAULT 0,
    total_files_deleted  INTEGER NOT NULL DEFAULT 0,
    reasoning_summary   TEXT NOT NULL DEFAULT '',
    created_at          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxAiAnalysisTask_StatusCreated ON AiAnalysisTask(status, created_at);`

	sqlCreateAiAnalysisLine = `CREATE TABLE IF NOT EXISTS AiAnalysisLine (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id         TEXT NOT NULL,
    file_path       TEXT NOT NULL,
    operation_type  TEXT NOT NULL DEFAULT 'read',
    start_line      INTEGER NOT NULL DEFAULT 0,
    end_line        INTEGER NOT NULL DEFAULT 0,
    line_count      INTEGER NOT NULL DEFAULT 0,
    content_snippet TEXT NOT NULL DEFAULT '',
    reasoning       TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES AiAnalysisTask(task_id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS IdxAiAnalysisLine_TaskIdFilePath ON AiAnalysisLine(task_id, file_path);
CREATE INDEX IF NOT EXISTS IdxAiAnalysisLine_TaskId ON AiAnalysisLine(task_id);`
)

// AiTask models an autonomous AI execution session or refactoring operation.
type AiTask struct {
	AiTaskId    int64  `json:"aiTaskId"`
	TaskUuid    string `json:"taskUuid"`
	Goal        string `json:"goal"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	Reasoning   string `json:"reasoning"`
	TotalFiles  int    `json:"totalFiles"`
	TotalLines  int    `json:"totalLines"`
	IsActive    bool   `json:"isActive"`
	HasFailed   bool   `json:"hasFailed"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
	Comments    string `json:"comments"`
	StartedAt   string `json:"startedAt"`
	CompletedAt string `json:"completedAt"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// AiTaskFile catalogs a file investigated, modified, or safely removed.
type AiTaskFile struct {
	AiTaskFileId int64  `json:"aiTaskFileId"`
	AiTaskId     int64  `json:"aiTaskId"`
	RelPath      string `json:"relPath"`
	AbsPath      string `json:"absPath"`
	Action       string `json:"action"`
	BeforeSha256 string `json:"beforeSha256"`
	AfterSha256  string `json:"afterSha256"`
	BackupPath   string `json:"backupPath"`
	LineCount    int    `json:"lineCount"`
	Reasoning    string `json:"reasoning"`
	IsModified   bool   `json:"isModified"`
	IsRemoved    bool   `json:"isRemoved"`
	Notes        string `json:"notes"`
	Comments     string `json:"comments"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// AiTaskLine catalogs line-level mutations with attached reasoning and rule citations.
type AiTaskLine struct {
	AiTaskLineId    int64  `json:"aiTaskLineId"`
	AiTaskFileId    int64  `json:"aiTaskFileId"`
	LineNumber      int    `json:"lineNumber"`
	OriginalContent string `json:"originalContent"`
	ProposedContent string `json:"proposedContent"`
	DiffKind        string `json:"diffKind"`
	RuleViolation   string `json:"ruleViolation"`
	Reasoning       string `json:"reasoning"`
	IsApplied       bool   `json:"isApplied"`
	Notes           string `json:"notes"`
	Comments        string `json:"comments"`
	CreatedAt       string `json:"createdAt"`
}

// AiAnalysisSummaryData holds aggregated counts for the AI analysis split database.
type AiAnalysisSummaryData struct {
	TotalTasks     int `json:"totalTasks"`
	ActiveTasks    int `json:"activeTasks"`
	CompletedTasks int `json:"completedTasks"`
	FailedTasks    int `json:"failedTasks"`
	TotalFiles     int `json:"totalFiles"`
	TotalLines     int `json:"totalLines"`
}

// AiAnalysisSplitDB wraps an isolated SQLite database connection for AI analysis sessions.
type AiAnalysisSplitDB struct {
	*DB
	Path string
}

// MigrateAiAnalysisSplitDb initializes schema and indexes for AI analysis.
func MigrateAiAnalysisSplitDb(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA busy_timeout=5000;",
	}
	for _, pragma := range pragmas {
		res := ExecWrapper(db, pragma)
		if res.IsFailure {
			return apperror.WrapSimple(res.Error, "ai_analysis_split.pragma")
		}
	}

	statements := []string{
		sqlCreateAiTask,
		sqlCreateAiTaskFile,
		sqlCreateAiTaskLine,
		sqlCreateAiAnalysisTask,
		sqlCreateAiAnalysisLine,
	}
	for _, stmt := range statements {
		res := ExecWrapper(db, stmt)
		if res.IsFailure {
			return apperror.WrapSimple(res.Error, "ai_analysis_split.migrate")
		}
	}

	return nil
}

// OpenAiAnalysisSplitDB opens or initializes the split DB for AI analysis.
func OpenAiAnalysisSplitDB(repoRoots ...string) (*AiAnalysisSplitDB, error) {
	repoRoot := ""
	if len(repoRoots) == 1 {
		repoRoot = repoRoots[0]
	} else if len(repoRoots) > 1 {
		repoRoot = repoRoots[1]
	}
	dbPath := ResolveAiAnalysisDbPath(repoRoot)

	return OpenAiAnalysisSplitDBAt(dbPath)
}

// OpenAiAnalysisSplitDBAt opens or creates an AI analysis split database at a specific path.
func OpenAiAnalysisSplitDBAt(dbPath string) (*AiAnalysisSplitDB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis_split.mkdir")
	}

	innerDB, err := OpenAt(dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis_split.open")
	}

	return initAiAnalysisConn(innerDB, dbPath)
}

func initAiAnalysisConn(innerDB *DB, dbPath string) (*AiAnalysisSplitDB, error) {
	if err := MigrateAiAnalysisSplitDb(innerDB.Conn()); err != nil {
		_ = innerDB.Close()

		return nil, err
	}

	db := &AiAnalysisSplitDB{DB: innerDB, Path: dbPath}
	db.registerWithRegistry()

	return db, nil
}

func (db *AiAnalysisSplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}
	defer master.Close()

	entry := SplitDatabaseEntry{
		DatabaseType:  "ai_analysis",
		DatabaseKey:   "ai_analysis_" + filepath.Base(filepath.Dir(db.Path)),
		DatabasePath:  db.Path,
		Status:        "active",
		IsActive:      true,
		Description:   "AI analysis and granular file reasoning split database",
		SchemaVersion: 1,
	}
	_ = master.RegisterSplitDB(entry)
}

// Close terminates the database connection.
func (db *AiAnalysisSplitDB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}

	return db.DB.Close()
}

// Conn returns the raw database connection.
func (db *AiAnalysisSplitDB) Conn() *sql.DB {
	if db == nil || db.DB == nil {
		return nil
	}

	return db.DB.Conn()
}

// CreateTask inserts a new AiTask record into the database.
func (db *AiAnalysisSplitDB) CreateTask(task *AiTask) (int64, error) {
	if task == nil {
		return 0, apperror.NewValidationError("task cannot be nil")
	}
	now := time.Now().UTC().Format(time.RFC3339)
	applyTaskDefaults(task, now)

	query := `INSERT INTO AiTask (
		TaskUuid, Goal, Status, Category, Reasoning,
		TotalFiles, TotalLines, IsActive, HasFailed,
		Description, Notes, Comments, StartedAt, CompletedAt,
		CreatedAt, UpdatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res := ExecWrapper(db.Conn(), query,
		task.TaskUuid, task.Goal, task.Status, task.Category, task.Reasoning,
		task.TotalFiles, task.TotalLines, boolToInt(task.IsActive), boolToInt(task.HasFailed),
		task.Description, task.Notes, task.Comments, task.StartedAt, task.CompletedAt,
		task.CreatedAt, task.UpdatedAt,
	)
	if res.IsFailure {
		return 0, apperror.WrapSimple(res.Error, "ai_analysis.CreateTask")
	}
	lastId, _ := res.Data.LastInsertId()
	task.AiTaskId = lastId

	return lastId, nil
}

func applyTaskDefaults(task *AiTask, now string) {
	if task.Status == "" {
		task.Status = "pending"
	}
	if task.Category == "" {
		task.Category = "refactor"
	}
	if task.StartedAt == "" {
		task.StartedAt = now
	}
	if task.CreatedAt == "" {
		task.CreatedAt = now
	}
	if task.UpdatedAt == "" {
		task.UpdatedAt = now
	}
}

// RecordFile inserts an AiTaskFile record and updates task file counter.
func (db *AiAnalysisSplitDB) RecordFile(file *AiTaskFile) (int64, error) {
	if file == nil {
		return 0, apperror.NewValidationError("file cannot be nil")
	}
	file.RelPath = filepath.ToSlash(filepath.Clean(file.RelPath))
	now := time.Now().UTC().Format(time.RFC3339)
	applyFileDefaults(file, now)

	query := `INSERT INTO AiTaskFile (
		AiTaskId, RelPath, AbsPath, Action, BeforeSha256, AfterSha256,
		BackupPath, LineCount, Reasoning, IsModified, IsRemoved,
		Notes, Comments, CreatedAt, UpdatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res := ExecWrapper(db.Conn(), query,
		file.AiTaskId, file.RelPath, file.AbsPath, file.Action, file.BeforeSha256, file.AfterSha256,
		file.BackupPath, file.LineCount, file.Reasoning, boolToInt(file.IsModified), boolToInt(file.IsRemoved),
		file.Notes, file.Comments, file.CreatedAt, file.UpdatedAt,
	)
	if res.IsFailure {
		return 0, apperror.WrapSimple(res.Error, "ai_analysis.RecordFile")
	}
	lastId, _ := res.Data.LastInsertId()
	file.AiTaskFileId = lastId
	_ = db.updateTaskFileCount(file.AiTaskId)

	return lastId, nil
}

func applyFileDefaults(file *AiTaskFile, now string) {
	if file.Action == "" {
		file.Action = "modified"
	}
	if file.CreatedAt == "" {
		file.CreatedAt = now
	}
	if file.UpdatedAt == "" {
		file.UpdatedAt = now
	}
}

func (db *AiAnalysisSplitDB) updateTaskFileCount(taskId int64) error {
	query := `UPDATE AiTask SET TotalFiles = (SELECT COUNT(*) FROM AiTaskFile WHERE AiTaskId = ?) WHERE AiTaskId = ?`
	res := ExecWrapper(db.Conn(), query, taskId, taskId)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "ai_analysis.updateTaskFileCount")
	}

	return nil
}

// RecordLine inserts an AiTaskLine record and updates task line counter.
func (db *AiAnalysisSplitDB) RecordLine(line *AiTaskLine) (int64, error) {
	if line == nil {
		return 0, apperror.NewValidationError("line cannot be nil")
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if line.CreatedAt == "" {
		line.CreatedAt = now
	}

	query := `INSERT INTO AiTaskLine (
		AiTaskFileId, LineNumber, OriginalContent, ProposedContent,
		DiffKind, RuleViolation, Reasoning, IsApplied,
		Notes, Comments, CreatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res := ExecWrapper(db.Conn(), query,
		line.AiTaskFileId, line.LineNumber, line.OriginalContent, line.ProposedContent,
		line.DiffKind, line.RuleViolation, line.Reasoning, boolToInt(line.IsApplied),
		line.Notes, line.Comments, line.CreatedAt,
	)
	if res.IsFailure {
		return 0, apperror.WrapSimple(res.Error, "ai_analysis.RecordLine")
	}
	lastId, _ := res.Data.LastInsertId()
	line.AiTaskLineId = lastId
	_ = db.updateTaskLineCountByFile(line.AiTaskFileId)

	return lastId, nil
}

func (db *AiAnalysisSplitDB) updateTaskLineCountByFile(fileId int64) error {
	query := `UPDATE AiTask SET TotalLines = TotalLines + 1 
		WHERE AiTaskId = (SELECT AiTaskId FROM AiTaskFile WHERE AiTaskFileId = ?)`
	res := ExecWrapper(db.Conn(), query, fileId)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "ai_analysis.updateTaskLineCount")
	}

	return nil
}

// GetTask retrieves an AiTask by its UUID.
func (db *AiAnalysisSplitDB) GetTask(taskUuid string) (*AiTask, error) {
	query := `SELECT 
		AiTaskId, TaskUuid, Goal, Status, Category, Reasoning,
		TotalFiles, TotalLines, IsActive, HasFailed, Description,
		Notes, Comments, StartedAt, CompletedAt, CreatedAt, UpdatedAt
	FROM AiTask WHERE TaskUuid = ? LIMIT 1`

	row := db.Conn().QueryRow(query, taskUuid)

	return scanAiTaskRow(row)
}

// GetTaskById retrieves an AiTask by its integer primary key.
func (db *AiAnalysisSplitDB) GetTaskById(taskId int64) (*AiTask, error) {
	query := `SELECT 
		AiTaskId, TaskUuid, Goal, Status, Category, Reasoning,
		TotalFiles, TotalLines, IsActive, HasFailed, Description,
		Notes, Comments, StartedAt, CompletedAt, CreatedAt, UpdatedAt
	FROM AiTask WHERE AiTaskId = ? LIMIT 1`

	row := db.Conn().QueryRow(query, taskId)

	return scanAiTaskRow(row)
}

func scanAiTaskRow(row *sql.Row) (*AiTask, error) {
	var t AiTask
	var isActiveInt, hasFailedInt int

	err := row.Scan(
		&t.AiTaskId, &t.TaskUuid, &t.Goal, &t.Status, &t.Category, &t.Reasoning,
		&t.TotalFiles, &t.TotalLines, &isActiveInt, &hasFailedInt, &t.Description,
		&t.Notes, &t.Comments, &t.StartedAt, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperror.NewNotFoundError("task not found")
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis.scanTask")
	}

	t.IsActive = isActiveInt == 1
	t.HasFailed = hasFailedInt == 1

	return &t, nil
}

// ListTasks retrieves a list of tasks optionally filtered by active status.
func (db *AiAnalysisSplitDB) ListTasks(isActiveOnly bool, limit int) ([]AiTask, error) {
	if limit <= 0 {
		limit = 50
	}
	query := buildListTasksQuery(isActiveOnly)
	rows, err := db.Conn().Query(query, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis.ListTasks")
	}
	defer rows.Close()

	return scanAiTasks(rows)
}

func buildListTasksQuery(isActiveOnly bool) string {
	if isActiveOnly {
		return `SELECT 
			AiTaskId, TaskUuid, Goal, Status, Category, Reasoning,
			TotalFiles, TotalLines, IsActive, HasFailed, Description,
			Notes, Comments, StartedAt, CompletedAt, CreatedAt, UpdatedAt
		FROM AiTask WHERE IsActive = 1 ORDER BY AiTaskId DESC LIMIT ?`
	}

	return `SELECT 
		AiTaskId, TaskUuid, Goal, Status, Category, Reasoning,
		TotalFiles, TotalLines, IsActive, HasFailed, Description,
		Notes, Comments, StartedAt, CompletedAt, CreatedAt, UpdatedAt
	FROM AiTask ORDER BY AiTaskId DESC LIMIT ?`
}

func scanAiTasks(rows *sql.Rows) ([]AiTask, error) {
	var list []AiTask
	for rows.Next() {
		var t AiTask
		var isActiveInt, hasFailedInt int
		err := rows.Scan(
			&t.AiTaskId, &t.TaskUuid, &t.Goal, &t.Status, &t.Category, &t.Reasoning,
			&t.TotalFiles, &t.TotalLines, &isActiveInt, &hasFailedInt, &t.Description,
			&t.Notes, &t.Comments, &t.StartedAt, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "ai_analysis.scanAiTasks")
		}
		t.IsActive = isActiveInt == 1
		t.HasFailed = hasFailedInt == 1
		list = append(list, t)
	}

	return list, nil
}

// UpdateTaskStatus updates a task's status, failure flag, and completion timestamp.
func (db *AiAnalysisSplitDB) UpdateTaskStatus(taskUuid string, status string, hasFailed bool) error {
	now := time.Now().UTC().Format(time.RFC3339)
	isActive := status == "pending" || status == "running" || status == "in_progress"
	completedAt := ""
	if !isActive {
		completedAt = now
	}

	query := `UPDATE AiTask SET 
		Status = ?, HasFailed = ?, IsActive = ?, CompletedAt = ?, UpdatedAt = ?
	WHERE TaskUuid = ?`

	res := ExecWrapper(db.Conn(), query,
		status, boolToInt(hasFailed), boolToInt(isActive), completedAt, now, taskUuid,
	)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "ai_analysis.UpdateTaskStatus")
	}

	return nil
}

// GetTaskFiles retrieves all AiTaskFile records for a given task ID.
func (db *AiAnalysisSplitDB) GetTaskFiles(taskId int64) ([]AiTaskFile, error) {
	query := `SELECT 
		AiTaskFileId, AiTaskId, RelPath, AbsPath, Action, BeforeSha256, AfterSha256,
		BackupPath, LineCount, Reasoning, IsModified, IsRemoved, Notes, Comments, CreatedAt, UpdatedAt
	FROM AiTaskFile WHERE AiTaskId = ? ORDER BY AiTaskFileId ASC`

	rows, err := db.Conn().Query(query, taskId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis.GetTaskFiles")
	}
	defer rows.Close()

	return scanTaskFiles(rows)
}

func scanTaskFiles(rows *sql.Rows) ([]AiTaskFile, error) {
	var files []AiTaskFile
	for rows.Next() {
		var f AiTaskFile
		var isModifiedInt, isRemovedInt int
		err := rows.Scan(
			&f.AiTaskFileId, &f.AiTaskId, &f.RelPath, &f.AbsPath, &f.Action, &f.BeforeSha256, &f.AfterSha256,
			&f.BackupPath, &f.LineCount, &f.Reasoning, &isModifiedInt, &isRemovedInt,
			&f.Notes, &f.Comments, &f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "ai_analysis.scanTaskFiles")
		}
		f.IsModified = isModifiedInt == 1
		f.IsRemoved = isRemovedInt == 1
		files = append(files, f)
	}

	return files, nil
}

// GetTaskFileByRelPath looks up a file record by task ID and relative path.
func (db *AiAnalysisSplitDB) GetTaskFileByRelPath(taskId int64, relPath string) (*AiTaskFile, error) {
	cleanRel := filepath.ToSlash(filepath.Clean(relPath))
	query := `SELECT 
		AiTaskFileId, AiTaskId, RelPath, AbsPath, Action, BeforeSha256, AfterSha256,
		BackupPath, LineCount, Reasoning, IsModified, IsRemoved, Notes, Comments, CreatedAt, UpdatedAt
	FROM AiTaskFile WHERE AiTaskId = ? AND RelPath = ? ORDER BY AiTaskFileId DESC LIMIT 1`

	row := db.Conn().QueryRow(query, taskId, cleanRel)
	var f AiTaskFile
	var isModifiedInt, isRemovedInt int
	err := row.Scan(
		&f.AiTaskFileId, &f.AiTaskId, &f.RelPath, &f.AbsPath, &f.Action, &f.BeforeSha256, &f.AfterSha256,
		&f.BackupPath, &f.LineCount, &f.Reasoning, &isModifiedInt, &isRemovedInt,
		&f.Notes, &f.Comments, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis.GetTaskFileByRelPath")
	}
	f.IsModified = isModifiedInt == 1
	f.IsRemoved = isRemovedInt == 1

	return &f, nil
}

// GetTaskLines retrieves all lines recorded for a specific file ID.
func (db *AiAnalysisSplitDB) GetTaskLines(fileId int64) ([]AiTaskLine, error) {
	query := `SELECT 
		AiTaskLineId, AiTaskFileId, LineNumber, OriginalContent, ProposedContent,
		DiffKind, RuleViolation, Reasoning, IsApplied, Notes, Comments, CreatedAt
	FROM AiTaskLine WHERE AiTaskFileId = ? ORDER BY LineNumber ASC`

	rows, err := db.Conn().Query(query, fileId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis.GetTaskLines")
	}
	defer rows.Close()

	return scanTaskLines(rows)
}

func scanTaskLines(rows *sql.Rows) ([]AiTaskLine, error) {
	var lines []AiTaskLine
	for rows.Next() {
		var l AiTaskLine
		var isAppliedInt int
		err := rows.Scan(
			&l.AiTaskLineId, &l.AiTaskFileId, &l.LineNumber, &l.OriginalContent, &l.ProposedContent,
			&l.DiffKind, &l.RuleViolation, &l.Reasoning, &isAppliedInt,
			&l.Notes, &l.Comments, &l.CreatedAt,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "ai_analysis.scanTaskLines")
		}
		l.IsApplied = isAppliedInt == 1
		lines = append(lines, l)
	}

	return lines, nil
}

// GetSummary calculates aggregate metrics across all AI analysis tasks.
func (db *AiAnalysisSplitDB) GetSummary() (*AiAnalysisSummaryData, error) {
	var summary AiAnalysisSummaryData
	queryTasks := `SELECT 
		COUNT(*),
		COALESCE(SUM(CASE WHEN IsActive = 1 THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN Status = 'completed' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN HasFailed = 1 THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(TotalFiles), 0),
		COALESCE(SUM(TotalLines), 0)
	FROM AiTask`

	row := db.Conn().QueryRow(queryTasks)
	err := row.Scan(
		&summary.TotalTasks,
		&summary.ActiveTasks,
		&summary.CompletedTasks,
		&summary.FailedTasks,
		&summary.TotalFiles,
		&summary.TotalLines,
	)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_analysis.GetSummary")
	}

	return &summary, nil
}
