# CI/CD Issue 10: `vscodepm` Cross-Platform Path Normalization

- **Stage**: CI `go test` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm`, `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

`TestUpdateRootPathAndRemoveEntry` and `TestPathsEqualCrossPlatformSlashes` failed exclusively on Linux/macOS CI runners while passing locally on Windows:
```
--- FAIL: TestUpdateRootPathAndRemoveEntry (0.00s)
update_remove_test.go:23: UpdateRootPathAt failed: project not found with rootPath D:/work/my-app
--- FAIL: TestPathsEqualCrossPlatformSlashes (0.00s)
update_remove_test.go:63: pathsEqual("D:\work\my-app", "D:/work/my-app") = false, want true
```

## 2. Root Cause Analysis

- The previous fix used `filepath.ToSlash(p)` followed by `filepath.Clean(p)`.
- `filepath.ToSlash` is OS-aware: it only replaces characters matching `filepath.Separator`.
- On Linux and macOS, `filepath.Separator` is `/`. Therefore, `filepath.ToSlash` does absolutely nothing to backslashes (`\`), as backslashes are technically valid filename characters in POSIX systems.
- As a result, `"D:\work\my-app"` remained unchanged and evaluated as a distinct string from `"D:/work/my-app"`.

## 3. Solution

- Replaced `filepath.ToSlash(p)` with an explicit `strings.ReplaceAll(p, "\\", "/")` before passing to `filepath.Clean()`.
- This forces backslashes in simulated Windows paths (which are common in config files or user-provided extra paths) to be consistently normalized to forward slashes regardless of the host OS's native separator.
- Applied this fix to both `vscodepm.normalizePath` and `cmd.pathKey`.

## 4. What NOT to Repeat

- **Never rely on `filepath.ToSlash` for cross-platform backslash normalization:** If you are dealing with paths that might contain backslashes (like user configs) on a POSIX host, `filepath.ToSlash` will silently skip them. Use `strings.ReplaceAll(p, "\\", "/")` instead.
