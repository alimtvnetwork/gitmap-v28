// Package termout — install.go: theme mode resolution for the
// synchronous output path.
//
// The old pipe-based os.Stdout/os.Stderr interception (wrap/forward
// goroutines + Drain-before-Exit) is gone: program 267 replaced it
// with the synchronous FilterWriter in cli/output, built once by
// cmd.Run after --theme stripping. os.Stdout/os.Stderr are never
// reassigned anymore.
//
// Install and Drain are kept as deprecated thin delegates so
// migration-window callers keep compiling; they only record the
// resolved mode (Install) or no-op (Drain).
package termout

import (
	"fmt"
	"os"
	"sync"
)

var (
	activeMode  ModeType
	stdoutIsTTY = detectTTY(os.Stdout)
	stderrIsTTY = detectTTY(os.Stderr)

	deprecateOnce sync.Once
)

// Install resolves the active theme mode for Active() consumers.
//
// Deprecated: use output.Build — Install no longer intercepts
// os.Stdout/os.Stderr. Emits a one-time stderr warning as the
// grep-able migration signal for the CI gate.
func Install() {
	deprecateOnce.Do(func() {
		fmt.Fprintln(os.Stderr, "output: termout.Install is deprecated; use output.Build")
	})

	activeMode = Resolve()
}

// Active returns the mode recorded by Install. Defaults to
// ModeBright when Install has not run.
func Active() ModeType {
	return activeMode
}

// Drain is a deprecated no-op. The installed-pipe registry it used
// to flush no longer exists — FilterWriter writes synchronously, so
// no flush is needed before os.Exit. Kept so migration-window
// callers (and the cliexit.RegisterFlusher hook) keep compiling.
func Drain() {}

// IsStdoutTTY reports whether the stdout handle is a real terminal.
// Kept for cli/render callers that gate ANSI pretty-rendering on it.
func IsStdoutTTY() bool { return stdoutIsTTY }

// IsStderrTTY is the stderr counterpart of IsStdoutTTY.
func IsStderrTTY() bool { return stderrIsTTY }

// detectTTY captures the TTY state of a handle at package init.
func detectTTY(f *os.File) bool {
	if f == nil {
		return false
	}

	info, err := f.Stat()
	if err != nil {
		return false
	}

	return (info.Mode() & os.ModeCharDevice) != 0
}
