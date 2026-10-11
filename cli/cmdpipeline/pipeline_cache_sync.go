package cmdpipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// PipelineCacheSyncMeta tracks synchronization telemetry to evaluate cache freshness.
type PipelineCacheSyncMeta struct {
	FetchedAt  time.Time `json:"fetchedAt"`
	FetchedSha string    `json:"fetchedSha,omitempty"`
	RunCount   int       `json:"runCount"`
}

func resolveCacheSyncFilePath(dbPath string) string {
	if len(dbPath) == 0 {
		return ""
	}

	return filepath.Join(filepath.Dir(dbPath), "cache_sync.json")
}

func writeCacheSyncMeta(dbPath string, runs []ghRunItem) {
	syncPath := resolveCacheSyncFilePath(dbPath)
	if len(syncPath) == 0 {
		return
	}

	meta := PipelineCacheSyncMeta{
		FetchedAt: time.Now().UTC(),
		RunCount:  len(runs),
	}
	if len(runs) > 0 {
		meta.FetchedSha = runs[0].HeadSha
	}

	data, err := json.Marshal(meta)
	if err == nil {
		_ = os.WriteFile(syncPath, data, 0644)
	}
}

func readCacheSyncMeta(syncPath string) (PipelineCacheSyncMeta, bool) {
	if len(syncPath) == 0 {
		return PipelineCacheSyncMeta{}, false
	}

	data, err := os.ReadFile(syncPath)
	if err != nil || len(data) == 0 {
		return PipelineCacheSyncMeta{}, false
	}

	var meta PipelineCacheSyncMeta
	if err := json.Unmarshal(data, &meta); err != nil || meta.FetchedAt.IsZero() {
		return PipelineCacheSyncMeta{}, false
	}

	return meta, true
}
