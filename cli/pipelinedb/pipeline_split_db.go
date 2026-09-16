package pipelinedb

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// PipelineSplitDb encapsulates an isolated SQLite database connection for a single repository's pipeline data.
type PipelineSplitDb struct {
	conn     *sql.DB
	RepoSlug string
	Path     string
}

// PipelineSplitDB is an alias to PipelineSplitDb for backward compatibility.
type PipelineSplitDB = PipelineSplitDb

// SanitizeRepoSlug converts a repository slug into a valid safe filesystem name.
func SanitizeRepoSlug(repo string) string {
	lower := strings.ToLower(strings.TrimSpace(repo))
	slug := lazyregex.SlugSanitizeRegex.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "pipeline-default"
	}

	return slug
}

// PipelineDbDir returns the dedicated directory where pipeline split DBs live.
func PipelineDbDir() string {
	dir := filepath.Join(store.BinaryDataDir(), "pipeline_db")
	_ = os.MkdirAll(dir, 0755)

	return dir
}

// PipelineDBDir is an alias to PipelineDbDir.
var PipelineDBDir = PipelineDbDir

// RepoScopedPipelineDbDir returns the repository-scoped data directory for pipeline db.
func RepoScopedPipelineDbDir(repoRoot string) string {
	return filepath.Join(repoRoot, ".gitmap", "data")
}

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}

	return true
}

func resolveRepoScopedPath(repoRoot string) string {
	primary := filepath.Join(RepoScopedPipelineDbDir(repoRoot), "pipeline.db")
	if isFileExisting(primary) {
		return primary
	}
	legacy := filepath.Join(repoRoot, ".gitmap", "pipeline.db")
	if isFileExisting(legacy) {
		return legacy
	}

	return primary
}

func fallbackBinaryPipelineDbPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)

	return filepath.Join(PipelineDbDir(), "pipeline_"+slug+".db")
}

func isTestRepoSlug(repoSlug string) bool {
	if strings.HasPrefix(repoSlug, "test-") {
		return true
	}
	if strings.Contains(repoSlug, "/test-") {
		return true
	}

	return false
}

func isDirExisting(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func stripVersionSuffix(name string) string {
	idx := strings.LastIndex(name, "-v")
	if idx > 0 {
		return name[:idx]
	}

	return name
}

func checkCandidateDir(candidate string) string {
	if !isDirExisting(candidate) {
		return ""
	}
	root, err := gitutil.RepoRoot(candidate)
	if err == nil && root != "" {
		return root
	}

	return ""
}

func findCandidateRepoRoot(repoSlug string) string {
	base := filepath.Base(repoSlug)
	candidates := []string{base, stripVersionSuffix(base)}
	for _, cand := range candidates {
		if root := checkCandidateDir(cand); root != "" {
			return root
		}
	}

	return ""
}

func isRepoMatchingSlug(repoRoot, repoSlug string) bool {
	remote, err := gitutil.RemoteURL(repoRoot)
	if err == nil && strings.Contains(strings.ToLower(remote), strings.ToLower(repoSlug)) {
		return true
	}
	base := filepath.Base(repoRoot)
	if strings.EqualFold(base, filepath.Base(repoSlug)) {
		return true
	}
	trimmed := stripVersionSuffix(filepath.Base(repoSlug))

	return strings.EqualFold(base, trimmed)
}

func findRepoRootInStore(repoSlug string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	records, err := db.FindBySlug(repoSlug)
	if err != nil || len(records) == 0 {
		return ""
	}

	for _, rec := range records {
		if isDirExisting(rec.AbsolutePath) {
			return rec.AbsolutePath
		}
	}

	return ""
}

func resolveTargetRepoRoot(repoSlug string) string {
	root, err := gitutil.RepoRoot(".")
	if err == nil && root != "" && isRepoMatchingSlug(root, repoSlug) {
		return root
	}
	if cand := findCandidateRepoRoot(repoSlug); cand != "" {
		return cand
	}
	if storeRoot := findRepoRootInStore(repoSlug); storeRoot != "" {
		return storeRoot
	}
	if err == nil && root != "" {
		return root
	}

	return ""
}

// ResolvePipelineDbPath resolves the repository-scoped SQLite database path.
func ResolvePipelineDbPath(repoSlug string) string {
	if isTestRepoSlug(repoSlug) {
		return fallbackBinaryPipelineDbPath(repoSlug)
	}
	if root := resolveTargetRepoRoot(repoSlug); root != "" {
		return resolveRepoScopedPath(root)
	}

	return resolveRepoScopedPath(".")
}

// PipelineDbPath returns the full SQLite database file path for a repository.
func PipelineDbPath(repoSlug string) string {
	return ResolvePipelineDbPath(repoSlug)
}

// PipelineDBPath is an alias to PipelineDbPath.
var PipelineDBPath = PipelineDbPath

// OpenPipelineSplitDb opens or initializes the split SQLite database for a repo.
func OpenPipelineSplitDb(repoSlug string) (*PipelineSplitDb, error) {
	dbPath := ResolvePipelineDbPath(repoSlug)
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open pipeline split db "+repoSlug)
	}

	return initPipelineSplitConn(conn, repoSlug, dbPath)
}

func (p *PipelineSplitDb) setupSchema() error {
	if err := p.InitSchema(); err != nil {
		_ = p.conn.Close()

		return err
	}

	return nil
}

func initPipelineSplitConn(conn *sql.DB, repoSlug, dbPath string) (*PipelineSplitDb, error) {
	if err := store.ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "configure pipeline split db "+repoSlug)
	}

	p := &PipelineSplitDb{conn: conn, RepoSlug: repoSlug, Path: dbPath}

	return p, p.setupSchema()
}

// OpenPipelineSplitDB is an alias to OpenPipelineSplitDb.
var OpenPipelineSplitDB = OpenPipelineSplitDb

func pipelineSchemaQueries() []string {
	return []string{
		sqlCreatePipelineRun,
		sqlCreatePipelineErrorLog,
		sqlCreatePipelineDetailErrorLog,
		sqlCreatePipelineCompactErrorLog,
		sqlCreatePipelineSegment,
	}
}

func (p *PipelineSplitDb) executeSchemaQueries(queries []string) error {
	for _, q := range queries {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "init pipeline db schema")
		}
	}

	return nil
}

// InitSchema ensures all pipeline tables exist.
func (p *PipelineSplitDb) InitSchema() error {
	return p.executeSchemaQueries(pipelineSchemaQueries())
}

// Close closes the underlying SQLite connection.
func (p *PipelineSplitDb) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}
