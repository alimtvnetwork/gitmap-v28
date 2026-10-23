package cmdcommit

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestIsFullSha(t *testing.T) {
	valid := "7a3b4c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b"
	if !isFullSha(valid) {
		t.Fatalf("expected valid 40-char sha to pass")
	}

	short := "7a3b4c8"
	if isFullSha(short) {
		t.Fatalf("expected short sha to fail isFullSha")
	}

	nonHex := "7a3b4c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4z"
	if isFullSha(nonHex) {
		t.Fatalf("expected non-hex string to fail isFullSha")
	}
}

func TestIsHexRune(t *testing.T) {
	if !isHexRune('0') || !isHexRune('9') || !isHexRune('a') || !isHexRune('F') {
		t.Fatalf("expected valid hex runes to pass")
	}

	if isHexRune('g') || isHexRune('Z') || isHexRune('-') {
		t.Fatalf("expected invalid hex runes to fail")
	}
}

func TestDetectCommitOriginLocalAI(t *testing.T) {
	host, _ := os.Hostname()
	rep := &CommitInfoReport{
		Provenance: CommitProvenance{
			IsAiCommitted:   true,
			MachineHostname: host,
			AiModel:         "Claude 3.5 Sonnet",
		},
	}
	prov := DetectCommitOrigin(rep)
	if prov.OriginType != OriginLocalAI || !prov.IsLocalMachine {
		t.Fatalf("expected OriginLocalAI, got %s", prov.OriginType)
	}

	if !strings.Contains(prov.AttributionText, "Committed by local AI agent on this machine") {
		t.Fatalf("unexpected attribution: %s", prov.AttributionText)
	}
}

func TestDetectCommitOriginForeignAI(t *testing.T) {
	rep := &CommitInfoReport{
		Provenance: CommitProvenance{
			IsAiCommitted:   true,
			MachineHostname: "REMOTE-WORKER-09",
			AiModel:         "GPT-4o",
		},
	}
	prov := DetectCommitOrigin(rep)
	if prov.OriginType != OriginForeignAI || prov.IsLocalMachine {
		t.Fatalf("expected OriginForeignAI, got %s", prov.OriginType)
	}

	if !strings.Contains(prov.AttributionText, "Committed by remote cluster worker AI") {
		t.Fatalf("unexpected attribution: %s", prov.AttributionText)
	}
}

func TestDetectCommitOriginHuman(t *testing.T) {
	rep := &CommitInfoReport{
		Provenance: CommitProvenance{
			IsAiCommitted: false,
		},
	}
	prov := DetectCommitOrigin(rep)
	if prov.OriginType != OriginHumanOrExternal {
		t.Fatalf("expected OriginHumanOrExternal, got %s", prov.OriginType)
	}

	expected := "The current machine's AI has not committed this (likely another worker/node/human)"
	if prov.AttributionText != expected {
		t.Fatalf("expected exact attribution text %q, got %q", expected, prov.AttributionText)
	}
}

func buildSampleReport() *CommitInfoReport {
	rep := &CommitInfoReport{
		FullSha:         "7a3b4c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b",
		ShortSha:        "7a3b4c8",
		AuthorName:      "Dev User",
		AuthorEmail:     "dev@example.com",
		CommittedAt:     time.Date(2026, 10, 23, 21, 0, 0, 0, time.UTC),
		FormattedDate:   "2026-10-23 21:00:00 +0000",
		RelativeDate:    "10 minutes ago",
		CommitMessage:   "feat(commit): implement universal sha resolver",
		BranchName:      "main",
		TotalFiles:      1,
		TotalInsertions: 10,
		TotalDeletions:  2,
	}
	attachSampleReportDetails(rep)

	return rep
}

func attachSampleReportDetails(rep *CommitInfoReport) {
	rep.Provenance = CommitProvenance{
		OriginType:      OriginLocalAI,
		AttributionText: "Committed by local AI agent on this machine",
		AiModel:         "Claude 3.5 Sonnet",
		IsAiCommitted:   true,
		IsLocalMachine:  true,
	}
	rep.FileDeltas = []CommitFileDelta{
		{FilePath: "cli/cmdcommit/commit_info.go", Insertions: 10, Deletions: 2},
	}
}

func TestRenderCommitInfoTree(t *testing.T) {
	rep := buildSampleReport()
	buf := new(bytes.Buffer)
	err := RenderCommitInfoReport(buf, rep, false)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	out := buf.String()
	assertTreeContainsElements(t, out)
}

func assertTreeContainsElements(t *testing.T, out string) {
	elements := []string{
		"commit 7a3b4c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b",
		constants.TreeBranch,
		constants.TreeCorner,
		"Author:",
		"Date:",
		"Message:",
		"Summary:",
	}
	for _, el := range elements {
		if !strings.Contains(out, el) {
			t.Errorf("missing element %q in tree output:\n%s", el, out)
		}
	}
}

func TestRenderCommitInfoJSON(t *testing.T) {
	rep := buildSampleReport()
	buf := new(bytes.Buffer)
	err := RenderCommitInfoReport(buf, rep, true)
	if err != nil {
		t.Fatalf("render json error: %v", err)
	}

	var parsed CommitInfoReport
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Fatalf("failed unmarshaling json output: %v", unmarshalErr)
	}

	if parsed.FullSha != rep.FullSha || parsed.ShortSha != rep.ShortSha {
		t.Fatalf("mismatched json fields: got %s, expected %s", parsed.FullSha, rep.FullSha)
	}
}

func TestParseCommitInfoFlags(t *testing.T) {
	args := []string{"7a3b4c8", "--json", "--repo=/custom/repo", "-v"}
	opts, err := parseCommitInfoFlags(args)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if opts.Query != "7a3b4c8" || !opts.IsJsonFormat || !opts.IsVerbose {
		t.Fatalf("flag parsing failed: %+v", opts)
	}

	if opts.RepoPath != "/custom/repo" {
		t.Fatalf("unexpected repo path: %s", opts.RepoPath)
	}
}

func TestStripLeadingInfoSubcommand(t *testing.T) {
	args := []string{"info", "7a3b4c8", "--json"}
	clean := stripLeadingInfoSubcommand(args)
	if len(clean) != 2 || clean[0] != "7a3b4c8" {
		t.Fatalf("stripLeadingInfoSubcommand failed: %+v", clean)
	}

	normal := []string{"7a3b4c8"}
	if len(stripLeadingInfoSubcommand(normal)) != 1 {
		t.Fatalf("unexpected strip on non-info args")
	}
}

func TestFindShaInRevListLines(t *testing.T) {
	lines := []string{
		"1111111111111111111111111111111111111111",
		"7a3b4c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b",
		"2222222222222222222222222222222222222222",
	}
	sha, hasMatch := findShaInRevListLines(lines, "9e0f1a2b")
	if !hasMatch || sha != "7a3b4c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b" {
		t.Fatalf("expected middle substring match, got %s (matched: %v)", sha, hasMatch)
	}

	_, hasMissing := findShaInRevListLines(lines, "ffffffff")
	if hasMissing {
		t.Fatalf("expected non-existent sha substring to fail")
	}
}
