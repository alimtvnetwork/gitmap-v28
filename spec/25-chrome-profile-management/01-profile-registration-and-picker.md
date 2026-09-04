# 01 — Chrome Profile Registration & Picker Visibility Specification

**Version:** 1.0.0
**Updated:** 2026-09-05
**Status:** Canonical
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Context & Problem Statement

Google Chrome and Chromium-based browsers maintain profile metadata across two distinct tiers:
1. **User Data Root (`Local State`)**: A central JSON file located at `%LOCALAPPDATA%\Google\Chrome\User Data\Local State` (Windows), `~/Library/Application Support/Google/Chrome/Local State` (macOS), or `~/.config/google-chrome/Local State` (Linux). This file stores the master registry of all installed profiles (`profile.info_cache`) and their UI display order (`profile.profiles_order`).
2. **Profile Subdirectories (`<ProfileDir>/Preferences`)**: Individual JSON files located inside each profile folder (e.g. `User Data/Default/Preferences`, `User Data/Profile 1/Preferences`, `User Data/Profile 5/Preferences`), storing local settings, extensions, bookmarks, and local identity tokens.

When an external CLI tool creates or restores a profile on disk without full registration in `Local State`, or writes to `Local State` while Chrome processes are actively running, **Chrome fails to display the profile in its Profile Picker UI (`chrome://settings/manageProfile`)**.

---

## 2. Root Cause Analysis (RCA)

Three root causes prevent imported or copied profiles from appearing in Chrome's Profile Picker:

### RCA-1: In-Memory Cache Overwrite by Active Chrome Processes
- **Mechanism:** Google Chrome maintains an in-memory image of `Local State` while any browser process (`chrome.exe`) is running. Chrome periodically flushes this in-memory cache to disk on background timers, window focus events, or graceful exit.
- **Failure Mode:** If `gitmap` updates `Local State` while Chrome is running, Chrome will shortly overwrite `Local State` on disk with its in-memory version, completely stripping any newly added profile keys from `info_cache` and `profiles_order`.
- **Mandatory Guard:** GitMap must check `isChromeRunning()` during import/copy operations. When Chrome is open, GitMap must either warn the user with actionable instructions or provide a safe reconciliation command once Chrome closes.

### RCA-2: Incomplete Attribute Schema in `info_cache`
- **Mechanism:** Modern Chromium profile pickers require a complete set of metadata fields to construct the profile picker tile and avatar badge.
- **Failure Mode:** Storing only `name` and `user_name` causes Chrome's profile picker parser to discard or bypass the entry during startup validation.
- **Mandatory Schema:** Every entry in `profile.info_cache[<dir>]` must supply:
  - `name`: String display name.
  - `shortcut_name`: String shortcut label (matches display name).
  - `user_name`: String email account or empty string.
  - `active_time`: Float64 epoch timestamp (e.g. `time.Now().Unix()`).
  - `avatar_icon`: String avatar theme identifier (`"chrome://theme/IDR_PROFILE_AVATAR_26"`).
  - `default_avatar_fill_color`: Int32 color value (e.g. `-13625057`).
  - `default_avatar_stroke_color`: Int32 stroke color value (e.g. `-1786428`).
  - `profile_highlight_color`: Int32 highlight color value (e.g. `-13625057`).
  - `is_using_default_avatar`: Boolean (`true`).
  - `is_using_default_name`: Boolean (`false`).
  - `is_ephemeral`: Boolean (`false`).
  - `is_consented_primary_account`: Boolean (`false`).
  - `signin.with_credential_provider`: Boolean (`false`).

### RCA-3: Per-Profile `Preferences` Desynchronization
- **Mechanism:** Chrome verifies profile integrity on startup by inspecting `<ProfileDir>/Preferences`.
- **Failure Mode:** If `Preferences` contains conflicting identity blobs from another machine (e.g. `custodian_email`, stale `gaia_id`, or a divergent `profile.name`), Chrome attempts to merge or quarantine the profile.
- **Mandatory Guard:** GitMap must patch `<ProfileDir>/Preferences` to set `profile.name` equal to the assigned display name and sanitize all conflicting GAIA account identifiers.

---

## 3. Local State Specification & Contracts

### 3.1 `profile.info_cache` Entry Schema

```json
{
  "active_time": 1788542727.0,
  "avatar_icon": "chrome://theme/IDR_PROFILE_AVATAR_26",
  "background_apps": false,
  "default_avatar_fill_color": -13625057,
  "default_avatar_stroke_color": -1786428,
  "enterprise_label": "",
  "force_signin_profile_locked": false,
  "gaia_given_name": "",
  "gaia_id": "",
  "gaia_name": "",
  "hosted_domain": "",
  "is_consented_primary_account": false,
  "is_ephemeral": false,
  "is_glic_eligible": false,
  "is_managed": 0,
  "is_using_default_avatar": true,
  "is_using_default_name": false,
  "managed_user_id": "",
  "metrics_bucket_index": 7,
  "name": "erfan.office.n",
  "profile_color_seed": -4385188,
  "profile_highlight_color": -13625057,
  "shortcut_name": "erfan.office.n",
  "signin.with_credential_provider": false,
  "user_name": "erfan.office.n@gmail.com"
}
```

### 3.2 `profile.profiles_order` Array
- Every registered profile directory name (e.g. `"Default"`, `"Profile 1"`, `"Profile 5"`) must be appended to the `profile.profiles_order` array in `Local State`.
- Duplicate directory names in `profiles_order` are strictly forbidden.
- Registration must preserve existing entries and append new directories at the tail.

---

## 4. Concurrency Guard & Process Life Cycle

```
[ gitmap chrome profile import ]
             │
             ▼
    Check isChromeRunning()
       ├── If Chrome is RUNNING:
       │     ├─ Issue actionable warning:
       │     │  "⚠ Chrome is currently running. Chrome caches Local State in memory."
       │     │  "To ensure your new profile is recognized in the profile picker,"
       │     │  "please close Chrome or run 'gitmap chrome profile reconcile'."
       │     └─ Proceed with disk copy and atomic Local State write.
       └── If Chrome is CLOSED:
             ├─ Perform atomic Local State write.
             └─ Confirm: "✓ Profile registered and ready for Chrome launch."
```

---

## 5. Automatic & Manual Profile Reconciliation

To handle cases where a profile folder exists on disk (e.g. `Profile 5`) but is missing from `Local State`:
1. **CLI Command:** `gitmap chrome profile reconcile` (alias `gitmap chrome profile repair`).
2. **Behavior:**
   - Scans all profile directories under Chrome User Data (`Default`, `Profile *`).
   - Compares on-disk directories against `profile.info_cache` in `Local State`.
   - For every orphan profile found on disk:
     - Extracts email and display name from `<ProfileDir>/Preferences`.
     - Synthesizes a valid `info_cache` entry.
     - Appends the directory to `profiles_order`.
     - Patches `<ProfileDir>/Preferences` to ensure consistent `profile.name`.
   - Writes `Local State` atomically.
   - Outputs a clean summary report of reconciled profiles.
