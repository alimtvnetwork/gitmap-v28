package cmd

// File enumeration + binary-detection helpers for `gitmap fix-repo`.
// Mirrors scripts/fix-repo/File-Scan.ps1: list tracked files via
// `git ls-files`, skip reparse points / oversized / binary-extension
// / NUL-byte-prefixed files.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// fixRepoBinaryExts is the suffix set we treat as opaque assets. The
// values match the PowerShell + Bash scripts so all three engines
// agree on which files are scanned.
var fixRepoBinaryExts = map[string]struct{}{
	".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {}, ".webp": {},
	".ico": {}, ".pdf": {}, ".zip": {}, ".tar": {}, ".gz": {},
	".tgz": {}, ".bz2": {}, ".xz": {}, ".7z": {}, ".rar": {},
	".woff": {}, ".woff2": {}, ".ttf": {}, ".otf": {}, ".eot": {},
	".mp3": {}, ".mp4": {}, ".mov": {}, ".wav": {}, ".ogg": {},
	".webm": {}, ".class": {}, ".jar": {}, ".so": {}, ".dylib": {},
	".dll": {}, ".exe": {}, ".pyc": {},
}

// fixRepoSweepResult aggregates one full sweep's counts.
type fixRepoSweepResult struct {
	scanned      int
	changed      int
	replacements int
	failed       bool
	goFiles      []string // absolute paths of modified .go files (for gofmt)
	// backup is the snapshot session opened lazily on first rewrite
	// (v5.40.0+). Finalized by runFixRepo after the sweep returns so
	// `gitmap undo` can restore pre-rewrite copies.
	backup *fixRepoBackupSession
}

// runFixRepoSweep enumerates tracked files and rewrites each.
func runFixRepoSweep(identity fixRepoIdentity, targets []int, opts fixRepoOptions) fixRepoSweepResult {
	files := listTrackedFiles(identity.root)
	result := fixRepoSweepResult{}
	if !opts.isDryRun {
		result.backup = newFixRepoBackupSession(identity)
	}
	for _, rel := range files {
		processFixRepoFile(rel, identity, targets, opts, &result)
	}

	return result
}

// processFixRepoFile is the per-file branch extracted from the sweep
// loop so runFixRepoSweep stays under the 15-line cap. Dry-run takes
// a separate path that emits a per-rule breakdown without writing.
func processFixRepoFile(rel string, identity fixRepoIdentity, targets []int,
	opts fixRepoOptions, result *fixRepoSweepResult,
) {
	full := filepath.Join(identity.root, rel)
	if isFixRepoIgnoredPath(rel) || !isFixRepoScannable(full) {
		return
	}
	result.scanned++
	if opts.isDryRun {
		previewOneFile(full, rel, identity, targets, opts, result)

		return
	}
	rewriteOneFile(full, rel, identity, targets, opts, result)
}

// rewriteOneFile is the write-path branch. Mutates disk and records
// the modified .go path for the gofmt + strict post-steps. Pre-write
// backup (v5.40.0+) is taken only when reps > 0 so untouched files
// never appear in the snapshot.
func rewriteOneFile(full, rel string, identity fixRepoIdentity, targets []int,
	opts fixRepoOptions, result *fixRepoSweepResult,
) {
	raw, err := os.ReadFile(full)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrWriteFmt, rel, err)
		result.failed = true

		return
	}
	updated, reps := applyAllTargetsR(string(raw), identity.base,
		identity.current, targets, opts.restrictNoVersion)
	if reps == 0 {
		return
	}
	persistRewrittenFile(full, rel, updated, reps, opts, result)
}

// persistRewrittenFile backs up the original, writes the new bytes,
// and records gofmt + verbose bookkeeping. Split out so rewriteOneFile
// stays under the 15-line cap.
func persistRewrittenFile(full, rel, updated string, reps int,
	opts fixRepoOptions, result *fixRepoSweepResult,
) {
	if result.backup != nil {
		result.backup.BackupFile(rel)
	}
	if err := os.WriteFile(full, []byte(updated), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrWriteFmt, rel, err)
		result.failed = true

		return
	}
	result.changed++
	result.replacements += reps
	if isGoSourceFile(rel) {
		result.goFiles = append(result.goFiles, full)
	}
	if opts.isVerbose {
		fmt.Printf(constants.FixRepoMsgModified, rel, reps)
	}
}

// previewOneFile is the dry-run branch. Always prints a per-file
// `[dry-run]` line plus a per-rule breakdown so the user sees every
// would-be substitution without touching disk.
func previewOneFile(full, rel string, identity fixRepoIdentity, targets []int,
	opts fixRepoOptions, result *fixRepoSweepResult,
) {
	reps, hits, err := previewFixRepoFile(full, identity.base, identity.current, targets, opts.restrictNoVersion)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrWriteFmt, rel, err)
		result.failed = true

		return
	}
	if reps == 0 {
		return
	}
	result.changed++
	result.replacements += reps
	fmt.Printf(constants.FixRepoMsgDryRunPreview, rel, reps, formatFixRepoHits(hits))
}

// formatFixRepoHits renders the per-rule breakdown as a compact
// comma-joined list, e.g. `v1×3, v2×1, bare×2`. Empty slice yields
// an empty string (caller still prints the count).
func formatFixRepoHits(hits []fixRepoTargetHit) string {
	if len(hits) == 0 {
		return ""
	}
	parts := make([]string, 0, len(hits))
	for _, h := range hits {
		parts = append(parts, formatOneFixRepoHit(h))
	}

	return strings.Join(parts, ", ")
}

// formatOneFixRepoHit renders a single rule hit. The bare-base
// sentinel (n == -1) is rendered as `bare` to distinguish it from
// the numbered `{base}-vN` rules in the dry-run breakdown.
func formatOneFixRepoHit(h fixRepoTargetHit) string {
	if h.n == fixRepoBareBaseSentinel {
		return fmt.Sprintf(constants.FixRepoMsgDryRunHitBare, h.count)
	}

	return fmt.Sprintf(constants.FixRepoMsgDryRunHit, h.n, h.count)
}

// listTrackedFiles runs `git ls-files` in repoRoot. Failures yield
// an empty list (the caller logs nothing because git already wrote
// to stderr) and an empty sweep is the natural no-op.
func listTrackedFiles(repoRoot string) []string {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = repoRoot
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n")
	files := make([]string, 0, len(lines))
	for _, l := range lines {
		if l != "" {
			files = append(files, l)
		}
	}

	return files
}

// isFixRepoScannable composes the per-file skip checks: reparse,
// oversize, binary-extension, NUL-byte prefix.
func isFixRepoScannable(fullPath string) bool {
	if isFixRepoSkippablePath(fullPath) {
		return false
	}
	if isFixRepoBinaryExt(fullPath) {
		return false
	}
	if hasFixRepoNullByte(fullPath) {
		return false
	}

	return true
}

// isFixRepoSkippablePath flags reparse points and >5 MiB files.
func isFixRepoSkippablePath(fullPath string) bool {
	info, err := os.Lstat(fullPath)
	if err != nil {
		return true
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	if info.Size() > constants.FixRepoMaxFileBytes {
		return true
	}

	return false
}

// isFixRepoBinaryExt reports whether fullPath has a known binary extension.
func isFixRepoBinaryExt(fullPath string) bool {
	ext := strings.ToLower(filepath.Ext(fullPath))
	_, ok := fixRepoBinaryExts[ext]

	return ok
}

// hasFixRepoNullByte checks the first 8 KiB for a NUL byte. A NUL
// in early bytes is the standard "this is binary" heuristic used by
// git, grep, etc., and matches the PowerShell script's behavior.
func hasFixRepoNullByte(fullPath string) bool {
	f, err := os.Open(fullPath)
	if err != nil {
		return true
	}
	defer func() { _ = f.Close() }()
	buf := make([]byte, constants.FixRepoBinarySniffMax)
	n, _ := f.Read(buf)
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}
