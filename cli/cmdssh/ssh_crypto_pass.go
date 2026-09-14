package cmdssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

const (
	prefixRSA = "rsa:"
	prefixAES = "aes:"
)

var defaultSSHKeyLocator = func() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	keyPath := filepath.Join(home, ".ssh", "id_rsa")
	if _, statErr := os.Stat(keyPath); statErr == nil {
		return keyPath, true
	}
	return "", false
}

func loadLocalSSHRSAKey() (*rsa.PrivateKey, error) {
	keyPath, hasKey := defaultSSHKeyLocator()
	if !hasKey {
		return nil, apperror.NewNotFoundError("local SSH RSA key (~/.ssh/id_rsa) not found")
	}
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "loadLocalSSHRSAKey_ReadFile")
	}
	return parseRawRSAPrivateKey(keyBytes)
}

func parseRawRSAPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	raw, err := ssh.ParseRawPrivateKey(keyBytes)
	if err != nil {
		return nil, apperror.WrapSimple(err, "parseRawRSAPrivateKey")
	}
	rsaKey, isRSA := raw.(*rsa.PrivateKey)
	if isRSA {
		return rsaKey, nil
	}
	return nil, apperror.NewValidationError("local SSH key is not an RSA key")
}

func encryptWithRSAPublicKey(pub *rsa.PublicKey, plain string) (string, error) {
	cipherBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, []byte(plain), nil)
	if err != nil {
		return "", apperror.WrapSimple(err, "encryptWithRSAPublicKey")
	}
	return prefixRSA + base64.StdEncoding.EncodeToString(cipherBytes), nil
}

func decryptWithRSAPrivateKey(priv *rsa.PrivateKey, cipherText string) (string, error) {
	rawBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", apperror.WrapSimple(err, "decryptWithRSAPrivateKey_Decode")
	}
	plainBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, rawBytes, nil)
	if err != nil {
		return "", apperror.WrapSimple(err, "decryptWithRSAPrivateKey_OAEP")
	}
	return string(plainBytes), nil
}

func deriveFallbackKey() []byte {
	return []byte("gitmap-ssh-fallback-secret-01234")
}

func encryptWithFallbackAES(plain string) (string, error) {
	enc, err := crypto.Encrypt([]byte(plain), deriveFallbackKey())
	if err != nil {
		return "", apperror.WrapSimple(err, "encryptWithFallbackAES")
	}
	return prefixAES + enc, nil
}

func decryptWithFallbackAES(enc string) (string, error) {
	plainBytes, err := crypto.Decrypt(enc, deriveFallbackKey())
	if err != nil {
		return "", apperror.WrapSimple(err, "decryptWithFallbackAES")
	}
	return string(plainBytes), nil
}

// EncryptSSHPassword encrypts a plaintext password using RSA-OAEP with the user's SSH key.
func EncryptSSHPassword(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	privKey, err := loadLocalSSHRSAKey()
	if err == nil {
		return encryptWithRSAPublicKey(&privKey.PublicKey, plain)
	}
	return encryptWithFallbackAES(plain)
}

func isRSACiphertext(cipherText string) bool {
	return strings.HasPrefix(cipherText, prefixRSA)
}

func isAESCiphertext(cipherText string) bool {
	return strings.HasPrefix(cipherText, prefixAES)
}

// DecryptSSHPassword decrypts an encrypted password using the user's SSH RSA key.
func DecryptSSHPassword(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}
	if isRSACiphertext(cipherText) {
		return decryptRSAPassword(strings.TrimPrefix(cipherText, prefixRSA))
	}
	if isAESCiphertext(cipherText) {
		return decryptWithFallbackAES(strings.TrimPrefix(cipherText, prefixAES))
	}
	return decryptWithFallbackAES(cipherText)
}

func decryptRSAPassword(body string) (string, error) {
	privKey, err := loadLocalSSHRSAKey()
	if err != nil {
		return "", err
	}
	return decryptWithRSAPrivateKey(privKey, body)
}
