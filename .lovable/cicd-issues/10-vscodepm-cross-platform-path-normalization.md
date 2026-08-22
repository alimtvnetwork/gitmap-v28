# CI/CD Issue 10: VS Code Project Manager Cross-Platform Slash Normalization

- **Stage**: CI `go test` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm`, `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

`TestUpdateRootPathAndRemoveEntry` in `vscodepm/update_remove_test.go` failed on Linux with `project not found with rootPath D:/work/my-app`.
On POSIX platforms (Linux/macOS), backslashes `\` are valid filename characters rather than directory separators. `vscodepm.normalizePath` used `filepath.Clean(p)` without converting slashes via `filepath.ToSlash(p)`. Consequently, Windows-style paths written into `projects.json` with backslashes (`D:\work\my-app`) failed to match target arguments passed with forward slashes (`D:/work/my-app`).

## 2. Root Cause Analysis

- `filepath.Clean` on POSIX systems preserves backslashes as verbatim characters.
- Path comparison in `vscodepm.pathsEqual` did not normalize backslashes to forward slashes before comparing.
- Windows drive letter paths (`D:/...`) were not case-folded when evaluated on non-Windows hosts.

## 3. Solution

1. **`filepath.ToSlash` Normalization**:
   - Updated `vscodepm.normalizePath` in `vscodepm/io.go` to use `filepath.Clean(filepath.ToSlash(p))`.
   - Added `isWindowsDrivePath` check to ensure Windows drive paths are case-folded across all platforms.
2. **Synchronized in `cmd/code.go`**:
   - Updated `pathKey` in `cmd/code.go` to follow the same `filepath.ToSlash` normalization contract.
3. **Unit Tests Added**:
   - Added `TestPathsEqualCrossPlatformSlashes` in `vscodepm/update_remove_test.go` to verify cross-platform comparison for mixed slashes, casing, and trailing slashes.

## 4. What NOT to Repeat

- **Never rely on `filepath.Clean` alone for cross-platform path equality**: Always convert slashes via `filepath.ToSlash` when parsing configuration or database entries that may store paths formatted on a different operating system.
