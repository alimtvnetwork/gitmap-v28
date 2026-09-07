package store

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	sqlInsertInstallationLog = `INSERT INTO InstallationLog
(Tool, Action, Version, PackageManager, DurationMs, IsSuccess, ExitCode, Stdout, Stderr, CommandLine, Notes, Comments, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlSelectInstallationLogs = `SELECT InstallationLogId, Tool, Action, Version, PackageManager, DurationMs, IsSuccess, ExitCode, Stdout, Stderr, CommandLine, Notes, Comments, CreatedAt
FROM InstallationLog ORDER BY InstallationLogId DESC LIMIT ?;`

	sqlSelectLogsByTool = `SELECT InstallationLogId, Tool, Action, Version, PackageManager, DurationMs, IsSuccess, ExitCode, Stdout, Stderr, CommandLine, Notes, Comments, CreatedAt
FROM InstallationLog WHERE Tool = ? ORDER BY InstallationLogId DESC LIMIT ?;`
)

// InstallationLogRecord encapsulates execution telemetry for an installation event.
type InstallationLogRecord struct {
	InstallationLogId int64  `json:"installationLogId"`
	Tool              string `json:"tool"`
	Action            string `json:"action"`
	Version           string `json:"version"`
	PackageManager    string `json:"packageManager"`
	DurationMs        int64  `json:"durationMs"`
	IsSuccess         bool   `json:"isSuccess"`
	ExitCode          int    `json:"exitCode"`
	Stdout            string `json:"stdout,omitempty"`
	Stderr            string `json:"stderr,omitempty"`
	CommandLine       string `json:"commandLine,omitempty"`
	Notes             string `json:"notes,omitempty"`
	Comments          string `json:"comments,omitempty"`
	CreatedAt         string `json:"createdAt"`
}

// RecordLog inserts an execution audit log into installation.db.
func (s *InstallationSplitDB) RecordLog(r InstallationLogRecord) error {
	now := r.CreatedAt
	if now == "" {
		now = time.Now().UTC().Format(time.RFC3339)
	}

	successInt := 0
	if r.IsSuccess {
		successInt = 1
	}

	_, err := s.conn.Exec(sqlInsertInstallationLog, r.Tool, r.Action, r.Version,
		r.PackageManager, r.DurationMs, successInt, r.ExitCode,
		r.Stdout, r.Stderr, r.CommandLine, r.Notes, r.Comments, now)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.recordLog")
	}

	return nil
}

// GetLogs returns the most recent installation logs up to limit.
func (s *InstallationSplitDB) GetLogs(limit int) ([]InstallationLogRecord, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.conn.Query(sqlSelectInstallationLogs, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.getLogs")
	}
	defer rows.Close()

	return scanInstallationLogs(rows)
}

// GetLogsByTool returns recent logs for a specific tool.
func (s *InstallationSplitDB) GetLogsByTool(tool string, limit int) ([]InstallationLogRecord, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.conn.Query(sqlSelectLogsByTool, tool, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.getLogsByTool")
	}
	defer rows.Close()

	return scanInstallationLogs(rows)
}

func scanInstallationLogs(rows *sql.Rows) ([]InstallationLogRecord, error) {
	var logs []InstallationLogRecord
	for rows.Next() {
		var r InstallationLogRecord
		var successInt int
		err := rows.Scan(&r.InstallationLogId, &r.Tool, &r.Action, &r.Version,
			&r.PackageManager, &r.DurationMs, &successInt, &r.ExitCode,
			&r.Stdout, &r.Stderr, &r.CommandLine, &r.Notes, &r.Comments, &r.CreatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "installation_split.scanLogs")
		}
		r.IsSuccess = successInt == 1
		logs = append(logs, r)
	}

	return logs, nil
}
