package appwriter_test

import (
	"context"
	"path/filepath"
	"testing"

	"coding-guidelines/common/pkg/appwriter"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/fileutil"
)

func TestNewFileWriter_Success(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "logs", "app.log")

	wrap := appwriter.NewFileWriter(appwriter.FileWriterOptions{
		Name:     "test-logger",
		FilePath: logPath,
		OpenMode: fileutil.FileOpenCreateAppend,
		PermMode: fileutil.FilePermStandard,
		IsLocked: true,
	})

	if wrap.IsFailed() {
		t.Fatalf("expected success, got fault: %v", wrap.Fault())
	}

	w := wrap.Data()
	defer w.Close()

	if w.Name() != "test-logger" {
		t.Fatalf("expected name 'test-logger', got: %s", w.Name())
	}

	isLocked := w.IsLocked()
	if !isLocked {
		t.Fatalf("expected writer to be locked")
	}

	ctx := context.Background()
	fault := w.Write(ctx, "event line 1\n")
	if fault != nil {
		t.Fatalf("write failed: %v", fault)
	}

	syncFault := w.Sync()
	if syncFault != nil {
		t.Fatalf("sync failed: %v", syncFault)
	}

	readRes := fileutil.ReadString(logPath)
	if readRes.IsFailed() {
		t.Fatalf("read failed: %v", readRes.Fault())
	}

	if readRes.Data() != "event line 1\n" {
		t.Fatalf("unexpected content: %s", readRes.Data())
	}
}

func TestNewFileWriter_EmptyPathFailureWithId(t *testing.T) {
	wrap := appwriter.NewFileWriter(appwriter.FileWriterOptions{
		FilePath: "",
	})

	if wrap.IsSuccess() {
		t.Fatalf("expected failure for empty file path")
	}

	fault := wrap.Fault()
	if fault == nil {
		t.Fatalf("expected non-nil fault")
	}

	if fault.Type() != errtype.Validation {
		t.Fatalf("expected errtype.Validation, got: %v", fault.Type())
	}
}

func TestWrapWriterFailureFromWrap(t *testing.T) {
	failedFile := fileutil.OpenFile("", fileutil.FileOpenReadOnly, fileutil.FilePermStandard)
	writerWrap := appwriter.WrapWriterFailureFromWrap(failedFile)

	if writerWrap.IsSuccess() {
		t.Fatalf("expected failed writer wrap")
	}

	if writerWrap.Fault().Type() != errtype.Validation {
		t.Fatalf("expected validation error type in propagated fault")
	}
}
