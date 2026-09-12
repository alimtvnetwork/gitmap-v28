package pipelinedb

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
	reg := regexp.MustCompile(`[^a-z0-9_-]+`)
	slug := reg.ReplaceAllString(lower, "-")
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

// PipelineDbPath returns the full SQLite database file path for a repository.
func PipelineDbPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)

	return filepath.Join(PipelineDbDir(), "pipeline_"+slug+".db")
}

// PipelineDBPath is an alias to PipelineDbPath.
var PipelineDBPath = PipelineDbPath

// OpenPipelineSplitDb opens or initializes the split SQLite database for a repo.
func OpenPipelineSplitDb(repoSlug string) (*PipelineSplitDb, error) {
	dbPath := PipelineDbPath(repoSlug)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open pipeline split db "+repoSlug)
	}

	return initPipelineSplitConn(conn, repoSlug, dbPath)
}

func initPipelineSplitConn(conn *sql.DB, repoSlug, dbPath string) (*PipelineSplitDb, error) {
	if err := store.ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "configure pipeline split db "+repoSlug)
	}

	p := &PipelineSplitDb{conn: conn, RepoSlug: repoSlug, Path: dbPath}
	if err := p.InitSchema(); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return p, nil
}

// OpenPipelineSplitDB is an alias to OpenPipelineSplitDb.
var OpenPipelineSplitDB = OpenPipelineSplitDb

// InitSchema ensures all pipeline tables exist.
func (p *PipelineSplitDb) InitSchema() error {
	queries := []string{
		sqlCreatePipelineRun,
		sqlCreatePipelineErrorLog,
		sqlCreatePipelineSegment,
	}

	for _, q := range queries {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "init pipeline db schema")
		}
	}

	return nil
}

// Close closes the underlying SQLite connection.
func (p *PipelineSplitDb) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}
