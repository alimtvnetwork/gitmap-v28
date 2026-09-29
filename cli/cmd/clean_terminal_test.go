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

	if isTerminalCleanSubToken("other") {
		t.Errorf("expected 'other' not to be recognized as terminal clean sub token")
	}
}
