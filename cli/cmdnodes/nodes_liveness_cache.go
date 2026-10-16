package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	NodesLivenessTTL       = 10 * time.Second
	NodeLivenessDefaultTTL = NodesLivenessTTL
	livenessCacheFileName  = "gitmap_nodes_liveness_cache.json"
)

type NodeLivenessCacheEntry struct {
	Alias        string `json:"alias,omitempty"`
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	Status       string `json:"status"`
	IsOnline     bool   `json:"is_online"`
	ProbedAtUnix int64  `json:"probed_at_unix"`
}

type NodeLivenessEntry = NodeLivenessCacheEntry

type NodesLivenessEnvelope struct {
	CachedAtUnix int64                             `json:"cached_at_unix"`
	TTLSeconds   int                               `json:"ttl_seconds"`
	Entries      map[string]NodeLivenessCacheEntry `json:"entries"`
}

type NodeLivenessCache = NodesLivenessEnvelope

var (
	livenessCacheMu  sync.RWMutex
	inMemoryEnvelope *NodesLivenessEnvelope
)

func resolveLivenessCacheFilePath() string {
	return filepath.Join(os.TempDir(), livenessCacheFileName)
}

func resolveLivenessCacheKey(host, alias string) string {
	trimmedHost := strings.TrimSpace(host)
	hasHost := trimmedHost != ""
	if hasHost {
		return strings.ToLower(trimmedHost)
	}
	return strings.ToLower(strings.TrimSpace(alias))
}

func isCacheEntryFresh(entry NodeLivenessCacheEntry, nowUnix int64, ttl time.Duration) bool {
	ageSeconds := nowUnix - entry.ProbedAtUnix
	maxAgeSeconds := int64(ttl.Seconds())
	isValidAge := ageSeconds >= 0
	isUnderTTL := ageSeconds < maxAgeSeconds
	return isValidAge && isUnderTTL
}

func findEntryByKeys(entries map[string]NodeLivenessCacheEntry, host, alias string) (NodeLivenessCacheEntry, bool) {
	hostKey := strings.ToLower(strings.TrimSpace(host))
	entry, hasHost := entries[hostKey]
	isHostMatch := hostKey != "" && hasHost
	if isHostMatch {
		return entry, true
	}

	aliasKey := strings.ToLower(strings.TrimSpace(alias))
	aliasEntry, hasAlias := entries[aliasKey]
	isAliasMatch := aliasKey != "" && aliasKey != "-" && hasAlias
	if isAliasMatch {
		return aliasEntry, true
	}

	return NodeLivenessCacheEntry{}, false
}

func lookupFreshEntry(env *NodesLivenessEnvelope, host, alias string) (NodeLivenessCacheEntry, bool) {
	isEnvNil := env == nil
	if isEnvNil {
		return NodeLivenessCacheEntry{}, false
	}
	isEntriesNil := env.Entries == nil
	if isEntriesNil {
		return NodeLivenessCacheEntry{}, false
	}

	nowUnix := time.Now().Unix()
	entry, hasEntry := findEntryByKeys(env.Entries, host, alias)
	isEntryMissing := !hasEntry
	if isEntryMissing {
		return NodeLivenessCacheEntry{}, false
	}

	isFresh := isCacheEntryFresh(entry, nowUnix, NodesLivenessTTL)
	if isFresh {
		return entry, true
	}

	return NodeLivenessCacheEntry{}, false
}

func loadLivenessEnvelopeFromDisk() *NodesLivenessEnvelope {
	cachePath := resolveLivenessCacheFilePath()
	data, err := os.ReadFile(cachePath)
	isReadErr := err != nil
	if isReadErr {
		return nil
	}

	var env NodesLivenessEnvelope
	errUnmarshal := json.Unmarshal(data, &env)
	isUnmarshalErr := errUnmarshal != nil
	if isUnmarshalErr {
		return nil
	}

	livenessCacheMu.Lock()
	inMemoryEnvelope = &env
	livenessCacheMu.Unlock()
	return &env
}

// GetCachedNodeLiveness retrieves a cached node liveness entry if present and fresh.
func GetCachedNodeLiveness(host, alias string) (NodeLivenessCacheEntry, bool) {
	isSkipEnvSet := os.Getenv("GITMAP_SKIP_NODE_CACHE") != ""
	if isSkipEnvSet {
		return NodeLivenessCacheEntry{}, false
	}

	livenessCacheMu.RLock()
	envelope := inMemoryEnvelope
	livenessCacheMu.RUnlock()

	isEnvelopeMissing := envelope == nil
	if isEnvelopeMissing {
		envelope = loadLivenessEnvelopeFromDisk()
	}

	isStillMissing := envelope == nil
	if isStillMissing {
		return NodeLivenessCacheEntry{}, false
	}

	return lookupFreshEntry(envelope, host, alias)
}

func readLivenessEnvelopeFromDiskUnsafe() *NodesLivenessEnvelope {
	cachePath := resolveLivenessCacheFilePath()
	data, err := os.ReadFile(cachePath)
	isReadErr := err != nil
	if isReadErr {
		return nil
	}

	var env NodesLivenessEnvelope
	errUnmarshal := json.Unmarshal(data, &env)
	isUnmarshalErr := errUnmarshal != nil
	if isUnmarshalErr {
		return nil
	}
	return &env
}

func fallbackDirectWrite(cachePath, tempPath string, data []byte) error {
	errDirect := os.WriteFile(cachePath, data, 0600)
	_ = os.Remove(tempPath)
	isDirectErr := errDirect != nil
	if isDirectErr {
		return apperror.WrapSimple(errDirect, "writeLivenessEnvelopeToDisk.directWrite")
	}
	return nil
}

func writeLivenessEnvelopeToDisk(envelope *NodesLivenessEnvelope) error {
	isEnvNil := envelope == nil
	if isEnvNil {
		return nil
	}

	data, err := json.MarshalIndent(envelope, "", "  ")
	isMarshalErr := err != nil
	if isMarshalErr {
		return apperror.WrapSimple(err, "writeLivenessEnvelopeToDisk.marshal")
	}

	cachePath := resolveLivenessCacheFilePath()
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("gitmap_nodes_liveness_%d_%d.tmp", os.Getpid(), time.Now().UnixNano()))

	errWrite := os.WriteFile(tempPath, data, 0600)
	isWriteErr := errWrite != nil
	if isWriteErr {
		return apperror.WrapSimple(errWrite, "writeLivenessEnvelopeToDisk.writeTemp")
	}

	_ = os.Remove(cachePath)
	errRename := os.Rename(tempPath, cachePath)
	isRenameErr := errRename != nil
	if isRenameErr {
		return fallbackDirectWrite(cachePath, tempPath, data)
	}
	return nil
}

// SetCachedNodesLiveness updates the cache in-memory and on disk.
func SetCachedNodesLiveness(entries map[string]NodeLivenessCacheEntry) {
	isEntriesEmpty := len(entries) == 0
	if isEntriesEmpty {
		return
	}

	livenessCacheMu.Lock()
	defer livenessCacheMu.Unlock()

	envelope := inMemoryEnvelope
	isEnvelopeMissing := envelope == nil
	if isEnvelopeMissing {
		envelope = readLivenessEnvelopeFromDiskUnsafe()
	}

	isStillNil := envelope == nil || envelope.Entries == nil
	if isStillNil {
		envelope = &NodesLivenessEnvelope{
			Entries: make(map[string]NodeLivenessCacheEntry),
		}
	}

	nowUnix := time.Now().Unix()
	envelope.CachedAtUnix = nowUnix
	envelope.TTLSeconds = int(NodesLivenessTTL.Seconds())

	for k, v := range entries {
		normalizedKey := strings.ToLower(strings.TrimSpace(k))
		isKeyEmpty := normalizedKey == ""
		if isKeyEmpty {
			continue
		}
		isTimeUnset := v.ProbedAtUnix == 0
		if isTimeUnset {
			v.ProbedAtUnix = nowUnix
		}
		envelope.Entries[normalizedKey] = v
	}

	inMemoryEnvelope = envelope
	_ = writeLivenessEnvelopeToDisk(envelope)
}

// BatchUpdateNodeLiveness updates entries in-memory and writes back to disk.
func BatchUpdateNodeLiveness(entries []NodeLivenessCacheEntry) error {
	m := make(map[string]NodeLivenessCacheEntry, len(entries))
	for _, e := range entries {
		k := resolveLivenessCacheKey(e.Host, e.Alias)
		m[k] = e
	}
	SetCachedNodesLiveness(m)
	return nil
}

// InvalidateNodesLivenessCache clears in-memory and disk liveness cache.
func InvalidateNodesLivenessCache() {
	livenessCacheMu.Lock()
	inMemoryEnvelope = nil
	livenessCacheMu.Unlock()

	cachePath := resolveLivenessCacheFilePath()
	_ = os.Remove(cachePath)
}

// InvalidateNodeLivenessCache is an alias for InvalidateNodesLivenessCache returning an error.
func InvalidateNodeLivenessCache() error {
	InvalidateNodesLivenessCache()
	return nil
}
