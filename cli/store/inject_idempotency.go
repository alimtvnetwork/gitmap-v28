// Package store — inject_idempotency.go: per-repo tracking of the last
// GitHub Desktop / VS Code injection. Backs `gitmap inject` and
// `gitmap open` to skip work that's already been done unless --force.
package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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

	row := QueryRowWrapper(db.conn, constants.SQLSelectInjectTimestamps, absPath)
	err := row.Scan(&ts.Desktop, &ts.VSCode)
	if errors.Is(err, sql.ErrNoRows) {
		return InjectTimestamps{}, nil
	}

	if err != nil {
		return InjectTimestamps{}, apperror.WrapSimple(err, "failed to query repos")
	}

	return ts, nil
}

// MarkInjected stamps the column for the given tool to CURRENT_TIMESTAMP.
// No-op (returns nil) when no row matches absPath — the upstream upsert
// is the only thing that creates Repo rows; injection alone never does.
func (db *DB) MarkInjected(absPath string, kind constants.InjectKindType) error {
	var stmt string
	switch kind {
	case constants.InjectKindDesktop:
		stmt = constants.SQLUpdateInjectDesktopTimestamp
	case constants.InjectKindVSCode:
		stmt = constants.SQLUpdateInjectVSCodeTimestamp
	default:
		return apperror.NewValidationError(fmt.Sprintf("inject_idempotency: unknown InjectKind %d", kind))
	}

	if _, err := ExecWrapper(db.conn, stmt, absPath).Destruct(); err != nil {
		return apperror.WrapSimple(err, "failed to upsert repo")
	}

	return nil
}
