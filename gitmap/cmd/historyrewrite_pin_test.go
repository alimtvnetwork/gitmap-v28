package cmd

import (
	"strings"
	"testing"
)

// TestBuildPinCallbackPythonUsesGlobalsCache ensures the emitted
// filter-repo callback does not depend on a function-object name such
// as `blob_callback`, which is not guaranteed to exist inside the
// wrapper body across filter-repo versions, and caches via builtins.
func TestBuildPinCallbackPythonUsesGlobalsCache(t *testing.T) {
	got := buildPinCallbackPython("/tmp/pin.json")
	if !strings.Contains(got, "getattr(builtins, '_gitmap_pin_lookup', None)") {
		t.Fatalf("callback missing builtins cache lookup: %q", got)
	}

	if !strings.Contains(got, "setattr(builtins, '_gitmap_pin_lookup', _pin_lookup)") {
		t.Fatalf("callback missing builtins cache store: %q", got)
	}

	if strings.Contains(got, "blob_callback") {
		t.Fatalf("callback must not reference blob_callback: %q", got)
	}
}

// TestParseBlobShasFromRawLogRequiresFullSha guards against a regression
// where `git log --raw` was invoked without `--no-abbrev`, returning
// 7-char abbreviated SHAs that the parser silently dropped — leaving
// the pin manifest's `blobs` slice empty and the callback a no-op.
func TestParseBlobShasFromRawLogRequiresFullSha(t *testing.T) {
	abbrev := ":100644 100644 ffe4cdf cbf1d7c M\tX\n"
	if got := parseBlobShasFromRawLog(abbrev); len(got) != 0 {
		t.Fatalf("abbreviated SHAs must be ignored, got %v", got)
	}

	full := ":100644 100644 " +
		strings.Repeat("a", 40) + " " + strings.Repeat("b", 40) + " M\tX\n"
	got := parseBlobShasFromRawLog(full)
	if len(got) != 1 || got[0] != strings.Repeat("b", 40) {
		t.Fatalf("full SHA not parsed, got %v", got)
	}
}
