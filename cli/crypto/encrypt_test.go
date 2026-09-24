package crypto

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
