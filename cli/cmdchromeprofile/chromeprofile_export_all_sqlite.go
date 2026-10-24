package cmdchromeprofile

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	_ "modernc.org/sqlite"
)

func insertProfileMetaToSQLite(tx *dbengine.TxWrapper, name string, exp chromeExport) *apperror.AppError {
	ctx := context.Background()
	displayName := chromeProfileDisplayName(name)
	query := "INSERT OR REPLACE INTO chrome_profiles (name, display_name, exported_at, extension_count) VALUES (?, ?, ?, ?)"
	if execRes := tx.Exec(ctx, query, name, displayName, exp.ExportedAt, len(exp.ExtensionIDs)); execRes.IsFailure() {
		return execRes.AppError()
	}

	if prefErr := insertOptionalJSONToSQLite(tx, "chrome_preferences", "preferences_json", name, exp.Preferences); prefErr != nil {
		return prefErr
	}

	return insertOptionalJSONToSQLite(tx, "chrome_bookmarks", "bookmarks_json", name, exp.Bookmarks)
}

func insertTokensToSQLite(tx *dbengine.TxWrapper, name string, vault *ChromeTokenVault) *apperror.AppError {
	if vault == nil || len(vault.Tokens) == 0 {
		return nil
	}

	ctx := context.Background()
	query := "INSERT OR REPLACE INTO chrome_tokens (profile_name, service, account_id, raw_base64, double_base64) VALUES (?, ?, ?, ?, ?)"
	for _, t := range vault.Tokens {
		if execRes := tx.Exec(ctx, query, name, t.Service, t.AccountID, t.RawBase64, t.DoubleBase64); execRes.IsFailure() {
			return execRes.AppError()
		}
	}

	return nil
}

func insertOptionalJSONToSQLite(tx *dbengine.TxWrapper, table, col, name string, raw json.RawMessage) *apperror.AppError {
	if len(raw) == 0 {
		return nil
	}

	query := fmt.Sprintf("INSERT OR REPLACE INTO %s (profile_name, %s) VALUES (?, ?)", table, col)
	execRes := tx.Exec(context.Background(), query, name, string(raw))

	return execRes.AppError()
}

func insertExtensionsAndBlobsToSQLite(tx *dbengine.TxWrapper, name string, extIDs []string) *apperror.AppError {
	ctx := context.Background()
	for _, id := range extIDs {
		if execRes := tx.Exec(ctx, "INSERT OR REPLACE INTO chrome_extensions (profile_name, extension_id) VALUES (?, ?)", name, id); execRes.IsFailure() {
			return execRes.AppError()
		}
	}

	return saveProfileBlobs(tx, name)
}

func saveProfileBlobs(tx *dbengine.TxWrapper, name string) *apperror.AppError {
	srcPath, _ := resolveChromeProfileDir(name)
	for _, blobFile := range constants.ChromeProfileSQLiteEntries {
		if err := saveBlobEntry(tx, name, srcPath, blobFile); err != nil {
			return err
		}
	}

	return nil
}

func saveBlobEntry(tx *dbengine.TxWrapper, name, srcPath, blobFile string) *apperror.AppError {
	filePath := filepath.Join(srcPath, blobFile)
	bytes, err := os.ReadFile(filePath)
	if err != nil || len(bytes) == 0 {
		return nil
	}

	query := "INSERT OR REPLACE INTO chrome_blobs (profile_name, file_name, payload) VALUES (?, ?, ?)"
	execRes := tx.Exec(context.Background(), query, name, blobFile, bytes)

	return execRes.AppError()
}

func writeAllChromeProfilesZIP(names []string, outPath string) (int, error) {
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	_ = addManifestToZip(zw, names)
	maxW := maxChromeProfileLabelWidth(names)
	for _, name := range names {
		writeProfileToZipWithProgress(zw, name, maxW)
	}

	_ = zw.Close()

	return getFileSize(outPath), nil
}

func writeProfileToZipWithProgress(zw *zip.Writer, name string, maxW int) {
	srcPath, hasDir := resolveChromeProfileDir(name)
	if !hasDir {
		return
	}

	_ = addProfileToZip(zw, srcPath, name)
	label := formatChromeProfileLabel(name, nil)
	pad := calculateLabelPadding(maxW, label)
	fmt.Printf("  \033[1;92m✓\033[0m %s%s → packaged in zip\n", label, strings.Repeat(" ", pad))
}

func addManifestToZip(zw *zip.Writer, names []string) error {
	m := chromeProfileManifest{
		GitMapVersion: constants.Version,
		SchemaVersion: chromeExportSchemaVersion,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		ProfileCount:  len(names),
		Profiles:      make([]chromeManifestProfile, 0, len(names)),
	}

	for _, name := range names {
		displayName, email := resolveProfileNameAndEmail(name, nil)
		srcPath, _ := resolveChromeProfileDir(name)
		extCount := len(listExtensionIDs(filepath.Join(srcPath, "Extensions")))
		m.Profiles = append(m.Profiles, chromeManifestProfile{
			Name:           name,
			DisplayName:    displayName,
			Email:          email,
			ExtensionCount: extCount,
		})
	}

	raw, err := json.MarshalIndent(m, "", constants.JSONIndent)
	if err != nil {
		return err
	}

	w, err := zw.Create("manifest.json")
	if err != nil {
		return err
	}

	_, err = w.Write(raw)

	return err
}

func addProfileToZip(zw *zip.Writer, srcPath, name string) error {
	tmpJSON := filepath.Join(os.TempDir(), name+"-snapshot.json")
	if _, err := writeChromeExport(srcPath, name, tmpJSON); err == nil {
		defer os.Remove(tmpJSON)
		_ = copyFileToZipPath(zw, tmpJSON, filepath.ToSlash(filepath.Join(name, name+".json")))
	}

	_ = copyFileToZipPath(zw, filepath.Join(srcPath, "Bookmarks"), filepath.ToSlash(filepath.Join(name, "Bookmarks")))
	_ = copyFileToZipPath(zw, filepath.Join(srcPath, "Preferences"), filepath.ToSlash(filepath.Join(name, "Preferences")))
	for _, dbName := range constants.ChromeProfileSQLiteEntries {
		_ = copyFileToZipPath(zw, filepath.Join(srcPath, dbName), filepath.ToSlash(filepath.Join(name, dbName)))
	}

	return nil
}

func copyFileToZipPath(zw *zip.Writer, srcPath, zipPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil || info.IsDir() {
		return nil
	}

	f, err := os.Open(srcPath)
	if err != nil {
		return nil
	}

	defer f.Close()
	w, err := zw.Create(zipPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)

	return err
}

func getFileSize(path string) int {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return int(info.Size())
}

func formatChromeProfileLabel(dirName string, prefsRaw json.RawMessage) string {
	displayName, email := resolveProfileNameAndEmail(dirName, prefsRaw)
	switch {
	case displayName != "" && displayName != dirName && displayName != email && email != "":
		return fmt.Sprintf("%s • %s (%s)", dirName, displayName, email)
	case email != "":
		return fmt.Sprintf("%s (%s)", dirName, email)
	case displayName != "" && displayName != dirName:
		return fmt.Sprintf("%s (%s)", dirName, displayName)
	default:
		return dirName
	}
}

func resolveProfileNameAndEmail(dirName string, prefsRaw json.RawMessage) (string, string) {
	displayName, email := fetchLocalStateProfileInfo(dirName)
	if email == "" && len(prefsRaw) > 0 {
		email = extractEmailFromPreferences(prefsRaw)
	}

	if email == "" {
		email = extractEmailFromProfileDisk(dirName)
	}

	return displayName, email
}

func fetchLocalStateProfileInfo(dirName string) (string, string) {
	state := readChromeLocalState()
	if state == nil {
		return "", ""
	}

	info, isFound := state.Profile.InfoCache[dirName]
	if !isFound {
		return "", ""
	}

	return info.Name, info.UserName
}

func extractEmailFromProfileDisk(dirName string) string {
	srcPath, hasDir := resolveChromeProfileDir(dirName)
	if !hasDir {
		return ""
	}

	prefBytes, err := os.ReadFile(filepath.Join(srcPath, "Preferences"))
	if err != nil || len(prefBytes) == 0 {
		return ""
	}

	return extractEmailFromPreferences(prefBytes)
}
