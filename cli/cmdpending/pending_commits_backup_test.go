package cmdpending

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

func setupTestBackupDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test_backup.db")
	conn, err := OpenPendingCommitsBackup(dbPath)
	if err != nil {
		t.Fatalf("OpenPendingCommitsBackup failed: %v", err)
	}

	return conn, dbPath
}

func makeSampleRepoRecord(name, path, branch, sha string) RepoPendingCommitRecord {
	return RepoPendingCommitRecord{
		RepoName:             name,
		RelativePath:         path,
		CurrentBranch:        branch,
		Version:              "v1.0.0",
		ShortVersionBranch:   "v1.0.0 (main)",
		IsDirty:              true,
		IsClean:              false,
		TotalUncommitted:     4,
		UnpushedCommitsCount: 2,
	}
}

func verifyBackupFields(t *testing.T, rec *PendingCommitBackupRecord, exp RepoPendingCommitRecord, sha string, now int64) {
	t.Helper()
	if rec.RepoPath != exp.RelativePath || rec.RepoName != exp.RepoName {
		t.Errorf("repo path/name mismatch: got %s/%s", rec.RepoPath, rec.RepoName)
	}
	if rec.CurrentBranch != exp.CurrentBranch || rec.HeadSHA != sha {
		t.Errorf("branch/sha mismatch: got %s/%s", rec.CurrentBranch, rec.HeadSHA)
	}
	if rec.UncommittedCount != exp.TotalUncommitted || rec.UnpushedCount != exp.UnpushedCommitsCount {
		t.Errorf("counts mismatch: got %d/%d", rec.UncommittedCount, rec.UnpushedCount)
	}
	if rec.BackedUpAtUnix != now || !rec.IsDirty {
		t.Errorf("backup time or dirty mismatch: at=%d dirty=%v", rec.BackedUpAtUnix, rec.IsDirty)
	}
}

func TestPendingCommitsBackup_SaveAndGet(t *testing.T) {
	conn, _ := setupTestBackupDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	sample := makeSampleRepoRecord("repo-alpha", "repos/alpha", "main", "sha-alpha-1")
	if err := SaveBackupPendingStatus(conn, sample, "sha-alpha-1", now); err != nil {
		t.Fatalf("save backup failed: %v", err)
	}
	rec, hasRecord, err := GetBackupPendingStatus(conn, "repos/alpha")
	if err != nil || !hasRecord || rec == nil {
		t.Fatalf("expected record hit, got rec=%v hasRecord=%v err=%v", rec, hasRecord, err)
	}
	verifyBackupFields(t, rec, sample, "sha-alpha-1", now)
}

func TestPendingCommitsBackup_Upsert(t *testing.T) {
	conn, _ := setupTestBackupDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	first := makeSampleRepoRecord("repo-beta", "repos/beta", "main", "sha-beta-1")
	_ = SaveBackupPendingStatus(conn, first, "sha-beta-1", now-100)
	second := makeSampleRepoRecord("repo-beta", "repos/beta", "feature", "sha-beta-2")
	second.TotalUncommitted = 7
	second.UnpushedCommitsCount = 5
	if err := SaveBackupPendingStatus(conn, second, "sha-beta-2", now); err != nil {
		t.Fatalf("upsert backup failed: %v", err)
	}
	rec, hasRecord, _ := GetBackupPendingStatus(conn, "repos/beta")
	if !hasRecord || rec.UncommittedCount != 7 || rec.HeadSHA != "sha-beta-2" {
		t.Fatalf("unexpected upsert record state: %+v", rec)
	}
}

func TestPendingCommitsBackup_ListAll(t *testing.T) {
	conn, _ := setupTestBackupDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	r1 := makeSampleRepoRecord("r1", "repos/1", "main", "sha-1")
	r2 := makeSampleRepoRecord("r2", "repos/2", "main", "sha-2")
	_ = SaveBackupPendingStatus(conn, r1, "sha-1", now-10)
	_ = SaveBackupPendingStatus(conn, r2, "sha-2", now)
	records, err := ListAllBackupPendingStatuses(conn)
	if err != nil {
		t.Fatalf("list all failed: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestPendingCommitsBackup_RestoreBackupToCache(t *testing.T) {
	backupConn, _ := setupTestBackupDB(t)
	defer backupConn.Close()
	cacheConn, _ := setupTestCacheDB(t)
	defer cacheConn.Close()
	now := time.Now().Unix()
	sample := makeSampleRepoRecord("repo-gamma", "repos/gamma", "main", "sha-gamma-1")
	_ = SaveBackupPendingStatus(backupConn, sample, "sha-gamma-1", now-50)
	if err := RestoreBackupToCache(backupConn, cacheConn, "repos/gamma", now); err != nil {
		t.Fatalf("restore backup to cache failed: %v", err)
	}
	cacheRec, hasHit, err := GetCachedPendingStatus(cacheConn, "repos/gamma", "sha-gamma-1", now)
	if err != nil || !hasHit || cacheRec == nil {
		t.Fatalf("expected cache hit after restore, got rec=%v hit=%v err=%v", cacheRec, hasHit, err)
	}
}

func TestPendingCommitsBackup_RestoreMissingRecord(t *testing.T) {
	backupConn, _ := setupTestBackupDB(t)
	defer backupConn.Close()
	cacheConn, _ := setupTestCacheDB(t)
	defer cacheConn.Close()
	err := RestoreBackupToCache(backupConn, cacheConn, "nonexistent/repo", time.Now().Unix())
	if err == nil {
		t.Fatalf("expected error when restoring missing record")
	}
}

func TestPendingCommitsBackup_PurgeOlderThan(t *testing.T) {
	conn, _ := setupTestBackupDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	rOld := makeSampleRepoRecord("old-repo", "repos/old", "main", "sha-old")
	rNew := makeSampleRepoRecord("new-repo", "repos/new", "main", "sha-new")
	_ = SaveBackupPendingStatus(conn, rOld, "sha-old", now-200)
	_ = SaveBackupPendingStatus(conn, rNew, "sha-new", now-10)
	purged, err := PurgeBackupOlderThan(conn, now-100)
	if err != nil || purged != 1 {
		t.Fatalf("expected 1 purged row, got %d err=%v", purged, err)
	}
	_, hasOld, _ := GetBackupPendingStatus(conn, "repos/old")
	_, hasNew, _ := GetBackupPendingStatus(conn, "repos/new")
	if hasOld || !hasNew {
		t.Fatalf("expected old purged and new retained: old=%v new=%v", hasOld, hasNew)
	}
}

func verifyBackupNilConnHelpers(t *testing.T, now int64) {
	t.Helper()
	rec, hasRecord, err := GetBackupPendingStatus(nil, "repos/none")
	if err != nil || hasRecord || rec != nil {
		t.Errorf("GetBackupPendingStatus nil conn unexpected result")
	}
	records, err := ListAllBackupPendingStatuses(nil)
	if err != nil || records != nil {
		t.Errorf("ListAllBackupPendingStatuses nil conn unexpected result")
	}
	purged, err := PurgeBackupOlderThan(nil, now)
	if err != nil || purged != 0 {
		t.Errorf("PurgeBackupOlderThan nil conn unexpected result")
	}
}

func TestPendingCommitsBackup_NilConnection(t *testing.T) {
	now := time.Now().Unix()
	sample := makeSampleRepoRecord("nil-repo", "repos/nil", "main", "sha-nil")
	if err := SaveBackupPendingStatus(nil, sample, "sha-nil", now); err != nil {
		t.Errorf("SaveBackupPendingStatus nil conn error: %v", err)
	}
	if err := RestoreBackupToCache(nil, nil, "repos/nil", now); err != nil {
		t.Errorf("RestoreBackupToCache nil conn error: %v", err)
	}
	verifyBackupNilConnHelpers(t, now)
}

func TestPendingCommitsBackup_DefaultPath(t *testing.T) {
	custom := resolvePendingBackupDBPath("custom/backup.db")
	if custom != "custom/backup.db" {
		t.Errorf("expected custom backup path, got %s", custom)
	}
	defaultPath := resolvePendingBackupDBPath("")
	if defaultPath == "" {
		t.Errorf("expected non-empty default backup path")
	}
}

func makeSampleRecordWithPayload(name, path, branch, sha string) RepoPendingCommitRecord {
	rec := makeSampleRepoRecord(name, path, branch, sha)
	rec.PendingFiles = []string{"pkg/auth.go", "pkg/token.go"}
	rec.UnpushedCommitSHAs = []string{"abc1234", "def5678"}
	rec.UntrackedFilesCount = 1
	rec.ModifiedFilesCount = 1

	return rec
}

func verifyDeserializedPayload(t *testing.T, rec RepoPendingCommitRecord) {
	t.Helper()
	if len(rec.PendingFiles) != 2 || rec.PendingFiles[0] != "pkg/auth.go" {
		t.Fatalf("expected 2 pending files with pkg/auth.go, got %v", rec.PendingFiles)
	}
	if len(rec.UnpushedCommitSHAs) != 2 || rec.UnpushedCommitSHAs[0] != "abc1234" {
		t.Fatalf("expected 2 unpushed SHAs with abc1234, got %v", rec.UnpushedCommitSHAs)
	}
	if rec.UntrackedFilesCount != 1 || rec.ModifiedFilesCount != 1 {
		t.Fatalf("expected 1 untracked and 1 modified, got %d, %d", rec.UntrackedFilesCount, rec.ModifiedFilesCount)
	}
}

func TestPendingCommitsBackup_Deserialization_WithPayload(t *testing.T) {
	conn, _ := setupTestBackupDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	sample := makeSampleRecordWithPayload("repo-payload", "repos/payload", "main", "sha-pay-1")
	if err := SaveBackupPendingStatus(conn, sample, "sha-pay-1", now); err != nil {
		t.Fatalf("save backup failed: %v", err)
	}
	backupRec, hasRecord, err := GetBackupPendingStatus(conn, "repos/payload")
	if err != nil || !hasRecord || backupRec == nil {
		t.Fatalf("expected record hit, got %v", err)
	}
	converted := ConvertBackupToPendingRecord(*backupRec)
	verifyDeserializedPayload(t, converted)
}

func TestPendingCommitsBackup_ConvertBackupToPendingRecords_Bulk(t *testing.T) {
	conn, _ := setupTestBackupDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	s1 := makeSampleRecordWithPayload("repo-b1", "repos/b1", "main", "sha-1")
	s2 := makeSampleRecordWithPayload("repo-b2", "repos/b2", "main", "sha-2")
	_ = SaveBackupPendingStatus(conn, s1, "sha-1", now)
	_ = SaveBackupPendingStatus(conn, s2, "sha-2", now)
	list, err := ListAllBackupPendingStatuses(conn)
	if err != nil || len(list) != 2 {
		t.Fatalf("expected 2 backup records, got %d", len(list))
	}
	inspected := ConvertBackupToPendingRecords(list, "all")
	if len(inspected) != 2 {
		t.Fatalf("expected 2 inspected records, got %d", len(inspected))
	}
	verifyDeserializedPayload(t, inspected[0])
}

func TestPendingCommitsBackup_Deserialization_MalformedPayload(t *testing.T) {
	backupRec := PendingCommitBackupRecord{
		RepoPath:         "repos/malformed",
		RepoName:         "malformed",
		CurrentBranch:    "main",
		UncommittedCount: 3,
		UnpushedCount:    2,
		IsDirty:          true,
		PayloadJSON:      "{malformed-json",
	}
	converted := ConvertBackupToPendingRecord(backupRec)
	if converted.TotalUncommitted != 3 || converted.UnpushedCommitsCount != 2 {
		t.Fatalf("expected preserved metrics on malformed payload, got %d, %d", converted.TotalUncommitted, converted.UnpushedCommitsCount)
	}
	if !converted.IsDirty {
		t.Fatalf("expected IsDirty preserved")
	}
}

