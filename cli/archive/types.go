// Package archive — types.go centralizes all archive formats, compression modes, domain models, and Result envelopes.
package archive

import (
	"context"
	"io"
	"io/fs"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/mholt/archives"
)

// FormatType is a string tag persisted in ArchiveHistory.ArchiveFormat. It
// reads cleanly in PascalCase logs ("Zip", "TarGz") yet round-trips
// through the canonical extension via FormatFromExt / FormatType.Extension.
type FormatType string

// Format is retained as a type alias for backwards compatibility.
type Format = FormatType

const (
	FormatZip     FormatType = "Zip"
	FormatTar     FormatType = "Tar"
	FormatTarGz   FormatType = "TarGz"
	FormatTarBz2  FormatType = "TarBz2"
	FormatTarXz   FormatType = "TarXz"
	FormatTarZst  FormatType = "TarZst"
	FormatGz      FormatType = "Gz"
	FormatBz2     FormatType = "Bz2"
	FormatXz      FormatType = "Xz"
	FormatZst     FormatType = "Zst"
	Format7z      FormatType = "SevenZip"
	FormatRar     FormatType = "Rar"
	FormatUnknown FormatType = ""
)

// CompressionModeType is the user-facing knob persisted in ArchiveHistory.CompressionMode.
type CompressionModeType string

// CompressionMode is retained as a type alias for backwards compatibility.
type CompressionMode = CompressionModeType

const (
	ModeStandard CompressionModeType = constants.CompressionStandard
	ModeBest     CompressionModeType = constants.CompressionBest
	ModeFast     CompressionModeType = constants.CompressionFast
)

// CreateOptions bundles every knob `gitmap zip` exposes.
type CreateOptions struct {
	OutputPath string
	Sources    []string // absolute local paths
	Mode       CompressionMode
	Includes   []string // optional glob list
	Excludes   []string // optional glob list
}

// CreateResult is returned to the cmd layer for printing + history rows.
type CreateResult struct {
	OutputPath     string
	Format         Format
	EntriesWritten int
}

// ArchiveWriteParams encapsulates parameters for writing an archive.
type ArchiveWriteParams struct {
	Ctx    context.Context
	Path   string
	Format Format
	Mode   CompressionMode
	Files  []archives.FileInfo
}

// ExtractResult is the final output of an extract operation.
type ExtractResult struct {
	OutputDir       string
	Format          Format
	EntriesWritten  int
	UsedTempDir     bool
	FlattenedLayers int
}

// CompactExtractParams encapsulates parameters for completing a compact extraction.
type CompactExtractParams struct {
	Ctx         context.Context
	SrcArchive  string
	DestBaseDir string
	TempDir     string
	Result      ExtractResult
}

// ArchiveExtractParams encapsulates parameters for running an archive extraction.
type ArchiveExtractParams struct {
	Ctx       context.Context
	Extractor archives.Extractor
	Stream    io.Reader
	DestDir   string
}

// CopyDirEntryParams encapsulates parameters for copying a single directory entry.
type CopyDirEntryParams struct {
	Src   string
	Dst   string
	Path  string
	Entry fs.DirEntry
}

// Entry is a single file or directory record inside an archive.
type Entry struct {
	Path  string
	Size  int64
	IsDir bool
}

// ListExtractParams encapsulates parameters for extracting list entries.
type ListExtractParams struct {
	Ctx       context.Context
	Extractor archives.Extractor
	Stream    io.Reader
	Format    archives.Format
}

// SourceKindType classifies one entry on the command line.
type SourceKindType int

// SourceKind is retained as a type alias for backwards compatibility.
type SourceKind = SourceKindType

const (
	SourceLocal SourceKindType = iota
	SourceHTTP
	SourceGit
)

// ResolvedSource is the materialized form of one user-supplied input.
type ResolvedSource struct {
	Original   string
	Kind       SourceKindType
	LocalPath  string
	CleanupDir string
}

// Aria2cDownloadParams encapsulates parameters for aria2c download.
type Aria2cDownloadParams struct {
	Ctx    context.Context
	RawURL string
	Dir    string
	Name   string
}

// BoolResult wraps a boolean in a Result envelope.
type BoolResult = result.Result[bool]
