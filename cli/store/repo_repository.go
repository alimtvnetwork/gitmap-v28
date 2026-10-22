// Package store manages repository database queries and persistence.
package store

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// RepoRepository encapsulates repository database queries returning typed AppErrors.
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
func (r *RepoRepository) ListByBranch(branch string) ([]model.ScanRecord, *apperror.AppError) {
	clean := strings.TrimSpace(branch)
	if clean == "" {
		return r.queryScanRecords(constants.SQLSelectAllRepos)
	}

	return r.queryScanRecords(constants.SQLSelectReposByBranch, clean)
}

// FindByBranch delegates to ListByBranch.
func (r *RepoRepository) FindByBranch(branch string) ([]model.ScanRecord, *apperror.AppError) {
	return r.ListByBranch(branch)
}

// ListDirty returns repositories flagged as dirty or requiring synchronization.
func (r *RepoRepository) ListDirty() ([]model.ScanRecord, *apperror.AppError) {
	return r.queryScanRecords(constants.SQLSelectDirtyRepos)
}

// FindDirty delegates to ListDirty.
func (r *RepoRepository) FindDirty() ([]model.ScanRecord, *apperror.AppError) {
	return r.ListDirty()
}

// ListByTag returns repositories matching a specific tag stored in metadata or notes.
func (r *RepoRepository) ListByTag(tag string) ([]model.ScanRecord, *apperror.AppError) {
	clean := strings.TrimSpace(tag)
	if clean == "" {
		return r.queryScanRecords(constants.SQLSelectAllRepos)
	}

	pattern := "%" + clean + "%"

	return r.queryScanRecords(constants.SQLSelectReposByTag, pattern)
}

// FindByTag delegates to ListByTag.
func (r *RepoRepository) FindByTag(tag string) ([]model.ScanRecord, *apperror.AppError) {
	return r.ListByTag(tag)
}

// ListUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (r *RepoRepository) ListUnpushed() ([]model.ScanRecord, *apperror.AppError) {
	return r.queryScanRecords(constants.SQLSelectUnpushedRepos)
}

// FindUnpushed delegates to ListUnpushed.
func (r *RepoRepository) FindUnpushed() ([]model.ScanRecord, *apperror.AppError) {
	return r.ListUnpushed()
}

// FindPathBySlug resolves the absolute path of a repository by full or short slug.
func (r *RepoRepository) FindPathBySlug(slug string) (string, *apperror.AppError) {
	rq, ok := r.conn.(sqlRowQueryer)
	if !ok {
		return "", apperror.NewSimple("database executor does not support row query", "E_DB_EXECUTOR")
	}

	parts := strings.Split(slug, "/")
	shortSlug := parts[len(parts)-1]
	row := QueryRowWrapper(rq, constants.SQLSelectRepoPathBySlug, slug, shortSlug)

	var absPath string
	scanErr := row.Scan(&absPath)
	if scanErr != nil {
		return "", apperror.WrapSimple(scanErr, "store.RepoRepository.FindPathBySlug")
	}

	return absPath, nil
}

func (r *RepoRepository) queryScanRecords(query string, args ...any) ([]model.ScanRecord, *apperror.AppError) {
	q, ok := r.conn.(sqlQueryer)
	if !ok {
		return nil, apperror.NewSimple("database executor does not support query", "E_DB_EXECUTOR")
	}

	res := QueryWrapper(q, query, args...)
	rows, err := res.Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.RepoRepository.queryScanRecords")
	}
	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return nil, apperror.WrapSimple(scanErr, "store.RepoRepository.scanRows")
	}

	return records, nil
}
