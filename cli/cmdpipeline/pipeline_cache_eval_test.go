package cmdpipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func TestIsPipelineCacheBypassed(t *testing.T) {
	flags := PipelineErrorFlags{HasForce: true}
	if !isPipelineCacheBypassed(flags) {
		t.Errorf("expected bypass for HasForce")
	}

	flagsTimeline := PipelineErrorFlags{HasTimeline: true}
	if !isPipelineCacheBypassed(flagsTimeline) {
		t.Errorf("expected bypass for HasTimeline")
	}

	flagsNormal := PipelineErrorFlags{}
	if isPipelineCacheBypassed(flagsNormal) {
		t.Errorf("expected no bypass for normal flags")
	}
}

func TestMatchCommitSha(t *testing.T) {
	if !matchCommitSha("abc1234567890", "abc1234") {
		t.Errorf("expected prefix match to succeed")
	}

	if !matchCommitSha("abc1234", "abc1234567890") {
		t.Errorf("expected reverse prefix match to succeed")
	}

	if matchCommitSha("abc1234", "def5678") {
		t.Errorf("expected different shas to fail")
	}

	if matchCommitSha("", "abc1234") {
		t.Errorf("expected empty sha to fail")
	}
}

func TestIsRunCompleted(t *testing.T) {
	if !isRunCompleted("completed", "success") {
		t.Errorf("expected completed run to be true")
	}

	if !isRunCompleted("finished", "failure") {
		t.Errorf("expected failure conclusion to be true")
	}

	if isRunCompleted("in_progress", "") {
		t.Errorf("expected in_progress run to be false")
	}

	if isRunCompleted("queued", "") {
		t.Errorf("expected queued run to be false")
	}
}

func TestCheckTargetIndexCacheHit(t *testing.T) {
	runs := []pipelinedb.PipelineRunRecord{
		{RunId: 101, Sha: "commit1"},
		{RunId: 102, Sha: "commit2"},
	}

	flagsValid := PipelineErrorFlags{HasIndex: true, Index: -1}
	if !checkTargetIndexCacheHit(runs, flagsValid) {
		t.Errorf("expected -1 index to hit cache")
	}

	flagsOutOfBounds := PipelineErrorFlags{HasIndex: true, Index: -5}
	if checkTargetIndexCacheHit(runs, flagsOutOfBounds) {
		t.Errorf("expected -5 index to miss cache")
	}
}

func TestEvaluatePipelineErrorsCache_ForceBypass(t *testing.T) {
	flags := PipelineErrorFlags{HasForce: true}
	decision := EvaluatePipelineErrorsCache("dummy/repo", flags)
	if decision.IsFromCache {
		t.Errorf("expected cache miss when force flag is set")
	}
	if decision.Reason != "bypassed" {
		t.Errorf("expected reason bypassed, got %s", decision.Reason)
	}
}

func TestCheckTtlCacheHit_WithFreshSyncMeta(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "sql.db")
	_ = os.WriteFile(dbPath, []byte("sqlite"), 0644)
	syncMeta := PipelineCacheSyncMeta{
		FetchedAt: time.Now().UTC(),
		RunCount:  3,
	}
	data, _ := json.Marshal(syncMeta)
	_ = os.WriteFile(filepath.Join(tempDir, "cache_sync.json"), data, 0644)

	if isHit := checkTtlCacheHit(dbPath); isHit == false {
		t.Errorf("expected cache hit for fresh sync meta within 5s")
	}
}

func TestCheckTtlCacheHit_WithExpiredSyncMeta(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "sql.db")
	_ = os.WriteFile(dbPath, []byte("sqlite"), 0644)
	syncMeta := PipelineCacheSyncMeta{
		FetchedAt: time.Now().UTC().Add(-10 * time.Second),
		RunCount:  3,
	}
	data, _ := json.Marshal(syncMeta)
	_ = os.WriteFile(filepath.Join(tempDir, "cache_sync.json"), data, 0644)

	if isHit := checkTtlCacheHit(dbPath); isHit {
		t.Errorf("expected cache miss for expired sync meta (10s old)")
	}
}

func TestCheckTtlCacheHit_FallbackDbStat(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "sql.db")
	_ = os.WriteFile(dbPath, []byte("sqlite"), 0644)

	if isHit := checkTtlCacheHit(dbPath); isHit == false {
		t.Errorf("expected cache hit for newly created db file fallback")
	}
}

func TestWriteAndReadCacheSyncMeta(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "sql.db")
	runs := []ghRunItem{{DatabaseId: 123, HeadSha: "abc"}}
	writeCacheSyncMeta(dbPath, runs)
	meta, isFound := readCacheSyncMeta(filepath.Join(tempDir, "cache_sync.json"))
	if isFound == false || meta.RunCount != 1 || meta.FetchedSha != "abc" {
		t.Errorf("expected sync meta to be written and read successfully")
	}
}
