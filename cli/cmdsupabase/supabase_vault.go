package cmdsupabase

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const (
	vaultSaltFileName = "vault.salt"
	saltByteLength    = 32
)

// resolveVaultSaltPath returns the canonical path to the installation vault salt.
func resolveVaultSaltPath() string {
	dbPath := store.ResolveSplitDbPath(store.SectionInstallation, "supabase", "")
	installDir := filepath.Dir(filepath.Dir(dbPath))

	return filepath.Clean(filepath.Join(installDir, vaultSaltFileName))
}

// readOrGenerateVaultSalt loads the vault salt from disk or generates a new one.
func readOrGenerateVaultSalt() ([]byte, error) {
	saltPath := resolveVaultSaltPath()
	if salt, err := os.ReadFile(saltPath); err == nil && len(salt) >= saltByteLength {
		return salt, nil
	}

	return generateAndSaveVaultSalt(saltPath)
}

func generateAndSaveVaultSalt(saltPath string) ([]byte, error) {
	newSalt := make([]byte, saltByteLength)
	if _, err := rand.Read(newSalt); err != nil {
		return nil, apperror.WrapSimple(err, "generateAndSaveVaultSalt_rand")
	}

	_ = os.MkdirAll(filepath.Dir(saltPath), 0700)
	if err := os.WriteFile(saltPath, newSalt, 0600); err != nil {
		return nil, apperror.WrapSimple(err, "generateAndSaveVaultSalt_writeFile")
	}

	return newSalt, nil
}

func getMachineFingerprint() string {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		return "gitmap-default-machine-fingerprint"
	}

	return hostname
}

// DeriveMasterVaultKey computes a 32-byte key using HMAC-SHA256 over machine identity and salt.
func DeriveMasterVaultKey() ([]byte, error) {
	salt, err := readOrGenerateVaultSalt()
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, salt)
	mac.Write([]byte(getMachineFingerprint()))

	return mac.Sum(nil), nil
}

// EncryptSecret encrypts a secret with AES-256-GCM using the master vault key.
func EncryptSecret(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key, err := DeriveMasterVaultKey()
	if err != nil {
		return "", err
	}

	ciphertext, err := secrets.Encrypt([]byte(plaintext), key)
	if err != nil {
		return "", apperror.WrapSimple(err, "EncryptSecret")
	}

	return ciphertext, nil
}

// DecryptSecret decrypts a secret using AES-256-GCM or SSH vault fallback.
func DecryptSecret(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	if strings.HasPrefix(ciphertext, "rsa:") || strings.HasPrefix(ciphertext, "aes:") {
		return cmdssh.DecryptSSHPassword(ciphertext)
	}

	key, err := DeriveMasterVaultKey()
	if err != nil {
		return "", err
	}

	plainBytes, err := secrets.Decrypt(ciphertext, key)
	if err == nil {
		return string(plainBytes), nil
	}

	return cmdssh.DecryptSSHPassword(ciphertext)
}

// MaskSecret masks sensitive credentials for safe display in terminal outputs.
func MaskSecret(secret string) string {
	if len(secret) == 0 {
		return "none"
	}

	if len(secret) <= 8 {
		return "****"
	}

	if len(secret) <= 12 {
		return secret[:4] + "...****"
	}

	return secret[:6] + "...****"
}
