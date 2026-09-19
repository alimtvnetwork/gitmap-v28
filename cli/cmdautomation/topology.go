package cmdautomation

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	extToLanguage = map[string]string{
		".go": "go", ".ts": "typescript", ".tsx": "typescript",
		".js": "javascript", ".jsx": "javascript", ".py": "python",
		".rs": "rust", ".cs": "csharp", ".php": "php",
		".sql": "sql", ".md": "markdown", ".sh": "shell", ".ps1": "powershell",
	}

	manifestPatterns = map[string]string{
		"go.mod": "go", "go.sum": "go",
		"package.json": "typescript", "tsconfig.json": "typescript",
		"requirements.txt": "python", "pyproject.toml": "python", "setup.py": "python",
		"cargo.toml": "rust", "cargo.lock": "rust",
		"composer.json": "php",
	}

	subsystemDirHints = map[string]string{
		"server": "backend", "backend": "backend", "api": "backend", "internal": "backend",
		"db": "database", "database": "database", "sql": "database", "migrations": "database",
		"store": "database", "pipelinedb": "database", "repodb": "database",
		"web": "frontend", "ui": "frontend", "frontend": "frontend",
		".github": "cicd", ".circleci": "cicd", "ci": "cicd",
		"docs": "docs", "spec": "docs", "02-spec": "docs",
		"cli": "cli", "cmd": "cli", "commands": "cli",
		"tests": "tests", "test": "tests",
	}
)

// RunTopology executes codebase topology discovery and caching.
func RunTopology(opts TopologyOptions) TopologyResultMonad {
	start := time.Now()
	root := opts.Dir
	if root == "" {
		root = "."
	}
	cached, hasValid := getCachedIfApplicable(root, opts.IsRefresh)
	if hasValid {
		return result.Ok(cached)
	}
	res := buildTopology(root, opts.TtlSec, start)
	_ = persistTopologyCache(res)
	return result.Ok(res)
}

func getCachedIfApplicable(root string, isRefresh bool) (TopologyResult, bool) {
	if isRefresh {
		return TopologyResult{}, false
	}
	return loadCachedTopology(root)
}

func loadCachedTopology(root string) (TopologyResult, bool) {
	db, err := OpenAutomationDb()
	if err != nil {
		return TopologyResult{}, false
	}
	defer db.Close()
	return queryTopologyCache(db, root)
}

func queryTopologyCache(db *sql.DB, root string) (TopologyResult, bool) {
	if appErr := initTopologyCacheTable(db); appErr != nil {
		return TopologyResult{}, false
	}
	query := `SELECT data_json, expires_at FROM codebase_topology_cache WHERE root_path = ?;`
	row := db.QueryRow(query, filepath.ToSlash(root))
	var dataJson, expiresAtStr string
	if err := row.Scan(&dataJson, &expiresAtStr); err != nil {
		return TopologyResult{}, false
	}
	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil || time.Now().After(expiresAt) {
		return TopologyResult{}, false
	}
	var res TopologyResult
	if err := json.Unmarshal([]byte(dataJson), &res); err != nil {
		return TopologyResult{}, false
	}
	res.IsValid = true
	return res, true
}

func initTopologyCacheTable(db *sql.DB) *apperror.AppError {
	schema := `CREATE TABLE IF NOT EXISTS codebase_topology_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		root_path TEXT NOT NULL UNIQUE,
		data_json TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return apperror.WrapSimple(err, "init topology cache table")
	}

	return nil
}

func persistTopologyCache(res TopologyResult) *apperror.AppError {
	db, err := OpenAutomationDb()
	if err != nil {
		return apperror.WrapSimple(err, "open automation db for topology cache")
	}
	defer db.Close()
	if appErr := initTopologyCacheTable(db); appErr != nil {
		return appErr
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		return apperror.WrapSimple(err, "serialize topology cache")
	}
	upsert := `INSERT INTO codebase_topology_cache (root_path, data_json, expires_at)
		VALUES (?, ?, ?)
		ON CONFLICT(root_path) DO UPDATE SET
			data_json = excluded.data_json,
			expires_at = excluded.expires_at,
			created_at = CURRENT_TIMESTAMP;`
	_, execErr := db.Exec(upsert, res.RootPath, string(bytes), res.ExpiresAt)
	if execErr != nil {
		return apperror.WrapSimple(execErr, "upsert topology cache")
	}
	return nil
}

func buildTopology(root string, ttlSec int, start time.Time) TopologyResult {
	if ttlSec <= 0 {
		ttlSec = 1800
	}
	now := time.Now().UTC()
	expires := now.Add(time.Duration(ttlSec) * time.Second)
	res := newEmptyTopologyResult(root, ttlSec, now, expires)
	walkCodebase(root, &res)
	res.Duration = time.Since(start)
	res.DurationMs = float64(res.Duration.Microseconds()) / 1000.0
	res.IsValid = true
	return res
}

func newEmptyTopologyResult(root string, ttlSec int, now, expires time.Time) TopologyResult {
	subsystems := map[string]SubsystemData{
		"backend": {}, "database": {}, "frontend": {},
		"cicd": {}, "docs": {}, "cli": {}, "tests": {},
	}
	return TopologyResult{
		Version:     "1.0.0",
		GeneratedAt: now.Format(time.RFC3339),
		ExpiresAt:   expires.Format(time.RFC3339),
		TtlSeconds:  ttlSec,
		RootPath:    filepath.ToSlash(root),
		Manifests:   make(map[string][]string),
		Languages:   make(map[string]int),
		LangRoots:   make(map[string][]string),
		Subsystems:  subsystems,
	}
}

func walkCodebase(root string, res *TopologyResult) {
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel := filepath.ToSlash(path)
		if info.IsDir() {
			return handleDirWalk(info.Name(), rel, res)
		}
		handleFileWalk(info.Name(), rel, res)
		return nil
	})
}

func handleDirWalk(dirName, relPath string, res *TopologyResult) error {
	if isTopologyExcludedDir(dirName) {
		return filepath.SkipDir
	}
	subsys, hasHint := subsystemDirHints[strings.ToLower(dirName)]
	if hasHint {
		data := res.Subsystems[subsys]
		data.Roots = appendUnique(data.Roots, relPath)
		res.Subsystems[subsys] = data
	}
	return nil
}

func isTopologyExcludedDir(name string) bool {
	return name == ".git" || name == ".gitmap" || name == "node_modules" ||
		name == "dist" || name == "build" || name == "bin" ||
		name == ".gemini" || name == ".pytest_cache" || name == "__pycache__"
}

func handleFileWalk(fileName, relPath string, res *TopologyResult) {
	res.TotalFiles++
	ext := strings.ToLower(filepath.Ext(fileName))
	lang, hasLang := extToLanguage[ext]
	if hasLang {
		res.Languages[lang]++
		dir := filepath.ToSlash(filepath.Dir(relPath))
		res.LangRoots[lang] = appendUnique(res.LangRoots[lang], dir)
	}
	inspectFileSubsystems(fileName, relPath, ext, res)
}

func inspectFileSubsystems(fileName, relPath, ext string, res *TopologyResult) {
	lowerName := strings.ToLower(fileName)
	checkManifestFile(lowerName, relPath, res)
	checkDatabaseFile(lowerName, relPath, ext, res)
	checkEntrypointFile(fileName, relPath, res)
	checkWorkflowFile(relPath, ext, res)
	checkTestFile(lowerName, relPath, res)
}

func checkManifestFile(lowerName, relPath string, res *TopologyResult) {
	lang, hasPattern := manifestPatterns[lowerName]
	if hasPattern {
		res.Manifests[lang] = appendUnique(res.Manifests[lang], relPath)
	}
}

func checkDatabaseFile(lowerName, relPath, ext string, res *TopologyResult) {
	isSql := ext == ".sql" || strings.Contains(lowerName, "schema") || strings.Contains(lowerName, "migration")
	if isSql {
		data := res.Subsystems["database"]
		data.SchemaFiles = appendUnique(data.SchemaFiles, relPath)
		res.Subsystems["database"] = data
	}
}

func checkEntrypointFile(fileName, relPath string, res *TopologyResult) {
	isEntry := fileName == "main.go" || fileName == "index.ts" || fileName == "app.py" || fileName == "main.rs"
	if isEntry {
		data := res.Subsystems["backend"]
		data.Entrypoints = appendUnique(data.Entrypoints, relPath)
		res.Subsystems["backend"] = data
	}
}

func checkWorkflowFile(relPath, ext string, res *TopologyResult) {
	isWorkflow := (ext == ".yml" || ext == ".yaml") && strings.Contains(relPath, ".github")
	if isWorkflow {
		data := res.Subsystems["cicd"]
		data.Workflows = appendUnique(data.Workflows, relPath)
		res.Subsystems["cicd"] = data
	}
}

func checkTestFile(lowerName, relPath string, res *TopologyResult) {
	isTest := strings.HasSuffix(lowerName, "_test.go") || strings.HasSuffix(lowerName, ".test.ts") || strings.HasPrefix(lowerName, "test_")
	if isTest {
		data := res.Subsystems["tests"]
		data.TestRunners = appendUnique(data.TestRunners, relPath)
		res.Subsystems["tests"] = data
	}
}

func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}
