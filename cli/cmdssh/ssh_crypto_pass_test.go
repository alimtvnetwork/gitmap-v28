package cmdssh

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestRSAOAEP_EncryptDecrypt(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	secret := "MyP@ssw0rd!2026"
	cipherText, err := encryptWithRSAPublicKey(&privKey.PublicKey, secret)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if !isRSACiphertext(cipherText) {
		t.Fatalf("expected rsa: prefix on ciphertext, got %s", cipherText)
	}

	plainText, err := decryptWithRSAPrivateKey(privKey, cipherText[4:])
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if plainText != secret {
		t.Fatalf("plaintext mismatch: got %q, want %q", plainText, secret)
	}
}

func TestFallbackAES_EncryptDecrypt(t *testing.T) {
	secret := "SecretFallbackPassword"
	cipherText, err := encryptWithFallbackAES(secret)
	if err != nil {
		t.Fatalf("fallback encryption failed: %v", err)
	}

	if !isAESCiphertext(cipherText) {
		t.Fatalf("expected aes: prefix on ciphertext, got %s", cipherText)
	}

	plainText, err := decryptWithFallbackAES(cipherText[4:])
	if err != nil {
		t.Fatalf("fallback decryption failed: %v", err)
	}

	if plainText != secret {
		t.Fatalf("plaintext mismatch: got %q, want %q", plainText, secret)
	}
}

func TestEncryptDecryptSSHPassword_Fallback(t *testing.T) {
	prevLocator := defaultSSHKeyLocator
	defaultSSHKeyLocator = func() (string, bool) { return "", false }
	defer func() { defaultSSHKeyLocator = prevLocator }()

	secret := "NoKeyPassword123"
	cipherText, err := EncryptSSHPassword(secret)
	if err != nil {
		t.Fatalf("EncryptSSHPassword failed: %v", err)
	}

	plainText, err := DecryptSSHPassword(cipherText)
	if err != nil {
		t.Fatalf("DecryptSSHPassword failed: %v", err)
	}

	if plainText != secret {
		t.Fatalf("plaintext mismatch: got %q, want %q", plainText, secret)
	}
}
