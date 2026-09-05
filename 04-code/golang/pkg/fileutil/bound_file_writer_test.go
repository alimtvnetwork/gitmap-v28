package fileutil_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/fileutil"
)

func TestBoundFileWriter_AutoLockWriteAndAppend(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "auto-lock.txt")
	ctx := context.Background()

	writer := fileutil.NewBoundFileWriter(targetPath)

	// 1. Automatic lock write
	err := writer.WriteString(ctx, "Initial Line\n")
	if err != nil {
		t.Fatalf("expected write success, got: %v", err)
	}

	if writer.BytesWritten() != 13 {
		t.Fatalf("expected 13 bytes written, got: %d", writer.BytesWritten())
	}

	// 2. Automatic lock append
	err = writer.AppendString(ctx, "Appended Line\n")
	if err != nil {
		t.Fatalf("expected append success, got: %v", err)
	}

	if writer.BytesAppended() != 14 {
		t.Fatalf("expected 14 bytes appended, got: %d", writer.BytesAppended())
	}

	if writer.WriteCount() != 2 {
		t.Fatalf("expected 2 write operations, got: %d", writer.WriteCount())
	}

	// 3. Verify content
	content, readErr := os.ReadFile(targetPath)
	expected := "Initial Line\nAppended Line\n"
	if readErr != nil || string(content) != expected {
		t.Fatalf("unexpected content: %s", string(content))
	}

	if closeErr := writer.Close(); closeErr != nil {
		t.Fatalf("close failed: %v", closeErr)
	}
}

func TestBoundFileWriter_AutoCloseMode(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "autoclose.txt")
	ctx := context.Background()

	writer := fileutil.NewBoundFileWriter(targetPath)
	writer.SetAutoClose(true)

	// 1. Write in autoClose mode
	err := writer.WriteString(ctx, "First Write\n")
	if err != nil {
		t.Fatalf("autoClose write failed: %v", err)
	}

	// File descriptor must be closed immediately after writing
	if writer.IsOpen() {
		t.Fatalf("expected file to be closed immediately after write in autoClose mode")
	}

	// 2. Append in autoClose mode
	err = writer.AppendString(ctx, "Second Append\n")
	if err != nil {
		t.Fatalf("autoClose append failed: %v", err)
	}

	if writer.IsOpen() {
		t.Fatalf("expected file to be closed immediately after append in autoClose mode")
	}

	content, readErr := os.ReadFile(targetPath)
	expected := "First Write\nSecond Append\n"
	if readErr != nil || string(content) != expected {
		t.Fatalf("unexpected content: %s", string(content))
	}
}

func TestBoundFileWriter_WriteAndClose_Explicit(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "explicit-close.txt")
	ctx := context.Background()

	// autoClose is false by default
	writer := fileutil.NewBoundFileWriter(targetPath)

	err := writer.WriteAndClose(ctx, []byte("Payload 1\n"))
	if err != nil {
		t.Fatalf("WriteAndClose failed: %v", err)
	}

	if writer.IsOpen() {
		t.Fatalf("expected file to be closed after WriteAndClose")
	}

	err = writer.AppendAndClose(ctx, []byte("Payload 2\n"))
	if err != nil {
		t.Fatalf("AppendAndClose failed: %v", err)
	}

	if writer.IsOpen() {
		t.Fatalf("expected file to be closed after AppendAndClose")
	}

	content, readErr := os.ReadFile(targetPath)
	expected := "Payload 1\nPayload 2\n"
	if readErr != nil || string(content) != expected {
		t.Fatalf("unexpected content: %s", string(content))
	}
}

func TestBoundFileWriter_WithLock_TransactionalBatch(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "batch.txt")
	ctx := context.Background()

	writer := fileutil.NewBoundFileWriter(targetPath)

	// Execute transactional batch of writes under single lock
	err := writer.WithLock(ctx, func(w *fileutil.BoundFileWriter) *appfault.AppError {
		_ = w.WriteLocked(ctx, []byte("Header: Batch 001\n"))
		_ = w.AppendLocked(ctx, []byte("Record: Item A\n"))
		_ = w.AppendLocked(ctx, []byte("Record: Item B\n"))
		_ = w.AppendLocked(ctx, []byte("Footer: Done\n"))

		return nil
	})

	if err != nil {
		t.Fatalf("WithLock failed: %v", err)
	}

	_ = writer.Close()

	content, readErr := os.ReadFile(targetPath)
	expected := "Header: Batch 001\nRecord: Item A\nRecord: Item B\nFooter: Done\n"
	if readErr != nil || string(content) != expected {
		t.Fatalf("unexpected batch content: %s", string(content))
	}
}

func TestBoundFileWriter_ModeShifting(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "shift.txt")
	ctx := context.Background()

	writer := fileutil.NewBoundFileWriter(targetPath)

	// 1. Direct mode
	_ = writer.WriteString(ctx, "Direct Content\n")

	// 2. Shift to Atomic mode
	writer.SetMode(fileutil.FileWriteModeAtomic)
	_ = writer.WriteString(ctx, "Atomic Swapped\n")

	// 3. Shift to Truncate mode with fsync
	writer.SetMode(fileutil.FileWriteModeTruncate).SetSyncOnWrite(true)
	_ = writer.WriteString(ctx, "Truncated Single Line")

	_ = writer.Close()

	content, readErr := os.ReadFile(targetPath)
	if readErr != nil || string(content) != "Truncated Single Line" {
		t.Fatalf("unexpected content after mode shifts: %s", string(content))
	}
}

func TestBoundFileWriter_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "concurrent.log")
	ctx := context.Background()

	writer := fileutil.NewBoundFileWriter(targetPath)
	writer.SetAutoClose(true) // Stress test opening and closing concurrently

	var wg sync.WaitGroup
	iterations := 50

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			line := []byte("Concurrent Log Entry\n")
			_ = writer.Append(ctx, line)
		}(i)
	}

	wg.Wait()

	if writer.BytesAppended() != int64(iterations*21) {
		t.Fatalf("expected %d bytes appended, got %d", iterations*21, writer.BytesAppended())
	}
}

func TestBoundFileWriter_StdWriterAndAppender(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "std-adapters.txt")

	writer := fileutil.NewBoundFileWriter(targetPath)

	stdW := writer.StdWriter()
	_, _ = stdW.Write([]byte("Standard Writer\n"))
	_ = stdW.Close()

	stdApp := writer.StdAppender()
	_, _ = stdApp.Write([]byte("Standard Appender\n"))
	_ = stdApp.Close()

	content, readErr := os.ReadFile(targetPath)
	expected := "Standard Writer\nStandard Appender\n"
	if readErr != nil || string(content) != expected {
		t.Fatalf("unexpected adapter content: %s", string(content))
	}
}
