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

func loadOrEnsureRSAPrivateKey() (*rsa.PrivateKey, error) {
	key, err := EnsureVaultRSAKeyPair()
	if err == nil {
		return key, nil
	}

	return loadLocalSSHRSAKey()
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

// EncryptWithRSAPublicKey encrypts plaintext using RSA-OAEP with SHA-256 and hardware random seed.
func EncryptWithRSAPublicKey(pub *rsa.PublicKey, plain string) (string, error) {
	cipherBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, []byte(plain), nil)
	if err != nil {
		return "", apperror.WrapSimple(err, "EncryptWithRSAPublicKey")
	}

	return prefixRSA + base64.StdEncoding.EncodeToString(cipherBytes), nil
}

func encryptWithRSAPublicKey(pub *rsa.PublicKey, plain string) (string, error) {
	return EncryptWithRSAPublicKey(pub, plain)
}

// DecryptWithRSAPrivateKey decrypts base64-encoded RSA-OAEP SHA-256 ciphertext using private key.
func DecryptWithRSAPrivateKey(priv *rsa.PrivateKey, cipherText string) (string, error) {
	rawBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", apperror.WrapSimple(err, "DecryptWithRSAPrivateKey_Decode")
	}

	plainBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, rawBytes, nil)
	if err != nil {
		return "", apperror.WrapSimple(err, "DecryptWithRSAPrivateKey_OAEP")
	}

	return string(plainBytes), nil
}

func decryptWithRSAPrivateKey(priv *rsa.PrivateKey, cipherText string) (string, error) {
	return DecryptWithRSAPrivateKey(priv, cipherText)
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

// EncryptSSHPassword encrypts a plaintext password using RSA-OAEP with the dedicated vault RSA key.
func EncryptSSHPassword(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}

	privKey, err := loadOrEnsureRSAPrivateKey()
	if err == nil {
		return encryptWithRSAPublicKey(&privKey.PublicKey, plain)
	}

	return encryptWithFallbackAES(plain)
}

// EncryptSaltedPassword encrypts a password using the salted rotation cipher.
func EncryptSaltedPassword(plain, salt string) string {
	return crypto.EncryptSalted(plain, salt)
}

// EncryptCaesarPassword encrypts a password using a caesar rotation cipher with optional salt.
func EncryptCaesarPassword(plain string, shift int, salt string) string {
	return crypto.EncryptCaesar(plain, shift, salt)
}

func isRSACiphertext(cipherText string) bool {
	return strings.HasPrefix(cipherText, prefixRSA)
}

func isAESCiphertext(cipherText string) bool {
	return strings.HasPrefix(cipherText, prefixAES)
}

// DecryptSSHPassword decrypts an encrypted password using RSA, AES, or legacy ciphers.
func DecryptSSHPassword(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}

	if plain, err, hasPrefix := tryDecryptPrefixedPassword(cipherText); hasPrefix {
		return plain, err
	}

	return decryptStoredFallback(cipherText)
}

func tryDecryptLegacyPrefix(cipherText string) (string, error, bool) {
	if strings.HasPrefix(cipherText, "salt:") {
		plain, err := crypto.DecryptSalted(cipherText)
		return plain, err, true
	}

	if strings.HasPrefix(cipherText, "caesar:") {
		plain, err := crypto.DecryptCaesar(cipherText)
		return plain, err, true
	}

	return "", nil, false
}

func tryDecryptStandardPrefix(cipherText string) (string, error, bool) {
	if isRSACiphertext(cipherText) {
		plain, err := decryptRSAPassword(strings.TrimPrefix(cipherText, prefixRSA))
		return plain, err, true
	}

	if isAESCiphertext(cipherText) {
		plain, err := decryptWithFallbackAES(strings.TrimPrefix(cipherText, prefixAES))
		return plain, err, true
	}

	return "", nil, false
}

func tryDecryptPrefixedPassword(cipherText string) (string, error, bool) {
	if plain, err, isLegacy := tryDecryptLegacyPrefix(cipherText); isLegacy {
		return plain, err, true
	}

	return tryDecryptStandardPrefix(cipherText)
}

func decryptStoredFallback(cipherText string) (string, error) {
	plain, decErr := crypto.DecryptStoredPassword(cipherText)
	if decErr == nil && plain != cipherText {
		return plain, nil
	}

	return decryptWithFallbackAES(cipherText)
}

func tryDecryptWithVaultKey(body string) (string, bool) {
	vaultKey, err := LoadVaultRSAPrivateKey()
	if err != nil {
		return "", false
	}

	plain, decErr := decryptWithRSAPrivateKey(vaultKey, body)

	return plain, decErr == nil
}

func decryptRSAPassword(body string) (string, error) {
	if plain, hasVaultDec := tryDecryptWithVaultKey(body); hasVaultDec {
		return plain, nil
	}

	sshKey, sshErr := loadLocalSSHRSAKey()
	if sshErr != nil {
		return "", sshErr
	}

	return decryptWithRSAPrivateKey(sshKey, body)
}
