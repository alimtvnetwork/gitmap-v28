package cmdautofix

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

// ---------------------------------------------------------------------------
// Scan cache
// ---------------------------------------------------------------------------

// cacheEntry is one cached file's scan record.
type cacheEntry struct {
	ModUnix    int64       `json:"modtime_unix"`
	Size       int64       `json:"size"`
	Categories []string    `json:"categories"`
	Violations []Violation `json:"violations"`
}

func cachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil || dir == "" {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "gitmap", "fix-scan-cache.json")
}

func loadCache() map[string]cacheEntry {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return map[string]cacheEntry{}
	}
	m := map[string]cacheEntry{}
	if json.Unmarshal(data, &m) != nil {
		return map[string]cacheEntry{}
	}
	return m
}

func saveCache(m map[string]cacheEntry) {
	p := cachePath()
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	data, err := json.Marshal(m)
	if err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0o644)
}

// cacheHit reports whether entry satisfies the hit rule: modtime+size match
// AND every requested category was covered by the cached scan.
func cacheHit(entry cacheEntry, info fs.FileInfo, requested []string) bool {
	if entry.ModUnix != info.ModTime().Unix() || entry.Size != info.Size() {
		return false
	}
	for _, want := range requested {
		if !extInList(want, entry.Categories) {
			return false
		}
	}
	return true
}

// filterViolations keeps only violations from the requested categories.
func filterViolations(in []Violation, requested []string) []Violation {
	out := make([]Violation, 0, len(in))
	for _, v := range in {
		if extInList(v.Category, requested) {
			out = append(out, v)
		}
	}
	return out
}

