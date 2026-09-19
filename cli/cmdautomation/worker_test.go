package cmdautomation

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildFileContext(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "sample.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	assertFileContextPreRead(t, testFile)
	assertFileContextNoPreRead(t, testFile)
}

func assertFileContextPreRead(t *testing.T, path string) {
	res := BuildFileContext(path, "utf-8", true)
	if res.IsFailure() {
		t.Fatalf("unexpected failure: %v", res.Err)
	}
	if res.Value.FileName != "sample.txt" || res.Value.Content != "test content" {
		t.Errorf("unexpected file context: %+v", res.Value)
	}
}

func assertFileContextNoPreRead(t *testing.T, path string) {
	res := BuildFileContext(path, "utf-8", false)
	if res.IsFailure() {
		t.Fatalf("unexpected failure: %v", res.Err)
	}
	if res.Value.Content != "" {
		t.Errorf("expected empty content when isPreRead is false, got: %s", res.Value.Content)
	}
}

func TestPartitionFiles(t *testing.T) {
	files := []string{"f1.go", "f2.go", "f3.go", "f4.go", "f5.go"}
	assertPartitionBalanced(t, files, 2)
	assertPartitionEdgeCases(t, files)
}

func assertPartitionBalanced(t *testing.T, files []string, workers int) {
	chunks := PartitionFiles(files, workers)
	if len(chunks) != workers {
		t.Fatalf("expected %d partitions, got %d", workers, len(chunks))
	}
	total := len(chunks[0]) + len(chunks[1])
	if total != len(files) {
		t.Errorf("expected %d files across partitions, got %d", len(files), total)
	}
}

func assertPartitionEdgeCases(t *testing.T, files []string) {
	empty := PartitionFiles([]string{}, 4)
	if len(empty) != 0 {
		t.Errorf("expected 0 partitions for empty list, got %d", len(empty))
	}
	zeroWorkers := PartitionFiles(files, 0)
	if len(zeroWorkers) == 0 {
		t.Error("expected at least 1 partition when workers <= 0")
	}
}

func TestEncodeToStream(t *testing.T) {
	ctx := FileContext{FileName: "root.go", FilePath: "cli/root.go", Encoding: "utf-8"}
	assertUtf8Encoding(t, ctx)
	assertUtf16Encoding(t, ctx)
}

func assertUtf8Encoding(t *testing.T, ctx FileContext) {
	res := EncodeToStream(ctx, "utf-8")
	if res.IsFailure() {
		t.Fatalf("unexpected UTF-8 encode error: %v", res.Err)
	}
	if !bytes.Contains(res.Value, []byte("root.go")) {
		t.Errorf("expected UTF-8 stream to contain 'root.go'")
	}
}

func assertUtf16Encoding(t *testing.T, ctx FileContext) {
	res := EncodeToStream(ctx, "utf16")
	if res.IsFailure() {
		t.Fatalf("unexpected UTF-16 encode error: %v", res.Err)
	}
	if len(res.Value) < 2 || res.Value[0] != 0xFF || res.Value[1] != 0xFE {
		t.Errorf("expected UTF-16LE BOM prefix (0xFF, 0xFE)")
	}
}

func TestDecodeFromStream(t *testing.T) {
	assertDecodeUtf8(t)
	assertDecodeUtf16(t)
}

func assertDecodeUtf8(t *testing.T) {
	raw := []byte("hello stream")
	res := DecodeFromStream(raw)
	if res.IsFailure() {
		t.Fatalf("unexpected decode error: %v", res.Err)
	}
	if res.Value != "hello stream" {
		t.Errorf("expected 'hello stream', got %q", res.Value)
	}
}

func assertDecodeUtf16(t *testing.T) {
	raw := []byte{0xFF, 0xFE, 'h', 0x00, 'i', 0x00}
	res := DecodeFromStream(raw)
	if res.IsFailure() {
		t.Fatalf("unexpected decode error: %v", res.Err)
	}
	if res.Value != "hi" {
		t.Errorf("expected 'hi', got %q", res.Value)
	}
}

func TestProbeRuntime_Fallback(t *testing.T) {
	assertMissingRuntimeProbe(t)
	assertValidRuntimeProbe(t)
}

func assertMissingRuntimeProbe(t *testing.T) {
	res := ProbeRuntime("nonexistent_lang_999")
	if res.IsSuccess() && res.Value.Status == "active" {
		t.Error("expected nonexistent runtime to fail or be marked missing")
	}
	if res.IsFailure() && res.Err.Code != "E7100" && res.Err.Code != "RUNTIME_MISSING" {
		t.Logf("probe returned failure as expected: %v", res.Err)
	}
}

func assertValidRuntimeProbe(t *testing.T) {
	res := ProbeRuntime("go")
	if res.IsSuccess() && res.Value.Name != "go" {
		t.Errorf("expected runtime name 'go', got: %s", res.Value.Name)
	}
}

func TestFormatMissingRuntimeMessage(t *testing.T) {
	assertPythonMissingMessage(t)
	assertRustMissingMessage(t)
}

func assertPythonMissingMessage(t *testing.T) {
	msg := FormatMissingRuntimeMessage("python")
	if !strings.Contains(msg, "E7100:RUNTIME_MISSING") {
		t.Errorf("expected E7100:RUNTIME_MISSING in message, got: %s", msg)
	}
	if !strings.Contains(msg, "gitmap install python") {
		t.Errorf("expected 'gitmap install python' in message")
	}
	if !strings.Contains(msg, "gitmap install profile dev-full") {
		t.Errorf("expected profile suggestion in message")
	}
}

func assertRustMissingMessage(t *testing.T) {
	msg := FormatMissingRuntimeMessage("rust")
	if !strings.Contains(msg, "gitmap install rust") {
		t.Errorf("expected 'gitmap install rust' in message")
	}
}
