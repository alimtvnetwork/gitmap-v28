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
	cases := []struct {
		input string
		want  string
	}{
		{"releas", "release"},
		{"scna", "scan"},
		{"fxi", "fix"},
		{"stauts", "status"},
	}

	for _, tc := range cases {
		suggestions := suggestTopLevelCommands(tc.input)
		hasMatch := false
		for _, s := range suggestions {
			if s == tc.want {
				hasMatch = true
				break
			}
		}

		if hasMatch == false {
			t.Errorf("suggestTopLevelCommands(%q) = %v, expected to contain %q", tc.input, suggestions, tc.want)
		}
	}
}

func TestFindRemediationItem_PathVariations(t *testing.T) {
	items := []RemediationItem{
		{RepoName: "scripts-fixer", RepoPath: "D:\\work\\scripts-fixer"},
		{RepoName: "gitmap", RepoPath: "D:\\work\\gitmap"},
	}

	assertFoundRepo(t, items, ".\\scripts-fixer\\", "scripts-fixer")
	assertFoundRepo(t, items, "./scripts-fixer/", "scripts-fixer")
	assertFoundRepo(t, items, "scripts-fixer", "scripts-fixer")
	assertFoundRepo(t, items, "scripts-fixer\\", "scripts-fixer")
}

func TestFindRemediationSuggestions(t *testing.T) {
	items := []RemediationItem{
		{RepoName: "scripts-fixer", RepoPath: "/path/to/scripts-fixer"},
		{RepoName: "gitmap-engine", RepoPath: "/path/to/gitmap-engine"},
	}

	suggs := FindRemediationSuggestions(items, "script-fix")
	if len(suggs) == 0 {
		t.Fatalf("expected suggestions for 'script-fix', got empty")
	}

	if suggs[0] != "scripts-fixer" {
		t.Errorf("top suggestion = %q, want 'scripts-fixer'", suggs[0])
	}
}

func TestBuildUnknownCommandMessage(t *testing.T) {
	msg := buildUnknownCommandMessage("pleae", []string{"pull", "release"})
	if strings.Contains(msg, "Unknown command: pleae") == false {
		t.Errorf("expected unknown command in message, got %q", msg)
	}

	if strings.Contains(msg, "Did you mean: pull, release") == false {
		t.Errorf("expected suggestion in message, got %q", msg)
	}
}
