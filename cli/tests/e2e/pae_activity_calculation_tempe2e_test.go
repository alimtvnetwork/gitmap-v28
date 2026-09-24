//go:build tempe2e

package e2e_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
)

func TestTempE2E_PAEActivityCalculationAndStability(t *testing.T) {
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("skipping temporary E2E test: RUN_TEMP_E2E=1 not set")
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-gitmap-pull-e2e.db")
	db, err := store.OpenPullSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenPullSplitDBAt failed: %v", err)
	}
	defer db.Close()

	// 1. Simulate 4 repos pulled in Run 1 (all returning up-to-date)
	repoA := filepath.Join(tempDir, "Antigravity-Manager")
	repoB := filepath.Join(tempDir, "gitmap-v28")
	repoC := filepath.Join(tempDir, "gstack")
	repoD := filepath.Join(tempDir, "riseup-asia-website-project-v6")

	runID, err := db.InsertPullRun(&store.PullRunRecord{
		CommandType: "pull all-efficient",
		TotalRepos:  4,
		PulledRepos: 4,
	})
	if err != nil {
		t.Fatalf("InsertPullRun failed: %v", err)
	}

	err = db.InsertPullRepoRuns(runID, []store.PullRepoRunRecord{
		{RepoPath: repoA, RepoName: "Antigravity-Manager", PullStatus: "up-to-date", FilesChanged: 0, HasChanges: false},
		{RepoPath: repoB, RepoName: "gitmap-v28", PullStatus: "up-to-date", FilesChanged: 0, HasChanges: false},
		{RepoPath: repoC, RepoName: "gstack", PullStatus: "up-to-date", FilesChanged: 0, HasChanges: false},
		{RepoPath: repoD, RepoName: "riseup-asia-website-project-v6", PullStatus: "up-to-date", FilesChanged: 0, HasChanges: false},
	})
	if err != nil {
		t.Fatalf("InsertPullRepoRuns failed: %v", err)
	}

	// 2. Simulate Run 2 (invoked seconds later)
	// Under the new logic, all 4 are in freshness cooldown (< 5m) -> all 4 must be inactive!
	repos := []string{repoA, repoB, repoC, repoD}
	for _, r := range repos {
		st, err := db.EvaluateRepoActivityStatus(r, 24, 5)
		if err != nil {
			t.Fatalf("EvaluateRepoActivityStatus failed for %s: %v", r, err)
		}
		if !st.IsInactive {
			t.Fatalf("expected repo %s to be inactive due to cooldown, got active (reason: %s)", r, st.Reason)
		}
	}

	// 3. Test layout stability with allRecords
	allTracked := []model.ScanRecord{
		{RepoName: "Antigravity-Manager"},
		{RepoName: "gitmap-v28"},
		{RepoName: "gstack"},
		{RepoName: "riseup-asia-website-project-v6"}, // 30 chars
	}
	activeSubset := []*cmdpull.PullRepoState{
		{RepoName: "Antigravity-Manager", Changes: "up-to-date"},
		{RepoName: "gitmap-v28", Changes: "up-to-date"},
	}

	var buf bytes.Buffer
	cmdpull.RenderConciseActiveResultsTo(&buf, activeSubset, allTracked)

	rawLines := strings.Split(buf.String(), "\n")
	var lines []string
	for _, l := range rawLines {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d:\n%s", len(lines), buf.String())
	}

	plain0 := termpad.StripAnsi(lines[0])
	plain1 := termpad.StripAnsi(lines[1])
	idx0 := strings.LastIndex(plain0, "up-to-date")
	idx1 := strings.LastIndex(plain1, "up-to-date")

	if idx0 != idx1 {
		t.Fatalf("column alignment mismatch: idx0=%d, idx1=%d\nLine 0: %s\nLine 1: %s", idx0, idx1, plain0, plain1)
	}
	if idx0 < 30 {
		t.Fatalf("expected status column >= 30 due to allTracked, got %d", idx0)
	}
}
