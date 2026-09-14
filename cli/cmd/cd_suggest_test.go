package cmd

import (
	"strings"
	"testing"
)

func TestEvaluateCandidateRepo_Fuzzy(t *testing.T) {
	score, hasMatch := evaluateCandidateRepo("gitamp", "gitmap")
	if !hasMatch {
		t.Fatalf("expected match between gitamp and gitmap")
	}
	if score.dist > 3 {
		t.Errorf("expected dist <= 3, got %d", score.dist)
	}
}

func TestEvaluateCandidateRepo_Substring(t *testing.T) {
	score, hasMatch := evaluateCandidateRepo("map", "gitmap")
	if !hasMatch {
		t.Fatalf("expected substring match")
	}
	if score.dist != 1 {
		t.Errorf("expected dist 1 for substring, got %d", score.dist)
	}
}

func TestSelectBestRepoSuggestions_Cap(t *testing.T) {
	scores := []repoScore{
		{name: "repo1", dist: 1},
		{name: "repo2", dist: 2},
		{name: "repo3", dist: 2},
		{name: "repo4", dist: 3},
	}
	res := selectBestRepoSuggestions(scores)
	if len(res) != 3 {
		t.Errorf("expected at most 3 suggestions, got %d", len(res))
	}
}

func TestFormatCDNotFoundMessage_WithSuggestions(t *testing.T) {
	msg := formatCDNotFoundMessage("gitamp", []string{"gitmap"})
	if !strings.Contains(msg, "no repo found matching 'gitamp'") {
		t.Errorf("missing base message: %s", msg)
	}
	if !strings.Contains(msg, "Did you mean: gitmap?") {
		t.Errorf("missing suggestion: %s", msg)
	}
}

func TestFormatCDNotFoundMessage_NoSuggestions(t *testing.T) {
	msg := formatCDNotFoundMessage("unknown", nil)
	if msg != "no repo found matching 'unknown'" {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestAppendUniqueCandidate(t *testing.T) {
	list := []string{"repo1"}
	list = appendUniqueCandidate(list, "repo1")
	if len(list) != 1 {
		t.Errorf("expected deduplication, got %v", list)
	}
	list = appendUniqueCandidate(list, "repo2")
	if len(list) != 2 {
		t.Errorf("expected new item appended, got %v", list)
	}
}
