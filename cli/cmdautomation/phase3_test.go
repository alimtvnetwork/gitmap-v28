package cmdautomation

import (
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func TestPreflightFiltering(t *testing.T) {
	opts := PreflightOptions{Filter: "format"}
	checks := resolvePreflightChecks(opts)
	if len(checks) == 0 {
		t.Errorf("expected at least 1 filtered check, got 0")
	}
}

func TestCheckScriptIssuesValid(t *testing.T) {
	content := "#!/bin/bash\nsha256sum file.tar.gz\nmv bin/tool /usr/local/bin\n"
	issues := checkScriptIssues("install.sh", content)
	if len(issues) > 0 {
		t.Errorf("expected 0 issues for valid script, got %d", len(issues))
	}
}

func TestCheckScriptIssuesInvalid(t *testing.T) {
	content := "#!/bin/bash\necho PLACEHOLDER\n"
	issues := checkScriptIssues("bad.sh", content)
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues for placeholder script, got %d", len(issues))
	}
}

func TestParseGoTestLineValid(t *testing.T) {
	line := "func TestSampleRunner(t *testing.T) {"
	item := parseGoTestLine(line, "cli/sample_test.go", "cli", line, 4.0)
	if item.Id != "cli.TestSampleRunner" {
		t.Errorf("unexpected test ID: %s", item.Id)
	}
	if item.TestFunc != "TestSampleRunner" {
		t.Errorf("unexpected test func: %s", item.TestFunc)
	}
}

func TestParseGoTestLineInvalid(t *testing.T) {
	line := "func helperFunc() {"
	item := parseGoTestLine(line, "cli/sample_test.go", "cli", line, 4.0)
	if len(item.Id) > 0 {
		t.Errorf("expected empty item for non-test func, got ID: %s", item.Id)
	}
}

func TestEstimateDurationFast(t *testing.T) {
	est := estimateDuration("TestFast", "a := 1\nb := 2", 4.0)
	if est.isSlow {
		t.Errorf("expected fast test, got slow")
	}
	if est.tier != "fast" {
		t.Errorf("expected fast tier, got %s", est.tier)
	}
}

func TestEstimateDurationSlow(t *testing.T) {
	est := estimateDuration("TestHeavyExecution", "exec.Command('sleep', '5')", 4.0)
	if !est.isSlow {
		t.Errorf("expected slow test, got fast")
	}
	if est.tier != "slow" {
		t.Errorf("expected slow tier, got %s", est.tier)
	}
}

func TestParseArtifactsJson(t *testing.T) {
	raw := []byte(`{"total_count": 1, "artifacts": [{"id": 12345, "name": "binaries", "size_in_bytes": 1048576, "created_at": "2026-09-19T00:00:00Z"}]}`)
	items := parseArtifactsJson(raw)
	if len(items) != 1 {
		t.Fatalf("expected 1 artifact item, got %d", len(items))
	}
	if items[0].Id != 12345 || items[0].SizeBytes != 1048576 {
		t.Errorf("parsed artifact mismatch: %+v", items[0])
	}
}

func TestSummarizeInventory(t *testing.T) {
	tests := map[string]TestInventoryItem{
		"cli.TestOne": {Id: "cli.TestOne", Package: "cli", DurationSec: 0.1, IsSlow: false},
		"cli.TestTwo": {Id: "cli.TestTwo", Package: "cli", DurationSec: 5.0, IsSlow: true},
	}
	summary := summarizeInventory(tests, 4.0)
	if summary.Total != 2 || summary.SlowTests != 1 || summary.FastTests != 1 {
		t.Errorf("unexpected summary values: %+v", summary)
	}
}

func TestPhase3Monads(t *testing.T) {
	prefRes := PreflightResult{TotalChecks: 5, PassedChecks: 5, IsPass: true}
	prefMonad := result.Ok(prefRes)
	if prefMonad.IsFailure() || !prefMonad.Value.IsPass {
		t.Errorf("expected successful PreflightResultMonad")
	}

	invRes := TestInventoryResult{TotalTests: 10, IsPass: true}
	invMonad := result.Ok(invRes)
	if invMonad.IsFailure() || invMonad.Value.TotalTests != 10 {
		t.Errorf("expected successful TestInventoryResultMonad")
	}
}

func TestAggregatePurgeResult(t *testing.T) {
	items := []PurgeArtifactItem{
		{Id: 1, Name: "art1", SizeBytes: 1048576, IsPurged: true},
		{Id: 2, Name: "art2", SizeBytes: 2097152, IsPurged: false},
	}
	res := aggregatePurgeResult("alimtvnetwork/gitmap-v28", items, time.Now())
	if res.PurgedArtifacts != 1 || res.FreedBytes != 1048576 {
		t.Errorf("unexpected purge result: %+v", res)
	}
}
