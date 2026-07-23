package cmd

import (
	"testing"
)

func TestParseSections_Empty(t *testing.T) {
	got := parseSections("")
	if got != nil {
		t.Errorf("expected nil for empty input, got %v", got)
	}
}

func TestParseSections_SingleValid(t *testing.T) {
	got := parseSections("commands")
	if len(got) != 1 || !got["commands"] {
		t.Errorf("expected {commands: true}, got %v", got)
	}
}

func TestParseSections_MultipleValid(t *testing.T) {
	got := parseSections("commands,architecture,flags")
	if len(got) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(got))
	}

	for _, s := range []string{"commands", "architecture", "flags"} {
		if !got[s] {
			t.Errorf("expected section %q to be present", s)
		}
	}
}

func TestParseSections_AllValid(t *testing.T) {
	got := parseSections("commands,architecture,flags,conventions,structure,database,installation,patterns")
	if len(got) != 8 {
		t.Errorf("expected 8 sections, got %d", len(got))
	}
}

func TestParseSections_TrailingComma(t *testing.T) {
	got := parseSections("commands,")
	if len(got) != 1 || !got["commands"] {
		t.Errorf("expected {commands: true}, got %v", got)
	}
}

func TestParseSections_WhitespaceHandling(t *testing.T) {
	got := parseSections(" commands , architecture ")
	if len(got) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got))
	}

	if !got["commands"] || !got["architecture"] {
		t.Errorf("expected trimmed sections, got %v", got)
	}
}

func TestWantSection_NilSetIncludesAll(t *testing.T) {
	if !wantSection(nil, "commands") {
		t.Error("nil set should include all sections")
	}

	if !wantSection(nil, "anything") {
		t.Error("nil set should include any name")
	}
}

func TestWantSection_FilteredSet(t *testing.T) {
	set := map[string]bool{"commands": true, "flags": true}

	if !wantSection(set, "commands") {
		t.Error("expected commands to be included")
	}

	if !wantSection(set, "flags") {
		t.Error("expected flags to be included")
	}

	if wantSection(set, "architecture") {
		t.Error("expected architecture to be excluded")
	}

	if wantSection(set, "database") {
		t.Error("expected database to be excluded")
	}
}

func TestWantSection_EmptySet(t *testing.T) {
	set := map[string]bool{}

	if wantSection(set, "commands") {
		t.Error("empty set should exclude all sections")
	}
}
