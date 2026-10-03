package cmdssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var vaultKeysDirLocator = func() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "resolveVaultKeysDir")
	}

	return filepath.Join(home, ".gitmap", "keys"), nil
}

func resolveVaultKeysDir() (string, error) {
	return vaultKeysDirLocator()
}

// ResolveVaultKeyPaths returns the paths to the vault private and public key files.
func ResolveVaultKeyPaths() (string, string, error) {
	dir, err := resolveVaultKeysDir()
	if err != nil {
		return "", "", err
	}

	return filepath.Join(dir, "vault_rsa"), filepath.Join(dir, "vault_rsa.pub"), nil
}

func resolveVaultKeyPaths() (string, string, error) {
	return ResolveVaultKeyPaths()
}

// HasVaultRSAKey checks whether the vault RSA private key file currently exists.
func HasVaultRSAKey() bool {
	privPath, _, err := resolveVaultKeyPaths()
	if err != nil {
		return false
	}

	_, statErr := os.Stat(privPath)

	return statErr == nil
}

func hasVaultRSAKey() bool {
	return HasVaultRSAKey()
}

func provisionVaultKeysDir() (string, string, error) {
	dir, err := resolveVaultKeysDir()
	if err != nil {
		return "", "", err
	}

	if mkdirErr := os.MkdirAll(dir, 0700); mkdirErr != nil {
		return "", "", apperror.WrapSimple(mkdirErr, "provisionVaultKeysDir")
	}

	return resolveVaultKeyPaths()
}

// EnsureVaultRSAKeyPair loads the vault RSA keypair, generating it if absent.
func EnsureVaultRSAKeyPair() (*rsa.PrivateKey, error) {
	if hasVaultRSAKey() {
		return LoadVaultRSAPrivateKey()
	}

	privPath, pubPath, err := provisionVaultKeysDir()
	if err != nil {
		return nil, err
	}

	return generateAndSaveVaultRSAKeyPair(privPath, pubPath)
}

// LoadVaultRSAPrivateKey loads the existing vault private key from disk.
func LoadVaultRSAPrivateKey() (*rsa.PrivateKey, error) {
	privPath, _, err := resolveVaultKeyPaths()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(privPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "LoadVaultRSAPrivateKey_ReadFile")
	}

	return parseRSAPrivateKeyPEM(data)
}

func generateAndSaveVaultRSAKeyPair(privPath, pubPath string) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, apperror.WrapSimple(err, "generateAndSaveVaultRSAKeyPair")
	}

	if saveErr := saveVaultKeyPairFiles(privPath, pubPath, key); saveErr != nil {
		return nil, saveErr
	}

	return key, nil
}

func saveVaultKeyPairFiles(privPath, pubPath string, key *rsa.PrivateKey) error {
	privPEM := encodeRSAPrivateKeyPEM(key)
	if err := saveVaultFile(privPath, privPEM, 0600); err != nil {
		return err
	}

	pubPEM, err := encodeRSAPublicKeyPEM(&key.PublicKey)
	if err != nil {
		return err
	}

	return saveVaultFile(pubPath, pubPEM, 0644)
}

func encodeRSAPrivateKeyPEM(priv *rsa.PrivateKey) []byte {
	der := x509.MarshalPKCS1PrivateKey(priv)

	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	})
}

func encodeRSAPublicKeyPEM(pub *rsa.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, apperror.WrapSimple(err, "encodeRSAPublicKeyPEM")
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}), nil
}

func saveVaultFile(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return apperror.WrapSimple(err, "saveVaultFile")
	}

	return nil
}

func parseRSAPrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, apperror.NewValidationError("invalid PEM block for vault RSA key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}

	return parsePKCS8RSAPrivateKey(block.Bytes)
}

func parsePKCS8RSAPrivateKey(derBytes []byte) (*rsa.PrivateKey, error) {
	raw, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, apperror.WrapSimple(err, "parsePKCS8RSAPrivateKey")
	}

	key, isRSA := raw.(*rsa.PrivateKey)
	if isRSA {
		return key, nil
	}

	return nil, apperror.NewValidationError("vault key is not an RSA private key")
}
