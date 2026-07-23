package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// ImportAll restores a DatabaseExport into the database using upsert/insert-or-ignore semantics.
func (db *DB) ImportAll(data model.DatabaseExport) error {
	if err := db.importRepos(data.Repos); err != nil {
		return fmt.Errorf("repos: %w", err)
	}

	if err := db.importGroups(data.Groups); err != nil {
		return fmt.Errorf("groups: %w", err)
	}

	if err := db.importReleases(data.Releases); err != nil {
		return fmt.Errorf("releases: %w", err)
	}

	if err := db.importHistory(data.History); err != nil {
		return fmt.Errorf("history: %w", err)
	}

	if err := db.importBookmarks(data.Bookmarks); err != nil {
		return fmt.Errorf("bookmarks: %w", err)
	}

	return nil
}

// importRepos upserts all repos by ID.
func (db *DB) importRepos(repos []model.ScanRecord) error {
	for _, r := range repos {
		_, err := db.conn.Exec(constants.SQLUpsertRepo,
			r.Slug, r.RepoName, r.HTTPSUrl, r.SSHUrl,
			r.Branch, r.RelativePath, r.AbsolutePath,
			r.CloneInstruction, r.Notes, r.Transport)
		if err != nil {
			return err
		}
	}

	return nil
}

// importGroups creates groups and links repos by slug.
func (db *DB) importGroups(groups []model.GroupExport) error {
	for _, ge := range groups {
		if err := db.importOneGroup(ge); err != nil {
			return err
		}
	}

	return nil
}

// importOneGroup creates a group and links its member repos.
func (db *DB) importOneGroup(ge model.GroupExport) error {
	_, err := db.conn.Exec(constants.SQLImportInsertGroup,
		ge.Name, ge.Description, ge.Color)
	if err != nil {
		return err
	}
	group, err := db.findGroupByName(ge.Name)
	if err != nil {
		return err
	}

	return db.linkGroupRepos(group.ID, ge.RepoSlugs)
}

// linkGroupRepos links repos to a group by resolving slugs.
func (db *DB) linkGroupRepos(groupID int64, slugs []string) error {
	for _, slug := range slugs {
		repos, err := db.FindBySlug(slug)
		if err != nil || len(repos) == 0 {
			continue
		}

		_, err = db.conn.Exec(constants.SQLInsertGroupRepo, groupID, repos[0].ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// importReleases upserts all release records.
func (db *DB) importReleases(releases []model.ReleaseRecord) error {
	for _, r := range releases {
		if err := db.UpsertRelease(r); err != nil {
			return err
		}
	}

	return nil
}

// importHistory inserts history records, ignoring duplicates.
func (db *DB) importHistory(records []model.CommandHistoryRecord) error {
	for _, r := range records {
		_, err := db.conn.Exec(
			"INSERT OR IGNORE INTO CommandHistory (Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			r.Command, r.Alias, r.Args, r.Flags,
			r.StartedAt, r.FinishedAt, r.DurationMs, r.ExitCode, r.Summary, r.RepoCount)
		if err != nil {
			return err
		}
	}

	return nil
}

// importBookmarks inserts bookmarks, ignoring duplicates by name.
func (db *DB) importBookmarks(bookmarks []model.BookmarkRecord) error {
	for _, b := range bookmarks {
		_, err := db.conn.Exec(constants.SQLImportInsertBookmark,
			b.Name, b.Command, b.Args, b.Flags)
		if err != nil {
			return err
		}
	}

	return nil
}
