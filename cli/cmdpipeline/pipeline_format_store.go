package cmdpipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveFormatsDir() string {
	return store.ResolveSplitDbDir(store.SectionPipeline, "formats", "")
}

// SavePEFormatProfile saves a format profile as JSON in the formats directory.
func SavePEFormatProfile(p PEFormatProfile) error {
	if strings.TrimSpace(p.Name) == "" {
		return apperror.NewValidationError("format profile name cannot be empty")
	}
	dir := resolveFormatsDir()
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal format profile")
	}
	filePath := filepath.Join(dir, p.Name+".json")
	if writeErr := os.WriteFile(filePath, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write format profile "+p.Name)
	}
	saveAliasProfileCopy(dir, p, data)

	return nil
}

func saveAliasProfileCopy(dir string, p PEFormatProfile, data []byte) {
	if strings.TrimSpace(p.Alias) == "" || strings.EqualFold(p.Alias, p.Name) {
		return
	}
	aliasPath := filepath.Join(dir, p.Alias+".json")
	_ = os.WriteFile(aliasPath, data, 0644)
}

// LoadPEFormatProfile loads a profile by direct path, registered name, or alias.
func LoadPEFormatProfile(target string) (*PEFormatProfile, error) {
	clean := strings.TrimSpace(target)
	if clean == "" || clean == "default" {
		p := DefaultPEFormatProfile()
		return &p, nil
	}
	if clean == "tauri" || clean == "rust" {
		p := BuiltinTauriProfile()
		return &p, nil
	}
	if profile, isFile := tryLoadProfileFromFile(clean); isFile {
		return profile, nil
	}
	return loadProfileFromFormatsDir(clean)
}

func tryLoadProfileFromFile(path string) (*PEFormatProfile, bool) {
	if !strings.HasSuffix(strings.ToLower(path), ".json") && !fileExists(path) {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var p PEFormatProfile
	if jsonErr := json.Unmarshal(data, &p); jsonErr != nil {
		return nil, false
	}
	return &p, true
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func loadProfileFromFormatsDir(nameOrAlias string) (*PEFormatProfile, error) {
	dir := resolveFormatsDir()
	candidate := filepath.Join(dir, nameOrAlias+".json")
	profile, hasLoaded := tryLoadCandidateProfile(candidate)
	if hasLoaded {
		return profile, nil
	}

	return scanProfilesForMatch(dir, nameOrAlias)
}

func tryLoadCandidateProfile(candidate string) (*PEFormatProfile, bool) {
	data, err := os.ReadFile(candidate)
	if err != nil {
		return nil, false
	}
	var p PEFormatProfile
	if jsonErr := json.Unmarshal(data, &p); jsonErr != nil {
		return nil, false
	}

	return &p, true
}

func scanProfilesForMatch(dir, target string) (*PEFormatProfile, error) {
	profiles, err := ListPEFormatProfiles()
	if err != nil {
		return nil, err
	}
	for _, p := range profiles {
		if strings.EqualFold(p.Name, target) || strings.EqualFold(p.Alias, target) {
			return &p, nil
		}
	}
	return nil, apperror.NewNotFoundError("format profile not found: " + target)
}

// DeletePEFormatProfile deletes a profile and its alias from the formats directory.
func DeletePEFormatProfile(target string) error {
	dir := resolveFormatsDir()
	candidate := filepath.Join(dir, target+".json")
	if fileExists(candidate) {
		_ = os.Remove(candidate)
	}
	return removeMatchedProfiles(dir, target)
}

func removeMatchedProfiles(dir, target string) error {
	profiles, err := ListPEFormatProfiles()
	if err != nil {
		return err
	}
	hasRemoved := false
	for _, p := range profiles {
		isMatch := strings.EqualFold(p.Name, target) || strings.EqualFold(p.Alias, target)
		if isMatch {
			removeProfileFiles(dir, p)
			hasRemoved = true
		}
	}
	if !hasRemoved {
		return apperror.NewNotFoundError("profile not found: " + target)
	}

	return nil
}

func removeProfileFiles(dir string, p PEFormatProfile) {
	_ = os.Remove(filepath.Join(dir, p.Name+".json"))
	if p.Alias != "" {
		_ = os.Remove(filepath.Join(dir, p.Alias+".json"))
	}
}

// ListPEFormatProfiles returns all built-in and registered profiles.
func ListPEFormatProfiles() ([]PEFormatProfile, error) {
	profiles := []PEFormatProfile{
		DefaultPEFormatProfile(),
		BuiltinTauriProfile(),
	}
	dir := resolveFormatsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return profiles, nil
	}
	for _, e := range entries {
		appendProfileFromEntry(&profiles, dir, e)
	}
	return deduplicateProfiles(profiles), nil
}

func appendProfileFromEntry(profiles *[]PEFormatProfile, dir string, e os.DirEntry) {
	if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
		return
	}
	data, err := os.ReadFile(filepath.Join(dir, e.Name()))
	if err != nil {
		return
	}
	var p PEFormatProfile
	if jsonErr := json.Unmarshal(data, &p); jsonErr == nil && p.Name != "" {
		*profiles = append(*profiles, p)
	}
}

func deduplicateProfiles(list []PEFormatProfile) []PEFormatProfile {
	seen := make(map[string]bool)
	var out []PEFormatProfile
	for _, p := range list {
		key := strings.ToLower(p.Name)
		if !seen[key] {
			seen[key] = true
			out = append(out, p)
		}
	}
	return out
}

// ImportPEFormatProfilesFromFolder imports all .json files in a folder.
func ImportPEFormatProfilesFromFolder(folder string) (int, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return 0, apperror.WrapSimple(err, "read directory "+folder)
	}
	count := 0
	for _, e := range entries {
		if tryImportEntry(folder, e) {
			count++
		}
	}

	return count, nil
}

func tryImportEntry(folder string, e os.DirEntry) bool {
	isJsonFile := !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".json")
	if !isJsonFile {
		return false
	}

	return importSingleFile(filepath.Join(folder, e.Name()))
}

func importSingleFile(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var p PEFormatProfile
	if jsonErr := json.Unmarshal(data, &p); jsonErr != nil {
		return false
	}
	if p.Name == "" {
		p.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	return SavePEFormatProfile(p) == nil
}
