package cmdagy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CachedSequenceEntry represents a single cached item row with its assigned sequence.
type CachedSequenceEntry struct {
	Seq      int    `json:"seq"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Messages int    `json:"messages"`
	Status   string `json:"status"`
	Queued   int    `json:"queued"`
	Preview  string `json:"preview,omitempty"`
	Category string `json:"category"`
}

// CachedSequenceManifest stores cached sequence mappings with a timestamp.
type CachedSequenceManifest struct {
	CreatedAt  time.Time             `json:"createdAt"`
	TTLSeconds int                   `json:"ttlSeconds"`
	Entries    []CachedSequenceEntry `json:"entries"`
}

func getSequenceCacheFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "user home dir")
	}

	dir := filepath.Join(home, ".gitmap", "cache")
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return "", apperror.WrapSimple(mkErr, "create cache dir")
	}

	return filepath.Join(dir, "agy_sequence_cache.json"), nil
}

func getSequenceCacheTTLSeconds() int {
	cfgPath, err := getAgyConfigPath()
	if err != nil {
		return constants.DefaultSequenceTTLSeconds
	}

	data, readErr := os.ReadFile(cfgPath)
	if readErr != nil {
		return constants.DefaultSequenceTTLSeconds
	}

	var m map[string]interface{}
	if jsonErr := json.Unmarshal(data, &m); jsonErr != nil {
		return constants.DefaultSequenceTTLSeconds
	}

	if ttl, ok := extractPositiveTTL(m, "agySequenceTTL"); ok {
		return ttl
	}

	if ttl, ok := extractPositiveTTL(m, "agy_sequence_ttl_seconds"); ok {
		return ttl
	}

	return constants.DefaultSequenceTTLSeconds
}

func extractPositiveTTL(m map[string]interface{}, key string) (int, bool) {
	val, hasKey := m[key]
	if !hasKey {
		return 0, false
	}
	fval, isFloat := val.(float64)
	if isFloat && fval > 0 {
		return int(fval), true
	}
	return 0, false
}

// SaveSequenceCache persists sequence rows with current timestamp and TTL.
func SaveSequenceCache(entries []CachedSequenceEntry) error {
	filePath, err := getSequenceCacheFilePath()
	if err != nil {
		return err
	}

	manifest := CachedSequenceManifest{
		CreatedAt:  time.Now().UTC(),
		TTLSeconds: getSequenceCacheTTLSeconds(),
		Entries:    entries,
	}

	bytes, mErr := json.MarshalIndent(manifest, "", "  ")
	if mErr != nil {
		return apperror.WrapSimple(mErr, "marshal sequence cache")
	}

	return os.WriteFile(filePath, bytes, constants.FilePermission)
}

// LoadSequenceCache loads cached sequences if still within the TTL window.
func LoadSequenceCache() (*CachedSequenceManifest, bool) {
	filePath, err := getSequenceCacheFilePath()
	if err != nil {
		return nil, false
	}

	data, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return nil, false
	}

	var manifest CachedSequenceManifest
	if jsonErr := json.Unmarshal(data, &manifest); jsonErr != nil {
		return nil, false
	}

	ttl := time.Duration(manifest.TTLSeconds) * time.Second
	hasValidTime := time.Since(manifest.CreatedAt) <= ttl
	if !hasValidTime {
		return nil, false
	}

	return &manifest, true
}

// ResolveCachedSequence retrieves a cached row by its sequence number.
func ResolveCachedSequence(seq int) (*CachedSequenceEntry, bool) {
	manifest, hasValidCache := LoadSequenceCache()
	if !hasValidCache || manifest == nil {
		return nil, false
	}

	for _, entry := range manifest.Entries {
		if entry.Seq == seq {
			return &entry, true
		}
	}

	return nil, false
}
