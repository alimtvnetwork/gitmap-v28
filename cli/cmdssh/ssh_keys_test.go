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

func TestFormatDisplayPublicKey_Unmasked(t *testing.T) {
	sampleKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExampleKeyBlobForTesting123456789 user@host"
	displayedNonRaw := formatDisplayPublicKey(sampleKey, false)
	if displayedNonRaw != sampleKey {
		t.Errorf("expected full unmasked key, got: %s", displayedNonRaw)
	}

	displayedRaw := formatDisplayPublicKey(sampleKey, true)
	if displayedRaw != sampleKey {
		t.Errorf("expected full unmasked key with raw flag, got: %s", displayedRaw)
	}

	if strings.Contains(displayedNonRaw, "redacted") {
		t.Errorf("expected key never to be redacted")
	}
}

func TestDispatchNodeSSH_Add(t *testing.T) {
	res := dispatchNodeSSH(nil, "add", []string{"invalid-target"}, nil)
	if !res.IsMatched() {
		t.Errorf("expected 'add' subcommand to match in dispatchNodeSSH")
	}
}

func TestDispatchSSH_AuthKeyAdd(t *testing.T) {
	res1 := dispatchPackageSSH("auth-key-add", []string{"--help"})
	if !res1.IsMatched() {
		t.Errorf("expected 'auth-key-add' to match in dispatchPackageSSH")
	}

	res2 := dispatchPackageSSH("auth-key", []string{"add", "--help"})
	if !res2.IsMatched() {
		t.Errorf("expected 'auth-key add' to match in dispatchPackageSSH")
	}

	res3 := dispatchKeyOpsSSH("ssh-key", []string{"add", "--help"})
	if !res3.IsMatched() {
		t.Errorf("expected 'ssh-key add' to match in dispatchKeyOpsSSH")
	}

	res4 := dispatchKeyOpsSSH("key", []string{"add", "--help"})
	if !res4.IsMatched() {
		t.Errorf("expected 'key add' to match in dispatchKeyOpsSSH")
	}
}
