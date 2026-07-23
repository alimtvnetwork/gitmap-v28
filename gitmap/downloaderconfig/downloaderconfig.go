// Package downloaderconfig owns the on-disk + in-DB representation of the
// gitmap downloader configuration.
//
// Slice 1 of the downloader feature only persists the config — actual
// download / install logic ships in later slices. Keeping the data layer
// isolated here means Slice 2 (aria2c installer + engine) and Slice 3
// (download / download-unzip commands) can both depend on a stable Load()
// without re-implementing JSON parsing.
//
// Storage model: a single JSON document under Setting[DownloaderConfig].
// We deliberately do NOT introduce a new SettingTypes table — the existing
// Setting(Key TEXT PK, Value TEXT) shape from constants_settings.go is
// reused, and the type discriminator lives in code as constants.SettingType.
package downloaderconfig

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Document is the top-level Seedable-Config envelope. Field names are
// PascalCase to match the spec and the JSON file shipped under
// gitmap/data/downloader-config.json.
type Document struct {
	DownloaderConfig DownloaderConfig `json:"DownloaderConfig"`
	DatabaseVersion  DatabaseVersion  `json:"DatabaseVersion"`
}

// DownloaderConfig is the per-downloader runtime config consumed by
// Slice 2 (aria2c installer + engine).
type DownloaderConfig struct {
	PreferredDownloader string `json:"PreferredDownloader"`
	FallbackDownloader  string `json:"FallbackDownloader"`
	ParallelDownloads   int    `json:"ParallelDownloads"`
	SplitConnections    int    `json:"SplitConnections"`
	DefaultSplitSize    string `json:"DefaultSplitSize"`
	LargeFileSplitSize  string `json:"LargeFileSplitSize"`
	LargeFileThreshold  string `json:"LargeFileThreshold"`
	TinyFileThreshold   string `json:"TinyFileThreshold"`
	TinyFileSplitSize   string `json:"TinyFileSplitSize"`
	TinyFileSplits      int    `json:"TinyFileSplits"`
	AllowFallback       bool   `json:"AllowFallback"`
	OverwriteUserConfig bool   `json:"OverwriteUserConfig"`
}

// DatabaseVersion records the last gitmap version that touched the DB.
// Stored as a string so we can keep the literal "auto" sentinel in the
// shipped seed file and resolve it at apply-time to constants.Version.
type DatabaseVersion struct {
	LastKnownVersion string `json:"LastKnownVersion"`
}

// Defaults returns a Document populated from the hard-coded constants.
// Used as the last-resort fallback when both the DB and the seed file are
// unavailable (e.g. first-run race before Migrate completes).
func Defaults() Document {
	return Document{
		DownloaderConfig: DownloaderConfig{
			PreferredDownloader: constants.DownloaderDefaultPreferred,
			FallbackDownloader:  constants.DownloaderDefaultFallback,
			ParallelDownloads:   constants.DownloaderDefaultParallel,
			SplitConnections:    constants.DownloaderDefaultSplits,
			DefaultSplitSize:    constants.DownloaderDefaultSplitSize,
			LargeFileSplitSize:  constants.DownloaderDefaultLargeSplitSize,
			LargeFileThreshold:  constants.DownloaderDefaultLargeThreshold,
			TinyFileThreshold:   constants.DownloaderDefaultTinyThreshold,
			TinyFileSplitSize:   constants.DownloaderDefaultTinySplitSize,
			TinyFileSplits:      constants.DownloaderDefaultTinySplits,
			AllowFallback:       constants.DownloaderDefaultAllowFallback,
			OverwriteUserConfig: constants.DownloaderDefaultOverwriteUser,
		},
		DatabaseVersion: DatabaseVersion{LastKnownVersion: constants.Version},
	}
}

// LoadFile reads + validates a Seedable-Config JSON file from disk.
// Used by `gitmap downloader-config <path>` and by the seeder.
func LoadFile(path string) (Document, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		// Preserve the underlying error (in particular fs.ErrNotExist) so
		// callers using os.IsNotExist / errors.Is can distinguish "missing
		// optional seed" from a real I/O failure. Without %w the seeder
		// printed a spurious "Could not read downloader seed" warning on
		// every fresh install where the seed file does not yet exist.
		return Document{}, fmt.Errorf(constants.ErrDownloaderConfigPathRequired+": %w", path, err)
	}

	return Parse(raw)
}

// Parse validates a raw JSON byte slice and returns the typed Document.
func Parse(raw []byte) (Document, error) {
	var doc Document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return Document{}, fmt.Errorf(constants.ErrDownloaderConfigInvalidJSON, err)
	}

	if err := Validate(doc); err != nil {
		return Document{}, err
	}

	// "auto" is a documented sentinel in the shipped seed: resolve to the
	// running binary's version so SettingDatabaseVersion always carries a
	// real semver string, never the literal "auto".
	if doc.DatabaseVersion.LastKnownVersion == "" || doc.DatabaseVersion.LastKnownVersion == "auto" {
		doc.DatabaseVersion.LastKnownVersion = constants.Version
	}

	return doc, nil
}

// Validate enforces the PascalCase + range rules. Required keys are
// checked first so error messages surface a missing key before a
// numeric range violation that may be a side effect.
func Validate(doc Document) error {
	dc := doc.DownloaderConfig
	if dc.PreferredDownloader == "" {
		return fmt.Errorf(constants.ErrDownloaderConfigMissingKey, "DownloaderConfig.PreferredDownloader")
	}
	if dc.FallbackDownloader == "" {
		return fmt.Errorf(constants.ErrDownloaderConfigMissingKey, "DownloaderConfig.FallbackDownloader")
	}
	if dc.ParallelDownloads < 1 || dc.ParallelDownloads > 64 {
		return fmt.Errorf(constants.ErrDownloaderConfigBadParallel, dc.ParallelDownloads)
	}
	if dc.SplitConnections < 1 || dc.SplitConnections > 64 {
		return fmt.Errorf(constants.ErrDownloaderConfigBadSplits, dc.SplitConnections)
	}
	for k, v := range map[string]string{
		"DownloaderConfig.DefaultSplitSize":   dc.DefaultSplitSize,
		"DownloaderConfig.LargeFileSplitSize": dc.LargeFileSplitSize,
		"DownloaderConfig.LargeFileThreshold": dc.LargeFileThreshold,
		"DownloaderConfig.TinyFileThreshold":  dc.TinyFileThreshold,
		"DownloaderConfig.TinyFileSplitSize":  dc.TinyFileSplitSize,
	} {
		if v == "" {
			return fmt.Errorf(constants.ErrDownloaderConfigMissingKey, k)
		}
	}

	return nil
}

// Marshal serializes a Document with deterministic 2-space indent, matching
// the project's JSONIndent convention so files written back round-trip
// cleanly with the seed.
func Marshal(doc Document) ([]byte, error) {
	return json.MarshalIndent(doc, "", constants.JSONIndent)
}

// SeedHash returns the SHA-256 of the canonical (re-marshaled) document.
// Hashing the re-marshaled form (not the raw bytes) means whitespace-only
// edits to the seed file do not falsely trigger a re-seed.
func SeedHash(doc Document) string {
	canon, err := Marshal(doc)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(canon)

	return hex.EncodeToString(sum[:])
}
