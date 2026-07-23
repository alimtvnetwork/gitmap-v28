package release

import (
	"bytes"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// filteredStderrWriter wraps an underlying writer (typically os.Stderr) and
// drops any newline-terminated line that contains a noise pattern (e.g. the
// CRLF/LF git warning that floods every release commit on Windows). Partial
// trailing data is buffered until the next newline arrives or Flush() is
// called, so we never split a real message in half.
type filteredStderrWriter struct {
	out      io.Writer
	patterns []string
	buf      bytes.Buffer
}

func newFilteredStderr(out io.Writer) *filteredStderrWriter {
	return &filteredStderrWriter{
		out:      out,
		patterns: constants.GitStderrNoisePatterns,
	}
}

// Write splits the input on '\n', filters each completed line, forwards the
// survivors, and buffers any trailing partial line.
func (w *filteredStderrWriter) Write(p []byte) (int, error) {
	w.buf.Write(p)
	for {
		idx := bytes.IndexByte(w.buf.Bytes(), '\n')
		if idx < 0 {
			return len(p), nil
		}

		line := w.buf.Next(idx + 1)
		if w.shouldDropLine(line) {
			continue
		}

		if _, err := w.out.Write(line); err != nil {
			return len(p), err
		}
	}
}

// Flush emits any buffered partial line (no trailing newline) so callers can
// drain the writer when the underlying process exits without a final '\n'.
func (w *filteredStderrWriter) Flush() error {
	if w.buf.Len() == 0 {
		return nil
	}

	line := w.buf.Bytes()
	if w.shouldDropLine(line) {
		w.buf.Reset()

		return nil
	}

	_, err := w.out.Write(line)
	w.buf.Reset()

	return err
}

func (w *filteredStderrWriter) shouldDropLine(line []byte) bool {
	s := string(line)
	for _, pat := range w.patterns {
		if strings.Contains(s, pat) {
			return true
		}
	}

	return false
}
