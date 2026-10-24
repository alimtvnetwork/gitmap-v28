// Package downloaderconfig — types.go centralizes downloader configuration models, keys, and Result aliases.
package downloaderconfig

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/downloaderconfig/enums"

// ConfigKeyType aliases the strongly-typed enum from the enums package.
type ConfigKeyType = enums.ConfigKeyType

const (
	KeyDefaultSplitSize   = enums.KeyDefaultSplitSize
	KeyLargeFileSplitSize = enums.KeyLargeFileSplitSize
	KeyLargeFileThreshold = enums.KeyLargeFileThreshold
	KeyTinyFileThreshold  = enums.KeyTinyFileThreshold
	KeyTinyFileSplitSize  = enums.KeyTinyFileSplitSize
)

// ConfigKey is a backward-compatible alias for ConfigKeyType.
type ConfigKey = enums.ConfigKeyType

// Document is the top-level Seedable-Config envelope. Field names are
// PascalCase to match the spec and the JSON file shipped under
// cli/data/downloader-config.json.
type Document struct {
	DownloaderConfig DownloaderConfig `json:",omitempty"`
	DatabaseVersion  DatabaseVersion  `json:",omitempty"`
}

// DownloaderConfig is the per-downloader runtime config consumed by
// Slice 2 (aria2c installer + engine).
type DownloaderConfig struct {
	PreferredDownloader string `json:",omitempty"`
	FallbackDownloader  string `json:",omitempty"`
	ParallelDownloads   int    `json:",omitempty"`
	SplitConnections    int    `json:",omitempty"`
	DefaultSplitSize    string `json:",omitempty"`
	LargeFileSplitSize  string `json:",omitempty"`
	LargeFileThreshold  string `json:",omitempty"`
	TinyFileThreshold   string `json:",omitempty"`
	TinyFileSplitSize   string `json:",omitempty"`
	TinyFileSplits      int    `json:",omitempty"`
	AllowFallback       bool   `json:",omitempty"`
	OverwriteUserConfig bool   `json:",omitempty"`
}

// DatabaseVersion records the last gitmap version that touched the DB.
// Stored as a string so we can keep the literal "auto" sentinel in the
// shipped seed file and resolve it at apply-time to constants.Version.
type DatabaseVersion struct {
	LastKnownVersion string `json:",omitempty"`
}

// DocumentResult wraps a Document in a Result envelope.
type DocumentResult = result.Wrap[Document]

// BytesResult wraps a byte slice in a Result envelope.
type BytesResult = result.Wrap[[]byte]
