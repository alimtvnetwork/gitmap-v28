package termpad

import (
	"bytes"
	"io"
	"sync"
)

// SmartPaddingWriter buffers and prepends 2 spaces to unpadded lines.
type SmartPaddingWriter struct {
	mu          sync.Mutex
	dest        io.Writer
	isLineStart bool
	hasWritten  bool
}

// NewSmartPaddingWriter creates a thread-safe padded output writer.
func NewSmartPaddingWriter(dest io.Writer) *SmartPaddingWriter {
	return &SmartPaddingWriter{
		dest:        dest,
		isLineStart: true,
	}
}

// Write processes byte slices, indenting fresh lines with two spaces.
func (w *SmartPaddingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(p) == 0 {
		return 0, nil
	}
	w.writePaddedBytes(p)
	w.hasWritten = true

	return len(p), nil
}

func (w *SmartPaddingWriter) writePaddedBytes(p []byte) {
	lines := bytes.Split(p, []byte("\n"))
	for i, line := range lines {
		w.writeSingleLine(line, i, len(lines))
	}
}

func (w *SmartPaddingWriter) writeSingleLine(line []byte, idx, total int) {
	if len(line) > 0 {
		w.writeIndentedLine(line)
	}
	if idx < total-1 {
		_, _ = w.dest.Write([]byte("\n"))
		w.isLineStart = true
	}
}

func (w *SmartPaddingWriter) writeIndentedLine(line []byte) {
	if w.isLineStart && !bytes.HasPrefix(line, []byte("  ")) {
		_, _ = w.dest.Write([]byte("  "))
		w.isLineStart = false
	}
	_, _ = w.dest.Write(line)
}

// Flush ensures trailing output is properly terminated.
func (w *SmartPaddingWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.hasWritten {
		return
	}
	SetPaddingApplied(true)
}
