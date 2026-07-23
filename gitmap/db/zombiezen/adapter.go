// Package zombiezen adapts the project's database/sql usage onto
// zombiezen.com/go/sqlite — a pure-Go SQLite driver that removes the
// cgo dependency, accelerates Windows builds, and simplifies cross-
// compilation.
//
// Migration strategy (incremental, behind a build tag):
//
//  1. New code paths import this package's Open() instead of the
//     mattn/go-sqlite3 driver.
//  2. Existing call sites are migrated one package at a time; tests
//     run under both drivers via the `cgosqlite` build tag.
//  3. Once parity tests pass for every migration in 001..008, the
//     mattn import is removed and the build tag retired.
//
// Until step 3, this package is opt-in.
package zombiezen

import (
	"context"
	"errors"
	"fmt"
)

// Conn is the abstract handle returned by Open. Mirrors the subset
// of database/sql.DB the rest of gitmap needs (Exec, Query, Close).
type Conn interface {
	Exec(ctx context.Context, query string, args ...any) error
	Close() error
}

// Open returns a pure-Go SQLite connection.
//
// Implementation note: the real driver wire-up lives in
// adapter_zombiezen.go behind the `purego_sqlite` build tag so this
// scaffold compiles without pulling the dependency until the
// migration is unblocked. Until then, Open returns ErrNotEnabled
// and existing mattn-backed code paths remain authoritative.
func Open(_ context.Context, _ string) (Conn, error) {
	return nil, fmt.Errorf("zombiezen sqlite adapter: %w", ErrNotEnabled)
}

// ErrNotEnabled signals that the project was built without the
// `purego_sqlite` build tag, so the pure-Go driver is unavailable.
var ErrNotEnabled = errors.New("pure-Go sqlite driver not compiled in (build with -tags purego_sqlite)")
