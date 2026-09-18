package cmdagy

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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

// RecordSentErrorSignature updates the persistent store with a new dispatch entry.
func RecordSentErrorSignature(storePath, sig, repo string, runID uint64, sha, errHash string) error {
	store := LoadSentAgyErrorsStore(storePath)
	count := resolveNextSentCount(store, sig)
	store.Records[sig] = newSentRecord(sig, repo, runID, sha, errHash, count)

	return SaveSentAgyErrorsStore(storePath, store)
}
