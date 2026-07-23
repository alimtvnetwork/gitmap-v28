// Package cmd — chromeprofile_export.go: JSON snapshot serialization.
// Captures bookmarks + extension IDs + preferences subset. See
// spec/04-generic-cli/40-chrome-profile-copy.md §4 for schema.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// chromeExport is the JSON snapshot format. Keep additive — new
// fields must default-zero so old exports remain importable.
type chromeExport struct {
	SchemaVersion int             `json:"schemaVersion"`
	Name          string          `json:"name"`
	ExportedAt    string          `json:"exportedAt"`
	Bookmarks     json.RawMessage `json:"bookmarks,omitempty"`
	Preferences   json.RawMessage `json:"preferences,omitempty"`
	ExtensionIDs  []string        `json:"extensionIds,omitempty"`
}

const chromeExportSchemaVersion = 1

// writeChromeExport reads the curated files from srcProfile and
// writes a JSON snapshot to outPath. Returns bytes written.
func writeChromeExport(srcProfile, name, outPath string) (int, error) {
	exp := chromeExport{
		SchemaVersion: chromeExportSchemaVersion,
		Name:          name,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	exp.Bookmarks = readOptionalJSON(filepath.Join(srcProfile, "Bookmarks"))
	exp.Preferences = readOptionalJSON(filepath.Join(srcProfile, "Preferences"))
	exp.ExtensionIDs = listExtensionIDs(filepath.Join(srcProfile, "Extensions"))

	raw, err := json.MarshalIndent(exp, "", constants.JSONIndent)
	if err != nil {
		return 0, fmt.Errorf("marshal export: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil {
		return 0, fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
	}
	if err := os.WriteFile(outPath, raw, constants.FilePermission); err != nil {
		return 0, fmt.Errorf("write %s: %w", outPath, err)
	}
	return len(raw), nil
}

// applyChromeExport writes the export's payloads into dstProfile.
// Existing files are overwritten. Extensions are recorded as a
// pending-install hint file; Chrome itself must reinstall them.
func applyChromeExport(exp *chromeExport, dstProfile string) error {
	if err := os.MkdirAll(dstProfile, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", dstProfile, err)
	}
	if err := writeOptional(filepath.Join(dstProfile, "Bookmarks"), exp.Bookmarks); err != nil {
		return err
	}
	if err := writeOptional(filepath.Join(dstProfile, "Preferences"), exp.Preferences); err != nil {
		return err
	}
	if len(exp.ExtensionIDs) > 0 {
		hint := filepath.Join(dstProfile, "gitmap-pending-extensions.txt")
		if err := os.WriteFile(hint, []byte(joinLines(exp.ExtensionIDs)), constants.FilePermission); err != nil {
			return err
		}
	}
	return nil
}

// readOptionalJSON reads path if present, else returns nil.
func readOptionalJSON(path string) json.RawMessage {
	raw, err := os.ReadFile(path) //nolint:gosec // curated path
	if err != nil {
		return nil
	}
	if !json.Valid(raw) {
		return nil
	}
	return json.RawMessage(raw)
}

// writeOptional writes payload to path when payload is non-empty.
func writeOptional(path string, payload json.RawMessage) error {
	if len(payload) == 0 {
		return nil
	}
	return os.WriteFile(path, payload, constants.FilePermission)
}

// listExtensionIDs returns the subdirectory names under Extensions/,
// which Chrome uses as extension IDs.
func listExtensionIDs(extDir string) []string {
	entries, err := os.ReadDir(extDir)
	if err != nil {
		return nil
	}
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	return ids
}

// joinLines joins items with newlines and a trailing newline.
func joinLines(items []string) string {
	out := ""
	for _, it := range items {
		out += it + "\n"
	}
	return out
}
