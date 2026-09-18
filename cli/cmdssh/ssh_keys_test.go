package cmdssh

import (
	"os"
	"strings"
	"testing"
)

func TestFormatMissingAuthAdvice(t *testing.T) {
	advice := formatMissingAuthAdvice("worker-1", "192.168.1.8", "ubuntu")

	if !strings.Contains(advice, "worker-1") {
		t.Errorf("expected advice to contain alias, got: %s", advice)
	}

	if !strings.Contains(advice, "192.168.1.8") {
		t.Errorf("expected advice to contain IP, got: %s", advice)
	}

	if !strings.Contains(advice, "add-with-pass") {
		t.Errorf("expected advice to contain add-with-pass hint, got: %s", advice)
	}
}

func TestFindDefaultUserSSHKey(t *testing.T) {
	keyPath := findDefaultUserSSHKey()
	if keyPath == "" {
		return
	}

	fi, err := os.Stat(keyPath)
	if err != nil || fi.IsDir() {
		t.Errorf("expected found key path to be an existing file: %s", keyPath)
	}
}

func TestConnectWithDefaultKey_Unreachable(t *testing.T) {
	client, isConnected := connectWithDefaultKey("192.0.2.1", "testuser", "[test]")
	if isConnected || client != nil {
		t.Errorf("expected connection failure on TEST-NET IP")
	}
}
