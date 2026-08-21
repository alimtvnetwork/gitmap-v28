// Package cmd — chromeprofile_import_csv.go: parses a CSV snapshot
// (Category,Key,Value rows produced by writeChromeExportCSV) back into
// a chromeExport struct. Lossy by design: bookmarks are not preserved
// in the CSV form, so they round-trip as empty.
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// loadChromeImport dispatches to JSON or CSV based on the file
// extension. Anything other than .csv is treated as JSON.
func loadChromeImport(path string) (*chromeExport, error) {
	if strings.HasSuffix(strings.ToLower(path), constants.ExtCSV) {
		fmt.Fprint(os.Stderr, constants.MsgChromeProfileImportCSV)
		return readChromeExportCSV(path)
	}
	return readChromeExport(path)
}

func parseCSVRows(path string) ([][]string, error) {
	f, err := os.Open(path) //nolint:gosec // user-supplied path
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return rows, nil
}

// readChromeExportCSV parses a Category/Key/Value CSV produced by
// writeChromeExportCSV. Unknown categories are ignored so the format
// stays additive.
func readChromeExportCSV(path string) (*chromeExport, error) {
	rows, err := parseCSVRows(path)
	if err != nil {
		return nil, err
	}
	exp := &chromeExport{SchemaVersion: chromeExportSchemaVersion}
	prefs := map[string]any{}
	populateExportRows(exp, prefs, rows)
	finalizeExport(exp, prefs)
	return exp, nil
}

func populateExportRows(exp *chromeExport, prefs map[string]any, rows [][]string) {
	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}
		assignChromeCSVRow(exp, prefs, row[0], row[1], row[2])
	}
}

func finalizeExport(exp *chromeExport, prefs map[string]any) {
	if len(prefs) > 0 {
		raw, _ := json.Marshal(prefs)
		exp.Preferences = raw
	}
	if exp.Name == "" {
		exp.Name = constants.ChromeDefaultImportedName
	}
}

// assignChromeCSVRow routes a single Category/Key/Value triple into
// the right field of the rebuilt export.
func assignChromeCSVRow(exp *chromeExport, prefs map[string]any, category, key, value string) {
	switch category {
	case constants.ChromeCSVCategoryMeta:
		if key == constants.ChromeCSVKeyName {
			exp.Name = value
		}
	case constants.ChromeCSVCategoryExtension:
		if key == constants.ChromeCSVKeyID && value != "" {
			exp.ExtensionIDs = append(exp.ExtensionIDs, value)
		}
	case constants.ChromeCSVCategoryPreference:
		prefs[key] = value
	}
}
