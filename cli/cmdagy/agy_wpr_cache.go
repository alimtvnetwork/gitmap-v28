// Package cmdagy — agy_wpr_cache.go manages local disk caching of remote machine aliases and hostnames.
package cmdagy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveWPRMachineCachePath() string {
	base := store.BinaryDataDir()
	if len(base) > 0 {
		return filepath.Join(base, "wpr_machine_cache.json")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "wpr_machine_cache.json")
	}

	return filepath.Join(home, ".gitmap", "wpr_machine_cache.json")
}

// LoadWPRMachineCache reads cached machine identities from local storage.
func LoadWPRMachineCache() []WPRMachineCacheEntry {
	path := resolveWPRMachineCachePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultLocalMachineCache()
	}

	var list []WPRMachineCacheEntry
	if json.Unmarshal(data, &list) != nil || len(list) == 0 {
		return defaultLocalMachineCache()
	}

	return list
}

func defaultLocalMachineCache() []WPRMachineCacheEntry {
	host, _ := os.Hostname()
	if len(host) == 0 {
		host = "localhost"
	}

	return []WPRMachineCacheEntry{
		{
			NodeId:      "local-01",
			Alias:       "local",
			MachineName: host,
			IPAddress:   "127.0.0.1",
			OS:          "windows",
			Status:      "active",
			CachedAt:    time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// SaveWPRMachineCache persists machine entries to atomic JSON disk storage.
func SaveWPRMachineCache(entries []WPRMachineCacheEntry) error {
	path := resolveWPRMachineCachePath()
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal wpr machine cache")
	}

	return os.WriteFile(path, data, 0644)
}

// UpdateWPRMachineCache updates or appends a machine entry in the local cache.
func UpdateWPRMachineCache(conn db.SSHConnection, machineName, status string) {
	current := LoadWPRMachineCache()
	mName := resolveEffectiveMachineName(conn.Alias, machineName)
	updated := false

	for i := range current {
		if strings.EqualFold(current[i].Alias, conn.Alias) || current[i].IPAddress == conn.IPAddress {
			current[i].MachineName = mName
			current[i].Status = status
			current[i].CachedAt = time.Now().UTC().Format(time.RFC3339)
			updated = true
			break
		}
	}

	if !updated {
		current = append(current, WPRMachineCacheEntry{
			NodeId:      "node-" + conn.Alias,
			Alias:       conn.Alias,
			MachineName: mName,
			IPAddress:   conn.IPAddress,
			OS:          conn.OS,
			Status:      status,
			CachedAt:    time.Now().UTC().Format(time.RFC3339),
		})
	}

	_ = SaveWPRMachineCache(current)
}

func resolveEffectiveMachineName(alias, machineName string) string {
	clean := strings.TrimSpace(machineName)
	if len(clean) > 0 {
		return clean
	}

	return alias
}

// FindCachedMachine returns a cached machine entry matching query by alias, machine name, or IP.
func FindCachedMachine(query string) (WPRMachineCacheEntry, bool) {
	clean := strings.ToLower(strings.TrimSpace(query))
	if len(clean) == 0 {
		return WPRMachineCacheEntry{}, false
	}

	for _, entry := range LoadWPRMachineCache() {
		if strings.ToLower(entry.Alias) == clean || strings.ToLower(entry.MachineName) == clean || entry.IPAddress == clean {
			return entry, true
		}
	}

	return WPRMachineCacheEntry{}, false
}
