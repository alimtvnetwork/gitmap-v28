package cmddownload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDownloadArgs(t *testing.T) {
	args := []string{"https://example.com/archive.zip", "-o", "my_file.zip", "-j", "-t", "8", "-s", "40", "-f"}
	opts, err := parseDownloadArgs(args)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if opts.URL != "https://example.com/archive.zip" {
		t.Errorf("expected URL to be https://example.com/archive.zip, got %s", opts.URL)
	}
	if opts.OutputPath != "my_file.zip" {
		t.Errorf("expected OutputPath to be my_file.zip, got %s", opts.OutputPath)
	}
	if !opts.IsJSON {
		t.Errorf("expected IsJSON to be true")
	}
	if opts.Threads != 8 {
		t.Errorf("expected Threads to be 8, got %d", opts.Threads)
	}
	if opts.Splits != 40 {
		t.Errorf("expected Splits to be 40, got %d", opts.Splits)
	}
	if !opts.IsForceOverwrite {
		t.Errorf("expected IsForceOverwrite to be true")
	}
}

func TestResolveDestination(t *testing.T) {
	_, _, errEmpty := resolveDestination("", "")
	if errEmpty == nil {
		t.Errorf("expected error for empty URL")
	}

	_, _, errInvalid := resolveDestination("ftp://example.com/test", "")
	if errInvalid == nil {
		t.Errorf("expected error for non-http URL")
	}

	dir, name, errOk := resolveDestination("https://example.com/path/to/archive.tar.gz", "")
	if errOk != nil {
		t.Fatalf("unexpected error: %v", errOk)
	}
	if name != "archive.tar.gz" {
		t.Errorf("expected archive.tar.gz, got %s", name)
	}
	if dir == "" {
		t.Errorf("expected non-empty dir")
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		actual := FormatBytes(tt.input)
		if actual != tt.expected {
			t.Errorf("FormatBytes(%d) = %s, expected %s", tt.input, actual, tt.expected)
		}
	}
}

func TestDownloadGoHTTP(t *testing.T) {
	content := []byte("Hello, GitMap download accelerator!")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	}))
	defer ts.Close()

	tempDir, errTemp := os.MkdirTemp("", "gitmap_dl_test_*")
	if errTemp != nil {
		t.Fatalf("failed to create temp dir: %v", errTemp)
	}
	defer os.RemoveAll(tempDir)

	opts := DownloadOptions{
		URL:              ts.URL + "/test.txt",
		OutputPath:       filepath.Join(tempDir, "test.txt"),
		Engine:           EngineGoHTTP,
		IsQuiet:          true,
		IsForceOverwrite: true,
	}

	res := ExecuteDownload(opts)
	if res.IsFailure() {
		t.Fatalf("download failed: %v", res.Err)
	}

	data, errRead := os.ReadFile(res.Value.DestinationPath)
	if errRead != nil {
		t.Fatalf("failed to read downloaded file: %v", errRead)
	}

	if string(data) != string(content) {
		t.Errorf("expected content '%s', got '%s'", string(content), string(data))
	}
}

func TestDownloadJSONSerialization(t *testing.T) {
	res := DownloadResult{
		Status:              "success",
		URL:                 "https://example.com/test.zip",
		DestinationPath:     "C:\\test\\test.zip",
		FileName:            "test.zip",
		FileSizeBytes:       1024,
		DownloadedBytes:     1024,
		DurationMs:          100,
		AverageSpeedBps:     10240,
		EngineUsed:          EngineAria2c,
		Tier:                1,
		HasFallbackOccurred: false,
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(bytes, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if m["engine_used"] != "aria2c" {
		t.Errorf("expected engine_used = aria2c, got %v", m["engine_used"])
	}
	if m["tier"] != float64(1) {
		t.Errorf("expected tier = 1, got %v", m["tier"])
	}
}
