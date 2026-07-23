// Package store — migrate_v15phase6.go performs the v17 release-repo FK migration.
//
// Spec: spec/04-generic-cli/24-release-repo-relationship.md
//
// Until v3.16.x the Release table was orphaned (no FK to Repo). v3.17.0
// adds:
//   - RepoId INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE
//   - Composite UNIQUE(RepoId, Tag) replacing the prior global Tag UNIQUE
//   - IdxRelease_RepoId index for per-repo filtering
//
// Migration policy (user-approved): WIPE existing Release rows and let the
// next `gitmap list-releases` re-import from .gitmap/release/v*.json.
// Backfilling RepoId from path heuristics would be fragile and the on-disk
// release metadata is the canonical source-of-truth.
package store

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// migrateV15Phase6 wipes the Release table when it lacks the RepoId column,
// allowing the standard CREATE pass to recreate it with the new FK schema.
// Fresh installs (no Release table yet) are no-ops.
func (db *DB) migrateV15Phase6() error {
	if !db.tableExists(constants.TableRelease) {
		return nil
	}

	if db.columnExists(constants.TableRelease, "RepoId") {
		return nil // already migrated.
	}

	fmt.Fprintln(os.Stderr, constants.MsgV15Phase6Start)
	fmt.Fprintln(os.Stderr, constants.MsgV15Phase6Wipe)

	if _, err := db.conn.Exec(constants.SQLDropRelease); err != nil {
		return fmt.Errorf("phase 1.6 Release drop: %w", err)
	}

	fmt.Fprintln(os.Stderr, constants.MsgV15Phase6Done)

	return nil
}
