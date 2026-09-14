package store

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	sqlInsertInstallStart = `INSERT INTO install_logs
(id, target_type, target_name, action, status, exit_code, error_message, log_path, started_at, ended_at)
VALUES (?, ?, ?, ?, 'running', 0, '', '', ?, '');`

	sqlUpdateInstallSuccess = `UPDATE install_logs
SET status = 'success', exit_code = 0, ended_at = ?
WHERE id = (
    SELECT id FROM install_logs
    WHERE target_type = ? AND target_name = ? AND status = 'running'
    ORDER BY started_at DESC, rowid DESC LIMIT 1
);`

	sqlInsertInstallSuccess = `INSERT INTO install_logs
(id, target_type, target_name, action, status, exit_code, error_message, log_path, started_at, ended_at)
VALUES (?, ?, ?, 'install', 'success', 0, '', '', ?, ?);`

	sqlUpdateInstallFailure = `UPDATE install_logs
SET status = 'failed', exit_code = ?, error_message = ?, ended_at = ?
WHERE id = (
    SELECT id FROM install_logs
    WHERE target_type = ? AND target_name = ? AND status = 'running'
    ORDER BY started_at DESC, rowid DESC LIMIT 1
);`

	sqlInsertInstallFailure = `INSERT INTO install_logs
(id, target_type, target_name, action, status, exit_code, error_message, log_path, started_at, ended_at)
VALUES (?, ?, ?, 'install', 'failed', ?, ?, '', ?, ?);`

	sqlUpdateInstallSkipped = `UPDATE install_logs
SET status = 'skipped', exit_code = 0, error_message = ?, ended_at = ?
WHERE id = (
    SELECT id FROM install_logs
    WHERE target_type = ? AND target_name = ? AND status = 'running'
    ORDER BY started_at DESC, rowid DESC LIMIT 1
);`

	sqlInsertInstallSkipped = `INSERT INTO install_logs
(id, target_type, target_name, action, status, exit_code, error_message, log_path, started_at, ended_at)
VALUES (?, ?, ?, 'install', 'skipped', 0, ?, '', ?, ?);`

	sqlSelectInstallLogByID = `SELECT id, target_type, target_name, action, status, exit_code,
COALESCE(error_message, ''), COALESCE(log_path, ''), started_at, COALESCE(ended_at, '')
FROM install_logs WHERE id = ?;`

	sqlSelectListInstallLogsByType = `SELECT id, target_type, target_name, action, status, exit_code,
COALESCE(error_message, ''), COALESCE(log_path, ''), started_at, COALESCE(ended_at, '')
FROM install_logs WHERE target_type = ? ORDER BY started_at DESC, rowid DESC LIMIT ?;`

	sqlSelectListAllInstallLogs = `SELECT id, target_type, target_name, action, status, exit_code,
COALESCE(error_message, ''), COALESCE(log_path, ''), started_at, COALESCE(ended_at, '')
FROM install_logs ORDER BY started_at DESC, rowid DESC LIMIT ?;`
)

// RecordInstallStart records the beginning of an installation attempt.
func (s *InstallationSplitDB) RecordInstallStart(targetType, targetName, action string) (string, error) {
	id := generateInstallLogID()
	now := time.Now().UTC().Format(time.RFC3339)
	act := resolveInstallAction(action)
	_, err := s.conn.Exec(sqlInsertInstallStart, id, targetType, targetName, act, now)
	if err != nil {
		return "", apperror.WrapSimple(err, "installation_split.recordInstallStart")
	}

	return id, nil
}

func resolveInstallAction(action string) string {
	if action != "" {
		return action
	}

	return "install"
}

// RecordInstallSuccess records the successful completion of an installation.
func (s *InstallationSplitDB) RecordInstallSuccess(targetType, targetName string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	hasUpdated, err := s.updateSuccessRecord(targetType, targetName, now)
	if err != nil {
		return err
	}
	if hasUpdated {
		return nil
	}

	return s.insertSuccessFallback(targetType, targetName, now)
}

func (s *InstallationSplitDB) updateSuccessRecord(targetType, targetName, now string) (bool, error) {
	res, err := s.conn.Exec(sqlUpdateInstallSuccess, now, targetType, targetName)
	if err != nil {
		return false, apperror.WrapSimple(err, "installation_split.updateSuccess")
	}

	return checkRowsUpdated(res, "installation_split.successRows")
}

func checkRowsUpdated(res sql.Result, op string) (bool, error) {
	affected, err := res.RowsAffected()
	if err != nil {
		return false, apperror.WrapSimple(err, op)
	}

	return affected > 0, nil
}

func (s *InstallationSplitDB) insertSuccessFallback(targetType, targetName, now string) error {
	id := generateInstallLogID()
	_, err := s.conn.Exec(sqlInsertInstallSuccess, id, targetType, targetName, now, now)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.insertSuccess")
	}

	return nil
}

// RecordInstallFailure records a failed installation with an exit code and error message.
func (s *InstallationSplitDB) RecordInstallFailure(targetType, targetName string, exitCode int, errorMsg string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	hasUpdated, err := s.updateFailureRecord(targetType, targetName, errorMsg, now, exitCode)
	if err != nil {
		return err
	}
	if hasUpdated {
		return nil
	}

	return s.insertFailureFallback(targetType, targetName, errorMsg, now, exitCode)
}

func (s *InstallationSplitDB) updateFailureRecord(targetType, targetName, errorMsg, now string, exitCode int) (bool, error) {
	res, err := s.conn.Exec(sqlUpdateInstallFailure, exitCode, errorMsg, now, targetType, targetName)
	if err != nil {
		return false, apperror.WrapSimple(err, "installation_split.updateFailure")
	}

	return checkRowsUpdated(res, "installation_split.failureRows")
}

func (s *InstallationSplitDB) insertFailureFallback(targetType, targetName, errorMsg, now string, exitCode int) error {
	id := generateInstallLogID()
	_, err := s.conn.Exec(sqlInsertInstallFailure, id, targetType, targetName, exitCode, errorMsg, now, now)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.insertFailure")
	}

	return nil
}

// RecordInstallSkipped records a skipped installation with a reason message.
func (s *InstallationSplitDB) RecordInstallSkipped(targetType, targetName, reason string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	hasUpdated, err := s.updateSkippedRecord(targetType, targetName, reason, now)
	if err != nil {
		return err
	}
	if hasUpdated {
		return nil
	}

	return s.insertSkippedFallback(targetType, targetName, reason, now)
}

func (s *InstallationSplitDB) updateSkippedRecord(targetType, targetName, reason, now string) (bool, error) {
	res, err := s.conn.Exec(sqlUpdateInstallSkipped, reason, now, targetType, targetName)
	if err != nil {
		return false, apperror.WrapSimple(err, "installation_split.updateSkipped")
	}

	return checkRowsUpdated(res, "installation_split.skippedRows")
}

func (s *InstallationSplitDB) insertSkippedFallback(targetType, targetName, reason, now string) error {
	id := generateInstallLogID()
	_, err := s.conn.Exec(sqlInsertInstallSkipped, id, targetType, targetName, reason, now, now)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.insertSkipped")
	}

	return nil
}

// GetInstallLog retrieves a single installation log record by its ID.
func (s *InstallationSplitDB) GetInstallLog(id string) (*InstallLogRecord, error) {
	row := s.conn.QueryRow(sqlSelectInstallLogByID, id)
	rec, err := scanInstallLogRow(row)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.getInstallLog")
	}

	return rec, nil
}

// ListInstallLogs returns installation log records, optionally filtered by targetType, up to limit.
func (s *InstallationSplitDB) ListInstallLogs(targetType string, limit int) ([]InstallLogRecord, error) {
	effectiveLimit := resolveInstallLogsLimit(limit)
	if targetType != "" {
		return s.queryInstallLogsByType(targetType, effectiveLimit)
	}

	return s.queryAllInstallLogs(effectiveLimit)
}

func resolveInstallLogsLimit(limit int) int {
	if limit > 0 {
		return limit
	}

	return 50
}

func (s *InstallationSplitDB) queryInstallLogsByType(targetType string, limit int) ([]InstallLogRecord, error) {
	rows, err := s.conn.Query(sqlSelectListInstallLogsByType, targetType, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.queryInstallLogsByType")
	}
	defer rows.Close()

	return scanInstallLogList(rows)
}

func (s *InstallationSplitDB) queryAllInstallLogs(limit int) ([]InstallLogRecord, error) {
	rows, err := s.conn.Query(sqlSelectListAllInstallLogs, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.queryAllInstallLogs")
	}
	defer rows.Close()

	return scanInstallLogList(rows)
}

func scanInstallLogList(rows *sql.Rows) ([]InstallLogRecord, error) {
	records := make([]InstallLogRecord, 0)
	for rows.Next() {
		rec, err := scanInstallLogRow(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "installation_split.scanInstallLogList")
		}
		records = append(records, *rec)
	}

	return records, nil
}

func scanInstallLogRow(scanner interface{ Scan(dest ...any) error }) (*InstallLogRecord, error) {
	var r InstallLogRecord
	err := scanner.Scan(
		&r.ID, &r.TargetType, &r.TargetName, &r.Action,
		&r.Status, &r.ExitCode, &r.ErrorMessage, &r.LogPath,
		&r.StartedAt, &r.EndedAt,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func generateInstallLogID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
