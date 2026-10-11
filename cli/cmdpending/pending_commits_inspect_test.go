// Package cmdpending — pending_commits_inspect_test.go verifies status parsing, repo inspection, and formatting.
package cmdpending

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

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

// parsePorcelainStatusLines parses `git status --porcelain` output and
// returns (untracked, modified, staged, files). Test-local helper.
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

func TestFormatShortVersionBranch(t *testing.T) {
	verifyFormatShortVersionBranchStandard(t)
	verifyFormatShortVersionBranchFallback(t)
	verifyFormatShortVersionBranchTruncation(t)
}

func verifyFormatShortVersionBranchStandard(t *testing.T) {
	result := formatShortVersionBranch("v1.2.3", "main")
	if result != "v1.2.3/main" {
		t.Errorf("expected 'v1.2.3/main', got '%s'", result)
	}
}

func verifyFormatShortVersionBranchFallback(t *testing.T) {
	result := formatShortVersionBranch("", "feature-login")
	if result != "feature-login" {
		t.Errorf("expected 'feature-login', got '%s'", result)
	}
}

func verifyFormatShortVersionBranchTruncation(t *testing.T) {
	result := formatShortVersionBranch("v6.523.1", "feature-long-name-here")
	if len([]rune(result)) > 16 {
		t.Errorf("expected rune length <= 16, got %d ('%s')", len([]rune(result)), result)
	}
	if !strings.HasSuffix(result, "…") {
		t.Errorf("expected trailing ellipsis, got '%s'", result)
	}
}

func TestBuildRepoRemediationOptions_Dirty(t *testing.T) {
	rec := RepoPendingCommitRecord{
		RepoName:     "gitmap",
		RelativePath: "gitmap",
		IsDirty:      true,
	}
	options := buildRepoRemediationOptions(rec)
	if len(options) != 2 {
		t.Fatalf("expected 2 remediation options for dirty repo, got %d", len(options))
	}
	verifyRemediationOptionsContent(t, options)
}

func verifyRemediationOptionsContent(t *testing.T, options []RemediationOption) {
	if options[0].OptionNumber != 1 {
		t.Errorf("expected OptionNumber 1, got %d", options[0].OptionNumber)
	}
	if !strings.Contains(options[0].Command, "gitmap sends cp") {
		t.Errorf("expected gitmap sends cp command in Option 1, got %s", options[0].Command)
	}
	if options[1].OptionNumber != 2 {
		t.Errorf("expected OptionNumber 2, got %d", options[1].OptionNumber)
	}
	if !strings.Contains(options[1].Command, "gitmap stash") {
		t.Errorf("expected gitmap stash command in Option 2, got %s", options[1].Command)
	}
}

func TestBuildRepoRemediationOptions_Clean(t *testing.T) {
	rec := RepoPendingCommitRecord{
		RepoName:     "gitmap",
		RelativePath: "gitmap",
		IsDirty:      false,
		IsClean:      true,
	}
	options := buildRepoRemediationOptions(rec)
	if len(options) != 0 {
		t.Errorf("expected 0 options for clean repo, got %d", len(options))
	}
}

func TestTotalUncommittedCalculation(t *testing.T) {
	rec := assembleRepoPendingCommitRecord(
		model.ScanRecord{RepoName: "test-repo"},
		"main", "v1.0.0", "v1.0.0/main",
		3, 5, 2, 10, 0,
		true, nil, nil,
	)
	if rec.TotalUncommitted != 10 {
		t.Errorf("expected TotalUncommitted 10, got %d", rec.TotalUncommitted)
	}
	expectedSum := rec.UntrackedFilesCount + rec.ModifiedFilesCount + rec.StagedFilesCount
	if rec.TotalUncommitted != expectedSum {
		t.Errorf("expected TotalUncommitted %d to match sum %d", rec.TotalUncommitted, expectedSum)
	}
}

func TestPendingCommitsJSON_NewFields(t *testing.T) {
	rec := RepoPendingCommitRecord{
		RepoName:           "gitmap",
		RelativePath:       ".",
		CurrentBranch:      "main",
		Version:            "v6.523.1",
		ShortVersionBranch: "v6.523.1/main",
		IsDirty:            true,
		TotalUncommitted:   4,
		RemediationOptions: []RemediationOption{
			{OptionNumber: 1, Label: "Commit & Push WIP", Command: "git add -A"},
		},
	}
	verifyNewFieldsJSONSerialization(t, rec)
}

func verifyNewFieldsJSONSerialization(t *testing.T, rec RepoPendingCommitRecord) {
	data, errMarshal := json.Marshal(rec)
	if errMarshal != nil {
		t.Fatalf("failed to marshal JSON: %v", errMarshal)
	}
	var parsed map[string]any
	if errUnmarshal := json.Unmarshal(data, &parsed); errUnmarshal != nil {
		t.Fatalf("failed to unmarshal JSON: %v", errUnmarshal)
	}
	checkJSONKeysPresent(t, parsed)
}

func checkJSONKeysPresent(t *testing.T, parsed map[string]any) {
	requiredKeys := []string{"totalUncommitted", "version", "shortVersionBranch", "remediationOptions"}
	for _, key := range requiredKeys {
		if _, hasKey := parsed[key]; !hasKey {
			t.Errorf("expected JSON key '%s' to be present", key)
		}
	}
}
