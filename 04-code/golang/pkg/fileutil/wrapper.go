package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// FileWrapper provides structured file operations with error object returns.
type FileWrapper struct {
	mu            sync.RWMutex
	name          string
	defaultAction FileActionType
	defaultMode   FileModeType
}

// NewFileWrapper creates a new FileWrapper with the provided identifier.
func NewFileWrapper(name string) *FileWrapper {
	if name == "" {
		name = "file-util-wrapper"
	}

	return NewFileWrapperWithOptions(name, FileActionTypeReadOnly, FileModeTypeDefaultFile)
}

// NewFileWrapperWithOptions creates a FileWrapper with initial defaults.
func NewFileWrapperWithOptions(name string, action FileActionType, mode FileModeType) *FileWrapper {
	return &FileWrapper{
		name:          name,
		defaultAction: action,
		defaultMode:   mode,
	}
}

// NewDefault creates a FileWrapper with standard production defaults.
func NewDefault() *FileWrapper {
	return NewFileWrapper("file-util-default")
}

// WrapFailure wraps any error into an AppError with an explicit errorId.
func (w *FileWrapper) WrapFailure(cause error, errorId string, message string) *appfault.AppError {
	if cause == nil {
		return nil
	}

	return appfault.WrapFailure(errtype.IO, errorId, cause, message)
}

// WrapReaderFailure wraps a read failure with an errorId and path context.
func (w *FileWrapper) WrapReaderFailure(cause error, errorId string, filePath string) *appfault.AppError {
	if cause == nil {
		return nil
	}
	msg := fmt.Sprintf("[%s] read failure on path: %s", w.name, filePath)

	return appfault.WrapReaderFailure(errorId, cause, msg)
}

// WrapWriterFailure wraps a write failure with an errorId and path context.
func (w *FileWrapper) WrapWriterFailure(cause error, errorId string, filePath string) *appfault.AppError {
	if cause == nil {
		return nil
	}
	msg := fmt.Sprintf("[%s] write failure on path: %s", w.name, filePath)

	return appfault.WrapWriterFailure(errorId, cause, msg)
}

// Read reads file contents based on action and mode parameters.
func (w *FileWrapper) Read(filePath string, action FileActionType, mode FileModeType) ([]byte, *appfault.AppError) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	flags := action.ToOSFlags()
	f, err := os.OpenFile(filePath, flags, mode.ToFileMode())
	if err != nil {
		return nil, w.WrapReaderFailure(err, "ERR_FILE_OPEN", filePath)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, w.WrapReaderFailure(err, "ERR_FILE_READ", filePath)
	}

	return data, nil
}

// ReadString reads file contents as string.
func (w *FileWrapper) ReadString(filePath string, action FileActionType, mode FileModeType) (string, *appfault.AppError) {
	data, appErr := w.Read(filePath, action, mode)
	if appErr != nil {
		return "", appErr
	}

	return string(data), nil
}

// Write writes data to filePath applying the specified action and mode.
func (w *FileWrapper) Write(filePath string, data []byte, action FileActionType, mode FileModeType) *appfault.AppError {
	w.mu.Lock()
	defer w.mu.Unlock()

	if appErr := ensureParentDirectory(filePath, w); appErr != nil {
		return appErr
	}

	flags := action.ToOSFlags()
	f, err := os.OpenFile(filePath, flags, mode.ToFileMode())
	if err != nil {
		return w.WrapWriterFailure(err, "ERR_FILE_WRITE_OPEN", filePath)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return w.WrapWriterFailure(err, "ERR_FILE_WRITE", filePath)
	}

	return nil
}

// WriteString writes a string content to filePath.
func (w *FileWrapper) WriteString(filePath string, content string, action FileActionType, mode FileModeType) *appfault.AppError {
	return w.Write(filePath, []byte(content), action, mode)
}

// Append appends data to filePath using FileModeTypeDefaultFile.
func (w *FileWrapper) Append(filePath string, data []byte, mode FileModeType) *appfault.AppError {
	return w.Write(filePath, data, FileActionTypeAppend, mode)
}

// Delete removes the specified file path.
func (w *FileWrapper) Delete(filePath string) *appfault.AppError {
	w.mu.Lock()
	defer w.mu.Unlock()

	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return w.WrapFailure(err, "ERR_FILE_DELETE", fmt.Sprintf("failed to delete file: %s", filePath))
	}

	return nil
}

// Create creates an empty file with given mode if it does not already exist.
func (w *FileWrapper) Create(filePath string, mode FileModeType) *appfault.AppError {
	return w.Write(filePath, []byte{}, FileActionTypeCreate, mode)
}

// Execute performs dynamic file I/O dispatch based on FileActionType.
func (w *FileWrapper) Execute(filePath string, action FileActionType, mode FileModeType, data []byte) ([]byte, *appfault.AppError) {
	if action.IsDelete() {
		return nil, w.Delete(filePath)
	}
	if action.IsRead() && len(data) == 0 {
		return w.Read(filePath, action, mode)
	}
	appErr := w.Write(filePath, data, action, mode)

	return data, appErr
}

// ensureParentDirectory creates parent directories if missing.
func ensureParentDirectory(filePath string, w *FileWrapper) *appfault.AppError {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return w.WrapWriterFailure(err, "ERR_DIR_CREATE", dir)
	}

	return nil
}
