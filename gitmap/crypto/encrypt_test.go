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
