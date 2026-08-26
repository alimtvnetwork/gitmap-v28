package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// writeChromeExportZIP writes a zip archive to outPath containing the profile JSON snapshot
// and the curated SQLite databases.
func writeChromeExportZIP(srcProfile, name, outPath string) (int, error) {
	return buildChromeProfileArchive(srcProfile, name, outPath, true, true)
}

// writeChromeExportSQLite writes a zip archive to outPath containing ONLY the curated SQLite databases.
func writeChromeExportSQLite(srcProfile, name, outPath string) (int, error) {
	return buildChromeProfileArchive(srcProfile, name, outPath, false, true)
}

// buildChromeProfileArchive builds the actual zip payload based on requested flags.
func buildChromeProfileArchive(srcProfile, name, outPath string, includeJSON, includeSQLite bool) (int, error) {
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil {
		return 0, fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return 0, fmt.Errorf("create %s: %w", outPath, err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	if includeJSON {
		// Just write to a temporary file, then copy it into the zip
		// Or build it in memory and write to the zip.
		// writeChromeExport currently writes to a file.
		// We can just use it with a temp file.
		tmpJSON := filepath.Join(os.TempDir(), name+"-snapshot.json")
		_, err := writeChromeExport(srcProfile, name, tmpJSON)
		if err != nil {
			return 0, err
		}
		defer os.Remove(tmpJSON)

		w, err := zw.Create(name + ".json")
		if err != nil {
			return 0, err
		}

		jsonFile, err := os.Open(tmpJSON)
		if err != nil {
			return 0, err
		}
		defer jsonFile.Close()

		if _, err := io.Copy(w, jsonFile); err != nil {
			return 0, err
		}
	}

	if includeSQLite {
		for _, dbName := range constants.ChromeProfileSQLiteEntries {
			src := filepath.Join(srcProfile, dbName)
			info, err := os.Stat(src)
			if err != nil {
				continue // skip missing files
			}
			if info.IsDir() {
				continue
			}

			w, err := zw.Create(dbName)
			if err != nil {
				return 0, err
			}

			dbFile, err := os.Open(src)
			if err != nil {
				continue // skip locked/inaccessible files gracefully
			}

			_, _ = io.Copy(w, dbFile)
			dbFile.Close()
		}
	}

	if err := zw.Close(); err != nil {
		return 0, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return int(info.Size()), nil
}
