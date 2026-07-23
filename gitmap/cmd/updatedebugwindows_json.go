// Package cmd — JSON sink for `--debug-windows` diagnostics.
//
// In addition to the human-readable `[debug-windows]` lines printed
// to stderr, every dump helper also emits a structured NDJSON event
// to a timestamped file under the project's output directory:
//
//	output/gitmap-debug-windows-YYYY-MM-DD_HH-MM-SS.jsonl
//
// One event per line — easy to `jq`/`grep`, easy to ship to a log
// aggregator, and (crucially) survives even when stdout/stderr are
// swallowed by a detached Windows launcher.
//
// Activation:
//  1. `--debug-windows-json` flag (boolean, defaults the path)
//  2. `--debug-windows-json=<path>` to override the file path
//  3. `GITMAP_DEBUG_WINDOWS_JSON=<path>` env var (also auto-forwarded
//     to the Phase 3 cleanup child so its events append to the same
//     file as the parent, giving one consolidated trace per handoff)
//
// The sink is OFF by default — `--debug-windows` alone keeps the
// console-only behavior from v3.86. You opt-in to the file sink
// because writing under the project tree has user-visible side effects.
//
// Failure policy: file open / write errors are swallowed and degrade
// to console-only. Diagnostics must NEVER block or fail the update.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

var (
	debugWinJSONOnce sync.Once
	debugWinJSONPath string
	debugWinJSONFile *os.File
	debugWinJSONMu   sync.Mutex
)

// emitDebugWindowsJSON appends one structured event to the JSON sink.
// Called from every dump helper alongside the existing fmt.Fprintf.
// No-op when the sink is disabled (default).
func emitDebugWindowsJSON(event string, fields map[string]any) {
	if !isDebugWindowsJSONRequested() {
		return
	}
	f := openDebugWindowsJSONFile()
	if f == nil {
		return
	}
	payload := buildDebugWindowsJSONPayload(event, fields)
	line, err := json.Marshal(payload)
	if err != nil {
		return
	}
	debugWinJSONMu.Lock()
	defer debugWinJSONMu.Unlock()
	_, _ = f.Write(append(line, '\n'))
}

// buildDebugWindowsJSONPayload assembles the standard envelope shared
// by every event so downstream tooling can filter/group consistently.
func buildDebugWindowsJSONPayload(event string,
	fields map[string]any) map[string]any {
	self, _ := os.Executable()
	payload := map[string]any{
		"ts":      time.Now().UTC().Format(time.RFC3339Nano),
		"event":   event,
		"pid":     os.Getpid(),
		"ppid":    os.Getppid(),
		"goos":    runtime.GOOS,
		"self":    self,
		"version": constants.Version,
	}
	for k, v := range fields {
		payload[k] = v
	}

	return payload
}

// isDebugWindowsJSONRequested returns true when the JSON sink is
// enabled via flag or env var.
func isDebugWindowsJSONRequested() bool {
	if len(os.Getenv(constants.EnvDebugWindowsJSON)) > 0 {
		return true
	}
	for _, arg := range os.Args[1:] {
		if arg == constants.FlagDebugWindowsJSON ||
			isDebugWindowsJSONFlagWithValue(arg) {
			return true
		}
	}

	return false
}

// isDebugWindowsJSONFlagWithValue handles `--debug-windows-json=<path>`.
func isDebugWindowsJSONFlagWithValue(arg string) bool {
	prefix := constants.FlagDebugWindowsJSON + "="

	return len(arg) > len(prefix) && arg[:len(prefix)] == prefix
}

// openDebugWindowsJSONFile opens (once per process) the sink file in
// append mode. Path resolution order:
//  1. Explicit `--debug-windows-json=<path>` value
//  2. `GITMAP_DEBUG_WINDOWS_JSON` env value
//  3. Default: output/gitmap-debug-windows-<ts>.jsonl
func openDebugWindowsJSONFile() *os.File {
	debugWinJSONOnce.Do(func() {
		path := resolveDebugWindowsJSONPath()
		_ = os.MkdirAll(filepath.Dir(path), constants.DirPermission)
		f, err := os.OpenFile(path,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				constants.MsgDebugWinJSONOpenFail, path, err)

			return
		}
		debugWinJSONFile = f
		debugWinJSONPath = path
		fmt.Fprintf(os.Stderr, constants.MsgDebugWinJSONFile, path)
		// Export the path so the Phase 3 cleanup child appends to the
		// same file rather than creating its own per-process trace.
		_ = os.Setenv(constants.EnvDebugWindowsJSON, path)
	})

	return debugWinJSONFile
}

// resolveDebugWindowsJSONPath picks the sink path from CLI value, env,
// or the default timestamped file under the output directory.
func resolveDebugWindowsJSONPath() string {
	for _, arg := range os.Args[1:] {
		if isDebugWindowsJSONFlagWithValue(arg) {
			return arg[len(constants.FlagDebugWindowsJSON)+1:]
		}
	}
	if envPath := os.Getenv(constants.EnvDebugWindowsJSON); len(envPath) > 0 {
		return envPath
	}
	ts := time.Now().Format("2006-01-02_15-04-05")
	name := fmt.Sprintf(constants.DebugWindowsJSONFileFmt, ts)

	return filepath.Join(constants.DefaultOutputFolder, name)
}

// debugWindowsJSONPath returns the resolved sink path (empty until
// the sink has been opened at least once). Used by the console dump
// to surface the path next to the existing handoff log line.
func debugWindowsJSONPath() string {
	return debugWinJSONPath
}
