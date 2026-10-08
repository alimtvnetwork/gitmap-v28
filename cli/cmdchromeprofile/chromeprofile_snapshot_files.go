package cmdchromeprofile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"gopkg.in/yaml.v3"
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
