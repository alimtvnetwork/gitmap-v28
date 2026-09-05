package fileutil_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"coding-guidelines/common/pkg/fileutil"
)

func TestFileWriter_BehaviorShifting(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "shift-test.txt")

	writer := fileutil.NewFileWriterEngine(targetPath)
	ctx := context.Background()

	// 1. Initial behavior: Direct Write
	if writer.Mode() != fileutil.FileWriteModeDirect {
		t.Fatalf("expected initial mode FileWriteModeDirect, got %v", writer.Mode())
	}

	err := writer.WriteString(ctx, "Direct Content 1\n")
	if err != nil {
		t.Fatalf("direct write failed: %v", err)
	}

	content, readErr := os.ReadFile(targetPath)
	if readErr != nil || string(content) != "Direct Content 1\n" {
		t.Fatalf("unexpected content after direct write: %s", string(content))
	}

	// 2. Behavior shift: Switch to Atomic Mode
	writer.SetMode(fileutil.FileWriteModeAtomic)
	if writer.Mode() != fileutil.FileWriteModeAtomic {
		t.Fatalf("expected shifted mode FileWriteModeAtomic")
	}

	err = writer.WriteString(ctx, "Atomic Swapped Content\n")
	if err != nil {
		t.Fatalf("atomic write failed: %v", err)
	}

	content, readErr = os.ReadFile(targetPath)
	if readErr != nil || string(content) != "Atomic Swapped Content\n" {
		t.Fatalf("unexpected content after atomic shift: %s", string(content))
	}

	// 3. Behavior shift: Switch to Truncate Mode with sync
	writer.SetMode(fileutil.FileWriteModeTruncate).SetSyncOnWrite(true)
	err = writer.WriteString(ctx, "Truncated Content")
	if err != nil {
		t.Fatalf("truncate write failed: %v", err)
	}

	content, _ = os.ReadFile(targetPath)
	if string(content) != "Truncated Content" {
		t.Fatalf("unexpected content after truncate shift: %s", string(content))
	}
}

func TestFileAppender_ContinuousAppend(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "nested", "sub", "appended.log")

	appender := fileutil.NewFileAppender(targetPath, fileutil.FilePermStandard)
	ctx := context.Background()

	// 1. Append creates nested directory and file automatically
	err := appender.AppendString(ctx, "Line 1: system started\n")
	if err != nil {
		t.Fatalf("first append failed: %v", err)
	}

	err = appender.AppendString(ctx, "Line 2: event received\n")
	if err != nil {
		t.Fatalf("second append failed: %v", err)
	}

	// 2. AutoSync behavior shift
	appender.SetAutoSync(true)
	err = appender.AppendString(ctx, "Line 3: critical alert\n")
	if err != nil {
		t.Fatalf("third append with sync failed: %v", err)
	}

	if appender.BytesAppended() <= 0 {
		t.Fatalf("expected positive BytesAppended count, got %d", appender.BytesAppended())
	}

	if err := appender.Sync(); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	if err := appender.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	content, readErr := os.ReadFile(targetPath)
	expected := "Line 1: system started\nLine 2: event received\nLine 3: critical alert\n"
	if readErr != nil || string(content) != expected {
		t.Fatalf("unexpected appended file content: %s", string(content))
	}
}

func TestFileWriterAndAppender_StdWriterAdapter(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "std-adapter.txt")

	writer := fileutil.NewFileWriterEngine(filePath)
	stdW := writer.StdWriter()

	n, err := stdW.Write([]byte("Standard Writer Line\n"))
	if err != nil || n != 21 {
		t.Fatalf("std writer write failed: %v, written: %d", err, n)
	}

	if err := stdW.Close(); err != nil {
		t.Fatalf("std writer close failed: %v", err)
	}

	appender := fileutil.NewFileAppender(filePath, fileutil.FilePermStandard)
	stdApp := appender.StdWriter()

	n, err = stdApp.Write([]byte("Standard Appender Line\n"))
	if err != nil || n != 23 {
		t.Fatalf("std appender write failed: %v, written: %d", err, n)
	}

	if err := stdApp.Close(); err != nil {
		t.Fatalf("std appender close failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	expected := "Standard Writer Line\nStandard Appender Line\n"
	if readErr != nil || string(content) != expected {
		t.Fatalf("unexpected adapter content: %s", string(content))
	}
}
