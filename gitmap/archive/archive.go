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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
)

// Format is a string tag persisted in ArchiveHistory.ArchiveFormat. It
// reads cleanly in PascalCase logs ("Zip", "TarGz") yet round-trips
// through the canonical extension via FormatFromExt / Format.Extension.
type Format string

const (
	FormatZip     Format = "Zip"
	FormatTar     Format = "Tar"
	FormatTarGz   Format = "TarGz"
	FormatTarBz2  Format = "TarBz2"
	FormatTarXz   Format = "TarXz"
	FormatTarZst  Format = "TarZst"
	FormatGz      Format = "Gz"
	FormatBz2     Format = "Bz2"
	FormatXz      Format = "Xz"
	FormatZst     Format = "Zst"
	Format7z      Format = "SevenZip"
	FormatRar     Format = "Rar"
	FormatUnknown Format = ""
)

// FormatFromPath inspects a file name and returns the matching Format,
// or FormatUnknown when nothing matches. Multi-extension forms
// (".tar.gz", ".tar.bz2", ".tar.xz", ".tar.zst") are checked first so a
// plain ".gz" never wins over ".tar.gz".
func FormatFromPath(p string) Format {
	lower := strings.ToLower(p)
	doubles := map[string]Format{
		".tar.gz":  FormatTarGz,
		".tgz":     FormatTarGz,
		".tar.bz2": FormatTarBz2,
		".tbz2":    FormatTarBz2,
		".tar.xz":  FormatTarXz,
		".txz":     FormatTarXz,
		".tar.zst": FormatTarZst,
		".tzst":    FormatTarZst,
	}
	for ext, f := range doubles {
		if strings.HasSuffix(lower, ext) {
			return f
		}
	}
	singles := map[string]Format{
		".zip": FormatZip,
		".tar": FormatTar,
		".gz":  FormatGz,
		".bz2": FormatBz2,
		".xz":  FormatXz,
		".zst": FormatZst,
		".7z":  Format7z,
		".rar": FormatRar,
	}
	for ext, f := range singles {
		if strings.HasSuffix(lower, ext) {
			return f
		}
	}

	return FormatUnknown
}

// Extension returns the canonical extension (with leading dot) the
// CreateArchive path uses to construct mholt/archives Format objects.
func (f Format) Extension() string {
	switch f {
	case FormatZip:
		return ".zip"
	case FormatTar:
		return ".tar"
	case FormatTarGz:
		return ".tar.gz"
	case FormatTarBz2:
		return ".tar.bz2"
	case FormatTarXz:
		return ".tar.xz"
	case FormatTarZst:
		return ".tar.zst"
	case FormatGz:
		return ".gz"
	case FormatBz2:
		return ".bz2"
	case FormatXz:
		return ".xz"
	case FormatZst:
		return ".zst"
	case Format7z:
		return ".7z"
	case FormatRar:
		return ".rar"
	case FormatUnknown:
		return ""
	}

	return ""
}

// IdentifyArchive opens the file and asks mholt/archives to sniff the
// magic bytes. Used as the authoritative format check after extension-
// based guesses, since a misnamed file (foo.zip that is really a tarball)
// would otherwise produce a misleading ArchiveHistory.ArchiveFormat row.
func IdentifyArchive(ctx context.Context, path string) (Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return FormatUnknown, err
	}
	defer f.Close()

	format, _, err := archives.Identify(ctx, filepath.Base(path), f)
	if err != nil {
		// mholt returns ErrNoMatch when nothing identifies — fall back to
		// the extension hint so a tar without magic bytes still works.
		return FormatFromPath(path), nil
	}

	return mholtToFormat(format), nil
}

// mholtToFormat maps mholt/archives Format objects back to our Format
// enum. The library exposes one struct per format, so a type switch is
// the cleanest way; comparing extensions would round-trip through strings.
func mholtToFormat(f archives.Format) Format {
	if f == nil {
		return FormatUnknown
	}
	switch f.(type) {
	case archives.Zip:
		return FormatZip
	case archives.Tar:
		return FormatTar
	case archives.CompressedArchive:
		// Composite (e.g. tar.gz) — fall back to extension parsing on the
		// declared name, which is what the caller passed in.
		return FormatFromPath(f.Extension())
	case archives.Gz:
		return FormatGz
	case archives.Bz2:
		return FormatBz2
	case archives.Xz:
		return FormatXz
	case archives.Zstd:
		return FormatZst
	case archives.SevenZip:
		return Format7z
	case archives.Rar:
		return FormatRar
	}

	return FormatUnknown
}

// ListEntries walks the archive and returns a flat list of entry names
// + sizes for the `--list` mode. Bounded internally to 50_000 entries to
// keep a malicious archive from exhausting memory.
type Entry struct {
	Path string
	Size int64
	Dir  bool
}

const maxListEntries = 50_000

// ListEntries returns up to maxListEntries entries plus the detected
// format. Used by `gitmap uzc --list <archive>`.
func ListEntries(ctx context.Context, path string) ([]Entry, Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, FormatUnknown, err
	}
	defer f.Close()

	format, stream, err := archives.Identify(ctx, filepath.Base(path), f)
	if err != nil {
		return nil, FormatFromPath(path), fmt.Errorf("archive identify: %w", err)
	}

	extractor, ok := format.(archives.Extractor)
	if !ok {
		return nil, mholtToFormat(format), fmt.Errorf("format %s is not extractable", format.Extension())
	}

	var out []Entry
	err = extractor.Extract(ctx, stream, func(_ context.Context, entry archives.FileInfo) error {
		if len(out) >= maxListEntries {
			return io.EOF // signal "stop walking" to mholt
		}
		out = append(out, Entry{
			Path: entry.NameInArchive,
			Size: entry.Size(),
			Dir:  entry.IsDir(),
		})

		return nil
	})
	if err != nil && !errors.Is(err, io.EOF) {
		return out, mholtToFormat(format), err
	}

	return out, mholtToFormat(format), nil
}
