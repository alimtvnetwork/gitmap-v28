package applogger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// ConsoleSink outputs log entries to an io.Writer (e.g. os.Stdout).
type ConsoleSink struct {
	mu      sync.Mutex
	writer  io.Writer
	useJSON bool
}

// NewConsoleSink constructs a ConsoleSink.
func NewConsoleSink(writer io.Writer, useJSON bool) *ConsoleSink {
	if writer == nil {
		writer = os.Stdout
	}

	return &ConsoleSink{
		writer:  writer,
		useJSON: useJSON,
	}
}

// formatTextEntry returns a formatted terminal line.
func (cs *ConsoleSink) formatTextEntry(e LogEntry) string {
	var fieldStr string
	if len(e.Fields) > 0 {
		fieldStr = fmt.Sprintf(" %s", e.Fields.Format())
	}

	return fmt.Sprintf("[%s] [%s] %s%s\n", e.Timestamp, e.Level.Name(), e.Message, fieldStr)
}

// writeJSON writes entry as json string.
func (cs *ConsoleSink) writeJSON(e LogEntry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(cs.writer, string(b))

	return err
}

// WriteEntry writes the log entry to the console.
func (cs *ConsoleSink) WriteEntry(e LogEntry) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.useJSON {
		return cs.writeJSON(e)
	}

	_, err := fmt.Fprint(cs.writer, cs.formatTextEntry(e))

	return err
}

// Sync flushes the writer if supported.
func (cs *ConsoleSink) Sync() error { return nil }

// Close closes the console sink.
func (cs *ConsoleSink) Close() error { return nil }
