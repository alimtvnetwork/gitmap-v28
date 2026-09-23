package cmdpull

import (
	"strings"
	"testing"
)

func TestFormatPullBannerString(t *testing.T) {
	banner := FormatPullBannerString("pull all-efficient", "pae", "/test/cwd", "6.307.0", true)
	if !strings.Contains(banner, "gitmap pull all-efficient") {
		t.Fatalf("expected command in banner, got %s", banner)
	}
	if !strings.Contains(banner, "v6.307.0") {
		t.Fatalf("expected version tag in banner, got %s", banner)
	}
	if !strings.Contains(banner, "/test/cwd") {
		t.Fatalf("expected cwd in banner, got %s", banner)
	}
	if !strings.Contains(banner, "[alias: pae → gitmap pull all-efficient]") {
		t.Fatalf("expected alias expansion in banner, got %s", banner)
	}
}

func TestFormatVersionTag(t *testing.T) {
	if got := formatVersionTag("6.307.0"); got != "v6.307.0" {
		t.Fatalf("expected v6.307.0, got %s", got)
	}
	if got := formatVersionTag("v1.2.3"); got != "v1.2.3" {
		t.Fatalf("expected v1.2.3, got %s", got)
	}
	if got := formatVersionTag(""); got != "v0.0.0" {
		t.Fatalf("expected v0.0.0, got %s", got)
	}
}
