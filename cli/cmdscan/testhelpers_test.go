package cmdscan

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

var stdIOMutex sync.Mutex

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	stdIOMutex.Lock()
	defer stdIOMutex.Unlock()

	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stderr = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	_ = w.Close()
	os.Stderr = orig

	return <-done
}

func scanPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("resolve cmdscan package dir")
	}

	return filepath.Dir(file)
}

func resolveSchemaDir() string {
	return filepath.Join(filepath.Dir(scanPackageDir()), "cmd", "testdata", "schemas")
}

func findSchemaFile(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(scanPackageDir())
	for i := 0; i < 8; i++ {
		candidateApp := filepath.Join(dir, "spec", "21-app", "08-json-schemas", filename)
		if _, err := os.Stat(candidateApp); err == nil {
			return candidateApp
		}

		candidateRoot := filepath.Join(dir, "spec", "08-json-schemas", filename)
		if _, err := os.Stat(candidateRoot); err == nil {
			return candidateRoot
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	t.Fatalf("could not locate %s in spec/21-app/08-json-schemas or spec/08-json-schemas", filename)
	return ""
}

func loadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	schemaPath := findSchemaFile(t, filename)
	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema %s: %v", schemaPath, err)
	}

	var schema map[string]any
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("parse schema %s: %v", schemaPath, err)
	}

	return schema
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func readEveryObjectKeys(t *testing.T, raw []byte) [][]string {
	t.Helper()
	var records []map[string]any
	if err := json.Unmarshal(raw, &records); err != nil {
		t.Fatalf("unmarshal records: %v", err)
	}

	var out [][]string
	for _, rec := range records {
		var keys []string
		for k := range rec {
			keys = append(keys, k)
		}
		out = append(out, keys)
	}

	return out
}
