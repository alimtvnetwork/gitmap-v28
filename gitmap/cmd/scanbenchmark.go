// Package cmd — scanbenchmark.go captures per-phase scan timings and writes
// them to a benchmark log so users can diagnose slow scans without us
// having to ask. Each scan invocation appends a fresh, timestamped block
// to .gitmap/output/scan-benchmark.log alongside the binary version.
//
// Why a file (not just stdout): users routinely report "scan is slow" with
// no reproducible numbers. The log gives them — and us — a record of
// exactly which phase ate the wall clock, across every run.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// scanBenchmarkFile is the fixed log file name appended in the scan
// output directory. A single file (not per-run) keeps history available
// for trend comparisons; rotating is the user's choice.
const scanBenchmarkFile = "scan-benchmark.log"

// benchPhase holds one labeled timing measurement.
type benchPhase struct {
	name     string
	duration time.Duration
}

// scanBenchmark accumulates per-phase timings during a scan run. It is
// safe for concurrent use; phases are appended under a mutex.
type scanBenchmark struct {
	mu        sync.Mutex
	startedAt time.Time
	scanDir   string
	phases    []benchPhase
}

// newScanBenchmark stamps the start time and the scan directory.
func newScanBenchmark(scanDir string) *scanBenchmark {
	return &scanBenchmark{
		startedAt: time.Now(),
		scanDir:   scanDir,
	}
}

// Phase runs fn, measures its wall-clock duration, and records it. The
// returned value is the duration so callers can short-circuit logging.
func (b *scanBenchmark) Phase(name string, fn func()) time.Duration {
	start := time.Now()
	fn()
	d := time.Since(start)
	b.mu.Lock()
	b.phases = append(b.phases, benchPhase{name: name, duration: d})
	b.mu.Unlock()

	return d
}

// Record appends a pre-measured phase without invoking a closure. Use this
// when the phase boundaries don't fit a simple wrapper (e.g. measuring
// only the inner loop of an existing function).
func (b *scanBenchmark) Record(name string, d time.Duration) {
	b.mu.Lock()
	b.phases = append(b.phases, benchPhase{name: name, duration: d})
	b.mu.Unlock()
}

// WriteLog appends the benchmark block to outputDir/scan-benchmark.log.
// Failures are reported to stderr but never fail the scan — benchmarking
// is observability, not a feature the user explicitly asked to gate on.
func (b *scanBenchmark) WriteLog(outputDir string) {
	if err := os.MkdirAll(outputDir, constants.DirPermission); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not create benchmark dir: %v\n", err)

		return
	}
	path := filepath.Join(outputDir, scanBenchmarkFile)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, constants.FilePermission)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not open benchmark log: %v\n", err)

		return
	}
	defer f.Close()

	b.writeBlock(f)
}

// writeBlock formats and writes one timestamped block to w. Format is
// human-readable but stable enough to grep / diff across runs.
func (b *scanBenchmark) writeBlock(f *os.File) {
	b.mu.Lock()
	defer b.mu.Unlock()

	total := time.Since(b.startedAt)
	fmt.Fprintf(f, "===== gitmap scan v%s @ %s =====\n",
		constants.Version, b.startedAt.Format(time.RFC3339))
	fmt.Fprintf(f, "  scanDir: %s\n", b.scanDir)
	fmt.Fprintf(f, "  goos:    %s/%s  cpus: %d\n",
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU())

	for _, p := range b.phases {
		fmt.Fprintf(f, "  %-28s %10s\n", p.name, formatBenchDuration(p.duration))
	}
	fmt.Fprintf(f, "  %-28s %10s\n", "TOTAL", formatBenchDuration(total))
	fmt.Fprintln(f)
}

// formatBenchDuration prints a duration with consistent precision.
// time.Duration's default formatter switches units, which makes columns
// hard to align in the benchmark log.
func formatBenchDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}

	return fmt.Sprintf("%.2fs", d.Seconds())
}
