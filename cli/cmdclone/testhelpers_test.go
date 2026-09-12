package cmdclone

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func captureStdoutForTest(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe failed: %v", err)
	}
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = oldStdout
	out := <-outC
	_ = r.Close()

	return out
}
