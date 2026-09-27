// Package store — errors_split_ops.go provides operations for recording and inspecting internal errors.
package store

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	sqlInsertInternalError = `INSERT INTO InternalErrorLog (
    ErrorCode, ErrorType, Command, Message, Details, SourceFile, ContextJson, StackTrace, GitMapVersion, IsResolved, Notes, Comments
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlSelectInternalErrors = `SELECT
    InternalErrorLogId, ErrorCode, ErrorType, Command, Message, COALESCE(Details, ''), COALESCE(SourceFile, ''),
    COALESCE(ContextJson, ''), COALESCE(StackTrace, ''), GitMapVersion, IsResolved, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt
FROM InternalErrorLog ORDER BY InternalErrorLogId DESC LIMIT ?`

	sqlSelectInternalErrorsUnresolved = `SELECT
    InternalErrorLogId, ErrorCode, ErrorType, Command, Message, COALESCE(Details, ''), COALESCE(SourceFile, ''),
    COALESCE(ContextJson, ''), COALESCE(StackTrace, ''), GitMapVersion, IsResolved, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt
FROM InternalErrorLog WHERE IsResolved = 0 ORDER BY InternalErrorLogId DESC LIMIT ?`

	sqlSelectInternalErrorByID = `SELECT
    InternalErrorLogId, ErrorCode, ErrorType, Command, Message, COALESCE(Details, ''), COALESCE(SourceFile, ''),
    COALESCE(ContextJson, ''), COALESCE(StackTrace, ''), GitMapVersion, IsResolved, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt
FROM InternalErrorLog WHERE InternalErrorLogId = ?`

	sqlClearInternalErrors  = `DELETE FROM InternalErrorLog`
	sqlResolveInternalError = `UPDATE InternalErrorLog SET IsResolved = 1 WHERE InternalErrorLogId = ?`
)

// RecordError inserts an internal error log entry directly using raw Exec.
func (s *ErrorsSplitDB) RecordError(rec InternalErrorRecord) (int64, error) {
	resolvedVal := 0
	if rec.IsResolved {
		resolvedVal = 1
	}

	res, err := s.conn.Exec(sqlInsertInternalError,
		rec.ErrorCode, rec.ErrorType, rec.Command, rec.Message,
		rec.Details, rec.SourceFile, rec.ContextJson, rec.StackTrace,
		rec.GitMapVersion, resolvedVal, rec.Notes, rec.Comments,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "errors_split.record")
	}

	return res.LastInsertId()
}

// ListErrors retrieves the most recent recorded internal errors.
func (s *ErrorsSplitDB) ListErrors(limit int, unresolvedOnly bool) ([]InternalErrorRecord, error) {
	if limit <= 0 {
		limit = 50
	}

	query := sqlSelectInternalErrors
	if unresolvedOnly {
		query = sqlSelectInternalErrorsUnresolved
	}

	rows, err := s.conn.Query(query, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "errors_split.list")
	}

	defer rows.Close()

	return scanErrorRecords(rows)
}

func scanErrorRecords(rows *sql.Rows) ([]InternalErrorRecord, error) {
	var records []InternalErrorRecord
	for rows.Next() {
		var rec InternalErrorRecord
		var isResolvedInt int
		err := rows.Scan(
			&rec.ID, &rec.ErrorCode, &rec.ErrorType, &rec.Command, &rec.Message,
			&rec.Details, &rec.SourceFile, &rec.ContextJson, &rec.StackTrace,
			&rec.GitMapVersion, &isResolvedInt, &rec.Notes, &rec.Comments, &rec.CreatedAt,
		)
		if err != nil {
			continue
		}

		rec.IsResolved = isResolvedInt > 0
		records = append(records, rec)
	}

	return records, nil
}

// GetError retrieves a single internal error by its ID.
func (s *ErrorsSplitDB) GetError(id int64) (*InternalErrorRecord, error) {
	row := s.conn.QueryRow(sqlSelectInternalErrorByID, id)
	var rec InternalErrorRecord
	var isResolvedInt int

	err := row.Scan(
		&rec.ID, &rec.ErrorCode, &rec.ErrorType, &rec.Command, &rec.Message,
		&rec.Details, &rec.SourceFile, &rec.ContextJson, &rec.StackTrace,
		&rec.GitMapVersion, &isResolvedInt, &rec.Notes, &rec.Comments, &rec.CreatedAt,
	)
	if err != nil {
		return nil, apperror.WrapSimple(err, "errors_split.get")
	}

	rec.IsResolved = isResolvedInt > 0

	return &rec, nil
}

// ClearErrors deletes all recorded internal errors.
func (s *ErrorsSplitDB) ClearErrors() error {
	_, err := s.conn.Exec(sqlClearInternalErrors)
	if err != nil {
		return apperror.WrapSimple(err, "errors_split.clear")
	}

	return nil
}

// ResolveError marks an error record as resolved.
func (s *ErrorsSplitDB) ResolveError(id int64) error {
	_, err := s.conn.Exec(sqlResolveInternalError, id)
	if err != nil {
		return apperror.WrapSimple(err, "errors_split.resolve")
	}

	return nil
}

// LogInternalError is a safe, non-panicking global helper to log internal errors.
func LogInternalError(errType, code, message, details, sourceFile string) {
	rec := InternalErrorRecord{
		ErrorCode:     code,
		ErrorType:     errType,
		Message:       message,
		Details:       details,
		SourceFile:    sourceFile,
		GitMapVersion: constants.Version,
	}

	LogInternalErrorRecord(rec)
}

// LogInternalErrorRecord writes an InternalErrorRecord to gitmap-errors.db safely.
func LogInternalErrorRecord(rec InternalErrorRecord) {
	if strings.TrimSpace(rec.Message) == "" {
		return
	}

	db, err := OpenErrorsSplitDB()
	if err != nil {
		return
	}

	defer db.Close()

	if rec.GitMapVersion == "" {
		rec.GitMapVersion = constants.Version
	}

	_, _ = db.RecordError(rec)
}
