// Package installer — export.go implements single installer zip export logic.
package installer

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// ExportToZip packages a specific installer script into a .zip archive.
func (m *Manager) ExportToZip(slug, targetPath string) error {
	if m == nil || m.db == nil {
		return apperror.New("ExportToZip", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(slug) == "" {
		return apperror.New("ExportToZip", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	script, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	outPath := targetPath
	if outPath == "" {
		outPath = fmt.Sprintf("%s-export.zip", slug)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil && filepath.Dir(outPath) != "." {
		return err
	}

	f, errCreate := os.Create(outPath)
	if errCreate != nil {
		return errCreate
	}

	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	data, errMarshal := json.MarshalIndent(script, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	w, errZip := zw.Create(fmt.Sprintf("%s.json", script.Slug))
	if errZip != nil {
		return errZip
	}

	_, errWrite := w.Write(data)

	return errWrite
}
