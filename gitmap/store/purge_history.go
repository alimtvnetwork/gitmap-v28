package store

import (
	"database/sql"
)

type PurgeHistoryLog struct {
	ID           int64
	RepoPath     string
	Pattern      string
	BackupBranch string
	TempDir      string
	Files        string
	Timestamp    int64
	Restored     bool
}

const sqlCreatePurgeHistory = `CREATE TABLE IF NOT EXISTS PurgeHistoryLog (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoPath TEXT NOT NULL,
	Pattern TEXT NOT NULL,
	BackupBranch TEXT NOT NULL,
	TempDir TEXT NOT NULL,
	Files TEXT NOT NULL,
	Timestamp INTEGER NOT NULL,
	Restored BOOLEAN NOT NULL DEFAULT 0
)`

func (db *DB) EnsurePurgeHistoryTable() error {
	_, err := db.conn.Exec(sqlCreatePurgeHistory)
	return err
}

func (db *DB) InsertPurgeHistoryLog(log *PurgeHistoryLog) error {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return err
	}

	res, err := db.conn.Exec(
		`INSERT INTO PurgeHistoryLog (RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, Restored) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.RepoPath, log.Pattern, log.BackupBranch, log.TempDir, log.Files, log.Timestamp, log.Restored,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	log.ID = id
	return nil
}

func (db *DB) GetLastPurgeHistoryLog(repoPath string) (*PurgeHistoryLog, error) {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return nil, err
	}

	row := db.conn.QueryRow(`SELECT ID, RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, Restored FROM PurgeHistoryLog WHERE RepoPath = ? AND Restored = 0 ORDER BY ID DESC LIMIT 1`, repoPath)

	var log PurgeHistoryLog
	err := row.Scan(&log.ID, &log.RepoPath, &log.Pattern, &log.BackupBranch, &log.TempDir, &log.Files, &log.Timestamp, &log.Restored)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (db *DB) MarkPurgeHistoryRestored(id int64) error {
	_, err := db.conn.Exec(`UPDATE PurgeHistoryLog SET Restored = 1 WHERE ID = ?`, id)
	return err
}
