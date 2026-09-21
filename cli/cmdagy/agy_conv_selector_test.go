package cmdagy

import (
	"bytes"
	"strings"
	"testing"
)

func assertEqualID(t *testing.T, got, want string) {
	isMatch := got == want
	if isMatch {
		return
	}

	t.Errorf("ID mismatch: got %q, want %q", got, want)
}

func makeSampleConvs() []AgyConvInfo {
	return []AgyConvInfo{
		{ID: "conv-1", UserSteps: 10, StepCount: 20},
		{ID: "conv-2", UserSteps: 5, StepCount: 15},
		{ID: "conv-3", UserSteps: 1, StepCount: 5},
	}
}

func TestPromptSelect_ZeroMatches(t *testing.T) {
	var buf bytes.Buffer
	_, err := PromptSelectConversation(nil, strings.NewReader(""), &buf, false)
	isSuccess := err == nil
	if isSuccess {
		t.Fatalf("expected error for zero matches, got nil")
	}
}

func TestPromptSelect_SingleMatch(t *testing.T) {
	var buf bytes.Buffer
	convs := []AgyConvInfo{{ID: "conv-single", UserSteps: 2, StepCount: 5}}
	res, err := PromptSelectConversation(convs, strings.NewReader(""), &buf, false)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqualID(t, res.ID, "conv-single")
}

func TestPromptSelect_MultipleMatches_NonInteractive(t *testing.T) {
	var buf bytes.Buffer
	convs := makeSampleConvs()
	res, err := PromptSelectConversation(convs, strings.NewReader(""), &buf, false)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqualID(t, res.ID, "conv-1")
}

func TestPromptSelect_MultipleMatches_Interactive_Index1(t *testing.T) {
	var buf bytes.Buffer
	convs := makeSampleConvs()
	res, err := PromptSelectConversation(convs, strings.NewReader("1\n"), &buf, true)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqualID(t, res.ID, "conv-1")
}

func TestPromptSelect_MultipleMatches_Interactive_Index2(t *testing.T) {
	var buf bytes.Buffer
	convs := makeSampleConvs()
	res, err := PromptSelectConversation(convs, strings.NewReader("2\n"), &buf, true)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqualID(t, res.ID, "conv-2")
}

func TestPromptSelect_MultipleMatches_Interactive_EmptyString(t *testing.T) {
	var buf bytes.Buffer
	convs := makeSampleConvs()
	res, err := PromptSelectConversation(convs, strings.NewReader("\n"), &buf, true)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqualID(t, res.ID, "conv-1")
}

func TestPromptSelect_MultipleMatches_Interactive_InvalidInput(t *testing.T) {
	var buf bytes.Buffer
	convs := makeSampleConvs()
	res, err := PromptSelectConversation(convs, strings.NewReader("invalid-99\n"), &buf, true)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqualID(t, res.ID, "conv-1")
}

func TestSortConversations(t *testing.T) {
	convs := []AgyConvInfo{
		{ID: "conv-b", UserSteps: 5, StepCount: 10},
		{ID: "conv-a", UserSteps: 5, StepCount: 10},
		{ID: "conv-c", UserSteps: 10, StepCount: 5},
		{ID: "conv-d", UserSteps: 5, StepCount: 20},
	}
	sortConversations(convs)
	assertEqualID(t, convs[0].ID, "conv-c")
	assertEqualID(t, convs[1].ID, "conv-d")
	assertEqualID(t, convs[2].ID, "conv-a")
	assertEqualID(t, convs[3].ID, "conv-b")
}

func TestSelectMatchingConversation_ZeroMatches(t *testing.T) {
	_, err := SelectMatchingConversation("/non/existent/repo/path/xyz")
	hasErr := err != nil
	if hasErr == false {
		t.Fatalf("expected error for non-existent workspace, got nil")
	}
}

func TestSelectMatchingConversation_LocalWorkspaces(t *testing.T) {
	conv, err := SelectMatchingConversation("d:\\work\\gitmap")
	hasErr := err != nil
	if hasErr {
		t.Skip("skipping local workspace test when summaries db not present")

		return
	}

	hasEmptyID := conv.ID == ""
	if hasEmptyID {
		t.Errorf("expected non-empty conversation ID, got empty")
	}
}
