package cmdcommit

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromptCommitMessageIfMissing_ExistingArgs(t *testing.T) {
	args := []string{"feat: existing message"}
	got, appErr := PromptCommitMessageIfMissing(args, "commit-push", "E9000")
	if appErr != nil {
		t.Fatalf("unexpected appErr: %v", appErr)
	}
	if len(got) != 1 || got[0] != "feat: existing message" {
		t.Fatalf("expected existing args, got %v", got)
	}
}

func TestPromptCommitMessageIfMissing_NonInteractive(t *testing.T) {
	prevCheck := isTerminalStdinFunc
	isTerminalStdinFunc = func() bool { return false }
	defer func() { isTerminalStdinFunc = prevCheck }()

	got, appErr := PromptCommitMessageIfMissing([]string{}, "commit-push", "E9000")
	if appErr == nil {
		t.Fatal("expected appErr in non-interactive environment, got nil")
	}
	if appErr.Code != "E9000" {
		t.Errorf("expected error code E9000, got %s", appErr.Code)
	}
	if got != nil {
		t.Errorf("expected nil got, got %v", got)
	}
}

func TestPromptCommitMessageIfMissing_InteractiveSuccess(t *testing.T) {
	prevCheck := isTerminalStdinFunc
	prevStdin := promptStdinReader
	prevStderr := promptStderr
	isTerminalStdinFunc = func() bool { return true }
	promptStdinReader = strings.NewReader("fix: prompt message\n")
	var errBuf bytes.Buffer
	promptStderr = &errBuf
	defer func() {
		isTerminalStdinFunc = prevCheck
		promptStdinReader = prevStdin
		promptStderr = prevStderr
	}()

	got, appErr := PromptCommitMessageIfMissing([]string{}, "commit-push", "E9000")
	if appErr != nil {
		t.Fatalf("unexpected appErr: %v", appErr)
	}
	if len(got) != 1 || got[0] != "fix: prompt message" {
		t.Fatalf("expected interactive input, got %v", got)
	}
	if !strings.Contains(errBuf.String(), "Enter commit message") {
		t.Errorf("expected prompt text in stderr buffer, got %q", errBuf.String())
	}
}

func TestPromptParameterIfMissing(t *testing.T) {
	val, err := PromptParameterIfMissing("already-present", "task", "Enter task: ")
	if err != nil || val != "already-present" {
		t.Fatalf("expected already-present, got %q, err %v", val, err)
	}

	prevCheck := isTerminalStdinFunc
	isTerminalStdinFunc = func() bool { return false }
	defer func() { isTerminalStdinFunc = prevCheck }()

	_, errNonInt := PromptParameterIfMissing("", "task", "Enter task: ")
	if errNonInt == nil {
		t.Fatal("expected error in non-interactive environment")
	}
}
