// Package cmdagy — agy_undo_cache.go restores conversations and cache from temp staging.
package cmdagy

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func restoreLatestCacheClearBackup() (int, bool, error) {
	baseDir := filepath.Join(os.TempDir(), "gitmap-agy-cache-backup")
	entries, err := os.ReadDir(baseDir)
	if err != nil || len(entries) == 0 {
		return 0, false, nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	latestDir := filepath.Join(baseDir, entries[0].Name())
	manifestPath := filepath.Join(latestDir, "manifest.json")
	m, readErr := loadCacheUndoManifest(manifestPath)
	if readErr != nil {
		return 0, false, readErr
	}

	restored := restoreFromCacheBackupDir(latestDir, m)
	return restored, true, nil
}

func loadCacheUndoManifest(path string) (*CacheUndoManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m CacheUndoManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func restoreFromCacheBackupDir(dir string, m *CacheUndoManifest) int {
	home, _ := os.UserHomeDir()
	convDst := filepath.Join(home, ".gemini", "antigravity", "conversations")
	brainDst := filepath.Join(home, ".gemini", "antigravity", "brain")
	_ = os.MkdirAll(convDst, constants.DirPermission)
	_ = os.MkdirAll(brainDst, constants.DirPermission)

	dbPath, _ := getConversationSummariesDBPath()
	conn, _ := store.OpenSQLiteDB(dbPath)
	if conn != nil {
		defer conn.Close()
	}

	count := 0
	for _, c := range m.Conversations {
		restoreSingleConv(dir, convDst, brainDst, c.ID)
		if conn != nil {
			restoreSummaryRecord(conn, c)
		}
		count++
	}
	return count
}

func restoreSingleConv(srcDir, convDst, brainDst, id string) {
	copyFileIfExists(filepath.Join(srcDir, "conversations", id+".db"), filepath.Join(convDst, id+".db"))
	copyFileIfExists(filepath.Join(srcDir, "conversations", id+".db-wal"), filepath.Join(convDst, id+".db-wal"))
	copyFileIfExists(filepath.Join(srcDir, "conversations", id+".db-shm"), filepath.Join(convDst, id+".db-shm"))
	copyDirIfExists(filepath.Join(srcDir, "brain", id), filepath.Join(brainDst, id))
}

func restoreSummaryRecord(conn *sql.DB, c ConvPruneCandidate) {
	query := "INSERT OR IGNORE INTO conversation_summaries (conversation_id, title, project_id, last_modified_time) VALUES (?, ?, ?, ?)"
	_ = store.ExecWrapper(conn, query, c.ID, c.Title, c.ProjectID, c.LastModified)
}
