package archive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ExtractResult is what a compact-extract returns to the caller so it can
// be persisted into ArchiveHistory and printed to the user.
type ExtractResult struct {
	OutputDir       string
	Format          Format
	EntriesWritten  int
	UsedTempDir     bool
	FlattenedLayers int
}

// CompactExtract extracts srcArchive into a single normalized directory
// under destBaseDir, named after the archive's base name (sans extension).
//
// Algorithm: temp-dir-then-move. We always extract into a fresh temp dir
// inside destBaseDir, then walk it to find the "real root" — the first
// directory that either holds >1 entry OR holds at least one non-dir
// entry. That real root is then moved (or its contents merged) into
// `<destBaseDir>/<archiveBaseName>/`. This guarantees:
//
//  1. xap.zip → xap/xap/<files>  becomes  destBaseDir/xap/<files>
//     (any number of duplicate-name layers up to MaxCompactFlattenLayers
//     is collapsed; we do not require the inner names to match xap —
//     we just promote single-child directories until we hit content.)
//
//  2. xlt.zip → <files>          becomes  destBaseDir/xlt/<files>
//     (no flatten, just a wrap.)
//
//  3. mixed.zip → README + src/  becomes  destBaseDir/mixed/{README,src}
//     (no flatten, the temp dir contents move directly under the wrap.)
//
// The temp dir is always cleaned, even on failure mid-extract.
func CompactExtract(ctx context.Context, srcArchive, destBaseDir string) (ExtractResult, error) {
	res := ExtractResult{UsedTempDir: true}

	format, err := IdentifyArchive(ctx, srcArchive)
	if err != nil {
		return res, fmt.Errorf("identify %q: %w", srcArchive, err)
	}
	res.Format = format

	if err := os.MkdirAll(destBaseDir, constants.DirPermission); err != nil {
		return res, err
	}

	tempDir, err := os.MkdirTemp(destBaseDir, ".gitmap-uzc-*")
	if err != nil {
		return res, err
	}
	defer os.RemoveAll(tempDir)

	written, err := extractAllIntoDir(ctx, srcArchive, tempDir)
	if err != nil {
		return res, fmt.Errorf("extract: %w", err)
	}
	res.EntriesWritten = written

	finalDir := filepath.Join(destBaseDir, archiveBaseName(srcArchive))
	if err := os.RemoveAll(finalDir); err != nil {
		return res, err
	}

	flattened, err := promoteRealRoot(tempDir, finalDir)
	if err != nil {
		return res, err
	}
	res.FlattenedLayers = flattened
	res.OutputDir = finalDir

	return res, nil
}

// extractAllIntoDir streams every entry from srcArchive into destDir
// using mholt/archives. Returns the entry count written. Symlinks are
// rejected (security: a malicious archive could otherwise escape destDir
// even after path sanitation).
func extractAllIntoDir(ctx context.Context, srcArchive, destDir string) (int, error) {
	f, err := os.Open(srcArchive)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	format, stream, err := archives.Identify(ctx, filepath.Base(srcArchive), f)
	if err != nil {
		return 0, fmt.Errorf("identify: %w", err)
	}

	extractor, ok := format.(archives.Extractor)
	if !ok {
		return 0, fmt.Errorf("format %s is not extractable", format.Extension())
	}

	written := 0
	handler := func(_ context.Context, entry archives.FileInfo) error {
		clean := safeJoin(destDir, entry.NameInArchive)
		if clean == "" {
			return fmt.Errorf("rejecting entry with unsafe path: %q", entry.NameInArchive)
		}

		if entry.IsDir() {
			return os.MkdirAll(clean, constants.DirPermission)
		}

		if err := os.MkdirAll(filepath.Dir(clean), constants.DirPermission); err != nil {
			return err
		}

		return writeArchiveFile(entry, clean, &written)
	}

	if err := extractor.Extract(ctx, stream, handler); err != nil {
		return written, err
	}

	return written, nil
}

// writeArchiveFile streams a single entry into destPath and bumps written.
// Split out so extractAllIntoDir stays under gocyclo's 15-complexity cap.
func writeArchiveFile(entry archives.FileInfo, destPath string, written *int) error {
	if entry.LinkTarget != "" {
		// Symlinks are skipped on purpose — see CompactExtract docstring.
		return nil
	}

	src, err := entry.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, constants.FilePermission)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	*written++

	return nil
}

// safeJoin sanitizes an in-archive path against destDir to prevent
// path-traversal (G305). Returns "" when the cleaned path escapes destDir.
func safeJoin(destDir, name string) string {
	clean := filepath.Clean("/" + name) // anchor at root, strip "..", "."
	clean = strings.TrimPrefix(clean, string(filepath.Separator))
	full := filepath.Join(destDir, clean)
	abs, err := filepath.Abs(full)
	if err != nil {
		return ""
	}
	destAbs, err := filepath.Abs(destDir)
	if err != nil {
		return ""
	}
	if !strings.HasPrefix(abs+string(filepath.Separator), destAbs+string(filepath.Separator)) && abs != destAbs {
		return ""
	}

	return full
}

// promoteRealRoot finds the deepest single-child directory chain inside
// tempDir (capped at MaxCompactFlattenLayers) and moves its contents to
// finalDir, returning the number of layers collapsed.
//
// Edge cases:
//
//   - Empty archive  → finalDir is created empty.
//   - One file only  → finalDir holds that file (no flatten).
//   - Single dir at root, with multiple children → finalDir holds those
//     children (1 layer flattened — the wrapping dir merges into the
//     name we already chose).
//   - Two dirs at root → finalDir holds both (no flatten possible).
func promoteRealRoot(tempDir, finalDir string) (int, error) {
	root := tempDir
	flattened := 0

	for layer := 0; layer < constants.MaxCompactFlattenLayers; layer++ {
		entries, err := os.ReadDir(root)
		if err != nil {
			return flattened, err
		}
		if len(entries) != 1 || !entries[0].IsDir() {
			break
		}
		root = filepath.Join(root, entries[0].Name())
		flattened++
	}

	if err := os.MkdirAll(finalDir, constants.DirPermission); err != nil {
		return flattened, err
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return flattened, err
	}
	for _, entry := range entries {
		from := filepath.Join(root, entry.Name())
		to := filepath.Join(finalDir, entry.Name())
		if err := moveOrCopy(from, to); err != nil {
			return flattened, err
		}
	}

	return flattened, nil
}

// moveOrCopy renames src to dst, falling back to a recursive copy when
// the rename crosses a filesystem boundary (EXDEV) or when dst already
// exists as a directory we need to merge into.
func moveOrCopy(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}

	return copyFile(src, dst, info.Mode())
}

// copyDir performs a deep copy of src into dst.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, constants.DirPermission)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}

		return copyFile(path, target, info.Mode())
	})
}

// copyFile streams src → dst preserving mode bits.
func copyFile(src, dst string, mode fs.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), constants.DirPermission); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}

// archiveBaseName strips every recognized archive extension off path's
// base name so a "foo.tar.gz" yields "foo", not "foo.tar".
func archiveBaseName(path string) string {
	base := filepath.Base(path)
	for _, ext := range []string{
		".tar.gz", ".tar.bz2", ".tar.xz", ".tar.zst",
		".tgz", ".tbz2", ".txz", ".tzst",
		".zip", ".tar", ".gz", ".bz2", ".xz", ".zst",
		".7z", ".rar",
	} {
		if strings.HasSuffix(strings.ToLower(base), ext) {
			return base[:len(base)-len(ext)]
		}
	}

	return base
}

// ErrUnknownFormat is returned by CreateArchive when the output extension
// is not recognized. Surfaced as a typed error so the cmd layer can
// translate it into a friendly user message.
var ErrUnknownFormat = errors.New("unknown archive format")
