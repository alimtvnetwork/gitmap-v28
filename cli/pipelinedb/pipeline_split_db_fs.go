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

func migrateOrFallbackPipelineDb(dir, slug, targetPrefixed string) string {
	legacyDir := filepath.Join(store.BinaryDataDir(), "pipeline_db")
	legacyFile := filepath.Join(legacyDir, "pipeline_"+slug+".db")
	hasLegacy := isFileExisting(legacyFile)
	if !hasLegacy {
		return targetPrefixed
	}

	err := os.Rename(legacyFile, targetPrefixed)
	if err == nil {
		return targetPrefixed
	}

	return legacyFile
}
