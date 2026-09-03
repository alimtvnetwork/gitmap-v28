// Package installer — export_global.go implements full system export packaging.
package installer

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// ExportGlobalState packages the full Gitmap state (installers, configurations) into a zip bundle.
func (m *Manager) ExportGlobalState(targetPath string) error {
	if m == nil || m.db == nil {
		return apperror.New("ExportGlobalState", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	outPath := targetPath
	if outPath == "" {
		outPath = "gitmap-export.zip"
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

	scripts, errList := m.db.ListInstallers()
	if errList != nil {
		return errList
	}

	for _, s := range scripts {
		data, errMarshal := json.MarshalIndent(s, "", "  ")
		if errMarshal != nil {
			return errMarshal
		}

		w, errZip := zw.Create(fmt.Sprintf("installers/%s.json", s.Slug))
		if errZip != nil {
			return errZip
		}

		if _, errWrite := w.Write(data); errWrite != nil {
			return errWrite
		}
	}

	return nil
}
