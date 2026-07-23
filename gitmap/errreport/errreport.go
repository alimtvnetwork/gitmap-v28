// Package errreport collects per-repo failures during long-running
// gitmap operations (currently `scan` and `clone-next`) and writes a
// single grouped JSON report at command exit.
//
// Design contract — locked via UX questions on 2026-04-26:
//
//   - Scope: ALL recoverable scan + cn failures (scanner ReadDir
//     errors, ls-remote probe failures, cn clone failures, cn
//     non-skip step errors). User-visible "skipped" outcomes
//     (already-exists, alias-match, user-quit) are NOT recorded —
//     only true errors land here.
//   - CLI: bare `--report-errors` boolean; output path is fixed at
//     `.gitmap/reports/errors-<unix-ts>.json` next to the binary.
//   - Timing: atomic write at command exit, ONLY if at least one
//     failure was recorded. Clean runs leave no file behind.
//   - Schema: grouped by phase — `{meta, scan: [...], clone: [...]}`
//     so a single report can cover a chained scan→cn pipeline
//     without ambiguity about which phase produced which entry.
//
// The collector is goroutine-safe so workers in scanner.walkParallel,
// probe.BackgroundRunner, and the cn batch worker pool can all hand
// it failures concurrently without external locking. Callers that
// disable the feature (the common case) construct nil and the
// per-call helpers no-op.
package errreport

import (
	"sync"
	"time"
)

// Phase tags which subsystem produced the failure. Stored as the
// JSON top-level array key (`scan` / `clone`) so consumers can
// branch on it without parsing the inner Entry first.
type Phase string

const (
	// PhaseScan covers everything that runs under `gitmap scan`:
	// scanner ReadDir failures, background ls-remote / shallow-clone
	// probe errors, project-detection failures we choose to surface.
	PhaseScan Phase = "scan"
	// PhaseClone covers everything under `gitmap clone-next` batch
	// mode (--all / --csv): per-repo clone failures, post-clone
	// hooks that errored, and any other non-skip cn outcome.
	PhaseClone Phase = "clone"
)

// Entry is one failure record. RepoPath is required (the whole point
// is per-repo attribution); the rest are best-effort and may be
// empty strings when the producer doesn't have them at hand.
type Entry struct {
	RepoPath        string `json:"repo_path"`
	RemoteURL       string `json:"remote_url,omitempty"`
	Step            string `json:"step,omitempty"`    // free-text producer-defined: "ls-remote", "clone", "readdir", …
	Error           string `json:"error"`             // human-readable error text
	TimestampUnixMS int64  `json:"timestamp_unix_ms"` // capture time, ms-precision
}

// Meta is the report header. Written even when both phase arrays are
// empty (the file itself is only emitted when at least one entry
// exists, but once it IS emitted the meta block is always present).
type Meta struct {
	Version    string `json:"gitmap_version"`
	StartedAt  int64  `json:"started_at_unix_ms"`
	EndedAt    int64  `json:"ended_at_unix_ms"`
	Command    string `json:"command"` // "scan" / "clone-next" / "scan+clone-next"
	TotalScan  int    `json:"total_scan_failures"`
	TotalClone int    `json:"total_clone_failures"`
}

// fileShape is the on-disk JSON layout. Kept private so callers can't
// hand-build malformed reports — the only path to a file is
// (*Collector).WriteIfAny.
type fileShape struct {
	Meta  Meta    `json:"meta"`
	Scan  []Entry `json:"scan"`
	Clone []Entry `json:"clone"`
}

// Collector accumulates failures across goroutines. The zero value is
// NOT usable — callers must use New so the start time is captured at
// construction. A nil *Collector is treated as "feature disabled" by
// every method (Add and WriteIfAny both no-op), so call sites only
// need a single nil-check at construction time, not at every Add.
type Collector struct {
	mu        sync.Mutex
	scanEnts  []Entry
	cloneEnts []Entry
	startedAt time.Time
	version   string
	command   string
}

// New returns a ready-to-use Collector. `version` is embedded in the
// report meta block (caller passes constants.Version); `command`
// describes the invoking subcommand for downstream tooling that
// pivots on it. Pass an empty string for either to omit them.
func New(version, command string) *Collector {
	return &Collector{
		startedAt: time.Now(),
		version:   version,
		command:   command,
	}
}

// Add records one failure under the given phase. Safe to call from
// any goroutine. A nil receiver is a silent no-op so call sites can
// pass an unconditionally-constructed (or unconditionally-nil)
// collector without branching on every error path.
func (c *Collector) Add(phase Phase, e Entry) {
	if c == nil {
		return
	}
	if e.TimestampUnixMS == 0 {
		e.TimestampUnixMS = time.Now().UnixMilli()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	switch phase {
	case PhaseScan:
		c.scanEnts = append(c.scanEnts, e)
	case PhaseClone:
		c.cloneEnts = append(c.cloneEnts, e)
	}
}

// Count returns the current per-phase totals. Useful for the
// end-of-command summary line ("3 failures recorded → <path>").
// Nil receiver returns zeros.
func (c *Collector) Count() (scan, clone int) {
	if c == nil {
		return 0, 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.scanEnts), len(c.cloneEnts)
}
