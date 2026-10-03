package cmdssh

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func setupTempVaultDir(t *testing.T) string {
	tempDir := t.TempDir()
	prevLocator := vaultKeysDirLocator
	vaultKeysDirLocator = func() (string, error) {
		return tempDir, nil
	}
	t.Cleanup(func() {
		vaultKeysDirLocator = prevLocator
	})

	return tempDir
}

func setupDisabledSSHLocator(t *testing.T) {
	prevLocator := defaultSSHKeyLocator
	defaultSSHKeyLocator = func() (string, bool) {
		return "", false
	}
	t.Cleanup(func() {
		defaultSSHKeyLocator = prevLocator
	})
}

func disableVaultAndSSH(t *testing.T) {
	prevVault := vaultKeysDirLocator
	vaultKeysDirLocator = func() (string, error) {
		return "", apperror.NewExecutionError("vault disabled")
	}
	t.Cleanup(func() {
		vaultKeysDirLocator = prevVault
	})

	setupDisabledSSHLocator(t)
}

func TestRSAOAEP_EncryptDecrypt(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	secret := "MyP@ssw0rd!2026"
	cipherText := assertEncryptRSA(t, &privKey.PublicKey, secret)
	assertDecryptRSA(t, privKey, cipherText, secret)
}

func assertEncryptRSA(t *testing.T, pub *rsa.PublicKey, secret string) string {
	cipherText, err := encryptWithRSAPublicKey(pub, secret)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if !isRSACiphertext(cipherText) {
		t.Fatalf("expected rsa: prefix on ciphertext, got %s", cipherText)
	}

	return cipherText
}

func assertDecryptRSA(t *testing.T, priv *rsa.PrivateKey, cipherText, want string) {
	plainText, err := decryptWithRSAPrivateKey(priv, cipherText[4:])
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if plainText != want {
		t.Fatalf("plaintext mismatch: got %q, want %q", plainText, want)
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

	assertDecryptFallbackAES(t, cipherText[4:], secret)
}

func assertDecryptFallbackAES(t *testing.T, body, want string) {
	plainText, err := decryptWithFallbackAES(body)
	if err != nil {
		t.Fatalf("fallback decryption failed: %v", err)
	}

	if plainText != want {
		t.Fatalf("plaintext mismatch: got %q, want %q", plainText, want)
	}
}

func TestEncryptDecryptSSHPassword_Fallback(t *testing.T) {
	disableVaultAndSSH(t)
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

func TestEnsureVaultRSAKeyPair_GenerationAndPermissions(t *testing.T) {
	tempDir := setupTempVaultDir(t)
	key1, err := EnsureVaultRSAKeyPair()
	if err != nil {
		t.Fatalf("EnsureVaultRSAKeyPair failed: %v", err)
	}

	assertVaultKeyFilesExist(t, tempDir)
	key2, err := EnsureVaultRSAKeyPair()
	if err != nil {
		t.Fatalf("second EnsureVaultRSAKeyPair failed: %v", err)
	}

	assertKeysMatch(t, key1, key2)
}

func assertVaultKeyFilesExist(t *testing.T, tempDir string) {
	privPath := filepath.Join(tempDir, "vault_rsa")
	pubPath := filepath.Join(tempDir, "vault_rsa.pub")
	if _, err := os.Stat(privPath); err != nil {
		t.Fatalf("vault_rsa not created: %v", err)
	}

	if _, err := os.Stat(pubPath); err != nil {
		t.Fatalf("vault_rsa.pub not created: %v", err)
	}

	if !hasVaultRSAKey() {
		t.Fatalf("expected hasVaultRSAKey() to be true")
	}
}

func assertKeysMatch(t *testing.T, k1, k2 *rsa.PrivateKey) {
	if k1.N.Cmp(k2.N) != 0 {
		t.Fatalf("keys differ on repeated Ensure call")
	}
}

func TestRSAOAEP_NonDeterministicSaltCiphertext(t *testing.T) {
	setupTempVaultDir(t)
	password := "MyP@ssw0rd!2026"

	cipher1, err1 := EncryptSSHPassword(password)
	cipher2, err2 := EncryptSSHPassword(password)
	if err1 != nil || err2 != nil {
		t.Fatalf("encryption failed: %v, %v", err1, err2)
	}

	assertCipherMismatch(t, cipher1, cipher2)
	assertDecryptedMatches(t, cipher1, password)
	assertDecryptedMatches(t, cipher2, password)
}

func assertCipherMismatch(t *testing.T, c1, c2 string) {
	if c1 == c2 {
		t.Fatalf("expected non-deterministic ciphertexts, but both are identical: %s", c1)
	}
}

func assertDecryptedMatches(t *testing.T, cipher, expected string) {
	dec, err := DecryptSSHPassword(cipher)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if dec != expected {
		t.Fatalf("decrypted mismatch: got %q, want %q", dec, expected)
	}
}

func TestVaultRSA_EncryptDecryptRoundtrip(t *testing.T) {
	setupTempVaultDir(t)
	passwords := []string{
		"SimplePass123",
		"P@$$w0rd with spaces & $peci@l ch@rs! 🚀",
		"Multi---Line\n\tPassword==#",
		"VeryLongPassphraseThatExceedsNormalLengthForSSHKeys1234567890!@#$%^&*()_+",
	}

	for _, pass := range passwords {
		assertRoundtrip(t, pass)
	}
}

func assertRoundtrip(t *testing.T, pass string) {
	cipher, err := EncryptSSHPassword(pass)
	if err != nil {
		t.Fatalf("EncryptSSHPassword failed for %q: %v", pass, err)
	}

	assertDecryptedMatches(t, cipher, pass)
}

func writeTempSSHKey(t *testing.T) string {
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "id_rsa")
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	pemBytes := encodeRSAPrivateKeyPEM(key)
	if err := os.WriteFile(keyPath, pemBytes, 0600); err != nil {
		t.Fatalf("failed to write id_rsa: %v", err)
	}

	return keyPath
}

func TestVaultRSA_FallbackToSSHRSA(t *testing.T) {
	keyPath := writeTempSSHKey(t)
	prevSSH := defaultSSHKeyLocator
	defaultSSHKeyLocator = func() (string, bool) {
		return keyPath, true
	}
	t.Cleanup(func() {
		defaultSSHKeyLocator = prevSSH
	})

	prevVault := vaultKeysDirLocator
	vaultKeysDirLocator = func() (string, error) {
		return "", apperror.NewExecutionError("simulated vault failure")
	}
	t.Cleanup(func() {
		vaultKeysDirLocator = prevVault
	})

	assertRoundtrip(t, "FallbackToSSHKeyPass")
}

func TestVaultRSA_FallbackToAES(t *testing.T) {
	disableVaultAndSSH(t)
	secret := "AESFallbackOnlyPass"
	cipher, err := EncryptSSHPassword(secret)
	if err != nil {
		t.Fatalf("EncryptSSHPassword failed: %v", err)
	}

	if !isAESCiphertext(cipher) {
		t.Fatalf("expected aes: prefix, got %s", cipher)
	}

	assertDecryptedMatches(t, cipher, secret)
}

func TestDecryptSSHPassword_LegacyCompatibility(t *testing.T) {
	setupTempVaultDir(t)
	testDecryptSalted(t)
	testDecryptCaesar(t)
	testDecryptEmpty(t)
}

func testDecryptSalted(t *testing.T) {
	plain := "SaltedSecret123"
	saltedCipher := EncryptSaltedPassword(plain, "mysalt")
	assertDecryptedMatches(t, saltedCipher, plain)
}

func testDecryptCaesar(t *testing.T) {
	plain := "CaesarSecret456"
	caesarCipher := EncryptCaesarPassword(plain, 5, "caesarsalt")
	assertDecryptedMatches(t, caesarCipher, plain)
}

func testDecryptEmpty(t *testing.T) {
	dec, err := DecryptSSHPassword("")
	if err != nil || dec != "" {
		t.Fatalf("expected empty result, got %q, %v", dec, err)
	}
}
