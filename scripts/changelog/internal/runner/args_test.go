package runner

import (
	"strings"
	"testing"
)

// TestParseArgsAcceptsSinceAndReleaseTag pins the new flag surface so a
// future refactor cannot silently drop the partial-update entrypoints.
func TestParseArgsAcceptsSinceAndReleaseTag(t *testing.T) {
	got, err := ParseArgs([]string{
		"-mode=check",
		"-version=v3.92.0",
		"-repo=/tmp/x",
		"-since=v3.90.0",
		"-release-tag=v3.91.0",
	})
	if err != nil {
		t.Fatalf("ParseArgs: %v", err)
	}

	want := Args{
		Mode: ModeCheck, Version: "v3.92.0", RepoRoot: "/tmp/x",
		Since: "v3.90.0", ReleaseTag: "v3.91.0",
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

// TestParseArgsRejectsUnknownMode keeps the validator honest.
func TestParseArgsRejectsUnknownMode(t *testing.T) {
	_, err := ParseArgs([]string{"-mode=publish"})
	if err == nil || !strings.Contains(err.Error(), "invalid -mode") {
		t.Fatalf("want invalid-mode error, got %v", err)
	}
}

// TestResolveVersionPriority documents the four-step fallback chain so
// it never silently changes: --version → --release-tag → <lower>+next
// → vNEXT.
func TestResolveVersionPriority(t *testing.T) {
	cases := []struct {
		name  string
		args  Args
		lower string
		want  string
	}{
		{"explicit version wins", Args{Version: "v9.9.9", ReleaseTag: "vRT"}, "vLT", "v9.9.9"},
		{"release-tag fallback", Args{ReleaseTag: "vRT"}, "vLT", "vRT"},
		{"lower-bound +next", Args{}, "vLT", "vLT+next"},
		{"vNEXT default", Args{}, "", "vNEXT"},
	}
	for _, tc := range cases {
		got := resolveVersion(tc.args, tc.lower)
		if got != tc.want {
			t.Fatalf("%s: got %q want %q", tc.name, got, tc.want)
		}
	}
}
