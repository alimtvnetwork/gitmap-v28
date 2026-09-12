package cmdchromeprofile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type DiscoveredProfileCandidate struct {
	SourcePath        string            `json:"source_path"`
	ProfileDirName    string            `json:"profile_dir_name"`
	DisplayName       string            `json:"display_name"`
	Email             string            `json:"email"`
	BookmarksCount    int               `json:"bookmarks_count"`
	ExtensionsCount   int               `json:"extensions_count"`
	IsFromZip         bool              `json:"is_from_zip"`
	ZipArchive        string            `json:"zip_archive,omitempty"`
	TargetDestination importDestination `json:"target_destination"`
}

func cleanCandidateTarget(target string) string {
	if target == "*.*" || target == "*" || target == "" {
		return "."
	}

	return filepath.Clean(target)
}

// DiscoverProfileCandidates finds all Chrome profile candidates in target file, zip, or directory.
func DiscoverProfileCandidates(target string) ([]DiscoveredProfileCandidate, error) {
	cleanTarget := cleanCandidateTarget(target)
	info, err := os.Stat(cleanTarget)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return discoverDirCandidates(cleanTarget)
	}

	if strings.HasSuffix(strings.ToLower(cleanTarget), constants.ExtZIP) {
		return discoverZipCandidates(cleanTarget)
	}

	return discoverSingleFileCandidate(cleanTarget)
}

func discoverDirCandidates(dir string) ([]DiscoveredProfileCandidate, error) {
	manifestCandidates, hasManifest := readManifestProfilesFromDir(dir)
	if hasManifest && len(manifestCandidates) > 0 {
		return manifestCandidates, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var candidates []DiscoveredProfileCandidate
	for _, e := range entries {
		c, ok := inspectDirEntryCandidate(dir, e)
		if ok {
			candidates = append(candidates, c)
		}
	}

	sortCandidates(candidates)

	return candidates, nil
}

func isIgnoredChromeProfileDir(name string) bool {
	if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "vendor" || name == "spec" || name == "tmp" {
		return true
	}

	return false
}

func inspectDirEntryCandidate(dir string, e os.DirEntry) (DiscoveredProfileCandidate, bool) {
	if isIgnoredChromeProfileDir(e.Name()) {
		return DiscoveredProfileCandidate{}, false
	}

	if e.IsDir() {
		return readSubdirProfileCandidate(dir, e.Name())
	}

	if isExcludedSystemJSON(strings.ToLower(e.Name())) {
		return DiscoveredProfileCandidate{}, false
	}

	p := filepath.Join(dir, e.Name())
	if !isSnapshotFileExtension(e.Name()) || !isValidChromeSnapshotFile(p) {
		return DiscoveredProfileCandidate{}, false
	}

	items, err := discoverSingleFileCandidate(p)
	if err != nil || len(items) == 0 {
		return DiscoveredProfileCandidate{}, false
	}

	return items[0], true
}

func readSubdirProfileCandidate(dir, subName string) (DiscoveredProfileCandidate, bool) {
	subPath := filepath.Join(dir, subName)
	jsonPath := filepath.Join(subPath, subName+".json")
	if !chromeProfilePathExists(jsonPath) {
		jsonPath = findFirstJSONInSubdir(subPath)
	}

	if cand, ok := probeSubdirJSONCandidate(subPath, subName, jsonPath); ok {
		return cand, true
	}

	return probeSubdirDiskCandidate(subPath, subName)
}

func probeSubdirJSONCandidate(subPath, subName, jsonPath string) (DiscoveredProfileCandidate, bool) {
	if jsonPath == "" || !isValidChromeSnapshotFile(jsonPath) {
		return DiscoveredProfileCandidate{}, false
	}

	items, err := discoverSingleFileCandidate(jsonPath)
	if err != nil || len(items) == 0 {
		return DiscoveredProfileCandidate{}, false
	}

	items[0].SourcePath = subPath
	items[0].ProfileDirName = subName

	return items[0], true
}

func findFirstJSONInSubdir(subPath string) string {
	entries, err := os.ReadDir(subPath)
	if err != nil {
		return ""
	}

	for _, e := range entries {
		lower := strings.ToLower(e.Name())
		if e.IsDir() || !strings.HasSuffix(lower, constants.ExtJSON) || isExcludedSystemJSON(lower) {
			continue
		}

		cand := filepath.Join(subPath, e.Name())
		if isValidChromeSnapshotFile(cand) {
			return cand
		}
	}

	return ""
}

func probeSubdirDiskCandidate(subPath, subName string) (DiscoveredProfileCandidate, bool) {
	hasPrefs := chromeProfilePathExists(filepath.Join(subPath, "Preferences"))
	hasBM := chromeProfilePathExists(filepath.Join(subPath, "Bookmarks"))
	if !hasPrefs && !hasBM {
		return DiscoveredProfileCandidate{}, false
	}

	bmCount := 0
	if rawBM, err := os.ReadFile(filepath.Join(subPath, "Bookmarks")); err == nil {
		bmCount = countBookmarks(rawBM)
	}

	extCount := len(listExtensionIDs(filepath.Join(subPath, "Extensions")))
	email, disp := "", subName
	if rawPref, err := os.ReadFile(filepath.Join(subPath, "Preferences")); err == nil {
		email = extractEmailFromPreferences(rawPref)
	}

	exp := &chromeExport{Name: subName, DisplayName: disp, Email: email}
	target := resolveImportDestination(exp, "", false)

	return DiscoveredProfileCandidate{
		SourcePath:        subPath,
		ProfileDirName:    subName,
		DisplayName:       disp,
		Email:             email,
		BookmarksCount:    bmCount,
		ExtensionsCount:   extCount,
		TargetDestination: target,
	}, true
}
