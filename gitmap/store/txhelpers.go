// Package store — txhelpers.go: tiny interface + helper used by store
// methods that wrap multi-statement writes in a sql.Tx. Lets the per-
// row loops accept either *sql.Tx or a mock in unit tests without
// dragging the database/sql dependency into the test surface.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §plan step 18.
package store

import (
	"database/sql"
	"fmt"
)

// txExecer is the minimal subset of *sql.Tx the per-row insert helpers
// need. Defined locally so tests can pass a fake without importing the
// real driver.
type txExecer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

// commitOrWrap runs tx.Commit() and wraps any error in the supplied
// Code Red format (operation + reason). Returns nil on success so
// callers can `return commitOrWrap(...)` directly.
func commitOrWrap(tx *sql.Tx, errFmt string) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf(errFmt, err, err.Error())
	}

	return nil
}
