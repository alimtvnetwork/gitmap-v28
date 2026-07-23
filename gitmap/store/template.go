// Package store — template.go manages CommitTemplates CRUD operations.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// CommitTemplate represents a single commit message template.
type CommitTemplate struct {
	ID        int64
	Kind      string
	Template  string
	CreatedAt string
}

// InsertTemplate adds a single template to the database.
func (db *DB) InsertTemplate(kind, template string) error {
	_, err := db.conn.Exec(constants.SQLInsertTemplate, kind, template)
	if err != nil {
		return fmt.Errorf(constants.ErrSEODBInsert, err)
	}

	return nil
}

// ListTemplatesByKind returns all templates of the given kind.
func (db *DB) ListTemplatesByKind(kind string) ([]CommitTemplate, error) {
	rows, err := db.conn.Query(constants.SQLSelectTemplatesByKind, kind)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}
	defer rows.Close()

	return scanTemplateRows(rows)
}

// CountTemplates returns the total number of templates in the database.
func (db *DB) CountTemplates() (int, error) {
	var count int
	err := db.conn.QueryRow(constants.SQLCountTemplates).Scan(&count)

	return count, err
}

// scanTemplateRows reads CommitTemplate values from query result rows.
func scanTemplateRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]CommitTemplate, error) {
	var results []CommitTemplate

	for rows.Next() {
		t, err := scanOneTemplate(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, t)
	}

	return results, nil
}

// scanOneTemplate reads a single CommitTemplate from a row.
func scanOneTemplate(row interface{ Scan(dest ...any) error }) (CommitTemplate, error) {
	var t CommitTemplate
	err := row.Scan(&t.ID, &t.Kind, &t.Template, &t.CreatedAt)

	return t, err
}
