// Package cmd — chromeprofile_export.go: JSON snapshot serialization.
// Captures bookmarks + extension IDs + preferences subset. See
// spec/04-generic-cli/40-chrome-profile-copy.md §4 for schema.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// chromeExport is the JSON snapshot format. Keep additive — new
// fields must default-zero so old exports remain importable.
type chromeExport struct {
	SchemaVersion int             `json:"schemaVersion" yaml:"schemaVersion"`
	GitMapVersion string          `json:"gitmapVersion,omitempty" yaml:"gitmapVersion,omitempty"`
	Name          string          `json:"name" yaml:"name"`
	DisplayName   string          `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Email         string          `json:"email,omitempty" yaml:"email,omitempty"`
	ExportedAt    string          `json:"exportedAt" yaml:"exportedAt"`
	Bookmarks     json.RawMessage `json:"bookmarks,omitempty" yaml:"bookmarks,omitempty"`
	Preferences   json.RawMessage `json:"preferences,omitempty" yaml:"preferences,omitempty"`
	ExtensionIDs  []string        `json:"extensionIds,omitempty" yaml:"extensionIds,omitempty"`
}

const chromeExportSchemaVersion = 1

// writeChromeExport reads the curated files from srcProfile and
// writes a JSON snapshot to outPath. Returns bytes written.
func writeChromeExport(srcProfile, name, outPath string) (int, error) {
	displayName, email := resolveProfileNameAndEmail(name, nil)
	exp := chromeExport{
		SchemaVersion: chromeExportSchemaVersion,
		GitMapVersion: constants.Version,
		Name:          name,
		DisplayName:   displayName,
		Email:         email,
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
// Existing files are overwritten. Extensions are merged into pending hints.
func applyChromeExport(exp *chromeExport, dstProfile string) error {
	checkSnapshotVersion(exp.GitMapVersion)
	if err := os.MkdirAll(dstProfile, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", dstProfile, err)
	}
	if err := writeOptional(filepath.Join(dstProfile, "Bookmarks"), exp.Bookmarks); err != nil {
		return err
	}
	if err := writeOptional(filepath.Join(dstProfile, "Preferences"), exp.Preferences); err != nil {
		return err
	}
	if err := writePendingExtensions(dstProfile, exp.ExtensionIDs); err != nil {
		return err
	}
	registerImportedProfileLocalState(exp, dstProfile)
	return nil
}

func writePendingExtensions(dstProfile string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	hint := filepath.Join(dstProfile, "gitmap-pending-extensions.txt")
	merged := mergePendingExtensions(hint, ids)
	return os.WriteFile(hint, []byte(joinLines(merged)), constants.FilePermission)
}

func checkSnapshotVersion(gitmapVersion string) {
	isMismatch := gitmapVersion != "" && gitmapVersion != constants.Version
	if isMismatch {
		fmt.Printf("  \033[1;93m⚠\033[0m Notice: snapshot was generated with gitmap v%s (current running version is v%s)\n", gitmapVersion, constants.Version)
	}
}

func registerImportedProfileLocalState(exp *chromeExport, dstProfile string) {
	if exp.DisplayName == "" {
		return
	}
	dstDir := filepath.Base(dstProfile)
	_ = registerChromeProfileInLocalState(exp.Name, dstDir, exp.DisplayName)
}

func mergePendingExtensions(hintPath string, newIDs []string) []string {
	seen := make(map[string]bool)
	var out []string
	appendExistingExtensions(hintPath, seen, &out)
	for _, id := range newIDs {
		trimmed := strings.TrimSpace(id)
		shouldAppend := trimmed != "" && !seen[trimmed]
		if shouldAppend {
			seen[trimmed] = true
			out = append(out, trimmed)
		}
	}
	return out
}

func appendExistingExtensions(hintPath string, seen map[string]bool, out *[]string) {
	raw, err := os.ReadFile(hintPath)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		shouldAppend := trimmed != "" && !seen[trimmed]
		if shouldAppend {
			seen[trimmed] = true
			*out = append(*out, trimmed)
		}
	}
}

// readOptionalJSON reads path if present, else returns nil.
func readOptionalJSON(path string) json.RawMessage {
	raw, err := os.ReadFile(path)
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
