package pipelinedb

import (
	"os"
	"path/filepath"
)

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func isDirExisting(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func migrateOrFallbackPipelineDb(dir, slug, targetPath string) string {
	canonical := store.ResolveSplitDbPath(store.SectionPipeline, slug, "")
	if isFileExisting(canonical) {
		return canonical
	}

	return findAndMigratePipelineCandidate(dir, slug, canonical)
}

func findAndMigratePipelineCandidate(dir, slug, canonical string) string {
	candidates := collectPipelineCandidates(dir, slug)
	for _, cand := range candidates {
		if isFileExisting(cand) && tryMigratePipelineDb(cand, canonical) {
			return canonical
		}
	}

	return canonical
}

func tryMigratePipelineDb(cand, target string) bool {
	return os.Rename(cand, target) == nil
}

func collectPipelineCandidates(dir, slug string) []string {
	return []string{
		filepath.Join(dir, slug, "sql.db"),
		filepath.Join(dir, slug, "pipeline.db"),
		filepath.Join(dir, "pipeline_"+slug+".db"),
		filepath.Join(dir, slug+".db"),
		filepath.Join(store.BinaryDataDir(), "pipeline_db", "pipeline_"+slug+".db"),
	}
}
