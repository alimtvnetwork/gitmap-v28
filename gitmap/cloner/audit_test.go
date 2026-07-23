package cloner

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// TestPlanCloneAudit_ClassifiesEveryRecordType drives the full planner
// against a JSON manifest with one record of each expected outcome:
// fresh clone, conflict, and invalid (no URL). The pull/cached cases
// require a real .git directory which is exercised separately below.
func TestPlanCloneAudit_ClassifiesEveryRecordType(t *testing.T) {
	target := t.TempDir()
	// Conflict case: a non-git directory at the would-be dest.
	mustMkdir(t, filepath.Join(target, "conflict-repo"))

	records := []model.ScanRecord{
		{
			RepoName: "fresh", RelativePath: "fresh", HTTPSUrl: "https://x/fresh.git",
			BranchSource: gitutil.BranchSourceHEAD, Branch: "main",
		},
		{
			RepoName: "conflict-repo", RelativePath: "conflict-repo", HTTPSUrl: "https://x/conflict.git",
			BranchSource: gitutil.BranchSourceUnknown,
		},
		{
			RepoName: "no-url", RelativePath: "no-url",
		},
	}
	source := writeJSONManifest(t, records)

	report, err := PlanCloneAudit(source, target)
	if err != nil {
		t.Fatalf("PlanCloneAudit: %v", err)
	}
	if len(report.Entries) != 3 {
		t.Fatalf("entries = %d, want 3", len(report.Entries))
	}

	wantActions := []AuditAction{AuditActionClone, AuditActionConflict, AuditActionInvalid}
	for i, w := range wantActions {
		if got := report.Entries[i].Action; got != w {
			t.Errorf("entry[%d].Action = %q, want %q", i, got, w)
		}
	}
	if report.Counts[AuditActionClone] != 1 ||
		report.Counts[AuditActionConflict] != 1 ||
		report.Counts[AuditActionInvalid] != 1 {
		t.Errorf("counts = %+v, want 1/1/1 for clone/conflict/invalid", report.Counts)
	}
}

// TestPlanCloneAudit_PullsExistingRepo confirms an existing .git
// directory at the target flips the action to "pull" (not clone).
// Cache is empty so IsUpToDate returns false → pull, not cached.
func TestPlanCloneAudit_PullsExistingRepo(t *testing.T) {
	target := t.TempDir()
	repoPath := filepath.Join(target, "live")
	mustMkdir(t, filepath.Join(repoPath, ".git"))

	records := []model.ScanRecord{{
		RepoName: "live", RelativePath: "live", HTTPSUrl: "https://x/live.git",
		BranchSource: gitutil.BranchSourceHEAD, Branch: "main",
	}}
	source := writeJSONManifest(t, records)

	report, err := PlanCloneAudit(source, target)
	if err != nil {
		t.Fatalf("PlanCloneAudit: %v", err)
	}
	if got := report.Entries[0].Action; got != AuditActionPull {
		t.Fatalf("Action = %q, want %q", got, AuditActionPull)
	}
}

// TestBuildAuditCommand_BranchOptIn locks the rendered git command to
// the live cloner's argv shape: `git clone -b <branch> <url> <dest>`
// when useBranch, otherwise `git clone <url> <dest>`. Drift here would
// silently make the audit lie about what the live path runs.
func TestBuildAuditCommand_BranchOptIn(t *testing.T) {
	withBranch := buildAuditCommand("https://x/r.git", "/d", cloneStrategy{useBranch: true, branch: "main"})
	if !strings.Contains(withBranch, "-b main https://x/r.git /d") {
		t.Errorf("withBranch = %q, missing branch flag form", withBranch)
	}
	noBranch := buildAuditCommand("https://x/r.git", "/d", cloneStrategy{})
	if strings.Contains(noBranch, "-b ") {
		t.Errorf("noBranch = %q, must not contain -b", noBranch)
	}
	if !strings.Contains(noBranch, "https://x/r.git /d") {
		t.Errorf("noBranch = %q, missing url+dest", noBranch)
	}
}

// TestCloneAuditReport_PrintFormatsRows verifies the diff-style output
// contract: header, one row per entry with the right marker, indented
// command line for actionable rows, and a final summary line. We grep
// for stable substrings rather than asserting on full byte equality so
// future cosmetic tweaks don't churn this test.
func TestCloneAuditReport_PrintFormatsRows(t *testing.T) {
	report := &CloneAuditReport{
		Source: "manifest.json", Target: "/t",
		Entries: []AuditEntry{
			{Action: AuditActionClone, RelativePath: "fresh", URL: "https://x/a.git", Reason: "r1", Command: "git clone https://x/a.git /t/fresh"},
			{Action: AuditActionInvalid, RelativePath: "bad", Reason: "no url"},
		},
		Counts: map[AuditAction]int{AuditActionClone: 1, AuditActionInvalid: 1},
	}
	var buf bytes.Buffer
	if err := report.Print(&buf); err != nil {
		t.Fatalf("Print: %v", err)
	}
	out := buf.String()
	wantSubs := []string{
		"clone audit:", "manifest.json", "records=2",
		"+ clone fresh", "git clone https://x/a.git /t/fresh",
		"! invalid bad",
		"+clone=1", "!invalid=1",
	}
	for _, s := range wantSubs {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q\nfull:\n%s", s, out)
		}
	}
}

// TestActionMarker_StableMapping freezes the diff-marker contract so
// downstream grep/awk pipelines stay stable across releases.
func TestActionMarker_StableMapping(t *testing.T) {
	cases := map[AuditAction]string{
		AuditActionClone:    "+",
		AuditActionPull:     "~",
		AuditActionCached:   "=",
		AuditActionInvalid:  "!",
		AuditActionConflict: "?",
	}
	for action, want := range cases {
		if got := actionMarker(action); got != want {
			t.Errorf("actionMarker(%q) = %q, want %q", action, got, want)
		}
	}
}

// --- helpers -------------------------------------------------------

func mustMkdir(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", p, err)
	}
}

func writeJSONManifest(t *testing.T, records []model.ScanRecord) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")
	data, err := json.Marshal(records)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}
