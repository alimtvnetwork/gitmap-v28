package errreport

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// reportDirRel is the on-disk location for emitted reports, relative
// to the binary's directory. Mirrors the existing `.gitmap/output/`
// convention used by scan artifacts so users find both kinds of
// generated files in the same parent.
const reportDirRel = ".gitmap/reports"

// reportFilePrefix is the per-report base name; the writer appends a
// unix timestamp + `.json` so concurrent invocations never collide.
const reportFilePrefix = "errors-"

// WriteIfAny writes the JSON report to `<binaryDir>/.gitmap/reports/
// errors-<unix-ts>.json` when the collector holds at least one
// failure. Returns the absolute path on success, "" if nothing was
// written (clean run OR nil receiver), and an error on actual I/O
// failure (which the caller should log to stderr but NOT exit on —
// the report is auxiliary, not load-bearing).
//
// `binaryDir` is supplied by the caller (resolved via
// cmd.resolveBinaryDir or equivalent) so this package stays free of
// the os.Executable / EvalSymlinks dance and remains trivially
// testable with a t.TempDir.
func (c *Collector) WriteIfAny(binaryDir string) (string, error) {
	if c == nil {
		return "", nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.scanEnts) == 0 && len(c.cloneEnts) == 0 {
		return "", nil
	}

	payload := c.buildPayloadLocked()
	dir := filepath.Join(binaryDir, reportDirRel)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create report dir: %w", err)
	}
	path := filepath.Join(dir, fmt.Sprintf("%s%d.json", reportFilePrefix, time.Now().Unix()))

	return path, writeJSONAtomic(path, payload)
}

// buildPayloadLocked snapshots the in-memory state into the on-disk
// shape. Caller must hold c.mu — both reads are serialized with Add
// so the totals match the slice lengths exactly (no torn updates).
func (c *Collector) buildPayloadLocked() fileShape {
	return fileShape{
		Meta: Meta{
			Version:    c.version,
			StartedAt:  c.startedAt.UnixMilli(),
			EndedAt:    time.Now().UnixMilli(),
			Command:    c.command,
			TotalScan:  len(c.scanEnts),
			TotalClone: len(c.cloneEnts),
		},
		Scan:  append([]Entry(nil), c.scanEnts...),
		Clone: append([]Entry(nil), c.cloneEnts...),
	}
}

// writeJSONAtomic writes `payload` to `path` via a temp-file +
// rename so a crash mid-write never leaves a half-formed report.
// Indented for human readability; reports are tiny so the size cost
// is negligible.
func writeJSONAtomic(path string, payload fileShape) error {
	tmp := path + ".tmp"
	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}
	if err := os.WriteFile(tmp, body, 0o644); err != nil {
		return fmt.Errorf("write temp report: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)

		return fmt.Errorf("rename report: %w", err)
	}

	return nil
}
