package pipelinedb

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/generated/db/pipelinedb/enums"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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

// ScanPipelineSplitDb maps a database row scanner to a PipelineSplitDb entity.
func ScanPipelineSplitDb(row dbengine.RowScanner) (*PipelineSplitDb, error) {
	var item PipelineSplitDb
	var (
		raw_RepoSlug any
		raw_Path     any
	)
	err := row.Scan(
		&raw_RepoSlug,
		&raw_Path,
	)
	if err != nil {
		return nil, err
	}

	item.RepoSlug = dbengine.ScanString(raw_RepoSlug)
	item.Path = dbengine.ScanString(raw_Path)
	return &item, nil
}

// PipelineSplitDbDbRepo provides typed database repository access for PipelineSplitDb.
type PipelineSplitDbDbRepo struct {
	db   *dbengine.DbWrapper
	repo *PipelineSplitDbRepository
}

// NewPipelineSplitDbDbRepo initializes a typed repository for PipelineSplitDb.
func NewPipelineSplitDbDbRepo(db *dbengine.DbWrapper) *PipelineSplitDbDbRepo {
	repo := dbengine.NewRepository[PipelineSplitDb, enums.PipelineSplitDbFieldType](
		db,
		enums.PipelineSplitDbTable,
		ScanPipelineSplitDb,
	)
	return &PipelineSplitDbDbRepo{
		db:   db,
		repo: repo,
	}
}

// Db returns the underlying DbWrapper.
func (r *PipelineSplitDbDbRepo) Db() *dbengine.DbWrapper {
	return r.db
}

// Repo returns the underlying generic Repository.
func (r *PipelineSplitDbDbRepo) Repo() *PipelineSplitDbRepository {
	return r.repo
}

// Query returns a fluent QueryBuilder initialized with all standard fields projected.
func (r *PipelineSplitDbDbRepo) Query() *PipelineSplitDbQueryBuilder {
	return r.repo.Query().Select(enums.PipelineSplitDbDb.All()...)
}

// QueryBare returns a fluent QueryBuilder without any pre-selected fields.
func (r *PipelineSplitDbDbRepo) QueryBare() *PipelineSplitDbQueryBuilder {
	return r.repo.Query()
}

// FindAll executes the query selecting all fields and returns a ListResult envelope.
func (r *PipelineSplitDbDbRepo) FindAll(ctx context.Context) dbengine.ListResult[PipelineSplitDb] {
	return r.Query().FindAll(ctx)
}

// First executes the query selecting all fields and returns the first record in an EntityResult envelope.
func (r *PipelineSplitDbDbRepo) First(ctx context.Context) dbengine.EntityResult[PipelineSplitDb] {
	return r.Query().First(ctx)
}

// Count returns the total number of records matching the query.
func (r *PipelineSplitDbDbRepo) Count(ctx context.Context) dbengine.Int64Result {
	return r.Query().Count(ctx)
}
