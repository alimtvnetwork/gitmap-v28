package cmdautomation

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// ExclusionEntry represents a persisted search exclusion rule.
type ExclusionEntry struct {
	ID        int64  `json:"id"`
	Pattern   string `json:"pattern"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"createdAt"`
}

var defaultAllowedLargeFiles = map[string]bool{
	"src/data/specTree.json":              true,
	"slides-app/dist.zip":                 true,
	"docs/demo.gif":                       true,
	".ai-memory/test-inventory.json":      true,
	".ai-memory/cicd/test-inventory.json": true,
}

var defaultBinaryExts = map[string]bool{
	".exe": true, ".dll": true, ".bin": true, ".zip": true, ".tar": true,
	".gz": true, ".tgz": true, ".7z": true, ".rar": true, ".iso": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true,
	".webp": true, ".pdf": true, ".wasm": true, ".db": true, ".sqlite": true,
	".sqlite3": true, ".pyc": true, ".class": true, ".o": true, ".a": true,
	".so": true, ".dylib": true,
}

// GetAutomationDbPath returns the standardized SQLite path for automation data.
func GetAutomationDbPath() string {
	root := findRepoRoot()
	slug := resolveRepoSlug(root)
	return store.ResolveSplitDbPath(store.SectionAutomation, slug, root)
}

func findRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if isDirExisting(filepath.Join(dir, ".gitmap")) || isDirExisting(filepath.Join(dir, ".git")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

func isDirExisting(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

// OpenAutomationDb opens the repository automation database with GitMap conventions.
func OpenAutomationDb() (*sql.DB, error) {
	dbPath := GetAutomationDbPath()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return initAutomationPragmas(db)
}

func initAutomationPragmas(db *sql.DB) (*sql.DB, error) {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA foreign_keys = ON;",
		`CREATE TABLE IF NOT EXISTS search_exclusions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pattern TEXT NOT NULL UNIQUE,
			reason TEXT NOT NULL DEFAULT 'user_excluded',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			_ = db.Close()
			return nil, err
		}
	}
	return db, nil
}

// IsBinaryExtension checks whether a file extension matches common binaries.
func IsBinaryExtension(ext string) bool {
	return defaultBinaryExts[strings.ToLower(ext)]
}

// IsAllowedLargeWaiver checks whether a relative path has a waiver.
func IsAllowedLargeWaiver(relPath string) bool {
	norm := filepath.ToSlash(strings.TrimPrefix(relPath, "./"))
	if defaultAllowedLargeFiles[norm] {
		return true
	}
	for allowed := range defaultAllowedLargeFiles {
		if strings.HasSuffix(norm, allowed) {
			return true
		}
	}
	return false
}

// IsPathExcluded checks whether a file path is covered by any exclusion.
func IsPathExcluded(relPath string, exclusions []string) bool {
	norm := filepath.ToSlash(strings.TrimPrefix(relPath, "./"))
	for _, exc := range exclusions {
		clean := filepath.ToSlash(strings.TrimSpace(exc))
		if clean == "" {
			continue
		}
		if norm == clean || strings.HasSuffix(norm, "/"+clean) || strings.Contains(norm, clean) {
			return true
		}
	}
	return false
}

// ListExclusions retrieves all persisted search exclusions from SQLite.
func ListExclusions() ([]ExclusionEntry, *apperror.AppError) {
	db, err := OpenAutomationDb()
	if err != nil {
		return nil, apperror.New("list_exclusions", "E_DB_OPEN", map[string]any{"err": err.Error()})
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, pattern, reason, created_at FROM search_exclusions ORDER BY pattern ASC")
	if err != nil {
		return nil, apperror.New("list_exclusions", "E_QUERY_FAILED", map[string]any{"err": err.Error()})
	}
	defer rows.Close()

	return scanExclusionRows(rows)
}

func scanExclusionRows(rows *sql.Rows) ([]ExclusionEntry, *apperror.AppError) {
	var entries []ExclusionEntry
	for rows.Next() {
		var e ExclusionEntry
		if err := rows.Scan(&e.ID, &e.Pattern, &e.Reason, &e.CreatedAt); err != nil {
			return nil, apperror.New("scan_exclusions", "E_SCAN_FAILED", map[string]any{"err": err.Error()})
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// AddExclusion persists a new pattern into the exclusions database.
func AddExclusion(pattern string, reason string) *apperror.AppError {
	clean := filepath.ToSlash(strings.TrimSpace(pattern))
	if clean == "" {
		return apperror.NewValidationError("exclusion pattern cannot be empty")
	}
	db, err := OpenAutomationDb()
	if err != nil {
		return apperror.New("add_exclusion", "E_DB_OPEN", map[string]any{"err": err.Error()})
	}
	defer db.Close()

	q := `INSERT INTO search_exclusions (pattern, reason) VALUES (?, ?)
		  ON CONFLICT(pattern) DO UPDATE SET reason=excluded.reason`
	if _, execErr := db.Exec(q, clean, reason); execErr != nil {
		return apperror.New("add_exclusion", "E_EXEC_FAILED", map[string]any{"err": execErr.Error()})
	}
	return nil
}

// RemoveExclusion deletes an exclusion pattern from the database.
func RemoveExclusion(pattern string) *apperror.AppError {
	clean := filepath.ToSlash(strings.TrimSpace(pattern))
	db, err := OpenAutomationDb()
	if err != nil {
		return apperror.New("remove_exclusion", "E_DB_OPEN", map[string]any{"err": err.Error()})
	}
	defer db.Close()

	if _, execErr := db.Exec("DELETE FROM search_exclusions WHERE pattern = ?", clean); execErr != nil {
		return apperror.New("remove_exclusion", "E_EXEC_FAILED", map[string]any{"err": execErr.Error()})
	}
	return nil
}

// ClearExclusions purges all persisted exclusions from the database.
func ClearExclusions() *apperror.AppError {
	db, err := OpenAutomationDb()
	if err != nil {
		return apperror.New("clear_exclusions", "E_DB_OPEN", map[string]any{"err": err.Error()})
	}
	defer db.Close()

	if _, execErr := db.Exec("DELETE FROM search_exclusions"); execErr != nil {
		return apperror.New("clear_exclusions", "E_EXEC_FAILED", map[string]any{"err": execErr.Error()})
	}
	return nil
}
