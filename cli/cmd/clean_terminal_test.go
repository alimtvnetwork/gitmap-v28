package cmd

import (
	"testing"
)

func TestTerminalCleanRouting(t *testing.T) {
	commands := []string{
		"terminal",
		"clear-terminal",
		"clean-terminal",
		"terminal-clear",
		"terminal-clean",
	}

	entries := toolingDispatchEntries()

	for _, c := range commands {
		matched := false
		for _, e := range entries {
			if matchAny(c, e.names) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("expected command %q to be registered in toolingDispatchEntries", c)
		}
	}
}

func TestIsTerminalCleanSubToken(t *testing.T) {
	validTokens := []string{
		"terminal",
		"term",
		"console",
		"history",
		"shell-history",
		"clear-terminal",
		"clean-terminal",
		"terminals",
	}

	for _, tok := range validTokens {
		if !isTerminalCleanSubToken(tok) {
			t.Errorf("expected %q to be recognized as terminal clean sub token", tok)
		}
	}

	if !isTerminalCleanSubToken("clear") {
		t.Errorf("expected 'clear' to be recognized as terminal clean sub token")
	}
	if !isTerminalCleanSubToken("terminal-clear") {
		t.Errorf("expected 'terminal-clear' to be recognized as terminal clean sub token")
	}

	if isTerminalCleanSubToken("other") {
		t.Errorf("expected 'other' not to be recognized as terminal clean sub token")
	}
}

func TestRunClearAndTerminalTopLevel_DryRun(t *testing.T) {
	if err := RunClearTopLevel([]string{"-n"}); err != nil {
		t.Errorf("expected RunClearTopLevel -n to succeed, got: %v", err)
	}
	if err := RunClearTopLevel([]string{"terminal", "-n"}); err != nil {
		t.Errorf("expected RunClearTopLevel terminal -n to succeed, got: %v", err)
	}
	if err := RunTerminalTopLevel([]string{"clear", "-n"}); err != nil {
		t.Errorf("expected RunTerminalTopLevel clear -n to succeed, got: %v", err)
	}
}
