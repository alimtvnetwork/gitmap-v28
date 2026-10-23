package store

import (
	"path/filepath"
	"testing"
	"time"
)

func setupTestCommitHistoryDB(t *testing.T) (*CommitHistorySplitDB, string) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "repodb", "commit_history.db")
	db, err := OpenCommitHistorySplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenCommitHistorySplitDBAt failed: %v", err)
	}

	return db, dbPath
}

func makeSampleCommitRecord(sha, shortSha, repoSlug, taskSlug string) CommitHistoryRecord {
	return CommitHistoryRecord{
		CommitSha:          sha,
		ShortSha:           shortSha,
		RepoName:           "gitmap",
		RepoSlug:           repoSlug,
		RepoRoot:           "/repos/gitmap",
		BranchName:         "main",
		AuthorName:         "Dev Agent",
		AuthorEmail:        "dev@example.com",
		CommitMessage:      "feat: commit tracking " + taskSlug,
		CommittedAtUnix:    time.Now().Unix(),
		AiModelName:        "gemini-2.5-pro",
		AiToolName:         "antigravity",
		AgentId:            "agy-worker-01",
		HeatmapSessionId:   "session-abc",
		HeatmapActionId:    42,
		TaskSlug:           taskSlug,
		SpecPath:           "02-spec/21-app/spec.md",
		PlansCount:         1,
		SubtasksCount:      6,
		IsAiCommitted:      true,
		IsCurrentMachineAi: true,
		Notes:              "automated tracking",
		FilesChanged:       2,
		LinesAdded:         150,
		LinesDeleted:       20,
		HasPushed:          false,
		CreatedAtUnix:      time.Now().Unix(),
	}
}

func makeSampleTouches(sha string) []CommitFileTouch {
	return []CommitFileTouch{
		{
			CommitSha:    sha,
			FilePath:     "cli/store/commit_history_split_db.go",
			ChangeType:   "M",
			LinesAdded:   100,
			LinesDeleted: 10,
		},
		{
			CommitSha:    sha,
			FilePath:     "cli/store/commit_history_ops.go",
			ChangeType:   "A",
			LinesAdded:   50,
			LinesDeleted: 10,
		},
	}
}

func TestCommitHistorySplitDB_Creation(t *testing.T) {
	db, path := setupTestCommitHistoryDB(t)
	defer db.Close()

	if db.Path != path {
		t.Fatalf("expected path %s, got %s", path, db.Path)
	}
	if !isFileExisting(path) {
		t.Fatalf("expected database file to exist at %s", path)
	}
}

func TestCommitHistorySplitDB_RecordAndGetBySha(t *testing.T) {
	db, _ := setupTestCommitHistoryDB(t)
	defer db.Close()

	rec := makeSampleCommitRecord("11223344556677889900aabbccddeeff11223344", "1122334", "alimtv/gitmap", "315-task")
	touches := makeSampleTouches(rec.CommitSha)
	if err := db.RecordCommit(rec, touches); err != nil {
		t.Fatalf("RecordCommit failed: %v", err)
	}

	assertRecordFoundBySha(t, db, rec.CommitSha)
	assertRecordFoundBySha(t, db, rec.ShortSha)
	assertRecordMissing(t, db, "unknown-sha")
}

func assertRecordFoundBySha(t *testing.T, db *CommitHistorySplitDB, sha string) {
	found, hasFound, err := db.GetBySha(sha)
	if err != nil {
		t.Fatalf("GetBySha(%s) returned error: %v", sha, err)
	}
	if !hasFound || found == nil {
		t.Fatalf("expected record for sha %s to be found", sha)
	}
	if !found.IsAiCommitted {
		t.Fatalf("expected IsAiCommitted true, got false")
	}
}

func assertRecordMissing(t *testing.T, db *CommitHistorySplitDB, sha string) {
	_, hasFound, err := db.GetBySha(sha)
	if err != nil {
		t.Fatalf("GetBySha unexpected error: %v", err)
	}
	if hasFound {
		t.Fatalf("expected hasFound false for missing sha")
	}
}

func TestCommitHistorySplitDB_SearchCommits(t *testing.T) {
	db, _ := setupTestCommitHistoryDB(t)
	defer db.Close()

	rec1 := makeSampleCommitRecord("abcdef1234567890abcdef1234567890abcdef12", "abcdef1", "alimtv/gitmap", "task-315")
	rec2 := makeSampleCommitRecord("987654fedcba0987654fedcba0987654fedcba09", "987654f", "alimtv/gitmap", "task-316")
	_ = db.RecordCommit(rec1, nil)
	_ = db.RecordCommit(rec2, nil)

	assertSearchResultsCount(t, db, "abcd", 1)
	assertSearchResultsCount(t, db, "cba09", 1)
	assertSearchResultsCount(t, db, "fedcba", 1)
	assertSearchResultsCount(t, db, "task-315", 1)
}

func assertSearchResultsCount(t *testing.T, db *CommitHistorySplitDB, query string, expected int) {
	results, err := db.SearchCommits(query, 10)
	if err != nil {
		t.Fatalf("SearchCommits(%s) failed: %v", query, err)
	}
	if len(results) != expected {
		t.Fatalf("SearchCommits(%s) expected %d results, got %d", query, expected, len(results))
	}
}

func TestCommitHistorySplitDB_ListFileTouches(t *testing.T) {
	db, _ := setupTestCommitHistoryDB(t)
	defer db.Close()

	rec := makeSampleCommitRecord("aabbccddee11223344556677889900aabbccddee", "aabbccd", "alimtv/gitmap", "task-315")
	touches := makeSampleTouches(rec.CommitSha)
	_ = db.RecordCommit(rec, touches)

	list, err := db.ListFileTouches(rec.CommitSha)
	if err != nil {
		t.Fatalf("ListFileTouches failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 touches, got %d", len(list))
	}
}

func TestCommitHistorySplitDB_ListRecentCommits(t *testing.T) {
	db, _ := setupTestCommitHistoryDB(t)
	defer db.Close()

	rec1 := makeSampleCommitRecord("1111111111111111111111111111111111111111", "1111111", "alimtv/gitmap", "task-1")
	rec1.CommittedAtUnix = 1000
	rec2 := makeSampleCommitRecord("2222222222222222222222222222222222222222", "2222222", "alimtv/gitmap", "task-2")
	rec2.CommittedAtUnix = 2000
	_ = db.RecordCommit(rec1, nil)
	_ = db.RecordCommit(rec2, nil)

	list, err := db.ListRecentCommits(10)
	if err != nil {
		t.Fatalf("ListRecentCommits failed: %v", err)
	}
	if len(list) != 2 || list[0].CommitSha != rec2.CommitSha {
		t.Fatalf("expected most recent commit first")
	}
}

func TestCommitHistorySplitDB_FilterCommits(t *testing.T) {
	db, _ := setupTestCommitHistoryDB(t)
	defer db.Close()

	rec := makeSampleCommitRecord("3333333333333333333333333333333333333333", "3333333", "alimtv/gitmap", "task-filtered")
	_ = db.RecordCommit(rec, nil)

	filter := CommitHistoryFilter{
		RepoSlug: "alimtv/gitmap",
		TaskSlug: "task-filtered",
		IsAiOnly: true,
		Limit:    10,
	}
	list, err := db.ListCommits(filter)
	if err != nil {
		t.Fatalf("ListCommits failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 result, got %d", len(list))
	}
}

func TestCommitHistorySplitDB_RegistrySync(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "root.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	db.syncCommitHistoryDB()

	entries, err := db.ListSplitDBs("")
	if err != nil {
		t.Fatalf("ListSplitDBs failed: %v", err)
	}

	assertRegistryEntryFound(t, entries, "commit_history")
}

func assertRegistryEntryFound(t *testing.T, entries []SplitDatabaseEntry, key string) {
	hasFound := false
	for _, e := range entries {
		if e.DatabaseKey == key {
			hasFound = true
			break
		}
	}

	if !hasFound {
		t.Fatalf("expected registry entry %s to be registered", key)
	}
}
