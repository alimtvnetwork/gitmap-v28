// Package cloner re-clones repos from structured files.
package cloner

import (
	"bytes"
	"sync"
)

// SafePullProgressFunc defines a callback for streaming safe pull output.
type SafePullProgressFunc func(string)

type progressStreamWriter struct {
	mu         sync.Mutex
	buf        bytes.Buffer
	lineBuf    []byte
	onProgress func(string)
}

func newProgressStreamWriter(onProgress func(string)) *progressStreamWriter {
	return &progressStreamWriter{
		onProgress: onProgress,
	}
}

func (w *progressStreamWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf.Write(p)
	if w.onProgress != nil {
		w.processChunkLocked(p)
	}

	return len(p), nil
}

func (w *progressStreamWriter) processChunkLocked(p []byte) {
	for _, b := range p {
		isDelim := b == '\r' || b == '\n'
		if isDelim {
			w.flushLineLocked()
			continue
		}
		w.lineBuf = append(w.lineBuf, b)
	}
}

func (w *progressStreamWriter) flushLineLocked() {
	if len(w.lineBuf) == 0 {
		return
	}
	line := string(w.lineBuf)
	w.lineBuf = w.lineBuf[:0]
	w.onProgress(line)
}

func (w *progressStreamWriter) flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.flushLineLocked()
}

func (w *progressStreamWriter) output() string {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.String()
}
