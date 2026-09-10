package store

import (
	"database/sql"
)

// PurgeHistoryLog represents a record of a git purge operation.
type PurgeHistoryLog struct {
	PurgeHistoryLogId int64  `json:"purgeHistoryLogId"`
	RepoPath          string `json:"repoPath"`
	Pattern           string `json:"pattern"`
	BackupBranch      string `json:"backupBranch"`
	TempDir           string `json:"tempDir"`
	Files             string `json:"files"`
	Timestamp         int64  `json:"timestamp"`
	IsRestored        bool   `json:"isRestored"`
	Notes             string `json:"notes,omitempty"`
	Comments          string `json:"comments,omitempty"`
}

const sqlCreatePurgeHistory = `CREATE TABLE IF NOT EXISTS PurgeHistoryLog (
	PurgeHistoryLogId INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoPath TEXT NOT NULL,
	Pattern TEXT NOT NULL,
	BackupBranch TEXT NOT NULL,
	TempDir TEXT NOT NULL,
	Files TEXT NOT NULL,
	Timestamp INTEGER NOT NULL,
	IsRestored INTEGER NOT NULL DEFAULT 0,
	Notes TEXT NULL,
	Comments TEXT NULL
)`

const sqlInsertPurgeHistory = `INSERT INTO PurgeHistoryLog (RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, IsRestored, Notes, Comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

// InsertPurgeHistoryLog inserts a new purge history entry.
func (db *DB) InsertPurgeHistoryLog(log *PurgeHistoryLog) error {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return err
	}
	res, err := ExecWrapper(db.conn, sqlInsertPurgeHistory,
		log.RepoPath, log.Pattern, log.BackupBranch, log.TempDir, log.Files, log.Timestamp, boolToInt(log.IsRestored), log.Notes, log.Comments,
	).Destruct()
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	log.PurgeHistoryLogId = id
	return nil
}

// GetLastPurgeHistoryLog retrieves the most recent unrestored purge history log for a repo.
func (db *DB) GetLastPurgeHistoryLog(repoPath string) (*PurgeHistoryLog, error) {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return nil, err
	}

	row := QueryRowWrapper(
		db.conn,
		`SELECT PurgeHistoryLogId, RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, IsRestored, COALESCE(Notes, ''), COALESCE(Comments, '') FROM PurgeHistoryLog WHERE RepoPath = ? AND IsRestored = 0 ORDER BY PurgeHistoryLogId DESC LIMIT 1`,
		repoPath,
	)
	return scanPurgeHistoryRow(row)
}

// GetPurgeHistoryLogById retrieves a purge history log by its ID.
func (db *DB) GetPurgeHistoryLogById(id int64) (*PurgeHistoryLog, error) {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return nil, err
	}

	row := QueryRowWrapper(
		db.conn,
		`SELECT PurgeHistoryLogId, RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, IsRestored, COALESCE(Notes, ''), COALESCE(Comments, '') FROM PurgeHistoryLog WHERE PurgeHistoryLogId = ?`,
		id,
	)
	return scanPurgeHistoryRow(row)
}

func scanPurgeHistoryRow(row *sql.Row) (*PurgeHistoryLog, error) {
	var log PurgeHistoryLog
	var isRestoredInt int
	err := row.Scan(
		&log.PurgeHistoryLogId,
		&log.RepoPath,
		&log.Pattern,
		&log.BackupBranch,
		&log.TempDir,
		&log.Files,
		&log.Timestamp,
		&isRestoredInt,
		&log.Notes,
		&log.Comments,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	log.IsRestored = isRestoredInt == 1
	return &log, nil
}

// MarkPurgeHistoryRestored marks a purge history log as restored.
func (db *DB) MarkPurgeHistoryRestored(id int64) error {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return err
	}

	_, err := ExecWrapper(db.conn, `UPDATE PurgeHistoryLog SET IsRestored = 1 WHERE PurgeHistoryLogId = ?`, id).Destruct()
	return err
}
