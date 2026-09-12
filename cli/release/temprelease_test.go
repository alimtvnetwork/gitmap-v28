package release

import (
	"testing"
)

func TestParseTempReleaseCommitLines_Basic(t *testing.T) {
	input := "abc123def456|Initial commit\n789012345678|Second commit"

	commits := parseTempReleaseCommitLines(input)
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}

	if commits[0].Message != "Initial commit" {
		t.Errorf("first message = %q, want %q", commits[0].Message, "Initial commit")
	}

	if commits[1].SHA != "789012345678" {
		t.Errorf("second SHA = %q, want %q", commits[1].SHA, "789012345678")
	}
}

func TestParseTempReleaseCommitLines_Empty(t *testing.T) {
	commits := parseTempReleaseCommitLines("")
	if commits != nil {
		t.Errorf("expected nil for empty input, got %v", commits)
	}
}

func TestParseOneCommitLine_Valid(t *testing.T) {
	line := "abcdef1234567890abcdef1234567890abcdef12|Fix bug in scanner"

	c := parseOneCommitLine(line)
	if c.SHA != "abcdef1234567890abcdef1234567890abcdef12" {
		t.Errorf("SHA = %q", c.SHA)
	}

	if c.Message != "Fix bug in scanner" {
		t.Errorf("Message = %q", c.Message)
	}

	if len(c.Short) > 7 {
		t.Errorf("Short should be truncated, got %q", c.Short)
	}
}

func TestParseOneCommitLine_NoPipe(t *testing.T) {
	c := parseOneCommitLine("no-pipe-here")
	if len(c.SHA) > 0 {
		t.Errorf("expected empty commit for invalid line, got SHA=%q", c.SHA)
	}
}

func TestParseBranchOutput_Basic(t *testing.T) {
	input := "  temp-release/v1.01\n* temp-release/v1.02\n  temp-release/v1.03"

	branches := parseBranchOutput(input)
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}

	if branches[1] != "temp-release/v1.02" {
		t.Errorf("branch[1] = %q, want temp-release/v1.02", branches[1])
	}
}

func TestParseBranchOutput_Empty(t *testing.T) {
	branches := parseBranchOutput("")
	if branches != nil {
		t.Errorf("expected nil for empty input, got %v", branches)
	}
}
