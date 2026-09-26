package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	TemplatesDBFileName = "gitmap-templates.db"

	sqlCreateTemplateCategory = `CREATE TABLE IF NOT EXISTS TemplateCategory (
    CategoryId  INTEGER PRIMARY KEY AUTOINCREMENT,
    Slug        TEXT NOT NULL UNIQUE,
    Name        TEXT NOT NULL,
    ParentSlug  TEXT NOT NULL DEFAULT '',
    Description TEXT NOT NULL DEFAULT '',
    IsDefault   INTEGER NOT NULL DEFAULT 0,
    CreatedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	sqlCreateTemplateItem = `CREATE TABLE IF NOT EXISTS TemplateItem (
    ItemId          TEXT PRIMARY KEY,
    CategorySlug    TEXT NOT NULL,
    SubCategorySlug TEXT NOT NULL DEFAULT '',
    Slug            TEXT NOT NULL UNIQUE,
    Title           TEXT NOT NULL DEFAULT '',
    Text            TEXT NOT NULL,
    AdditionalJson  TEXT NOT NULL DEFAULT '{}',
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxTemplateItem_CategorySlug ON TemplateItem(CategorySlug);`

	sqlCreateTemplateVariable = `CREATE TABLE IF NOT EXISTS TemplateVariable (
    VarKey      TEXT NOT NULL,
    Scope       TEXT NOT NULL DEFAULT 'global',
    VarValue    TEXT NOT NULL,
    Description TEXT NOT NULL DEFAULT '',
    UpdatedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (VarKey, Scope)
);`

	sqlCreateTemplateImportHistory = `CREATE TABLE IF NOT EXISTS TemplateImportHistory (
    ExportHashId TEXT PRIMARY KEY,
    SourcePath   TEXT NOT NULL DEFAULT '',
    ItemCount    INTEGER NOT NULL DEFAULT 0,
    VarCount     INTEGER NOT NULL DEFAULT 0,
    ImportedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	sqlSeedDefaultCategories = `INSERT OR IGNORE INTO TemplateCategory (Slug, Name, ParentSlug, Description, IsDefault) VALUES
    ('seo', 'SEO', '', 'SEO and commit reasoning templates', 1),
    ('prompts', 'Prompts', '', 'AI and developer prompt templates', 1),
    ('ui-ux', 'UI/UX', 'prompts', 'UI and UX design guidance templates', 1),
    ('prefix', 'Prefix', '', 'Commit and message prefix templates', 1),
    ('pr-descriptions', 'PR Descriptions', '', 'Pull request description templates', 1);`

	sqlSeedDefaultTemplates = `INSERT OR IGNORE INTO TemplateItem (ItemId, CategorySlug, SubCategorySlug, Slug, Title, Text, AdditionalJson) VALUES
    ('tpl-prompt-ui-ux-audit', 'prompts', 'ui-ux', 'ui-ux-responsive-audit', '# How should UI/UX components be structured?', 'Because responsive design tokens, WCAG AA contrast ratios, and keyboard navigation states eliminate layout shift and accessibility regressions.', '{"version":"1.0","scope":"frontend"}'),
    ('tpl-prefix-standard', 'prefix', '', 'standard-commit-prefix', '# Why enforce structured commit prefixes?', 'Because deterministic conventional commit prefixes accelerate changelog generation and semantic release automation.', '{"version":"1.0","scope":"git"}');`
)

// TemplatesSplitDB manages the dedicated SQLite split database for state templates and variables.
type TemplatesSplitDB struct {
	conn *sql.DB
	Path string
}

// TemplatesDBPath returns the canonical file path to gitmap-templates.db in the binary data folder.
func TemplatesDBPath() string {
	if override := strings.TrimSpace(os.Getenv("GITMAP_TEMPLATES_DB_PATH")); override != "" {
		return override
	}

	dir := BinaryDataDir()
	_ = os.MkdirAll(dir, 0755)

	return filepath.Join(dir, TemplatesDBFileName)
}

// OpenTemplatesSplitDB opens the canonical split database for templates and variables.
func OpenTemplatesSplitDB() (*TemplatesSplitDB, error) {
	return OpenTemplatesSplitDBAt(TemplatesDBPath())
}

// OpenTemplatesSplitDBAt opens or creates a templates split database at a specific path.
func OpenTemplatesSplitDBAt(dbPath string) (*TemplatesSplitDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "templates_split.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "templates_split.open")
	}

	return initTemplatesSplitConn(conn, dbPath)
}

func initTemplatesSplitConn(conn *sql.DB, dbPath string) (*TemplatesSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "templates_split.config")
	}

	db := &TemplatesSplitDB{conn: conn, Path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return db, nil
}

// InitSchema creates the 4 template state tables and seeds default categories and templates.
func (s *TemplatesSplitDB) InitSchema() error {
	stmts := []string{
		sqlCreateTemplateCategory,
		sqlCreateTemplateItem,
		sqlCreateTemplateVariable,
		sqlCreateTemplateImportHistory,
		sqlSeedDefaultCategories,
		sqlSeedDefaultTemplates,
	}
	for _, stmt := range stmts {
		if _, err := s.conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "templates_split.init_schema")
		}
	}

	return nil
}

// Close closes the underlying SQLite database connection.
func (s *TemplatesSplitDB) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

// Conn returns the raw database connection handle.
func (s *TemplatesSplitDB) Conn() *sql.DB {
	return s.conn
}
