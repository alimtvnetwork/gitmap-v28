// Package archive owns local archive operations: format identification,
// compact-extraction (with up to MaxCompactFlattenLayers of duplicate-
// folder flattening), creation, and listing.
//
// The package is deliberately isolated from CLI concerns — the cmd layer
// resolves sources (URL fetch, git clone, cwd auto-pick) and then feeds
// concrete local paths into the functions defined here. That separation
// is what lets the same archive engine power Slice 3 of the downloader
// feature later without an import cycle.
//
// Format coverage (via github.com/mholt/archives):
//
//	read  + write : zip, tar, tar.gz, tar.bz2, tar.xz, tar.zst, gz, bz2, xz, zst
//	read-only     : 7z, rar
//
// 7z/rar writing is rejected with a clear error in CreateArchive so the
// CLI can surface "use zip/tar.* for outputs" without crashing.
package archive

import (
	"context"
	"os"
	"path/filepath"
	"strings"

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

var doubleExtMap = map[string]FormatType{
	".tar.gz":  FormatTarGz,
	".tgz":     FormatTarGz,
	".tar.bz2": FormatTarBz2,
	".tbz2":    FormatTarBz2,
	".tar.xz":  FormatTarXz,
	".txz":     FormatTarXz,
	".tar.zst": FormatTarZst,
	".tzst":    FormatTarZst,
}

var singleExtMap = map[string]FormatType{
	".zip": FormatZip,
	".tar": FormatTar,
	".gz":  FormatGz,
	".bz2": FormatBz2,
	".xz":  FormatXz,
	".zst": FormatZst,
	".7z":  Format7z,
	".rar": FormatRar,
}

var formatExtMap = map[FormatType]string{
	FormatZip:    ".zip",
	FormatTar:    ".tar",
	FormatTarGz:  ".tar.gz",
	FormatTarBz2: ".tar.bz2",
	FormatTarXz:  ".tar.xz",
	FormatTarZst: ".tar.zst",
	FormatGz:     ".gz",
	FormatBz2:    ".bz2",
	FormatXz:     ".xz",
	FormatZst:    ".zst",
	Format7z:     ".7z",
	FormatRar:    ".rar",
}

func matchExtMap(lower string, m map[string]FormatType) FormatType {
	for ext, f := range m {
		if strings.HasSuffix(lower, ext) {
			return f
		}
	}

	return FormatUnknown
}

// FormatFromPath inspects a file name and returns the matching FormatType,
// or FormatUnknown when nothing matches. Multi-extension forms
// (".tar.gz", ".tar.bz2", ".tar.xz", ".tar.zst") are checked first so a
// plain ".gz" never wins over ".tar.gz".
func FormatFromPath(p string) FormatType {
	lower := strings.ToLower(p)
	if f := matchExtMap(lower, doubleExtMap); f != FormatUnknown {
		return f
	}
	if f := matchExtMap(lower, singleExtMap); f != FormatUnknown {
		return f
	}

	return FormatUnknown
}

// Extension returns the canonical extension (with leading dot) the
// CreateArchive path uses to construct mholt/archives Format objects.
func (f FormatType) Extension() string {
	return formatExtMap[f]
}

// IdentifyArchive opens the file and asks mholt/archives to sniff the
// magic bytes. Used as the authoritative format check after extension-
// based guesses, since a misnamed file (foo.zip that is really a tarball)
// would otherwise produce a misleading ArchiveHistory.ArchiveFormat row.

func IdentifyArchive(ctx context.Context, path string) (FormatType, error) {
	f, err := os.Open(path)
	if err != nil {
		return FormatUnknown, err
	}
	defer f.Close()

	format, _, err := archives.Identify(ctx, filepath.Base(path), f)
	if err != nil {
		return FormatFromPath(path), nil
	}

	return mholtToFormat(format), nil
}

// mholtToFormat maps mholt/archives Format objects back to our FormatType
// enum. The library exposes one struct per format, so a type switch is
// the cleanest way; comparing extensions would round-trip through strings.
func mholtToFormat(f archives.Format) FormatType {
	if f == nil {
		return FormatUnknown
	}
	if res := mholtArchiveType(f); res != FormatUnknown {
		return res
	}

	return mholtCompressType(f)
}

func mholtArchiveType(f archives.Format) FormatType {
	switch f.(type) {
	case archives.Zip:
		return FormatZip
	case archives.Tar:
		return FormatTar
	case archives.SevenZip:
		return Format7z
	case archives.Rar:
		return FormatRar
	case archives.CompressedArchive:
		return FormatFromPath(f.Extension())
	}
	return FormatUnknown
}

func mholtCompressType(f archives.Format) FormatType {
	switch f.(type) {
	case archives.Gz:
		return FormatGz
	case archives.Bz2:
		return FormatBz2
	case archives.Xz:
		return FormatXz
	case archives.Zstd:
		return FormatZst
	}
	return FormatUnknown
}
