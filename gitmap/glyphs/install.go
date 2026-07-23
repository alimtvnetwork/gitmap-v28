// Package glyphs — install.go: optional pipe-wrap that runs every
// stdout / stderr byte through Filter. Skipped entirely in ModeRich
// for true zero-overhead passthrough.
package glyphs

import (
	"io"
	"os"
	"sync"
)

var (
	installOnce sync.Once
	activeMode  Mode

	// pipeMu guards installedPipes against concurrent Install / Drain.
	// In practice Install runs once and Drain runs at process exit,
	// but the lock keeps the contract explicit for future callers.
	pipeMu         sync.Mutex
	installedPipes []installedPipe
)

// installedPipe records one pipe-wrap so Drain can flush it before the
// process exits. Without Drain, bytes written to os.Stdout/os.Stderr
// just before os.Exit can be discarded on Windows: the forwarding
// goroutine never gets scheduled to copy them from the pipe buffer to
// the real fd inherited from the parent (root cause of the cliexit
// subprocess-test flakiness on Windows CI).
type installedPipe struct {
	w    *os.File
	done chan struct{}
}

// Install resolves the active mode and (when ModeSafe) replaces
// os.Stdout / os.Stderr with pipe-backed writers that filter glyphs.
// Idempotent across calls.
func Install() {
	installOnce.Do(func() {
		activeMode = Resolve()
		if activeMode == ModeRich {
			return
		}
		os.Stdout = wrap(os.Stdout, activeMode)
		os.Stderr = wrap(os.Stderr, activeMode)
	})
}

// Active returns the resolved mode (defaults to ModeRich pre-Install).
func Active() Mode { return activeMode }

// Drain closes every installed pipe writer and waits for the matching
// forwarder goroutine to flush its buffered bytes to the underlying
// destination fd. MUST be called before os.Exit when output integrity
// matters (e.g. cliexit.Fail) — otherwise the last failure message
// can vanish on Windows.
//
// After Drain returns the wrapped os.Stdout / os.Stderr are closed
// and any further writes will fail with ErrClosed. That is the
// correct shape for an exit-path-only flusher.
func Drain() {
	pipeMu.Lock()
	pipes := installedPipes
	installedPipes = nil
	pipeMu.Unlock()
	for _, p := range pipes {
		_ = p.w.Close()
		<-p.done
	}
}

// wrap returns a *os.File whose writes are filtered before reaching dst.
func wrap(dst *os.File, mode Mode) *os.File {
	r, w, err := os.Pipe()
	if err != nil {
		return dst
	}
	done := make(chan struct{})
	go func() {
		forward(r, dst, mode)
		close(done)
	}()
	pipeMu.Lock()
	installedPipes = append(installedPipes, installedPipe{w: w, done: done})
	pipeMu.Unlock()

	return w
}

// forward streams r → Filter → dst until EOF.
func forward(r io.ReadCloser, dst io.Writer, mode Mode) {
	defer func() { _ = r.Close() }()

	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, _ = dst.Write(Filter(buf[:n], mode))
		}
		if err != nil {
			return
		}
	}
}
