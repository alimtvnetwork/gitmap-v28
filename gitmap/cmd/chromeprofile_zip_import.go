package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// applyChromeExportZIP extracts a ZIP archive into a target profile directory.
func applyChromeExportZIP(zipPath, dstProfile string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip %s: %w", zipPath, err)
	}
	defer r.Close()

	if err := os.MkdirAll(dstProfile, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", dstProfile, err)
	}

	var jsonFile *zip.File
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".json") {
			jsonFile = f
			continue
		}
		
		// Extract SQLite databases
		if isAllowedSQLiteDB(f.Name) {
			if err := extractZipFile(f, filepath.Join(dstProfile, f.Name)); err != nil {
				return err
			}
		}
	}

	if jsonFile != nil {
		rc, err := jsonFile.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		
		var exp chromeExport
		if err := json.NewDecoder(rc).Decode(&exp); err != nil {
			return fmt.Errorf("decode json: %w", err)
		}
		
		if err := applyChromeExport(&exp, dstProfile); err != nil {
			return err
		}
	}
	return nil
}

func extractZipFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, constants.FilePermission)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

func isAllowedSQLiteDB(name string) bool {
	for _, allowed := range constants.ChromeProfileSQLiteEntries {
		if name == allowed {
			return true
		}
	}
	return false
}
