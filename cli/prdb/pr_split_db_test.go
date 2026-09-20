package prdb

import (
	"path/filepath"
	"strings"
	"testing"
)

func setupTestPrDb(t *testing.T) *PrSplitDb {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "sql.db")
	res := OpenPrSplitDbAt(dbPath)
	if res.IsFailure() {
		t.Fatalf("failed to open test pr db: %v", res.Err)
	}

	return res.Value
}

func TestOpenPrSplitDb_SchemaAndPragmas(t *testing.T) {
	db := setupTestPrDb(t)
	defer db.Close()

	if db.Conn() == nil {
		t.Fatal("expected non-nil db connection")
	}

	var journalMode string
	if err := db.Conn().QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}
	if strings.ToLower(journalMode) != "wal" {
		t.Errorf("expected WAL journal_mode, got %s", journalMode)
	}
}

func testCreatePullRequest(t *testing.T, db *PrSplitDb) int64 {
	rec := PullRequestRecord{
		PrNumber:     101,
		Title:        "Feature Alpha",
		Description:  "Initial PR implementation",
		SourceBranch: "feature/alpha",
		TargetBranch: "main",
	}
	res := db.CreatePullRequest(rec)
	if res.IsFailure() || res.Value <= 0 {
		t.Fatalf("CreatePullRequest failed: %v", res.Err)
	}

	return res.Value
}

func testUpdateAndVerifyPR(t *testing.T, db *PrSplitDb) {
	updateRes := db.UpdatePullRequestStatus(101, "merged", "abc1234")
	if updateRes.IsFailure() || !updateRes.Value {
		t.Fatalf("UpdatePullRequestStatus failed: %v", updateRes.Err)
	}

	getRes := db.GetPullRequestByNumber(101)
	if getRes.IsFailure() || getRes.Value == nil {
		t.Fatalf("GetPullRequestByNumber failed: %v", getRes.Err)
	}
	pr := getRes.Value
	if pr.Status != "merged" || pr.MergeCommitSha != "abc1234" || pr.MergedAt <= 0 {
		t.Errorf("unexpected PR state: %+v", pr)
	}
}

func testNonExistentPR(t *testing.T, db *PrSplitDb) {
	getRes := db.GetPullRequestByNumber(999)
	if getRes.IsFailure() {
		t.Fatalf("GetPullRequestByNumber for missing PR failed: %v", getRes.Err)
	}
	if getRes.Value != nil {
		t.Errorf("expected nil PR record for missing PR, got %+v", getRes.Value)
	}
}

func TestPullRequest_Lifecycle(t *testing.T) {
	db := setupTestPrDb(t)
	defer db.Close()

	_ = testCreatePullRequest(t, db)
	testUpdateAndVerifyPR(t, db)
	testNonExistentPR(t, db)
}

func TestPrRelease_CRUD(t *testing.T) {
	db := setupTestPrDb(t)
	defer db.Close()

	prId := testCreatePullRequest(t, db)
	relRec := PrReleaseRecord{
		PullRequestId: prId,
		ReleaseTag:    "v1.0.0",
		CommitSha:     "deadbeef12345678",
		Notes:         "Release notes for v1.0.0",
	}
	res := db.AddPrRelease(relRec)
	if res.IsFailure() || res.Value <= 0 {
		t.Fatalf("AddPrRelease failed: %v", res.Err)
	}
}

func testUpsertAndListActiveBranches(t *testing.T, db *PrSplitDb) {
	rec := PrBranchRecord{
		BranchName: "feature/alpha",
		BranchType: "feature",
	}
	upsertRes := db.UpsertPrBranch(rec)
	if upsertRes.IsFailure() || !upsertRes.Value {
		t.Fatalf("UpsertPrBranch failed: %v", upsertRes.Err)
	}

	listRes := db.ListActivePrBranches()
	if listRes.IsFailure() || len(listRes.Value) == 0 {
		t.Fatalf("ListActivePrBranches failed: %v", listRes.Err)
	}
	if listRes.Value[0].BranchName != "feature/alpha" {
		t.Errorf("expected feature/alpha, got %s", listRes.Value[0].BranchName)
	}
}

func testMarkBranchDeletedFlow(t *testing.T, db *PrSplitDb) {
	delRes := db.MarkPrBranchDeleted("feature/alpha")
	if delRes.IsFailure() || !delRes.Value {
		t.Fatalf("MarkPrBranchDeleted failed: %v", delRes.Err)
	}

	listRes := db.ListActivePrBranches()
	if listRes.IsFailure() {
		t.Fatalf("ListActivePrBranches after delete failed: %v", listRes.Err)
	}
	if len(listRes.Value) != 0 {
		t.Errorf("expected 0 active branches, got %d", len(listRes.Value))
	}
}

func testMergedBranchesFlow(t *testing.T, db *PrSplitDb) {
	rec := PrBranchRecord{
		BranchName: "feature/beta",
		BranchType: "feature",
		IsMerged:   true,
	}
	_ = db.UpsertPrBranch(rec)

	activeList := db.ListActivePrBranches()
	if len(activeList.Value) != 0 {
		t.Errorf("expected 0 active branches for merged branch, got %d", len(activeList.Value))
	}

	mergedList := db.ListMergedPrBranches()
	if len(mergedList.Value) != 1 || mergedList.Value[0].BranchName != "feature/beta" {
		t.Errorf("expected 1 merged branch feature/beta, got %+v", mergedList.Value)
	}
}

func TestPrBranch_Lifecycle(t *testing.T) {
	db := setupTestPrDb(t)
	defer db.Close()

	testUpsertAndListActiveBranches(t, db)
	testMarkBranchDeletedFlow(t, db)
	testMergedBranchesFlow(t, db)
}

func TestPrViews_Compatibility(t *testing.T) {
	db := setupTestPrDb(t)
	defer db.Close()

	_ = testCreatePullRequest(t, db)
	var (
		viewId   int64
		viewPrNo int
		viewTit  string
	)
	err := db.Conn().QueryRow("SELECT id, pr_number, title FROM pull_requests WHERE pr_number = 101;").
		Scan(&viewId, &viewPrNo, &viewTit)
	if err != nil {
		t.Fatalf("query pull_requests view failed: %v", err)
	}
	if viewPrNo != 101 || viewTit != "Feature Alpha" {
		t.Errorf("unexpected view values: id=%d pr_number=%d title=%s", viewId, viewPrNo, viewTit)
	}
}

func TestSanitizeRepoSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"github.com/alimtvnetwork/gitmap-v28.git", "alimtvnetwork-gitmap-v28"},
		{"gitlab.com:org/repo.git", "org-repo"},
		{"My-Project_Repo!", "my-project-repo"},
		{"", "pr-default"},
	}

	for _, tc := range tests {
		got := SanitizeRepoSlug(tc.input)
		if got != tc.expected {
			t.Errorf("SanitizeRepoSlug(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestResolvePrDbPath(t *testing.T) {
	repoPath := ResolvePrDbPath("owner/repo", "/my/repo")
	if !strings.Contains(repoPath, ".gitmap/data/pr/owner-repo/sql.db") {
		t.Errorf("unexpected repo path: %s", repoPath)
	}

	globalPath := ResolvePrDbPath("owner/repo", "")
	if !strings.Contains(globalPath, "pr/owner-repo/sql.db") {
		t.Errorf("unexpected global path: %s", globalPath)
	}
}
