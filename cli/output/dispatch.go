package output

import (
	"io"
	"os"
	"sync"
)

var (
	// rawStdout/rawStderr are the true process handles, captured at
	// package init. Nothing in the program reassigns os.Stdout /
	// os.Stderr anymore, so these are the byte-faithful
	// destinations; capturing them here also immunizes Raw()
	// against transient swaps (e.g. nodes_clone captureOutput).
	rawStdout = os.Stdout
	rawStderr = os.Stderr

	dispatchMu     sync.RWMutex
	dispatchOut    io.Writer = os.Stdout
	dispatchErr    io.Writer = os.Stderr
	dispatchWereUp bool
)

// SetDispatchWriters installs the filtered UI writers for this
// dispatch. Called once by cmd.Run after --theme/--glyphs stripping;
// the writers then travel with the dispatch via UI()/UIErr().
func SetDispatchWriters(stdout, stderr io.Writer) {
	dispatchMu.Lock()
	defer dispatchMu.Unlock()

	if stdout != nil {
		dispatchOut = stdout
	}

	if stderr != nil {
		dispatchErr = stderr
	}

	dispatchWereUp = true
}

// UI returns the dispatch's filtered stdout writer. Falls back to
// the raw handle when no dispatch writers were installed (library /
// test contexts where cmd.Run never executed).
func UI() io.Writer {
	dispatchMu.RLock()
	defer dispatchMu.RUnlock()

	if !dispatchWereUp {
		return rawStdout
	}

	return dispatchOut
}

// UIErr returns the dispatch's filtered stderr writer, with the same
// fallback contract as UI.
func UIErr() io.Writer {
	dispatchMu.RLock()
	defer dispatchMu.RUnlock()

	if !dispatchWereUp {
		return rawStderr
	}

	return dispatchErr
}

// Raw returns the unfiltered original os.Stdout captured at package
// init, before any dispatch wiring. Byte-faithful commands write
// file payloads through it so no filter ever touches file bytes.
func Raw() *os.File {
	return rawStdout
}

// RawErr returns the unfiltered original os.Stderr captured at
// package init.
func RawErr() *os.File {
	return rawStderr
}
