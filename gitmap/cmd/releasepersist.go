package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// persistReleaseToDB saves the release metadata to SQLite if available.
// v17: requires the current repo to be registered in the Repo table so the
// new Release.RepoId FK can be satisfied.
func persistReleaseToDB() {
	meta := release.LastMeta
	if meta == nil {
		return
	}

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not cache release to database: %v\n", err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Release DB migration failed: %v\n", err)
	}

	repoID, err := resolveOrRegisterCurrentRepoID(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve repo for release: %v\n", err)

		return
	}

	record := releaseMetaToRecord(*meta)
	record.RepoID = repoID
	if err := db.UpsertRelease(record); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not cache release metadata: %v\n", err)
	}
}

// resolveCurrentRepoID returns the RepoId for the cwd. Caller must handle
// the "repo not scanned" error path explicitly.
func resolveCurrentRepoID(db *store.DB) (int64, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	return db.ResolveCurrentRepoID(cwd)
}

// resolveOrRegisterCurrentRepoID first tries to resolve the cwd's RepoId;
// when the repo is not yet registered, it auto-registers (parent dir =
// ScanFolder, cwd = Repo) and re-resolves. This makes `gitmap r` self-
// healing on a fresh clone where `gitmap scan` was never run.
func resolveOrRegisterCurrentRepoID(db *store.DB) (int64, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	if id, err := db.ResolveCurrentRepoID(cwd); err == nil {
		return id, nil
	}

	if err := autoRegisterCurrentRepo(db, cwd); err != nil {
		return 0, fmt.Errorf("auto-register failed: %w", err)
	}

	return db.ResolveCurrentRepoID(cwd)
}

// releaseMetaToRecord converts a ReleaseMeta to a ReleaseRecord for DB storage.
func releaseMetaToRecord(m release.ReleaseMeta) model.ReleaseRecord {
	return model.ReleaseRecord{
		Version:      m.Version,
		Tag:          m.Tag,
		Branch:       m.Branch,
		SourceBranch: m.SourceBranch,
		CommitSha:    m.Commit,
		Changelog:    store.JoinChangelog(m.Changelog),
		Notes:        m.Notes,
		IsDraft:      m.IsDraft,
		IsPreRelease: m.IsPreRelease,
		IsLatest:     m.IsLatest,
		Source:       model.SourceRelease,
		CreatedAt:    m.CreatedAt,
	}
}
