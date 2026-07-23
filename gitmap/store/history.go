package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// InsertHistory inserts a new command history record and returns the auto-generated ID.
func (db *DB) InsertHistory(r model.CommandHistoryRecord) (int64, error) {
	result, err := db.conn.Exec(constants.SQLInsertHistory,
		r.Command, r.Alias, r.Args, r.Flags,
		r.StartedAt, r.FinishedAt, r.DurationMs, r.ExitCode, r.Summary, r.RepoCount)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrHistoryQuery, err)
	}

	return result.LastInsertId()
}

// UpdateHistory updates a history record with completion details.
func (db *DB) UpdateHistory(r model.CommandHistoryRecord) error {
	_, err := db.conn.Exec(constants.SQLUpdateHistory,
		r.FinishedAt, r.DurationMs, r.ExitCode, r.Summary, r.RepoCount, r.ID)
	if err != nil {
		return fmt.Errorf(constants.ErrHistoryQuery, err)
	}

	return nil
}

// ListHistory returns all command history records, newest first.
func (db *DB) ListHistory() ([]model.CommandHistoryRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllHistory)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrHistoryQuery, err)
	}
	defer rows.Close()

	return scanHistoryRows(rows)
}

// ListHistoryByCommand returns history filtered by command name.
func (db *DB) ListHistoryByCommand(command string) ([]model.CommandHistoryRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectHistoryByCommand, command)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrHistoryQuery, err)
	}
	defer rows.Close()

	return scanHistoryRows(rows)
}

// ClearHistory deletes all command history records.
func (db *DB) ClearHistory() error {
	_, err := db.conn.Exec(constants.SQLDeleteAllHistory)

	return err
}

// scanHistoryRows reads all rows into CommandHistoryRecord slices.
func scanHistoryRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.CommandHistoryRecord, error) {
	var results []model.CommandHistoryRecord

	for rows.Next() {
		r, err := scanOneHistory(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// scanOneHistory scans a single row into a CommandHistoryRecord.
func scanOneHistory(row interface{ Scan(dest ...any) error }) (model.CommandHistoryRecord, error) {
	var r model.CommandHistoryRecord
	err := row.Scan(&r.ID, &r.Command, &r.Alias, &r.Args, &r.Flags,
		&r.StartedAt, &r.FinishedAt, &r.DurationMs, &r.ExitCode,
		&r.Summary, &r.RepoCount, &r.CreatedAt)

	return r, err
}
