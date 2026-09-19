package pipelinedb

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
	candidates := []string{
		filepath.Join(dir, "pipeline_"+slug+".db"),
		filepath.Join(dir, slug+".db"),
		filepath.Join(store.BinaryDataDir(), "pipeline_db", "pipeline_"+slug+".db"),
	}
	for _, cand := range candidates {
		if isFileExisting(cand) {
			if err := os.Rename(cand, targetPath); err == nil {
				return targetPath
			}

			return cand
		}
	}

	return targetPath
}
