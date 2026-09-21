package cmdpipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
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

// PipelineCacheDecision encapsulates the result of evaluating SQLite DB cache for pipeline telemetry.
type PipelineCacheDecision struct {
	IsFromCache bool
	CacheSource string
	CachedRuns  []ghRunItem
	TargetSha   string
	Reason      string
}

// EvaluatePipelineErrorsCache checks whether pipeline error telemetry can be served from SQLite.
func EvaluatePipelineErrorsCache(repo string, flags PipelineErrorFlags) PipelineCacheDecision {
	if isPipelineCacheBypassed(flags) {
		return PipelineCacheDecision{IsFromCache: false, Reason: "bypassed"}
	}

	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return PipelineCacheDecision{IsFromCache: false, Reason: "db_error"}
	}
	defer db.Close()

	return evaluateOpenDbCache(db, flags)
}

func evaluateOpenDbCache(db *pipelinedb.PipelineSplitDb, flags PipelineErrorFlags) PipelineCacheDecision {
	runRes := db.QueryRecentRuns(20)
	if runRes.IsFailure() || len(runRes.Data) == 0 {
		return PipelineCacheDecision{IsFromCache: false, Reason: "empty_db"}
	}

	return evaluateDecisionFromRuns(db, runRes.Data, flags)
}

func isPipelineCacheBypassed(flags PipelineErrorFlags) bool {
	if flags.HasForce || flags.HasTimeline {
		return true
	}

	return false
}

func evaluateDecisionFromRuns(db *pipelinedb.PipelineSplitDb, dbRuns []pipelinedb.PipelineRunRecord, flags PipelineErrorFlags) PipelineCacheDecision {
	latest := dbRuns[0]
	if checkTargetIndexCacheHit(dbRuns, flags) {
		return buildCacheHitDecision(dbRuns, latest.Sha, "target_index_matched")
	}
	if checkTtlCacheHit(db.Path) {
		return buildCacheHitDecision(dbRuns, latest.Sha, "within_ttl")
	}

	return PipelineCacheDecision{IsFromCache: false, Reason: "cache_expired"}
}

func resolvePipelineCacheTTL() time.Duration {
	return config.ResolvePipelineCacheTTL()
}

func checkTtlCacheHit(dbPath string) bool {
	syncPath := resolveCacheSyncFilePath(dbPath)
	if meta, isFound := readCacheSyncMeta(syncPath); isFound {
		return time.Since(meta.FetchedAt) < resolvePipelineCacheTTL()
	}

	info, err := os.Stat(dbPath)
	if err != nil {
		return false
	}

	return time.Since(info.ModTime()) < resolvePipelineCacheTTL()
}

func isRunCompleted(status, conclusion string) bool {
	if status == "completed" {
		return true
	}

	lower := strings.ToLower(conclusion)
	return lower == "success" || lower == "failure"
}

func matchCommitSha(shaA, shaB string) bool {
	if len(shaA) == 0 || len(shaB) == 0 {
		return false
	}

	return strings.EqualFold(shaA, shaB) || strings.HasPrefix(shaA, shaB) || strings.HasPrefix(shaB, shaA)
}

func checkTargetIndexCacheHit(dbRuns []pipelinedb.PipelineRunRecord, flags PipelineErrorFlags) bool {
	if !flags.HasIndex || flags.Index >= 0 {
		return false
	}

	offset := NormalizeNegativeIndex(flags.Index)
	return offset < len(dbRuns)
}

func buildCacheHitDecision(dbRuns []pipelinedb.PipelineRunRecord, sha, reason string) PipelineCacheDecision {
	return PipelineCacheDecision{
		IsFromCache: true,
		CacheSource: "sqlite",
		CachedRuns:  mapDbRunsToGhRuns(dbRuns),
		TargetSha:   sha,
		Reason:      reason,
	}
}
