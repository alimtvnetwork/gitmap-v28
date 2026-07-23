// Package archive — write side. Builds zip / tar / tar.* archives from
// a heterogeneous list of local source paths using mholt/archives.
//
// Compression mode → library knobs:
//
//	Best     → DEFLATE max  / gzip 9 / bz2 9
//	Standard → DEFLATE def  / gzip default / bz2 default
//	Fast     → DEFLATE 1    / gzip 1 / bz2 1
//
// Filtering: optional include / exclude glob lists run against the
// in-archive name (NameInArchive). An entry survives when either no
// includes are set OR it matches at least one include, AND it does NOT
// match any exclude.
package archive

import (
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mholt/archives"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// CompressionMode is the user-facing knob persisted in
// ArchiveHistory.CompressionMode.
type CompressionMode string

const (
	ModeStandard CompressionMode = constants.CompressionStandard
	ModeBest     CompressionMode = constants.CompressionBest
	ModeFast     CompressionMode = constants.CompressionFast
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

// CreateArchive walks every source, applies include/exclude filters, and
// writes the archive to opts.OutputPath using the format derived from
// the output extension.
func CreateArchive(ctx context.Context, opts CreateOptions) (CreateResult, error) {
	res := CreateResult{OutputPath: opts.OutputPath}

	format := FormatFromPath(opts.OutputPath)
	if format == FormatUnknown {
		return res, fmt.Errorf("%w: %q", ErrUnknownFormat, opts.OutputPath)
	}
	if format == Format7z || format == FormatRar {
		return res, fmt.Errorf("%s archives are read-only in this build (use zip or tar.*)", format)
	}
	res.Format = format

	files, err := gatherFiles(ctx, opts.Sources)
	if err != nil {
		return res, fmt.Errorf("gather sources: %w", err)
	}
	files = filterFiles(files, opts.Includes, opts.Excludes)
	res.EntriesWritten = len(files)

	if err := os.MkdirAll(filepath.Dir(opts.OutputPath), constants.DirPermission); err != nil {
		return res, err
	}
	out, err := os.OpenFile(opts.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, constants.FilePermission)
	if err != nil {
		return res, err
	}
	defer out.Close()

	writer, err := buildArchiver(format, opts.Mode)
	if err != nil {
		return res, err
	}

	if err := writer.Archive(ctx, out, files); err != nil {
		return res, err
	}

	return res, nil
}

// gatherFiles converts each source root into mholt FileInfo entries.
// Roots are mapped to "<basename>/" so multi-source archives stay tidy
// (e.g. zip foo bar → archive contains foo/... + bar/...).
func gatherFiles(ctx context.Context, sources []string) ([]archives.FileInfo, error) {
	mapping := make(map[string]string, len(sources))
	for _, src := range sources {
		info, err := os.Stat(src)
		if err != nil {
			return nil, err
		}
		base := filepath.Base(src)
		if info.IsDir() {
			mapping[src] = base + "/"
		} else {
			mapping[src] = base
		}
	}

	return archives.FilesFromDisk(ctx, nil, mapping)
}

// filterFiles applies include/exclude globs against NameInArchive.
func filterFiles(in []archives.FileInfo, includes, excludes []string) []archives.FileInfo {
	if len(includes) == 0 && len(excludes) == 0 {
		return in
	}
	out := in[:0]
	for _, f := range in {
		if !matchAny(f.NameInArchive, includes, true) {
			continue
		}
		if matchAny(f.NameInArchive, excludes, false) {
			continue
		}
		out = append(out, f)
	}

	return out
}

// matchAny returns true when name matches any pattern. emptyDefault is
// what we return when patterns is empty (true for includes = "match all
// when no filter set", false for excludes = "exclude nothing").
func matchAny(name string, patterns []string, emptyDefault bool) bool {
	if len(patterns) == 0 {
		return emptyDefault
	}
	for _, p := range patterns {
		ok, err := filepath.Match(p, name)
		if err == nil && ok {
			return true
		}
		// Also try matching just the basename so users can write *.go
		// without worrying about the in-archive prefix.
		ok, err = filepath.Match(p, filepath.Base(name))
		if err == nil && ok {
			return true
		}
	}

	return false
}

// buildArchiver returns the mholt writer pre-tuned for the requested
// compression mode. Tar (uncompressed) and Zip get bespoke handling;
// the tar.* family flows through CompressedArchive.
func buildArchiver(format Format, mode CompressionMode) (archives.Archiver, error) {
	switch format {
	case FormatZip:
		return archives.Zip{Compression: zip.Deflate, SelectiveCompression: true}, nil
	case FormatTar:
		return archives.Tar{}, nil
	case FormatTarGz:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Gz{CompressionLevel: gzipLevel(mode)}}, nil
	case FormatTarBz2:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Bz2{CompressionLevel: bz2Level(mode)}}, nil
	case FormatTarXz:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Xz{}}, nil
	case FormatTarZst:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Zstd{}}, nil
	case FormatGz, FormatBz2, FormatXz, FormatZst, Format7z, FormatRar, FormatUnknown:
		return nil, errors.New("format not supported as archiver")
	}

	return nil, errors.New("format not supported as archiver")
}

// gzipLevel maps mode → compress/gzip level.
func gzipLevel(mode CompressionMode) int {
	switch mode {
	case ModeFast:
		return gzip.BestSpeed
	case ModeBest:
		return gzip.BestCompression
	case ModeStandard:
		return gzip.DefaultCompression
	}

	return gzip.DefaultCompression
}

// bz2Level maps mode → klauspost bzip2 level (1..9).
func bz2Level(mode CompressionMode) int {
	switch mode {
	case ModeFast:
		return 1
	case ModeBest:
		return 9
	case ModeStandard:
		return 6
	}

	return 6
}

// flateLevel exists for callers that want to log the resolved deflate
// level alongside the zip method (kept here so install/test sites can
// assert against a single source of truth).
func flateLevel(mode CompressionMode) int {
	switch mode {
	case ModeFast:
		return flate.BestSpeed
	case ModeBest:
		return flate.BestCompression
	case ModeStandard:
		return flate.DefaultCompression
	}

	return flate.DefaultCompression
}

// FlateLevelForMode is the exported helper for the cmd layer's --list
// banner so users can see what they signed up for.
func FlateLevelForMode(mode CompressionMode) int { return flateLevel(mode) }
