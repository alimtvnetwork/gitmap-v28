// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ReleaseMeta holds metadata for a single release.
// v15: boolean fields use the IsX prefix convention. ReadReleaseMeta below
// adds backward-compat aliases so old "draft"/"preRelease" JSON files load
// without breakage; new files are written using the IsX form.
type ReleaseMeta struct {
	Version           string            `json:"version"`
	Branch            string            `json:"branch"`
	SourceBranch      string            `json:"sourceBranch"`
	Commit            string            `json:"commit"`
	Tag               string            `json:"tag"`
	Assets            []string          `json:"assets"`
	Changelog         []string          `json:"changelog,omitempty"`
	Notes             string            `json:"notes,omitempty"`
	ZipGroups         []string          `json:"zipGroups,omitempty"`
	ZipGroupChecksums map[string]string `json:"zipGroupChecksums,omitempty"`
	IsDraft           bool              `json:"isDraft"`
	IsPreRelease      bool              `json:"isPreRelease"`
	CreatedAt         string            `json:"createdAt"`
	IsLatest          bool              `json:"isLatest"`
}

// LatestMeta holds the pointer to the latest stable release.
type LatestMeta struct {
	Version string `json:"version"`
	Tag     string `json:"tag"`
	Branch  string `json:"branch"`
}

// VersionFile represents the version.json file format.
type VersionFile struct {
	Version string `json:"version"`
}

// ReleaseExists checks if a release metadata file already exists.
func ReleaseExists(version Version) bool {
	path := metaFilePath(version)
	_, err := os.Stat(path)

	return err == nil
}

// metaFilePath returns the path for a release metadata file.
func metaFilePath(v Version) string {
	filename := v.String() + constants.ExtJSON

	return filepath.Join(constants.DefaultReleaseDir, filename)
}

// WriteReleaseMeta writes release metadata to .gitmap/release/vX.Y.Z.json.
func WriteReleaseMeta(meta ReleaseMeta) error {
	err := os.MkdirAll(constants.DefaultReleaseDir, constants.DirPermission)
	if err != nil {
		return fmt.Errorf("create release dir: %w", err)
	}

	v, err := Parse(meta.Tag)
	if err != nil {
		return fmt.Errorf("parse tag for path: %w", err)
	}

	return writeJSON(metaFilePath(v), meta)
}

// WriteLatest updates .gitmap/release/latest.json if the version is the highest stable.
func WriteLatest(v Version) error {
	if latestIsHigher(v) {
		return nil
	}

	latest := LatestMeta{
		Version: v.CoreString(),
		Tag:     v.String(),
		Branch:  constants.ReleaseBranchPrefix + v.String(),
	}

	return writeJSON(latestFilePath(), latest)
}

// latestIsHigher returns true when the existing latest version is >= candidate.
func latestIsHigher(candidate Version) bool {
	current, err := ReadLatest()
	if err != nil {
		return false
	}
	currentVer, parseErr := Parse(current.Tag)
	if parseErr != nil {
		return false
	}
	if candidate.GreaterThan(currentVer) {
		return false
	}

	return true
}

// ReadLatest reads .gitmap/release/latest.json.
func ReadLatest() (LatestMeta, error) {
	data, err := os.ReadFile(latestFilePath())
	if err != nil {
		return LatestMeta{}, err
	}

	var latest LatestMeta
	err = json.Unmarshal(data, &latest)

	return latest, err
}

// ReadVersionFile reads version.json from the project root.
func ReadVersionFile() (string, error) {
	data, err := os.ReadFile(constants.DefaultVersionFile)
	if err != nil {
		return "", err
	}

	var vf VersionFile
	err = json.Unmarshal(data, &vf)
	if err != nil {
		return "", err
	}

	return vf.Version, nil
}

// ReadReleaseMeta reads and unmarshals a single .gitmap/release/vX.Y.Z.json file.
// Provides backward-compat for the pre-v15 "draft" / "preRelease" JSON keys
// by overlaying a legacy struct after the primary unmarshal.
func ReadReleaseMeta(path string) (ReleaseMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ReleaseMeta{}, err
	}

	var meta ReleaseMeta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return meta, err
	}

	// Legacy field names (v3.4.x and earlier). Only override when the new
	// IsX field is the zero value AND the legacy field is set, so a
	// re-saved (post-v15) file is never downgraded.
	var legacy struct {
		Draft      *bool `json:"draft"`
		PreRelease *bool `json:"preRelease"`
	}
	if jsonErr := json.Unmarshal(data, &legacy); jsonErr == nil {
		if !meta.IsDraft && legacy.Draft != nil && *legacy.Draft {
			meta.IsDraft = true
		}
		if !meta.IsPreRelease && legacy.PreRelease != nil && *legacy.PreRelease {
			meta.IsPreRelease = true
		}
	}

	return meta, nil
}

// ListReleaseMetaFiles reads all .gitmap/release/v*.json files and returns parsed metadata.
func ListReleaseMetaFiles() ([]ReleaseMeta, error) {
	pattern := filepath.Join(constants.DefaultReleaseDir, constants.ReleaseGlob)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	results := make([]ReleaseMeta, 0, len(matches))

	for _, path := range matches {
		if filepath.Base(path) == constants.DefaultLatestFile {
			continue
		}

		meta, readErr := ReadReleaseMeta(path)
		if readErr != nil {
			continue
		}

		results = append(results, meta)
	}

	return results, nil
}

// latestFilePath returns the path to latest.json.
func latestFilePath() string {
	return filepath.Join(constants.DefaultReleaseDir, constants.DefaultLatestFile)
}

// writeJSON marshals data to a JSON file with indentation.
func writeJSON(path string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", constants.JSONIndent)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, constants.DirPermission)
}
