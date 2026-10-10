package cmdchromeprofile

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func isValidChromeSnapshotFile(path string) bool {
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, constants.ExtJSON) {
		return isChromeSnapshotJSON(path)
	}

	if strings.HasSuffix(lower, constants.ExtZIP) {
		return true
	}

	if isSQLiteSnapshot(lower) {
		return true
	}

	if strings.HasSuffix(lower, constants.ExtYAML) || strings.HasSuffix(lower, constants.ExtYML) {
		return isChromeSnapshotYAML(path)
	}

	return false
}

func isImportableSnapshot(path string) bool {
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, constants.ExtJSON) || strings.HasSuffix(lower, constants.ExtZIP) {
		return true
	}

	if strings.HasSuffix(lower, constants.ExtYAML) || strings.HasSuffix(lower, constants.ExtYML) {
		return true
	}

	return isSQLiteSnapshot(lower)
}

func isSQLiteSnapshot(lower string) bool {
	if strings.HasSuffix(lower, constants.ExtDB) || strings.HasSuffix(lower, constants.ExtSQLite) {
		return true
	}

	return strings.HasSuffix(lower, ".sqlite3")
}

func isChromeSnapshotJSON(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	if isExcludedSystemJSON(base) {
		return false
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err == nil && isExportPopulated(&exp) {
		return true
	}

	var all chromeAllProfilesExport
	if err := json.Unmarshal(raw, &all); err == nil && len(all.Profiles) > 0 {
		return true
	}

	if isChromeProfileNamedJSON(base) {
		var obj map[string]any

		return json.Unmarshal(raw, &obj) == nil && len(obj) > 0
	}

	return false
}

func isExcludedSystemJSON(base string) bool {
	switch strings.ToLower(base) {
	case "manifest.json", "package.json", "package-lock.json", "tsconfig.json", "version.json", "composer.json":
		return true
	}

	return false
}

func isChromeProfileNamedJSON(base string) bool {
	if base == "default.json" {
		return true
	}

	return strings.HasPrefix(base, "profile") && strings.HasSuffix(base, ".json")
}

func isExportPopulated(exp *chromeExport) bool {
	return exp.SchemaVersion > 0 ||
		exp.ExportedAt != "" ||
		exp.Email != "" ||
		exp.DisplayName != "" ||
		len(exp.Bookmarks) > 0 ||
		len(exp.Preferences) > 0 ||
		len(exp.ExtensionIDs) > 0
}

func isChromeSnapshotYAML(path string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var exp chromeExport
	if err := yaml.Unmarshal(raw, &exp); err == nil {
		return exp.SchemaVersion > 0 || exp.ExportedAt != ""
	}

	return false
}

func registerImportedProfileToLocalState(dstDir, displayName, email string) error {
	dstPath := chromeProfilePath(dstDir)
	prefsRaw := readOptionalJSON(filepath.Join(dstPath, "Preferences"))
	gaiaInfo := resolveProfileGaiaInfo(dstDir, prefsRaw, dstPath)
	if email == "" && gaiaInfo.Email != "" {
		email = gaiaInfo.Email
	}

	return registerChromeProfileWithFullSchemaAndGAIA(
		dstDir,
		displayName,
		email,
		gaiaInfo.GaiaID,
		gaiaInfo.GaiaName,
		gaiaInfo.GaiaGivenName,
	)
}

func extractEmailIfMissing(exp *chromeExport, raw []byte) {
	if exp.Email == "" && len(exp.Preferences) > 0 {
		exp.Email = extractEmailFromPreferences(exp.Preferences)
	}

	if exp.Email == "" && len(raw) > 0 {
		exp.Email = extractEmailFromPreferences(json.RawMessage(raw))
	}

	if exp.Email == "" && len(raw) > 0 {
		exp.Email = extractEmailFromRawString(raw)
	}
}

var rawEmailRegex = regexp.MustCompile(`"(?:email|user_name|username)"\s*:\s*"([^"\\]+@[^"\\]+\.[a-zA-Z]{2,})"`)

func extractEmailFromRawString(raw []byte) string {
	matches := rawEmailRegex.FindSubmatch(raw)
	if len(matches) > 1 {
		return string(matches[1])
	}

	return ""
}

func readYAMLSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	raw, err := os.ReadFile(srcFile)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", srcFile, err)
	}

	var exp chromeExport
	if err := yaml.Unmarshal(raw, &exp); err != nil {
		return nil, fmt.Errorf("parse %s: %w", srcFile, err)
	}

	extractEmailIfMissing(&exp, raw)
	info, _ := os.Stat(srcFile)
	var size int64
	if info != nil {
		size = info.Size()
	}

	target := resolveImportDestination(&exp, "", false)

	return &snapshotMetadata{
		FilePath:          srcFile,
		FileName:          filepath.Base(srcFile),
		FileSize:          size,
		Export:            &exp,
		BookmarksCount:    countBookmarks(exp.Bookmarks),
		ExtensionsCount:   len(exp.ExtensionIDs),
		HasEmail:          exp.Email != "",
		Email:             exp.Email,
		DisplayName:       exp.DisplayName,
		ProfileName:       exp.Name,
		ExportedAt:        exp.ExportedAt,
		TargetDestination: target,
	}, nil
}

func readGenericSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	base := filepath.Base(srcFile)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	info, _ := os.Stat(srcFile)
	var size int64
	if info != nil {
		size = info.Size()
	}

	exp := &chromeExport{Name: name, DisplayName: name}
	target := resolveImportDestination(exp, "", false)

	return &snapshotMetadata{
		FilePath:          srcFile,
		FileName:          base,
		FileSize:          size,
		Export:            exp,
		BookmarksCount:    0,
		ExtensionsCount:   0,
		HasEmail:          false,
		Email:             "",
		DisplayName:       name,
		ProfileName:       name,
		ExportedAt:        "",
		TargetDestination: target,
	}, nil
}

func resolveDirSnapshotPath(p, name string) string {
	candJSON := filepath.Join(p, name+".json")
	if chromeProfilePathExists(candJSON) {
		return candJSON
	}

	candFirst := findFirstJSONInSubdir(p)
	if candFirst != "" {
		return candFirst
	}

	if chromeProfilePathExists(filepath.Join(p, "Preferences")) || chromeProfilePathExists(filepath.Join(p, "Bookmarks")) {
		return p
	}

	return ""
}

func resolveSnapshotEntryPath(dir string, e os.DirEntry) string {
	if isIgnoredChromeProfileDir(e.Name()) {
		return ""
	}

	p := filepath.Join(dir, e.Name())
	if e.IsDir() {
		return resolveDirSnapshotPath(p, e.Name())
	}

	if isImportableSnapshot(p) && isValidChromeSnapshotFile(p) {
		return p
	}

	return ""
}

func scanSnapshotFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var files []string
	for _, e := range entries {
		p := resolveSnapshotEntryPath(dir, e)
		if p != "" {
			files = append(files, p)
		}
	}

	sort.Strings(files)

	return files
}
