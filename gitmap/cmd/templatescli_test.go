package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/render"
)

// TestIsMarkdownTemplatePathRecognizesMarkdownExtensions locks the file-
// extension allow-list. Adding a new markdown extension is a deliberate
// choice — this test will scream if someone broadens it accidentally.
// The non-markdown rows are the load-bearing assertions: they protect
// every existing `templates show ignore go > .gitignore` redirect on
// disk in users' repos from suddenly gaining ANSI bytes.
func TestIsMarkdownTemplatePathRecognizesMarkdownExtensions(t *testing.T) {
	cases := map[string]bool{
		"assets/notes/intro.md":              true,
		"assets/notes/intro.MD":              true,
		"assets/notes/intro.markdown":        true,
		"assets/notes/intro.MARKDOWN":        true,
		"assets/ignore/go.gitignore":         false,
		"assets/attributes/go.gitattributes": false,
		"assets/lfs/common.gitattributes":    false,
		"plain":                              false,
		"":                                   false,
	}
	for path, want := range cases {
		if got := isMarkdownTemplatePath(path); got != want {
			t.Errorf("isMarkdownTemplatePath(%q) = %v, want %v", path, got, want)
		}
	}
}

// TestParseTemplatesShowFlagsPrefersPretty checks the new --pretty path:
// positional <kind> <lang> survive intact and the flag maps to the
// correct render.PrettyMode regardless of arg order.
func TestParseTemplatesShowFlagsPrefersPretty(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want render.PrettyMode
	}{
		{"no flag → auto", []string{"ignore", "go"}, render.PrettyAuto},
		{"--pretty → on", []string{"ignore", "go", "--pretty"}, render.PrettyOn},
		{"--no-pretty → off", []string{"--no-pretty", "ignore", "go"}, render.PrettyOff},
		{"--pretty=false → off", []string{"ignore", "--pretty=false", "go"}, render.PrettyOff},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rest, mode := parseTemplatesShowFlags(tc.args)
			if mode != tc.want {
				t.Errorf("mode = %v, want %v", mode, tc.want)
			}
			if len(rest) != 2 || rest[0] != "ignore" || rest[1] != "go" {
				t.Errorf("rest = %v, want [ignore go]", rest)
			}
		})
	}
}

// TestParseTemplatesShowFlagsRawAliases is the back-compat guard: the
// legacy --raw flag from v3.23.x must still downgrade to PrettyOff so
// scripts that already use it keep working without modification.
func TestParseTemplatesShowFlagsRawAliases(t *testing.T) {
	rest, mode := parseTemplatesShowFlags([]string{"ignore", "go", "--raw"})
	if mode != render.PrettyOff {
		t.Errorf("--raw → %v, want PrettyOff (back-compat alias)", mode)
	}
	if len(rest) != 2 || rest[0] != "ignore" || rest[1] != "go" {
		t.Errorf("rest = %v, want [ignore go]", rest)
	}
}

// TestParseTemplatesShowFlagsPrettyBeatsRaw guards the conflict-resolution
// rule: when a user passes both --pretty and the deprecated --raw, the
// preferred new flag wins. Otherwise --raw could silently override an
// explicit --pretty in shell aliases that bake one of them in.
func TestParseTemplatesShowFlagsPrettyBeatsRaw(t *testing.T) {
	_, mode := parseTemplatesShowFlags([]string{"--raw", "--pretty", "ignore", "go"})
	if mode != render.PrettyOn {
		t.Fatalf("--pretty must beat --raw: got %v, want PrettyOn", mode)
	}
}
