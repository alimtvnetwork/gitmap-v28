package cmdnodes

import (
	"strings"
	"testing"
)

func TestBuildWindowsRemoteIDECommand(t *testing.T) {
	cmd := buildWindowsRemoteIDECommand("my-repo", "D:/work/my-repo", true)
	if !strings.Contains(cmd, "powershell.exe -NoProfile -EncodedCommand") {
		t.Fatalf("expected EncodedCommand powershell invocation, got: %s", cmd)
	}
}

func TestBuildWindowsFallbackScript(t *testing.T) {
	fallback := buildWindowsFallbackScript("my-repo", "D:\\work\\my-repo", true)
	if strings.Contains(strings.ToLower(fallback), "python") {
		t.Errorf("expected fallback script to contain ZERO python, got: %s", fallback)
	}

	if strings.Contains(strings.ToLower(fallback), "state.vscdb") {
		t.Errorf("expected fallback script not to touch state.vscdb: %s", fallback)
	}
}

func TestBuildLinuxRemoteIDECommand(t *testing.T) {
	cmd := buildLinuxRemoteIDECommand("my-repo", "/home/user/work/my-repo", true)
	if strings.Contains(strings.ToLower(cmd), "python") {
		t.Fatalf("expected linux remote command to contain ZERO python, but found python in: %s", cmd)
	}

	if !strings.Contains(cmd, "gitmap ide register") {
		t.Errorf("expected command to invoke 'gitmap ide register', got: %s", cmd)
	}

	if !strings.Contains(cmd, "--pinned") {
		t.Errorf("expected command to include '--pinned', got: %s", cmd)
	}
}

func TestParseRegisteredIDEs(t *testing.T) {
	output := "Some output...\n" + RegisteredHeaderAll + "\nDone."
	ides := parseRegisteredIDEs(output)
	if len(ides) != 4 {
		t.Fatalf("expected 4 registered IDEs, got %d", len(ides))
	}

	if ides[0] != TokenVSCode || ides[1] != TokenCursor || ides[2] != TokenAntigravity || ides[3] != TokenDesktop {
		t.Errorf("unexpected registered IDEs: %v", ides)
	}
}

func TestParseRegisteredIDEsNoDesktop(t *testing.T) {
	output := "Some output...\n" + RegisteredHeaderNoDesktop + "\nDone."
	ides := parseRegisteredIDEs(output)
	if len(ides) != 3 {
		t.Fatalf("expected 3 registered IDEs, got %d", len(ides))
	}

	if ides[0] != TokenVSCode || ides[1] != TokenCursor || ides[2] != TokenAntigravity {
		t.Errorf("unexpected registered IDEs: %v", ides)
	}
}

func TestRegisteredTokensCollision(t *testing.T) {
	tokens := []string{TokenVSCode, TokenCursor, TokenAntigravity, TokenDesktop}
	hasSeen := make(map[string]bool)
	for _, tok := range tokens {
		if hasSeen[tok] {
			t.Errorf("detected duplicate collision for IDE token: %s", tok)
		}

		hasSeen[tok] = true
	}
}
