package pipelinedb

import (
	"database/sql"
	"os"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"fmt"
	"path/filepath"
	"strings"
)

// Reset drops all tables and re-initializes the schema.
func (p *PipelineSplitDb) Reset() error {
	for _, q := range dropTableQueries() {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "reset pipeline split db")
		}
	}

	return p.InitSchema()
}

func (p *PipelineSplitDb) optimizePragmas() *apperror.AppError {
	if _, err := p.conn.Exec("PRAGMA wal_checkpoint(TRUNCATE);"); err != nil {
		return apperror.WrapSimple(err, "wal checkpoint pipeline db")
	}

	if _, err := p.conn.Exec("VACUUM;"); err != nil {
		return apperror.WrapSimple(err, "vacuum pipeline db")
	}

	if _, err := p.conn.Exec("PRAGMA optimize;"); err != nil {
		return apperror.WrapSimple(err, "optimize pipeline db")
	}

	return nil
}

// Optimize executes WAL checkpoint and VACUUM, returning reclaimed bytes.
func (p *PipelineSplitDb) Optimize() (int64, error) {
	sizeBefore := getFileSize(p.Path)
	if err := p.optimizePragmas(); err != nil {
		return 0, err
	}

	sizeAfter := getFileSize(p.Path)
	if sizeBefore <= sizeAfter {
		return 0, nil
	}

	return sizeBefore - sizeAfter, nil
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	if info.Size() < 0 {
		return 0
	}

	return info.Size()
}

func safeInt64ToUint64(val int64) uint64 {
	if val < 0 {
		return 0
	}

	return uint64(val)
}

func countQuery(conn *sql.DB, query string) (int, *apperror.AppError) {
	var count int
	if err := conn.QueryRow(query).Scan(&count); err != nil {
		return 0, apperror.WrapSimple(err, "count query: "+query)
	}

	return count, nil
}

func queryLastUpdated(conn *sql.DB) (string, *apperror.AppError) {
	var lastUpdated string
	query := "SELECT COALESCE(MAX(UpdatedAt), '') FROM PipelineRun;"
	if err := conn.QueryRow(query).Scan(&lastUpdated); err != nil {
		return "", apperror.WrapSimple(err, "query last updated")
	}

	return lastUpdated, nil
}

func (p *PipelineSplitDb) purgeCacheFiles(runIds []uint64) {
	dir := filepath.Dir(p.Path)
	purgeRunIdCacheFiles(dir, runIds)
	purgeRepoCacheDir(filepath.Join(dir, strings.ReplaceAll(p.RepoSlug, "/", "_")))
	purgeRepoCacheDir(filepath.Join(dir, SanitizeRepoSlug(p.RepoSlug)))
	purgeRepoMatchingJsonFiles(dir, p.RepoSlug)
	purgePipelineReports(dir)
	purgeParentLegacyFiles(dir, p.RepoSlug, runIds)
}

func purgePipelineReports(dir string) {
	_ = os.Remove(filepath.Join(dir, "pipeline_errors.log"))
	_ = os.Remove(filepath.Join(dir, "last_error.log"))
}

func purgeParentLegacyFiles(dir, repoSlug string, runIds []uint64) {
	parent := filepath.Dir(dir)
	if filepath.Base(parent) == "pipeline" {
		slug := SanitizeRepoSlug(repoSlug)
		_ = os.Remove(filepath.Join(parent, "pipeline_"+slug+".db"))
		_ = os.Remove(filepath.Join(parent, slug+".db"))
		purgeRunIdCacheFiles(parent, runIds)
	}
}

func purgeRunIdCacheFiles(dir string, runIds []uint64) {
	for _, id := range runIds {
		_ = os.Remove(filepath.Join(dir, fmt.Sprintf("%d.log", id)))
		_ = os.Remove(filepath.Join(dir, fmt.Sprintf("%d.json", id)))
		_ = os.Remove(filepath.Join(dir, fmt.Sprintf("%d.jobs.json", id)))
	}
}

func purgeRepoCacheDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() && isPurgeableCacheFile(e.Name()) {
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
	_ = os.Remove(dir)
}

func isPurgeableCacheFile(name string) bool {
	lower := strings.ToLower(name)

	return strings.HasSuffix(lower, ".log") || strings.HasSuffix(lower, ".json")
}

func purgeRepoMatchingJsonFiles(dir, repoSlug string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
			purgeIfRepoMatches(filepath.Join(dir, e.Name()), repoSlug)
		}
	}
}

func purgeIfRepoMatches(jsonPath, repoSlug string) {
	content, err := os.ReadFile(jsonPath)
	if err != nil || !strings.Contains(string(content), repoSlug) {
		return
	}
	_ = os.Remove(jsonPath)
	baseNoExt := strings.TrimSuffix(jsonPath, ".json")
	_ = os.Remove(baseNoExt + ".log")
	_ = os.Remove(strings.TrimSuffix(baseNoExt, ".jobs") + ".log")
}

func dropTableQueries() []string {
	return []string{
		"DROP TABLE IF EXISTS PipelineCompactErrorLog;",
		"DROP TABLE IF EXISTS PipelineDetailErrorLog;",
		"DROP TABLE IF EXISTS PipelineErrorLog;",
		"DROP TABLE IF EXISTS PipelineSegment;",
		"DROP TABLE IF EXISTS PipelineJob;",
		"DROP TABLE IF EXISTS PipelineRun;",
	}
}
