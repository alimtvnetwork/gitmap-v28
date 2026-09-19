package cmdautomation

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

const (
	queryRuntimeByName = `SELECT id, name, binary_name, binary_path, version, status,
		install_cmd, profile_suggestion, fallback_cmd, discovered_at, last_verified_at, is_valid
		FROM runtimes WHERE name = ? LIMIT 1;`

	upsertRuntimeSQL = `INSERT INTO runtimes (
		name, binary_name, binary_path, version, status, install_cmd,
		profile_suggestion, fallback_cmd, discovered_at, last_verified_at, is_valid
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(name) DO UPDATE SET
		binary_name=excluded.binary_name, binary_path=excluded.binary_path, version=excluded.version,
		status=excluded.status, install_cmd=excluded.install_cmd, profile_suggestion=excluded.profile_suggestion,
		fallback_cmd=excluded.fallback_cmd, last_verified_at=excluded.last_verified_at, is_valid=excluded.is_valid;`
)

var automationSchemaDDL = []string{
	`CREATE TABLE IF NOT EXISTS runtimes (
		id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL UNIQUE,
		binary_name TEXT NOT NULL, binary_path TEXT NOT NULL, version TEXT NOT NULL,
		status TEXT NOT NULL, install_cmd TEXT NOT NULL, profile_suggestion TEXT NOT NULL,
		fallback_cmd TEXT NOT NULL, discovered_at DATETIME NOT NULL,
		last_verified_at DATETIME NOT NULL, is_valid INTEGER NOT NULL DEFAULT 1
	);`,
	`CREATE INDEX IF NOT EXISTS idx_runtimes_name ON runtimes(name);`,
	`CREATE INDEX IF NOT EXISTS idx_runtimes_status ON runtimes(status);`,
	`CREATE TABLE IF NOT EXISTS file_manifest (
		id INTEGER PRIMARY KEY AUTOINCREMENT, relative_path TEXT NOT NULL UNIQUE,
		file_name TEXT NOT NULL, file_extension TEXT NOT NULL, parent_folder_path TEXT NOT NULL,
		absolute_file_path TEXT NOT NULL, absolute_parent_folder_path TEXT NOT NULL,
		file_size INTEGER NOT NULL, modified_timestamp INTEGER NOT NULL,
		content_hash TEXT, is_binary INTEGER NOT NULL DEFAULT 0, scanned_at DATETIME NOT NULL
	);`,
	`CREATE INDEX IF NOT EXISTS idx_manifest_ext ON file_manifest(file_extension);`,
	`CREATE INDEX IF NOT EXISTS idx_manifest_parent ON file_manifest(parent_folder_path);`,
	`CREATE INDEX IF NOT EXISTS idx_manifest_size ON file_manifest(file_size);`,
	`CREATE TABLE IF NOT EXISTS execution_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT, runtime TEXT NOT NULL, command_type TEXT NOT NULL,
		script_preview TEXT NOT NULL, files_matched INTEGER NOT NULL, files_processed INTEGER NOT NULL,
		workers_used INTEGER NOT NULL, threads_per_worker INTEGER NOT NULL, encoding TEXT NOT NULL,
		duration_ms INTEGER NOT NULL, exit_code INTEGER NOT NULL, executed_at DATETIME NOT NULL
	);`,
	`CREATE INDEX IF NOT EXISTS idx_exec_runtime ON execution_history(runtime);`,
	`CREATE INDEX IF NOT EXISTS idx_exec_date ON execution_history(executed_at);`,
	`CREATE TABLE IF NOT EXISTS search_exclusions (
		id INTEGER PRIMARY KEY AUTOINCREMENT, pattern TEXT NOT NULL UNIQUE,
		reason TEXT NOT NULL DEFAULT 'user_excluded', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE INDEX IF NOT EXISTS idx_exclusions_pattern ON search_exclusions(pattern);`,
}

// OpenAutomationDB opens the repository automation database with WAL mode and pragmas.
func OpenAutomationDB(repoPath string) result.Result[*sql.DB] {
	dbPath := ResolveAutomationDbPath(repoPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return result.Fail[*sql.DB](apperror.WrapSimple(err, "OpenAutomationDB.open"))
	}
	if err := store.ConfigureSQLiteConn(db); err != nil {
		_ = db.Close()
		return result.Fail[*sql.DB](apperror.WrapSimple(err, "OpenAutomationDB.configure"))
	}
	if schemaErr := initAutomationSchema(db); schemaErr != nil {
		_ = db.Close()
		return result.Fail[*sql.DB](schemaErr)
	}
	return result.Ok(db)
}

// ResolveAutomationDbPath resolves the repository-scoped SQLite path.
func ResolveAutomationDbPath(repoPath string) string {
	slug := resolveRepoSlug(repoPath)
	root := resolveTargetRepoRoot(repoPath)
	dir := filepath.Join(root, ".gitmap", "data", slug, "automation")
	if !isDirExisting(filepath.Join(root, ".gitmap")) {
		dir = filepath.Join(store.BinaryDataDir(), "automation", slug)
	}
	_ = os.MkdirAll(dir, 0755)
	return filepath.ToSlash(filepath.Join(dir, "sql.db"))
}

func resolveTargetRepoRoot(repoPath string) string {
	if repoPath != "" {
		return repoPath
	}
	return findRepoRoot()
}

func resolveRepoSlug(repoPath string) string {
	if repoPath != "" {
		return sanitizeAutomationSlug(filepath.Base(repoPath))
	}
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err == nil && len(out) > 0 {
		return sanitizeAutomationSlug(strings.TrimSpace(string(out)))
	}
	return "alimtvnetwork-gitmap-v28"
}

func sanitizeAutomationSlug(repo string) string {
	lower := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(repo), ".git"))
	var sb strings.Builder
	for _, r := range lower {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			sb.WriteRune(r)
		} else if r == '/' || r == '\\' || r == '.' || r == ':' || r == '@' {
			sb.WriteRune('-')
		}
	}
	res := strings.Trim(sb.String(), "-")
	if res == "" {
		return "automation-default"
	}
	return res
}

func initAutomationSchema(db *sql.DB) *apperror.AppError {
	for _, stmt := range automationSchemaDDL {
		if _, err := db.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "initAutomationSchema")
		}
	}
	return nil
}

// GetCachedRuntime retrieves a cached runtime record by name.
func GetCachedRuntime(db *sql.DB, name string) result.Result[RuntimeRecord] {
	if db == nil {
		return result.Fail[RuntimeRecord](apperror.NewSimple("nil db", "E7102"))
	}
	row := db.QueryRow(queryRuntimeByName, name)
	return scanRuntimeRecord(row)
}

func scanRuntimeRecord(row *sql.Row) result.Result[RuntimeRecord] {
	var r RuntimeRecord
	var validInt int
	err := row.Scan(&r.ID, &r.Name, &r.BinaryName, &r.BinaryPath, &r.Version,
		&r.Status, &r.InstallCmd, &r.ProfileSuggestion, &r.FallbackCmd,
		&r.DiscoveredAt, &r.LastVerifiedAt, &validInt)
	if err != nil {
		return result.Fail[RuntimeRecord](apperror.WrapSimple(err, "scanRuntimeRecord"))
	}
	r.IsValid = (validInt == 1)
	return result.Ok(r)
}

// SaveCachedRuntime inserts or updates a runtime record in SQLite cache.
func SaveCachedRuntime(db *sql.DB, rec RuntimeRecord) *apperror.AppError {
	if db == nil {
		return apperror.NewSimple("nil db", "E7102")
	}
	rec = normalizeRuntimeTimestamps(rec)
	validInt := boolToInt(rec.IsValid)
	_, err := db.Exec(upsertRuntimeSQL,
		rec.Name, rec.BinaryName, rec.BinaryPath, rec.Version, rec.Status,
		rec.InstallCmd, rec.ProfileSuggestion, rec.FallbackCmd,
		rec.DiscoveredAt, rec.LastVerifiedAt, validInt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "SaveCachedRuntime")
	}
	return nil
}

func normalizeRuntimeTimestamps(rec RuntimeRecord) RuntimeRecord {
	now := time.Now().UTC().Format(time.RFC3339)
	if rec.DiscoveredAt == "" {
		rec.DiscoveredAt = now
	}
	rec.LastVerifiedAt = now
	return rec
}

func boolToInt(val bool) int {
	if val {
		return 1
	}
	return 0
}
