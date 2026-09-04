# 97 — Chrome Profile Management Acceptance Criteria

**Version:** 1.0.0
**Updated:** 2026-09-05
**Status:** Canonical
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Automated Verification Commands

```bash
# 1. Unit testing for Chrome profile registration & reconciliation
go test -v ./cmd/... -run TestChromeProfile

# 2. Zero-nesting architectural validation
python linter-scripts/check-nested-ifs.py

# 3. Error management conformance
python linter-scripts/check-error-management.py

# 4. Specification cross-links validation
python linter-scripts/check-spec-cross-links.py --root spec --repo-root .

# 5. Full CI/CD local runner suite
python 03-ai-scripts/06-cicd-local-runner.py
```

---

## 2. Acceptance Criteria (Given / When / Then)

### AC-01: Complete `info_cache` Attribute Registration
- **GIVEN** a profile snapshot is imported into a new profile directory (e.g. `Profile 5`)
- **WHEN** GitMap registers the profile inside Chrome's `Local State`
- **THEN** the entry in `profile.info_cache[<dir>]` contains all mandatory attributes:
  - `name`, `shortcut_name`, `user_name`
  - `active_time`, `avatar_icon`, `default_avatar_fill_color`, `default_avatar_stroke_color`, `profile_highlight_color`
  - `is_using_default_avatar`, `is_using_default_name`, `is_ephemeral`
- **AND** the values conform to Chromium's profile picker schema.

### AC-02: Profile Directory Display Order
- **GIVEN** a new profile directory (e.g. `Profile 5`) is registered in `Local State`
- **WHEN** `Local State` is written to disk
- **THEN** `profile.profiles_order` contains the new profile directory key
- **AND** no duplicate profile entries exist in `profile.profiles_order`.

### AC-03: Concurrency Process Guard & Warning
- **GIVEN** Google Chrome is actively running (`chrome.exe` processes present)
- **WHEN** a profile import or registration is executed
- **THEN** GitMap detects active Chrome processes via `isChromeRunning()`
- **AND** GitMap prints an actionable warning notifying the user that Chrome must be restarted or reconciled to prevent in-memory cache overwrites.

### AC-04: Preferences Sanitization & Name Alignment
- **GIVEN** a profile directory is created or updated
- **WHEN** the profile's internal `Preferences` file is written
- **THEN** `profile.name` matches the assigned profile display name
- **AND** `profile.using_default_name` is `false`
- **AND** foreign account authentication tokens and custodian identifiers are sanitized.

### AC-05: Orphan Profile Reconciliation
- **GIVEN** a profile folder exists on disk (e.g. `Profile 5`) but is missing from `Local State`
- **WHEN** `gitmap chrome profile reconcile` is executed
- **THEN** GitMap scans all profile directories under `User Data`
- **AND** automatically restores missing entries into `profile.info_cache` and `profile.profiles_order`
- **AND** reports the reconciled profiles clearly in terminal output.
