package cmdssh

import "testing"

const sampleTestPubKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl user@laptop"

func TestValidateSSHPublicKey_Valid(t *testing.T) {
	key, blob, err := validateSSHPublicKey(sampleTestPubKey)
	if err != nil {
		t.Fatalf("expected valid key, got: %v", err)
	}

	if key != sampleTestPubKey {
		t.Errorf("expected key to match trimmed input")
	}

	if blob == "" {
		t.Errorf("expected non-empty base64 blob")
	}
}

func TestValidateSSHPublicKey_Invalid(t *testing.T) {
	_, _, err := validateSSHPublicKey("invalid-ssh-key-data")
	if err == nil {
		t.Fatalf("expected error for invalid key")
	}
}

func TestValidateSSHPublicKey_Empty(t *testing.T) {
	_, _, err := validateSSHPublicKey("   ")
	if err == nil {
		t.Fatalf("expected error for empty key")
	}
}

func TestIsKeyInAuthorizedKeys_Match(t *testing.T) {
	_, blob, _ := validateSSHPublicKey(sampleTestPubKey)
	existing := "# Some comment\nssh-rsa AAAAB3NzaC1yc2E... other@host\n" + sampleTestPubKey + "\n"
	hasKey := isKeyInAuthorizedKeys(existing, blob)
	if !hasKey {
		t.Errorf("expected key to be detected in authorized_keys")
	}
}

func TestIsKeyInAuthorizedKeys_NoMatch(t *testing.T) {
	_, blob, _ := validateSSHPublicKey(sampleTestPubKey)
	existing := "# Some comment\nssh-rsa AAAAB3NzaC1yc2E... other@host\n"
	hasKey := isKeyInAuthorizedKeys(existing, blob)
	if hasKey {
		t.Errorf("expected key NOT to be detected")
	}
}

func TestFormatKeyAppend(t *testing.T) {
	formattedEmpty := formatKeyAppend("", sampleTestPubKey)
	if formattedEmpty != sampleTestPubKey+"\n" {
		t.Errorf("expected key with newline, got: %q", formattedEmpty)
	}

	formattedWithNL := formatKeyAppend("key1\n", sampleTestPubKey)
	if formattedWithNL != "key1\n"+sampleTestPubKey+"\n" {
		t.Errorf("expected key appended, got: %q", formattedWithNL)
	}

	formattedNoNL := formatKeyAppend("key1", sampleTestPubKey)
	if formattedNoNL != "key1\n"+sampleTestPubKey+"\n" {
		t.Errorf("expected newline added before key, got: %q", formattedNoNL)
	}
}

func TestRunSSHAuthKeyAddCLI_NonInteractiveEmpty(t *testing.T) {
	err := RunSSHAuthKeyAddCLI([]string{})
	if err == nil {
		t.Fatalf("expected validation error when running empty args non-interactively")
	}
}

func TestRunSSHAuthKeyAddCLI_Help(t *testing.T) {
	err := RunSSHAuthKeyAddCLI([]string{"--help"})
	if err != nil {
		t.Fatalf("expected nil for help, got: %v", err)
	}
}
