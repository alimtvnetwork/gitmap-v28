// Package store — amendment.go manages Amendments CRUD operations.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// AmendmentRow represents a single row from the Amendments table.
type AmendmentRow struct {
	ID            int64
	Branch        string
	FromCommit    string
	ToCommit      string
	TotalCommits  int
	PreviousName  string
	PreviousEmail string
	NewName       string
	NewEmail      string
	Mode          string
	ForcePushed   int
	CreatedAt     string
}

// InsertAmendment saves an amendment record to the database.
func (db *DB) InsertAmendment(branch, fromCommit, toCommit string, total int, prevName, prevEmail, newName, newEmail, mode string, forcePushed bool) error {
	fp := boolToIntAmend(forcePushed)

	_, err := db.conn.Exec(constants.SQLInsertAmendment,
		branch, fromCommit, toCommit, total,
		prevName, prevEmail, newName, newEmail, mode, fp)

	if err != nil {
		return fmt.Errorf(constants.ErrDBUpsert, err)
	}

	return nil
}

// ListAmendments returns all amendment records, most recent first.
func (db *DB) ListAmendments() ([]AmendmentRow, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllAmendments)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}
	defer rows.Close()

	return scanAmendmentRows(rows)
}

// ListAmendmentsByBranch returns amendments for a specific branch.
func (db *DB) ListAmendmentsByBranch(branch string) ([]AmendmentRow, error) {
	rows, err := db.conn.Query(constants.SQLSelectAmendmentsByBranch, branch)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}
	defer rows.Close()

	return scanAmendmentRows(rows)
}

// scanAmendmentRows reads AmendmentRow values from query result rows.
func scanAmendmentRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]AmendmentRow, error) {
	var results []AmendmentRow

	for rows.Next() {
		var r AmendmentRow
		err := rows.Scan(&r.ID, &r.Branch, &r.FromCommit, &r.ToCommit,
			&r.TotalCommits, &r.PreviousName, &r.PreviousEmail,
			&r.NewName, &r.NewEmail, &r.Mode, &r.ForcePushed, &r.CreatedAt)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}

// boolToIntAmend converts a bool to SQLite integer (0/1).
func boolToIntAmend(b bool) int {
	if b {
		return 1
	}

	return 0
}
