package cmdfixreleasetags

import (
	"strings"
	"testing"
	"time"
)

type testMockExecutor struct {
	responses map[string][]byte
	errors    map[string]error
	calls     [][]string
}

func newTestMockExecutor() *testMockExecutor {
	return &testMockExecutor{
		responses: make(map[string][]byte),
		errors:    make(map[string]error),
	}
}

func (m *testMockExecutor) Run(dir string, name string, args ...string) ([]byte, error) {
	key := name + " " + strings.Join(args, " ")
	m.calls = append(m.calls, append([]string{name}, args...))

	if err, ok := m.errors[key]; ok {
		return m.responses[key], err
	}

	if out, ok := m.responses[key]; ok {
		return out, nil
	}

	return nil, nil
}

func TestAuditCriteria_OrphanTag(t *testing.T) {
	reason := DetermineAuditReason(false, false, 0, false, false, nil)
	if reason != ReasonOrphanTag {
		t.Fatalf("expected ReasonOrphanTag, got %s", reason)
	}
}

func TestAuditCriteria_DraftRelease(t *testing.T) {
	reason := DetermineAuditReason(true, true, 2, true, false, nil)
	if reason != ReasonDraftRelease {
		t.Fatalf("expected ReasonDraftRelease, got %s", reason)
	}
}

func TestAuditCriteria_MissingAssets(t *testing.T) {
	reasonNoAssets := DetermineAuditReason(true, false, 0, false, false, nil)
	if reasonNoAssets != ReasonMissingAssets {
		t.Fatalf("expected ReasonMissingAssets for zero assets, got %s", reasonNoAssets)
	}

	reasonNoChecksum := DetermineAuditReason(true, false, 2, false, false, nil)
	if reasonNoChecksum != ReasonMissingAssets {
		t.Fatalf("expected ReasonMissingAssets for missing checksums, got %s", reasonNoChecksum)
	}
}

func TestAuditCriteria_CorruptedAssets(t *testing.T) {
	reason := DetermineAuditReason(true, false, 2, true, true, nil)
	if reason != ReasonCorruptedAssets {
		t.Fatalf("expected ReasonCorruptedAssets, got %s", reason)
	}
}

func TestAuditCriteria_CICDFailed(t *testing.T) {
	failingRuns := []CIWorkflowRunInfo{
		{RunId: 101, Status: "completed", Conclusion: "failure"},
	}

	reason := DetermineAuditReason(true, false, 2, true, false, failingRuns)
	if reason != ReasonCICDFailed {
		t.Fatalf("expected ReasonCICDFailed, got %s", reason)
	}
}

func TestSafetyGuard_ActiveVersion(t *testing.T) {
	status, eligible := EvaluateTagSafety("v6.529.0", "6.529.0", "v6.528.0", nil, time.Now(), 15*time.Minute)
	if status != StatusProtectedActive {
		t.Fatalf("expected StatusProtectedActive, got %s", status)
	}
	if eligible {
		t.Fatalf("expected eligible=false for active version")
	}

	rec := &ReleaseTagAuditRecord{
		Tag:         "v6.529.0",
		AuditReason: ReasonMissingAssets,
	}
	EvaluateSafety(rec, "6.529.0", "v6.528.0", 15*time.Minute)
	if !rec.IsProtected || rec.IsEligibleForDeletion {
		t.Fatalf("record safety violation: isProtected=%v eligible=%v", rec.IsProtected, rec.IsEligibleForDeletion)
	}
}

func TestSafetyGuard_LatestHealthy(t *testing.T) {
	status, eligible := EvaluateTagSafety("v6.528.0", "6.529.0", "v6.528.0", nil, time.Now(), 15*time.Minute)
	if status != StatusProtectedLatest {
		t.Fatalf("expected StatusProtectedLatest, got %s", status)
	}
	if eligible {
		t.Fatalf("expected eligible=false for latest healthy")
	}

	rec := &ReleaseTagAuditRecord{
		Tag:         "v6.528.0",
		AuditReason: ReasonHealthy,
	}
	EvaluateSafety(rec, "6.529.0", "v6.528.0", 15*time.Minute)
	if !rec.IsProtected || rec.IsEligibleForDeletion {
		t.Fatalf("record safety violation: isProtected=%v eligible=%v", rec.IsProtected, rec.IsEligibleForDeletion)
	}
}

func TestSafetyGuard_GraceWindow(t *testing.T) {
	now := time.Now()
	runs := []CIWorkflowRunInfo{
		{RunId: 999, Status: "in_progress", CreatedAt: now.Add(-5 * time.Minute), IsInProgress: true},
	}

	status, eligible := EvaluateTagSafety("v6.530.0", "6.529.0", "v6.528.0", runs, now, 15*time.Minute)
	if status != StatusProtectedGrace {
		t.Fatalf("expected StatusProtectedGrace, got %s", status)
	}
	if eligible {
		t.Fatalf("expected eligible=false during grace window")
	}
}

func TestAuditReleaseTags_CompleteScenario(t *testing.T) {
	mockExec := newTestMockExecutor()

	mockExec.responses["git tag -l"] = []byte("v6.529.0\nv6.528.0\nv6.525.0\nv6.520.0\n")
	mockExec.responses["git ls-remote --tags origin"] = []byte(
		"sha1\trefs/tags/v6.529.0\nsha2\trefs/tags/v6.528.0\nsha3\trefs/tags/v6.520.0\n",
	)

	ghReleasesJSON := `[
		{
			"tagName": "v6.529.0",
			"isDraft": false,
			"createdAt": "2026-10-01T00:00:00Z",
			"assets": [
				{"name": "gitmap-win.zip", "size": 1000, "state": "uploaded"},
				{"name": "checksums.txt", "size": 100, "state": "uploaded"}
			]
		},
		{
			"tagName": "v6.528.0",
			"isDraft": false,
			"createdAt": "2026-09-01T00:00:00Z",
			"assets": [
				{"name": "gitmap-win.zip", "size": 1000, "state": "uploaded"},
				{"name": "checksums.txt", "size": 100, "state": "uploaded"}
			]
		},
		{
			"tagName": "v6.525.0",
			"isDraft": true,
			"createdAt": "2026-08-01T00:00:00Z",
			"assets": []
		}
	]`
	mockExec.responses["gh release list --limit 100 --json tagName,isDraft,createdAt,assets"] = []byte(ghReleasesJSON)

	opts := AuditFilterOptions{
		ActiveVersionOverride: "6.529.0",
		GracePeriod:           15 * time.Minute,
		CommandExecutor:       mockExec,
	}

	report, err := AuditReleaseTags(".", opts)
	if err != nil {
		t.Fatalf("AuditReleaseTags failed: %v", err)
	}

	if report.LatestHealthyTag != "v6.528.0" && report.LatestHealthyTag != "v6.529.0" {
		t.Fatalf("expected healthy tag, got %s", report.LatestHealthyTag)
	}

	var orphanFound, draftFound bool
	for _, rec := range report.Records {
		if rec.Tag == "v6.520.0" && rec.AuditReason == ReasonOrphanTag {
			orphanFound = true
			if !rec.IsEligibleForDeletion {
				t.Errorf("orphan tag should be eligible for deletion")
			}
		}
		if rec.Tag == "v6.525.0" && rec.AuditReason == ReasonDraftRelease {
			draftFound = true
			if !rec.IsEligibleForDeletion {
				t.Errorf("draft release should be eligible for deletion")
			}
		}
	}

	if !orphanFound {
		t.Errorf("expected v6.520.0 orphan tag detection")
	}
	if !draftFound {
		t.Errorf("expected v6.525.0 draft release detection")
	}
}

func TestSemverComparisonAndSorting(t *testing.T) {
	tags := []string{"v1.0.0", "v2.5.0", "v2.10.1", "v0.9.0"}
	SortTagsDescending(tags)

	expected := []string{"v2.10.1", "v2.5.0", "v1.0.0", "v0.9.0"}
	for i, exp := range expected {
		if tags[i] != exp {
			t.Errorf("at index %d: expected %s, got %s", i, exp, tags[i])
		}
	}

	if CompareSemver("v6.529.0", "6.528.0") != 1 {
		t.Errorf("expected v6.529.0 > 6.528.0")
	}
	if CompareSemver("6.528.0", "v6.529.0") != -1 {
		t.Errorf("expected 6.528.0 < v6.529.0")
	}
	if CompareSemver("v6.529.0", "v6.529.0") != 0 {
		t.Errorf("expected equality to yield 0")
	}
}

func TestInspectAssets(t *testing.T) {
	raw := []ghAssetItem{
		{Name: "gitmap.zip", Size: 500, State: "uploaded"},
		{Name: "checksums.txt", Size: 64, State: "uploaded"},
		{Name: "empty.tar.gz", Size: 0, State: "uploaded"},
	}

	assets, hasChecksum, hasBinary, hasCorrupt := InspectAssets(raw)
	if !hasChecksum {
		t.Errorf("expected hasChecksum=true")
	}
	if !hasBinary {
		t.Errorf("expected hasBinary=true")
	}
	if !hasCorrupt {
		t.Errorf("expected hasCorrupt=true for 0-byte asset")
	}
	if len(assets) != 3 {
		t.Errorf("expected 3 assets, got %d", len(assets))
	}
}
