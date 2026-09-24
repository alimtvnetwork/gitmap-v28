package cmdagy

import (
	"strings"
	"testing"
)

func TestBuildNewConvArgsWithOptions(t *testing.T) {
	args := buildNewConvArgsWithOptions("My Title", "pro", "dev-profile", "Hello world")
	joined := strings.Join(args, " ")

	if !strings.Contains(joined, "--model=pro") {
		t.Errorf("expected --model=pro in args, got %v", args)
	}
	if !strings.Contains(joined, "--title=My Title") {
		t.Errorf("expected --title='My Title' in args, got %v", args)
	}
	if !strings.Contains(joined, "--profile=dev-profile") {
		t.Errorf("expected --profile=dev-profile in args, got %v", args)
	}
	if args[len(args)-1] != "Hello world" {
		t.Errorf("expected prompt as last argument, got %q", args[len(args)-1])
	}
}

func TestIsAntigravityProcessName(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{"Antigravity-default-copy-4978.exe", true},
		{"antigravity.exe", true},
		{"Antigravity.EXE", true},
		{"language_server.exe", true},
		{"notepad.exe", false},
		{"code.exe", false},
	}

	for _, tc := range testCases {
		res := isAntigravityProcessName(tc.name)
		if res != tc.expected {
			t.Errorf("for process name %q: expected %v, got %v", tc.name, tc.expected, res)
		}
	}
}

func TestResolveRecipientID_Tokens(t *testing.T) {
	id, err := resolveRecipientID("custom-id-999")
	if err != nil || id != "custom-id-999" {
		t.Errorf("expected 'custom-id-999', got %q (err: %v)", id, err)
	}

	// 'first' and 'latest' should query latest active conversation from db
	latestID, latestErr := getLatestActiveConversationID()
	if latestErr == nil && latestID != "" {
		verifyRecipientToken(t, "first", latestID)
		verifyRecipientToken(t, "latest", latestID)
	}
}

func verifyRecipientToken(t *testing.T, token, expectedID string) {
	resolved, err := resolveRecipientID(token)
	if err != nil || resolved != expectedID {
		t.Errorf("expected %q to resolve to %q, got %q (err: %v)", token, expectedID, resolved, err)
	}
}

func TestParseSendMessageArgs_TwoArgs(t *testing.T) {
	rec, content, err := parseSendMessageArgs([]string{"my-conv-123", "hello", "there"}, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != "my-conv-123" {
		t.Errorf("expected recipient 'my-conv-123', got %q", rec)
	}
	if content != "hello there" {
		t.Errorf("expected content 'hello there', got %q", content)
	}
}
