package cmdssh

import (
	"os"
	"strings"
	"testing"
)

func TestCreateAskPassScript(t *testing.T) {
	scriptPath, cleanup, err := CreateAskPassScript()
	if err != nil {
		t.Fatalf("CreateAskPassScript failed: %v", err)
	}
	defer cleanup()

	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("askpass script not found on disk: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("askpass script is empty")
	}

	content, _ := os.ReadFile(scriptPath)
	if !strings.Contains(string(content), "GITMAP_SSH_PASS") {
		t.Fatalf("expected script to reference GITMAP_SSH_PASS, got: %s", string(content))
	}
}

func TestBuildAskPassEnv(t *testing.T) {
	baseEnv := []string{"PATH=/usr/bin"}
	env := BuildAskPassEnv(baseEnv, "/tmp/askpass.sh", "Secret123")

	hasAskPass := false
	hasRequire := false
	hasPass := false

	for _, e := range env {
		if strings.HasPrefix(e, "SSH_ASKPASS=") {
			hasAskPass = true
		}
		if e == "SSH_ASKPASS_REQUIRE=force" {
			hasRequire = true
		}
		if e == "GITMAP_SSH_PASS=Secret123" {
			hasPass = true
		}
	}

	if !hasAskPass || !hasRequire || !hasPass {
		t.Fatalf("missing required askpass env vars in: %v", env)
	}
}
