package cmdpipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func TestParseNegativeIndex(t *testing.T) {
	val, hasVal := ParseNegativeIndex("-2")
	if !hasVal || val != -2 {
		t.Fatalf("expected -2, true; got %d, %v", val, hasVal)
	}

	_, hasInvalid := ParseNegativeIndex("2")
	if hasInvalid {
		t.Fatalf("expected false for positive number 2")
	}

	_, hasInvalidStr := ParseNegativeIndex("invalid")
	if hasInvalidStr {
		t.Fatalf("expected false for invalid string")
	}
}

func TestNormalizeNegativeIndex(t *testing.T) {
	idx1 := NormalizeNegativeIndex(-1)
	if idx1 != 0 {
		t.Fatalf("expected 0 for -1, got %d", idx1)
	}

	idx2 := NormalizeNegativeIndex(-2)
	if idx2 != 1 {
		t.Fatalf("expected 1 for -2, got %d", idx2)
	}

	idxPos := NormalizeNegativeIndex(5)
	if idxPos != 5 {
		t.Fatalf("expected 5 for 5, got %d", idxPos)
	}
}

func TestParseLastFailuresFlag(t *testing.T) {
	args := []string{"--last-failures", "10"}
	val, hasVal := ParseLastFailuresFlag(args)
	if !hasVal || val != 10 {
		t.Fatalf("expected 10, true; got %d, %v", val, hasVal)
	}

	valEmpty, hasEmpty := ParseLastFailuresFlag([]string{"status"})
	if hasEmpty || valEmpty != 0 {
		t.Fatalf("expected 0, false; got %d, %v", valEmpty, hasEmpty)
	}
}

func TestIsErrorLogsSubcmd(t *testing.T) {
	aliases := []string{
		"error-logs", "errorlogs", "error-log", "errorlog",
		"errors", "err", "errorslogs", "errors-log", "errors-logs",
	}

	for _, a := range aliases {
		if !isErrorLogsSubcmd(a) {
			t.Fatalf("expected isErrorLogsSubcmd(%s) to be true", a)
		}
	}

	if isErrorLogsSubcmd("status") {
		t.Fatalf("expected status to not be an errorlogs subcmd")
	}
}

func TestFormatRelativeDbPath(t *testing.T) {
	rel := FormatRelativeDbPath("/some/path/to/pipeline_test.db")
	if strings.Contains(rel, "\\") && !strings.Contains(rel, "/") {
		t.Fatalf("expected forward slashes in relative db path, got %s", rel)
	}

	defRel := FormatRelativeDbPath("")
	if defRel == "" {
		t.Fatalf("expected non-empty fallback for empty db path")
	}
}

func TestInspectPositionalRunPassing(t *testing.T) {
	runs := []ghRunItem{
		{
			DatabaseId: 1001,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			HeadBranch: "main",
			HeadSha:    "abc1234",
			CreatedAt:  "2026-09-11T12:00:00Z",
			UpdatedAt:  "2026-09-11T12:02:00Z",
		},
	}

	err := InspectPositionalRun("repo/test", runs, -1, false)
	if err != nil {
		t.Fatalf("expected nil err for InspectPositionalRun, got %v", err)
	}
}

func TestInspectPositionalRunFailing(t *testing.T) {
	runs := []ghRunItem{
		{
			DatabaseId: 1002,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "failure",
			HeadBranch: "main",
			HeadSha:    "def5678",
			CreatedAt:  "2026-09-11T12:00:00Z",
			UpdatedAt:  "2026-09-11T12:03:00Z",
		},
	}

	err := InspectPositionalRun("repo/test", runs, -1, false)
	if err != nil {
		t.Fatalf("expected nil err for InspectPositionalRun failing, got %v", err)
	}
}

func TestRenderHistorySummaryTable(t *testing.T) {
	runs := []ghRunItem{
		{
			DatabaseId: 2001,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			HeadBranch: "main",
			HeadSha:    "1111111",
			CreatedAt:  "2026-09-11T10:00:00Z",
			UpdatedAt:  "2026-09-11T10:02:00Z",
		},
		{
			DatabaseId: 2002,
			Name:       "Release",
			Status:     "completed",
			Conclusion: "failure",
			HeadBranch: "main",
			HeadSha:    "2222222",
			CreatedAt:  "2026-09-11T11:00:00Z",
			UpdatedAt:  "2026-09-11T11:05:00Z",
		},
	}

	RenderHistorySummaryTable(runs)
}

func TestRenderLastCachedFailuresEmpty(t *testing.T) {
	err := RenderLastCachedFailures("test/nonexistent-repo-slug", 5, false)
	if err != nil {
		t.Fatalf("expected nil err for RenderLastCachedFailures, got %v", err)
	}
}

func TestPipelineSyncResultCount(t *testing.T) {
	res := &PipelineSyncResult{
		DownloadedLogs: 2,
		CachedErrors:   []uint64{1, 2, 3, 4},
	}

	cached := calculateAlreadyCachedCount(res)
	if cached != 2 {
		t.Fatalf("expected 2 cached, got %d", cached)
	}
}

func TestPipelineDBOpenAndIndexes(t *testing.T) {
	db, err := pipelinedb.OpenPipelineSplitDb("test-repo-indexer")
	if err != nil {
		t.Fatalf("failed to open test split db: %v", err)
	}

	defer db.Close()

	if db.Path == "" {
		t.Fatalf("expected non-empty db path")
	}
}

func TestSyncAllRunsIntoDbIncremental(t *testing.T) {
	db, err := pipelinedb.OpenPipelineSplitDb("test-incremental-sync-repo")
	if err != nil {
		t.Fatalf("failed to open test split db: %v", err)
	}

	defer db.Close()

	run := ghRunItem{
		DatabaseId: 999901,
		Name:       "CI",
		Status:     "completed",
		Conclusion: "failure",
		CreatedAt:  "2026-09-11T12:00:00Z",
		UpdatedAt:  "2026-09-11T12:05:00Z",
	}

	testPassOne(t, db, run)
	testPassTwo(t, db, run)
}

func testPassOne(t *testing.T, db *pipelinedb.PipelineSplitDb, run ghRunItem) {
	res1 := initSyncResult("test-repo", db.Path, 1)
	syncAllRunsIntoDb(db, "test-repo", []ghRunItem{run}, res1)
	if len(res1.FailedRuns) != 1 {
		t.Fatalf("expected 1 failed run on pass 1, got %d", len(res1.FailedRuns))
	}
}

func testPassTwo(t *testing.T, db *pipelinedb.PipelineSplitDb, run ghRunItem) {
	res2 := initSyncResult("test-repo", db.Path, 1)
	syncAllRunsIntoDb(db, "test-repo", []ghRunItem{run}, res2)
	cached := calculateAlreadyCachedCount(res2)
	if cached < 1 && res2.DownloadedLogs > 0 {
		t.Fatalf("expected download to be skipped on pass 2, downloaded: %d", res2.DownloadedLogs)
	}
}

func TestParsePipelineErrorFlags_Combined(t *testing.T) {
	args := []string{"-3", "--last-failures", "5", "--json"}
	flags := ParsePipelineErrorFlags(args)
	if !flags.HasIndex || flags.Index != -3 {
		t.Fatalf("expected Index=-3, HasIndex=true; got %d, %v", flags.Index, flags.HasIndex)
	}

	if !flags.HasLastFailures || flags.LastFailures != 5 {
		t.Fatalf("expected LastFailures=5, HasLastFailures=true; got %d, %v", flags.LastFailures, flags.HasLastFailures)
	}

	if !flags.IsJSON {
		t.Fatalf("expected IsJSON=true")
	}
}

func TestParsePipelineErrorFlags_LastFailedLogs(t *testing.T) {
	args := []string{"last-failed-logs"}
	flags := ParsePipelineErrorFlags(args)
	if !flags.HasLastFailedLogs {
		t.Fatalf("expected HasLastFailedLogs=true")
	}

	if !flags.HasLastFailures || flags.LastFailures != 20 {
		t.Fatalf("expected LastFailures=20 default, got %d", flags.LastFailures)
	}
}

func TestVisibleLenAndPadRightVisible(t *testing.T) {
	if visibleLen("PASS") != 4 {
		t.Fatalf("expected visibleLen('PASS') == 4")
	}
	ansiPass := constants.ColorGreen + "PASS" + constants.ColorReset
	if visibleLen(ansiPass) != 4 {
		t.Fatalf("expected visibleLen(ansiPass) == 4, got %d", visibleLen(ansiPass))
	}
	padded := padRightVisible(ansiPass, 10)
	if visibleLen(padded) != 10 {
		t.Fatalf("expected visibleLen(padded) == 10, got %d", visibleLen(padded))
	}
}

func TestFormatGroupWorkflowsSummary_FitsWithinLimit(t *testing.T) {
	wfs := []CommitWorkflowItem{
		{Name: "CI", Conclusion: "success"},
		{Name: "Lint", Conclusion: "success"},
	}
	got := formatGroupWorkflowsSummary(wfs, 32)
	want := "CI [PASS], Lint [PASS]"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormatGroupWorkflowsSummary_TruncatesSafely(t *testing.T) {
	wfs := []CommitWorkflowItem{
		{Name: "VeryLongWorkflowNameAlpha", Conclusion: "success"},
		{Name: "VeryLongWorkflowNameBeta", Conclusion: "failure"},
	}
	got := formatGroupWorkflowsSummary(wfs, 32)
	if strings.Contains(got, "[PASS") && !strings.Contains(got, "[PASS]") {
		t.Fatalf("detected severed bracket in workflow summary: %q", got)
	}
	if !strings.Contains(got, "(+1)") {
		t.Fatalf("expected remaining indicator (+1) in summary: %q", got)
	}
}

func TestRecentCommitsSummaryTableAlignment(t *testing.T) {
	g := CommitPipelineGroup{
		HeadSha:         "a1b2c3d4e5",
		HeadBranch:      "main",
		Conclusion:      "success",
		Status:          "completed",
		FailedWorkflows: 0,
		Workflows: []CommitWorkflowItem{
			{Name: "CI", Conclusion: "success"},
		},
	}
	var sb strings.Builder
	printRecentCommitRow(&sb, g, 0)
	plainRow := stripANSI(sb.String())
	parts := strings.Fields(plainRow)
	if len(parts) < 7 {
		t.Fatalf("expected at least 7 columns, got %d in %q", len(parts), plainRow)
	}
	if parts[3] != "-" {
		t.Fatalf("expected Release column to be -, got %s", parts[3])
	}
	if parts[4] != "PASS" {
		t.Fatalf("expected Status column to be PASS, got %s", parts[4])
	}
}

func TestFormatDbPathWithSize_ExistingAndMissing(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test_pipeline.db")
	_ = os.WriteFile(tempFile, []byte("sqlite database test payload"), 0644)
	withSize := FormatDbPathWithSize(tempFile)
	if !strings.Contains(withSize, " B)") && !strings.Contains(withSize, " KB)") {
		t.Fatalf("expected file size in formatted path, got %s", withSize)
	}
	missingPath := filepath.Join(t.TempDir(), "non_existent.db")
	withoutSize := FormatDbPathWithSize(missingPath)
	if strings.Contains(withoutSize, "(") {
		t.Fatalf("expected no size for missing db file, got %s", withoutSize)
	}
}

func TestResolveDbFileSize(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test_size.db")
	_ = os.WriteFile(tempFile, make([]byte, 2048), 0644)
	size := ResolveDbFileSize(tempFile)
	if !strings.Contains(size, "KB") && !strings.Contains(size, "B") {
		t.Fatalf("expected valid size string, got %s", size)
	}

	missing := filepath.Join(t.TempDir(), "missing.db")
	missingSize := ResolveDbFileSize(missing)
	if missingSize != "0 B" {
		t.Fatalf("expected '0 B' for missing file, got %s", missingSize)
	}
}

func TestFormatRelativeDbPath_RepoScoped(t *testing.T) {
	p1 := filepath.Join("project", "repo", ".gitmap", "data", "pipeline.db")
	if rel := FormatRelativeDbPath(p1); rel != ".gitmap/data/pipeline.db" {
		t.Errorf("expected .gitmap/data/pipeline.db, got %s", rel)
	}

	p2 := filepath.Join("home", "user", "projects", "repo", ".gitmap", "pipeline.db")
	if rel := FormatRelativeDbPath(p2); rel != ".gitmap/pipeline.db" {
		t.Errorf("expected .gitmap/pipeline.db, got %s", rel)
	}
}

func TestResolveCommitRelease(t *testing.T) {
	gReleaseBranch := CommitPipelineGroup{HeadBranch: "release/v6.290.0"}
	if rel := resolveCommitRelease(gReleaseBranch); rel != "v6.290.0" {
		t.Errorf("expected v6.290.0 from release branch, got: %s", rel)
	}

	gTagBranch := CommitPipelineGroup{HeadBranch: "v6.289.0"}
	if rel := resolveCommitRelease(gTagBranch); rel != "v6.289.0" {
		t.Errorf("expected v6.289.0 from tag branch, got: %s", rel)
	}

	gNoRelease := CommitPipelineGroup{HeadBranch: "feature/login", HeadSha: "0000000000"}
	if rel := resolveCommitRelease(gNoRelease); rel != "-" {
		t.Errorf("expected '-' for unreleased commit, got: %s", rel)
	}
}
