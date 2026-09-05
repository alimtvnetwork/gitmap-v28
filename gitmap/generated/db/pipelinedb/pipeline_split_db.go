package pipelinedb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
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

// PipelineSplitDbFieldType represents column name enums for PipelineSplitDb.
type PipelineSplitDbFieldType string

// Name returns the identifier name of the field enum.
func (e PipelineSplitDbFieldType) Name() string {
	return string(e)
}

// String returns the string representation of the field enum.
func (e PipelineSplitDbFieldType) String() string {
	return string(e)
}

// Value returns the raw string value of the field enum.
func (e PipelineSplitDbFieldType) Value() string {
	return string(e)
}

// IsCompare checks equality against another field enum object.
func (e PipelineSplitDbFieldType) IsCompare(target PipelineSplitDbFieldType) bool {
	return e == target
}

// IsEnum checks whether this field enum exists in the valid enum map.
func (e PipelineSplitDbFieldType) IsEnum() bool {
	return pipelineSplitDbValidMap[e]
}

// MarshalJSON implements json.Marshaler.
func (e PipelineSplitDbFieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// UnmarshalJSON implements json.Unmarshaler with strict map validation.
func (e *PipelineSplitDbFieldType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	target := PipelineSplitDbFieldType(s)
	if !pipelineSplitDbValidMap[target] {
		return fmt.Errorf("invalid %s enum: %s", "PipelineSplitDbFieldType", s)
	}
	*e = target
	return nil
}

// ToJSON converts the field enum to a JSON string representation, returning an AppError on failure.
func (e PipelineSplitDbFieldType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(e))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize field to json")
	}
	return string(b), nil
}

// FromJSON parses a field enum from a JSON string representation, returning an AppError on failure.
func (e *PipelineSplitDbFieldType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize field from json")
	}
	target := PipelineSplitDbFieldType(str)
	if !pipelineSplitDbValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid %s enum: %s", "PipelineSplitDbFieldType", str), "validate field enum from json")
	}
	*e = target
	return nil
}

// IsRepoSlug checks whether this field enum instance is RepoSlug.
func (e PipelineSplitDbFieldType) IsRepoSlug() bool {
	return e == PipelineSplitDbDb.RepoSlug
}

// IsPath checks whether this field enum instance is Path.
func (e PipelineSplitDbFieldType) IsPath() bool {
	return e == PipelineSplitDbDb.Path
}

type pipelineSplitDbDbRegistry struct {
	RepoSlug PipelineSplitDbFieldType
	Path PipelineSplitDbFieldType
}

// All returns a slice of all field enums in PipelineSplitDb.
func (r pipelineSplitDbDbRegistry) All() []PipelineSplitDbFieldType {
	return []PipelineSplitDbFieldType{
		r.RepoSlug,
		r.Path,
	}
}

// Names returns a slice of string names for all fields in PipelineSplitDb.
func (r pipelineSplitDbDbRegistry) Names() []string {
	return []string{
		"RepoSlug",
		"Path",
	}
}

// IsEnum checks whether the target object matches any registered field enum in PipelineSplitDb.
func (r pipelineSplitDbDbRegistry) IsEnum(target PipelineSplitDbFieldType) bool {
	return pipelineSplitDbValidMap[target]
}

// IsRepoSlug checks whether the target object is RepoSlug.
func (r pipelineSplitDbDbRegistry) IsRepoSlug(target PipelineSplitDbFieldType) bool {
	return target == r.RepoSlug
}

// IsPath checks whether the target object is Path.
func (r pipelineSplitDbDbRegistry) IsPath(target PipelineSplitDbFieldType) bool {
	return target == r.Path
}

// ToJSON converts the field registry to a JSON string representation, returning an AppError on failure.
func (r pipelineSplitDbDbRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize registry to json")
	}
	return string(b), nil
}

// PipelineSplitDbDb provides scoped access to field enums: PipelineSplitDbDb.<Field>.
var PipelineSplitDbDb = pipelineSplitDbDbRegistry{
	RepoSlug: "RepoSlug",
	Path: "Path",
}

// pipelineSplitDbValidMap provides O(1) map validation for field enums.
var pipelineSplitDbValidMap = map[PipelineSplitDbFieldType]bool{
	PipelineSplitDbDb.RepoSlug: true,
	PipelineSplitDbDb.Path: true,
}

// PipelineSplitDbField is an alias to PipelineSplitDbDb.
var PipelineSplitDbField = PipelineSplitDbDb

// ScanPipelineSplitDb maps a database row scanner to a PipelineSplitDb entity.
func ScanPipelineSplitDb(row dbengine.RowScanner) (*PipelineSplitDb, error) {
	var item PipelineSplitDb
	var (
		raw_RepoSlug any
		raw_Path any
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
	repo := dbengine.NewRepository[PipelineSplitDb, PipelineSplitDbFieldType](
		db,
		PipelineSplitDbTable,
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
	return r.repo.Query().Select(PipelineSplitDbDb.All()...)
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
