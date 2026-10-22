package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// InsertBookmark saves a new bookmark record.
func (db *DB) InsertBookmark(r model.BookmarkRecord) error {
	_, err := ExecWrapper(db.conn, constants.SQLInsertBookmark,
		r.Name, r.Command, r.Args, r.Flags).Destruct()
	if err != nil {
		return fmt.Errorf(constants.ErrBookmarkQuery, err)
	}

	return nil
}

// ListBookmarks returns all saved bookmarks ordered by name.
func (db *DB) ListBookmarks() ([]model.BookmarkRecord, error) {
	res := db.ListBookmarksWrap()
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

func (db *DB) ListBookmarksWrap() BookmarkRecordSliceWrap {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllBookmarks).Destruct()
	if err != nil {
		return result.FailSlice[model.BookmarkRecord](apperror.WrapSimple(err, "ListBookmarks"))
	}

	defer rows.Close()

	items, errScan := scanBookmarkRows(rows)
	if errScan != nil {
		return result.FailSlice[model.BookmarkRecord](apperror.WrapSimple(errScan, "ListBookmarks"))
	}

	return result.OkSlice(items)
}

// FindBookmarkByName returns a single bookmark by name.
func (db *DB) FindBookmarkByName(name string) (model.BookmarkRecord, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectBookmarkByName, name)

	var r model.BookmarkRecord
	err := row.Scan(&r.ID, &r.Name, &r.Command, &r.Args, &r.Flags, &r.CreatedAt)
	if err != nil {
		return r, fmt.Errorf(constants.ErrBookmarkQuery, err)
	}

	return r, nil
}

// DeleteBookmark removes a bookmark by name.
func (db *DB) DeleteBookmark(name string) error {
	_, err := ExecWrapper(db.conn, constants.SQLDeleteBookmark, name).Destruct()
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
