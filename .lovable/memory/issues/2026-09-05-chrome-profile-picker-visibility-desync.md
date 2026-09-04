# Root Cause Analysis: Chrome Profile Picker Visibility & Local State Desynchronization

**Date:** 2026-09-05
**Status:** Documented / Remediating
**Component:** `gitmap chrome profile import` / `gitmap chrome-profile-copy`
**Affected Artifacts:** `.lovable/assets/screenshots/19-chrome-profile-import-terminal.png`, `.lovable/assets/screenshots/20-chrome-profile-picker-missing-profile5.png`

---

## 1. Problem Description

During invocation of:
```bash
gitmap chrome profile import "erfan.office.n@gmail.com"
```
The CLI reported successful allocation and import into `Profile 5` (`"Default"`), creating the folder `%LOCALAPPDATA%\Google\Chrome\User Data\Profile 5` with `Bookmarks`, `Preferences`, and extension manifests.

However, when Google Chrome was launched or focused, `Profile 5` was completely absent from the Chrome Profile Picker interface (`chrome://settings/manageProfile`). Only pre-existing profiles (`Default`, `Profile 1`, `Profile 2`, `Profile 3`, `Profile 4`, `Profile 6`) were rendered.

---

## 2. Root Cause Analysis

### Cause 1: Chrome In-Memory Cache Overwrite (Concurrency Collision)
Google Chrome maintains an in-memory representation of `%LOCALAPPDATA%\Google\Chrome\User Data\Local State`. When Chrome processes (`chrome.exe`) are running concurrently with `gitmap`:
- GitMap writes the new profile key to `Local State` on disk.
- Chrome's internal background worker flushes its in-memory state back to disk (on timer or event), which does NOT contain `Profile 5`.
- Consequently, `Profile 5` is wiped from `profile.info_cache` and `profile.profiles_order`.

### Cause 2: Incomplete Chromium Attribute Schema
Chromium's profile picker parser validates multiple attributes before declaring a profile valid for UI rendering:
- `avatar_icon`
- `default_avatar_fill_color` & `default_avatar_stroke_color`
- `profile_highlight_color`
- `active_time`
- `shortcut_name`
- `is_using_default_avatar`
The existing `registerImportedProfileToLocalState` only set `name`, `user_name`, `is_using_default_name`, and `is_ephemeral`. Missing fields caused Chrome to discard the profile during tile layout construction.

### Cause 3: Unsanitized Profile Preferences
The restored `Profile 5/Preferences` file contained foreign account blobs and `name: "Person 1"`, conflicting with the allocated display name and triggering Chrome's quarantine or profile merge logic.

---

## 3. Corrective Measures

1. **Process Guard:** Integrate `isChromeRunning()` check before modifying `Local State`. Warn the user and provide actionable guidance.
2. **Schema Invariant:** Ensure `buildChromeDestinationInfoEntry` and `registerImportedProfileToLocalState` populate all 13 mandatory Chromium attributes with valid fallback values.
3. **Preferences Sanitization:** Synchronize `profile.name` in `<ProfileDir>/Preferences` and scrub obsolete `custodian_*` and foreign `gaia_*` fields.
4. **Orphan Reconciliation Command:** Implement `gitmap chrome profile reconcile` to inspect on-disk profile folders and re-register any unlinked directories into `Local State`.
