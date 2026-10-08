# Chrome Profile & Extension Import/Export E2E Verification Report

- **Date:** 2026-10-08T17:14:00+08:00
- **GitMap Version:** `v6.509.0`
- **Host OS:** Windows Server (Administrator)
- **Status:** **PASS (100% Verified)**

---

## 1. Scope & Objective

Verify the complete lifecycle of Google Chrome profile and extension export and smart-import in GitMap (`gitmap cpe` and `gitmap cpi`):
1. **Sample Selection**: Target non-destructive test profile (`Profile 6`, display name: `"Personal"`, email: `"test@test.com"`).
2. **Extension Staging**: Add test extension (`cjpalhdlnbpafiamejdnhcphjbkeiagm`) under `Extensions/` to test extension detection and hint staging.
3. **Safety Backup**: Create safety backup in `repo-secrets/01-gitmap/02-chrome-import-export-test/backup-profile-6/`.
4. **Export Phase**: Snapshot profile metadata, bookmarks, preferences, and extension IDs to `repo-secrets/01-gitmap/02-chrome-import-export-test/profile-6-export.json` and `.csv`.
5. **Removal Phase**: Delete `Profile 6` from Chrome User Data directory on disk (`C:\Users\Administrator\AppData\Local\Google\Chrome\User Data\Profile 6`).
6. **Smart-Import Phase**: Re-import the profile using `gitmap cpi` targeting `"Profile 6"`.
7. **Post-Import Verification**: Verify directory recreation, file contents (`Bookmarks`, `Preferences`, `gitmap-pending-extensions.txt`), database tracking, and profile discovery via `gitmap cpl`.
8. **Logging**: Full command execution logs saved to `test-execution.log` in this directory.

---

## 2. Test Execution Steps & Results

### Step 1: Pre-Test Inspection & Safety Backup
- Inspected `Profile 6` directory on disk:
  - `Bookmarks` (211 bytes, Site1 -> `https://example.com`)
  - `Preferences` (185 bytes, name: `"Personal"`, email: `"test@test.com"`)
  - Added test extension `Extensions/cjpalhdlnbpafiamejdnhcphjbkeiagm/1.0.0`
- Safety backup copied to `repo-secrets/01-gitmap/02-chrome-import-export-test/backup-profile-6/`.

### Step 2: Export (`gitmap cpe`)
```powershell
gitmap cpe "Profile 6" "repo-secrets\01-gitmap\02-chrome-import-export-test\profile-6-export.json"
```
**Output:**
```text
Artifacts
  json: repo-secrets\01-gitmap\02-chrome-import-export-test\profile-6-export.json
  csv:  repo-secrets\01-gitmap\02-chrome-import-export-test\profile-6-export.csv
✔ db synced  profile Profile 6
```
- Validated JSON payload schema:
  - `schemaVersion`: 1
  - `gitmapVersion`: `6.509.0`
  - `name`: `"Profile 6"`
  - `displayName`: `"Personal"`
  - `email`: `"test@test.com"`
  - `extensionIds`: `["cjpalhdlnbpafiamejdnhcphjbkeiagm"]`
  - `bookmarks` and `preferences` preserved.

### Step 3: Deletion
```powershell
Remove-Item -Path "C:\Users\Administrator\AppData\Local\Google\Chrome\User Data\Profile 6" -Recurse -Force
```
- Verified `Test-Path` returned `False`.

### Step 4: Smart-Import (`gitmap cpi`)
```powershell
gitmap cpi "repo-secrets\01-gitmap\02-chrome-import-export-test\profile-6-export.json" "Profile 6"
```
**Output:**
```text
Starting Chrome Profile Import (1 snapshot(s))

  [Step 1/5] Inspecting snapshot: repo-secrets\01-gitmap\02-chrome-import-export-test\profile-6-export.json
        -> Name: "Profile 6" | Display: "Personal" | Email: "test@test.com" | Bookmarks: 1 | Extensions: 1
      [Step 2/5] Using explicit profile target: "Profile 6"
  [Step 3/5] Restoring Bookmarks, Preferences & Auth into "Profile 6"...
  [Step 4/5] Staging extensions (1 pending hints)...
  [Step 5/5] Registering "Profile 6" in Chrome Local State...
  ✔ Successfully imported repo-secrets\01-gitmap\02-chrome-import-export-test\profile-6-export.json -> Profile 6 ("Personal")

✔ Chrome Profile Import Complete: 1 of 1 profile(s) imported successfully.
```

### Step 5: Post-Import Verification
- `Profile 6` recreated on disk (`Test-Path = True`).
- `Bookmarks` restored with 100% fidelity (`Site1` -> `https://example.com`).
- `Preferences` restored (`name: "Personal"`, `email: "test@test.com"`).
- `gitmap-pending-extensions.txt` restored with `cjpalhdlnbpafiamejdnhcphjbkeiagm`.
- `gitmap cpl` confirms `Profile 6` is registered and active in Chrome Local State and SQLite split-DB.

---

## 3. Verdict
The import/export and extension lifecycle pipeline is fully operational, verified, and logged in `test-execution.log`.
