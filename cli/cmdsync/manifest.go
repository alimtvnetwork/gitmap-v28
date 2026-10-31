// Package cmdsync provides manifest generation and version.json metadata tracking.
package cmdsync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// SyncManifest represents the cryptographic inventory recorded in version.json.
type SyncManifest struct {
	ManifestVersion string            `json:"manifestVersion"`
	GeneratedAt     string            `json:"generatedAt"`
	SourceAuthority string            `json:"sourceAuthority"`
	AssetCount      int               `json:"assetCount"`
	SectionHashes   map[string]string `json:"sectionHashes,omitempty"`
	Assets          map[string]string `json:"assets"`
}

// NormalizedFileHash computes a stable SHA-256 digest ignoring CRLF differences.
func NormalizedFileHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	hasher := sha256.New()
	hasher.Write([]byte(normalized))

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// LoadTargetManifest reads syncManifest from a target repository's version.json.
func LoadTargetManifest(repoRoot string) (*SyncManifest, error) {
	verPath := filepath.Join(repoRoot, "version.json")
	data, err := os.ReadFile(verPath)
	if err != nil {
		return nil, err
	}

	var rootMap map[string]interface{}
	if err := json.Unmarshal(data, &rootMap); err != nil {
		return nil, err
	}

	rawManifest, hasManifest := rootMap["syncManifest"]
	if !hasManifest || rawManifest == nil {
		return &SyncManifest{Assets: make(map[string]string)}, nil
	}

	manifestBytes, _ := json.Marshal(rawManifest)
	var manifest SyncManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return &SyncManifest{Assets: make(map[string]string)}, nil
	}

	if manifest.Assets == nil {
		manifest.Assets = make(map[string]string)
	}

	return &manifest, nil
}

// SaveTargetManifest writes or updates syncManifest in a target repository's version.json.
func SaveTargetManifest(repoRoot string, manifest *SyncManifest) error {
	verPath := filepath.Join(repoRoot, "version.json")
	var rootMap map[string]interface{}

	data, err := os.ReadFile(verPath)
	if err == nil {
		_ = json.Unmarshal(data, &rootMap)
	}
	if rootMap == nil {
		rootMap = make(map[string]interface{})
	}

	manifest.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	manifest.ManifestVersion = "1.0.0"
	manifest.SourceAuthority = "gitmap"
	manifest.AssetCount = len(manifest.Assets)

	rootMap["syncManifest"] = manifest

	outBytes, err := json.MarshalIndent(rootMap, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(verPath, append(outBytes, '\n'), 0644)
}

// ComputeSectionHashes calculates Merkle roots for each managed folder section.
func ComputeSectionHashes(assets map[string]string) map[string]string {
	sectionFiles := make(map[string][]string)
	for relPath := range assets {
		section := resolveAssetSection(relPath)
		sectionFiles[section] = append(sectionFiles[section], relPath)
	}

	hashes := make(map[string]string)
	for sec, files := range sectionFiles {
		sort.Strings(files)
		hasher := sha256.New()
		for _, f := range files {
			hasher.Write([]byte(f + ":" + assets[f] + "\n"))
		}
		hashes[sec] = hex.EncodeToString(hasher.Sum(nil))
	}

	return hashes
}

func resolveAssetSection(relPath string) string {
	slashPath := filepath.ToSlash(relPath)
	switch {
	case strings.HasPrefix(slashPath, "01-prompts/"):
		return "prompts"
	case strings.HasPrefix(slashPath, ".agents/skills/") || strings.HasPrefix(slashPath, ".cursor/skills/"):
		return "skills"
	case strings.HasPrefix(slashPath, "02-spec/"):
		return "specs"
	case strings.HasPrefix(slashPath, "03-ai-scripts/") || strings.HasPrefix(slashPath, ".agents/scripts/"):
		return "scripts"
	default:
		return "other"
	}
}
