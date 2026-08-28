package archive

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

	"context"
	"errors"
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
// `<destBaseDir>/<archiveBaseName>/`.
func CompactExtract(ctx context.Context, srcArchive, destBaseDir string) (ExtractResult, error) {
	res := ExtractResult{UsedTempDir: true}
	format, err := prepareExtractDest(ctx, srcArchive, destBaseDir)
	if err != nil {
		return res, err
	}
	res.Format = format

	tempDir, err := os.MkdirTemp(destBaseDir, ".gitmap-uzc-*")
	if err != nil {
		return res, err
	}
	defer os.RemoveAll(tempDir)

	return completeCompactExtract(ctx, srcArchive, destBaseDir, tempDir, res)
}

func prepareExtractDest(ctx context.Context, srcArchive, destBaseDir string) (Format, error) {
	format, err := IdentifyArchive(ctx, srcArchive)
	if err != nil {
		return format, apperror.Wrap(err, "identify archive", map[string]any{"src": srcArchive})
	}
	if err := os.MkdirAll(destBaseDir, constants.DirPermission); err != nil {
		return format, err
	}
	return format, nil
}

func completeCompactExtract(ctx context.Context, srcArchive, destBaseDir, tempDir string, res ExtractResult) (ExtractResult, error) {
	written, err := extractAllIntoDir(ctx, srcArchive, tempDir)
	if err != nil {
		return res, apperror.WrapSimple(err, "extract")
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
		return 0, apperror.WrapSimple(err, "identify")
	}

	extractor, ok := format.(archives.Extractor)
	isMissingExtractor := !ok
	if isMissingExtractor == true {
		return 0, apperror.New("extract", "ERR_UNSUPPORTED_FORMAT", map[string]any{"format": format.Extension()})
	}

	return runArchiveExtraction(ctx, extractor, stream, destDir)
}

func runArchiveExtraction(ctx context.Context, extractor archives.Extractor, stream io.Reader, destDir string) (int, error) {
	written := 0
	handler := func(_ context.Context, entry archives.FileInfo) error {
		return extractArchiveEntry(destDir, entry, &written)
	}

	if err := extractor.Extract(ctx, stream, handler); err != nil {
		return written, err
	}
	return written, nil
}

func extractArchiveEntry(destDir string, entry archives.FileInfo, written *int) error {
	clean := safeJoin(destDir, entry.NameInArchive)
	isEmptyClean := clean == ""
	if isEmptyClean == true {
		return apperror.New("extract entry", "ERR_UNSAFE_PATH", map[string]any{"path": entry.NameInArchive})
	}

	isDir := entry.IsDir()
	if isDir == true {
		return os.MkdirAll(clean, constants.DirPermission)
	}

	if err := os.MkdirAll(filepath.Dir(clean), constants.DirPermission); err != nil {
		return err
	}
	return writeArchiveFile(entry, clean, written)
}

// writeArchiveFile streams a single entry into destPath and bumps written.
// Split out so extractAllIntoDir stays under gocyclo's 15-complexity cap.
func writeArchiveFile(entry archives.FileInfo, destPath string, written *int) error {
	hasLinkTarget := entry.LinkTarget != ""
	if hasLinkTarget == true {
		// Symlinks are skipped on purpose — see CompactExtract docstring.
		return nil
	}

	src, err := entry.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	return copyEntryToFile(src, destPath, written)
}

func copyEntryToFile(src io.Reader, destPath string, written *int) error {
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
	isOutsideDest := !strings.HasPrefix(abs+string(filepath.Separator), destAbs+string(filepath.Separator))
	isSame := abs == destAbs
	isEscaped := isOutsideDest == true && isSame == false
	if isEscaped == true {
		return ""
	}

	return full
}

// promoteRealRoot finds the deepest single-child directory chain inside
// tempDir (capped at MaxCompactFlattenLayers) and moves its contents to
// finalDir, returning the number of layers collapsed.
func promoteRealRoot(tempDir, finalDir string) (int, error) {
	root, flattened, err := findDeepestRoot(tempDir)
	if err != nil {
		return flattened, err
	}

	if err := os.MkdirAll(finalDir, constants.DirPermission); err != nil {
		return flattened, err
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return flattened, err
	}
	return flattened, moveEntries(root, finalDir, entries)
}

func findDeepestRoot(tempDir string) (string, int, error) {
	root := tempDir
	flattened := 0
	for layer := 0; layer < constants.MaxCompactFlattenLayers; layer++ {
		entries, err := os.ReadDir(root)
		if err != nil {
			return root, flattened, err
		}
		hasSingleEntry := len(entries) == 1
		if hasSingleEntry == false {
			break
		}
		isDir := entries[0].IsDir()
		isNonDir := !isDir
		if isNonDir == true {
			break
		}
		root = filepath.Join(root, entries[0].Name())
		flattened++
	}
	return root, flattened, nil
}

func moveEntries(root, finalDir string, entries []fs.DirEntry) error {
	for _, entry := range entries {
		from := filepath.Join(root, entry.Name())
		to := filepath.Join(finalDir, entry.Name())
		if err := moveOrCopy(from, to); err != nil {
			return err
		}
	}
	return nil
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
	isDir := info.IsDir()
	if isDir == true {
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
		return copyDirEntry(src, dst, path, d)
	})
}

func copyDirEntry(src, dst, path string, d fs.DirEntry) error {
	rel, err := filepath.Rel(src, path)
	if err != nil {
		return err
	}
	target := filepath.Join(dst, rel)
	isDir := d.IsDir()
	if isDir == true {
		return os.MkdirAll(target, constants.DirPermission)
	}
	info, err := d.Info()
	if err != nil {
		return err
	}

	return copyFile(path, target, info.Mode())
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
	return streamToFile(in, dst, mode)
}

func streamToFile(in io.Reader, dst string, mode fs.FileMode) error {
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
		hasExt := strings.HasSuffix(strings.ToLower(base), ext)
		if hasExt == true {
			return base[:len(base)-len(ext)]
		}
	}

	return base
}

// ErrUnknownFormat is returned by CreateArchive when the output extension
// is not recognized. Surfaced as a typed error so the cmd layer can
// translate it into a friendly user message.
var ErrUnknownFormat = errors.New("unknown archive format")
