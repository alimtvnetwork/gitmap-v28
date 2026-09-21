package cmdpipeline

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

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
	if checkTtlCacheHit(db.Path) {
		return buildCacheHitDecision(dbRuns, latest.Sha, "within_5s_ttl")
	}
	if checkCommitMatchCacheHit(latest, resolveLocalCommitSHA()) {
		return buildCacheHitDecision(dbRuns, latest.Sha, "commit_match_completed")
	}
	if checkTargetIndexCacheHit(dbRuns, flags) {
		return buildCacheHitDecision(dbRuns, latest.Sha, "target_index_matched")
	}

	return PipelineCacheDecision{IsFromCache: false, Reason: "no_cache_match"}
}

func resolvePipelineCacheTTL() time.Duration {
	envVal := os.Getenv("GITMAP_PIPELINE_CACHE_TTL_SEC")
	if len(envVal) == 0 {
		return 5 * time.Second
	}

	sec, err := strconv.Atoi(envVal)
	if err == nil && sec > 0 {
		return time.Duration(sec) * time.Second
	}

	return 5 * time.Second
}

func checkTtlCacheHit(dbPath string) bool {
	info, err := os.Stat(dbPath)
	if err != nil {
		return false
	}

	return time.Since(info.ModTime()) < resolvePipelineCacheTTL()
}

func checkCommitMatchCacheHit(latest pipelinedb.PipelineRunRecord, localSha string) bool {
	if !matchCommitSha(latest.Sha, localSha) {
		return false
	}

	return isRunCompleted(latest.Status, latest.Conclusion)
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
