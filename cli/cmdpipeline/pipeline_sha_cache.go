package cmdpipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// TryRenderFromShaCache attempts to retrieve and render pipeline error telemetry from SQLite PipelineShaCache.
func TryRenderFromShaCache(repo string, flags PipelineErrorFlags) (bool, error) {
	if isPipelineCacheBypassed(flags) {
		return false, nil
	}

	targetSha := resolveTargetShaFast(repo, flags)
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

func resolveTargetShaFast(repo string, flags PipelineErrorFlags) string {
	if len(flags.CommitTarget) > 0 && !flags.HasIndex {
		return flags.CommitTarget
	}

	if len(flags.CommitTarget) > 0 {
		resolved := resolveTargetCommitCandidate(flags.CommitTarget)
		if len(resolved) > 0 {
			return resolved
		}
	}

	if sha := fastResolveHeadCommitSha(repo); len(sha) > 0 {
		return sha
	}

	return resolveTargetOrHeadSha(repo, flags)
}

func fastResolveHeadCommitSha(repo string) string {
	if !isCurrentRepoMatch(repo) {
		repoDir := findLocalRepoDirForSlug(repo)
		if len(repoDir) == 0 {
			return ""
		}
		return readGitHeadSha(repoDir)
	}

	return readGitHeadSha(".")
}

func readGitHeadSha(dir string) string {
	headPath := filepath.Join(dir, ".git", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}

	content := strings.TrimSpace(string(data))
	if len(content) == 40 && !strings.HasPrefix(content, "ref:") {
		return content
	}

	if !strings.HasPrefix(content, "ref: ") {
		return ""
	}

	refRel := strings.TrimSpace(content[5:])
	refPath := filepath.Join(dir, ".git", filepath.FromSlash(refRel))
	refData, err := os.ReadFile(refPath)
	if err != nil {
		return ""
	}

	sha := strings.TrimSpace(string(refData))
	if len(sha) == 40 {
		return sha
	}

	return ""
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
