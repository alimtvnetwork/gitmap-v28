// Package store — inject_idempotency.go: per-repo tracking of the last
// GitHub Desktop / VS Code injection. Backs `gitmap inject` and
// `gitmap open` to skip work that's already been done unless --force.
package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// InjectTimestamps holds the per-tool stamps for a single repo. Empty
// strings mean "never injected"; the zero-value is therefore the
// natural "do everything" answer.
type InjectTimestamps struct {
	Desktop string
	VSCode  string
}

// GetInjectTimestamps reads both stamps for the row matching absPath.
// Returns the zero value (and nil) when the path is unknown so callers
// can treat unknown == never-injected without special-casing.
func (db *DB) GetInjectTimestamps(absPath string) (InjectTimestamps, error) {
	var ts InjectTimestamps

	row := db.conn.QueryRow(constants.SQLSelectInjectTimestamps, absPath)
	if err := row.Scan(&ts.Desktop, &ts.VSCode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return InjectTimestamps{}, nil
		}

		return InjectTimestamps{}, fmt.Errorf(constants.ErrDBQuery, err)
	}

	return ts, nil
}

// MarkInjected stamps the column for the given tool to CURRENT_TIMESTAMP.
// No-op (returns nil) when no row matches absPath — the upstream upsert
// is the only thing that creates Repo rows; injection alone never does.
func (db *DB) MarkInjected(absPath string, kind constants.InjectKind) error {
	col := injectColumnFor(kind)
	if col == "" {
		return fmt.Errorf("inject_idempotency: unknown InjectKind %d", kind)
	}

	stmt := fmt.Sprintf(constants.SQLUpdateInjectTimestampFmt, col)
	if _, err := db.conn.Exec(stmt, absPath); err != nil {
		return fmt.Errorf(constants.ErrDBUpsert, err)
	}

	return nil
}

// injectColumnFor maps the typed kind to its physical column name.
func injectColumnFor(kind constants.InjectKind) string {
	switch kind {
	case constants.InjectKindDesktop:
		return constants.ColInjectDesktop
	case constants.InjectKindVSCode:
		return constants.ColInjectVSCode
	default:
		return ""
	}
}
