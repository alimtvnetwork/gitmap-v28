package cmdvscode

import (
	"bytes"
	"io"
	"os"
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
