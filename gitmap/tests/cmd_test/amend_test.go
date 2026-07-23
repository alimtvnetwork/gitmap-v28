package cmd_test

import (
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// --- extractAmendSHA tests ---

func extractAmendSHA(args []string, commitHash *string) []string {
	if len(args) == 0 {
		return args
	}

	if args[0] == "" || args[0][0] == '-' {
		return args
	}

	*commitHash = args[0]

	return args[1:]
}

func TestExtractAmendSHA_Empty(t *testing.T) {
	var sha string
	remaining := extractAmendSHA([]string{}, &sha)
	if sha != "" {
		t.Errorf("expected empty SHA, got %q", sha)
	}
	if len(remaining) != 0 {
		t.Errorf("expected empty remaining, got %v", remaining)
	}
}

func TestExtractAmendSHA_WithSHA(t *testing.T) {
	var sha string
	remaining := extractAmendSHA([]string{"abc123", "--name", "Test"}, &sha)
	if sha != "abc123" {
		t.Errorf("expected sha=abc123, got %q", sha)
	}
	if len(remaining) != 2 || remaining[0] != "--name" {
		t.Errorf("expected [--name Test], got %v", remaining)
	}
}

func TestExtractAmendSHA_FlagFirst(t *testing.T) {
	var sha string
	remaining := extractAmendSHA([]string{"--name", "Test"}, &sha)
	if sha != "" {
		t.Errorf("expected empty SHA, got %q", sha)
	}
	if len(remaining) != 2 {
		t.Errorf("expected 2 remaining args, got %d", len(remaining))
	}
}

func TestExtractAmendSHA_HEAD(t *testing.T) {
	var sha string
	remaining := extractAmendSHA([]string{"HEAD", "--email", "a@b.com"}, &sha)
	if sha != "HEAD" {
		t.Errorf("expected sha=HEAD, got %q", sha)
	}
	if len(remaining) != 2 {
		t.Errorf("expected 2 remaining, got %d", len(remaining))
	}
}

// --- resolveAmendMode tests ---

func resolveAmendMode(commitHash string) string {
	if commitHash == "" {
		return constants.AmendModeAll
	}

	if commitHash == "HEAD" {
		return constants.AmendModeHead
	}

	return constants.AmendModeRange
}

func TestResolveAmendMode_Empty(t *testing.T) {
	mode := resolveAmendMode("")
	if mode != "all" {
		t.Errorf("expected all, got %q", mode)
	}
}

func TestResolveAmendMode_HEAD(t *testing.T) {
	mode := resolveAmendMode("HEAD")
	if mode != "head" {
		t.Errorf("expected head, got %q", mode)
	}
}

func TestResolveAmendMode_SHA(t *testing.T) {
	mode := resolveAmendMode("abc123")
	if mode != "range" {
		t.Errorf("expected range, got %q", mode)
	}
}

// --- parseCommitLines tests ---

func parseCommitLines(output string) []model.CommitEntry {
	if output == "" {
		return nil
	}

	lines := splitLines(output)
	entries := make([]model.CommitEntry, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := splitFirst(line, " ")
		msg := ""
		if len(parts) > 1 {
			msg = parts[1]
		}

		entries = append(entries, model.CommitEntry{
			SHA:     parts[0],
			Message: msg,
		})
	}

	return entries
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}

	return lines
}

func splitFirst(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep[0] {
			return []string{s[:i], s[i+1:]}
		}
	}

	return []string{s}
}

func TestParseCommitLines_Multiple(t *testing.T) {
	input := "abc1234 Fix login page\ndef5678 Add dashboard"
	entries := parseCommitLines(input)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].SHA != "abc1234" || entries[0].Message != "Fix login page" {
		t.Errorf("entry 0 mismatch: %+v", entries[0])
	}
	if entries[1].SHA != "def5678" || entries[1].Message != "Add dashboard" {
		t.Errorf("entry 1 mismatch: %+v", entries[1])
	}
}

func TestParseCommitLines_Empty(t *testing.T) {
	entries := parseCommitLines("")
	if entries != nil {
		t.Errorf("expected nil, got %v", entries)
	}
}

func TestParseCommitLines_SingleLine(t *testing.T) {
	entries := parseCommitLines("abc1234 Solo commit")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].SHA != "abc1234" || entries[0].Message != "Solo commit" {
		t.Errorf("entry mismatch: %+v", entries[0])
	}
}

func TestParseCommitLines_NoMessage(t *testing.T) {
	entries := parseCommitLines("abc1234")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].SHA != "abc1234" || entries[0].Message != "" {
		t.Errorf("entry mismatch: %+v", entries[0])
	}
}

// --- buildAuditRecord tests ---

type testAmendFlags struct {
	name      string
	email     string
	forcePush bool
}

func buildAuditRecord(f testAmendFlags, commits []model.CommitEntry, branch, mode, prevName, prevEmail string, ts time.Time) model.AmendmentRecord {
	fromCommit := ""
	toCommit := ""

	if len(commits) > 0 {
		fromCommit = commits[0].SHA
		toCommit = commits[len(commits)-1].SHA
	}

	return model.AmendmentRecord{
		Timestamp:    ts.Format(time.RFC3339),
		Branch:       branch,
		FromCommit:   fromCommit,
		ToCommit:     toCommit,
		TotalCommits: len(commits),
		PreviousAuthor: model.AmendAuthor{
			Name:  prevName,
			Email: prevEmail,
		},
		NewAuthor: model.AmendAuthor{
			Name:  f.name,
			Email: f.email,
		},
		Mode:        mode,
		ForcePushed: f.forcePush,
		Commits:     commits,
	}
}

func TestBuildAuditRecord_Range(t *testing.T) {
	commits := []model.CommitEntry{
		{SHA: "aaa1111", Message: "First"},
		{SHA: "bbb2222", Message: "Second"},
	}
	flags := testAmendFlags{name: "New Name", email: "new@test.com", forcePush: true}
	ts := time.Date(2026, 3, 9, 14, 30, 0, 0, time.UTC)

	record := buildAuditRecord(flags, commits, "develop", "range", "Old Name", "old@test.com", ts)

	if record.Branch != "develop" {
		t.Errorf("expected branch=develop, got %q", record.Branch)
	}
	if record.FromCommit != "aaa1111" {
		t.Errorf("expected fromCommit=aaa1111, got %q", record.FromCommit)
	}
	if record.ToCommit != "bbb2222" {
		t.Errorf("expected toCommit=bbb2222, got %q", record.ToCommit)
	}
	if record.TotalCommits != 2 {
		t.Errorf("expected totalCommits=2, got %d", record.TotalCommits)
	}
	if record.PreviousAuthor.Name != "Old Name" {
		t.Errorf("expected prevName=Old Name, got %q", record.PreviousAuthor.Name)
	}
	if record.NewAuthor.Email != "new@test.com" {
		t.Errorf("expected newEmail=new@test.com, got %q", record.NewAuthor.Email)
	}
	if record.Mode != "range" {
		t.Errorf("expected mode=range, got %q", record.Mode)
	}
	if record.ForcePushed != true {
		t.Error("expected forcePushed=true")
	}
	if record.Timestamp != "2026-03-09T14:30:00Z" {
		t.Errorf("expected timestamp 2026-03-09T14:30:00Z, got %q", record.Timestamp)
	}
}

func TestBuildAuditRecord_EmptyCommits(t *testing.T) {
	flags := testAmendFlags{name: "Test", email: "t@t.com"}

	record := buildAuditRecord(flags, nil, "main", "all", "", "", time.Now())

	if record.FromCommit != "" || record.ToCommit != "" {
		t.Error("expected empty from/to commits for nil slice")
	}
	if record.TotalCommits != 0 {
		t.Errorf("expected 0 total, got %d", record.TotalCommits)
	}
}

func TestBuildAuditRecord_HeadMode(t *testing.T) {
	commits := []model.CommitEntry{{SHA: "head123", Message: "Latest"}}
	flags := testAmendFlags{name: "Bot", email: "bot@ci.com"}

	record := buildAuditRecord(flags, commits, "main", "head", "Dev", "dev@co.com", time.Now())

	if record.Mode != "head" {
		t.Errorf("expected mode=head, got %q", record.Mode)
	}
	if record.FromCommit != "head123" || record.ToCommit != "head123" {
		t.Error("expected from=to for single HEAD commit")
	}
	if record.ForcePushed != false {
		t.Error("expected forcePushed=false")
	}
}

// --- formatAuditTimestamp tests ---

func formatAuditTimestamp(ts time.Time) string {
	return ts.Format("2006-01-02T15-04-05")
}

func TestFormatAuditTimestamp(t *testing.T) {
	ts := time.Date(2026, 3, 9, 14, 30, 45, 0, time.UTC)
	result := formatAuditTimestamp(ts)
	expected := "2026-03-09T14-30-45"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestFormatAuditTimestamp_Midnight(t *testing.T) {
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	result := formatAuditTimestamp(ts)
	expected := "2026-01-01T00-00-00"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
