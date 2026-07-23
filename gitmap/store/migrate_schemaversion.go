// Package store — migrate_schemaversion.go: the one-time schema-version
// marker fast-path.
//
// Migrate() is called by openDB() on EVERY gitmap subcommand. Before this
// fast-path, every command paid the full v15 phase pipeline cost — dozens
// of PRAGMA + tableExists() round-trips even when the schema was already
// current. After this fast-path, a single SELECT against Setting decides
// whether to skip everything.
//
// The marker is intentionally stored in the existing Setting key/value
// table (key = SettingSchemaVersion, value = stringified int) so:
//
//   - No new table is needed (Setting is already in the standard CREATE
//     pass and survives db-reset → re-create).
//   - GetSetting()/SetSetting() can be reused without new infrastructure.
//   - Legacy databases with no Setting table return "" from GetSetting(),
//     which parses to schema version 0 — guaranteed to be < current,
//     so the full pipeline runs exactly once on first launch after the
//     fast-path is introduced.
//
// The marker is cleared by db-reset (the whole DB is wiped) and by
// migrateLegacyIDs() when it rebuilds the Repos table, so any database
// requiring genuine repair always re-runs the full pipeline.
package store

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// readSchemaVersion returns the persisted schema version, or 0 when the
// marker is missing / unreadable / the Setting table does not exist yet.
// 0 is the safe sentinel: it always compares less-than the current target,
// so the full migration pipeline runs.
func (db *DB) readSchemaVersion() int {
	if !db.tableExists(constants.TableSetting) {
		return 0
	}

	raw := db.GetSetting(constants.SettingSchemaVersion)
	if raw == "" {
		return 0
	}

	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}

	return n
}

// writeSchemaVersion records the current schema version. Failures are
// warned to stderr but never fatal — a missed write just means the next
// invocation pays the pipeline cost again, which is correct behavior.
func (db *DB) writeSchemaVersion(version int) {
	err := db.SetSetting(constants.SettingSchemaVersion, strconv.Itoa(version))
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnSchemaVersionWriteFmt, version, err)
	}
}

// isSchemaUpToDate reports whether the persisted marker matches the
// current target. A true return means Migrate() can safely no-op.
func (db *DB) isSchemaUpToDate() bool {
	return db.readSchemaVersion() == constants.SchemaVersionCurrent
}
