package cmddb

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func initCacheTestSchema(t *testing.T, conn *sql.DB) {
	t.Helper()
	schema := `
	CREATE TABLE Users (id INTEGER PRIMARY KEY, name TEXT);
	INSERT INTO Users (name) VALUES ('Alice'), ('Bob');
	CREATE TABLE Posts (id INTEGER PRIMARY KEY, title TEXT);
	`
	if _, err := conn.Exec(schema); err != nil {
		t.Fatalf("init schema: %v", err)
	}
}

func createCacheTestDb(t *testing.T, filename string) (*sql.DB, string) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), filename)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	initCacheTestSchema(t, conn)

	return conn, dbPath
}

func TestTableCacheHitAndMiss(t *testing.T) {
	ClearAllTableCaches()
	conn, dbPath := createCacheTestDb(t, "hit_miss.db")
	defer conn.Close()

	tables1, err1 := GetCachedTablesAndViews(conn, dbPath)
	if err1 != nil {
		t.Fatalf("first call failed: %v", err1)
	}
	if len(tables1) != 2 {
		t.Fatalf("expected 2 tables, got %d", len(tables1))
	}

	key := normalizeDbCacheKey(dbPath)
	entry, hasEntry := lookupTableCache(key)
	if !hasEntry || entry == nil {
		t.Fatalf("expected cached entry in map for %s", key)
	}

	verifySecondCacheHit(t, conn, dbPath, entry)
}

func verifySecondCacheHit(t *testing.T, conn *sql.DB, dbPath string, firstEntry *CachedTableMetadata) {
	t.Helper()
	tables2, err2 := GetCachedTablesAndViews(conn, dbPath)
	if err2 != nil {
		t.Fatalf("second call failed: %v", err2)
	}
	if len(tables2) != 2 {
		t.Fatalf("expected 2 tables on second call, got %d", len(tables2))
	}

	key := normalizeDbCacheKey(dbPath)
	secondEntry, hasEntry := lookupTableCache(key)
	if !hasEntry || secondEntry != firstEntry {
		t.Fatalf("expected same cached entry pointer indicating cache hit")
	}
}

func TestTableCacheTtlExpiration(t *testing.T) {
	ClearAllTableCaches()
	conn, dbPath := createCacheTestDb(t, "ttl_expire.db")
	defer conn.Close()

	_, _ = GetCachedTablesAndViews(conn, dbPath)
	key := normalizeDbCacheKey(dbPath)
	entry, hasEntry := lookupTableCache(key)
	if !hasEntry {
		t.Fatalf("expected entry cached")
	}

	entry.LoadedAt = time.Now().Add(-15 * time.Second)
	hasValidTtl := isCacheTtlValid(entry.LoadedAt, DefaultTableCacheTTL)
	if hasValidTtl {
		t.Fatalf("expected TTL to be expired for 15s ago loadedAt")
	}

	isFresh := ValidateCacheFreshness(entry, conn, dbPath)
	if isFresh {
		t.Fatalf("expected freshness validation to fail on expired TTL")
	}
}

func TestTableCacheExternalFileMutation(t *testing.T) {
	ClearAllTableCaches()
	conn, dbPath := createCacheTestDb(t, "file_mod.db")
	defer conn.Close()

	_, _ = GetCachedTablesAndViews(conn, dbPath)
	key := normalizeDbCacheKey(dbPath)
	entry, hasEntry := lookupTableCache(key)
	if !hasEntry {
		t.Fatalf("expected cached entry")
	}

	futureTime := time.Now().Add(5 * time.Minute)
	if err := os.Chtimes(dbPath, futureTime, futureTime); err != nil {
		t.Fatalf("failed to update chtimes: %v", err)
	}

	isFresh := ValidateCacheFreshness(entry, conn, dbPath)
	if isFresh {
		t.Fatalf("expected freshness validation to fail after file modtime change")
	}
}

func externalWriteToDb(t *testing.T, dbPath string) {
	t.Helper()
	extConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open ext conn: %v", err)
	}
	defer extConn.Close()

	if _, err := extConn.Exec("INSERT INTO Users (name) VALUES ('Charlie');"); err != nil {
		t.Fatalf("failed to exec ext insert: %v", err)
	}
}

func TestTableCachePragmaDataVersion(t *testing.T) {
	ClearAllTableCaches()
	conn, dbPath := createCacheTestDb(t, "data_version.db")
	defer conn.Close()

	_, _ = GetCachedTablesAndViews(conn, dbPath)
	key := normalizeDbCacheKey(dbPath)
	entry, hasEntry := lookupTableCache(key)
	if !hasEntry {
		t.Fatalf("expected cached entry")
	}

	externalWriteToDb(t, dbPath)

	isFresh := ValidateCacheFreshness(entry, conn, dbPath)
	if isFresh {
		t.Fatalf("expected freshness to fail after external db write modified data_version")
	}
}

func TestTableCacheProactiveInvalidation(t *testing.T) {
	ClearAllTableCaches()
	conn, dbPath := createCacheTestDb(t, "proactive_inval.db")
	defer conn.Close()

	_, _ = GetCachedTablesAndViews(conn, dbPath)
	key := normalizeDbCacheKey(dbPath)
	if _, hasEntry := lookupTableCache(key); !hasEntry {
		t.Fatalf("expected entry cached")
	}

	InvalidateDbTableCache(dbPath)

	if _, hasEntry := lookupTableCache(key); hasEntry {
		t.Fatalf("expected entry removed after InvalidateDbTableCache")
	}

	_, _ = GetCachedTablesAndViews(conn, dbPath)
	ClearAllTableCaches()

	if _, hasEntry := lookupTableCache(key); hasEntry {
		t.Fatalf("expected entry removed after ClearAllTableCaches")
	}
}

func BenchmarkTableCacheFreshnessValidation(b *testing.B) {
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "bench.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		b.Fatalf("open sqlite: %v", err)
	}
	defer conn.Close()

	_, _ = conn.Exec("CREATE TABLE T (id INTEGER);")
	_, _ = GetCachedTablesAndViews(conn, dbPath)
	key := normalizeDbCacheKey(dbPath)
	entry, _ := lookupTableCache(key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateCacheFreshness(entry, conn, dbPath)
	}
}
