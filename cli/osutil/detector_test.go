package osutil

import (
	"testing"
)

func TestDetectHostOS(t *testing.T) {
	info := DetectHostOS()
	if info.Family == "" || info.PackageManager == "" {
		t.Fatalf("expected populated OS info, got %+v", info)
	}
}
