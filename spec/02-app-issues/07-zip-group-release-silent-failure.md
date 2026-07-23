# Issue 07: Zip Group Not Created or Uploaded During Release

## Symptom

Running `gitmap-v27 release v1.61.1 --zip-group "chrome extension-v2"` completes the release (branch, tag, push, metadata) but produces **zero** zip-related output — no success messages, no error messages. The zip archive is neither created locally nor uploaded to GitHub.

## Evidence

- The `.gitmap/zip-groups.json` file was created in the auto-commit, proving the zip group **exists** in the DB.
- Between `✓ Pushed branch and tag to origin` and `Release v1.61.1 complete.`, there are no zip or asset messages at all.
- No `GITHUB_TOKEN`-related error was printed, suggesting `uploadToGitHub` had no assets to report on.

## Root Cause Analysis

### Cause 1 (Primary): `EnsureStagingDir()` silent failure

In `release/workflowfinalize.go` lines 365–368:

```go
stagingDir, err := EnsureStagingDir()
if err != nil {
    return nil  // ← SILENT: no error message printed
}
```

If `EnsureStagingDir()` fails (permissions, disk, path issue), `buildZipGroupAssets` returns `nil` without printing any diagnostic. This is the **only code path** that explains zero output — every other failure path in `BuildZipGroupArchives` → `buildOneZipGroup` prints to either stdout or stderr.

### Cause 2 (Contributing): No `GITHUB_TOKEN`

Even if the zip archive were successfully created in the staging directory, `uploadToGitHub()` requires `GITHUB_TOKEN` to create a GitHub release and upload assets. Without it:
- If `assets` is empty → returns silently (no error printed)
- If `assets` is non-empty → prints `ErrAssetNoToken` to stderr

The release command creates local archives but has no mechanism to attach them to GitHub releases without the token.

### Cause 3 (Possible): DB context mismatch

`buildZipGroupAssets` opens the DB via `store.OpenDefault()`, which resolves the DB path relative to the **binary's physical location** (not CWD). If the zip group was created using a different gitmap-v27 binary instance or the binary was moved/updated between `z create` and `release`, the DB queried during release may not contain the group. The error would go to stderr via `ErrZGCompress`.

## Impact

- Zip group archives are silently skipped during release with no user feedback.
- Users cannot diagnose whether the issue is DB lookup, staging dir, or upload.

## Proposed Fixes

1. **Add error logging to `EnsureStagingDir()` failure** in `buildZipGroupAssets`:
   ```go
   stagingDir, err := EnsureStagingDir()
   if err != nil {
       fmt.Fprintf(os.Stderr, "  ✗ Staging dir failed: %v\n", err)
       return nil
   }
   ```

2. **Add summary output for zip group processing** — print a "Processing N zip group(s)..." message before attempting, so the user knows the flag was recognized.

3. **Print explicit warning when zip groups are requested but no `GITHUB_TOKEN` is set** — currently the token check only warns when assets are non-empty, but users expect the zip to appear on the release.

4. **Add `--verbose` output** showing the resolved DB path and group lookup results during release.

## Files Involved

| File | Role |
|------|------|
| `release/workflowfinalize.go` | `buildZipGroupAssets` — silent failure on staging dir |
| `release/ziparchive.go` | `BuildZipGroupArchives` / `buildOneZipGroup` — archive creation |
| `release/assets.go` | `EnsureStagingDir` — staging directory creation |
| `store/zipgroup.go` | `ListZipGroupItems` — DB item retrieval |
| `store/location.go` | `OpenDefault` — DB path resolution |

## Status

Fixed — added error logging for `EnsureStagingDir` failure, "Processing N zip group(s)..." message before build, and "No zip archives were produced" warning when groups yield zero archives.
