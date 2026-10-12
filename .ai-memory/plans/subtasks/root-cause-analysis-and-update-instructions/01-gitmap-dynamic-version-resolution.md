# Subtask 01: GitMap Dynamic Version Resolution & Default Constants

- **Subtask ID:** `SUB-01`
- **Parent Plan:** `.ai-memory/plans/root-cause-analysis-and-update-instructions.md`
- **Status:** `COMPLETED`
- **Target Files:**
  - `cli/cmdinstall/installantigravity_types.go`
  - `cli/cmdinstall/installantigravity_fetch.go`
  - `cli/cmdinstall/installantigravity_test.go`

---

## Technical Tasks

1. **Update Default Constants (`installantigravity_types.go`):**
   - Change `AntigravityDefaultVersion` to `"2.22.0"`.
   - Change `AntigravityDefaultBuildID` to `"5446056071266304"`.
   - Add `AntigravityManifestURL` pointing to `https://antigravity-hub-auto-updater-974169037036.us-central1.run.app/manifest/latest-x64-linux.yml`.

2. **Implement Dynamic Manifest Resolution (`installantigravity_fetch.go`):**
   - Create `fetchLiveReleaseMetadata(manifestURL string) (string, string, error)` with a strict 3-second timeout.
   - Parse `version` and `buildID` from manifest lines or URL components.
   - Fall back smoothly to default constants on timeout or network error.
   - Update `buildAntigravityURL` and `getAntigravityDesktopDownloadUrl`.

3. **Update Unit Tests (`installantigravity_test.go`):**
   - Update expected test URLs from `2.13.0` / `2.19.1` to `2.22.0`.
   - Add unit test validating dynamic manifest resolution and graceful fallback.
