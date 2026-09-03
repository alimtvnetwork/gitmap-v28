package folder

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var seqRegex = regexp.MustCompile(`^(\d+)[-_](.*)$`)

// FileMeta holds rich analytical metadata for a scanned file.
type FileMeta struct {
	Path          string `json:"path"`
	Filename      string `json:"filename"`
	Directory     string `json:"directory"`
	Extension     string `json:"extension"`
	SizeBytes     int64  `json:"sizeBytes"`
	SizeFormatted string `json:"sizeFormatted"`
	LinesOfCode   int    `json:"linesOfCode"`
	IsBinary      bool   `json:"isBinary"`
	Sequence      int    `json:"sequence"`
}

// ExtractMetadata scans an individual file for size, LOC, encoding, and sequence.
func ExtractMetadata(absPath, relPath string) (*FileMeta, error) {
	fi, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	fname := filepath.Base(absPath)
	ext := filepath.Ext(fname)
	dir := filepath.Dir(relPath)
	dir = strings.ReplaceAll(dir, "\\", "/")
	relNorm := strings.ReplaceAll(relPath, "\\", "/")

	isBin, loc := inspectContent(absPath, size)
	seq := extractSequence(fname)

	return &FileMeta{
		Path:          relNorm,
		Filename:      fname,
		Directory:     dir,
		Extension:     ext,
		SizeBytes:     size,
		SizeFormatted: formatBytes(size),
		LinesOfCode:   loc,
		IsBinary:      isBin,
		Sequence:      seq,
	}, nil
}

func extractSequence(fname string) int {
	matches := seqRegex.FindStringSubmatch(fname)
	if len(matches) < 2 {
		return 0
	}

	seq, _ := strconv.Atoi(matches[1])

	return seq
}

func inspectContent(absPath string, size int64) (bool, int) {
	if size == 0 {
		return false, 0
	}

	f, err := os.Open(absPath)
	if err != nil {
		return false, 0
	}

	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	if bytes.Contains(buf[:n], []byte{0}) {
		return true, 0
	}

	_, _ = f.Seek(0, 0)
	loc := countLines(f)

	return false, loc
}

func countLines(r io.Reader) int {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		if err != nil {
			break
		}
	}

	return count
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
