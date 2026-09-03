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

// applyChromeExportZIP extracts a ZIP archive into target profile directories.
func applyChromeExportZIP(zipPath, dstProfile string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip %s: %w", zipPath, err)
	}
	defer r.Close()

	targetDir := filepath.Join(chromeUserDataDir(), dstProfile)
	_, _ = snapshotChromeProfile(targetDir, "pre-import")

	m := readZipManifest(r)
	if isMultiProfileArchive(m, r) {
		return extractMultiProfileZip(r)
	}
	return extractSingleProfileZip(r, dstProfile)
}

func isMultiProfileArchive(m *chromeProfileManifest, r *zip.ReadCloser) bool {
	if m != nil {
		checkSnapshotVersion(m.GitMapVersion)
		return len(m.Profiles) > 1 || isMultiProfileZip(r)
	}
	return isMultiProfileZip(r)
}

func readZipManifest(r *zip.ReadCloser) *chromeProfileManifest {
	for _, f := range r.File {
		if f.Name == "manifest.json" {
			return decodeZipManifestFile(f)
		}
	}
	return nil
}

func decodeZipManifestFile(f *zip.File) *chromeProfileManifest {
	rc, err := f.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()
	var m chromeProfileManifest
	if json.NewDecoder(rc).Decode(&m) == nil {
		return &m
	}
	return nil
}

func isMultiProfileZip(r *zip.ReadCloser) bool {
	for _, f := range r.File {
		if strings.Contains(f.Name, "/") && !strings.HasPrefix(f.Name, "__") {
			return true
		}
	}
	return false
}

func extractMultiProfileZip(r *zip.ReadCloser) error {
	importedProfiles := make(map[string]bool)
	for _, f := range r.File {
		processMultiProfileZipEntry(f, importedProfiles)
	}
	for profName := range importedProfiles {
		label := formatChromeProfileLabel(profName, nil)
		fmt.Printf("  \033[1;92m✓\033[0m %s → imported from zip\n", label)
	}
	return nil
}

func processMultiProfileZipEntry(f *zip.File, imported map[string]bool) {
	parts := strings.Split(filepath.ToSlash(f.Name), "/")
	if len(parts) < 2 {
		return
	}
	profName, fileName := parts[0], parts[len(parts)-1]
	profDir := chromeProfilePath(profName)
	_ = os.MkdirAll(profDir, constants.DirPermission)
	if strings.HasSuffix(fileName, ".json") {
		_ = importChromeJSONFile(f, profDir)
	} else if isAllowedSQLiteDB(fileName) || fileName == "Bookmarks" || fileName == "Preferences" {
		_ = extractZipFile(f, filepath.Join(profDir, fileName))
	}
	imported[profName] = true
}

func extractSingleProfileZip(r *zip.ReadCloser, dstProfile string) error {
	if err := os.MkdirAll(dstProfile, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", dstProfile, err)
	}
	var jsonFile *zip.File
	for _, f := range r.File {
		if f.Name == "manifest.json" {
			continue
		}
		if strings.HasSuffix(f.Name, ".json") {
			jsonFile = f
			continue
		}
		if isAllowedSQLiteDB(f.Name) || f.Name == "Bookmarks" || f.Name == "Preferences" {
			_ = extractZipFile(f, filepath.Join(dstProfile, f.Name))
		}
	}
	if jsonFile != nil {
		return importChromeJSONFile(jsonFile, dstProfile)
	}
	return nil
}

func importChromeJSONFile(jsonFile *zip.File, dstProfile string) error {
	rc, err := jsonFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	var exp chromeExport
	if err := json.NewDecoder(rc).Decode(&exp); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	return applyChromeExport(&exp, dstProfile)
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
