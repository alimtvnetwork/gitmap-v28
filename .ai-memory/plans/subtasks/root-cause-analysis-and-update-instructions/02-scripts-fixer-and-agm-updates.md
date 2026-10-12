# Subtask 02: Scripts-Fixer & AGM Update Code Synchronization

- **Subtask ID:** `SUB-02`
- **Parent Plan:** `.ai-memory/plans/root-cause-analysis-and-update-instructions.md`
- **Status:** `COMPLETED`
- **Target Files:**
  - `../scripts-fixer/scripts/os/ubuntu/install-antigravity.sh`
  - `../scripts-fixer/scripts-linux/69-install-antigravity/config.json`
  - `../antigravity-manager/scripts/update-antigravity-ide.sh`

---

## Technical Tasks

1. **Update `scripts-fixer` Download Functions:**
   - In `scripts/os/ubuntu/install-antigravity.sh`:
     - Update `fetch_ide_url()` to perform curl against `https://antigravity-hub-auto-updater-974169037036.us-central1.run.app/manifest/latest-x64-linux.yml` with 3s timeout.
     - Fallback default to `2.22.0-5446056071266304`.
   - In `scripts-linux/69-install-antigravity/config.json`:
     - Update `"tarball"` field to `2.22.0-5446056071266304`.

2. **Update `antigravity-manager` Update Script:**
   - In `antigravity-manager/scripts/update-antigravity-ide.sh`:
     - Update version reference to `2.22.0-5446056071266304`.

3. **Verify Shell Scripts:**
   - Run `bash -n` on all updated shell scripts to ensure zero syntax errors.
