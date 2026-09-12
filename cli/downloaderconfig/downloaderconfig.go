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

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

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
func LoadFile(path string) DocumentResult {
	raw, err := os.ReadFile(path)
	if err != nil {
		// Preserve the underlying error (in particular fs.ErrNotExist) so
		// callers using os.IsNotExist / errors.Is can distinguish "missing
		// optional seed" from a real I/O failure. Without %w the seeder
		// printed a spurious "Could not read downloader seed" warning on
		// every fresh install where the seed file does not yet exist.
		appErr := apperror.WrapSimple(err, constants.ErrDownloaderConfigPathRequired).
			WithContext("path", path)

		return result.FailureResult[Document](appErr)
	}

	return Parse(raw)
}

// Parse validates a raw JSON byte slice and returns the typed Document.
func Parse(raw []byte) DocumentResult {
	var doc Document
	if err := json.Unmarshal(raw, &doc); err != nil {
		appErr := apperror.WrapSimple(err, constants.ErrDownloaderConfigInvalidJSON)

		return result.FailureResult[Document](appErr)
	}

	if err := Validate(doc); err != nil {
		return result.FailureResult[Document](err)
	}

	// "auto" is a documented sentinel in the shipped seed: resolve to the
	// running binary's version so SettingDatabaseVersion always carries a
	// real semver string, never the literal "auto".
	if doc.DatabaseVersion.LastKnownVersion == "" || doc.DatabaseVersion.LastKnownVersion == "auto" {
		doc.DatabaseVersion.LastKnownVersion = constants.Version
	}

	return result.SuccessResult(doc)
}

// Validate enforces the PascalCase + range rules. Required keys are
// checked first so error messages surface a missing key before a
// numeric range violation that may be a side effect.
func Validate(doc Document) *apperror.AppError {
	dc := doc.DownloaderConfig
	if dc.PreferredDownloader == "" {
		return apperror.NewSimple("downloader-config", "E1000").
			WithContext("error", fmt.Sprintf(constants.ErrDownloaderConfigMissingKey, "DownloaderConfig.PreferredDownloader"))
	}

	if dc.FallbackDownloader == "" {
		return apperror.NewSimple("downloader-config", "E1000").
			WithContext("error", fmt.Sprintf(constants.ErrDownloaderConfigMissingKey, "DownloaderConfig.FallbackDownloader"))
	}

	if dc.ParallelDownloads < 1 || dc.ParallelDownloads > 64 {
		return apperror.NewSimple("downloader-config", "E1000").
			WithContext("error", fmt.Sprintf(constants.ErrDownloaderConfigBadParallel, dc.ParallelDownloads))
	}

	if dc.SplitConnections < 1 || dc.SplitConnections > 64 {
		return apperror.NewSimple("downloader-config", "E1000").
			WithContext("error", fmt.Sprintf(constants.ErrDownloaderConfigBadSplits, dc.SplitConnections))
	}

	for k, v := range map[ConfigKey]string{
		KeyDefaultSplitSize:   dc.DefaultSplitSize,
		KeyLargeFileSplitSize: dc.LargeFileSplitSize,
		KeyLargeFileThreshold: dc.LargeFileThreshold,
		KeyTinyFileThreshold:  dc.TinyFileThreshold,
		KeyTinyFileSplitSize:  dc.TinyFileSplitSize,
	} {
		if v == "" {
			return apperror.NewSimple("downloader-config", "E1000").
				WithContext("error", fmt.Sprintf(constants.ErrDownloaderConfigMissingKey, string(k)))
		}
	}

	return nil
}

// Marshal serializes a Document with deterministic 2-space indent, matching
// the project's JSONIndent convention so files written back round-trip
// cleanly with the seed.
func Marshal(doc Document) BytesResult {
	b, err := json.MarshalIndent(doc, "", constants.JSONIndent)
	if err != nil {
		appErr := apperror.WrapSimple(err, "downloaderconfig.Marshal")

		return result.FailureResult[[]byte](appErr)
	}

	return result.SuccessResult(b)
}

// SeedHash returns the SHA-256 of the canonical (re-marshaled) document.
// Hashing the re-marshaled form (not the raw bytes) means whitespace-only
// edits to the seed file do not falsely trigger a re-seed.
func SeedHash(doc Document) string {
	res := Marshal(doc)
	if res.IsFailure() {
		return ""
	}

	canon, _ := res.Unwrap()
	sum := sha256.Sum256(canon)

	return hex.EncodeToString(sum[:])
}
