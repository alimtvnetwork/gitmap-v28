package cmd

import (
	"strings"
	"testing"
)

func TestPromptPreflightConfirmation_Affirmative(t *testing.T) {
	pairs := []RenamePair{
		{OldBase: "README.md", NewBase: "readme.md"},
	}
	opts := LowerCaseFixOptions{Patterns: []string{"*.md"}}

	inputs := []string{"confirm\n", "CONFIRM\n", "yes\n", "YES\n", "y\n", "Y\n"}
	for _, input := range inputs {
		r := strings.NewReader(input)
		isConfirmed, err := promptPreflightConfirmationWithReader(r, pairs, opts)
		if err != nil {
			t.Fatalf("unexpected error for input %q: %v", input, err)
		}
		if !isConfirmed {
			t.Errorf("expected isConfirmed=true for input %q", input)
		}
	}
}

func TestPromptPreflightConfirmation_Negative(t *testing.T) {
	pairs := []RenamePair{
		{OldBase: "README.md", NewBase: "readme.md"},
	}
	opts := LowerCaseFixOptions{Patterns: []string{"*.md"}}

	inputs := []string{"no\n", "n\n", "cancel\n", "\n", "something else\n", ""}
	for _, input := range inputs {
		r := strings.NewReader(input)
		isConfirmed, err := promptPreflightConfirmationWithReader(r, pairs, opts)
		if err != nil {
			t.Fatalf("unexpected error for input %q: %v", input, err)
		}
		if isConfirmed {
			t.Errorf("expected isConfirmed=false for input %q", input)
		}
	}
}

func TestCheckPreflightConfirmation_BypassConditions(t *testing.T) {
	pairs := []RenamePair{
		{OldBase: "README.md", NewBase: "readme.md"},
	}

	// 1. Non-git repository bypasses confirmation
	optsNoGit := LowerCaseFixOptions{IsYes: false}
	isProceed, err := checkPreflightConfirmation(pairs, optsNoGit, false)
	if err != nil || !isProceed {
		t.Errorf("expected bypass for non-git directory, got proceed=%v, err=%v", isProceed, err)
	}

	// 2. IsYes flag bypasses confirmation even in git
	optsYes := LowerCaseFixOptions{IsYes: true}
	isProceed, err = checkPreflightConfirmation(pairs, optsYes, true)
	if err != nil || !isProceed {
		t.Errorf("expected bypass for IsYes=true, got proceed=%v, err=%v", isProceed, err)
	}
}

func TestCheckPreflightConfirmation_InteractiveInGit(t *testing.T) {
	pairs := []RenamePair{
		{OldBase: "README.md", NewBase: "readme.md"},
	}
	opts := LowerCaseFixOptions{IsYes: false}

	// Test confirmed
	oldReader := lcfStdinReader
	defer func() { lcfStdinReader = oldReader }()

	lcfStdinReader = strings.NewReader("confirm\n")
	isProceed, err := checkPreflightConfirmation(pairs, opts, true)
	if err != nil || !isProceed {
		t.Errorf("expected proceed=true on 'confirm', got %v, err=%v", isProceed, err)
	}

	// Test canceled
	lcfStdinReader = strings.NewReader("no\n")
	isProceed, err = checkPreflightConfirmation(pairs, opts, true)
	if err != nil {
		t.Fatalf("unexpected error on 'no': %v", err)
	}
	if isProceed {
		t.Errorf("expected proceed=false on 'no'")
	}
}
