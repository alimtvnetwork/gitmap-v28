package store

import (
	"database/sql"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	maxLogOutputBytes = 65536
	defaultLogLimit   = 50

	sqlInsertInstallationLog = `INSERT INTO InstallationLog
(Tool, Action, Version, PackageManager, DurationMs, IsSuccess, ExitCode, Stdout, Stderr, CommandLine, Notes, Comments, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlSelectInstallationLogs = `SELECT InstallationLogId, Tool, Action, Version, PackageManager, DurationMs, IsSuccess, ExitCode, Stdout, Stderr, CommandLine, Notes, Comments, CreatedAt
FROM InstallationLog ORDER BY InstallationLogId DESC LIMIT ?;`

	sqlSelectLogsByTool = `SELECT InstallationLogId, Tool, Action, Version, PackageManager, DurationMs, IsSuccess, ExitCode, Stdout, Stderr, CommandLine, Notes, Comments, CreatedAt
FROM InstallationLog WHERE Tool = ? ORDER BY InstallationLogId DESC LIMIT ?;`

	sqlSelectFailedLogs = `SELECT InstallationLogId, Tool, Action, Version, PackageManager, DurationMs, IsSuccess, ExitCode, Stdout, Stderr, CommandLine, Notes, Comments, CreatedAt
FROM InstallationLog WHERE IsSuccess = 0 ORDER BY InstallationLogId DESC LIMIT ?;`
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

// IsFailed reports whether the installation log record represents a failure.
func (r InstallationLogRecord) IsFailed() bool {
	return !r.IsSuccess
}

// IsFail reports whether the installation log record represents a failure.
func (r InstallationLogRecord) IsFail() bool {
	return !r.IsSuccess
}

// SanitizeLogOutput truncates log output to 64 KB to prevent SQLite database bloat.
func SanitizeLogOutput(s string) string {
	if len(s) > maxLogOutputBytes {
		return strings.ToValidUTF8(s[:maxLogOutputBytes], "")
	}

	return s
}

// RecordLog inserts an execution audit log into installation.db.
func (s *InstallationSplitDB) RecordLog(r InstallationLogRecord) error {
	now := resolveCreatedAt(r.CreatedAt)
	successInt := boolToInt(r.IsSuccess)
	stdout := SanitizeLogOutput(r.Stdout)
	stderr := SanitizeLogOutput(r.Stderr)

	_, err := s.conn.Exec(sqlInsertInstallationLog, r.Tool, r.Action, r.Version,
		r.PackageManager, r.DurationMs, successInt, r.ExitCode,
		stdout, stderr, r.CommandLine, r.Notes, r.Comments, now)

	if err != nil {
		return apperror.WrapSimple(err, "installation_split.recordLog")
	}

	return nil
}

// RecordExecution inserts a sanitized execution telemetry record into installation.db.
func (s *InstallationSplitDB) RecordExecution(
	tool, action, version, manager string,
	durationMs int64, isSuccess bool, exitCode int,
	stdout, stderr, cmdLine, notes, comments string,
) error {
	rec := newExecutionRecord(tool, action, version, manager, durationMs, isSuccess, exitCode)
	rec.Stdout, rec.Stderr, rec.CommandLine = stdout, stderr, cmdLine
	rec.Notes, rec.Comments = notes, comments

	return s.RecordLog(rec)
}

// GetLogs returns the most recent installation logs up to limit.
func (s *InstallationSplitDB) GetLogs(limit int) ([]InstallationLogRecord, error) {
	effectiveLimit := resolveLogLimit(limit)

	rows, err := s.conn.Query(sqlSelectInstallationLogs, effectiveLimit)

	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.getLogs")
	}

	defer rows.Close()

	return scanInstallationLogs(rows)
}

// GetLogsByTool returns recent logs for a specific tool.
func (s *InstallationSplitDB) GetLogsByTool(tool string, limit int) ([]InstallationLogRecord, error) {
	effectiveLimit := resolveLogLimit(limit)

	rows, err := s.conn.Query(sqlSelectLogsByTool, tool, effectiveLimit)

	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.getLogsByTool")
	}

	defer rows.Close()

	return scanInstallationLogs(rows)
}

// GetFailedLogs returns the most recent failed installation logs up to limit.
func (s *InstallationSplitDB) GetFailedLogs(limit int) ([]InstallationLogRecord, error) {
	effectiveLimit := resolveLogLimit(limit)

	rows, err := s.conn.Query(sqlSelectFailedLogs, effectiveLimit)

	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.getFailedLogs")
	}

	defer rows.Close()

	return scanInstallationLogs(rows)
}

func scanInstallationLogs(rows *sql.Rows) ([]InstallationLogRecord, error) {
	var logs []InstallationLogRecord

	for rows.Next() {
		var r InstallationLogRecord

		err := scanSingleLogRecord(rows, &r)

		if err != nil {
			return nil, apperror.WrapSimple(err, "installation_split.scanLogs")
		}

		logs = append(logs, r)
	}

	return logs, nil
}

func scanSingleLogRecord(rows *sql.Rows, r *InstallationLogRecord) error {
	var successInt int

	err := rows.Scan(&r.InstallationLogId, &r.Tool, &r.Action, &r.Version,
		&r.PackageManager, &r.DurationMs, &successInt, &r.ExitCode,
		&r.Stdout, &r.Stderr, &r.CommandLine, &r.Notes, &r.Comments, &r.CreatedAt)

	if err != nil {
		return err
	}

	r.IsSuccess = successInt == 1

	return nil
}

func newExecutionRecord(
	tool, action, version, manager string,
	durationMs int64, isSuccess bool, exitCode int,
) InstallationLogRecord {
	return InstallationLogRecord{
		Tool:           tool,
		Action:         action,
		Version:        version,
		PackageManager: manager,
		DurationMs:     durationMs,
		IsSuccess:      isSuccess,
		ExitCode:       exitCode,
	}
}

func resolveCreatedAt(createdAt string) string {
	if createdAt != "" {
		return createdAt
	}

	return time.Now().UTC().Format(time.RFC3339)
}

func resolveLogLimit(limit int) int {
	if limit > 0 {
		return limit
	}

	return defaultLogLimit
}
