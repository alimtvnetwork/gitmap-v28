package cmd

import (
	"strings"
	"testing"
)

func TestSuggestTopLevelCommands_Pleae(t *testing.T) {
	suggestions := suggestTopLevelCommands("pleae")
	if len(suggestions) < 3 {
		t.Fatalf("expected at least 3 suggestions for 'pleae', got %v", suggestions)
	}

	expected := []string{"pull", "release", "pull-release"}
	for i, want := range expected {
		if suggestions[i] != want {
			t.Errorf("suggestions[%d] = %q, want %q", i, suggestions[i], want)
		}
	}
}

func TestSuggestTopLevelCommands_Please(t *testing.T) {
	suggestions := suggestTopLevelCommands("please")
	if len(suggestions) < 3 {
		t.Fatalf("expected at least 3 suggestions for 'please', got %v", suggestions)
	}

	if suggestions[0] != "pull" || suggestions[1] != "release" || suggestions[2] != "pull-release" {
		t.Errorf("unexpected suggestions: %v", suggestions)
	}
}

func TestSuggestTopLevelCommands_Typos(t *testing.T) {
	cases := []struct{ input, want string }{
		{"releas", "release"},
		{"scna", "scan"},
		{"fxi", "fix"},
		{"stauts", "status"},
	}

	for _, tc := range cases {
		verifyTypoSuggestion(t, tc.input, tc.want)
	}
}

func verifyTypoSuggestion(t *testing.T, input, want string) {
	suggestions := suggestTopLevelCommands(input)
	if !containsSuggestion(suggestions, want) {
		t.Errorf("suggestTopLevelCommands(%q) = %v, expected %q", input, suggestions, want)
	}
}

func containsSuggestion(suggestions []string, target string) bool {
	for _, s := range suggestions {
		if s == target {
			return true
		}
	}

	return false
}

func TestBuildUnknownCommandMessage(t *testing.T) {
	msg := buildUnknownCommandMessage("pleae", []string{"pull", "release"})
	if !strings.Contains(msg, "Unknown command: pleae") {
		t.Errorf("expected unknown command in message, got %q", msg)
	}

	if !strings.Contains(msg, "Did you mean: pull, release") {
		t.Errorf("expected suggestion in message, got %q", msg)
	}
}
