package cmdpipeline

import (
	"testing"
)

func TestTryRenderFromShaCache_Bypassed(t *testing.T) {
	flags := PipelineErrorFlags{
		HasForce: true,
	}
	hit, err := TryRenderFromShaCache("test/repo", flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hit {
		t.Fatalf("expected cache miss when bypassed")
	}
}

func TestResolveTargetShaFast_DirectSha(t *testing.T) {
	flags := PipelineErrorFlags{
		CommitTarget: "0e6da8da9a5ddf8888a2e737cfbb2c65a480bebd",
		HasIndex:     false,
	}
	sha := resolveTargetShaFast("alimtvnetwork/gitmap-v28", flags)
	if sha != "0e6da8da9a5ddf8888a2e737cfbb2c65a480bebd" {
		t.Fatalf("expected direct sha, got: %s", sha)
	}
}

func TestFastResolveHeadCommitSha(t *testing.T) {
	sha := fastResolveHeadCommitSha("")
	if len(sha) == 0 {
		t.Log("no local head sha resolved from disk (might be CI or packed refs)")
		return
	}
	if len(sha) != 40 {
		t.Fatalf("expected 40-char commit sha, got: %s", sha)
	}
}

func TestResolveJobLogConcurrency(t *testing.T) {
	if c := resolveJobLogConcurrency(0); c != 1 {
		t.Fatalf("expected 1, got %d", c)
	}
	if c := resolveJobLogConcurrency(4); c != 4 {
		t.Fatalf("expected 4, got %d", c)
	}
	if c := resolveJobLogConcurrency(16); c != 8 {
		t.Fatalf("expected 8, got %d", c)
	}
}
