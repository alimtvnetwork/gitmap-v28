package cmdssh

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const defaultVaultSalt = "gitmap-ssh-vault-salt-v1"

func deriveVaultMasterKey() []byte {
	hostname, err := os.Hostname()
	hasErr := err != nil || hostname == ""
	if hasErr {
		hostname = "localhost"
	}

	mac := hmac.New(sha256.New, []byte(defaultVaultSalt))
	mac.Write([]byte(hostname))

	return mac.Sum(nil)
}

func createGCMCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, apperror.WrapSimple(err, "createGCMCipher.NewCipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, apperror.WrapSimple(err, "createGCMCipher.NewGCM")
	}

	return gcm, nil
}

// EncryptPrivateKey encrypts a private key PEM using AES-256-GCM.
func EncryptPrivateKey(plainPEM []byte) (string, error) {
	key := deriveVaultMasterKey()
	gcm, err := createGCMCipher(key)
	if err != nil {
		return "", err
	}

	return sealPEMWithGCM(gcm, plainPEM)
}

func sealPEMWithGCM(gcm cipher.AEAD, plainPEM []byte) (string, error) {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", apperror.WrapSimple(err, "EncryptPrivateKey.ReadNonce")
	}

	sealed := gcm.Seal(nonce, nonce, plainPEM, nil)

	return base64.StdEncoding.EncodeToString(sealed), nil
}

// DecryptPrivateKey decrypts an AES-256-GCM base64 ciphertext to private key PEM.
func DecryptPrivateKey(cipherB64 string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return nil, apperror.WrapSimple(err, "DecryptPrivateKey.Decode")
	}

	key := deriveVaultMasterKey()
	gcm, err := createGCMCipher(key)
	if err != nil {
		return nil, err
	}

	return decryptWithGCM(gcm, raw)
}

func decryptWithGCM(gcm cipher.AEAD, raw []byte) ([]byte, error) {
	nonceSize := gcm.NonceSize()
	isShort := len(raw) < nonceSize
	if isShort {
		return nil, apperror.NewValidationError("ciphertext too short")
	}

	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, apperror.WrapSimple(err, "DecryptPrivateKey.Open")
	}

	return plain, nil
}

// StoreEncryptedSSHKey encrypts and stores an SSH private key in SQLite.
func StoreEncryptedSSHKey(name, privateKeyPEM, publicKey, email string) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "StoreEncryptedSSHKey.openDB")
	}

	defer dbConn.Close()

	return StoreEncryptedSSHKeyWithDB(dbConn, name, privateKeyPEM, publicKey, email)
}

// StoreEncryptedSSHKeyWithDB persists an encrypted SSH key into the provided DB.
func StoreEncryptedSSHKeyWithDB(db *store.DB, name, privateKeyPEM, publicKey, email string) error {
	encPriv, err := EncryptPrivateKey([]byte(privateKeyPEM))
	if err != nil {
		return err
	}

	return saveEncryptedKeyRecord(db, name, encPriv, publicKey, email)
}

func saveEncryptedKeyRecord(db *store.DB, name, encPriv, publicKey, email string) error {
	privPath := resolveDefaultKeyPath(name)
	fp := computePublicKeyFingerprint(publicKey)
	hasExisting := db.SSHKeyExists(name)
	if hasExisting {
		return db.UpdateSSHKey(name, privPath, publicKey, fp, email, encPriv)
	}

	_, err := db.InsertSSHKey(name, privPath, publicKey, fp, email, encPriv)

	return err
}

func resolveDefaultKeyPath(name string) string {
	home, err := os.UserHomeDir()
	hasErr := err != nil
	if hasErr {
		return filepath.Join(".ssh", name)
	}

	return filepath.Join(home, ".ssh", name)
}

func computePublicKeyFingerprint(pubKeyStr string) string {
	pubBytes := []byte(strings.TrimSpace(pubKeyStr))
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(pubBytes)
	hasErr := err != nil
	if hasErr {
		return "SHA256:unknown"
	}

	return ssh.FingerprintSHA256(parsed)
}

// RetrieveDecryptedSSHKey retrieves and decrypts an SSH private key from SQLite.
func RetrieveDecryptedSSHKey(name string) ([]byte, error) {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return nil, apperror.WrapSimple(err, "RetrieveDecryptedSSHKey.openDB")
	}

	defer dbConn.Close()

	return RetrieveDecryptedSSHKeyWithDB(dbConn, name)
}

// RetrieveDecryptedSSHKeyWithDB fetches and decrypts an SSH key from the given DB.
func RetrieveDecryptedSSHKeyWithDB(db *store.DB, name string) ([]byte, error) {
	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		return nil, apperror.WrapSimple(err, "RetrieveDecryptedSSHKeyWithDB.Find")
	}

	hasEnc := key.EncryptedPrivateKey != ""
	if hasEnc {
		return DecryptPrivateKey(key.EncryptedPrivateKey)
	}

	return readPlainDiskFallback(key.PrivatePath)
}

func readPlainDiskFallback(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "readPlainDiskFallback")
	}

	return data, nil
}
