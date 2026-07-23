package cmd

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
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

func TestBuildPinCallbackPythonExecutesWithoutFunctionSymbol(t *testing.T) {
	manifest := filepath.Join(t.TempDir(), "pin.json")
	entries := []pinManifestEntry{{
		Path:    "X",
		DataB64: base64.StdEncoding.EncodeToString([]byte("pinned")),
		Blobs:   []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}}
	data, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if err := os.WriteFile(manifest, data, 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	code := buildPinCallbackPython(manifest)
	script := filepath.Join(t.TempDir(), "run.py")
	py := "" +
		"callback_globals = {'__builtins__': __builtins__}\n" +
		"callback_locals = {}\n" +
		"code = " + reprForPython(code) + "\n" +
		"exec('def callback(blob, metadata=None):\\n' + '  ' + '\\n  '.join(code.splitlines()), callback_globals, callback_locals)\n" +
		"callback = callback_locals['callback']\n" +
		"class Blob:\n" +
		"    def __init__(self, oid):\n" +
		"        self.original_id = oid\n" +
		"        self.data = b'orig'\n" +
		"blob = Blob(b'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa')\n" +
		"callback(blob, {})\n" +
		"assert blob.data == b'pinned', blob.data\n" +
		"callback(blob, {})\n" +
		"assert blob.data == b'pinned', blob.data\n"
	if err := os.WriteFile(script, []byte(py), 0o600); err != nil {
		t.Fatalf("write python script: %v", err)
	}
	cmd := exec.Command("python3", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("python callback execution failed: %v\n%s", err, string(out))
	}
}

func reprForPython(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
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
