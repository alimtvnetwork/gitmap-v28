package render

import (
	"os"
	"testing"
)

// TestDecidePrettyOffAlwaysFalse locks the strongest gate: an explicit
// --pretty=false from the user must trump TTY auto-detect, env state,
// and even content type. This is the user's break-glass option.
func TestDecidePrettyOffAlwaysFalse(t *testing.T) {
	t.Setenv(EnvNoPretty, "")
	for _, isTTY := range []bool{true, false} {
		for _, isMD := range []bool{true, false} {
			if Decide(PrettyOff, isTTY, isMD) {
				t.Errorf("PrettyOff returned true for isTTY=%v isMarkdown=%v", isTTY, isMD)
			}
		}
	}
}

// TestDecideNonMarkdownAlwaysFalse guards the content gate: even with
// PrettyOn + TTY, plain text (e.g. .gitignore body) must never get
// routed through the markdown renderer.
func TestDecideNonMarkdownAlwaysFalse(t *testing.T) {
	t.Setenv(EnvNoPretty, "")
	if Decide(PrettyOn, true, false) {
		t.Fatal("non-markdown content must never be pretty-rendered")
	}
}

// TestDecidePrettyOnOverridesEnvAndTTY verifies the explicit opt-in
// beats the shared GITMAP_NO_PRETTY env var and a non-TTY stdout.
// Use case: `gitmap help foo --pretty | less -R`.
func TestDecidePrettyOnOverridesEnvAndTTY(t *testing.T) {
	t.Setenv(EnvNoPretty, "1")
	if !Decide(PrettyOn, false, true) {
		t.Fatal("PrettyOn must override GITMAP_NO_PRETTY and non-TTY stdout")
	}
}

// TestDecideAutoFollowsEnvAndTTY locks the default ladder for the
// PrettyAuto case: env opt-out wins, otherwise TTY decides.
func TestDecideAutoFollowsEnvAndTTY(t *testing.T) {
	t.Setenv(EnvNoPretty, "1")
	if Decide(PrettyAuto, true, true) {
		t.Fatal("PrettyAuto + GITMAP_NO_PRETTY must yield false even on a TTY")
	}
	_ = os.Unsetenv(EnvNoPretty)
	if !Decide(PrettyAuto, true, true) {
		t.Fatal("PrettyAuto + TTY + no env opt-out must yield true")
	}
	if Decide(PrettyAuto, false, true) {
		t.Fatal("PrettyAuto + non-TTY must yield false")
	}
}
