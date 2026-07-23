package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// resolveScanTarget converts a user-supplied scan directory argument into
// a clean, absolute path. It explicitly supports:
//
//	"."          → current working directory
//	".."         → parent directory
//	"../.."      → grandparent directory
//	"../../x"    → grandparent's "x" subdirectory
//	"~/foo"      → expands to $HOME/foo (or %USERPROFILE%/foo on Windows)
//	"./sub"      → CWD/sub
//	absolute     → returned unchanged after Clean
//
// The returned absolute path is validated to exist and to be a directory.
// On failure the process exits with a clear, actionable message — we never
// silently fall back to CWD when the user typed a path that does not exist.
func resolveScanTarget(raw string) string {
	original := raw
	expanded := expandHome(strings.TrimSpace(raw))
	if expanded == "" {
		expanded = constants.DefaultDir
	}

	abs, err := filepath.Abs(expanded)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScanFailed, original, err)
		os.Exit(1)
	}
	abs = filepath.Clean(abs)

	info, err := os.Stat(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScanDirNotFound, original, abs)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, constants.ErrScanDirNotDir, original, abs)
		os.Exit(1)
	}

	if shouldAnnounceResolve(original, abs) {
		fmt.Fprintf(os.Stderr, constants.MsgScanResolvedDir, original, abs)
	}

	return abs
}

// expandHome rewrites a leading "~" segment to the user's home directory.
// We only expand the literal "~" or "~/..." form — "~user" is intentionally
// not supported because Go has no portable resolver for it on Windows.
func expandHome(p string) string {
	if p != "~" && !strings.HasPrefix(p, "~/") && !strings.HasPrefix(p, `~\`) {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return p
	}
	if p == "~" {
		return home
	}

	return filepath.Join(home, p[2:])
}

// shouldAnnounceResolve decides whether to print the "resolved" hint. We
// only print it when the user typed something that materially differs from
// the absolute target (relative segments, "~", trailing dots) — printing
// it for an already-absolute path would be noise.
func shouldAnnounceResolve(original, abs string) bool {
	if original == "" || original == abs {
		return false
	}
	if filepath.IsAbs(original) {
		return false
	}

	return true
}

// resolveRelativeRoot validates the user-supplied --relative-root and
// returns an absolute, cleaned path suitable for passing to mapper.
// Empty input returns "" — the mapper treats that as "fall back to the
// per-repo RelativePath the scanner already computed against the scan
// dir". A non-existent or non-directory path is fatal: silently falling
// back would yield surprising output paths and break the whole point of
// the flag (cwd-stable artifacts).
func resolveRelativeRoot(raw, scanDir string, quiet bool) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	expanded := expandHome(trimmed)
	abs, err := filepath.Abs(expanded)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScanRelativeRootInvalid, raw, err)
		os.Exit(1)
	}
	abs = filepath.Clean(abs)
	info, statErr := os.Stat(abs)
	if statErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScanRelativeRootInvalid, raw, statErr)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, constants.ErrScanRelativeRootInvalid, raw,
			fmt.Errorf("not a directory"))
		os.Exit(1)
	}
	if !quiet && abs != scanDir {
		fmt.Fprintf(os.Stderr, constants.MsgScanRelativeRoot, abs)
	}

	return abs
}
