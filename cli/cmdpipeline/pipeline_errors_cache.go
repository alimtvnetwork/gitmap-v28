package cmdpipeline

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// ResolveLocalHeadCommitSha reads the latest Git commit SHA from local repository.
func ResolveLocalHeadCommitSha() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil || len(out) == 0 {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// QueryLocalPipelineDbBySha queries the local repository SQLite database for a matching commit.
func QueryLocalPipelineDbBySha(repo, sha string) (*pipelinedb.PipelineRunRecord, error) {
	if len(sha) < 7 {
		return nil, apperror.NewValidationError("commit sha too short for cache query")
	}
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db.QueryRunBySha(sha)
}

// CheckHeadCommitCachedInDb checks if current local HEAD commit SHA is recorded in SQLite DB.
func CheckHeadCommitCachedInDb(repo string) (*pipelinedb.PipelineRunRecord, bool) {
	headSha := ResolveLocalHeadCommitSha()
	if headSha == "" {
		return nil, false
	}
	record, err := QueryLocalPipelineDbBySha(repo, headSha)
	if err != nil || record == nil {
		return nil, false
	}
	if !isRunCompleted(record.Status, record.Conclusion) {
		return nil, false
	}
	return record, true
}

// TouchPipelineCacheIfHeadMatches refreshes the local cache sync metadata if HEAD commit matches DB.
func TouchPipelineCacheIfHeadMatches(repo string) bool {
	record, isCached := CheckHeadCommitCachedInDb(repo)
	if !isCached || record == nil {
		return false
	}
	return updateSyncMetaForRepo(repo)
}

func updateSyncMetaForRepo(repo string) bool {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return false
	}
	defer db.Close()
	runRes := db.QueryRecentRuns(20)
	if runRes.IsFailure() || len(runRes.Data) == 0 {
		return false
	}
	writeCacheSyncMeta(db.Path, mapDbRunsToGhRuns(runRes.Data))
	return true
}

// FetchCachedPipelineErrorReportFast retrieves cached error/clean report from DB if HEAD matches.
func FetchCachedPipelineErrorReportFast(repo string, isDetailed bool) (PipelineErrorLogsPayload, string, bool, error) {
	record, isCached := CheckHeadCommitCachedInDb(repo)
	if !isCached || record == nil {
		return PipelineErrorLogsPayload{}, "", false, apperror.NewValidationError("cache miss for head sha")
	}
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return PipelineErrorLogsPayload{}, "", false, err
	}
	defer db.Close()
	return buildFastCachedReport(db, repo, isDetailed)
}

func buildFastCachedReport(db *pipelinedb.PipelineSplitDb, repo string, isDetailed bool) (PipelineErrorLogsPayload, string, bool, error) {
	runRes := db.QueryRecentRuns(20)
	if runRes.IsFailure() || len(runRes.Data) == 0 {
		return PipelineErrorLogsPayload{}, "", false, apperror.NewValidationError("empty pipeline database")
	}
	payload := createCachedPayload(repo, runRes.Data, isDetailed)
	return formatCachedReportOutcome(payload)
}

func createCachedPayload(repo string, records []pipelinedb.PipelineRunRecord, isDetailed bool) PipelineErrorLogsPayload {
	payload := buildErrorLogsPayload(repo, mapDbRunsToGhRuns(records))
	payload.IsFromCache = true
	payload.CacheSource = "sqlite_head_sha"
	if !isDetailed {
		compactErrorPayload(&payload)
	}
	return payload
}

func formatCachedReportOutcome(payload PipelineErrorLogsPayload) (PipelineErrorLogsPayload, string, bool, error) {
	hasFail := payload.Conclusion == "failure" || len(payload.FailedRuns) > 0
	if hasFail {
		return payload, buildClipboardErrorReport(payload), true, nil
	}
	return payload, buildClipboardCleanReport(payload), false, nil
}

func init() {
	EnsureAgyFixRunnerExplicitOnly()
	targetRepo := resolveCurrentRepoSlug()
	if targetRepo != "" {
		_ = TouchPipelineCacheIfHeadMatches(targetRepo)
	}
}
