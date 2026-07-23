package cmd

// Shared test helper: capture os.Stderr output produced by fn.
//
// Lifted out of scanworkersalias_test.go and clonepmsync_debugpaths_test.go
// (which each defined their own captureStderr) to eliminate the duplicate-
// symbol redeclaration risk in the flat cmd/ test namespace. Both call
// sites now share this single definition.
//
// Implementation note: a goroutine drains the pipe concurrently with fn,
// which prevents a deadlock when fn writes more than one OS pipe buffer
// (~64 KiB on Linux) before the helper gets a chance to read. The simpler
// "fn(); io.ReadAll(r)" form is safe only for tiny outputs and was the
// weaker of the two original implementations; this is the stronger one.

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// captureStderr swaps os.Stderr for a pipe, runs fn, and returns
// whatever fn wrote to stderr.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

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
