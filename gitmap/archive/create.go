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
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/result"

	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"context"
	"errors"
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

func validateCreateFormat(path string) (Format, error) {
	format := FormatFromPath(path)
	if format == FormatUnknown {
		return FormatUnknown, apperror.Wrap(ErrUnknownFormat, "validate format", map[string]any{"path": path})
	}
	if format == Format7z || format == FormatRar {
		return FormatUnknown, apperror.New("validate format", "ERR_READONLY_FORMAT", map[string]any{"format": format})
	}

	return format, nil
}

func prepareArchiveFiles(ctx context.Context, opts CreateOptions) ([]archives.FileInfo, error) {
	files, err := gatherFiles(ctx, opts.Sources)
	if err != nil {
		return nil, apperror.WrapSimple(err, "gather sources")
	}

	return filterFiles(files, opts.Includes, opts.Excludes), nil
}

func createOutputFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), constants.DirPermission); err != nil {
		return nil, err
	}

	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, constants.FilePermission)
}

func writeArchive(
	ctx context.Context,
	path string,
	format Format,
	mode CompressionMode,
	files []archives.FileInfo,
) error {
	out, err := createOutputFile(path)
	if err != nil {
		return err
	}
	defer out.Close()

	writer, err := buildArchiver(format, mode)
	if err != nil {
		return err
	}

	return writer.Archive(ctx, out, files)
}

// CreateArchive walks every source, applies include/exclude filters, and
// writes the archive to opts.OutputPath using the format derived from
// the output extension.

func CreateArchive(ctx context.Context, opts CreateOptions) (CreateResult, error) {
	res := CreateResult{OutputPath: opts.OutputPath}
	format, err := validateCreateFormat(opts.OutputPath)
	if err != nil {
		return res, err
	}
	res.Format = format

	files, err := prepareArchiveFiles(ctx, opts)
	if err != nil {
		return res, err
	}
	res.EntriesWritten = len(files)

	if err := writeArchive(ctx, opts.OutputPath, format, opts.Mode, files); err != nil {
		return res, err
	}

	return res, nil
}

func mapSourceEntry(src string) (string, error) {
	info, err := os.Stat(src)
	if err != nil {
		return "", err
	}
	base := filepath.Base(src)
	if info.IsDir() {
		return base + "/", nil
	}

	return base, nil
}

// gatherFiles converts each source root into mholt FileInfo entries.
// Roots are mapped to "<basename>/" so multi-source archives stay tidy
// (e.g. zip foo bar → archive contains foo/... + bar/...).

func gatherFiles(ctx context.Context, sources []string) ([]archives.FileInfo, error) {
	mapping := make(map[string]string, len(sources))
	for _, src := range sources {
		target, err := mapSourceEntry(src)
		if err != nil {
			return nil, err
		}
		mapping[src] = target
	}

	return archives.FilesFromDisk(ctx, nil, mapping)
}

func isEntryIncluded(name string, includes, excludes []string) result.Result[bool] {
	if matchAny(name, includes, !true).Data {
		return result.NewSuccess(false)
	}
	if matchAny(name, excludes, false).Data {
		return result.NewSuccess(false)
	}

	return result.NewSuccess(true)
}

// filterFiles applies include/exclude globs against NameInArchive.
func filterFiles(in []archives.FileInfo, includes, excludes []string) []archives.FileInfo {
	if len(includes) == 0 && len(excludes) == 0 {
		return in
	}
	out := in[:0]
	for _, f := range in {
		if isEntryIncluded(f.NameInArchive, includes, excludes).Data {
			out = append(out, f)
		}
	}

	return out
}

func matchPattern(pattern, name string) result.Result[bool] {
	if ok, err := filepath.Match(pattern, name); err == nil && ok {
		return result.NewSuccess(true)
	}
	if ok, err := filepath.Match(pattern, filepath.Base(name)); err == nil && ok {
		return result.NewSuccess(true)
	}

	return result.NewSuccess(false)
}

// matchAny returns true when name matches any pattern. emptyDefault is
// what we return when patterns is empty (true for includes = "match all
// when no filter set", false for excludes = "exclude nothing").
func matchAny(name string, patterns []string, emptyDefault bool) result.Result[bool] {
	if len(patterns) == 0 {
		return result.NewSuccess(emptyDefault)
	}
	for _, p := range patterns {
		if matchPattern(p, name).Data {
			return result.NewSuccess(true)
		}
	}

	return result.NewSuccess(false)
}

func buildCompressedTarArchiver(format Format, mode CompressionMode) (archives.Archiver, error) {
	switch format {
	case FormatTarGz:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Gz{CompressionLevel: gzipLevel(mode)}}, nil
	case FormatTarBz2:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Bz2{CompressionLevel: bz2Level(mode)}}, nil
	case FormatTarXz:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Xz{}}, nil
	case FormatTarZst:
		return archives.CompressedArchive{Archival: archives.Tar{}, Compression: archives.Zstd{}}, nil
	}

	return nil, errors.New("format not supported as archiver")
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
	case FormatTarGz, FormatTarBz2, FormatTarXz, FormatTarZst:
		return buildCompressedTarArchiver(format, mode)
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
