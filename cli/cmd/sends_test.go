// Package cmd — sends_test.go verifies semantic verb parsing, prefix normalization, dry-run safety, and clean repo skipping.
package cmd

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestParseSendsArguments_SemanticVerbs(t *testing.T) {
	testCases := []struct {
		args           []string
		expectedVerb   string
		expectedPfx    string
		expectedMsg    string
		expectedPushed bool
	}{
		{
			args:           []string{"cpf", "all", "add feature X"},
			expectedVerb:   "cpf",
			expectedPfx:    "Feature: ",
			expectedMsg:    "add feature X",
			expectedPushed: true,
		},
		{
			args:           []string{"cpb", "gitmap", "resolve nil pointer"},
			expectedVerb:   "cpb",
			expectedPfx:    "Bug: ",
			expectedMsg:    "resolve nil pointer",
			expectedPushed: true,
		},
		{
			args:           []string{"cpr", "all", "v1.2.3 release"},
			expectedVerb:   "cpr",
			expectedPfx:    "Release: ",
			expectedMsg:    "v1.2.3 release",
			expectedPushed: true,
		},
		{
			args:           []string{"commit-fix", "my-repo", "merge resolution"},
			expectedVerb:   "commit-fix",
			expectedPfx:    "Fix: ",
			expectedMsg:    "merge resolution",
			expectedPushed: true,
		},
		{
			args:           []string{"cp", "all", "regular commit and push"},
			expectedVerb:   "cp",
			expectedPfx:    "",
			expectedMsg:    "regular commit and push",
			expectedPushed: true,
		},
		{
			args:           []string{"cm", "all", "local commit only"},
			expectedVerb:   "cm",
			expectedPfx:    "",
			expectedMsg:    "local commit only",
			expectedPushed: false,
		},
	}

	for _, tc := range testCases {
		opts, err := parseSendsArguments(tc.args)
		if err != nil {
			t.Fatalf("unexpected error parsing %v: %v", tc.args, err)
		}
		if opts.Verb != tc.expectedVerb {
			t.Errorf("expected verb %s, got %s", tc.expectedVerb, opts.Verb)
		}
		if opts.RawMessage != tc.expectedMsg {
			t.Errorf("expected message '%s', got '%s'", tc.expectedMsg, opts.RawMessage)
		}
		pfx := resolveSemanticCommitPrefix(opts.Verb)
		if pfx != tc.expectedPfx {
			t.Errorf("expected prefix '%s', got '%s'", tc.expectedPfx, pfx)
		}
		if opts.IsPushed != tc.expectedPushed {
			t.Errorf("expected IsPushed %v, got %v", tc.expectedPushed, opts.IsPushed)
		}
	}
}

func TestParseSendsArguments_InvalidUsage(t *testing.T) {
	_, errFewArgs := parseSendsArguments([]string{"cpf", "all"})
	if errFewArgs == nil {
		t.Errorf("expected error for fewer than 3 arguments, got nil")
	}

	_, errUnknownVerb := parseSendsArguments([]string{"invalid-verb", "all", "message"})
	if errUnknownVerb == nil {
		t.Errorf("expected error for unknown verb, got nil")
	}
}

func TestParseSendsArguments_Flags(t *testing.T) {
	args := []string{"cpf", "all", "multi", "word", "message", "-n", "--no-push", "-j", "-v"}
	opts, err := parseSendsArguments(args)
	if err != nil {
		t.Fatalf("unexpected error parsing arguments with flags: %v", err)
	}

	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun true with -n")
	}
	if opts.IsPushed {
		t.Errorf("expected IsPushed false with --no-push")
	}
	if !opts.IsJSON {
		t.Errorf("expected IsJSON true with -j")
	}
	if !opts.IsVerbose {
		t.Errorf("expected IsVerbose true with -v")
	}
	if opts.RawMessage != "multi word message" {
		t.Errorf("expected RawMessage 'multi word message', got '%s'", opts.RawMessage)
	}
}

func TestProcessSingleRepoSend_DirtyRepo_Mock(t *testing.T) {
	origExecutor := currentSendsGitExecutor
	defer func() { currentSendsGitExecutor = origExecutor }()

	executedCommands := []string{}
	currentSendsGitExecutor = func(dir string, args ...string) (string, error) {
		cmdStr := strings.Join(args, " ")
		executedCommands = append(executedCommands, cmdStr)
		switch {
		case cmdStr == "status --porcelain":
			return " M file1.go\n?? file2.go\n", nil
		case cmdStr == "rev-parse --abbrev-ref HEAD":
			return "main\n", nil
		case cmdStr == "rev-parse --short HEAD":
			return "54c19736\n", nil
		default:
			return "", nil
		}
	}

	repo := model.ScanRecord{
		RepoName:     "gitmap",
		AbsolutePath: "/mock/gitmap",
		RelativePath: ".",
	}

	opts := SendsOptions{
		Verb:       "cpf",
		Target:     "all",
		RawMessage: "add pending commits",
		IsDryRun:   false,
		IsPushed:   true,
	}

	finalMsg := "Feature: add pending commits"
	rec := processSingleRepoSend(repo, opts, finalMsg)

	if rec.Status != "pushed" {
		t.Errorf("expected Status 'pushed', got %s", rec.Status)
	}
	if rec.IsSuccess {
	} else {
		t.Errorf("expected IsSuccess true, got false")
	}
	if rec.FilesStaged != 2 {
		t.Errorf("expected FilesStaged 2, got %d", rec.FilesStaged)
	}
	if rec.HeadSHA != "54c19736" {
		t.Errorf("expected HeadSHA '54c19736', got %s", rec.HeadSHA)
	}

	hasAdd := false
	hasCommit := false
	hasPush := false
	for _, c := range executedCommands {
		if strings.HasPrefix(c, "add -A") {
			hasAdd = true
		}
		if strings.HasPrefix(c, "commit -m") {
			hasCommit = true
		}
		if c == "push" {
			hasPush = true
		}
	}

	if !hasAdd || !hasCommit || !hasPush {
		t.Errorf("expected add, commit, and push in executed commands: %v", executedCommands)
	}
}

func TestProcessSingleRepoSend_CleanRepoTargetAll_Mock(t *testing.T) {
	origExecutor := currentSendsGitExecutor
	defer func() { currentSendsGitExecutor = origExecutor }()

	currentSendsGitExecutor = func(dir string, args ...string) (string, error) {
		cmdStr := strings.Join(args, " ")
		if cmdStr == "status --porcelain" {
			return "", nil // clean
		}
		if cmdStr == "rev-parse --abbrev-ref HEAD" {
			return "main\n", nil
		}
		return "", nil
	}

	repo := model.ScanRecord{
		RepoName:     "clean-project",
		AbsolutePath: "/mock/clean-project",
		RelativePath: "clean-project",
	}

	opts := SendsOptions{
		Verb:       "cpf",
		Target:     "all",
		RawMessage: "test msg",
		IsDryRun:   false,
		IsPushed:   true,
	}

	rec := processSingleRepoSend(repo, opts, "Feature: test msg")

	if rec.Status != "clean-skipped" {
		t.Errorf("expected clean repository in 'all' mode to be 'clean-skipped', got %s", rec.Status)
	}
	if rec.IsSuccess {
	} else {
		t.Errorf("expected clean repository skip to be IsSuccess true")
	}
	if rec.FilesStaged != 0 {
		t.Errorf("expected FilesStaged 0 for skipped clean repo, got %d", rec.FilesStaged)
	}
}

func TestProcessSingleRepoSend_DryRun_Mock(t *testing.T) {
	origExecutor := currentSendsGitExecutor
	defer func() { currentSendsGitExecutor = origExecutor }()

	currentSendsGitExecutor = func(dir string, args ...string) (string, error) {
		cmdStr := strings.Join(args, " ")
		if cmdStr == "status --porcelain" {
			return " M modified.go\n", nil
		}
		if cmdStr == "rev-parse --abbrev-ref HEAD" {
			return "main\n", nil
		}
		return "", nil
	}

	repo := model.ScanRecord{
		RepoName:     "dry-run-repo",
		AbsolutePath: "/mock/dry-run-repo",
		RelativePath: "dry-run-repo",
	}

	opts := SendsOptions{
		Verb:       "cpb",
		Target:     "dry-run-repo",
		RawMessage: "bugfix",
		IsDryRun:   true,
		IsPushed:   true,
	}

	rec := processSingleRepoSend(repo, opts, "Bug: bugfix")

	if rec.Status != "dry-run-simulated" {
		t.Errorf("expected Status 'dry-run-simulated', got %s", rec.Status)
	}
	if rec.IsSuccess {
	} else {
		t.Errorf("expected IsSuccess true for dry run simulation")
	}
	if rec.FilesStaged != 1 {
		t.Errorf("expected FilesStaged 1, got %d", rec.FilesStaged)
	}
}

func TestSendsExecutionPayload_JSONSchema(t *testing.T) {
	payload := SendsExecutionPayload{
		Timestamp:          time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC),
		Verb:               "cpf",
		PrefixApplied:      "Feature: ",
		RawMessage:         "add tests",
		FinalCommitMessage: "Feature: add tests",
		TargetScope:        "all",
		IsDryRun:           false,
		IsPushed:           true,
		TotalProcessed:     2,
		TotalCommitted:     1,
		TotalSkipped:       1,
		Results: []RepoSendResultRecord{
			{
				RepoName:     "gitmap",
				RelativePath: ".",
				Status:       "pushed",
				Branch:       "main",
				HeadSHA:      "54c19736",
				FilesStaged:  2,
				IsSuccess:    true,
			},
			{
				RepoName:     "scripts-fixer",
				RelativePath: "scripts-fixer",
				Status:       "clean-skipped",
				Branch:       "main",
				FilesStaged:  0,
				IsSuccess:    true,
			},
		},
	}

	b, errMarshal := json.Marshal(payload)
	if errMarshal != nil {
		t.Fatalf("failed to marshal SendsExecutionPayload JSON: %v", errMarshal)
	}

	var parsed map[string]any
	if errUnmarshal := json.Unmarshal(b, &parsed); errUnmarshal != nil {
		t.Fatalf("failed to unmarshal JSON: %v", errUnmarshal)
	}

	expectedKeys := []string{
		"timestamp", "verb", "prefixApplied", "rawMessage",
		"finalCommitMessage", "targetScope", "isDryRun", "isPushed",
		"totalProcessed", "totalCommitted", "totalSkipped", "results",
	}

	for _, k := range expectedKeys {
		if _, exists := parsed[k]; !exists {
			t.Errorf("missing expected top-level JSON key: %s", k)
		}
	}

	results, ok := parsed["results"].([]any)
	if !ok || len(results) != 2 {
		t.Fatalf("expected 2 result records in JSON array")
	}

	firstRes := results[0].(map[string]any)
	resKeys := []string{"repoName", "relativePath", "status", "branch", "filesStaged", "isSuccess"}
	for _, rk := range resKeys {
		if _, exists := firstRes[rk]; !exists {
			t.Errorf("missing expected result JSON key: %s", rk)
		}
	}
}

func TestTruncateText(t *testing.T) {
	short := "hello"
	if truncateText(short, 10) != "hello" {
		t.Errorf("expected 'hello', got '%s'", truncateText(short, 10))
	}

	long := "this is a very long string that should be truncated"
	truncated := truncateText(long, 15)
	if len(truncated) != 15 || !strings.HasSuffix(truncated, "...") {
		t.Errorf("expected truncated string with '...', got '%s'", truncated)
	}
}
