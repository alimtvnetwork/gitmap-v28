package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// importReleases discovers .gitmap/release/v*.json files and upserts them into the DB.
// v17: stamps each record with the current repo's RepoId.
func importReleases(scanDir, outputDir string) {
	releaseDir := filepath.Join(scanDir, constants.DefaultReleaseDir)
	files := discoverReleaseFiles(releaseDir)
	if len(files) == 0 {
		return
	}

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)

		return
	}

	repoID, err := resolveImportRepoID(db, scanDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Skipping release import: %v\n", err)

		return
	}

	count := upsertReleaseFiles(db, files, repoID)
	if count > 0 {
		fmt.Printf(constants.MsgReleasesImported, count)
	}
}

// resolveImportRepoID looks up the RepoId for the scan target directory.
func resolveImportRepoID(db *store.DB, scanDir string) (int64, error) {
	abs, err := filepath.Abs(scanDir)
	if err != nil {
		return 0, err
	}

	return db.ResolveCurrentRepoID(abs)
}

// discoverReleaseFiles returns paths to all v*.json files in the release dir.
func discoverReleaseFiles(releaseDir string) []string {
	pattern := filepath.Join(releaseDir, constants.ReleaseGlob)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	return matches
}

// upsertReleaseFiles reads and upserts each release file, returning the count.
func upsertReleaseFiles(db *store.DB, files []string, repoID int64) int {
	count := 0
	for _, f := range files {
		if isLatestFile(f) {
			continue
		}
		if importOneRelease(db, f, repoID) {
			count++
		}
	}

	return count
}

// isLatestFile checks if the file is latest.json (not a release file).
func isLatestFile(path string) bool {
	return filepath.Base(path) == constants.DefaultLatestFile
}

// importOneRelease reads a single release file and upserts it.
func importOneRelease(db *store.DB, path string, repoID int64) bool {
	meta, err := release.ReadReleaseMeta(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnReleaseImportSkip, filepath.Base(path), err)

		return false
	}

	record := mapMetaToRecord(meta)
	record.RepoID = repoID
	if err := db.UpsertRelease(record); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnReleaseImportSkip, filepath.Base(path), err)

		return false
	}

	return true
}

// mapMetaToRecord converts a ReleaseMeta to a ReleaseRecord.
func mapMetaToRecord(m release.ReleaseMeta) model.ReleaseRecord {
	return model.ReleaseRecord{
		Version:      m.Version,
		Tag:          m.Tag,
		Branch:       m.Branch,
		SourceBranch: m.SourceBranch,
		CommitSha:    m.Commit,
		Changelog:    strings.Join(m.Changelog, "\n"),
		IsDraft:      m.IsDraft,
		IsPreRelease: m.IsPreRelease,
		IsLatest:     m.IsLatest,
		Source:       model.SourceImport,
		CreatedAt:    m.CreatedAt,
	}
}
