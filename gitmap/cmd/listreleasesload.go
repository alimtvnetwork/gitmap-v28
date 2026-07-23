package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// loadReleasesFromRepo reads .gitmap/release/v*.json files and converts to records.
func loadReleasesFromRepo() []model.ReleaseRecord {
	metas, err := release.ListReleaseMetaFiles()
	if err != nil || len(metas) == 0 {
		return nil
	}

	records := convertMetasToRecords(metas)
	sortRecordsByDate(records)
	markLatestRecord(records)

	return records
}

// convertMetasToRecords converts ReleaseMeta slices to ReleaseRecord slices.
func convertMetasToRecords(metas []release.ReleaseMeta) []model.ReleaseRecord {
	records := make([]model.ReleaseRecord, 0, len(metas))

	for _, m := range metas {
		records = append(records, metaToRecord(m))
	}

	return records
}

// metaToRecord converts a single ReleaseMeta to a ReleaseRecord.
func metaToRecord(m release.ReleaseMeta) model.ReleaseRecord {
	return model.ReleaseRecord{
		Version:      m.Version,
		Tag:          m.Tag,
		Branch:       m.Branch,
		SourceBranch: m.SourceBranch,
		CommitSha:    m.Commit,
		Changelog:    strings.Join(m.Changelog, "\n"),
		Notes:        m.Notes,
		IsDraft:      m.IsDraft,
		IsPreRelease: m.IsPreRelease,
		IsLatest:     m.IsLatest,
		Source:       model.SourceRepo,
		CreatedAt:    m.CreatedAt,
	}
}

// sortRecordsByDate sorts records by CreatedAt descending (newest first).
func sortRecordsByDate(records []model.ReleaseRecord) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt > records[j].CreatedAt
	})
}

// markLatestRecord sets IsLatest on the first record matching latest.json.
func markLatestRecord(records []model.ReleaseRecord) {
	latest, err := release.ReadLatest()
	if err != nil {
		return
	}

	for i := range records {
		if records[i].Tag == latest.Tag {
			records[i].IsLatest = true

			return
		}
	}
}

// loadReleasesFromDB opens the DB and fetches all releases.
func loadReleasesFromDB() []model.ReleaseRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrNoDatabase)
		os.Exit(1)
	}
	defer db.Close()

	releases, err := db.ListReleases()
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrListReleasesFailed, err)
		os.Exit(1)
	}

	return releases
}

// loadReleasesFromTags scans git tags and creates minimal records for tags
// that do not already have a matching metadata file or existing record.
func loadReleasesFromTags(existing []model.ReleaseRecord) []model.ReleaseRecord {
	tags := release.ListVersionTags()
	if len(tags) == 0 {
		return nil
	}

	seen := buildTagSet(existing)
	var added []model.ReleaseRecord

	for _, t := range tags {
		if seen[t.Tag] {
			continue
		}
		added = append(added, tagToRecord(t))
	}

	return added
}

// buildTagSet creates a set of tags from existing records.
func buildTagSet(records []model.ReleaseRecord) map[string]bool {
	m := make(map[string]bool, len(records))
	for _, r := range records {
		m[r.Tag] = true
	}

	return m
}

// tagToRecord creates a minimal ReleaseRecord from a git tag.
func tagToRecord(t release.TagEntry) model.ReleaseRecord {
	v, _ := release.Parse(t.Tag)
	branchName := constants.ReleaseBranchPrefix + v.String()

	return model.ReleaseRecord{
		Version:   v.String(),
		Tag:       t.Tag,
		Branch:    branchName,
		Source:    model.SourceTag,
		CreatedAt: t.CreatedAt,
	}
}

// cacheReleasesToDB upserts all records into the SQLite database.
// v17: stamps every record with the current repo's RepoId before upsert.
func cacheReleasesToDB(records []model.ReleaseRecord) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not cache releases to database: %v\n", err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Release cache DB migration failed: %v\n", err)
	}

	repoID, err := resolveCurrentRepoID(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Skipping DB cache: %v\n", err)

		return
	}
	upsertRecords(db, records, repoID)
}

// upsertRecords persists each record to the database, stamping RepoId.
func upsertRecords(db *store.DB, records []model.ReleaseRecord, repoID int64) {
	for _, r := range records {
		r.RepoID = repoID
		if err := db.UpsertRelease(r); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not cache release %s: %v\n", r.Version, err)
		}
	}
}
