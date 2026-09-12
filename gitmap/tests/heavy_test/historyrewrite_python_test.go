package heavy_test

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
)

type testPinManifestEntry struct {
	Path    string   `json:"path"`
	DataB64 string   `json:"data_b64"`
	Blobs   []string `json:"blobs"`
}

func TestBuildPinCallbackPythonExecutesWithoutFunctionSymbol(t *testing.T) {
	manifest := filepath.Join(t.TempDir(), "pin.json")
	entries := []testPinManifestEntry{{
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

	code := cmd.BuildPinCallbackPython(manifest)
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

	pyBin := "python3"
	_, err3 := exec.LookPath("python3")
	_, errPy := exec.LookPath("python")

	if err3 != nil && errPy == nil {
		pyBin = "python"
	}

	if err3 != nil && errPy != nil {
		t.Skip("python not found on PATH; skipping python callback execution test")
	}

	c := exec.Command(pyBin, script)
	out, err := c.CombinedOutput()
	isStub := strings.Contains(string(out), "Python was not found")

	if err != nil && isStub {
		t.Skip("python stub found but python is not installed; skipping")
	}

	if err != nil {
		t.Fatalf("python callback execution failed: %v\n%s", err, string(out))
	}
}

func reprForPython(s string) string {
	b, _ := json.Marshal(s)

	return string(b)
}
