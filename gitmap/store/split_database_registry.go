package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// SplitDatabaseEntry represents a registered split SQLite database tracked in the root DB.
type SplitDatabaseEntry struct {
	ID                      int64  `json:"id"`
	SplitDatabaseRegistryID int64  `json:"splitDatabaseRegistryId"`
	DatabaseType            string `json:"databaseType"`
	DatabaseKey             string `json:"databaseKey"`
	DatabasePath            string `json:"databasePath"`
	SizeBytes               int64  `json:"sizeBytes"`
	TableCount              int    `json:"tableCount"`
	RecordCount             int64  `json:"recordCount"`
	SchemaVersion           int    `json:"schemaVersion"`
	Status                  string `json:"status"`
	IsActive                bool   `json:"isActive"`
	IsAttached              bool   `json:"isAttached"`
	Description             string `json:"description,omitempty"`
	Notes                   string `json:"notes,omitempty"`
	Comments                string `json:"comments,omitempty"`
	LastAccessedAt          int64  `json:"lastAccessedAt"`
	LastSyncedAt            int64  `json:"lastSyncedAt"`
	CreatedAt               int64  `json:"createdAt"`
	UpdatedAt               int64  `json:"updatedAt"`
}

// RegisterSplitDB inserts or updates a split database registration in the root DB.
func (db *DB) RegisterSplitDB(entry SplitDatabaseEntry) error {
	entry = normalizeSplitDBEntry(entry)
	args := entryToUpsertArgs(entry)
	_, err := ExecWrapper(db.conn, constants.SQLUpsertSplitDatabaseRegistry, args...).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "store.RegisterSplitDB")
	}

	return nil
}

// GetSplitDB retrieves a split database entry by type and key.
func (db *DB) GetSplitDB(dbType, dbKey string) (*SplitDatabaseEntry, error) {
	row := db.conn.QueryRow(constants.SQLSelectSplitDB, dbType, dbKey)
	entry, err := scanSplitEntry(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil, apperror.WrapSimple(err, "store.GetSplitDB")
	}

	return entry, nil
}

// ListSplitDBs lists split database entries, optionally filtered by database type.
func (db *DB) ListSplitDBs(dbType string) ([]SplitDatabaseEntry, error) {
	rows, err := db.querySplitDBs(dbType)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListSplitDBs")
	}

	defer rows.Close()

	return scanSplitDBRows(rows)
}

func (db *DB) querySplitDBs(dbType string) (*sql.Rows, error) {
	if dbType != "" {
		return db.conn.Query(constants.SQLSelectListSplitDBByType, dbType)
	}

	return db.conn.Query(constants.SQLSelectListAllSplitDB)
}

func scanSplitDBRows(rows *sql.Rows) ([]SplitDatabaseEntry, error) {
	var list []SplitDatabaseEntry
	for rows.Next() {
		entry, err := scanSplitEntry(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "store.scanSplitDBRows")
		}

		list = append(list, *entry)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "store.scanSplitDBRows.iter")
	}

	return list, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSplitEntry(s rowScanner) (*SplitDatabaseEntry, error) {
	var e SplitDatabaseEntry
	var active, attached int
	var desc, notes, comm sql.NullString
	dest := entryScanDest(&e, &active, &attached, &desc, &notes, &comm)
	if err := s.Scan(dest...); err != nil {
		return nil, err
	}

	populateEntryAux(&e, active, attached, desc, notes, comm)

	return &e, nil
}

func entryScanDest(e *SplitDatabaseEntry, act, att *int, d, n, c *sql.NullString) []any {
	return []any{
		&e.ID, &e.DatabaseType, &e.DatabaseKey, &e.DatabasePath,
		&e.SizeBytes, &e.TableCount, &e.RecordCount, &e.SchemaVersion,
		&e.Status, act, att, d, n, c,
		&e.LastAccessedAt, &e.LastSyncedAt, &e.CreatedAt, &e.UpdatedAt,
	}
}

func populateEntryAux(e *SplitDatabaseEntry, act, att int, d, n, c sql.NullString) {
	e.SplitDatabaseRegistryID = e.ID
	e.IsActive = act == 1
	e.IsAttached = att == 1
	e.Description = d.String
	e.Notes = n.String
	e.Comments = c.String
}

func normalizeSplitDBEntry(e SplitDatabaseEntry) SplitDatabaseEntry {
	if e.Status == "" {
		e.Status = "active"
	}

	if e.SchemaVersion == 0 {
		e.SchemaVersion = 1
	}

	return fillEntryTimestamps(e)
}

func fillEntryTimestamps(e SplitDatabaseEntry) SplitDatabaseEntry {
	now := time.Now().Unix()
	if e.LastSyncedAt == 0 {
		e.LastSyncedAt = now
	}

	if e.CreatedAt == 0 {
		e.CreatedAt = now
	}

	if e.UpdatedAt == 0 {
		e.UpdatedAt = now
	}

	return e
}

func entryToUpsertArgs(e SplitDatabaseEntry) []any {
	return []any{
		e.DatabaseType, e.DatabaseKey, e.DatabasePath,
		e.SizeBytes, e.TableCount, e.RecordCount,
		e.SchemaVersion, e.Status, boolToInt(e.IsActive), boolToInt(e.IsAttached),
		e.Description, e.Notes, e.Comments,
		e.LastAccessedAt, e.LastSyncedAt, e.CreatedAt, e.UpdatedAt,
	}
}
