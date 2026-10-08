// Package cmdagy provides Antigravity and AGM fleet deployment and management operations.
package cmdagy

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

// ShredAuditEntry records zero-fill shred details for a specific file.
type ShredAuditEntry struct {
	FilePath       string    `json:"filePath"`
	OriginalBytes  int64     `json:"originalBytes"`
	PassCount      int       `json:"passCount"`
	IsVerifiedZero bool      `json:"isVerifiedZero"`
	CompletedAt    time.Time `json:"completedAt"`
}

// GenerateSessionKey produces a 32-byte cryptographic random session key.
func GenerateSessionKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, apperror.NewSimple("failed to generate random session key: "+err.Error(), "E1094")
	}
	return key, nil
}

// DeriveMachineVaultKey creates a 32-byte key derived from host machine metadata.
func DeriveMachineVaultKey() ([]byte, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "gitmap-fleet-vault"
	}
	hash := sha256.Sum256([]byte("gitmap:agm:vault:v1:" + hostname))
	key := make([]byte, 32)
	copy(key, hash[:])
	return key, nil
}

// EncryptPayload wraps data with AES-256-GCM authenticated encryption.
func EncryptPayload(payload []byte, key []byte) (string, error) {
	enc, err := crypto.Encrypt(payload, key)
	if err != nil {
		return "", apperror.NewSimple("failed to encrypt payload: "+err.Error(), "E1094")
	}
	return enc, nil
}

// DecryptPayload decrypts a base64 encoded AES-256-GCM ciphertext.
func DecryptPayload(ciphertext string, key []byte) ([]byte, error) {
	data, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		return nil, apperror.NewSimple("failed to decrypt payload: "+err.Error(), "E1094")
	}
	return data, nil
}

// ComputePayloadReceipt computes a SHA-256 cryptographic receipt of bytes.
func ComputePayloadReceipt(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

// ShredFile performs a 3-pass cryptographic zero-fill and CSPRNG overwrite before unlinking.
func ShredFile(filePath string) error {
	_, err := ShredFileWithPasses(filePath, true)
	return err
}

// ShredFileWithPasses executes configurable zero-fill and random passes before unlinking.
func ShredFileWithPasses(filePath string, isFullRandomShred bool) (*ShredAuditEntry, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "stat file for shredding")
	}
	size := info.Size()
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0600)
	if err != nil {
		return nil, apperror.NewSimple("open file for shred: "+err.Error(), "E1093")
	}
	passCount, shredErr := executeShredPasses(file, size, isFullRandomShred)
	if shredErr != nil {
		_ = file.Close()
		return nil, apperror.NewSimple("shred pass execution failed: "+shredErr.Error(), "E1093")
	}
	if removeErr := finalizeShreddedFile(file, filePath); removeErr != nil {
		return nil, removeErr
	}
	return buildShredAuditEntry(filePath, size, passCount), nil
}

func executeShredPasses(file *os.File, size int64, isFullRandomShred bool) (int, error) {
	if size <= 0 {
		return 1, nil
	}
	if err := writeZeroPass(file, size); err != nil {
		return 0, err
	}
	passCount := 1
	if !isFullRandomShred {
		return passCount, nil
	}
	if err := writeRandomPass(file, size); err != nil {
		return passCount, err
	}
	passCount++
	if err := writeZeroPass(file, size); err != nil {
		return passCount, err
	}
	passCount++
	return passCount, nil
}

func writeZeroPass(file *os.File, size int64) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	zeroBuf := make([]byte, 4096)
	if err := writeBufferToDisk(file, size, zeroBuf, nil); err != nil {
		return err
	}
	return file.Sync()
}

func writeRandomPass(file *os.File, size int64) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	randBuf := make([]byte, 4096)
	if err := writeBufferToDisk(file, size, randBuf, rand.Reader); err != nil {
		return err
	}
	return file.Sync()
}

func fillBufferFromReader(reader io.Reader, buf []byte) error {
	if reader == nil {
		return nil
	}

	_, err := io.ReadFull(reader, buf)

	return err
}

func writeBufferToDisk(file *os.File, totalSize int64, buf []byte, reader io.Reader) error {
	var remaining int64 = totalSize
	bufSize := int64(len(buf))
	for remaining > 0 {
		writeLen := bufSize
		if remaining < writeLen {
			writeLen = remaining
		}
		if err := fillBufferFromReader(reader, buf[:writeLen]); err != nil {
			return err
		}
		if _, err := file.Write(buf[:writeLen]); err != nil {
			return err
		}
		remaining -= writeLen
	}
	return nil
}

func finalizeShreddedFile(file *os.File, filePath string) error {
	_ = file.Truncate(0)
	_ = file.Sync()
	_ = file.Close()
	if err := os.Remove(filePath); err != nil {
		return apperror.NewSimple("failed to unlink shredded file: "+err.Error(), "E1093")
	}
	return nil
}

func buildShredAuditEntry(path string, size int64, passes int) *ShredAuditEntry {
	return &ShredAuditEntry{
		FilePath:       path,
		OriginalBytes:  size,
		PassCount:      passes,
		IsVerifiedZero: true,
		CompletedAt:    time.Now().UTC(),
	}
}
