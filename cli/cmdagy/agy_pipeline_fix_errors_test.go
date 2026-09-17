package cmdagy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAgyFixArgs_FlagsAndSubcommands(t *testing.T) {
	args := []string{"fix", "errors", "agy", "--force", "-v", "--no-release"}
	opts := parseAgyFixArgs(args)

	if !opts.IsForce {
		t.Fatalf("expected IsForce = true")
	}
	if !opts.IsDetailed {
		t.Fatalf("expected IsDetailed = true")
	}
	if !opts.IsNoRelease {
		t.Fatalf("expected IsNoRelease = true")
	}
	if opts.Repo != "" {
		t.Fatalf("expected empty Repo, got %s", opts.Repo)
	}
}

func TestParseAgyFixArgs_ShortForceFlagAndTargetRepo(t *testing.T) {
	args := []string{"aef", "alimtvnetwork/gitmap-v28", "-f"}
	opts := parseAgyFixArgs(args)

	if !opts.IsForce {
		t.Fatalf("expected IsForce = true for -f")
	}
	if opts.Repo != "alimtvnetwork/gitmap-v28" {
		t.Fatalf("expected Repo = alimtvnetwork/gitmap-v28, got %s", opts.Repo)
	}
}

func TestComputeErrorSignature(t *testing.T) {
	sig1, hash1 := ComputeErrorSignature("owner/repo", 998877, "abc1234", "error output line 1")
	if sig1 != "owner/repo:998877" {
		t.Fatalf("expected sig1 owner/repo:998877, got %s", sig1)
	}
	if len(hash1) != 64 {
		t.Fatalf("expected 64-char sha256 hash, got %d chars", len(hash1))
	}

	sig2, _ := ComputeErrorSignature("owner/repo", 0, "def5678", "error output line 2")
	if sig2 != "owner/repo:def5678" {
		t.Fatalf("expected sig2 owner/repo:def5678, got %s", sig2)
	}
}

func TestDeduplicationStore_Hermetic(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "sent_errors.json")

	store := LoadSentAgyErrorsStore(storePath)
	if len(store.Records) != 0 {
		t.Fatalf("expected empty store on initial read")
	}

	sig := "test/repo:12345"
	isDup, rec := CheckSentErrorDuplicate(sig, store, false)
	if isDup || rec != nil {
		t.Fatalf("expected not duplicate on first check")
	}

	err := RecordSentErrorSignature(storePath, sig, "test/repo", 12345, "sha111", "hash222")
	if err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	reloaded := LoadSentAgyErrorsStore(storePath)
	isDupAfter, recAfter := CheckSentErrorDuplicate(sig, reloaded, false)
	if !isDupAfter || recAfter == nil {
		t.Fatalf("expected duplicate after recording without force")
	}
	if recAfter.SentCount != 1 {
		t.Fatalf("expected SentCount = 1, got %d", recAfter.SentCount)
	}

	isDupForced, recForced := CheckSentErrorDuplicate(sig, reloaded, true)
	if isDupForced || recForced == nil {
		t.Fatalf("expected isDup=false when isForce=true")
	}

	_ = RecordSentErrorSignature(storePath, sig, "test/repo", 12345, "sha111", "hash222")
	reloaded2 := LoadSentAgyErrorsStore(storePath)
	if reloaded2.Records[sig].SentCount != 2 {
		t.Fatalf("expected SentCount = 2 after resend, got %d", reloaded2.Records[sig].SentCount)
	}
}

func TestPromptAssembly_RcaHeaderAndSections(t *testing.T) {
	hdr := FormatRCAHeader("alimtvnetwork/gitmap-v28", 1234, "deadbeef")
	if !strings.Contains(hdr, "# Fix CI/CD Pipeline Errors with 4-Part Root Cause Analysis (RCA)") {
		t.Fatalf("missing RCA title in header: %s", hdr)
	}
	if !strings.Contains(hdr, "alimtvnetwork/gitmap-v28") || !strings.Contains(hdr, "#1234") {
		t.Fatalf("missing repo or runId in header: %s", hdr)
	}

	gitLogSec := FormatGitLogSection("commit 123: fix something")
	if !strings.Contains(gitLogSec, "## Recent Git Commit History") {
		t.Fatalf("missing git log heading: %s", gitLogSec)
	}

	errSec := FormatPipelineErrorSection("FAIL: test_something failed")
	if !strings.Contains(errSec, "## Failing CI/CD Pipeline Error Logs") {
		t.Fatalf("missing error logs heading: %s", errSec)
	}

	payload := AssembleRcaFixPayload("repo", 123, "sha", "log lines", "error lines", "# Prompt template")
	if !strings.Contains(payload, "log lines") || !strings.Contains(payload, "error lines") || !strings.Contains(payload, "# Prompt template") {
		t.Fatalf("payload missing required components")
	}
}

func TestVerificationFollowupPrompt_ContainsIsItFixed(t *testing.T) {
	prompt := BuildVerificationFollowupPrompt("alimtvnetwork/gitmap-v28", 9999, "feedface")
	if !strings.Contains(prompt, "# Verification Check: Is It Fixed?") {
		t.Fatalf("expected verification title in prompt: %s", prompt)
	}
	if !strings.Contains(prompt, "Is it fixed?") && !strings.Contains(prompt, "Is It Fixed?") {
		t.Fatalf("missing 'Is it fixed?' text in prompt: %s", prompt)
	}
	if !strings.Contains(prompt, "#9999") || !strings.Contains(prompt, "feedface") {
		t.Fatalf("missing runId or sha in prompt")
	}
}

func TestUpdatePromptQueueFile_Hermetic(t *testing.T) {
	tmpDir := t.TempDir()
	queuePath := filepath.Join(tmpDir, "agy-prompt-queue.json")

	primary := "# Primary RCA Fix Prompt"
	followup := "# Is It Fixed Verification"

	err := updatePromptQueueFile(queuePath, primary, followup)
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}

	data, err := os.ReadFile(queuePath)
	if err != nil {
		t.Fatalf("failed to read back queue file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "primary_fix_rca") || !strings.Contains(content, "followup_verification") {
		t.Fatalf("queue json missing expected entry types: %s", content)
	}
	if !strings.Contains(content, "dispatched") || !strings.Contains(content, "queued") {
		t.Fatalf("queue json missing status values: %s", content)
	}
}

func TestNormalizeAgyArgs_CompoundFixPhrases(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "errors_fix",
			input:    []string{"errors", "fix", "--force"},
			expected: []string{"fix-pipeline", "--force"},
		},
		{
			name:     "fix_errors",
			input:    []string{"fix", "errors", "-v"},
			expected: []string{"fix-pipeline", "-v"},
		},
		{
			name:     "fix_pipeline",
			input:    []string{"fix", "pipeline"},
			expected: []string{"fix-pipeline"},
		},
		{
			name:     "aef_alias",
			input:    []string{"aef"},
			expected: []string{"fix-pipeline"},
		},
		{
			name:     "clean_cache",
			input:    []string{"clean-cache"},
			expected: []string{"clean-cache"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := normalizeAgyArgs(tc.input)
			if len(actual) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d: %v", len(tc.expected), len(actual), actual)
			}
			for i := range actual {
				if actual[i] != tc.expected[i] {
					t.Fatalf("at index %d: expected %s, got %s", i, tc.expected[i], actual[i])
				}
			}
		})
	}
}

