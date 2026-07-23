package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// InsertBookmark saves a new bookmark record.
func (db *DB) InsertBookmark(r model.BookmarkRecord) error {
	_, err := db.conn.Exec(constants.SQLInsertBookmark,
		r.Name, r.Command, r.Args, r.Flags)
	if err != nil {
		return fmt.Errorf(constants.ErrBookmarkQuery, err)
	}

	return nil
}

// ListBookmarks returns all saved bookmarks ordered by name.
func (db *DB) ListBookmarks() ([]model.BookmarkRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllBookmarks)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrBookmarkQuery, err)
	}
	defer rows.Close()

	return scanBookmarkRows(rows)
}

// FindBookmarkByName returns a single bookmark by name.
func (db *DB) FindBookmarkByName(name string) (model.BookmarkRecord, error) {
	row := db.conn.QueryRow(constants.SQLSelectBookmarkByName, name)

	var r model.BookmarkRecord
	err := row.Scan(&r.ID, &r.Name, &r.Command, &r.Args, &r.Flags, &r.CreatedAt)
	if err != nil {
		return r, fmt.Errorf(constants.ErrBookmarkQuery, err)
	}

	return r, nil
}

// DeleteBookmark removes a bookmark by name.
func (db *DB) DeleteBookmark(name string) error {
	_, err := db.conn.Exec(constants.SQLDeleteBookmark, name)
	if err != nil {
		return fmt.Errorf(constants.ErrBookmarkQuery, err)
	}

	return nil
}

// scanBookmarkRows reads all rows into BookmarkRecord slices.
func scanBookmarkRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.BookmarkRecord, error) {
	var results []model.BookmarkRecord

	for rows.Next() {
		var r model.BookmarkRecord
		err := rows.Scan(&r.ID, &r.Name, &r.Command, &r.Args, &r.Flags, &r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrBookmarkQuery, err)
		}
		results = append(results, r)
	}

	return results, nil
}
