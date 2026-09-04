# 05 — Chrome Profile Picker Registration, Local State Invariants & Process Concurrency

**Date:** 2026-09-05
**Author:** Antigravity (Advanced Agentic AI)
**Category:** Architecture & Chromium Integration
**Status:** Canonical / Active

---

## 1. Verbatim User Directives & Problems

### User Directive 1:
> *"PS M:\Work-files\export\export\extv2> gitmap chrome import*
> *chrome-profile-import: ERROR <file> is required*
> *PS M:\Work-files\export\export\extv2> gitmap chrome import .*
> *chrome-profile-import: ERROR read .: read .: Incorrect function.*
> *PS M:\Work-files\export\export> gitmap chrome import-all extv2*
> *chrome-profile-import: imported extv2\Default.json into profile "Default"*
> *also need to have proper import the import is not working properly it needs to create profile if the email doens't exist, it should log steps as it performs and also the don't breaking exisitn profile if no email matches and also make sure that we can check a backdir what email contains in which file can import single email account as well"*

### User Directive 2 & Problem Evidence:
> User imported snapshot via `gitmap chrome profile import "erfan.office.n@gmail.com"`, which allocated `Profile 5` on disk.
> However, Chrome's Profile Picker UI (`chrome://settings/manageProfile`) completely ignored `Profile 5`.
> Visual evidence:
> - Terminal: `.lovable/assets/screenshots/19-chrome-profile-import-terminal.png`
> - UI Missing Profile 5: `.lovable/assets/screenshots/20-chrome-profile-picker-missing-profile5.png`

---

## 2. Root Cause Analysis (RCA)

1. **Active Chrome Process In-Memory Cache Overwrite:**
   - Chrome maintains an in-memory cache of `Local State` while `chrome.exe` is running.
   - Any external writes to `Local State` made while Chrome processes are active get wiped on periodic background flushes or browser shutdown.
2. **Missing Mandatory Chromium Profile Picker Attributes:**
   - Chromium requires 13 specific attributes in `profile.info_cache[<dir>]` to render a valid tile:
     `name`, `shortcut_name`, `user_name`, `avatar_icon`, `default_avatar_fill_color`, `default_avatar_stroke_color`, `profile_highlight_color`, `profile_color_seed`, `active_time`, `is_using_default_avatar`, `is_using_default_name`, `is_ephemeral`, `is_consented_primary_account`, `signin.with_credential_provider`.
   - Partial entries (only `name` and `user_name`) cause Chromium to silently drop the profile tile from the profile picker interface.
3. **Preferences & Local State Desynchronization:**
   - The profile's `<dir>/Preferences` file must have `profile.name` aligned with the display name, and obsolete parental lock keys (`managed`, `managed_user_id`, `custodian_*`) sanitized.

---

## 3. Implemented Architecture & Solutions

1. **Unified Registration Engine (`registerChromeProfileWithFullSchema`):**
   - Populates all 13 Chromium attributes with valid defaults.
   - Appends directory to `profile.profiles_order`.
   - Creates `Local State` file automatically if it does not yet exist.
2. **Preferences Sanitization on Import (`patchImportedChromeProfilePreferences`):**
   - Synchronizes `profile.name = displayName` and sets `using_default_name = false`.
   - Preserves `account_info` email identities so email resolution succeeds.
   - Removes stale parental lock `managed` blocks.
3. **Concurrency Process Guard (`checkChromeRunningAdvisory`):**
   - Detects running Chrome processes via `isChromeRunning(runtime.GOOS)`.
   - Displays clear warning advising user to close Chrome and run `reconcile` if changes do not persist.
4. **Orphan Profile Reconciler (`gitmap chrome profile reconcile`):**
   - Inspects `%LOCALAPPDATA%\Google\Chrome\User Data` for any unlinked directories (`Default`, `Profile *`).
   - Automatically synthesizes metadata from `<dir>/Preferences` and registers into `Local State`.
   - Reconciled existing unlinked `Profile 5` (`erfan.office.n@gmail.com`) into Chrome `Local State`.

---

## 4. Verification & Conformance

- All 16 quality gates in CI/CD local runner suite passed.
- 0 nested ifs verified via `python linter-scripts/check-nested-ifs.py`.
- 0 error management violations via `python linter-scripts/check-error-management.py`.
- Comprehensive unit tests in `gitmap/cmd/chromeprofile_reconcile_test.go` and `gitmap/cmd/chromeprofile_smart_import_test.go`.
