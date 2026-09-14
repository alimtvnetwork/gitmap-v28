package macro

import (
	"bytes"
	"io"
	"sync"
)

// SmartPaddedWriter buffers and pads process output with two spaces and smart newlines.
type SmartPaddedWriter struct {
	mu               sync.Mutex
	dest             io.Writer
	isFirstChunk     bool
	isLineStart      bool
	hasWritten       bool
	trailingNewlines int
}

// NewSmartPaddedWriter creates an output writer with 2-space padding and smart gaps.
func NewSmartPaddedWriter(dest io.Writer) *SmartPaddedWriter {
	return &SmartPaddedWriter{
		dest:         dest,
		isFirstChunk: true,
		isLineStart:  true,
		hasWritten:   false,
	}
}

// Write processes incoming chunks, applying leading gap and 2-space line indentation.
func (w *SmartPaddedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(p) == 0 {
		return 0, nil
	}

	w.handleLeadingGap(p)
	w.writePaddedBytes(p)
	w.hasWritten = true
	w.updateTrailingNewlines(p)

	return len(p), nil
}

func (w *SmartPaddedWriter) updateTrailingNewlines(p []byte) {
	for _, b := range p {
		if b == '\n' {
			w.trailingNewlines++
		} else if b != '\r' {
			w.trailingNewlines = 0
		}
	}
}

func (w *SmartPaddedWriter) handleLeadingGap(p []byte) {
	if !w.isFirstChunk {
		return
	}

	w.isFirstChunk = false
	if p[0] != '\n' && p[0] != '\r' {
		_, _ = w.dest.Write([]byte("\n"))
	}
}

func (w *SmartPaddedWriter) writePaddedBytes(p []byte) {
	lines := bytes.Split(p, []byte("\n"))
	for i, line := range lines {
		w.writeSinglePaddedLine(line, i, len(lines))
	}
}

func (w *SmartPaddedWriter) writeSinglePaddedLine(line []byte, idx, total int) {
	if len(line) > 0 {
		w.writeIndentedLineContent(line)
	}

	hasRemainingLines := idx < total-1
	if hasRemainingLines {
		_, _ = w.dest.Write([]byte("\n"))
		w.isLineStart = true
	}
}

func (w *SmartPaddedWriter) writeIndentedLineContent(line []byte) {
	if w.isLineStart {
		_, _ = w.dest.Write([]byte("  "))
		w.isLineStart = false
	}

	_, _ = w.dest.Write(line)
}

// Flush ensures the trailing gap is emitted without causing double blank lines.
func (w *SmartPaddedWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.hasWritten {
		return
	}

	if w.trailingNewlines >= 2 {
		return
	}

	if w.trailingNewlines == 1 {
		_, _ = w.dest.Write([]byte("\n"))
		return
	}

	_, _ = w.dest.Write([]byte("\n\n"))
}
