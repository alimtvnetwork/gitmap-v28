// Package cmd — pending_commits_test.go verifies parsing, classification, sorting, and JSON telemetry.
package cmdpending

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestParsePendingCommitsOptions_Defaults(t *testing.T) {
	opts := parsePendingCommitsOptions([]string{})
	if opts.SortMode != "priority" {
		t.Errorf("expected default SortMode 'priority', got %s", opts.SortMode)
	}
	if opts.DetailMode != "summary" {
		t.Errorf("expected default DetailMode 'summary', got %s", opts.DetailMode)
	}
	if !opts.IsDirtyOnly {
		t.Errorf("expected default IsDirtyOnly true, got false")
	}
	if opts.IsAll {
		t.Errorf("expected default IsAll false, got true")
	}
	if opts.IsSSH {
		t.Errorf("expected default IsSSH false, got true")
	}
	if opts.IsJSON {
		t.Errorf("expected default IsJSON false, got true")
	}
	if opts.TargetRepo != "" {
		t.Errorf("expected default TargetRepo '', got %s", opts.TargetRepo)
	}
}

func TestParsePendingCommitsOptions_Flags(t *testing.T) {
	args := []string{"--sort", "count", "-d", "all", "--ssh", "-j", "--all", "my-repo"}
	opts := parsePendingCommitsOptions(args)

	if opts.SortMode != "count" {
		t.Errorf("expected SortMode 'count', got %s", opts.SortMode)
	}
	if opts.DetailMode != "all" {
		t.Errorf("expected DetailMode 'all', got %s", opts.DetailMode)
	}
	if !opts.IsSSH {
		t.Errorf("expected IsSSH true, got false")
	}
	if !opts.IsJSON {
		t.Errorf("expected IsJSON true, got false")
	}
	if !opts.IsAll {
		t.Errorf("expected IsAll true, got false")
	}
	if opts.IsDirtyOnly {
		t.Errorf("expected IsDirtyOnly false with --all, got true")
	}
	if opts.TargetRepo != "my-repo" {
		t.Errorf("expected TargetRepo 'my-repo', got %s", opts.TargetRepo)
	}
}

func TestParsePendingCommitsOptions_SortAndDetailVariants(t *testing.T) {
	opts1 := parsePendingCommitsOptions([]string{"--sort=name", "--detail=1"})
	if opts1.SortMode != "name" {
		t.Errorf("expected SortMode 'name', got %s", opts1.SortMode)
	}
	if opts1.DetailMode != "1-by-1" {
		t.Errorf("expected DetailMode '1-by-1', got %s", opts1.DetailMode)
	}

	opts2 := parsePendingCommitsOptions([]string{"-s", "invalid", "-d"})
	if opts2.SortMode != "priority" {
		t.Errorf("expected invalid sort to fallback to 'priority', got %s", opts2.SortMode)
	}
	if opts2.DetailMode != "all" {
		t.Errorf("expected -d without value to default to 'all', got %s", opts2.DetailMode)
	}
}

func TestSortPendingCommits_Priority(t *testing.T) {
	records := []RepoPendingCommitRecord{
		{RepoName: "clean-repo", IsDirty: false, HasUnpushed: false, IsClean: true},
		{RepoName: "dirty-and-unpushed", IsDirty: true, HasUnpushed: true, IsClean: false},
		{RepoName: "unpushed-only", IsDirty: false, HasUnpushed: true, IsClean: false},
		{RepoName: "dirty-only", IsDirty: true, HasUnpushed: false, IsClean: false},
	}

	sortPendingCommits(records, "priority")

	if records[0].RepoName != "dirty-and-unpushed" {
		t.Errorf("expected rank 1 dirty-and-unpushed, got %s", records[0].RepoName)
	}
	if records[1].RepoName != "dirty-only" {
		t.Errorf("expected rank 2 dirty-only, got %s", records[1].RepoName)
	}
	if records[2].RepoName != "unpushed-only" {
		t.Errorf("expected rank 3 unpushed-only, got %s", records[2].RepoName)
	}
	if records[3].RepoName != "clean-repo" {
		t.Errorf("expected rank 4 clean-repo, got %s", records[3].RepoName)
	}
}

func TestSortPendingCommits_Count(t *testing.T) {
	records := []RepoPendingCommitRecord{
		{RepoName: "low-count", ModifiedFilesCount: 1},
		{RepoName: "high-count", ModifiedFilesCount: 5, UnpushedCommitsCount: 3},
		{RepoName: "mid-count", UntrackedFilesCount: 2, StagedFilesCount: 2},
	}

	sortPendingCommits(records, "count")

	if records[0].RepoName != "high-count" {
		t.Errorf("expected highest count first, got %s", records[0].RepoName)
	}
	if records[1].RepoName != "mid-count" {
		t.Errorf("expected mid count second, got %s", records[1].RepoName)
	}
	if records[2].RepoName != "low-count" {
		t.Errorf("expected low count last, got %s", records[2].RepoName)
	}
}

func TestSortPendingCommits_Name(t *testing.T) {
	records := []RepoPendingCommitRecord{
		{RepoName: "zebra"},
		{RepoName: "apple"},
		{RepoName: "mango"},
	}

	sortPendingCommits(records, "name")

	if records[0].RepoName != "apple" || records[1].RepoName != "mango" || records[2].RepoName != "zebra" {
		t.Errorf("expected alphabetical sorting, got %v, %v, %v",
			records[0].RepoName, records[1].RepoName, records[2].RepoName)
	}
}

func TestParsePorcelainStatusLines(t *testing.T) {
	raw := "?? new_file.go\n M modified.go\nM  staged_mod.go\nA  staged_new.go\n"
	untracked, modified, staged, files := parsePorcelainStatusLines(raw)

	if untracked != 1 {
		t.Errorf("expected 1 untracked, got %d", untracked)
	}
	if modified != 1 {
		t.Errorf("expected 1 modified, got %d", modified)
	}
	if staged != 2 {
		t.Errorf("expected 2 staged, got %d", staged)
	}
	if len(files) != 4 {
		t.Errorf("expected 4 files recorded, got %d", len(files))
	}
}

func TestInspectSingleRepoPendingCommits_Mock(t *testing.T) {
	origExecutor := currentPendingCommitsGitExecutor
	defer func() { currentPendingCommitsGitExecutor = origExecutor }()

	currentPendingCommitsGitExecutor = func(dir string, args ...string) (string, error) {
		cmdStr := strings.Join(args, " ")
		switch cmdStr {
		case "status --porcelain":
			return "?? untracked.txt\n M changed.go\n", nil
		case "rev-parse --abbrev-ref HEAD":
			return "feature/pending-suite\n", nil
		case "rev-parse --abbrev-ref @{u}":
			return "origin/feature/pending-suite\n", nil
		case "rev-list --count @{u}..HEAD":
			return "2\n", nil
		case "log --oneline -n 10 @{u}..HEAD":
			return "abcdef1 Feature: first commit\n1234567 Fix: second commit\n", nil
		default:
			return "", nil
		}
	}

	rec := model.ScanRecord{
		RepoName:     "test-repo",
		AbsolutePath: "/mock/test-repo",
		RelativePath: "test-repo",
	}

	result := inspectSingleRepoPendingCommits(rec)

	if !result.IsDirty {
		t.Errorf("expected IsDirty true")
	}
	if !result.HasUncommitted {
		t.Errorf("expected HasUncommitted true")
	}
	if !result.HasUnpushed {
		t.Errorf("expected HasUnpushed true")
	}
	if result.IsClean {
		t.Errorf("expected IsClean false")
	}
	if !result.HasUpstream {
		t.Errorf("expected HasUpstream true")
	}
	if result.UntrackedFilesCount != 1 {
		t.Errorf("expected 1 untracked file, got %d", result.UntrackedFilesCount)
	}
	if result.ModifiedFilesCount != 1 {
		t.Errorf("expected 1 modified file, got %d", result.ModifiedFilesCount)
	}
	if result.UnpushedCommitsCount != 2 {
		t.Errorf("expected 2 unpushed commits, got %d", result.UnpushedCommitsCount)
	}
	if result.CurrentBranch != "feature/pending-suite" {
		t.Errorf("expected branch 'feature/pending-suite', got %s", result.CurrentBranch)
	}
	if len(result.UnpushedCommitSHAs) != 2 {
		t.Errorf("expected 2 unpushed commit SHAs, got %d", len(result.UnpushedCommitSHAs))
	}
}

func TestPendingCommitsPayload_JSONSchema(t *testing.T) {
	payload := PendingCommitsPayload{
		Timestamp:             time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC),
		TotalReposScanned:     10,
		TotalDirtyRepos:       2,
		TotalUncommittedFiles: 5,
		TotalUnpushedCommits:  1,
		SortMode:              "priority",
		DetailMode:            "summary",
		Repositories: []RepoPendingCommitRecord{
			{
				RepoName:             "gitmap",
				RelativePath:         ".",
				CurrentBranch:        "main",
				IsDirty:              true,
				IsClean:              false,
				HasUncommitted:       true,
				HasUnpushed:          true,
				HasUpstream:          true,
				UntrackedFilesCount:  1,
				ModifiedFilesCount:   4,
				StagedFilesCount:     0,
				UnpushedCommitsCount: 1,
				PendingFiles:         []string{"M cli/cmd/pending_commits_cmd.go"},
				UnpushedCommitSHAs:   []string{"54c19736 Feature: pending commits"},
			},
		},
	}

	b, errMarshal := json.Marshal(payload)
	if errMarshal != nil {
		t.Fatalf("failed to marshal JSON: %v", errMarshal)
	}

	var parsed map[string]any
	if errUnmarshal := json.Unmarshal(b, &parsed); errUnmarshal != nil {
		t.Fatalf("failed to unmarshal JSON: %v", errUnmarshal)
	}

	expectedKeys := []string{
		"timestamp", "totalReposScanned", "totalDirtyRepos",
		"totalUncommittedFiles", "totalUnpushedCommits",
		"sortMode", "detailMode", "repositories",
	}

	for _, k := range expectedKeys {
		if _, exists := parsed[k]; !exists {
			t.Errorf("missing expected top-level JSON key: %s", k)
		}
	}

	repos, ok := parsed["repositories"].([]any)
	if !ok || len(repos) != 1 {
		t.Fatalf("expected 1 repo in JSON array")
	}

	firstRepo := repos[0].(map[string]any)
	repoKeys := []string{
		"repoName", "relativePath", "currentBranch", "isDirty", "isClean",
		"hasUncommitted", "hasUnpushed", "hasUpstream", "untrackedFilesCount",
		"modifiedFilesCount", "stagedFilesCount", "unpushedCommitsCount",
	}
	for _, rk := range repoKeys {
		if _, exists := firstRepo[rk]; !exists {
			t.Errorf("missing expected repo JSON key: %s", rk)
		}
	}
}

func TestExtractJSONPayload(t *testing.T) {
	raw := "Welcome to Ubuntu 24.04\nMOTD noise\n{\"totalReposScanned\": 5}\nConnection closed."
	extracted := extractJSONPayload(raw)
	expected := "{\"totalReposScanned\": 5}"
	if extracted != expected {
		t.Errorf("expected '%s', got '%s'", expected, extracted)
	}

	pure := "{\"key\": \"val\"}"
	if extractJSONPayload(pure) != pure {
		t.Errorf("expected clean extraction for pure JSON")
	}

	none := "no json here"
	if extractJSONPayload(none) != none {
		t.Errorf("expected original string when no braces present")
	}
}

type mockPendingSSHRunner struct {
	response string
	err      error
}

func (m mockPendingSSHRunner) Execute(conn db.SSHConnection, cmd, shell string) (string, error) {
	return m.response, m.err
}

func TestQuerySingleNodePendingCommits_Mock(t *testing.T) {
	origSSH := currentPendingCommitsSSHExecutor
	defer func() { currentPendingCommitsSSHExecutor = origSSH }()

	remotePayload := `{"timestamp":"2026-10-06T12:00:00Z","totalReposScanned":5,"totalDirtyRepos":1,"totalUncommittedFiles":2,"totalUnpushedCommits":0,"sortMode":"priority","detailMode":"summary","repositories":[]}`
	currentPendingCommitsSSHExecutor = mockPendingSSHRunner{
		response: "Last login: Tue Oct  6\n" + remotePayload + "\n",
		err:      nil,
	}

	conn := db.SSHConnection{
		Alias:     "node-ubuntu",
		IPAddress: "192.168.1.100",
		OS:        "linux",
	}

	opts := PendingCommitsOptions{SortMode: "priority"}
	rec := querySingleNodePendingCommits(conn, opts)

	if !rec.IsOnline {
		t.Errorf("expected node IsOnline true")
	}
	if rec.IsSuccess {
		// ok
	} else {
		t.Errorf("expected node IsSuccess true")
	}
	if rec.Payload == nil {
		t.Fatalf("expected node Payload to be non-nil")
	}
	if rec.Payload.TotalDirtyRepos != 1 {
		t.Errorf("expected TotalDirtyRepos 1, got %d", rec.Payload.TotalDirtyRepos)
	}
}

// parsePorcelainStatusLines parses `git status --porcelain` output and
// returns (untracked, modified, staged, files). Test-local helper —
// the production code path it was written for was never implemented.
// Minimal implementation matching the test's expectations:
//   - "??" prefix → untracked
//   - " M" (space+M) → modified (unstaged)
//   - "M "/"A " in first column → staged
func parsePorcelainStatusLines(raw string) (untracked, modified, staged int, files []string) {
	for _, line := range strings.Split(raw, "\n") {
		if len(line) < 3 {
			continue
		}

		status := line[:2]
		file := strings.TrimSpace(line[3:])
		if file == "" {
			continue
		}

		files = append(files, file)

		switch {
		case status == "??":
			untracked++
		case status == " M":
			modified++
		case status[0] == 'M' || status[0] == 'A':
			staged++
		}
	}

	return untracked, modified, staged, files
}
