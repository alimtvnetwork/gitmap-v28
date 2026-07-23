// Package store — migrate_v15phase5.go performs the Phase 1.5 v15 column rename:
//
//	Release.Draft       → Release.IsDraft
//	Release.PreRelease  → Release.IsPreRelease
//
// This handles the upgrade path from v3.4.x where the singular Release table
// already exists with the old column names. SQLite does not let us rename
// columns inside CREATE-with-FKs idempotently across all versions, so we
// detect-then-act: if the old Draft column exists on Release, rebuild the
// table copying Draft → IsDraft and PreRelease → IsPreRelease.
//
// Fresh installs and upgrades from older versions (where Phase 1.2 wrote
// IsDraft/IsPreRelease directly via the new SQLCreateRelease) are no-ops.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// migrateV15Phase5 renames Release.Draft → IsDraft and Release.PreRelease → IsPreRelease.
func (db *DB) migrateV15Phase5() error {
	if !db.tableExists("Release") {
		return nil // fresh install — Phase 1.2 + standard CREATE pass own this.
	}

	if !db.columnExists("Release", "Draft") {
		return nil // already migrated, or freshly created with IsDraft.
	}

	if db.columnExists("Release", "IsDraft") {
		// Both columns exist — refuse to clobber. Should not happen given
		// the rebuild drops the old table, but defensive.
		return nil
	}

	spec := v15RebuildSpec{
		OldTable:      "Release",
		NewTable:      "Release_v15phase5",
		NewCreateSQL:  `CREATE TABLE IF NOT EXISTS Release_v15phase5 (` + releasePhase5Body() + `)`,
		OldColumnList: "ReleaseId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt",
		NewColumnList: "ReleaseId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt",
		StartMsg:      "→ Migrating Release.Draft → IsDraft, Release.PreRelease → IsPreRelease...",
		DoneMsg:       "✓ Renamed Release.Draft → IsDraft and Release.PreRelease → IsPreRelease.",
	}

	if err := db.runV15Rebuild(spec); err != nil {
		return fmt.Errorf("phase 1.5 Release: %w", err)
	}

	// Final swap: Release_v15phase5 → Release. The standard CREATE pass
	// after Migrate() will then no-op against the canonical Release name.
	finalSpec := v15RebuildSpec{
		OldTable:      "Release_v15phase5",
		NewTable:      "Release",
		NewCreateSQL:  constants.SQLCreateRelease,
		OldColumnList: "ReleaseId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt",
		NewColumnList: "ReleaseId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt",
	}

	return db.runV15Rebuild(finalSpec)
}

// releasePhase5Body is the inline schema for the staging table. We can't
// reuse SQLCreateRelease verbatim because that targets the canonical name
// "Release" — the staging table needs a unique name during the rebuild.
func releasePhase5Body() string {
	return `
	ReleaseId    INTEGER PRIMARY KEY AUTOINCREMENT,
	Version      TEXT NOT NULL,
	Tag          TEXT NOT NULL UNIQUE,
	Branch       TEXT NOT NULL,
	SourceBranch TEXT NOT NULL,
	CommitSha    TEXT NOT NULL,
	Changelog    TEXT DEFAULT '',
	Notes        TEXT DEFAULT '',
	IsDraft      INTEGER DEFAULT 0,
	IsPreRelease INTEGER DEFAULT 0,
	IsLatest     INTEGER DEFAULT 0,
	Source       TEXT DEFAULT 'release',
	CreatedAt    TEXT DEFAULT CURRENT_TIMESTAMP`
}
