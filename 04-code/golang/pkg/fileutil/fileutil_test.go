package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFile_CreateAppend(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nested", "test.log")

	res := OpenFile(filePath, FileOpenCreateAppend, FilePermStandard)
	if res.IsFailed() {
		t.Fatalf("expected success, got error: %v", res.Fault())
	}

	f := res.Data()
	_, err := f.WriteString("hello log\n")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	f.Close()

	readRes := ReadString(filePath)
	if readRes.IsFailed() {
		t.Fatalf("read failed: %v", readRes.Fault())
	}

	if readRes.Data() != "hello log\n" {
		t.Fatalf("unexpected content: %s", readRes.Data())
	}
}

func TestOpenFile_NotFound(t *testing.T) {
	nonExistent := filepath.Join(t.TempDir(), "missing.txt")
	res := OpenFile(nonExistent, FileOpenReadOnly, FilePermStandard)
	if res.IsSuccess() {
		t.Fatalf("expected failure for non-existent file")
	}

	if res.Fault() == nil {
		t.Fatalf("expected fault object")
	}
}

func TestEnums_NamesAndFlags(t *testing.T) {
	if FileOpenCreateAppend.Flags() == 0 {
		t.Fatalf("expected non-zero flags for FileOpenCreateAppend")
	}

	if FileOpenCreateAppend.Name() != "CreateAppend" {
		t.Fatalf("unexpected name: %s", FileOpenCreateAppend.Name())
	}

	if FileOpenReadOrCreateOnly.Name() != "ReadOrCreateOnly" {
		t.Fatalf("unexpected name: %s", FileOpenReadOrCreateOnly.Name())
	}

	if FileOpenWriteOrCreateOnly.Name() != "WriteOrCreateOnly" {
		t.Fatalf("unexpected name: %s", FileOpenWriteOrCreateOnly.Name())
	}

	if FileOpenReadWriteOrCreateOnly.Name() != "ReadWriteOrCreateOnly" {
		t.Fatalf("unexpected name: %s", FileOpenReadWriteOrCreateOnly.Name())
	}

	if FilePermStandard.Mode() != os.FileMode(0644) {
		t.Fatalf("unexpected mode: %v", FilePermStandard.Mode())
	}
}

func TestFilePermType_StringsAndOctal(t *testing.T) {
	if FilePermStandard.OctalString() != "0644" {
		t.Fatalf("expected 0644, got: %s", FilePermStandard.OctalString())
	}

	if FilePermStandard.PosixString() != "rw-r--r--" {
		t.Fatalf("expected rw-r--r--, got: %s", FilePermStandard.PosixString())
	}

	if FilePermExecutable.PosixString() != "rwxr-xr-x" {
		t.Fatalf("expected rwxr-xr-x, got: %s", FilePermExecutable.PosixString())
	}

	if FilePermPrivate.PosixString() != "rw-------" {
		t.Fatalf("expected rw-------, got: %s", FilePermPrivate.PosixString())
	}
}

func TestFilePermType_Inspections(t *testing.T) {
	isPrivate := FilePermPrivate.IsPrivate()
	if !isPrivate {
		t.Fatalf("expected FilePermPrivate to be private")
	}

	isPublic := FilePermStandard.IsPublic()
	if !isPublic {
		t.Fatalf("expected FilePermStandard to be public")
	}

	isExec := FilePermExecutable.IsExecutable()
	if !isExec {
		t.Fatalf("expected FilePermExecutable to be executable")
	}
}

func TestFilePermType_Modifiers(t *testing.T) {
	standard := FilePermStandard
	private := standard.WithPrivate()
	if private.OctalString() != "0600" {
		t.Fatalf("expected 0600 after WithPrivate, got: %s", private.OctalString())
	}

	readOnly := standard.WithReadOnly()
	if readOnly.OctalString() != "0444" {
		t.Fatalf("expected 0444 after WithReadOnly, got: %s", readOnly.OctalString())
	}

	exec := standard.WithExecutable()
	if exec.OctalString() != "0755" {
		t.Fatalf("expected 0755 after WithExecutable, got: %s", exec.OctalString())
	}
}

func TestFilePermType_ParseAndConvert(t *testing.T) {
	wrap := ParsePerm("0644")
	if wrap.IsFailed() {
		t.Fatalf("parse failed: %v", wrap.Fault())
	}

	if wrap.Data() != FilePermStandard {
		t.Fatalf("expected FilePermStandard, got: %v", wrap.Data())
	}

	invalidWrap := ParsePerm("invalid")
	if invalidWrap.IsSuccess() {
		t.Fatalf("expected parse failure on invalid string")
	}

	fromMode := FromFileMode(os.FileMode(0755))
	if fromMode != FilePermExecutable {
		t.Fatalf("expected FilePermExecutable, got: %v", fromMode)
	}
}

func TestFileOpType_Enums(t *testing.T) {
	if FileOpReadOnly.Name() != "ReadOnly" {
		t.Fatalf("unexpected name: %s", FileOpReadOnly.Name())
	}

	if FileOpDelete.Name() != "Delete" {
		t.Fatalf("unexpected name: %s", FileOpDelete.Name())
	}

	if !FileOpDelete.IsDelete() {
		t.Fatalf("expected FileOpDelete to be delete")
	}

	if !FileOpReadOnly.IsReadOnly() {
		t.Fatalf("expected FileOpReadOnly to be read only")
	}

	if !FileOpAppend.IsAppend() {
		t.Fatalf("expected FileOpAppend to be append")
	}

	if FileOpCreate.OpenMode() != FileOpenCreateNew {
		t.Fatalf("expected FileOpenCreateNew open mode")
	}
}

func TestFileutil_DeleteAndStat(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sample.txt")

	writeRes := WriteFile(filePath, []byte("hello world"), FilePermStandard)
	if writeRes.IsFailed() {
		t.Fatalf("write failed: %v", writeRes.Fault())
	}

	readRes := ReadFile(filePath)
	if readRes.IsFailed() {
		t.Fatalf("read failed: %v", readRes.Fault())
	}

	if string(readRes.Data()) != "hello world" {
		t.Fatalf("unexpected content: %s", string(readRes.Data()))
	}

	statRes := Stat(filePath)
	if statRes.IsFailed() {
		t.Fatalf("stat failed: %v", statRes.Fault())
	}

	sizeRes := FileSize(filePath)
	if sizeRes.IsFailed() {
		t.Fatalf("size failed: %v", sizeRes.Fault())
	}

	if sizeRes.Data() != int64(11) {
		t.Fatalf("expected size 11, got %d", sizeRes.Data())
	}

	delRes := DeleteFile(filePath)
	if delRes.IsFailed() {
		t.Fatalf("delete failed: %v", delRes.Fault())
	}

	statAfter := Stat(filePath)
	if statAfter.IsSuccess() {
		t.Fatalf("expected failure after delete")
	}
}

func TestFileutil_ExecuteOp(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "op_test.txt")

	createRes := ExecuteOp(filePath, FileOpCreateAppend, FilePermStandard, []byte("step1\n"))
	if createRes.IsFailed() {
		t.Fatalf("create failed: %v", createRes.Fault())
	}

	appendRes := ExecuteOp(filePath, FileOpAppend, FilePermStandard, []byte("step2\n"))
	if appendRes.IsFailed() {
		t.Fatalf("append failed: %v", appendRes.Fault())
	}

	readRes := ExecuteOp(filePath, FileOpReadOnly, FilePermStandard, nil)
	if readRes.IsFailed() {
		t.Fatalf("read op failed: %v", readRes.Fault())
	}

	if string(readRes.Data()) != "step1\nstep2\n" {
		t.Fatalf("unexpected content: %s", string(readRes.Data()))
	}

	delRes := ExecuteOp(filePath, FileOpDelete, FilePermStandard, nil)
	if delRes.IsFailed() {
		t.Fatalf("delete op failed: %v", delRes.Fault())
	}
}
