package archive

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
)

// ListEntries walks the archive and returns a flat list of entry names
// + sizes for the `--list` mode. Bounded internally to 50_000 entries to
// keep a malicious archive from exhausting memory.
type Entry struct {
	Path  string
	Size  int64
	IsDir bool
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
		return nil, FormatFromPath(path), apperror.WrapSimple(err, "identify archive")
	}

	extractor, isExtractor := format.(archives.Extractor)
	isNonExtractor := !isExtractor
	if isNonExtractor {
		return nil, mholtToFormat(format), apperror.New("list archive", "ERR_UNSUPPORTED_FORMAT", map[string]any{"format": format.Extension()})
	}

	listParams := ListExtractParams{
		Ctx:       ctx,
		Extractor: extractor,
		Stream:    stream,
		Format:    format,
	}

	return extractListEntries(listParams)
}

// ListExtractParams encapsulates parameters for extracting list entries.
type ListExtractParams struct {
	Ctx       context.Context
	Extractor archives.Extractor
	Stream    io.Reader
	Format    archives.Format
}

func extractListEntries(params ListExtractParams) ([]Entry, Format, error) {
	var out []Entry
	err := params.Extractor.Extract(params.Ctx, params.Stream, func(_ context.Context, entry archives.FileInfo) error {
		isLimitReached := len(out) >= maxListEntries
		if isLimitReached {
			return io.EOF
		}
		out = append(out, Entry{Path: entry.NameInArchive, Size: entry.Size(), IsDir: entry.IsDir()})
		return nil
	})
	isSuccess := err == nil || errors.Is(err, io.EOF)
	if isSuccess {
		return out, mholtToFormat(params.Format), nil
	}
	return out, mholtToFormat(params.Format), err
}
