package cmdpipeline

import (
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// TryRenderFromShaCache attempts to retrieve and render pipeline error telemetry from SQLite PipelineShaCache.
func TryRenderFromShaCache(repo string, flags PipelineErrorFlags) (bool, error) {
	if isPipelineCacheBypassed(flags) {
		return false, nil
	}

	targetSha := resolveTargetOrHeadSha(repo, flags)
	if len(targetSha) == 0 {
		return false, nil
	}

	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return false, nil
	}
	defer db.Close()

	payloadJson, err := db.QueryShaCache(targetSha, repo)
	if err != nil || len(payloadJson) == 0 {
		return false, nil
	}

	var payload PipelineErrorLogsPayload
	if err := json.Unmarshal([]byte(payloadJson), &payload); err != nil {
		return false, nil
	}

	activeLogLineLimit = flags.Limit
	payload.IsFromCache = true
	payload.CacheSource = "sha_cache"
	applyPayloadOptions(&payload, payload.Runs, flags)
	_ = RecordPipelineCheckInTask(repo, payload)

	return true, writePayloadWithFlags(payload, flags)
}

// SavePayloadToShaCache persists completed pipeline error telemetry payload into SQLite PipelineShaCache.
func SavePayloadToShaCache(repo string, payload PipelineErrorLogsPayload) {
	if payload.IsRunning || len(payload.LastHash) == 0 || payload.RunId == 0 {
		return
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return
	}
	defer db.Close()

	payloadJson := string(jsonBytes)
	_ = db.RecordShaCache(payload.LastHash, payload.RunId, repo, payloadJson)
	if len(payload.Sha) > 0 && payload.Sha != payload.LastHash {
		_ = db.RecordShaCache(payload.Sha, payload.RunId, repo, payloadJson)
	}
}
