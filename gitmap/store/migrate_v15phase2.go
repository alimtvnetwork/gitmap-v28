// Package store — migrate_v15phase2.go performs the Phase 1.2 v15 renames:
//
//	Groups     → Group     (Id → GroupId)
//	Releases   → Release   (Id → ReleaseId)
//	Aliases    → Alias     (Id → AliasId)
//	Bookmarks  → Bookmark  (Id → BookmarkId)
//
// Each migration is idempotent (detect-then-act on the legacy plural name)
// and uses the shared runV15Rebuild helper.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// migrateV15Phase2 runs all Phase 1.2 table rebuilds in dependency-safe order.
// Order: Group first (referenced by GroupRepo), then Release/Alias/Bookmark
// (independent). FKs are off during each rebuild via runV15Rebuild.
func (db *DB) migrateV15Phase2() error {
	specs := []v15RebuildSpec{
		{
			OldTable:      "Groups",
			NewTable:      "Group",
			NewCreateSQL:  constants.SQLCreateGroup,
			OldColumnList: "Id, Name, Description, Color, CreatedAt",
			NewColumnList: "GroupId, Name, Description, Color, CreatedAt",
			StartMsg:      "→ Migrating Groups → Group (GroupId PK)...",
			DoneMsg:       "✓ Migrated Groups → Group.",
		},
		{
			OldTable:      "Releases",
			NewTable:      "Release",
			NewCreateSQL:  constants.SQLCreateRelease,
			OldColumnList: "Id, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt",
			NewColumnList: "ReleaseId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt",
			StartMsg:      "→ Migrating Releases → Release (ReleaseId PK + IsDraft/IsPreRelease)...",
			DoneMsg:       "✓ Migrated Releases → Release.",
		},
		{
			OldTable:      "Aliases",
			NewTable:      "Alias",
			NewCreateSQL:  constants.SQLCreateAlias,
			OldColumnList: "Id, Alias, RepoId, CreatedAt",
			NewColumnList: "AliasId, Alias, RepoId, CreatedAt",
			StartMsg:      "→ Migrating Aliases → Alias (AliasId PK)...",
			DoneMsg:       "✓ Migrated Aliases → Alias.",
		},
		{
			OldTable:      "Bookmarks",
			NewTable:      "Bookmark",
			NewCreateSQL:  constants.SQLCreateBookmark,
			OldColumnList: "Id, Name, Command, Args, Flags, CreatedAt",
			NewColumnList: "BookmarkId, Name, Command, Args, Flags, CreatedAt",
			StartMsg:      "→ Migrating Bookmarks → Bookmark (BookmarkId PK)...",
			DoneMsg:       "✓ Migrated Bookmarks → Bookmark.",
		},
	}

	for _, spec := range specs {
		if err := db.runV15Rebuild(spec); err != nil {
			return fmt.Errorf("phase 1.2 %s: %w", spec.OldTable, err)
		}
	}

	// After all four singulars exist, rebuild GroupRepo so its FKs point at
	// the new "Group"(GroupId) and Repo(RepoId). The plural "Groups" is gone
	// at this point, so the legacy GroupRepo (which references Groups(Id))
	// would break on next FK validation. Drop and recreate.
	if db.tableExists("GroupRepo") {
		// Preserve any existing rows by copying through a temp table.
		if err := db.rebuildGroupRepoFK(); err != nil {
			return fmt.Errorf("rebuild GroupRepo FK: %w", err)
		}
	}

	return nil
}

// rebuildGroupRepoFK rewrites the GroupRepo CREATE so its FK references the
// new singular tables. SQLite stores the FK clause as text in sqlite_master,
// so renaming the parent table does NOT auto-update the child FK definition.
func (db *DB) rebuildGroupRepoFK() error {
	if _, err := db.conn.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return err
	}

	defer func() {
		_, _ = db.conn.Exec("PRAGMA foreign_keys = ON")
	}()

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS GroupRepo_v15new (
			GroupId INTEGER NOT NULL REFERENCES "Group"(GroupId) ON DELETE CASCADE,
			RepoId  INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
			PRIMARY KEY (GroupId, RepoId)
		)`,
		`INSERT OR IGNORE INTO GroupRepo_v15new (GroupId, RepoId) SELECT GroupId, RepoId FROM GroupRepo`,
		`DROP TABLE GroupRepo`,
		`ALTER TABLE GroupRepo_v15new RENAME TO GroupRepo`,
	}

	for _, s := range stmts {
		if _, err := db.conn.Exec(s); err != nil {
			return fmt.Errorf("groupRepo rebuild step (%s): %w", firstWords(s, 4), err)
		}
	}

	return nil
}

// firstWords returns the first n whitespace-separated words of s for error context.
func firstWords(s string, n int) string {
	count := 0

	for i, r := range s {
		if r == ' ' || r == '\n' || r == '\t' {
			count++
			if count >= n {
				return s[:i]
			}
		}
	}

	return s
}
