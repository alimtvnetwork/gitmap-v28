package cmdagy

import (
	"crypto/sha256"
	"encoding/hex"
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
