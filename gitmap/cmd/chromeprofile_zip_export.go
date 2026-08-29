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
func buildChromeProfileArchive(
	srcProfile,
	name,
	outPath string,
	includeJSON,
	includeSQLite bool,
) (int, error) {
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

	if err := maybeExportJSON(zw, srcProfile, name, includeJSON); err != nil {
		return 0, err
	}

	if err := maybeExportSQLite(zw, srcProfile, includeSQLite); err != nil {
		return 0, err
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

func maybeExportJSON(zw *zip.Writer, srcProfile, name string, shouldInclude bool) error {
	if !shouldInclude {
		return nil
	}
	return exportChromeProfileJSON(zw, srcProfile, name)
}

func maybeExportSQLite(zw *zip.Writer, srcProfile string, shouldInclude bool) error {
	if !shouldInclude {
		return nil
	}
	return exportChromeProfileSQLite(zw, srcProfile)
}

func exportChromeProfileJSON(zw *zip.Writer, srcProfile, name string) error {
	tmpJSON := filepath.Join(os.TempDir(), name+"-snapshot.json")
	if _, err := writeChromeExport(srcProfile, name, tmpJSON); err != nil {
		return err
	}
	defer os.Remove(tmpJSON)

	w, err := zw.Create(name + ".json")
	if err != nil {
		return err
	}

	jsonFile, err := os.Open(tmpJSON)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	_, err = io.Copy(w, jsonFile)
	return err
}

func exportChromeProfileSQLite(zw *zip.Writer, srcProfile string) error {
	for _, dbName := range constants.ChromeProfileSQLiteEntries {
		if err := copySQLiteEntryToZip(zw, srcProfile, dbName); err != nil {
			return err
		}
	}
	return nil
}

func copySQLiteEntryToZip(zw *zip.Writer, srcProfile, dbName string) error {
	src := filepath.Join(srcProfile, dbName)
	info, err := os.Stat(src)
	if err != nil || info.IsDir() {
		return nil
	}

	w, err := zw.Create(dbName)
	if err != nil {
		return err
	}

	dbFile, err := os.Open(src)
	if err != nil {
		return nil
	}
	defer dbFile.Close()

	_, err = io.Copy(w, dbFile)
	return err
}
