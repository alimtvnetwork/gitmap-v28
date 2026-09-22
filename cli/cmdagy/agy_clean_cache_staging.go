// Package cmdagy — agy_clean_cache_staging.go handles temp staging for undo.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type CacheUndoManifest struct {
	Timestamp     string               `json:"timestamp"`
	Conversations []ConvPruneCandidate `json:"conversations"`
}

func stagePrunedConvs(pruned []ConvPruneCandidate) (string, error) {
	if len(pruned) == 0 {
		return "", nil
	}

	ts := time.Now().UTC().Format("20060102-150405")
	stagingDir := filepath.Join(os.TempDir(), "gitmap-agy-cache-backup", ts)
	if err := os.MkdirAll(stagingDir, constants.DirPermission); err != nil {
		return "", err
	}

	copyPrunedArtifacts(stagingDir, pruned)
	saveCacheUndoManifest(stagingDir, ts, pruned)
	return stagingDir, nil
}

func copyPrunedArtifacts(stagingDir string, pruned []ConvPruneCandidate) {
	home, _ := os.UserHomeDir()
	convSrc := filepath.Join(home, ".gemini", "antigravity", "conversations")
	brainSrc := filepath.Join(home, ".gemini", "antigravity", "brain")

	convDst := filepath.Join(stagingDir, "conversations")
	brainDst := filepath.Join(stagingDir, "brain")
	_ = os.MkdirAll(convDst, constants.DirPermission)
	_ = os.MkdirAll(brainDst, constants.DirPermission)

	for _, c := range pruned {
		copyFileIfExists(filepath.Join(convSrc, c.ID+".db"), filepath.Join(convDst, c.ID+".db"))
		copyFileIfExists(filepath.Join(convSrc, c.ID+".db-wal"), filepath.Join(convDst, c.ID+".db-wal"))
		copyFileIfExists(filepath.Join(convSrc, c.ID+".db-shm"), filepath.Join(convDst, c.ID+".db-shm"))
		copyDirIfExists(filepath.Join(brainSrc, c.ID), filepath.Join(brainDst, c.ID))
	}
}

func saveCacheUndoManifest(dir, ts string, pruned []ConvPruneCandidate) {
	m := CacheUndoManifest{Timestamp: ts, Conversations: pruned}
	data, _ := json.MarshalIndent(m, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "manifest.json"), data, constants.FilePermission)
}

func copyFileIfExists(src, dst string) {
	content, err := os.ReadFile(src)
	if err == nil {
		_ = os.WriteFile(dst, content, constants.FilePermission)
	}
}

func copyDirIfExists(src, dst string) {
	if info, err := os.Stat(src); err == nil && info.IsDir() {
		_ = copyDirContents(src, dst)
	}
}
