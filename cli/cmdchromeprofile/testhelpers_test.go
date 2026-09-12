package cmdchromeprofile

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var stdIOMutex sync.Mutex

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	stdIOMutex.Lock()
	defer stdIOMutex.Unlock()

	origStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stderr = w

	outC := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	fn()
	_ = w.Close()
	os.Stderr = origStderr
	res := <-outC
	_ = r.Close()

	return res
}

func seedChromeProfileTree(t *testing.T, root string) {
	t.Helper()
	paths := []string{
		filepath.Join(root, "Default"),
		filepath.Join(root, "Profile 1"),
	}
	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			t.Fatalf("seed profile dir: %v", err)
		}
		f := filepath.Join(p, "Preferences")
		if err := os.WriteFile(f, []byte(`{"profile":{"name":"Test"}}`), 0644); err != nil {
			t.Fatalf("seed preferences: %v", err)
		}
	}
}
