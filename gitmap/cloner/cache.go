// Package cloner — clone-cache.go
//
// Idempotent clone cache. Stores per-record fingerprints under
// <targetDir>/.gitmap/clone-cache.json so repeated `gitmap clone` runs can
// detect repos that are already cloned at the desired URL/branch and skip
// them when both the local HEAD and the remote tip match the cached values.
//
// Cache schema (versioned for forward compatibility):
//
//	{
//	  "version": 1,
//	  "entries": {
//	    "<relativePath>": {
//	      "url":       "https://github.com/owner/repo.git",
//	      "branch":    "main",
//	      "headSHA":   "abc123...",
//	      "remoteSHA": "abc123...",
//	      "updatedAt": "2025-04-21T10:00:00Z"
//	    }
//	  }
//	}
package cloner

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// CloneCacheVersion is the on-disk schema version.
const CloneCacheVersion = 1

// cloneCacheRelPath is the cache file path relative to the clone target dir.
const cloneCacheRelPath = ".gitmap/clone-cache.json"

// CloneCacheEntry is a single per-repo cache record.
type CloneCacheEntry struct {
	URL       string    `json:"url"`
	Branch    string    `json:"branch"`
	HeadSHA   string    `json:"headSHA"`
	RemoteSHA string    `json:"remoteSHA"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CloneCache holds all entries keyed by record.RelativePath.
type CloneCache struct {
	Version int                        `json:"version"`
	Entries map[string]CloneCacheEntry `json:"entries"`

	path string     // resolved on-disk path
	mu   sync.Mutex // guards Entries during concurrent writes
}

// LoadCloneCache reads the cache from <targetDir>/.gitmap/clone-cache.json.
// A missing file or unparseable contents yields a fresh empty cache; this is
// intentional so the cache never blocks a clone run.
func LoadCloneCache(targetDir string) *CloneCache {
	c := &CloneCache{
		Version: CloneCacheVersion,
		Entries: map[string]CloneCacheEntry{},
		path:    filepath.Join(targetDir, cloneCacheRelPath),
	}

	data, err := os.ReadFile(c.path)
	if err != nil {
		return c
	}

	var loaded CloneCache
	if err := json.Unmarshal(data, &loaded); err != nil {
		return c
	}
	if loaded.Version != CloneCacheVersion || loaded.Entries == nil {
		return c
	}

	c.Entries = loaded.Entries

	return c
}

// Save persists the cache to disk. Errors are non-fatal — a failure to write
// the cache must never abort or corrupt a clone run.
func (c *CloneCache) Save() error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(c.path), constants.DirPermission); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0o644)
}

// IsUpToDate reports whether the repo at dest matches the cached fingerprint
// for rec and the remote tip still matches what we cached. When true, the
// caller can safely skip the clone/pull.
//
// All git lookups here are read-only and bounded; on any error we report
// "not up to date" so the caller falls back to the normal clone/pull path.
func (c *CloneCache) IsUpToDate(rec model.ScanRecord, dest string) bool {
	if c == nil {
		return false
	}
	if !IsGitRepo(dest) {
		return false
	}

	c.mu.Lock()
	entry, ok := c.Entries[rec.RelativePath]
	c.mu.Unlock()
	if !ok {
		return false
	}

	if entry.URL != pickURL(rec) || entry.Branch != rec.Branch {
		return false
	}

	localSHA, err := readLocalHead(dest)
	if err != nil || localSHA != entry.HeadSHA {
		return false
	}

	remoteSHA, err := readRemoteHead(dest, rec.Branch)
	if err != nil {
		// Remote unreachable (offline, auth, etc.) — trust the local cache
		// only when the local HEAD still matches the cached HEAD. This keeps
		// `gitmap clone` idempotent in offline reruns without ever returning
		// a false positive when the user has changed the URL/branch.
		return true
	}

	if remoteSHA != entry.RemoteSHA {
		return false
	}

	return entry.RemoteSHA == entry.HeadSHA
}

// Record stores or updates the cache entry for rec at dest. Best-effort:
// failures to introspect the repo are logged via skip and do not abort.
func (c *CloneCache) Record(rec model.ScanRecord, dest string) {
	if c == nil {
		return
	}
	if !IsGitRepo(dest) {
		return
	}

	localSHA, err := readLocalHead(dest)
	if err != nil {
		return
	}
	remoteSHA, _ := readRemoteHead(dest, rec.Branch)
	if remoteSHA == "" {
		// Fall back to the local SHA so we still benefit from offline reruns.
		remoteSHA = localSHA
	}

	c.mu.Lock()
	c.Entries[rec.RelativePath] = CloneCacheEntry{
		URL:       pickURL(rec),
		Branch:    rec.Branch,
		HeadSHA:   localSHA,
		RemoteSHA: remoteSHA,
		UpdatedAt: time.Now().UTC(),
	}
	c.mu.Unlock()
}

// readLocalHead returns the SHA of HEAD in the working tree at dest.
func readLocalHead(dest string) (string, error) {
	cmd := exec.Command(constants.GitBin,
		constants.GitDirFlag, dest,
		constants.GitRevParse, constants.GitHEAD)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// readRemoteHead returns the SHA the remote currently advertises for branch.
// Uses `git ls-remote origin <branch>` so no fetch occurs and no local refs
// are mutated.
func readRemoteHead(dest, branch string) (string, error) {
	if len(branch) == 0 {
		branch = constants.DefaultBranch
	}
	cmd := exec.Command(constants.GitBin,
		constants.GitDirFlag, dest,
		constants.GitLsRemote, constants.GitOrigin, branch)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	line := strings.TrimSpace(string(out))
	if len(line) == 0 {
		return "", nil
	}
	// Output: "<sha>\trefs/heads/<branch>"
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", nil
	}

	return fields[0], nil
}
