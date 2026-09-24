package cmdvscode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

type jsonRepairStats struct {
	TotalCount   int
	MissingPaths int
	BackupPath   string
	WasCorrupt   bool
}

func resolveProjectJSONPaths() (string, string, error) {
	root, err := vscodepm.UserDataRoot()
	if err != nil {
		return "", "", err
	}

	modern := filepath.Join(root, "User", "globalStorage", "alefragnani.project-manager", "projects.json")
	legacy := filepath.Join(root, "User", "projects.json")

	return modern, legacy, nil
}

func repairProjectManagerJSON() (jsonRepairStats, error) {
	modernPath, legacyPath, err := resolveProjectJSONPaths()
	if err != nil {
		return jsonRepairStats{}, err
	}

	activePath := pickActiveJSONPath(modernPath, legacyPath)
	entries, stats := loadAndSanitizeEntries(activePath)

	_ = vscodepm.WriteEntries(modernPath, entries)
	_ = vscodepm.WriteEntries(legacyPath, entries)

	return stats, nil
}

func pickActiveJSONPath(modern, legacy string) string {
	if _, err := os.Stat(modern); err == nil {
		return modern
	}
	if _, err := os.Stat(legacy); err == nil {
		return legacy
	}

	return modern
}

func loadAndSanitizeEntries(path string) ([]vscodepm.Entry, jsonRepairStats) {
	var stats jsonRepairStats
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return []vscodepm.Entry{}, stats
	}

	var entries []vscodepm.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		stats.WasCorrupt = true
		stats.BackupPath = backupCorruptJSON(path, data)

		return []vscodepm.Entry{}, stats
	}

	stats.TotalCount = len(entries)
	stats.MissingPaths = sanitizeEntries(entries)

	return entries, stats
}

func sanitizeEntries(entries []vscodepm.Entry) int {
	missing := 0
	for i := range entries {
		if entries[i].Paths == nil {
			entries[i].Paths = []string{}
		}
		if entries[i].Tags == nil {
			entries[i].Tags = []string{"gitmap"}
		}
		if isMissingRootPath(entries[i].RootPath) {
			missing++
		}
	}

	return missing
}

func isMissingRootPath(rootPath string) bool {
	if rootPath == "" {
		return false
	}
	_, err := os.Stat(rootPath)

	return os.IsNotExist(err)
}

func backupCorruptJSON(path string, data []byte) string {
	ts := time.Now().Format("20060102150405")
	bak := fmt.Sprintf("%s.corrupt.%s.bak", path, ts)
	_ = os.WriteFile(bak, data, 0644)

	return bak
}
