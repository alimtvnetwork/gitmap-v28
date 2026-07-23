# gitmap update-cleanup

Remove leftover temporary binaries and `.old` backup files from previous updates.

## Alias

None

## Usage

    gitmap update-cleanup

## Flags

None.

## Prerequisites

- None (safe to run at any time)

## Examples

### Example 1: Clean up after a successful update

    gitmap update-cleanup

**Output:**

    ■ Cleaning up update artifacts...
    ✓ Removed temp copy: gitmap-update-tmp-20260319.exe
    ✓ Removed old backup: gitmap.exe.old
    ✓ Cleaned up 2 leftover files

### Example 2: Nothing to clean

    gitmap update-cleanup

**Output:**

    ■ Cleaning up update artifacts...
    ✓ No leftover update files found

### Example 3: Multiple temp copies from interrupted updates

    gitmap update-cleanup

**Output:**

    ■ Cleaning up update artifacts...
    ✓ Removed temp copy: gitmap-update-tmp-20260318.exe
    ✓ Removed temp copy: gitmap-update-tmp-20260319.exe
    ✓ Removed old backup: gitmap.exe.old
    ✓ Removed old backup: gitmap.old
    ✓ Cleaned up 4 leftover files

### Example 4: Run after a failed update left artifacts behind

    gitmap update
    # ✗ Build failed — leftover temp binary remains

    gitmap update-cleanup

**Output:**

    ■ Cleaning up update artifacts...
    ✓ Removed temp copy: gitmap-update-tmp-20260319.exe
    ✓ Cleaned up 1 leftover file
    → Safe to retry: gitmap update

## See Also

- [update](update.md) — Self-update gitmap to the latest version
- [version](version.md) — Check current installed version
- [doctor](doctor.md) — Diagnose installation issues

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter update-cleanup
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
