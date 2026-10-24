package cmddb

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// RestTableInfo holds table and view metadata.
type RestTableInfo struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	ColumnCount int    `json:"ColumnCount"`
	RowCount    int64  `json:"RowCount"`
	IsView      bool   `json:"IsView"`
}

// RestColumnInfo holds column schema introspection data.
type RestColumnInfo struct {
	Table        string `json:"Table"`
	Name         string `json:"Name"`
	Type         string `json:"Type"`
	IsNullable   bool   `json:"IsNullable"`
	IsPrimaryKey bool   `json:"IsPrimaryKey"`
	DefaultValue any    `json:"DefaultValue"`
}

// CachedSchema stores introspected schema entities for a database.
type CachedSchema struct {
	Tables       []RestTableInfo
	ColumnsByTab map[string][]RestColumnInfo
	LoadedAt     time.Time
}

// MultiDbManager manages thread-safe SQLite connections and schema caches.
type MultiDbManager struct {
	mu      sync.RWMutex
	aliases map[string]string
	conns   map[string]*sql.DB
	schemas map[string]*CachedSchema
}

// NewMultiDbManager constructs a MultiDbManager instance.
func NewMultiDbManager(dbs map[string]string) *MultiDbManager {
	mgr := &MultiDbManager{
		aliases: make(map[string]string),
		conns:   make(map[string]*sql.DB),
		schemas: make(map[string]*CachedSchema),
	}
	for k, v := range dbs {
		mgr.aliases[strings.ToLower(k)] = v
	}

	return mgr
}

// Mount dynamically binds a new alias to a database file path.
func (m *MultiDbManager) Mount(alias string, path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cleanAlias := strings.ToLower(strings.TrimSpace(alias))
	if cleanAlias == "" {
		return apperror.NewValidationError("alias cannot be empty")
	}
	m.aliases[cleanAlias] = filepath.Clean(path)
	delete(m.schemas, cleanAlias)

	return nil
}

// GetPath returns the file path for the given alias.
func (m *MultiDbManager) GetPath(alias string) (string, error) {
	m.mu.RLock()
	cleanAlias := strings.ToLower(strings.TrimSpace(alias))
	p, ok := m.aliases[cleanAlias]
	m.mu.RUnlock()
	if ok {
		return p, nil
	}

	resolved, err := ResolveDbPath(alias)
	if err == nil && resolved != "" {
		_ = m.Mount(cleanAlias, resolved)

		return resolved, nil
	}

	return "", apperror.NewNotFoundError(fmt.Sprintf("alias not found: %s", alias))
}

// GetDB returns an open connection and path for an alias.
func (m *MultiDbManager) GetDB(alias string) (*sql.DB, string, error) {
	path, err := m.GetPath(alias)
	if err != nil {
		return nil, "", err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	cleanAlias := strings.ToLower(strings.TrimSpace(alias))
	if conn, exists := m.conns[cleanAlias]; exists {
		return conn, path, nil
	}

	conn, openErr := m.openAndConfigureConn(path)
	if openErr != nil {
		return nil, path, openErr
	}
	m.conns[cleanAlias] = conn

	return conn, path, nil
}

func (m *MultiDbManager) openAndConfigureConn(path string) (*sql.DB, error) {
	openRes := store.OpenSQLiteDB(path)
	if openRes.IsFailure() {
		return nil, apperror.WrapSimple(openRes.AppError(), "open sqlite "+path)
	}
	conn := openRes.Value
	conn.SetMaxOpenConns(1)

	return conn, nil
}

// InvalidateCache clears the cached schema for the specified alias.
func (m *MultiDbManager) InvalidateCache(alias string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cleanAlias := strings.ToLower(strings.TrimSpace(alias))
	delete(m.schemas, cleanAlias)
}

// GetTables returns table metadata with buffer-before-next introspection.
func (m *MultiDbManager) GetTables(alias string) ([]RestTableInfo, error) {
	cleanAlias := strings.ToLower(strings.TrimSpace(alias))
	m.mu.RLock()
	cached, ok := m.schemas[cleanAlias]
	m.mu.RUnlock()
	if ok && cached != nil {
		return cached.Tables, nil
	}

	return m.loadAndCacheTables(cleanAlias)
}

func (m *MultiDbManager) loadAndCacheTables(alias string) ([]RestTableInfo, error) {
	conn, _, err := m.GetDB(alias)
	if err != nil {
		return nil, err
	}

	tables, err := introspectAllTables(conn)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.schemas[alias] = &CachedSchema{
		Tables:       tables,
		ColumnsByTab: make(map[string][]RestColumnInfo),
		LoadedAt:     time.Now(),
	}

	return tables, nil
}

func introspectAllTables(conn *sql.DB) ([]RestTableInfo, error) {
	raw, err := queryRawTableNames(conn)
	if err != nil {
		return nil, err
	}

	tables := make([]RestTableInfo, 0, len(raw))
	for _, item := range raw {
		rowCount := countTableRowSafe(conn, item[0], item[1])
		cols, _ := queryTableColumnsRaw(conn, item[0])
		tables = append(tables, RestTableInfo{
			Name:        item[0],
			Type:        item[1],
			ColumnCount: len(cols),
			RowCount:    rowCount,
			IsView:      item[1] == "view",
		})
	}

	return tables, nil
}

func queryRawTableNames(conn *sql.DB) ([][2]string, error) {
	q := `SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY type, name`
	rows, err := conn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out [][2]string
	for rows.Next() {
		var n, k string
		if scanErr := rows.Scan(&n, &k); scanErr == nil {
			out = append(out, [2]string{n, k})
		}
	}

	return out, rows.Err()
}

func countTableRowSafe(conn *sql.DB, tableName, kind string) int64 {
	if kind == "view" {
		return 0
	}
	var count int64
	q := fmt.Sprintf(`SELECT COUNT(*) FROM "%s"`, strings.ReplaceAll(tableName, `"`, `""`))
	if err := conn.QueryRow(q).Scan(&count); err != nil {
		return 0
	}
	return count
}

// GetColumns returns column schema metadata with buffer-before-next reading.
func (m *MultiDbManager) GetColumns(alias, tableName string) ([]RestColumnInfo, error) {
	conn, _, err := m.GetDB(alias)
	if err != nil {
		return nil, err
	}

	return queryTableColumnsRaw(conn, tableName)
}

func queryTableColumnsRaw(conn *sql.DB, tableName string) ([]RestColumnInfo, error) {
	q := fmt.Sprintf(`PRAGMA table_info("%s")`, strings.ReplaceAll(tableName, `"`, `""`))
	rows, err := conn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPragmaColumns(rows, tableName)
}

func scanPragmaColumns(rows *sql.Rows, table string) ([]RestColumnInfo, error) {
	var cols []RestColumnInfo
	for rows.Next() {
		var cid, notNull, pk int
		var name, colType string
		var dflt any
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			return nil, err
		}
		cols = append(cols, RestColumnInfo{
			Table:        table,
			Name:         name,
			Type:         colType,
			IsNullable:   notNull == 0,
			IsPrimaryKey: pk > 0,
			DefaultValue: dflt,
		})
	}

	return cols, rows.Err()
}

// ListAliases returns a copy of all mounted database aliases.
func (m *MultiDbManager) ListAliases() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	aliases := make([]string, 0, len(m.aliases))
	for a := range m.aliases {
		aliases = append(aliases, a)
	}

	return aliases
}

// Close closes all managed database connections.
func (m *MultiDbManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, conn := range m.conns {
		if conn != nil {
			_ = conn.Close()
		}
	}
	m.conns = make(map[string]*sql.DB)

	return nil
}
