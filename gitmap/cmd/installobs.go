package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// runOBSSettingsOnly syncs OBS Studio settings from the bundled settings folder.
func runOBSSettingsOnly() {
	fmt.Println("  Syncing OBS Studio settings...")
	syncOBSSettings()
}

// syncOBSSettings copies OBS Studio settings from the bundled folder.
// It looks for a .zip file first (scene collections + profiles), then
// falls back to copying loose files/directories.
func syncOBSSettings() {
	target := obsSettingsTarget()
	if target == "" {
		return
	}

	fmt.Printf("  -> OBS settings target: %s\n", target)

	sourcePath := resolveSettingsPath(
		filepath.Join("settings", "03 - obs"),
		filepath.Join("data", "obs-settings"),
	)

	info, err := os.Stat(sourcePath)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "  Error: OBS settings source not found at %s: %v\n", sourcePath, err)

		return
	}

	// Look for a .zip file to extract.
	zipFile := findFirstZip(sourcePath)
	if zipFile != "" {
		extractOBSSettingsZip(zipFile, target)

		return
	}

	// Fallback: copy loose files/directories.
	copied, copyErr := copyDirRecursive(sourcePath, target)
	if copyErr != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to copy OBS settings: %v\n", copyErr)

		return
	}

	fmt.Printf("  -> Synced %d file(s) to %s\n", copied, target)
}

// obsSettingsTarget returns the OBS Studio config directory.
func obsSettingsTarget() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			fmt.Fprint(os.Stderr, "  Error: APPDATA environment variable not set\n")

			return ""
		}

		return filepath.Join(appData, "obs-studio")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: could not resolve home directory: %v\n", err)

			return ""
		}

		return filepath.Join(home, "Library", "Application Support", "obs-studio")
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: could not resolve home directory: %v\n", err)

			return ""
		}

		return filepath.Join(home, ".config", "obs-studio")
	}
}

// findFirstZip returns the path of the first .zip file in a directory.
func findFirstZip(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(entry.Name()), ".zip") {
			return filepath.Join(dir, entry.Name())
		}
	}

	return ""
}

// extractOBSSettingsZip extracts an OBS settings zip and routes files to
// the correct OBS subdirectories:
//   - .json files -> basic/scenes/
//   - directories -> basic/profiles/
func extractOBSSettingsZip(zipPath, target string) {
	fmt.Printf("  -> Extracting OBS settings from: %s\n", filepath.Base(zipPath))

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to open OBS settings zip: %v\n", err)

		return
	}
	defer reader.Close()

	// Create temp directory for extraction.
	tmpDir, err := os.MkdirTemp("", "gitmap-obs-extract-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to create temp directory: %v\n", err)

		return
	}
	defer os.RemoveAll(tmpDir)

	// Extract all files to temp.
	for _, file := range reader.File {
		extractOBSZipEntry(tmpDir, file)
	}

	// Route extracted files to the correct OBS directories.
	scenesDir := filepath.Join(target, "basic", "scenes")
	profilesDir := filepath.Join(target, "basic", "profiles")

	if err := os.MkdirAll(scenesDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to create scenes directory: %v\n", err)

		return
	}

	if err := os.MkdirAll(profilesDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to create profiles directory: %v\n", err)

		return
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to read extracted files: %v\n", err)

		return
	}

	scenes := 0
	profiles := 0

	for _, entry := range entries {
		srcPath := filepath.Join(tmpDir, entry.Name())

		if entry.IsDir() {
			// Directories are profile folders.
			dstPath := filepath.Join(profilesDir, entry.Name())
			n, dirErr := copyDirRecursive(srcPath, dstPath)

			if dirErr != nil {
				fmt.Fprintf(os.Stderr, "  ! Failed to copy profile %s: %v\n", entry.Name(), dirErr)

				continue
			}

			profiles++
			_ = n
		} else if strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			// JSON files are scene collections.
			dstPath := filepath.Join(scenesDir, entry.Name())
			copyErr := copyFile(srcPath, dstPath)

			if copyErr != nil {
				fmt.Fprintf(os.Stderr, "  ! Failed to copy scene %s: %v\n", entry.Name(), copyErr)

				continue
			}

			scenes++
		}
	}

	fmt.Printf("  -> Synced %d scene(s) and %d profile(s) to %s\n", scenes, profiles, target)
}

// extractOBSZipEntry extracts a single entry from the OBS zip.
func extractOBSZipEntry(target string, file *zip.File) {
	cleanName := filepath.FromSlash(file.Name)
	destPath := filepath.Join(target, cleanName)

	// Path traversal protection.
	absTarget, absErr := filepath.Abs(target)
	if absErr != nil {
		absTarget = target
	}
	absDest, destErr := filepath.Abs(destPath)
	if destErr != nil {
		absDest = destPath
	}

	if !strings.HasPrefix(absDest, absTarget+string(os.PathSeparator)) {
		return
	}

	if file.FileInfo().IsDir() {
		if mkErr := os.MkdirAll(destPath, 0o755); mkErr != nil {
			fmt.Fprintf(os.Stderr, "  ! Failed to create directory %s: %v\n", destPath, mkErr)

			return
		}

		return
	}

	if mkErr := os.MkdirAll(filepath.Dir(destPath), 0o755); mkErr != nil {
		fmt.Fprintf(os.Stderr, "  ! Failed to create parent directory for %s: %v\n", destPath, mkErr)

		return
	}

	src, err := file.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ! Failed to open zip entry %s: %v\n", file.Name, err)

		return
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ! Failed to create file %s: %v\n", destPath, err)

		return
	}
	defer dst.Close()

	if _, copyErr := io.Copy(dst, io.LimitReader(src, 50*1024*1024)); copyErr != nil {
		fmt.Fprintf(os.Stderr, "  ! Failed to extract %s: %v\n", file.Name, copyErr)

		return
	}
}

// copyDirRecursive copies all files from src to dst recursively.
func copyDirRecursive(src, dst string) (int, error) {
	copied := 0

	entries, err := os.ReadDir(src)
	if err != nil {
		return 0, err
	}

	err = os.MkdirAll(dst, 0o755)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// Skip readme files.
		name := strings.ToLower(entry.Name())
		if name == "readme.txt" || name == "readme.md" {
			continue
		}

		if entry.IsDir() {
			n, dirErr := copyDirRecursive(srcPath, dstPath)
			if dirErr != nil {
				fmt.Fprintf(os.Stderr, "  ! Failed to copy directory %s: %v\n", entry.Name(), dirErr)

				continue
			}

			copied += n

			continue
		}

		copyErr := copyFile(srcPath, dstPath)
		if copyErr != nil {
			fmt.Fprintf(os.Stderr, "  ! Failed to copy %s: %v\n", entry.Name(), copyErr)

			continue
		}

		copied++
	}

	return copied, nil
}

// NOTE: copyFile is defined in update.go and shared across the cmd package.
