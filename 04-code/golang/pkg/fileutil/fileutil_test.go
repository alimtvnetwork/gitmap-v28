package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileActionTypeEnums(t *testing.T) {
	if !FileActionTypeReadOnly.IsRead() {
		t.Errorf("expected ReadOnly to be read")
	}
	if !FileActionTypeWrite.IsWrite() {
		t.Errorf("expected Write to be write")
	}
	if !FileActionTypeAppend.IsAppend() {
		t.Errorf("expected Append to be append")
	}
	if !FileActionTypeCreate.IsCreate() {
		t.Errorf("expected Create to be create")
	}
	if !FileActionTypeDelete.IsDelete() {
		t.Errorf("expected Delete to be delete")
	}
	if FileActionTypeReadWrite.ToOSFlags()&os.O_RDWR == 0 {
		t.Errorf("expected ReadWrite flags to include O_RDWR")
	}
}

func TestFileModeTypeEnums(t *testing.T) {
	modeFile := FileModeTypeDefaultFile
	if modeFile.ToFileMode() != 0644 {
		t.Errorf("expected 0644, got %v", modeFile.ToFileMode())
	}
	if !modeFile.IsWritable() {
		t.Errorf("expected 0644 to be writable")
	}
	if modeFile.IsExecutable() {
		t.Errorf("expected 0644 not to be executable")
	}

	modeDir := FileModeTypeDefaultDir
	if !modeDir.IsExecutable() {
		t.Errorf("expected 0755 to be executable")
	}

	modeRo := FileModeTypeReadOnly
	if !modeRo.IsReadOnly() {
		t.Errorf("expected 0444 to be read only")
	}
}

func TestFileWrapper_ReadWriteCycle(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sub", "test.txt")

	wrapper := NewFileWrapper("test-wrapper")

	// 1. Write
	content := "hello world"
	appErr := wrapper.WriteString(filePath, content, FileActionTypeWrite, FileModeTypeDefaultFile)
	if appErr != nil {
		t.Fatalf("write failed: %v", appErr)
	}

	// 2. ReadString
	readContent, appErr := wrapper.ReadString(filePath, FileActionTypeReadOnly, FileModeTypeDefaultFile)
	if appErr != nil {
		t.Fatalf("read failed: %v", appErr)
	}
	if readContent != content {
		t.Errorf("got %q, want %q", readContent, content)
	}

	// 3. Append
	appendContent := " appended"
	appErr = wrapper.Append(filePath, []byte(appendContent), FileModeTypeDefaultFile)
	if appErr != nil {
		t.Fatalf("append failed: %v", appErr)
	}

	// 4. Verify combined
	finalContent, appErr := wrapper.ReadString(filePath, FileActionTypeReadOnly, FileModeTypeDefaultFile)
	if appErr != nil {
		t.Fatalf("read after append failed: %v", appErr)
	}
	if finalContent != content+appendContent {
		t.Errorf("got %q, want %q", finalContent, content+appendContent)
	}

	// 5. Delete
	appErr = wrapper.Delete(filePath)
	if appErr != nil {
		t.Fatalf("delete failed: %v", appErr)
	}
}

func TestFileWrapper_WrapFailuresWithErrorId(t *testing.T) {
	wrapper := NewDefault()

	rawErr := errors.New("underlying disk failure")

	// WrapFailure
	appErr := wrapper.WrapFailure(rawErr, "ERR_CUSTOM_DISK", "disk error occurred")
	if appErr == nil {
		t.Fatalf("expected non-nil AppError")
	}
	if appErr.ErrorId() != "ERR_CUSTOM_DISK" {
		t.Errorf("got errorId %q, want ERR_CUSTOM_DISK", appErr.ErrorId())
	}

	// WrapReaderFailure
	readErr := wrapper.WrapReaderFailure(rawErr, "ERR_READ_TEST", "/path/to/file")
	if readErr.ErrorId() != "ERR_READ_TEST" {
		t.Errorf("got errorId %q, want ERR_READ_TEST", readErr.ErrorId())
	}

	// WrapWriterFailure
	writeErr := wrapper.WrapWriterFailure(rawErr, "ERR_WRITE_TEST", "/path/to/file")
	if writeErr.ErrorId() != "ERR_WRITE_TEST" {
		t.Errorf("got errorId %q, want ERR_WRITE_TEST", writeErr.ErrorId())
	}

	// Nil safety
	if wrapper.WrapFailure(nil, "NONE", "msg") != nil {
		t.Errorf("expected nil for nil cause")
	}
}

func TestFileWrapper_NonExistentFileError(t *testing.T) {
	wrapper := NewDefault()
	_, appErr := wrapper.Read("non-existent-path/file.txt", FileActionTypeReadOnly, FileModeTypeDefaultFile)
	if appErr == nil {
		t.Fatalf("expected error on non-existent file")
	}
	if appErr.ErrorId() != "ERR_FILE_OPEN" {
		t.Errorf("got errorId %q, want ERR_FILE_OPEN", appErr.ErrorId())
	}
}

func TestFileWrapper_Execute(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "exec.txt")

	wrapper := NewDefault()

	// Execute write
	data := []byte("exec-data")
	_, appErr := wrapper.Execute(filePath, FileActionTypeWrite, FileModeTypeDefaultFile, data)
	if appErr != nil {
		t.Fatalf("execute write failed: %v", appErr)
	}

	// Execute read
	readData, appErr := wrapper.Execute(filePath, FileActionTypeRead, FileModeTypeDefaultFile, nil)
	if appErr != nil {
		t.Fatalf("execute read failed: %v", appErr)
	}
	if string(readData) != "exec-data" {
		t.Errorf("got %q, want exec-data", string(readData))
	}

	// Execute delete
	_, appErr = wrapper.Execute(filePath, FileActionTypeDelete, FileModeTypeDefaultFile, nil)
	if appErr != nil {
		t.Fatalf("execute delete failed: %v", appErr)
	}
}

func TestPackageLevelConvenience(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "conv.txt")

	appErr := WriteFile(filePath, []byte("pkg-data"), FileActionTypeWrite, FileModeTypeDefaultFile)
	if appErr != nil {
		t.Fatalf("WriteFile failed: %v", appErr)
	}

	str, appErr := ReadFileString(filePath, FileActionTypeReadOnly, FileModeTypeDefaultFile)
	if appErr != nil {
		t.Fatalf("ReadFileString failed: %v", appErr)
	}
	if str != "pkg-data" {
		t.Errorf("got %q, want pkg-data", str)
	}

	appErr = DeleteFile(filePath)
	if appErr != nil {
		t.Fatalf("DeleteFile failed: %v", appErr)
	}
}
