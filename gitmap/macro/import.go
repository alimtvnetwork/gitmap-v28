package macro

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ImportOptions configures macro import parsing and persistence.
type ImportOptions struct {
	Format     string
	FilePath   string
	IsForce    bool
	IsDryRun   bool
	ExceptList []string
}

// ImportResult summarizes the outcome of a macro import batch.
type ImportResult struct {
	TotalFound  int      `json:"total_found" yaml:"total_found"`
	Imported    int      `json:"imported" yaml:"imported"`
	Skipped     int      `json:"skipped" yaml:"skipped"`
	Overwritten int      `json:"overwritten" yaml:"overwritten"`
	Names       []string `json:"names" yaml:"names"`
}

// ValidateMacro validates macro name against path traversal and verifies step count.
func ValidateMacro(m *Macro) error {
	cleanName := strings.TrimSpace(m.Name)
	if cleanName == "" {
		return apperror.NewValidationError("macro name cannot be empty")
	}

	if hasInvalidNameChars(cleanName) {
		return apperror.NewValidationError("invalid characters in macro name: " + cleanName)
	}

	if len(m.Steps) == 0 {
		return apperror.NewValidationError("macro contains no steps: " + cleanName)
	}

	return nil
}

func hasInvalidNameChars(name string) bool {
	return strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\")
}

// ParseImportJSON deserializes JSON bytes using polymorphic dual-shape detection.
func ParseImportJSON(data []byte) ([]Macro, error) {
	var list []Macro
	if err := json.Unmarshal(data, &list); err == nil && len(list) > 0 {
		return list, nil
	}

	var single Macro
	if err := json.Unmarshal(data, &single); err == nil && single.Name != "" {
		return []Macro{single}, nil
	}

	return nil, apperror.NewSimple("invalid json macro format", "E6022")
}

// ParseImportYAML deserializes YAML bytes using polymorphic dual-shape detection.
func ParseImportYAML(data []byte) ([]Macro, error) {
	var list []Macro
	if err := yaml.Unmarshal(data, &list); err == nil && len(list) > 0 {
		return list, nil
	}

	var single Macro
	if err := yaml.Unmarshal(data, &single); err == nil && single.Name != "" {
		return []Macro{single}, nil
	}

	return nil, apperror.NewSimple("invalid yaml macro format", "E6023")
}

// ParseImportZIP unpacks and deserializes all JSON macro definitions inside a ZIP archive.
func ParseImportZIP(filePath string) ([]Macro, error) {
	zr, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open zip macro archive")
	}

	defer zr.Close()

	return extractMacrosFromZip(&zr.Reader)
}

func extractMacrosFromZip(zr *zip.Reader) ([]Macro, error) {
	var list []Macro
	for _, f := range zr.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".json") {
			continue
		}

		m, err := readZipMacroEntry(f)
		if err == nil && m != nil && m.Name != "" {
			list = append(list, *m)
		}
	}

	if len(list) == 0 {
		return nil, apperror.NewSimple("no valid macro definitions found in zip", "E6024")
	}

	return list, nil
}

func readZipMacroEntry(f *zip.File) (*Macro, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open zip file entry")
	}

	defer rc.Close()
	raw, err := io.ReadAll(rc)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read zip file entry")
	}

	var m Macro
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, apperror.WrapSimple(err, "unmarshal zip entry")
	}

	return &m, nil
}

// InferMacroFormat detects the format based on file extension.
func InferMacroFormat(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		return constants.OutputYAML
	case ".db", ".sqlite", ".sqlite3":
		return "sqlite"
	case ".zip":
		return "zip"
	default:
		return constants.OutputJSON
	}
}

// ParseImportFile loads and parses macros from disk with auto format inference.
func ParseImportFile(filePath string, explicitFormat string) ([]Macro, error) {
	format := strings.ToLower(strings.TrimSpace(explicitFormat))
	if format == "" {
		format = InferMacroFormat(filePath)
	}

	if format == "sqlite" || format == "db" {
		return ParseImportSQLite(filePath)
	}

	if format == "zip" {
		return ParseImportZIP(filePath)
	}

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read import file")
	}

	if format == constants.OutputYAML || format == "yml" {
		return ParseImportYAML(raw)
	}

	return ParseImportJSON(raw)
}

// ImportMacros validates, filters, and saves parsed macros into the local store.
func ImportMacros(macros []Macro, opts ImportOptions) (*ImportResult, error) {
	res := &ImportResult{TotalFound: len(macros)}
	for _, m := range macros {
		if isExcludedMacro(m.Name, opts.ExceptList) {
			continue
		}

		if err := processSingleMacroImport(&m, opts, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func processSingleMacroImport(m *Macro, opts ImportOptions, res *ImportResult) error {
	if err := ValidateMacro(m); err != nil {
		return err
	}

	isExisting := MacroExists(m.Name)
	if isExisting && !opts.IsForce {
		res.Skipped++

		return nil
	}

	recordImportMetrics(isExisting, res, m.Name)
	if opts.IsDryRun {
		return nil
	}

	return SaveMacro(m)
}

func recordImportMetrics(isExisting bool, res *ImportResult, name string) {
	if isExisting {
		res.Overwritten++
	} else {
		res.Imported++
	}

	res.Names = append(res.Names, name)
}
