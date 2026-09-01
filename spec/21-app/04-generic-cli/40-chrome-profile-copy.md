# 40 — Chrome Profile Copy / Export / Import

Status: **DRAFT** — spec only, not yet implemented.

## Goal

Let a user clone a Chrome browser profile (extensions, bookmarks,
settings, flags, preferences) into a new local profile, and round-trip
the same payload through CSV/JSON for backup, diff, and replay across
machines. All operations are tracked in the gitmap SQLite DB so they
can be listed, audited, and re-applied.

## Task Breakdown

1. **CLI surface (new top-level commands)**
   - `gitmap chrome-profile-copy <from> <to>` (alias `cpc`)
     - `<from>`: existing Chrome profile name (e.g. `Default`, `Profile 1`).
     - `<to>`: destination profile name. Created if missing; merged if present.
     - Created profiles are **offline** (no Google sign-in). User can sign in later from chrome://settings manually.
   - `gitmap chrome-profile-export <profile> [file]` (alias `cpe`)
     - Default file: `chrome-profile-<name>-<UTC>.json`.
     - `--csv` flag writes flat CSV instead of JSON.
   - `gitmap chrome-profile-import <file> [profile]` (alias `cpi`)
     - If `[profile]` omitted, restores into the source profile name embedded in the file.
     - `--confirm` required when target exists.
   - `gitmap chrome-profile-list` (alias `cpl`)
     - Lists installed Chrome profiles + last export timestamp from DB.

2. **What gets copied / serialized**
   - Bookmarks (`Bookmarks` JSON file).
   - Extensions (folder copy + `Preferences > extensions.settings` slice).
   - Preferences (`Preferences`, `Secure Preferences` — sanitized; OS-bound encrypted blobs flagged, not copied).
   - Flags (`Local State > browser.enabled_labs_experiments`).
   - Site settings, search engines, startup pages, autofill profiles (non-secret fields only).
   - **Excluded**: cookies, login data, history, cache, GPU cache, sync tokens — these are device/account bound.

3. **Storage layout**
   - Source root (Windows): `%LOCALAPPDATA%\Google\Chrome\User Data\<profile>`.
   - macOS: `~/Library/Application Support/Google/Chrome/<profile>`.
   - Linux: `~/.config/google-chrome/<profile>`.
   - Chrome MUST be fully closed before copy/import (detect running process; refuse with actionable hint).

4. **DB schema (new tables, migration 015)**
   ```sql
   CREATE TABLE ChromeProfile (
     ChromeProfileId INTEGER PRIMARY KEY AUTOINCREMENT,
     Name            TEXT NOT NULL UNIQUE,
     SourcePath      TEXT NOT NULL DEFAULT '',
     IsOffline       INTEGER NOT NULL DEFAULT 1,
     CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
     UpdatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
   );
   CREATE TABLE ChromeProfileExport (
     ChromeProfileExportId INTEGER PRIMARY KEY AUTOINCREMENT,
     ChromeProfileId       INTEGER NOT NULL REFERENCES ChromeProfile(ChromeProfileId) ON DELETE CASCADE,
     Format                TEXT NOT NULL,           -- 'json' | 'csv'
     FilePath              TEXT NOT NULL,
     ByteSize              INTEGER NOT NULL DEFAULT 0,
     Sha256                TEXT NOT NULL DEFAULT '',
     CreatedAt             TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
   );
   GRANT statements N/A — local SQLite.
   ```
   - Every `cpc` / `cpe` / `cpi` invocation upserts `ChromeProfile` and inserts a row in `ChromeProfileExport` (FilePath empty for `cpc`).

5. **CSV format**
   Flat key/value with category column:
   ```
   category,key,value
   bookmark,/Bookmarks Bar/GitHub,https://github.com
   extension,nmmhkkegccagdldgiimedpiccmgmieda,enabled
   flag,enable-quic,true
   preference,profile.default_content_setting_values.notifications,2
   ```

6. **JSON format**
   ```json
   {
     "schemaVersion": 1,
     "exportedAt": "2026-06-19T12:00:00Z",
     "sourceProfile": "Default",
     "bookmarks": [...],
     "extensions": [{"id":"...","name":"...","version":"..."}],
     "flags": ["enable-quic"],
     "preferences": { ... },
     "searchEngines": [...]
   }
   ```

7. **UI (src/pages)**
   - New page `ChromeProfileSpec.tsx` describing commands + sample invocations.
   - Route added to App router.
   - Docs page index entry.

8. **README + help text**
   - Add a "Chrome profile" section to `readme.txt` with the four commands.
   - Add `Help*` constants in `gitmap/constants/constants_chromeprofile.go`.
   - Register handlers in `dispatchTooling` (or a new `dispatchBrowser`).
   - Mark const block with `// gitmap:cmd top-level` for completion generator.

9. **Tests**
   - Unit: serializer round-trip (JSON → struct → JSON byte-identical).
   - Integration: tmp-dir fake Chrome profile → export → import into second tmp dir → assert file tree equality (modulo timestamps).
   - DB: `ChromeProfile` upsert idempotency.

10. **Version bump**
    - Minor bump on landing (`6.29.x` → `6.30.0`).
    - CHANGELOG entry under `v6.30.0` listing all four commands.

## Out of scope (v1)

- Sync with Google account / cloud restore.
- Cookie / saved-password migration (security boundary).
- Edge / Brave / Firefox profiles (separate spec when requested).
- GUI dialogs — CLI-only for v1.

## Open questions

- **Q1**: For `cpc <from> <to>` when `<to>` exists, do we **merge** (default) or require `--overwrite`? Spec assumes merge.
- **Q2**: Encrypted Chrome blobs (`Login Data`, OS-DPAPI-bound entries) — flag-and-skip vs. hard-fail? Spec assumes flag-and-skip with stderr warning.
- **Q3**: Should `cpi` auto-launch Chrome after import? Spec assumes no; print hint instead.
