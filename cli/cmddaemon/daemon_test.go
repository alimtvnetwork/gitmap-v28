package cmddaemon

import (
	"testing"
)

func TestResolveTimeoutMs(t *testing.T) {
	if resolveTimeoutMs(0) != 30000 {
		t.Errorf("expected 30000 for 0, got %d", resolveTimeoutMs(0))
	}
	if resolveTimeoutMs(-50) != 30000 {
		t.Errorf("expected 30000 for negative, got %d", resolveTimeoutMs(-50))
	}
	if resolveTimeoutMs(15000) != 15000 {
		t.Errorf("expected 15000, got %d", resolveTimeoutMs(15000))
	}
}

func TestBuildDaemonCmd(t *testing.T) {
	cmd := buildDaemonCmd("echo", []string{"hello", "world"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}
}
