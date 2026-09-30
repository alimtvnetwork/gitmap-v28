// Package store — failed_commands_ops.go provides operations for recording and inspecting failed/unknown CLI commands.
package store

import (
	"database/sql"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	sqlFindFailedCommand = `SELECT FailedCommandId, HitCount FROM FailedCommand WHERE LOWER(Command) = LOWER(?) AND LOWER(Domain) = LOWER(?) LIMIT 1`

	sqlUpdateFailedCommand = `UPDATE FailedCommand SET
    FullArgs = ?,
    ErrorCode = ?,
    Message = ?,
    Suggestions = ?,
    HitCount = HitCount + 1,
    WorkingDir = ?,
    GitMapVersion = ?,
    IsResolved = 0,
    LastSeenAt = CURRENT_TIMESTAMP
WHERE FailedCommandId = ?`

	sqlInsertFailedCommand = `INSERT INTO FailedCommand (
    Command, FullArgs, Domain, ErrorCode, Message, Suggestions, HitCount, WorkingDir, GitMapVersion, IsResolved, Notes, Comments
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlSelectFailedCommands = `SELECT
    FailedCommandId, Command, FullArgs, Domain, ErrorCode, Message, Suggestions, HitCount,
    COALESCE(WorkingDir, ''), GitMapVersion, IsResolved, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt, LastSeenAt
FROM FailedCommand ORDER BY HitCount DESC, FailedCommandId DESC LIMIT ?`

	sqlCountFailedCommands = `SELECT COUNT(*), COALESCE(SUM(HitCount), 0) FROM FailedCommand`

	sqlClearFailedCommands = `DELETE FROM FailedCommand`
)

// FailedCommandRecord represents a recorded failed or undetected command entry in gitmap-errors.db.
type FailedCommandRecord struct {
	ID            int64  `json:"id"`
	Command       string `json:"command"`
	FullArgs      string `json:"fullArgs"`
	Domain        string `json:"domain"`
	ErrorCode     string `json:"errorCode"`
	Message       string `json:"message"`
	Suggestions   string `json:"suggestions"`
	HitCount      int64  `json:"hitCount"`
	WorkingDir    string `json:"workingDir,omitempty"`
	GitMapVersion string `json:"gitMapVersion,omitempty"`
	IsResolved    bool   `json:"isResolved"`
	Notes         string `json:"notes,omitempty"`
	Comments      string `json:"comments,omitempty"`
	CreatedAt     string `json:"createdAt"`
	LastSeenAt    string `json:"lastSeenAt"`
}

// FailedCommandSummary aggregates failed command counts and records for CLI reporting.
type FailedCommandSummary struct {
	TotalDistinctCommands int64                 `json:"totalDistinctCommands"`
	TotalFailedAttempts   int64                 `json:"totalFailedAttempts"`
	DatabasePath          string                `json:"databasePath"`
	Records               []FailedCommandRecord `json:"records"`
}

// RecordFailedCommand inserts or increments a failed command entry in FailedCommand.
func (s *ErrorsSplitDB) RecordFailedCommand(rec FailedCommandRecord) (int64, error) {
	rec = normalizeFailedCommandRecord(rec)

	var existingID int64
	var existingHits int64
	findErr := s.conn.QueryRow(sqlFindFailedCommand, rec.Command, rec.Domain).Scan(&existingID, &existingHits)
	if findErr == nil && existingID > 0 {
		return s.incrementFailedCommand(existingID, rec)
	}

	return s.insertFailedCommand(rec)
}

func normalizeFailedCommandRecord(rec FailedCommandRecord) FailedCommandRecord {
	rec.Command = strings.TrimSpace(rec.Command)
	if rec.Domain == "" {
		rec.Domain = "root"
	}
	if rec.ErrorCode == "" {
		rec.ErrorCode = "E1001"
	}
	if rec.HitCount <= 0 {
		rec.HitCount = 1
	}
	if rec.GitMapVersion == "" {
		rec.GitMapVersion = constants.Version
	}
	if rec.WorkingDir == "" {
		wd, _ := os.Getwd()
		rec.WorkingDir = wd
	}

	return rec
}

func (s *ErrorsSplitDB) incrementFailedCommand(existingID int64, rec FailedCommandRecord) (int64, error) {
	_, err := s.conn.Exec(sqlUpdateFailedCommand,
		rec.FullArgs, rec.ErrorCode, rec.Message, rec.Suggestions,
		rec.WorkingDir, rec.GitMapVersion, existingID,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "errors_split.update_failed_command")
	}

	return existingID, nil
}

func (s *ErrorsSplitDB) insertFailedCommand(rec FailedCommandRecord) (int64, error) {
	resolvedVal := 0
	if rec.IsResolved {
		resolvedVal = 1
	}

	res, err := s.conn.Exec(sqlInsertFailedCommand,
		rec.Command, rec.FullArgs, rec.Domain, rec.ErrorCode, rec.Message,
		rec.Suggestions, rec.HitCount, rec.WorkingDir, rec.GitMapVersion,
		resolvedVal, rec.Notes, rec.Comments,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "errors_split.insert_failed_command")
	}

	return res.LastInsertId()
}

// ListFailedCommands returns the top recorded failed commands ordered by hit count and recency.
func (s *ErrorsSplitDB) ListFailedCommands(limit int) ([]FailedCommandRecord, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.conn.Query(sqlSelectFailedCommands, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "errors_split.list_failed_commands")
	}

	defer rows.Close()

	return scanFailedCommandRecords(rows)
}

func scanFailedCommandRecords(rows *sql.Rows) ([]FailedCommandRecord, error) {
	var records []FailedCommandRecord
	for rows.Next() {
		var rec FailedCommandRecord
		var isResolvedInt int
		err := rows.Scan(
			&rec.ID, &rec.Command, &rec.FullArgs, &rec.Domain, &rec.ErrorCode,
			&rec.Message, &rec.Suggestions, &rec.HitCount, &rec.WorkingDir,
			&rec.GitMapVersion, &isResolvedInt, &rec.Notes, &rec.Comments,
			&rec.CreatedAt, &rec.LastSeenAt,
		)
		if err != nil {
			continue
		}

		rec.IsResolved = isResolvedInt > 0
		records = append(records, rec)
	}

	return records, nil
}

// CountFailedCommands returns both distinct failed command count and total failed attempts.
func (s *ErrorsSplitDB) CountFailedCommands() (int64, int64, error) {
	var distinctCount int64
	var totalHits int64
	err := s.conn.QueryRow(sqlCountFailedCommands).Scan(&distinctCount, &totalHits)
	if err != nil {
		return 0, 0, apperror.WrapSimple(err, "errors_split.count_failed_commands")
	}

	return distinctCount, totalHits, nil
}

// GetFailedCommandSummary returns aggregated statistics and records for failed commands.
func (s *ErrorsSplitDB) GetFailedCommandSummary(limit int) (FailedCommandSummary, error) {
	distinctCount, totalHits, err := s.CountFailedCommands()
	if err != nil {
		return FailedCommandSummary{}, err
	}

	records, listErr := s.ListFailedCommands(limit)
	if listErr != nil {
		return FailedCommandSummary{}, listErr
	}

	return FailedCommandSummary{
		TotalDistinctCommands: distinctCount,
		TotalFailedAttempts:   totalHits,
		DatabasePath:          s.Path,
		Records:               records,
	}, nil
}

// ClearFailedCommands deletes all entries from the FailedCommand table.
func (s *ErrorsSplitDB) ClearFailedCommands() error {
	_, err := s.conn.Exec(sqlClearFailedCommands)
	if err != nil {
		return apperror.WrapSimple(err, "errors_split.clear_failed_commands")
	}

	return nil
}

// LogFailedCommand safely records an unknown or failed-to-detect command into gitmap-errors.db.
func LogFailedCommand(command, fullArgs, domain, errorCode, message string, suggestions []string) {
	cleanCmd := strings.TrimSpace(command)
	if cleanCmd == "" {
		return
	}

	db, err := OpenErrorsSplitDB()
	if err != nil {
		return
	}

	defer db.Close()

	rec := FailedCommandRecord{
		Command:       cleanCmd,
		FullArgs:      strings.TrimSpace(fullArgs),
		Domain:        domain,
		ErrorCode:     errorCode,
		Message:       strings.TrimSpace(message),
		Suggestions:   strings.Join(suggestions, ", "),
		GitMapVersion: constants.Version,
	}

	_, _ = db.RecordFailedCommand(rec)
}
