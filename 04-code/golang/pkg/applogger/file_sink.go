package applogger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileSink writes log entries to a file path.
type FileSink struct {
	mu       sync.Mutex
	filePath string
	file     *os.File
}

// openLogFile creates directories and opens file.
func openLogFile(filePath string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

// NewFileSink creates and opens a log file sink.
func NewFileSink(filePath string) (*FileSink, error) {
	f, err := openLogFile(filePath)
	if err != nil {
		return nil, err
	}

	return &FileSink{filePath: filePath, file: f}, nil
}

// WriteEntry serializes and writes entry to file.
func (fs *FileSink) WriteEntry(e LogEntry) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(fs.file, string(b))

	return err
}

// Sync flushes the file buffer to disk.
func (fs *FileSink) Sync() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.file != nil {
		return fs.file.Sync()
	}

	return nil
}

// Close closes the file descriptor.
func (fs *FileSink) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.file != nil {
		err := fs.file.Close()
		fs.file = nil

		return err
	}

	return nil
}
