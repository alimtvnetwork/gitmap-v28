// Package installer — import.go implements zip archive import business logic.
package installer

import (
	"archive/zip"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// ImportFromZip unzips and registers all json script files in the archive.
func (m *Manager) ImportFromZip(zipPath string) error {
	if m == nil || m.db == nil {
		return apperror.New("ImportFromZip", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(zipPath) == "" {
		return apperror.New("ImportFromZip", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "zipPath cannot be empty"})
	}

	if _, err := os.Stat(zipPath); err != nil {
		return err
	}

	zr, errZip := zip.OpenReader(zipPath)
	if errZip != nil {
		return errZip
	}
	defer zr.Close()

	for _, f := range zr.File {
		m.importZipFileEntry(f)
	}

	return nil
}

func (m *Manager) importZipFileEntry(f *zip.File) {
	if !strings.HasSuffix(strings.ToLower(f.Name), ".json") {
		return
	}
	rc, errOpen := f.Open()
	if errOpen != nil {
		return
	}
	defer rc.Close()

	data, errRead := io.ReadAll(rc)
	if errRead == nil {
		m.ImportFromJson(string(data))
	}
}
