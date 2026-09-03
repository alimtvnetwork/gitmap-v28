// Package cmd — macro_add_test.go: unit tests for CLI macro creation.
package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func TestHandleMacroAdd(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	args := []string{"test-build", "go build -o app.exe .", "go test ./...", "--desc=Build and test pipeline", "--tag=ci"}
	if err := handleMacroAdd(args); err != nil {
		t.Fatalf("handleMacroAdd failed: %v", err)
	}

	loaded, err := macro.LoadMacro("test-build")
	if err != nil {
		t.Fatalf("load created macro failed: %v", err)
	}
	if loaded.Name != "test-build" || len(loaded.Steps) != 2 {
		t.Fatalf("unexpected macro: %+v", loaded)
	}
	if loaded.Description != "Build and test pipeline" || loaded.Tags != "ci" {
		t.Errorf("unexpected metadata: desc=%q, tag=%q", loaded.Description, loaded.Tags)
	}
}

func TestHandleMacroAddChainedSteps(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	args := []string{"chained-macro", "echo step1 && echo step2 && echo step3"}
	if err := handleMacroAdd(args); err != nil {
		t.Fatalf("handleMacroAdd chained failed: %v", err)
	}

	loaded, err := macro.LoadMacro("chained-macro")
	if err != nil {
		t.Fatalf("load chained macro failed: %v", err)
	}
	if len(loaded.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(loaded.Steps))
	}
}

func TestHandleMacroAddValidationErrors(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	// Zero args
	if err := handleMacroAdd([]string{}); err == nil {
		t.Fatal("expected error for empty args, got nil")
	}

	// Name only in non-interactive environment (EOF immediately)
	if err := handleMacroAdd([]string{"alim"}); err != nil {
		t.Fatalf("expected nil when no commands entered interactively, got: %v", err)
	}

	// Verify macro was not saved
	if _, err := macro.LoadMacro("alim"); err == nil {
		t.Fatal("expected macro 'alim' not to be saved, but found it")
	}
}
