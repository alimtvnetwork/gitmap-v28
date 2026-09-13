// Package store — startup_split_ops.go: CRUD operations on StartupItem and StartupLog tables.
package store

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SaveStartupItem inserts or updates a startup item record in the split DB.
func (db *StartupSplitDB) SaveStartupItem(item *StartupItemRecord) error {
	now := time.Now().Unix()
	item.UpdatedAt = now
	if item.CreatedAt == 0 {
		item.CreatedAt = now
	}

	query := `INSERT INTO StartupItem (
		Name, TargetType, TargetPath, CommandArgs, IconPath,
		RunFrequency, IsActive, Description, CreatedAt, UpdatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(Name) DO UPDATE SET
		TargetType = excluded.TargetType,
		TargetPath = excluded.TargetPath,
		CommandArgs = excluded.CommandArgs,
		IconPath = excluded.IconPath,
		RunFrequency = excluded.RunFrequency,
		IsActive = excluded.IsActive,
		Description = excluded.Description,
		UpdatedAt = excluded.UpdatedAt;`

	activeInt := 0
	if item.IsActive {
		activeInt = 1
	}

	_, err := db.conn.Exec(query,
		item.Name, item.TargetType, item.TargetPath, item.CommandArgs,
		item.IconPath, item.RunFrequency, activeInt, item.Description,
		item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "save startup item "+item.Name)
	}

	return nil
}

// GetStartupItem retrieves a startup item by unique name.
func (db *StartupSplitDB) GetStartupItem(name string) (*StartupItemRecord, error) {
	query := `SELECT StartupItemId, Name, TargetType, TargetPath, CommandArgs,
		IconPath, RunFrequency, IsActive, Description, CreatedAt, UpdatedAt
	FROM StartupItem WHERE Name = ?;`

	row := db.conn.QueryRow(query, name)
	var it StartupItemRecord
	var activeInt int
	var args, icon, desc sql.NullString

	err := row.Scan(
		&it.StartupItemId, &it.Name, &it.TargetType, &it.TargetPath,
		&args, &icon, &it.RunFrequency, &activeInt, &desc,
		&it.CreatedAt, &it.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "get startup item "+name)
	}

	it.CommandArgs = args.String
	it.IconPath = icon.String
	it.Description = desc.String
	it.IsActive = activeInt == 1

	return &it, nil
}

// ListStartupItems returns all registered startup items.
func (db *StartupSplitDB) ListStartupItems() ([]StartupItemRecord, error) {
	query := `SELECT StartupItemId, Name, TargetType, TargetPath, CommandArgs,
		IconPath, RunFrequency, IsActive, Description, CreatedAt, UpdatedAt
	FROM StartupItem ORDER BY Name ASC;`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "list startup items")
	}
	defer rows.Close()

	var items []StartupItemRecord
	for rows.Next() {
		it, scanErr := scanStartupItemRow(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, it)
	}

	return items, nil
}

func scanStartupItemRow(rows *sql.Rows) (StartupItemRecord, error) {
	var it StartupItemRecord
	var activeInt int
	var args, icon, desc sql.NullString

	err := rows.Scan(
		&it.StartupItemId, &it.Name, &it.TargetType, &it.TargetPath,
		&args, &icon, &it.RunFrequency, &activeInt, &desc,
		&it.CreatedAt, &it.UpdatedAt,
	)
	if err != nil {
		return it, apperror.WrapSimple(err, "scan startup item row")
	}

	it.CommandArgs = args.String
	it.IconPath = icon.String
	it.Description = desc.String
	it.IsActive = activeInt == 1

	return it, nil
}

// DeleteStartupItem removes a startup item and its cascade logs.
func (db *StartupSplitDB) DeleteStartupItem(name string) error {
	query := `DELETE FROM StartupItem WHERE Name = ?;`
	_, err := db.conn.Exec(query, name)
	if err != nil {
		return apperror.WrapSimple(err, "delete startup item "+name)
	}

	return nil
}

// RecordStartupLog inserts a new execution log entry for a startup item.
func (db *StartupSplitDB) RecordStartupLog(log *StartupLogRecord) error {
	if log.RunAt == 0 {
		log.RunAt = time.Now().Unix()
	}

	successInt := 0
	if log.IsSuccess {
		successInt = 1
	}

	query := `INSERT INTO StartupLog (
		StartupItemId, RunAt, DurationMs, IsSuccess, ExitCode,
		OutputSummary, Notes, Comments
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := db.conn.Exec(query,
		log.StartupItemId, log.RunAt, log.DurationMs, successInt,
		log.ExitCode, log.OutputSummary, log.Notes, log.Comments,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record startup log")
	}

	return nil
}

// ListStartupLogs returns recent execution logs for an item.
func (db *StartupSplitDB) ListStartupLogs(itemId int64, limit int) ([]StartupLogRecord, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `SELECT StartupLogId, StartupItemId, RunAt, DurationMs,
		IsSuccess, ExitCode, OutputSummary, Notes, Comments
	FROM StartupLog WHERE StartupItemId = ? ORDER BY RunAt DESC LIMIT ?;`

	rows, err := db.conn.Query(query, itemId, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "list startup logs")
	}
	defer rows.Close()

	var logs []StartupLogRecord
	for rows.Next() {
		var l StartupLogRecord
		var successInt int
		var out, notes, comments sql.NullString

		if err := rows.Scan(
			&l.StartupLogId, &l.StartupItemId, &l.RunAt, &l.DurationMs,
			&successInt, &l.ExitCode, &out, &notes, &comments,
		); err != nil {
			return nil, apperror.WrapSimple(err, "scan startup log")
		}

		l.OutputSummary = out.String
		l.Notes = notes.String
		l.Comments = comments.String
		l.IsSuccess = successInt == 1
		logs = append(logs, l)
	}

	return logs, nil
}
