// Package cmd — chromeprofile_csv.go: CSV serialization of a Chrome
// profile snapshot. Companion to chromeprofile_export.go (JSON). The
// CSV mirrors the same curated subset so spreadsheet tooling can audit
// exports without parsing JSON. Schema: Category,Key,Value.
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// writeChromeExportCSV writes a flat Category/Key/Value CSV alongside
// the JSON snapshot. Returns bytes written.
func writeChromeExportCSV(srcProfile, name, outPath string) (int, error) {
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil {
		return 0, fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
	}
	f, err := os.Create(outPath) //nolint:gosec // curated output path
	if err != nil {
		return 0, fmt.Errorf("create %s: %w", outPath, err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	if err := w.Write([]string{"Category", "Key", "Value"}); err != nil {
		return 0, err
	}
	rows := buildChromeCSVRows(srcProfile, name)
	for _, r := range rows {
		if err := w.Write(r); err != nil {
			return 0, err
		}
	}
	w.Flush()
	info, _ := os.Stat(outPath)
	return int(info.Size()), nil
}

// buildChromeCSVRows assembles the flat row set from the on-disk profile.
func buildChromeCSVRows(srcProfile, name string) [][]string {
	rows := [][]string{
		{"meta", "name", name},
		{"meta", "sourcePath", srcProfile},
	}
	for _, id := range listExtensionIDs(filepath.Join(srcProfile, "Extensions")) {
		rows = append(rows, []string{"extension", "id", id})
	}
	rows = append(rows, flattenPreferences(filepath.Join(srcProfile, "Preferences"))...)
	rows = append(rows, bookmarkSummary(filepath.Join(srcProfile, "Bookmarks"))...)
	return rows
}

// flattenPreferences extracts a small allow-list of non-secret keys.
func flattenPreferences(prefsPath string) [][]string {
	raw, err := os.ReadFile(prefsPath) //nolint:gosec // curated path
	if err != nil {
		return nil
	}
	var doc map[string]any
	if json.Unmarshal(raw, &doc) != nil {
		return nil
	}
	keys := []string{"homepage", "homepage_is_newtabpage", "browser.show_home_button"}
	var out [][]string
	for _, k := range keys {
		if v := lookupDotted(doc, k); v != "" {
			out = append(out, []string{"preference", k, v})
		}
	}
	return out
}

// lookupDotted resolves a dotted key against a parsed JSON map.
func lookupDotted(doc map[string]any, dotted string) string {
	parts := strings.Split(dotted, ".")
	var cur any = doc
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return ""
		}
		cur = m[p]
	}
	if cur == nil {
		return ""
	}
	return fmt.Sprintf("%v", cur)
}

// bookmarkSummary returns one row per top-level bookmark folder.
func bookmarkSummary(bookmarksPath string) [][]string {
	raw, err := os.ReadFile(bookmarksPath) //nolint:gosec // curated path
	if err != nil {
		return nil
	}
	var doc struct {
		Roots map[string]struct {
			Name     string `json:"name"`
			Children []any  `json:"children"`
		} `json:"roots"`
	}
	if json.Unmarshal(raw, &doc) != nil {
		return nil
	}
	var out [][]string
	for key, root := range doc.Roots {
		out = append(out, []string{"bookmark", key, fmt.Sprintf("%s (%d items)", root.Name, len(root.Children))})
	}
	return out
}
