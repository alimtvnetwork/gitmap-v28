package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

var (
	sshSecretKeyPrimary  = []byte("gitmap-ssh-secret-key-0123456789")
	sshSecretKeyFallback = []byte("gitmap-ssh-fallback-secret-01234")
)

// DecryptStoredPassword decrypts stored ciphertext into a plaintext password.
// If the input is not encrypted or decryption fails, it safely falls back to the original string.
func DecryptStoredPassword(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}
	if strings.HasPrefix(cipherText, "rsa:") {
		return decryptRSACipher(strings.TrimPrefix(cipherText, "rsa:"))
	}
	toDecrypt := strings.TrimPrefix(cipherText, "aes:")
	if plain, isOk := tryDecryptAES(toDecrypt, sshSecretKeyPrimary); isOk {
		return plain, nil
	}
	if plain, isOk := tryDecryptAES(toDecrypt, sshSecretKeyFallback); isOk {
		return plain, nil
	}
	if plain, err := decryptRSACipher(cipherText); err == nil && plain != "" {
		return plain, nil
	}
	return cipherText, nil
}

func tryDecryptAES(cipherText string, key []byte) (string, bool) {
	bytes, err := Decrypt(cipherText, key)
	if err != nil || len(bytes) == 0 {
		return "", false
	}
	return string(bytes), true
}

func decryptRSACipher(body string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	keyPath := filepath.Join(home, ".ssh", "id_rsa")
	keyBytes, readErr := os.ReadFile(keyPath)
	if readErr != nil {
		return "", readErr
	}
	raw, parseErr := ssh.ParseRawPrivateKey(keyBytes)
	if parseErr != nil {
		return "", parseErr
	}
	rsaKey, isRSA := raw.(*rsa.PrivateKey)
	if !isRSA {
		return "", nil
	}
	rawBytes, decErr := base64.StdEncoding.DecodeString(body)
	if decErr != nil {
		return "", decErr
	}
	plainBytes, errOAEP := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaKey, rawBytes, nil)
	if errOAEP != nil {
		return "", errOAEP
	}
	return string(plainBytes), nil
}
