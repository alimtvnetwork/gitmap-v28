package macro

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ExportOptions specifies parameters for macro export operations.
type ExportOptions struct {
	Format     string
	FilePath   string
	TargetName string
	IsAll      bool
	ExceptList []string
}

// ExportToJSON serializes data to indented JSON bytes.
func ExportToJSON(data any) ([]byte, error) {
	bytes, err := json.MarshalIndent(data, "", constants.JSONIndent)
	if err != nil {
		return nil, apperror.WrapSimple(err, "marshal macro json")
	}

	return bytes, nil
}

// ExportToYAML serializes data to YAML bytes.
func ExportToYAML(data any) ([]byte, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, apperror.WrapSimple(err, "marshal macro yaml")
	}

	return bytes, nil
}

// ExportToZIP packages macros into a ZIP archive containing individual JSON files.
func ExportToZIP(macros []Macro, filePath string) error {
	_ = os.MkdirAll(filepath.Dir(filePath), constants.DirPermission)
	file, err := os.Create(filePath)
	if err != nil {
		return apperror.WrapSimple(err, "create zip archive")
	}

	defer file.Close()
	zw := zip.NewWriter(file)
	defer zw.Close()

	return writeAllZipEntries(zw, macros)
}

func writeAllZipEntries(zw *zip.Writer, macros []Macro) error {
	for _, m := range macros {
		if err := writeZipEntry(zw, m); err != nil {
			return err
		}
	}

	return nil
}

func writeZipEntry(zw *zip.Writer, m Macro) error {
	raw, err := json.MarshalIndent(m, "", constants.JSONIndent)
	if err != nil {
		return apperror.WrapSimple(err, "marshal zip entry json")
	}

	entryWriter, err := zw.Create(m.Name + ".json")
	if err != nil {
		return apperror.WrapSimple(err, "create zip entry")
	}

	_, writeErr := entryWriter.Write(raw)

	return writeErr
}

// FilterMacrosForExport filters macros by name and exclusion list.
func FilterMacrosForExport(macros []Macro, opts ExportOptions) []Macro {
	var result []Macro
	for _, m := range macros {
		if isExcludedMacro(m.Name, opts.ExceptList) {
			continue
		}

		if opts.IsAll || strings.EqualFold(m.Name, opts.TargetName) {
			result = append(result, m)
		}
	}

	return result
}

func isExcludedMacro(name string, exceptList []string) bool {
	for _, exp := range exceptList {
		if strings.EqualFold(strings.TrimSpace(exp), strings.TrimSpace(name)) {
			return true
		}
	}

	return false
}

// SerializeMacros serializes macros into the specified format (json or yaml).
func SerializeMacros(macros []Macro, format string) ([]byte, error) {
	cleanFormat := strings.ToLower(strings.TrimSpace(format))
	if cleanFormat == constants.OutputYAML || cleanFormat == "yml" {
		return ExportToYAML(macros)
	}

	return ExportToJSON(macros)
}

// SerializeSingleOrAll serializes macros as a single object when isSingle is true, or array otherwise.
func SerializeSingleOrAll(macros []Macro, isSingle bool, format string) ([]byte, error) {
	if isSingle && len(macros) == 1 {
		return serializeSingleMacro(macros[0], format)
	}

	return SerializeMacros(macros, format)
}

func serializeSingleMacro(m Macro, format string) ([]byte, error) {
	cleanFormat := strings.ToLower(strings.TrimSpace(format))
	if cleanFormat == constants.OutputYAML || cleanFormat == "yml" {
		return ExportToYAML(m)
	}

	return ExportToJSON(m)
}

// WriteExportPayload writes raw bytes to destination path ensuring parent directory exists.
func WriteExportPayload(filePath string, payload []byte) error {
	if err := ensureExportDirExists(filePath); err != nil {
		return err
	}

	tmpPath := fmt.Sprintf("%s.%d.tmp", filePath, time.Now().UnixNano())
	if err := writeAndSyncPayload(tmpPath, payload); err != nil {
		_ = os.Remove(tmpPath)

		return err
	}

	return replaceExportFile(tmpPath, filePath)
}

func ensureExportDirExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "" || dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, constants.DirPermission); err != nil {
		return apperror.WrapSimple(err, "create export directory")
	}

	return nil
}

func writeAndSyncPayload(tmpPath string, payload []byte) error {
	file, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, constants.FilePermission)
	if err != nil {
		return apperror.WrapSimple(err, "create temp export file")
	}

	defer file.Close()
	if _, writeErr := file.Write(payload); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write temp export payload")
	}

	return file.Sync()
}

func replaceExportFile(tmpPath, filePath string) error {
	_ = os.Remove(filePath)
	if err := os.Rename(tmpPath, filePath); err != nil {
		_ = os.Remove(tmpPath)

		return apperror.WrapSimple(err, "rename temp export file")
	}

	return nil
}
