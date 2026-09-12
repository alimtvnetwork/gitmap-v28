package macro

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ImportOptions configures macro import parsing and persistence.
type ImportOptions struct {
	TargetName string
	Format     string
	FilePath   string
	IsForce    bool
	IsDryRun   bool
	IsSingle   bool
	IsAll      bool
	RenameAs   string
	ExceptList []string
}

// ImportResult summarizes the outcome of a macro import batch.
type ImportResult struct {
	TotalFound   int      `json:"total_found" yaml:"total_found"`
	Imported     int      `json:"imported" yaml:"imported"`
	Skipped      int      `json:"skipped" yaml:"skipped"`
	Overwritten  int      `json:"overwritten" yaml:"overwritten"`
	Names        []string `json:"names" yaml:"names"`
	SkippedNames []string `json:"skipped_names,omitempty" yaml:"skipped_names,omitempty"`
}

// ValidateMacro validates macro name against path traversal and verifies step count and commands.
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

	return validateMacroSteps(m.Steps, cleanName)
}

func validateMacroSteps(steps []MacroStep, macroName string) error {
	for i, step := range steps {
		if strings.TrimSpace(step.CommandLine) == "" {
			return apperror.NewValidationError("macro " + macroName + " step command line cannot be empty at index " + strconv.Itoa(i+1))
		}
	}

	return nil
}

func hasInvalidNameChars(name string) bool {
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return true
	}

	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return true
	}

	if hasReservedWin32Chars(name) {
		return true
	}

	return isReservedDeviceName(name)
}

func hasReservedWin32Chars(name string) bool {
	for _, r := range name {
		if r < 32 || strings.ContainsRune(`:*?"<>|`, r) {
			return true
		}
	}

	return false
}

func isReservedDeviceName(name string) bool {
	base := strings.ToUpper(strings.TrimSuffix(name, filepath.Ext(name)))
	switch base {
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return true
	default:
		return false
	}
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
	case ".db", ".sqlite", ".sqlite3", ".sqlitedb":
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

	if format == "sqlite" || format == "db" || format == "sqlitedb" {
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
	filtered, err := filterAndValidateImportTargets(macros, opts)
	if err != nil {
		return nil, err
	}

	res := &ImportResult{TotalFound: len(filtered)}
	for _, m := range filtered {
		if err := processSingleMacroImport(&m, opts, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func filterAndValidateImportTargets(macros []Macro, opts ImportOptions) ([]Macro, error) {
	matched := matchImportMacros(macros, opts)
	if opts.TargetName != "" && len(matched) == 0 {
		return nil, apperror.NewValidationError("macro target " + opts.TargetName + " not found in import archive")
	}

	if opts.IsSingle && len(matched) > 1 && opts.TargetName == "" {
		return nil, apperror.NewValidationError("multiple macros found in archive; please specify target macro name")
	}

	if opts.RenameAs != "" && len(matched) > 1 {
		return nil, apperror.NewValidationError("cannot rename multiple macros with a single name")
	}

	return matched, nil
}

func matchImportMacros(macros []Macro, opts ImportOptions) []Macro {
	var matched []Macro
	for _, m := range macros {
		if isExcludedMacro(m.Name, opts.ExceptList) {
			continue
		}

		if opts.TargetName != "" && !strings.EqualFold(m.Name, opts.TargetName) {
			continue
		}

		matched = append(matched, m)
	}

	return matched
}

func processSingleMacroImport(m *Macro, opts ImportOptions, res *ImportResult) error {
	if opts.RenameAs != "" {
		m.Name = strings.TrimSpace(opts.RenameAs)
	}

	if err := ValidateMacro(m); err != nil {
		return err
	}

	isExisting := MacroExists(m.Name)
	if isExisting && !opts.IsForce {
		res.Skipped++
		res.SkippedNames = append(res.SkippedNames, m.Name)

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
