// Package installer — import_global.go implements global state archive restoration.
package installer

import (
	"archive/zip"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// ImportGlobalState unpacks a full global state export archive into the database.
func (m *Manager) ImportGlobalState(zipPath string) error {
	if m == nil || m.db == nil {
		return apperror.New("ImportGlobalState", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(zipPath) == "" {
		return apperror.New("ImportGlobalState", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "zipPath cannot be empty"})
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
		m.importGlobalZipEntry(f)
	}

	return nil
}

func (m *Manager) importGlobalZipEntry(f *zip.File) {
	if !strings.HasPrefix(f.Name, "installers/") || !strings.HasSuffix(f.Name, ".json") {
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
