package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// extractNppSettingsZip extracts the bundled settings zip to the target.
func extractNppSettingsZip(target string) {
	// Try the current filename first, then legacy name.
	zipPath := resolveNppDataPath("02. Notepad++ settings.zip")

	fmt.Printf(constants.MsgInstallNppExtract, target)
	fmt.Printf("  -> Settings zip: %s\n", zipPath)

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppZipNotFound, zipPath, err)
		syncNppSettingsFallback(target)

		return
	}
	defer reader.Close()

	extracted := 0

	for _, file := range reader.File {
		extractZipEntry(target, file)
		extracted++
	}

	fmt.Printf("  ✓ Extracted %d files\n", extracted)
	fmt.Printf(constants.MsgNppSettingsSynced, target)
}

// extractZipEntry writes a single zip entry to the target directory.
func extractZipEntry(target string, file *zip.File) {
	cleanName := filepath.FromSlash(file.Name)
	destPath := filepath.Join(target, cleanName)

	absTarget, absErr := filepath.Abs(target)
	if absErr != nil {
		absTarget = target
	}
	absDest, destErr := filepath.Abs(destPath)
	if destErr != nil {
		absDest = destPath
	}

	if !strings.HasPrefix(absDest, absTarget+string(os.PathSeparator)) {
		fmt.Fprintf(os.Stderr, constants.ErrNppExtractEntry, file.Name, destPath, fmt.Errorf("path traversal detected"))

		return
	}

	if file.FileInfo().IsDir() {
		err := os.MkdirAll(destPath, 0o755)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrNppDirCreate, destPath, err)
		}

		return
	}

	writeZipFile(target, file, destPath)
}

// writeZipFile extracts a single file from the zip archive.
func writeZipFile(target string, file *zip.File, destPath string) {
	err := os.MkdirAll(filepath.Dir(destPath), 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppDirCreate, filepath.Dir(destPath), err)

		return
	}

	src, err := file.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppExtractEntry, file.Name, destPath, err)

		return
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileCreate, destPath, err)

		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, io.LimitReader(src, maxNppFileSize))
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileCopy, file.Name, destPath, err)
	}
}

// syncNppSettingsFallback copies loose settings files as a fallback.
func syncNppSettingsFallback(target string) {
	source := resolveNppDataPath("npp-settings")

	entries, err := os.ReadDir(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppSourceDir, source, err)

		return
	}

	for _, entry := range entries {
		if entry.Name() == "npp-settings.zip" {
			continue
		}
		copySettingsFile(source, target, entry.Name())
	}

	fmt.Printf(constants.MsgNppSettingsFallback, target)
}

// copySettingsFile copies a single settings file to the target.
func copySettingsFile(source, target, name string) {
	src := filepath.Join(source, name)
	dst := filepath.Join(target, name)

	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileRead, src, err)

		return
	}

	err = os.WriteFile(dst, data, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileWrite, dst, err)
	}
}
