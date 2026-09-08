package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	_ "modernc.org/sqlite"
)

const (
	sitesDBFileName       = "sites.db"
	sqlCreateSiteRegistry = `CREATE TABLE IF NOT EXISTS SiteRegistry (
    SiteRegistryId  INTEGER PRIMARY KEY AUTOINCREMENT,
    Domain          TEXT NOT NULL UNIQUE,
    SiteType        TEXT NOT NULL,
    DocumentRoot    TEXT NOT NULL,
    NginxConfigPath TEXT NOT NULL,
    PhpVersion      TEXT NULL,
    PhpSocketPath   TEXT NULL,
    ListenPort      INTEGER NOT NULL DEFAULT 80,
    IsSslEnabled    INTEGER NOT NULL DEFAULT 0,
    IsActive        INTEGER NOT NULL DEFAULT 1,
    Description     TEXT NULL,
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxSiteRegistry_Domain ON SiteRegistry(Domain);
CREATE INDEX IF NOT EXISTS IdxSiteRegistry_SiteType ON SiteRegistry(SiteType);
CREATE INDEX IF NOT EXISTS IdxSiteRegistry_IsActive ON SiteRegistry(IsActive);`
)

// SitesSplitDB wraps an isolated SQLite database connection for Nginx virtual host sites.
type SitesSplitDB struct {
	conn *sql.DB
	Path string
}

// SitesDBPath returns the full path to the split sites database.
func SitesDBPath() string {
	return filepath.Join(BinaryDataDir(), sitesDBFileName)
}

// OpenSitesSplitDB opens the canonical split database for Nginx virtual host sites.
func OpenSitesSplitDB() (*SitesSplitDB, error) {
	return OpenSitesSplitDBAt(SitesDBPath())
}

// OpenSitesSplitDBAt opens or creates a split sites database at a specific path.
func OpenSitesSplitDBAt(dbPath string) (*SitesSplitDB, error) {
	parentDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return nil, apperror.WrapSimple(err, "sites_split.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "sites_split.open")
	}

	conn.SetMaxOpenConns(1)

	db := &SitesSplitDB{conn: conn, Path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return db, nil
}

// InitSchema creates the SiteRegistry table and indexes if absent.
func (s *SitesSplitDB) InitSchema() error {
	if _, err := s.conn.Exec(sqlCreateSiteRegistry); err != nil {
		return apperror.WrapSimple(err, "sites_split.initSiteRegistry")
	}

	return nil
}

// Close closes the underlying split database connection.
func (s *SitesSplitDB) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

// Conn returns the raw database connection.
func (s *SitesSplitDB) Conn() *sql.DB {
	return s.conn
}
