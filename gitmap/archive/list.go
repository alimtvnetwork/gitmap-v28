package archive

import (
	"context"
	"errors"
	"fmt"
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
		return nil, FormatFromPath(path), fmt.Errorf("archive identify: %w", err)
	}

	extractor, isExtractor := format.(archives.Extractor)
	if !isExtractor {
		return nil, mholtToFormat(format), fmt.Errorf("format %s is not extractable", format.Extension())
	}

	var out []Entry
	err = extractor.Extract(ctx, stream, func(_ context.Context, entry archives.FileInfo) error {
		if len(out) >= maxListEntries {
			return io.EOF // signal "stop walking" to mholt
		}
		out = append(out, Entry{
			Path:  entry.NameInArchive,
			Size:  entry.Size(),
			IsDir: entry.IsDir(),
		})

		return nil
	})
	if err == nil || errors.Is(err, io.EOF) {
		return out, mholtToFormat(format), nil
	}
	return out, mholtToFormat(format), err
}
