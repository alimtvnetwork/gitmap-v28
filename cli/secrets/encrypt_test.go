package secrets

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("thisis32bitlongpassphraseimusing") // 32 bytes
	plaintext := []byte("supersecretpassword")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Expected %s, got %s", string(plaintext), string(decrypted))
	}
}

func TestDecrypt_InvalidKey(t *testing.T) {
	key := []byte("thisis32bitlongpassphraseimusing")
	wrongKey := []byte("thisis32bitlongpassphraseimwrong")
	plaintext := []byte("supersecretpassword")

	ciphertext, _ := Encrypt(plaintext, key)

	_, err := Decrypt(ciphertext, wrongKey)
	if err == nil {
		t.Error("Expected error with wrong key, got nil")
	}
}

func TestDecryptStoredPassword(t *testing.T) {
	rawSecret := "ClusterSecretPassword123!"
	cipherText, err := Encrypt([]byte(rawSecret), sshSecretKeyPrimary)
	if err != nil {
		t.Fatalf("Failed to encrypt with primary key: %v", err)
	}

	decrypted, decErr := DecryptStoredPassword(cipherText)
	if decErr != nil {
		t.Fatalf("DecryptStoredPassword failed: %v", decErr)
	}
	if decrypted != rawSecret {
		t.Errorf("Expected %q, got %q", rawSecret, decrypted)
	}

	empty, emptyErr := DecryptStoredPassword("")
	if emptyErr != nil || empty != "" {
		t.Errorf("Expected empty string, got %q, err=%v", empty, emptyErr)
	}

	fallback, _ := DecryptStoredPassword("plainPasswordNoEncryption")
	if fallback != "plainPasswordNoEncryption" {
		t.Errorf("Expected plaintext fallback, got %q", fallback)
	}
}

func TestConnectWithFallback_NoCredentials(t *testing.T) {
	_, err := ConnectWithFallback("127.0.0.1", "testuser", "", "")
	if err == nil {
		t.Error("Expected error when no credentials provided, got nil")
	}
}

func TestConnectWithFallback_InvalidKeyWithPassword(t *testing.T) {
	_, err := ConnectWithFallback("127.0.0.1:9999", "testuser", "non-existent-key-file", "secret")
	if err == nil {
		t.Error("Expected dial error to invalid port, got nil")
	}
}

func TestWrapCommandForShell_NoDoubleWrap(t *testing.T) {
	cases := []struct {
		cmd       string
		shellType string
		expected  string
	}{
		{"echo hi", "ps", "powershell -NoProfile -Command \"echo hi\""},
		{"powershell -NoProfile -Command \"echo hi\"", "ps", "powershell -NoProfile -Command \"echo hi\""},
		{"pwsh -c echo hi", "ps", "pwsh -c echo hi"},
		{"echo hi", "cmd", "cmd.exe /c \"echo hi\""},
		{"cmd.exe /c echo hi", "cmd", "cmd.exe /c echo hi"},
		{"echo hi", "bash", "bash -c \"echo hi\""},
		{"bash -c 'echo hi'", "bash", "bash -c 'echo hi'"},
	}

	for _, c := range cases {
		actual := wrapCommandForShell(c.cmd, c.shellType)
		if actual != c.expected {
			t.Errorf("wrapCommandForShell(%q, %q) = %q, want %q", c.cmd, c.shellType, actual, c.expected)
		}
	}
}
