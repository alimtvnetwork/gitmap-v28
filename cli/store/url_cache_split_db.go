package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// UrlCacheSplitDB manages the dedicated split database for cached web responses.
type UrlCacheSplitDB struct {
	conn  *sql.DB
	group urlSingleflight
	Path  string
}

// UrlCacheDB is a type alias for UrlCacheSplitDB.
type UrlCacheDB = UrlCacheSplitDB

type urlCall struct {
	wg  sync.WaitGroup
	val any
	err error
}

type urlSingleflight struct {
	mu    sync.Mutex
	calls map[string]*urlCall
}

func (g *urlSingleflight) Do(key string, fn func() (any, error)) (any, error) {
	c, isNew := g.registerCall(key)
	if !isNew {
		c.wg.Wait()

		return c.val, c.err
	}

	return g.executeCall(key, c, fn)
}

func (g *urlSingleflight) registerCall(key string) (*urlCall, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.ensureCallsMap()

	return g.getOrCreateCall(key)
}

func (g *urlSingleflight) ensureCallsMap() {
	if g.calls == nil {
		g.calls = make(map[string]*urlCall)
	}
}

func (g *urlSingleflight) getOrCreateCall(key string) (*urlCall, bool) {
	if c, hasExisting := g.calls[key]; hasExisting {
		return c, false
	}

	c := new(urlCall)
	c.wg.Add(1)
	g.calls[key] = c

	return c, true
}

func (g *urlSingleflight) executeCall(key string, c *urlCall, fn func() (any, error)) (any, error) {
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()

	return c.val, c.err
}

// HashUrl generates a deterministic SHA-256 hex string for a given URL.
func HashUrl(rawUrl string) string {
	sum := sha256.Sum256([]byte(rawUrl))

	return hex.EncodeToString(sum[:])
}

// OpenUrlCacheDB creates or opens the dedicated url_cache.db split database.
func OpenUrlCacheDB(repoRoot string) (*UrlCacheSplitDB, error) {
	return OpenUrlCacheSplitDB(repoRoot)
}

// OpenUrlCacheSplitDB creates or opens the dedicated url_cache.db split database.
func OpenUrlCacheSplitDB(repoRoot string) (*UrlCacheSplitDB, error) {
	dbPath := ResolveUrlCacheDbPath(repoRoot)

	return OpenUrlCacheSplitDBAt(dbPath)
}

// OpenUrlCacheSplitDBAt opens or creates the URL cache database at an explicit file path.
func OpenUrlCacheSplitDBAt(dbPath string) (*UrlCacheSplitDB, error) {
	conn, err := openConfiguredCacheConn(dbPath)
	if err != nil {
		return nil, err
	}

	db := &UrlCacheSplitDB{conn: conn, Path: dbPath}
	db.registerWithRegistry()

	return db, nil
}

func openConfiguredCacheConn(dbPath string) (*sql.DB, error) {
	conn, err := openUrlCacheConn(dbPath)
	if err != nil {
		return nil, err
	}

	if err := setupUrlCacheDB(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return conn, nil
}

func openUrlCacheConn(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenUrlCacheSplitDBAt.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenUrlCacheSplitDBAt.open")
	}

	conn.SetMaxOpenConns(1)

	return conn, nil
}

func setupUrlCacheDB(conn *sql.DB) error {
	if err := initUrlCachePragmas(conn); err != nil {
		return err
	}

	return initUrlCacheSchema(conn)
}

func initUrlCachePragmas(conn *sql.DB) error {
	for _, pragma := range urlCachePragmas {
		if _, err := conn.Exec(pragma); err != nil {
			return apperror.WrapSimple(err, "store.initUrlCachePragmas")
		}
	}

	return nil
}

func initUrlCacheSchema(conn *sql.DB) error {
	for _, ddl := range urlCacheSchemaDDLs {
		if _, err := conn.Exec(ddl); err != nil {
			return apperror.WrapSimple(err, "store.initUrlCacheSchema")
		}
	}

	return nil
}

func (db *UrlCacheSplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}

	defer master.Close()

	_ = master.RegisterSplitDB(makeRegistryEntry(db.Path))
}

func makeRegistryEntry(dbPath string) SplitDatabaseEntry {
	return SplitDatabaseEntry{
		DatabaseType:  "repodb",
		DatabaseKey:   "url_cache",
		DatabasePath:  dbPath,
		Status:        "active",
		IsActive:      true,
		Description:   "HTTP/HTTPS URL payload cache and access log",
		SchemaVersion: 1,
	}
}

// Close terminates the database connection.
func (db *UrlCacheSplitDB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}

	return db.conn.Close()
}
