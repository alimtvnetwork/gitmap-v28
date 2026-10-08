package pipelinedb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
