# Issue 05: `gitmap-v28 ls` Returns Empty After Scan

## Symptom

Running `gitmap-v28 ls` shows "No repos tracked. Run 'gitmap-v28 scan' first."
even though a scan was previously completed successfully.

## Root Cause

The `OpenDefault()` function in `store/location.go` passes
`BinaryDataDir()` (which already resolves to `<binary-dir>/data/`)
into `ActiveProfileDBFile(dir)`, which internally calls
`profileConfigPath(dir)` → `filepath.Join(dir, "data", "profiles.json")`.

This **double-nests** the `data` segment:
- Expected: `<binary-dir>/data/profiles.json`
- Actual:   `<binary-dir>/data/data/profiles.json`

While this doesn't break the default profile (file-not-found falls back
to `"default"` → `gitmap.db`), it means:

1. Profile switching via `OpenDefault()` is silently broken.
2. If the user scanned with an older binary version (pre-v2.15.1) that
   used CWD-relative paths (`.gitmap/output/data/gitmap.db`), the new
   binary looks at `<binary-dir>/data/gitmap.db` — a completely different
   file.

## Fix

1. **`store/location.go`**: Pass `filepath.Dir(resolved)` (the binary
   directory, WITHOUT `/data`) to `ActiveProfileDBFile()` so profile
   config resolves correctly.
2. **`cmd/list.go`**: Add `--verbose` diagnostic that prints the resolved
   DB path so users can verify which database is being queried.
3. **`store/profile.go`**: `ActiveProfileDBFile` should accept the base
   directory (binary dir), not the data dir, to avoid double-nesting.

## Affected Versions

v2.15.1 through v2.19.0

## Verification

After fix: `gitmap-v28 ls --verbose` should print the resolved DB path.
Re-scanning should persist to the same location that `ls` reads from.
