package cmdos

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestGetLocalOSInfo(t *testing.T) {
	info := GetLocalOSInfo()
	if info.OSType == "" {
		t.Error("expected non-empty OSType")
	}
	if info.Architecture == "" {
		t.Error("expected non-empty Architecture")
	}
	if info.Platform == "" {
		t.Error("expected non-empty Platform")
	}
	if info.NumCPU <= 0 {
		t.Errorf("expected positive NumCPU, got %d", info.NumCPU)
	}
}

func TestRenderOSInfo_Text(t *testing.T) {
	info := OSInfoReport{
		OSType:       "windows",
		OSVersion:    "Windows 11 Pro (Build 22631)",
		Architecture: "amd64",
		Platform:     "windows/amd64",
		Hostname:     "test-host",
		NumCPU:       8,
	}

	var buf bytes.Buffer
	err := RenderOSInfo(&buf, info, false)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "windows") || !strings.Contains(out, "Windows 11 Pro") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestRenderOSInfo_JSON(t *testing.T) {
	info := OSInfoReport{
		OSType:       "ubuntu",
		OSVersion:    "Ubuntu 22.04 LTS",
		Architecture: "amd64",
		Platform:     "linux/amd64",
		Hostname:     "u1",
		NumCPU:       4,
	}

	var buf bytes.Buffer
	err := RenderOSInfo(&buf, info, true)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	var decoded OSInfoReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid json generated: %v", err)
	}
	if decoded.OSType != "ubuntu" || decoded.Architecture != "amd64" {
		t.Errorf("unexpected decoded content: %+v", decoded)
	}
}
