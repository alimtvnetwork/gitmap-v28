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

	// Verify the script content before base64 does not contain any python
	cleanRoot := "D:\\work\\my-repo"
	fallback := buildWindowsFallbackScript("my-repo", cleanRoot, true)
	if strings.Contains(strings.ToLower(fallback), "python") {
		t.Errorf("expected fallback script to contain ZERO python, but found python in: %s", fallback)
	}
	if strings.Contains(strings.ToLower(fallback), "state.vscdb") {
		t.Errorf("expected fallback script not to touch state.vscdb via shell script: %s", fallback)
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
	output := "Some output...\nREGISTERED:vscode,cursor,antigravity,github-desktop\nDone."
	ides := parseRegisteredIDEs(output)
	if len(ides) != 4 {
		t.Fatalf("expected 4 registered IDEs, got %d", len(ides))
	}
	if ides[0] != "vscode" || ides[1] != "cursor" || ides[2] != "antigravity" || ides[3] != "github-desktop" {
		t.Errorf("unexpected registered IDEs: %v", ides)
	}
}
