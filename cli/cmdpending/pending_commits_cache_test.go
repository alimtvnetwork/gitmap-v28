package cmdpending

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

func setupTestCacheDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test_cache.db")
	conn, err := OpenPendingCommitsCache(dbPath)
	if err != nil {
		t.Fatalf("OpenPendingCommitsCache failed: %v", err)
	}

	return conn, dbPath
}

func queryTableRowCount(t *testing.T, conn *sql.DB, repoPath string) int {
	t.Helper()
	var count int
	query := "SELECT COUNT(*) FROM pending_commits_cache WHERE repo_path = ?;"
	if err := conn.QueryRow(query, repoPath).Scan(&count); err != nil {
		t.Fatalf("query row count failed: %v", err)
	}

	return count
}

func queryTableTotalRows(t *testing.T, conn *sql.DB) int {
	t.Helper()
	var count int
	query := "SELECT COUNT(*) FROM pending_commits_cache;"
	if err := conn.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("query total rows failed: %v", err)
	}

	return count
}

func makeSampleRecord(repoPath, headSHA string, cachedAt int64) PendingCommitCacheRecord {
	return PendingCommitCacheRecord{
		RepoPath:         repoPath,
		HeadSHA:          headSHA,
		UncommittedCount: 3,
		UnpushedCount:    2,
		IsDirty:          true,
		CachedAtUnix:     cachedAt,
	}
}

func verifyRecordMatch(t *testing.T, rec *PendingCommitCacheRecord, expected PendingCommitCacheRecord) {
	t.Helper()
	if rec.RepoPath != expected.RepoPath || rec.HeadSHA != expected.HeadSHA {
		t.Errorf("repo/head mismatch: got %s/%s, want %s/%s", rec.RepoPath, rec.HeadSHA, expected.RepoPath, expected.HeadSHA)
	}
	if rec.UncommittedCount != expected.UncommittedCount || rec.UnpushedCount != expected.UnpushedCount {
		t.Errorf("counts mismatch: got %d/%d, want %d/%d", rec.UncommittedCount, rec.UnpushedCount, expected.UncommittedCount, expected.UnpushedCount)
	}
	if rec.IsDirty != expected.IsDirty {
		t.Errorf("isDirty mismatch: got %v, want %v", rec.IsDirty, expected.IsDirty)
	}
}

func assertCacheMiss(t *testing.T, rec *PendingCommitCacheRecord, hasHit bool) {
	t.Helper()
	if hasHit {
		t.Fatalf("expected cache miss for expired entry")
	}
	if rec != nil {
		t.Fatalf("expected nil record on miss")
	}
}

func assertRepoRowDeleted(t *testing.T, conn *sql.DB, repoPath string) {
	t.Helper()
	count := queryTableRowCount(t, conn, repoPath)
	if count > 0 {
		t.Fatalf("expected immediate purge of expired record, but row exists")
	}
}

func TestPendingCommitsCache_HitWithinTTL(t *testing.T) {
	conn, _ := setupTestCacheDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	expected := makeSampleRecord("repo/alpha", "commit-sha-123", now-30)
	if err := SaveCachedPendingStatus(conn, expected); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	rec, hasHit, err := GetCachedPendingStatus(conn, "repo/alpha", "commit-sha-123", now)
	if err != nil {
		t.Fatalf("get cached failed: %v", err)
	}
	if !hasHit {
		t.Fatalf("expected cache hit within TTL")
	}
	verifyRecordMatch(t, rec, expected)
}

func TestPendingCommitsCache_ExpiredPastTTL_Purged(t *testing.T) {
	conn, _ := setupTestCacheDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	sample := makeSampleRecord("repo/beta", "commit-sha-old", now-95)
	if err := SaveCachedPendingStatus(conn, sample); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	rec, hasHit, err := GetCachedPendingStatus(conn, "repo/beta", "commit-sha-old", now)
	if err != nil {
		t.Fatalf("get cached failed: %v", err)
	}
	assertCacheMiss(t, rec, hasHit)
	assertRepoRowDeleted(t, conn, "repo/beta")
}

func TestPendingCommitsCache_ClockSkew_Purged(t *testing.T) {
	conn, _ := setupTestCacheDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	sample := makeSampleRecord("repo/future", "commit-sha-skew", now+100)
	if err := SaveCachedPendingStatus(conn, sample); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	rec, hasHit, err := GetCachedPendingStatus(conn, "repo/future", "commit-sha-skew", now)
	if err != nil {
		t.Fatalf("get cached failed: %v", err)
	}
	assertCacheMiss(t, rec, hasHit)
	assertRepoRowDeleted(t, conn, "repo/future")
}

func TestPendingCommitsCache_HeadSHAMismatch_Invalidated(t *testing.T) {
	conn, _ := setupTestCacheDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	sample := makeSampleRecord("repo/gamma", "commit-sha-orig", now-10)
	if err := SaveCachedPendingStatus(conn, sample); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	rec, hasHit, err := GetCachedPendingStatus(conn, "repo/gamma", "commit-sha-new", now)
	if err != nil {
		t.Fatalf("get cached failed: %v", err)
	}
	assertCacheMiss(t, rec, hasHit)
	assertRepoRowDeleted(t, conn, "repo/gamma")
}

func seedBulkRecords(t *testing.T, conn *sql.DB, now int64) {
	t.Helper()
	r1 := makeSampleRecord("repo/bulk1", "sha1", now-120)
	r2 := makeSampleRecord("repo/bulk2", "sha2", now-100)
	r3 := makeSampleRecord("repo/bulk3", "sha3", now-10)
	_ = SaveCachedPendingStatus(conn, r1)
	_ = SaveCachedPendingStatus(conn, r2)
	_ = SaveCachedPendingStatus(conn, r3)
}

func verifyRemainingBulkRows(t *testing.T, conn *sql.DB) {
	t.Helper()
	if count := queryTableTotalRows(t, conn); count != 1 {
		t.Fatalf("expected 1 remaining row, got %d", count)
	}
	if count := queryTableRowCount(t, conn, "repo/bulk3"); count != 1 {
		t.Fatalf("expected fresh record bulk3 to remain")
	}
}

func TestPendingCommitsCache_PurgeExpiredBulk(t *testing.T) {
	conn, _ := setupTestCacheDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	seedBulkRecords(t, conn, now)
	purged, err := PurgeExpiredPendingCache(conn, now, DefaultCacheTTLSeconds)
	if err != nil {
		t.Fatalf("purge failed: %v", err)
	}
	if purged != 2 {
		t.Fatalf("expected 2 purged records, got %d", purged)
	}
	verifyRemainingBulkRows(t, conn)
}

func verifyUpsertResult(t *testing.T, conn *sql.DB, repoPath string, expected PendingCommitCacheRecord, now int64) {
	t.Helper()
	if count := queryTableTotalRows(t, conn); count != 1 {
		t.Fatalf("expected exactly 1 row after upsert, got %d", count)
	}
	rec, hasHit, err := GetCachedPendingStatus(conn, repoPath, expected.HeadSHA, now)
	if err != nil {
		t.Fatalf("query after upsert failed: %v", err)
	}
	if !hasHit {
		t.Fatalf("expected hit on updated record")
	}
	verifyRecordMatch(t, rec, expected)
}

func TestPendingCommitsCache_Upsert(t *testing.T) {
	conn, _ := setupTestCacheDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	first := makeSampleRecord("repo/delta", "sha-initial", now-20)
	_ = SaveCachedPendingStatus(conn, first)
	second := PendingCommitCacheRecord{
		RepoPath:         "repo/delta",
		HeadSHA:          "sha-updated",
		UncommittedCount: 10,
		UnpushedCount:    5,
		IsDirty:          true,
		CachedAtUnix:     now,
	}
	_ = SaveCachedPendingStatus(conn, second)
	verifyUpsertResult(t, conn, "repo/delta", second, now)
}

func TestPendingCommitsCache_InvalidatePendingCache(t *testing.T) {
	conn, _ := setupTestCacheDB(t)
	defer conn.Close()
	now := time.Now().Unix()
	sample := makeSampleRecord("repo/invalidate", "sha-inv", now-10)
	_ = SaveCachedPendingStatus(conn, sample)
	if err := InvalidatePendingCache(conn, "repo/invalidate"); err != nil {
		t.Fatalf("invalidate failed: %v", err)
	}
	assertRepoRowDeleted(t, conn, "repo/invalidate")
}

func verifyNilConnMutations(t *testing.T, now int64) {
	t.Helper()
	if err := SaveCachedPendingStatus(nil, PendingCommitCacheRecord{}); err != nil {
		t.Errorf("SaveCachedPendingStatus nil conn: %v", err)
	}
	rows, err := PurgeExpiredPendingCache(nil, now, 90)
	if err != nil {
		t.Errorf("PurgeExpiredPendingCache nil conn: %v", err)
	}
	if rows != 0 {
		t.Errorf("PurgeExpiredPendingCache nil conn rows=%d", rows)
	}
	if err := InvalidatePendingCache(nil, "repo/nil"); err != nil {
		t.Errorf("InvalidatePendingCache nil conn: %v", err)
	}
}

func TestPendingCommitsCache_NilConnection(t *testing.T) {
	now := time.Now().Unix()
	rec, hasHit, err := GetCachedPendingStatus(nil, "repo/nil", "sha", now)
	if err != nil {
		t.Errorf("GetCachedPendingStatus nil conn error: %v", err)
	}
	if hasHit {
		t.Errorf("GetCachedPendingStatus nil conn unexpected hit")
	}
	if rec != nil {
		t.Errorf("GetCachedPendingStatus nil conn unexpected record")
	}
	verifyNilConnMutations(t, now)
}

func TestPendingCommitsCache_DefaultPath(t *testing.T) {
	custom := resolvePendingCacheDBPath("custom/path.db")
	if custom != "custom/path.db" {
		t.Errorf("expected custom path, got %s", custom)
	}
	defaultPath := resolvePendingCacheDBPath("")
	if defaultPath == "" {
		t.Errorf("expected non-empty default path")
	}
}

func assertDualWriteResults(t *testing.T, cacheConn, backupConn *sql.DB, path, sha string, now int64) {
	t.Helper()
	cRec, hasCache, _ := GetCachedPendingStatus(cacheConn, path, sha, now)
	if !hasCache || cRec == nil {
		t.Fatalf("expected cache record to exist after dual write")
	}
	bRec, hasBackup, _ := GetBackupPendingStatus(backupConn, path)
	if !hasBackup || bRec == nil {
		t.Fatalf("expected backup record to exist after dual write")
	}
}

func TestPendingCommitsCache_SaveStatusWithBackup_DualWrite(t *testing.T) {
	cacheConn, _ := setupTestCacheDB(t)
	defer cacheConn.Close()
	backupConn, _ := setupTestBackupDB(t)
	defer backupConn.Close()
	now := time.Now().Unix()
	cacheRec := makeSampleRecord("repo/dual", "sha-dual", now)
	fullRec := makeSampleRepoRecord("dual", "repo/dual", "main", "sha-dual")
	if err := SaveStatusWithBackup(cacheConn, backupConn, cacheRec, fullRec, now); err != nil {
		t.Fatalf("dual write failed: %v", err)
	}
	assertDualWriteResults(t, cacheConn, backupConn, "repo/dual", "sha-dual", now)
}

func TestPendingCommitsCache_SaveStatusWithBackup_NilBackup(t *testing.T) {
	cacheConn, _ := setupTestCacheDB(t)
	defer cacheConn.Close()
	now := time.Now().Unix()
	cacheRec := makeSampleRecord("repo/cacheonly", "sha-c1", now)
	fullRec := makeSampleRepoRecord("cacheonly", "repo/cacheonly", "main", "sha-c1")
	if err := SaveStatusWithBackup(cacheConn, nil, cacheRec, fullRec, now); err != nil {
		t.Fatalf("cache only write failed: %v", err)
	}
	rec, hasCache, _ := GetCachedPendingStatus(cacheConn, "repo/cacheonly", "sha-c1", now)
	if !hasCache || rec == nil {
		t.Fatalf("expected cache record to exist")
	}
}

func TestPendingCommitsCache_SaveStatusWithBackup_NilConns(t *testing.T) {
	now := time.Now().Unix()
	cacheRec := makeSampleRecord("repo/nilwrite", "sha-nil", now)
	fullRec := makeSampleRepoRecord("nilwrite", "repo/nilwrite", "main", "sha-nil")
	if err := SaveStatusWithBackup(nil, nil, cacheRec, fullRec, now); err != nil {
		t.Fatalf("nil write failed: %v", err)
	}
}

func seedExpiredAndFreshRecords(t *testing.T, conn *sql.DB, now int64) {
	t.Helper()
	expiredRec := makeSampleRecord("repo/expired-on-open", "sha-expired", now-100)
	freshRec := makeSampleRecord("repo/fresh-on-open", "sha-fresh", now-10)
	if err := SaveCachedPendingStatus(conn, expiredRec); err != nil {
		t.Fatalf("save expiredRec failed: %v", err)
	}
	if err := SaveCachedPendingStatus(conn, freshRec); err != nil {
		t.Fatalf("save freshRec failed: %v", err)
	}
}

func verifyOpenPurgeResults(t *testing.T, conn *sql.DB) {
	t.Helper()
	if count := queryTableRowCount(t, conn, "repo/expired-on-open"); count != 0 {
		t.Fatalf("expected expired row to be purged on open, got count %d", count)
	}
	if count := queryTableRowCount(t, conn, "repo/fresh-on-open"); count != 1 {
		t.Fatalf("expected fresh row to remain on open, got count %d", count)
	}
}

func TestPendingCommitsCache_OpenPurgesExpiredOnInit(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_open_purge.db")
	conn1, err := OpenPendingCommitsCache(dbPath)
	if err != nil {
		t.Fatalf("OpenPendingCommitsCache failed: %v", err)
	}
	now := time.Now().Unix()
	seedExpiredAndFreshRecords(t, conn1, now)
	_ = conn1.Close()

	conn2, err := OpenPendingCommitsCache(dbPath)
	if err != nil {
		t.Fatalf("reopen OpenPendingCommitsCache failed: %v", err)
	}
	defer conn2.Close()
	verifyOpenPurgeResults(t, conn2)
}


