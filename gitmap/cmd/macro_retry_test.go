// Package cmd — macro_retry_test.go: unit tests for the macro retry-until-success engine.
package cmd

import (
	"path/filepath"
	"testing"
	"time"
)

func TestParseMacroRetryConfig(t *testing.T) {
	args := []string{"my-macro", "--sleep=3s", "--max-retries=5", "--backoff=linear", "--ai"}
	cfg := parseMacroRetryConfig(args)
	if cfg.Target != "my-macro" {
		t.Errorf("expected target my-macro, got %s", cfg.Target)
	}
	if cfg.Delay != 3*time.Second {
		t.Errorf("expected delay 3s, got %v", cfg.Delay)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("expected max retries 5, got %d", cfg.MaxRetries)
	}
	if cfg.Backoff != "linear" {
		t.Errorf("expected backoff linear, got %s", cfg.Backoff)
	}
	if !cfg.IsAI {
		t.Errorf("expected IsAI to be true")
	}
}

func TestCalculateBackoff(t *testing.T) {
	base := 2 * time.Second
	if val := calculateBackoff(base, "fixed", 3); val != 2*time.Second {
		t.Errorf("expected fixed 2s, got %v", val)
	}
	if val := calculateBackoff(base, "linear", 3); val != 6*time.Second {
		t.Errorf("expected linear 6s, got %v", val)
	}
	if val := calculateBackoff(base, "exponential", 3); val != 8*time.Second {
		t.Errorf("expected exponential 8s, got %v", val)
	}
}

func TestMacroRetrySuccessFlow(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	aiFile := filepath.Join(tempDir, "ai-error.md")
	args := []string{"echo Success!", "--sleep=10ms", "--max-retries=2", "--ai-file=" + aiFile}
	if err := runMacroUntilSuccess(args); err != nil {
		t.Fatalf("runMacroUntilSuccess failed: %v", err)
	}
}
