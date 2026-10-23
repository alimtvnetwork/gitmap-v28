package cmdpipeline

import (
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// Regression: the empty-runs path used to fabricate Conclusion="success",
// which made `gitmap pe` render CLEAN / 100% green whenever the gh CLI was
// unauthenticated and the local pipeline DB was empty. Empty runs must be
// reported as no-data with an unknown conclusion instead.
func TestBuildLocalOrEmptyErrorPayloadMarksNoData(t *testing.T) {
	payload := buildLocalOrEmptyErrorPayload(PipelineErrorLogsPayload{Repo: "example/no-data-repo"})

	if payload.Conclusion == "failure" {
		t.Skip("a local .gitmap/last_error.log is present in the test working directory")
	}
	if !payload.IsNoData {
		t.Fatalf("expected IsNoData to be set, got payload %+v", payload)
	}
	if payload.Conclusion != "unknown" {
		t.Fatalf("expected conclusion %q, got %q", "unknown", payload.Conclusion)
	}
	if isFailingConclusion(payload.Conclusion) {
		t.Fatalf("unknown conclusion must not be classified as failing")
	}
}

func TestBuildErrorLogsPayloadWithEmptyRunsIsNoData(t *testing.T) {
	t.Setenv("GITMAP_MOCK_GH", "1")

	payload := buildErrorLogsPayloadWithFlags("example/no-data-repo", nil, PipelineErrorFlags{})

	if !payload.IsNoData {
		t.Fatalf("expected no-data payload for empty runs, got conclusion %q", payload.Conclusion)
	}
	if payload.Conclusion != "unknown" {
		t.Fatalf("expected conclusion %q, got %q", "unknown", payload.Conclusion)
	}
}

// Regression: a cached payload that carries no runs must never be served
// as a verdict from the SHA cache.
func TestTryRenderFromShaCacheRejectsNoDataPayload(t *testing.T) {
	repo := "example/no-data-cache"
	sha := "abc123abc123abc123abc123abc123abc123abcd"

	payload := PipelineErrorLogsPayload{
		Repo:       repo,
		Conclusion: "unknown",
		IsNoData:   true,
		RunId:      42,
		LastHash:   sha,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		t.Fatalf("open pipeline db: %v", err)
	}
	_ = db.RecordShaCache(sha, 42, repo, string(jsonBytes))
	db.Close()

	hit, err := TryRenderFromShaCache(repo, PipelineErrorFlags{CommitTarget: sha})
	if err != nil {
		t.Fatalf("TryRenderFromShaCache returned error: %v", err)
	}
	if hit {
		t.Fatalf("expected no-data payload to be rejected from the SHA cache")
	}
}
