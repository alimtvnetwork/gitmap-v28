package cmdpipeline

import (
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/config"
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

	return evaluateOpenDbCache(db, repo, flags)
}

func evaluateOpenDbCache(db *pipelinedb.PipelineSplitDb, repo string, flags PipelineErrorFlags) PipelineCacheDecision {
	runRes := db.QueryRecentRuns(20)
	if runRes.IsFailure() || len(runRes.Data) == 0 {
		return PipelineCacheDecision{IsFromCache: false, Reason: "empty_db"}
	}

	return evaluateDecisionFromRuns(db, repo, runRes.Data, flags)
}

func isPipelineCacheBypassed(flags PipelineErrorFlags) bool {
	if flags.HasForce || flags.HasTimeline || flags.HasUntilDone {
		return true
	}

	return false
}

func evaluateDecisionFromRuns(db *pipelinedb.PipelineSplitDb, repo string, dbRuns []pipelinedb.PipelineRunRecord, flags PipelineErrorFlags) PipelineCacheDecision {
	latest := dbRuns[0]
	if decision, isTargetSpecified := evaluateExplicitTargetDecision(repo, dbRuns, flags, latest.Sha); isTargetSpecified {
		return decision
	}

	if hasAnyActiveRun(dbRuns, latest.Sha) {
		return PipelineCacheDecision{IsFromCache: false, Reason: "active_run_in_progress"}
	}

	if checkCompletedCommitCacheHit(repo, dbRuns) {
		return buildCacheHitDecision(dbRuns, latest.Sha, "latest_commit_completed")
	}

	headSha := ResolveRepoHeadCommitSha(repo)
	if len(headSha) > 0 && !matchCommitSha(headSha, latest.Sha) {
		return PipelineCacheDecision{IsFromCache: false, Reason: "cache_expired"}
	}

	if checkTtlCacheHit(db.Path) && !hasAnyActiveRun(dbRuns, latest.Sha) {
		return buildCacheHitDecision(dbRuns, latest.Sha, "within_ttl")
	}

	return PipelineCacheDecision{IsFromCache: false, Reason: "cache_expired"}
}

func hasAnyActiveRun(dbRuns []pipelinedb.PipelineRunRecord, sha string) bool {
	for _, r := range dbRuns {
		if sha != "" && !matchCommitSha(r.Sha, sha) {
			continue
		}

		if !isRunCompleted(r.Status, r.Conclusion) {
			return true
		}
	}

	return false
}

func evaluateExplicitTargetDecision(repo string, dbRuns []pipelinedb.PipelineRunRecord, flags PipelineErrorFlags, latestSha string) (PipelineCacheDecision, bool) {
	if flags.CommitTarget == "" && !flags.HasIndex {
		return PipelineCacheDecision{}, false
	}

	if hit, decision := evaluateTargetOrHeadCacheHit(repo, dbRuns, flags, latestSha); hit {
		return decision, true
	}

	return PipelineCacheDecision{IsFromCache: false, Reason: "target_not_in_cache"}, true
}

func evaluateTargetOrHeadCacheHit(repo string, dbRuns []pipelinedb.PipelineRunRecord, flags PipelineErrorFlags, latestSha string) (bool, PipelineCacheDecision) {
	if checkCommitTargetCacheHit(dbRuns, flags.CommitTarget) {
		return true, buildCacheHitDecision(dbRuns, flags.CommitTarget, "commit_target_matched")
	}

	if checkTargetIndexCacheHit(dbRuns, flags) {
		targetSha := resolveTargetIndexSha(flags.CommitTarget, latestSha)
		return true, buildCacheHitDecision(dbRuns, targetSha, "target_index_matched")
	}

	if checkHeadShaCacheHit(repo, dbRuns) {
		return true, buildCacheHitDecision(dbRuns, latestSha, "head_sha_matched")
	}

	return false, PipelineCacheDecision{}
}

func checkCommitTargetCacheHit(dbRuns []pipelinedb.PipelineRunRecord, target string) bool {
	if len(target) == 0 {
		return false
	}

	hasMatch := false
	for _, r := range dbRuns {
		if matchCommitSha(r.Sha, target) {
			hasMatch = true
			break
		}
	}

	return hasMatch && isCommitRunsCompleted(dbRuns, target)
}

func checkHeadShaCacheHit(repo string, dbRuns []pipelinedb.PipelineRunRecord) bool {
	headSha := ResolveRepoHeadCommitSha(repo)
	if len(headSha) == 0 {
		return false
	}

	hasMatch := false
	for _, r := range dbRuns {
		if matchCommitSha(r.Sha, headSha) {
			hasMatch = true
			break
		}
	}

	return hasMatch && isCommitRunsCompleted(dbRuns, headSha)
}

func checkCompletedCommitCacheHit(repo string, dbRuns []pipelinedb.PipelineRunRecord) bool {
	if len(dbRuns) == 0 {
		return false
	}
	headSha := ResolveRepoHeadCommitSha(repo)
	latestSha := dbRuns[0].Sha
	if len(headSha) > 0 && !matchCommitSha(headSha, latestSha) {
		return false
	}
	return isCommitRunsCompleted(dbRuns, latestSha)
}

func isCommitRunsCompleted(dbRuns []pipelinedb.PipelineRunRecord, sha string) bool {
	hasCompleted := false
	for _, r := range dbRuns {
		if !matchCommitSha(r.Sha, sha) {
			continue
		}
		if !isRunCompleted(r.Status, r.Conclusion) {
			return false
		}
		hasCompleted = true
	}

	return hasCompleted
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
	if offset >= len(dbRuns) {
		return false
	}

	targetSha := resolveTargetCommitCandidate(flags.CommitTarget)
	if len(targetSha) > 0 {
		return checkCommitTargetCacheHit(dbRuns, targetSha)
	}

	return true
}

func resolveTargetIndexSha(commitTarget, fallbackSha string) string {
	target := resolveTargetCommitCandidate(commitTarget)
	if len(target) > 0 {
		return target
	}

	return fallbackSha
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
