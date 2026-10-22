// Package store manages repository database queries and persistence.
package store

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RepoRepository encapsulates repository database queries returning typed Result monads.
type RepoRepository struct {
	conn sqlExecutor
}

// DefaultRepoRepository aliases RepoRepository for unified repository naming.
type DefaultRepoRepository = RepoRepository

// NewRepoRepository constructs a new RepoRepository with the provided sqlExecutor.
func NewRepoRepository(conn sqlExecutor) *RepoRepository {
	return &RepoRepository{conn: conn}
}

// DefaultDB returns the open default database instance.
func DefaultDB() (*DB, error) {
	return OpenDefault()
}

// ListByBranch returns repositories configured for the specified branch.
func (r *RepoRepository) ListByBranch(branch string) result.Result[[]model.ScanRecord] {
	clean := strings.TrimSpace(branch)
	if clean == "" {
		return r.queryScanRecords(constants.SQLSelectAllRepos)
	}

	return r.queryScanRecords(constants.SQLSelectReposByBranch, clean)
}

// FindByBranch delegates to ListByBranch.
func (r *RepoRepository) FindByBranch(branch string) result.Result[[]model.ScanRecord] {
	return r.ListByBranch(branch)
}

// ListDirty returns repositories flagged as dirty or requiring synchronization.
func (r *RepoRepository) ListDirty() result.Result[[]model.ScanRecord] {
	return r.queryScanRecords(constants.SQLSelectDirtyRepos)
}

// FindDirty delegates to ListDirty.
func (r *RepoRepository) FindDirty() result.Result[[]model.ScanRecord] {
	return r.ListDirty()
}

// ListByTag returns repositories matching a specific tag stored in metadata or notes.
func (r *RepoRepository) ListByTag(tag string) result.Result[[]model.ScanRecord] {
	clean := strings.TrimSpace(tag)
	if clean == "" {
		return r.queryScanRecords(constants.SQLSelectAllRepos)
	}

	pattern := "%" + clean + "%"

	return r.queryScanRecords(constants.SQLSelectReposByTag, pattern)
}

// FindByTag delegates to ListByTag.
func (r *RepoRepository) FindByTag(tag string) result.Result[[]model.ScanRecord] {
	return r.ListByTag(tag)
}

// ListUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (r *RepoRepository) ListUnpushed() result.Result[[]model.ScanRecord] {
	return r.queryScanRecords(constants.SQLSelectUnpushedRepos)
}

// FindUnpushed delegates to ListUnpushed.
func (r *RepoRepository) FindUnpushed() result.Result[[]model.ScanRecord] {
	return r.ListUnpushed()
}

// FindPathBySlug resolves the absolute path of a repository by full or short slug.
func (r *RepoRepository) FindPathBySlug(slug string) result.Result[string] {
	rq, ok := r.conn.(sqlRowQueryer)
	if !ok {
		return result.Fail[string](apperror.NewSimple("database executor does not support row query", "E_DB_EXECUTOR"))
	}

	parts := strings.Split(slug, "/")
	shortSlug := parts[len(parts)-1]
	row := QueryRowWrapper(rq, constants.SQLSelectRepoPathBySlug, slug, shortSlug)

	var absPath string
	scanErr := row.Scan(&absPath)
	if scanErr != nil {
		return result.Fail[string](apperror.WrapSimple(scanErr, "store.RepoRepository.FindPathBySlug"))
	}

	return result.Ok(absPath)
}

func (r *RepoRepository) queryScanRecords(query string, args ...any) result.Result[[]model.ScanRecord] {
	q, ok := r.conn.(sqlQueryer)
	if !ok {
		return result.Fail[[]model.ScanRecord](apperror.NewSimple("database executor does not support query", "E_DB_EXECUTOR"))
	}

	res := QueryWrapper(q, query, args...)
	rows, err := res.Destruct()
	if err != nil {
		return result.Fail[[]model.ScanRecord](apperror.WrapSimple(err, "store.RepoRepository.queryScanRecords"))
	}
	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return result.Fail[[]model.ScanRecord](apperror.WrapSimple(scanErr, "store.RepoRepository.scanRows"))
	}

	return result.Ok(records)
}
