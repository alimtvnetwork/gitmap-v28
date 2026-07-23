// Package cmd — durable on-disk handoff log for the self-update
// Phase 3 cleanup chain.
//
// The verbose logger (verbose.Get) only writes when --verbose is on,
// and stdout/stderr from a detached Windows cleanup child can be
// swallowed by intermediate launchers (run.ps1 wrappers, hidden
// process attrs, etc.). To make these failures forensically
// recoverable, every Phase 3 lifecycle event also goes to a small,
// always-on log file under the same temp directory used for the
// handoff copy and update script:
//
//	<TMP>/gitmap-update-handoff-YYYYMMDD.log
//
// The file is opened in append mode (O_APPEND|O_CREATE|O_WRONLY)
// with line-oriented entries:
//
//	2026-04-24T12:34:56Z pid=12345 ppid=12000 phase=phase-3 event=resolve source=config target=C:\bin\gitmap.exe
//
// We never rotate — daily filename is enough to keep the file
// bounded for the typical update cadence. If the file cannot be
// opened (read-only volume, etc.) writes degrade silently; this
// logger must NEVER block or fail the update flow.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var handoffLogMu sync.Mutex

// handoffLogPath returns the absolute path to today's handoff log file.
// Exposed (lowercase but package-internal) so the --debug-windows dump
// can surface the path in its output.
func handoffLogPath() string {
	day := time.Now().UTC().Format("20060102")
	name := fmt.Sprintf(constants.UpdateHandoffLogNameFmt, day)

	return filepath.Join(os.TempDir(), name)
}

// logHandoffEvent appends one line to the handoff log. `event` is a
// short snake_case identifier (resolve, inline, start_ok, start_fail,
// target_missing, cleanup_start, cleanup_done, cleanup_delay_invalid).
// fields is a flat key=value map appended after the standard prefix.
//
// All logger errors are swallowed — this writer must never disturb
// the update flow. Concurrent calls from the dispatcher and the
// (potentially separate) cleanup child are serialized within the
// process via handoffLogMu; cross-process appends rely on the OS
// O_APPEND atomicity guarantee for small writes (Windows does honor
// this for sub-PIPE_BUF-sized writes when O_APPEND is used).
func logHandoffEvent(phase, event string, fields map[string]string) {
	line := formatHandoffLogLine(phase, event, fields)
	handoffLogMu.Lock()
	defer handoffLogMu.Unlock()

	path := handoffLogPath()
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(line)
}

// formatHandoffLogLine renders the single-line entry. Split out so
// tests can assert the line shape without touching disk.
func formatHandoffLogLine(phase, event string,
	fields map[string]string) string {
	var b strings.Builder
	b.WriteString(time.Now().UTC().Format(time.RFC3339))
	fmt.Fprintf(&b, " pid=%d ppid=%d goos=%s phase=%s event=%s",
		os.Getpid(), os.Getppid(), runtime.GOOS, phase, event)
	for _, k := range sortedStringMapKeys(fields) {
		fmt.Fprintf(&b, " %s=%s", k, escapeHandoffLogValue(fields[k]))
	}
	b.WriteByte('\n')

	return b.String()
}

// sortedStringMapKeys returns the field keys in stable order so log lines from
// different processes diff cleanly.
func sortedStringMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Tiny n — insertion sort keeps the helper allocation-free without
	// dragging in sort.Strings for a hot, no-fail logger.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}

	return keys
}

// escapeHandoffLogValue quotes values that contain spaces or equals
// signs so the line stays unambiguously parseable by `awk` / `grep`.
func escapeHandoffLogValue(v string) string {
	if !strings.ContainsAny(v, " \t=\"") {
		return v
	}
	escaped := strings.ReplaceAll(v, `"`, `\"`)

	return `"` + escaped + `"`
}
