package cmdpipeline

import (
	"testing"
)

func TestFindPreviousFailingRunInHistory_ReturnsFirstFailedRun(t *testing.T) {
	runs := []ghRunItem{
		{DatabaseId: 101, Conclusion: "success", HeadSha: "aaa111"},
		{DatabaseId: 102, Conclusion: "success", HeadSha: "aaa111"},
		{DatabaseId: 100, Conclusion: "failure", HeadSha: "bbb222"},
		{DatabaseId: 99, Conclusion: "failure", HeadSha: "ccc333"},
	}

	failedRun, isFound := FindPreviousFailingRunInHistory(runs)
	if !isFound {
		t.Fatalf("expected failing run to be found, got false")
	}
	if failedRun.DatabaseId != 100 {
		t.Fatalf("expected run ID 100, got %d", failedRun.DatabaseId)
	}
	if failedRun.HeadSha != "bbb222" {
		t.Fatalf("expected SHA bbb222, got %s", failedRun.HeadSha)
	}
}

func TestFindPreviousFailingRunInHistory_ReturnsFalseWhenAllSuccess(t *testing.T) {
	runs := []ghRunItem{
		{DatabaseId: 101, Conclusion: "success", HeadSha: "aaa111"},
		{DatabaseId: 102, Conclusion: "success", HeadSha: "aaa111"},
	}

	_, isFound := FindPreviousFailingRunInHistory(runs)
	if isFound {
		t.Fatalf("expected no failing run found, got true")
	}
}

func TestResolveFallbackTargetGroup_ReturnsSecondGroupWhenAvailable(t *testing.T) {
	groups := []*CommitPipelineGroup{
		{CommitSha: "sha-latest", ShortSha: "latest"},
		{CommitSha: "sha-prev", ShortSha: "prev"},
	}

	fallback, isFound := ResolveFallbackTargetGroup(groups)
	if !isFound || fallback == nil {
		t.Fatalf("expected fallback group to be found")
	}
	if fallback.CommitSha != "sha-prev" {
		t.Fatalf("expected sha-prev, got %s", fallback.CommitSha)
	}
}

func TestResolveFallbackTargetGroup_ReturnsFalseWhenSingleGroup(t *testing.T) {
	groups := []*CommitPipelineGroup{
		{CommitSha: "sha-only", ShortSha: "only"},
	}

	_, isFound := ResolveFallbackTargetGroup(groups)
	if isFound {
		t.Fatalf("expected false for single group, got true")
	}
}
