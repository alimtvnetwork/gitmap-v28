package store

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ArchiveHistoryRow is the persisted shape of a `gitmap zip` /
// `gitmap unzip-compact` invocation. One row is inserted at command
// start (Status = "" until finished) so a crash mid-extract still
// leaves a forensic trace.
type ArchiveHistoryRow struct {
	ID                     int64
	CommandName            string
	InputSources           []string
	OutputPath             string
	ArchiveFormat          string
	CompressionMode        string
	UsedTemporaryDirectory bool
	Status                 string
	ErrorMessage           string
	StartedAt              string
	CompletedAt            string
}

// StartArchiveHistory inserts a partial row at command start and returns
// the new ArchiveHistoryId so the caller can finalize it later. Failures
// are surfaced (not swallowed) — callers may choose to continue without a
// history row, but the choice is theirs.
func (db *DB) StartArchiveHistory(cmd string, inputs []string, mode string) (int64, error) {
	raw, err := json.Marshal(inputs)
	if err != nil {
		raw = []byte("[]")
	}

	res, err := db.conn.Exec(constants.SQLInsertArchiveHistory,
		cmd,
		string(raw),
		"",   // OutputPath filled at finish
		"",   // ArchiveFormat filled at finish
		mode, // CompressionMode known up-front
		0,    // UsedTemporaryDirectory updated at finish
		"",   // Status — empty means "in flight"
		"",   // ErrorMessage
		time.Now().UTC().Format(time.RFC3339),
		"", // CompletedAt
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// FinishArchiveHistory updates an in-flight row with the final outcome.
// Pass status = constants.ArchiveStatusSuccess / ArchiveStatusFailed.
func (db *DB) FinishArchiveHistory(
	id int64,
	outputPath, format, status, errMsg string,
	usedTempDir bool,
) error {
	tempInt := 0
	if usedTempDir {
		tempInt = 1
	}

	_, err := db.conn.Exec(constants.SQLUpdateArchiveHistoryFinish,
		outputPath, format, tempInt, status, errMsg,
		time.Now().UTC().Format(time.RFC3339),
		id,
	)

	return err
}

// RecentArchiveHistory returns the N most recent rows for `gitmap` UI
// surfaces (e.g. a future `gitmap history archive` view). Capped to 200
// to keep memory bounded even when callers pass a huge limit.
func (db *DB) RecentArchiveHistory(limit int) ([]ArchiveHistoryRow, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}

	rows, err := db.conn.Query(constants.SQLSelectArchiveHistoryRecent, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ArchiveHistoryRow
	for rows.Next() {
		row, err := scanArchiveHistory(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}

	return out, rows.Err()
}

// scanArchiveHistory unmarshals the JSON-encoded InputSources column back
// into a string slice and decodes the boolean tempdir flag.
func scanArchiveHistory(rows *sql.Rows) (ArchiveHistoryRow, error) {
	var (
		row     ArchiveHistoryRow
		rawIn   string
		tempInt int
	)

	if err := rows.Scan(
		&row.ID, &row.CommandName, &rawIn, &row.OutputPath,
		&row.ArchiveFormat, &row.CompressionMode, &tempInt,
		&row.Status, &row.ErrorMessage, &row.StartedAt, &row.CompletedAt,
	); err != nil {
		return row, err
	}

	if rawIn != "" {
		_ = json.Unmarshal([]byte(rawIn), &row.InputSources)
	}
	row.UsedTemporaryDirectory = tempInt != 0

	return row, nil
}
