package cmdurl

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestHashUrl(t *testing.T) {
	url := "https://example.com/api/test"
	hash1 := store.HashUrl(url)
	hash2 := store.HashUrl(url)
	if hash1 != hash2 {
		t.Fatalf("expected deterministic hash, got %s and %s", hash1, hash2)
	}

	if len(hash1) != 64 {
		t.Fatalf("expected 64 char hex hash, got %d chars", len(hash1))
	}
}

func TestParseUrlArgs_ImplicitRead(t *testing.T) {
	args := []string{"https://example.com/data.json", "--fresh", "--json"}
	opts, err := ParseUrlArgs(args)
	if err != nil {
		t.Fatalf("unexpected error parsing args: %v", err)
	}

	if opts.Subcommand != "read" || opts.URL != "https://example.com/data.json" {
		t.Fatalf("expected read subcommand with url, got %s, %s", opts.Subcommand, opts.URL)
	}

	if !opts.IsForceFresh || !opts.IsJSON {
		t.Fatalf("expected IsForceFresh and IsJSON true")
	}
}

func TestParseUrlArgs_ExplicitReadWithTTL(t *testing.T) {
	args := []string{"read", "https://example.com", "--ttl", "6h", "-q"}
	opts, err := ParseUrlArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Subcommand != "read" || opts.URL != "https://example.com" {
		t.Fatalf("unexpected subcommand or url: %s, %s", opts.Subcommand, opts.URL)
	}

	if opts.TTL != 6*time.Hour || !opts.IsQuiet {
		t.Fatalf("expected TTL 6h and IsQuiet true, got %v, %v", opts.TTL, opts.IsQuiet)
	}
}

func TestParseUrlArgs_Subcommands(t *testing.T) {
	assertSubcommand(t, []string{"ls", "--limit", "10", "--json"}, "ls")
	assertSubcommand(t, []string{"info", "https://example.com"}, "info")
	assertSubcommand(t, []string{"purge", "--all"}, "purge")
	assertSubcommand(t, []string{"rm", "https://example.com"}, "rm")
	assertSubcommand(t, []string{"help"}, "help")
}

func assertSubcommand(t *testing.T, args []string, expected string) {
	opts, err := ParseUrlArgs(args)
	if err != nil {
		t.Fatalf("unexpected error parsing %v: %v", args, err)
	}

	if opts.Subcommand != expected {
		t.Fatalf("expected subcommand %s, got %s", expected, opts.Subcommand)
	}
}

func TestUrlCacheDB_HitAndMiss(t *testing.T) {
	db := createTestUrlCacheDB(t)
	defer db.Close()

	testUrl := "https://example.com/resource"
	assertCacheMiss(t, db, testUrl)
	saveTestEntry(t, db, testUrl, []byte("test payload"), 3600)
	assertCacheHit(t, db, testUrl, "test payload")
}

func assertCacheMiss(t *testing.T, db *store.UrlCacheSplitDB, testUrl string) {
	_, hasHit, err := db.GetUrlCache(testUrl)
	if err != nil || hasHit {
		t.Fatalf("expected miss for uncached url, got hit=%v err=%v", hasHit, err)
	}
}

func assertCacheHit(t *testing.T, db *store.UrlCacheSplitDB, testUrl, expected string) {
	entry, hasHit, err := db.GetUrlCache(testUrl)
	if err != nil || !hasHit {
		t.Fatalf("expected hit for cached url, got hit=%v err=%v", hasHit, err)
	}

	if string(entry.PayloadBlob) != expected {
		t.Fatalf("unexpected payload: %s", string(entry.PayloadBlob))
	}
}

func TestUrlCacheDB_ExpirationAndPurge(t *testing.T) {
	db := createTestUrlCacheDB(t)
	defer db.Close()

	seedExpirationEntries(t, db)
	purged, err := db.PurgeExpiredUrlCache()
	if err != nil || purged != 1 {
		t.Fatalf("expected 1 purged entry, got %d, err=%v", purged, err)
	}

	assertUrlExistence(t, db, "https://example.com/expired", false)
	assertUrlExistence(t, db, "https://example.com/fresh", true)
}

func seedExpirationEntries(t *testing.T, db *store.UrlCacheSplitDB) {
	saveTestEntry(t, db, "https://example.com/expired", []byte("expired payload"), -60)
	saveTestEntry(t, db, "https://example.com/fresh", []byte("fresh payload"), 3600)
}

func TestUrlCacheDB_DeleteAndTouch(t *testing.T) {
	db := createTestUrlCacheDB(t)
	defer db.Close()

	testUrl := "https://example.com/touch"
	saveTestEntry(t, db, testUrl, []byte("touch payload"), 3600)
	verifyTouchAndDeletion(t, db, testUrl)
}

func verifyTouchAndDeletion(t *testing.T, db *store.UrlCacheSplitDB, testUrl string) {
	entry, hasHit, _ := db.GetUrlCache(testUrl)
	if !hasHit {
		t.Fatalf("expected entry to exist")
	}

	touchEntry(t, db, entry.UrlCacheEntryId)
	deleteAndVerify(t, db, testUrl)
}

func touchEntry(t *testing.T, db *store.UrlCacheSplitDB, entryId int64) {
	if err := db.TouchAccess(entryId); err != nil {
		t.Fatalf("failed to touch access: %v", err)
	}
}

func deleteAndVerify(t *testing.T, db *store.UrlCacheSplitDB, testUrl string) {
	if err := db.DeleteUrlCache(testUrl); err != nil {
		t.Fatalf("failed to delete url cache: %v", err)
	}

	assertUrlExistence(t, db, testUrl, false)
}

func TestUrlCacheDB_AccessLog(t *testing.T) {
	db := createTestUrlCacheDB(t)
	defer db.Close()

	testUrl := "https://example.com"
	saveTestEntry(t, db, testUrl, []byte("test payload"), 3600)

	log := makeTestAccessLog(testUrl)
	if err := db.LogUrlAccess(log); err != nil {
		t.Fatalf("failed to log url access: %v", err)
	}
}

func makeTestAccessLog(url string) store.UrlAccessLog {
	return store.UrlAccessLog{
		UrlHash:      store.HashUrl(url),
		RawUrl:       url,
		IsCacheHit:   true,
		StatusCode:   200,
		BytesServed:  1024,
		DurationMs:   1,
		CallerSource: "test",
	}
}

func createTestUrlCacheDB(t *testing.T) *store.UrlCacheSplitDB {
	dbPath := filepath.Join(t.TempDir(), "url_cache.db")
	db, err := store.OpenUrlCacheSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test url cache db: %v", err)
	}

	return db
}

func saveTestEntry(t *testing.T, db *store.UrlCacheSplitDB, rawUrl string, payload []byte, ttlSeconds int64) {
	entry := makeTestCacheEntry(rawUrl, payload, ttlSeconds)
	if err := db.SaveUrlCache(entry); err != nil {
		t.Fatalf("failed to save test entry: %v", err)
	}
}

func makeTestCacheEntry(rawUrl string, payload []byte, ttlSeconds int64) store.UrlCacheEntry {
	now := time.Now().Unix()
	entry := initTestCachePayload(rawUrl, payload)
	entry.ExpiresAt = now + ttlSeconds
	entry.CreatedAt = now
	entry.UpdatedAt = now

	return entry
}

func initTestCachePayload(rawUrl string, payload []byte) store.UrlCacheEntry {
	return store.UrlCacheEntry{
		UrlHash:       store.HashUrl(rawUrl),
		RawUrl:        rawUrl,
		UrlText:       rawUrl,
		HttpMethod:    "GET",
		StatusCode:    200,
		ContentType:   "text/plain",
		ContentLength: int64(len(payload)),
		PayloadBlob:   payload,
		IsActive:      true,
	}
}

func assertUrlExistence(t *testing.T, db *store.UrlCacheSplitDB, rawUrl string, hasExpected bool) {
	_, hasHit, err := db.GetUrlCache(rawUrl)
	if err != nil {
		t.Fatalf("unexpected error checking url: %v", err)
	}

	if hasHit != hasExpected {
		t.Fatalf("expected existence %v for %s, got %v", hasExpected, rawUrl, hasHit)
	}
}
