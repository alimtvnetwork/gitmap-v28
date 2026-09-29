// Package store — watch_prompts_split_db.go manages SQLite Split-DB storage for watch-prompts-running.
package store

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	sqlCreateWatchPromptRecord = `CREATE TABLE IF NOT EXISTS WatchPromptRecord (
    RecordId       TEXT PRIMARY KEY,
    RepoSlug       TEXT NOT NULL,
    ProjectName    TEXT NOT NULL,
    ProjectPath    TEXT NOT NULL,
    ProjectId      TEXT NOT NULL,
    SequenceId     TEXT NOT NULL DEFAULT '',
    ConversationId TEXT NOT NULL DEFAULT '',
    PromptText     TEXT NOT NULL,
    PromptStatus   TEXT NOT NULL,
    WordCount      INTEGER NOT NULL DEFAULT 0,
    MediaPathsJson TEXT NOT NULL DEFAULT '[]',
    IsActive       INTEGER NOT NULL DEFAULT 1,
    UpdatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_watch_prompt_repo ON WatchPromptRecord(RepoSlug);
CREATE INDEX IF NOT EXISTS idx_watch_prompt_active ON WatchPromptRecord(IsActive);`

	sqlCreateWatchConfig = `CREATE TABLE IF NOT EXISTS WatchConfig (
    ConfigKey   TEXT PRIMARY KEY,
    ConfigValue TEXT NOT NULL,
    UpdatedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	sqlCreateWatchLog = `CREATE TABLE IF NOT EXISTS WatchLog (
    LogId     INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoSlug  TEXT NOT NULL DEFAULT '',
    Event     TEXT NOT NULL,
    Message   TEXT NOT NULL,
    Status    TEXT NOT NULL DEFAULT 'info',
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_watch_log_repo ON WatchLog(RepoSlug);`
)

// WatchPromptsSplitDB manages dedicated SQLite storage for watched prompts in a repo slug.
type WatchPromptsSplitDB struct {
	conn *sql.DB
	path string
	slug string
}

// ResolveWatchPromptsRootDir resolves the root watch-prompts directory beside data folder.
func ResolveWatchPromptsRootDir() string {
	base := filepath.Dir(BinaryDataDir())
	standard := filepath.Join(base, "watch-prompts")
	if isDirExisting(standard) {
		return standard
	}

	spaceNamed := filepath.Join(base, "watch prompts")
	if isDirExisting(spaceNamed) {
		return spaceNamed
	}

	return standard
}

// ResolveWatchPromptsDbPath resolves the canonical database path for a repo slug.
func ResolveWatchPromptsDbPath(repoSlug string) string {
	cleanSlug := SanitizeSlug(repoSlug)
	rootDir := ResolveWatchPromptsRootDir()
	return filepath.ToSlash(filepath.Join(rootDir, cleanSlug, "sql.db"))
}

// OpenWatchPromptsSplitDB opens or creates the split database for a given repo slug.
func OpenWatchPromptsSplitDB(repoSlug string) (*WatchPromptsSplitDB, error) {
	dbPath := ResolveWatchPromptsDbPath(repoSlug)
	mkErr := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if mkErr != nil {
		return nil, apperror.WrapSimple(mkErr, "mkdir watch-prompts db")
	}

	conn, openErr := sql.Open("sqlite", dbPath)
	if openErr != nil {
		return nil, apperror.WrapSimple(openErr, "open watch-prompts sqlite")
	}

	return initWatchPromptsConn(conn, dbPath, SanitizeSlug(repoSlug))
}

func initWatchPromptsConn(conn *sql.DB, dbPath, slug string) (*WatchPromptsSplitDB, error) {
	cfgErr := ConfigureSQLiteConn(conn)
	if cfgErr != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(cfgErr, "configure watch-prompts sqlite")
	}

	db := &WatchPromptsSplitDB{conn: conn, path: dbPath, slug: slug}
	initErr := db.InitSchema()
	if initErr != nil {
		_ = conn.Close()
		return nil, initErr
	}

	return db, nil
}

// InitSchema creates the prompt record, config, and log tables.
func (db *WatchPromptsSplitDB) InitSchema() error {
	stmts := []string{sqlCreateWatchPromptRecord, sqlCreateWatchConfig, sqlCreateWatchLog}
	for _, stmt := range stmts {
		res := ExecWrapper(db.conn, stmt)
		if res.IsFailure {
			return apperror.WrapSimple(res.Error, "init watch-prompts schema")
		}
	}

	return nil
}

// Close closes the underlying database connection.
func (db *WatchPromptsSplitDB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}

	return db.conn.Close()
}

// Path returns the database file path.
func (db *WatchPromptsSplitDB) Path() string {
	return db.path
}

// RepoSlug returns the associated repository slug.
func (db *WatchPromptsSplitDB) RepoSlug() string {
	return db.slug
}

// PrunePreviousPromptRecords removes older backup prompts so only recent prompt is preserved.
func (db *WatchPromptsSplitDB) PrunePreviousPromptRecords(repoSlug string) error {
	cleanSlug := SanitizeSlug(repoSlug)
	query := `DELETE FROM WatchPromptRecord WHERE RepoSlug = ?`
	res := ExecWrapper(db.conn, query, cleanSlug)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "prune previous watch prompts")
	}

	return nil
}

// SavePromptRecord stores the prompt record after serializing media paths.
func (db *WatchPromptsSplitDB) SavePromptRecord(rec WatchPromptRecord) error {
	mediaJson, _ := json.Marshal(rec.MediaPaths)
	isActiveVal := 0
	if rec.IsActive {
		isActiveVal = 1
	}

	query := `INSERT OR REPLACE INTO WatchPromptRecord (
		RecordId, RepoSlug, ProjectName, ProjectPath, ProjectId,
		SequenceId, ConversationId, PromptText, PromptStatus,
		WordCount, MediaPathsJson, IsActive, UpdatedAt, CreatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res := ExecWrapper(
		db.conn, query,
		rec.RecordId, rec.RepoSlug, rec.ProjectName, rec.ProjectPath, rec.ProjectId,
		rec.SequenceId, rec.ConversationId, rec.PromptText, rec.PromptStatus,
		rec.WordCount, string(mediaJson), isActiveVal, rec.UpdatedAt, rec.CreatedAt,
	)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "save watch prompt record")
	}

	return nil
}

// ListPromptRecords returns all prompt records for the current repo slug.
func (db *WatchPromptsSplitDB) ListPromptRecords(repoSlug string) ([]WatchPromptRecord, error) {
	cleanSlug := SanitizeSlug(repoSlug)
	query := `SELECT RecordId, RepoSlug, ProjectName, ProjectPath, ProjectId,
		SequenceId, ConversationId, PromptText, PromptStatus,
		WordCount, MediaPathsJson, IsActive, UpdatedAt, CreatedAt
		FROM WatchPromptRecord WHERE RepoSlug = ? ORDER BY UpdatedAt DESC`

	rows, err := db.conn.Query(query, cleanSlug)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query watch prompt records")
	}
	defer rows.Close()

	return scanWatchPromptRecords(rows)
}

func scanWatchPromptRecords(rows *sql.Rows) ([]WatchPromptRecord, error) {
	var list []WatchPromptRecord
	for rows.Next() {
		rec, err := scanSingleWatchPromptRecord(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, rec)
	}

	return list, nil
}

func scanSingleWatchPromptRecord(rows *sql.Rows) (WatchPromptRecord, error) {
	var rec WatchPromptRecord
	var mediaJson string
	var isActiveInt int

	err := rows.Scan(
		&rec.RecordId, &rec.RepoSlug, &rec.ProjectName, &rec.ProjectPath, &rec.ProjectId,
		&rec.SequenceId, &rec.ConversationId, &rec.PromptText, &rec.PromptStatus,
		&rec.WordCount, &mediaJson, &isActiveInt, &rec.UpdatedAt, &rec.CreatedAt,
	)
	if err != nil {
		return rec, apperror.WrapSimple(err, "scan watch prompt record")
	}

	rec.IsActive = isActiveInt == 1
	_ = json.Unmarshal([]byte(mediaJson), &rec.MediaPaths)

	return rec, nil
}

// InsertWatchLog logs an event into the watch log table.
func (db *WatchPromptsSplitDB) InsertWatchLog(repoSlug, event, message, status string) error {
	query := `INSERT INTO WatchLog (RepoSlug, Event, Message, Status, CreatedAt) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`
	res := ExecWrapper(db.conn, query, SanitizeSlug(repoSlug), event, message, status)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "insert watch log")
	}

	return nil
}

// ListWatchLogs returns recent log entries for the current repo slug.
func (db *WatchPromptsSplitDB) ListWatchLogs(repoSlug string, limit int) ([]WatchLogRecord, error) {
	maxRows := limit
	if maxRows <= 0 {
		maxRows = 50
	}

	query := `SELECT LogId, RepoSlug, Event, Message, Status, CreatedAt FROM WatchLog WHERE RepoSlug = ? OR ? = '' ORDER BY LogId DESC LIMIT ?`
	cleanSlug := SanitizeSlug(repoSlug)
	rows, err := db.conn.Query(query, cleanSlug, cleanSlug, maxRows)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query watch logs")
	}
	defer rows.Close()

	return scanWatchLogRecords(rows)
}

func scanWatchLogRecords(rows *sql.Rows) ([]WatchLogRecord, error) {
	var list []WatchLogRecord
	for rows.Next() {
		var l WatchLogRecord
		err := rows.Scan(&l.LogId, &l.RepoSlug, &l.Event, &l.Message, &l.Status, &l.CreatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scan watch log")
		}
		list = append(list, l)
	}

	return list, nil
}

// SaveWatchConfig persists a key-value setting in WatchConfig.
func (db *WatchPromptsSplitDB) SaveWatchConfig(key, val string) error {
	query := `INSERT OR REPLACE INTO WatchConfig (ConfigKey, ConfigValue, UpdatedAt) VALUES (?, ?, CURRENT_TIMESTAMP)`
	res := ExecWrapper(db.conn, query, strings.TrimSpace(key), strings.TrimSpace(val))
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "save watch config")
	}

	return nil
}

// GetWatchConfig retrieves a key value from WatchConfig.
func (db *WatchPromptsSplitDB) GetWatchConfig(key string) (string, error) {
	query := `SELECT ConfigValue FROM WatchConfig WHERE ConfigKey = ?`
	var val string
	err := db.conn.QueryRow(query, strings.TrimSpace(key)).Scan(&val)
	if err != nil {
		return "", err
	}

	return val, nil
}
