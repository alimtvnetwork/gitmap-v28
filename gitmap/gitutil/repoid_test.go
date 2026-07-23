package gitutil

import "testing"

// TestCanonicalRepoID_Equivalence pins the contract that the three
// transport shapes for the same repo collapse to the same identifier.
// Adding a new transport handler must not break this equivalence.
func TestCanonicalRepoID_Equivalence(t *testing.T) {
	want := "github.com/acme/widget"
	cases := []string{
		"https://github.com/acme/widget.git",
		"https://github.com/acme/widget",
		"http://github.com/acme/widget.git",
		"git@github.com:acme/widget.git",
		"git@github.com:acme/widget",
		"ssh://git@github.com/acme/widget.git",
		"  https://github.com/acme/widget.git/  ",
		"https://GitHub.com/Acme/Widget.git",
	}
	for _, in := range cases {
		got := CanonicalRepoID(in)
		if got != want {
			t.Errorf("CanonicalRepoID(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestCanonicalRepoID_Edges covers empty input and unfamiliar shapes
// to confirm graceful degradation rather than a panic or empty silent
// match between unrelated inputs.
func TestCanonicalRepoID_Edges(t *testing.T) {
	if got := CanonicalRepoID(""); got != "" {
		t.Errorf("empty: got %q, want \"\"", got)
	}
	if got := CanonicalRepoID("not-a-url"); got != "not-a-url" {
		t.Errorf("plain: got %q, want %q", got, "not-a-url")
	}
	a := CanonicalRepoID("https://gitlab.com/x/a")
	b := CanonicalRepoID("https://gitlab.com/x/b")
	if a == b {
		t.Errorf("distinct repos collapsed: %q == %q", a, b)
	}
}
