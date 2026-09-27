// Package store — agy_log_split_db.go manages SQLite Split-DB connections and schemas for AGY decision audit logs.
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	sqlCreateAgyDecisionLog = `CREATE TABLE IF NOT EXISTS AgyDecisionLog (
    LogId TEXT PRIMARY KEY,
    Command TEXT NOT NULL,
    TargetProject TEXT NOT NULL DEFAULT '',
    ProjectPath TEXT NOT NULL DEFAULT '',
    ConversationId TEXT NOT NULL DEFAULT '',
    DecisionReason TEXT NOT NULL DEFAULT '',
    Status TEXT NOT NULL DEFAULT 'success',
    Node TEXT NOT NULL DEFAULT 'local',
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_agy_decision_created ON AgyDecisionLog(CreatedAt DESC);
CREATE INDEX IF NOT EXISTS idx_agy_decision_project ON AgyDecisionLog(TargetProject);
CREATE INDEX IF NOT EXISTS idx_agy_decision_cmd ON AgyDecisionLog(Command);`
)

// AgyDecisionLogRecord represents an auditable decision trace for an AGY operation.
type AgyDecisionLogRecord struct {
	LogId          string `json:"logId"`
	Command        string `json:"command"`
	TargetProject  string `json:"targetProject"`
	ProjectPath    string `json:"projectPath"`
	ConversationId string `json:"conversationId"`
	DecisionReason string `json:"decisionReason"`
	Status         string `json:"status"`
	Node           string `json:"node"`
	CreatedAt      string `json:"createdAt"`
}

// AgyLogQueryOptions filters decision log query results.
type AgyLogQueryOptions struct {
	Project string
	Command string
	Node    string
	Limit   int
}

// AgyLogSplitDB manages dedicated SQLite storage for AGY decision logs.
type AgyLogSplitDB struct {
	conn *sql.DB
	path string
}

func resolveAgyLogDbPath(customPath string) string {
	if len(strings.TrimSpace(customPath)) > 0 {
		return customPath
	}
	return filepath.Join(BinaryDataDir(), "gitmap-agy-log.db")
}

// OpenAgyLogSplitDB opens or creates the split-db at customPath or default location.
func OpenAgyLogSplitDB(customPath string) (*AgyLogSplitDB, error) {
	dbPath := resolveAgyLogDbPath(customPath)
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "create agy log db dir")
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open agy log db")
	}
	ConfigureSQLiteConn(conn)
	if err := initAgyLogSchema(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &AgyLogSplitDB{conn: conn, path: dbPath}, nil
}

func initAgyLogSchema(conn *sql.DB) error {
	_, err := conn.Exec(sqlCreateAgyDecisionLog)
	if err != nil {
		return apperror.WrapSimple(err, "init agy log schema")
	}
	return nil
}

// Close closes the underlying database connection.
func (db *AgyLogSplitDB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// InsertAgyDecisionLog persists a decision record into the database.
func (db *AgyLogSplitDB) InsertAgyDecisionLog(rec AgyDecisionLogRecord) error {
	if rec.LogId == "" {
		rec.LogId = fmt.Sprintf("log-%d", time.Now().UnixNano())
	}
	if rec.CreatedAt == "" {
		rec.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	query := `INSERT INTO AgyDecisionLog (
		LogId, Command, TargetProject, ProjectPath, ConversationId, DecisionReason, Status, Node, CreatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query,
		rec.LogId, rec.Command, rec.TargetProject, rec.ProjectPath,
		rec.ConversationId, rec.DecisionReason, rec.Status, rec.Node, rec.CreatedAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "insert agy decision log")
	}
	return nil
}

// QueryAgyDecisionLogs returns filtered decision log records in reverse chronological order.
func (db *AgyLogSplitDB) QueryAgyDecisionLogs(opts AgyLogQueryOptions) ([]AgyDecisionLogRecord, error) {
	query, args := buildAgyLogQuery(opts)
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query agy decision logs")
	}
	defer rows.Close()
	return scanAgyLogRows(rows)
}

func buildAgyLogQuery(opts AgyLogQueryOptions) (string, []any) {
	query := `SELECT LogId, Command, TargetProject, ProjectPath, ConversationId, DecisionReason, Status, Node, CreatedAt
		FROM AgyDecisionLog WHERE 1=1`
	var args []any
	if opts.Project != "" {
		query += " AND (LOWER(TargetProject) LIKE ? OR LOWER(ProjectPath) LIKE ?)"
		pattern := "%" + strings.ToLower(opts.Project) + "%"
		args = append(args, pattern, pattern)
	}
	if opts.Command != "" {
		query += " AND LOWER(Command) = ?"
		args = append(args, strings.ToLower(opts.Command))
	}
	if opts.Node != "" {
		query += " AND LOWER(Node) = ?"
		args = append(args, strings.ToLower(opts.Node))
	}
	query += " ORDER BY CreatedAt DESC"
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	query += " LIMIT ?"
	args = append(args, limit)
	return query, args
}

func scanAgyLogRows(rows *sql.Rows) ([]AgyDecisionLogRecord, error) {
	var records []AgyDecisionLogRecord
	for rows.Next() {
		var rec AgyDecisionLogRecord
		err := rows.Scan(
			&rec.LogId, &rec.Command, &rec.TargetProject, &rec.ProjectPath,
			&rec.ConversationId, &rec.DecisionReason, &rec.Status, &rec.Node, &rec.CreatedAt,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scan agy log row")
		}
		records = append(records, rec)
	}
	return records, nil
}

// RecordAgyDecision writes an audit log entry using the default split-db.
func RecordAgyDecision(cmd, targetProj, projPath, convID, reason, status string) {
	db, err := OpenAgyLogSplitDB("")
	if err != nil {
		return
	}
	defer db.Close()
	nodeName, _ := os.Hostname()
	if nodeName == "" {
		nodeName = "local"
	}
	_ = db.InsertAgyDecisionLog(AgyDecisionLogRecord{
		Command:        cmd,
		TargetProject:  targetProj,
		ProjectPath:    projPath,
		ConversationId: convID,
		DecisionReason: reason,
		Status:         status,
		Node:           nodeName,
	})
}
