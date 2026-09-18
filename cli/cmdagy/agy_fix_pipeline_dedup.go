package cmdagy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ComputeErrorSignature generates a unique fingerprint for a pipeline error batch.
func ComputeErrorSignature(repo string, runID uint64, sha, errorLogs string) (string, string) {
	hasher := sha256.New()
	hasher.Write([]byte(errorLogs))
	errHash := hex.EncodeToString(hasher.Sum(nil))

	if runID > 0 {
		return fmt.Sprintf("%s:%d", repo, runID), errHash
	}

	if len(sha) > 0 {
		return fmt.Sprintf("%s:%s", repo, sha), errHash
	}

	return fmt.Sprintf("%s:%s", repo, errHash[:16]), errHash
}

func sentAgyErrorsStorePath() string {
	rootDir := resolveProjectRootDir()
	primaryPath := filepath.Join(rootDir, ".gitmap", "pipeline", "sent_agy_errors.json")
	if isDirWritable(filepath.Dir(primaryPath)) {
		return primaryPath
	}

	return filepath.Join(rootDir, ".ai-memory", "temp", "sent_agy_errors.json")
}

func isDirWritable(dir string) bool {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false
	}

	return true
}

// LoadSentAgyErrorsStore loads historical sent error fingerprints from disk.
func LoadSentAgyErrorsStore(path string) SentAgyErrorsStore {
	store := SentAgyErrorsStore{Records: make(map[string]SentAgyErrorRecord)}
	content, err := os.ReadFile(path)
	if err != nil {
		return store
	}

	_ = json.Unmarshal(content, &store)
	if store.Records == nil {
		store.Records = make(map[string]SentAgyErrorRecord)
	}

	return store
}

// SaveSentAgyErrorsStore writes historical sent error fingerprints to disk.
func SaveSentAgyErrorsStore(path string, store SentAgyErrorsStore) error {
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	marshaled, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, marshaled, 0644)
}

// CheckSentErrorDuplicate checks if errors have been dispatched previously.
func CheckSentErrorDuplicate(sig string, store SentAgyErrorsStore, isForce bool) (bool, *SentAgyErrorRecord) {
	existing, hasRecord := store.Records[sig]
	if !hasRecord {
		return false, nil
	}

	if isForce {
		return false, &existing
	}

	return true, &existing
}

func newSentRecord(sig, repo string, runID uint64, sha, errHash string, count int) SentAgyErrorRecord {
	return SentAgyErrorRecord{
		Signature: sig,
		Repo:      repo,
		RunID:     runID,
		SHA:       sha,
		ErrorHash: errHash,
		SentAt:    time.Now().UTC().Format(time.RFC3339),
		SentCount: count,
	}
}

func resolveNextSentCount(store SentAgyErrorsStore, sig string) int {
	if rec, hasExisting := store.Records[sig]; hasExisting {
		return rec.SentCount + 1
	}

	return 1
}

// RecordSentErrorSignature updates the persistent store with a new dispatch entry.
func RecordSentErrorSignature(storePath, sig, repo string, runID uint64, sha, errHash string) error {
	store := LoadSentAgyErrorsStore(storePath)
	count := resolveNextSentCount(store, sig)
	store.Records[sig] = newSentRecord(sig, repo, runID, sha, errHash, count)

	return SaveSentAgyErrorsStore(storePath, store)
}
