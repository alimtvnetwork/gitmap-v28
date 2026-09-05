package fileutil_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/fileutil"
)

func TestWriteAtomic(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "atomic-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(tempDir)

	targetFile := filepath.Join(tempDir, "atomic-test.txt")
	initialData := []byte("first atomic revision")

	// 1. Initial atomic write
	res := fileutil.WriteAtomic(targetFile, initialData, fileutil.FilePermStandard)
	if res.IsFailed() {
		t.Fatalf("WriteAtomic failed: %v", res.Fault())
	}

	readRes := fileutil.ReadString(targetFile)
	if readRes.IsFailed() {
		t.Fatalf("ReadString failed: %v", readRes.Fault())
	}

	if readRes.Data() != "first atomic revision" {
		t.Fatalf("unexpected read data: %s", readRes.Data())
	}

	// 2. Overwrite atomically
	updatedData := []byte("second atomic revision - updated cleanly")
	res = fileutil.WriteAtomic(targetFile, updatedData, fileutil.FilePermStandard)
	if res.IsFailed() {
		t.Fatalf("WriteAtomic overwrite failed: %v", res.Fault())
	}

	readRes = fileutil.ReadString(targetFile)
	if readRes.IsFailed() {
		t.Fatalf("ReadString after overwrite failed: %v", readRes.Fault())
	}

	if readRes.Data() != "second atomic revision - updated cleanly" {
		t.Fatalf("unexpected updated data: %s", readRes.Data())
	}
}

func TestReadChunked(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "chunked-read-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(tempDir)

	targetFile := filepath.Join(tempDir, "chunked.bin")
	payload := bytes.Repeat([]byte("0123456789ABCDEF"), 1024) // 16KB

	writeRes := fileutil.WriteFile(targetFile, payload, fileutil.FilePermStandard)
	if writeRes.IsFailed() {
		t.Fatalf("WriteFile failed: %v", writeRes.Fault())
	}

	var accumulated bytes.Buffer
	chunkCount := 0

	readRes := fileutil.ReadChunked(targetFile, 4096, func(chunk []byte) *appfault.AppError {
		chunkCount++
		accumulated.Write(chunk)

		return nil
	})

	if readRes.IsFailed() {
		t.Fatalf("ReadChunked failed: %v", readRes.Fault())
	}

	if readRes.Data() != int64(len(payload)) {
		t.Fatalf("expected %d bytes, got %d", len(payload), readRes.Data())
	}

	if chunkCount != 4 {
		t.Fatalf("expected 4 chunks of 4KB, got %d", chunkCount)
	}

	if !bytes.Equal(accumulated.Bytes(), payload) {
		t.Fatalf("accumulated bytes do not match original payload")
	}
}

func TestWriteChunked(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "chunked-write-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(tempDir)

	targetFile := filepath.Join(tempDir, "write-chunked.bin")
	payload := bytes.Repeat([]byte("CHUNKED-DATA-BLOCK"), 512)
	reader := bytes.NewReader(payload)

	res := fileutil.WriteChunked(targetFile, fileutil.FilePermStandard, reader, 2048)
	if res.IsFailed() {
		t.Fatalf("WriteChunked failed: %v", res.Fault())
	}

	if res.Data() != int64(len(payload)) {
		t.Fatalf("expected %d bytes written, got %d", len(payload), res.Data())
	}

	readRes := fileutil.ReadAll(targetFile)
	if readRes.IsFailed() {
		t.Fatalf("ReadAll failed: %v", readRes.Fault())
	}

	if !bytes.Equal(readRes.Data(), payload) {
		t.Fatalf("read bytes do not match written payload")
	}
}

func TestNewFileWriter(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "file-writer-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(tempDir)

	targetFile := filepath.Join(tempDir, "writer-output.txt")
	writerRes := fileutil.NewFileWriter(targetFile, fileutil.FileOpenCreateAppend, fileutil.FilePermStandard)
	if writerRes.IsFailed() {
		t.Fatalf("NewFileWriter failed: %v", writerRes.Fault())
	}

	writer := writerRes.Data()
	appErr := writer.Write(context.Background(), "log event 1")
	if appErr != nil {
		t.Fatalf("writer.Write failed: %v", appErr)
	}

	// Close writer to flush
	_ = writer.Close()

	readRes := fileutil.ReadString(targetFile)
	if readRes.IsFailed() {
		t.Fatalf("ReadString failed: %v", readRes.Fault())
	}

	if len(readRes.Data()) == 0 {
		t.Fatalf("expected non-empty written file")
	}
}
