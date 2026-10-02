package store

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestFindDuplicateRepos(t *testing.T) {
	db := openTempDB(t)

	records := []model.ScanRecord{
		{
			Slug:         "repo-a-1",
			RepoName:     "repo-a",
			HTTPSUrl:     "https://github.com/org/repo-a",
			AbsolutePath: "/tmp/test/repo-a-1",
			Branch:       "main",
		},
		{
			Slug:         "repo-a-2",
			RepoName:     "repo-a",
			HTTPSUrl:     "https://github.com/org/repo-a.git/",
			AbsolutePath: "/tmp/test/repo-a-2",
			Branch:       "main",
		},
		{
			Slug:         "repo-b-1",
			RepoName:     "repo-b",
			SSHUrl:       "git@github.com:org/repo-b.git",
			AbsolutePath: "/tmp/test/repo-b-1",
			Branch:       "main",
		},
		{
			Slug:         "repo-b-2",
			RepoName:     "repo-b",
			SSHUrl:       "git@github.com:org/repo-b",
			AbsolutePath: "/tmp/test/repo-b-2",
			Branch:       "main",
		},
		{
			Slug:         "repo-c",
			RepoName:     "repo-c",
			HTTPSUrl:     "https://github.com/org/repo-c",
			AbsolutePath: "/tmp/test/repo-c",
			Branch:       "main",
		},
	}

	if err := db.UpsertRepos(records); err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}

	groups, err := db.FindDuplicateRepos()
	if err != nil {
		t.Fatalf("FindDuplicateRepos failed: %v", err)
	}

	if len(groups) != 2 {
		t.Fatalf("expected 2 duplicate groups, got %d", len(groups))
	}

	for _, g := range groups {
		if g.Count != 2 {
			t.Errorf("group %q expected count 2, got %d", g.CleanKey, g.Count)
		}
		if len(g.Duplicates) != 1 {
			t.Errorf("group %q expected 1 duplicate record, got %d", g.CleanKey, len(g.Duplicates))
		}
		if g.Keeper.ID == 0 {
			t.Errorf("group %q keeper ID is 0", g.CleanKey)
		}
	}
}

func TestDeduplicateRepos(t *testing.T) {
	db := openTempDB(t)

	records := []model.ScanRecord{
		{
			Slug:         "repo-dup-1",
			RepoName:     "repo-dup",
			HTTPSUrl:     "https://github.com/example/dup",
			AbsolutePath: "/tmp/dup1",
			Branch:       "main",
		},
		{
			Slug:         "repo-dup-2",
			RepoName:     "repo-dup",
			HTTPSUrl:     "https://github.com/example/dup.git",
			AbsolutePath: "/tmp/dup2",
			Branch:       "main",
		},
	}

	if err := db.UpsertRepos(records); err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}

	// Verify initial repos
	all, err := db.ListRepos()
	if err != nil || len(all) != 2 {
		t.Fatalf("expected 2 repos initially, got %d (err: %v)", len(all), err)
	}

	keeperID := all[0].ID
	dupID := all[1].ID

	// Create a group and link duplicate repo
	group, err := db.CreateGroup("backend", "backend services", "blue")
	if err != nil {
		t.Fatalf("CreateGroup failed: %v", err)
	}
	if err := db.AddRepoToGroup(group.Name, dupID); err != nil {
		t.Fatalf("AddRepoToGroup failed: %v", err)
	}

	// Add release to duplicate repo
	rel := model.ReleaseRecord{
		RepoID:       dupID,
		Version:      "v1.0.0",
		Tag:          "v1.0.0",
		Branch:       "main",
		SourceBranch: "main",
		CommitSha:    "abc1234",
	}
	if err := db.UpsertRelease(rel); err != nil {
		t.Fatalf("UpsertRelease failed: %v", err)
	}

	// Run DeduplicateRepos
	summary, err := db.DeduplicateRepos(false)
	if err != nil {
		t.Fatalf("DeduplicateRepos failed: %v", err)
	}

	if summary.GroupsFound != 1 {
		t.Errorf("expected 1 group found, got %d", summary.GroupsFound)
	}
	if summary.RowsPurged != 1 {
		t.Errorf("expected 1 row purged, got %d", summary.RowsPurged)
	}

	// Verify only 1 repo remains and it is the keeper
	remaining, err := db.ListRepos()
	if err != nil || len(remaining) != 1 {
		t.Fatalf("expected 1 repo remaining, got %d (err: %v)", len(remaining), err)
	}
	if remaining[0].ID != keeperID {
		t.Errorf("expected remaining repo ID %d, got %d", keeperID, remaining[0].ID)
	}

	// Verify GroupRepo remapped to keeper
	groupRepos, err := db.ShowGroup(group.Name)
	if err != nil {
		t.Fatalf("ShowGroup failed: %v", err)
	}
	if len(groupRepos) != 1 || groupRepos[0].ID != keeperID {
		t.Errorf("expected group to contain keeper repo %d, got %+v", keeperID, groupRepos)
	}

	// Verify Release remapped to keeper
	releases, err := db.ListReleases()
	if err != nil {
		t.Fatalf("ListReleases failed: %v", err)
	}
	if len(releases) != 1 || releases[0].RepoID != keeperID {
		t.Errorf("expected release to be remapped to keeper repo %d, got %+v", keeperID, releases)
	}

	// Idempotent re-run should find 0 groups
	summary2, err := db.DeduplicateRepos(false)
	if err != nil {
		t.Fatalf("Second DeduplicateRepos failed: %v", err)
	}
	if summary2.GroupsFound != 0 || summary2.RowsPurged != 0 {
		t.Errorf("expected 0 groups and 0 rows on re-run, got %d groups, %d rows", summary2.GroupsFound, summary2.RowsPurged)
	}
}
