# Seedable Config Architecture: Acceptance Criteria

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**Format:** GIVEN/WHEN/THEN (E2E-test-ready)

---

## Core Seeding Flow (00-overview.md)

### SC-01: First-Run Seeding

**GIVEN** the CLI starts for the first time with no existing configuration in the database  
**WHEN** the seed initialization runs  
**THEN** all values from `config.seed.json` are inserted into the settings table  
**AND** each setting has `IsUserModified: false`  
**AND** the config version is set to "1.0.0" in the database

**Edge Cases:**
- **GIVEN** `config.seed.json` is missing **WHEN** first-run seeding is attempted **THEN** the tool exits with a clear error referencing the expected file path
- **GIVEN** `config.seed.json` fails JSON schema validation **WHEN** seeding is attempted **THEN** all validation errors are listed and seeding is aborted

### SC-02: Version-Gated Seeding (Golden Rule)

**GIVEN** the database already has configuration with version "1.2.0"  
**WHEN** the CLI starts with a `config.seed.json` at version "1.3.0"  
**THEN** only NEW keys (not present in DB) are seeded from the new version  
**AND** existing keys with `IsUserModified: true` are NOT overwritten  
**AND** the version is updated to "1.3.0"

**Edge Cases:**
- **GIVEN** the seed version is LOWER than the DB version **WHEN** the CLI starts **THEN** seeding is skipped entirely and a debug log is emitted
- **GIVEN** the seed version EQUALS the DB version **WHEN** the CLI starts **THEN** seeding is skipped (no-op)
- **GIVEN** two CLI instances start simultaneously with the same seed file **WHEN** both attempt seeding **THEN** the database transaction ensures only one succeeds and the other detects the version is already current

### SC-03: User Modification Tracking

**GIVEN** a setting exists with `IsUserModified: false`  
**WHEN** the user updates the setting via API or CLI  
**THEN** the value is updated and `IsUserModified` is set to `true`  
**AND** subsequent seeding at a higher version will NOT overwrite this setting

### SC-04: Changelog Versioning

**GIVEN** a configuration change is persisted (user or seed)  
**WHEN** the version is incremented  
**THEN** a new entry is appended to `changelog.md` with version, date, and change description

**Edge Cases:**
- **GIVEN** `changelog.md` does not exist **WHEN** a version change occurs **THEN** the file is created with a header and the first entry
- **GIVEN** `changelog.md` is read-only **WHEN** an append is attempted **THEN** the change proceeds but a warning is logged about the changelog failure

---

## Settings Cache (sync.Map)

### SC-05: In-Memory Cache

**GIVEN** settings are loaded from the database  
**WHEN** a setting is requested via `GetSetting(key)`  
**THEN** the value is returned from the `sync.Map` cache without a database query

**Edge Cases:**
- **GIVEN** the cache is empty (cold start) **WHEN** `GetSetting(key)` is called **THEN** all settings are loaded from DB into cache first, then the value is returned
- **GIVEN** a setting is updated via API **WHEN** the update is confirmed **THEN** the cache is invalidated for that key and refreshed on next read
- **GIVEN** two goroutines concurrently write different values for the same key via `sync.Map.Store()` **WHEN** both complete **THEN** the last write wins and subsequent reads return a consistent value (no panic or corruption)

---

## RAG Chunk Settings (02-rag-chunk-settings.md)

### RC-01: Chunk Size Validation

**GIVEN** a RAG chunk size setting is being configured  
**WHEN** the value is set to 0 or negative  
**THEN** validation rejects the value with a specific error message listing the valid range (100–10000)

### RC-02: Chunk Overlap Validation

**GIVEN** chunk size is set to 500  
**WHEN** chunk overlap is set to 500 or higher  
**THEN** validation rejects with "Overlap must be less than ChunkSize"

---

## Validation Data Seeding (06-validation-data-seeding.md)

### VD-01: Validation Arrays Loaded

**GIVEN** the CLI starts with seedable config containing validation arrays  
**WHEN** the seeding process runs  
**THEN** validation arrays (e.g., allowed file extensions, blocked keywords) are loaded from CW Config into the root DB  
**AND** they are accessible via the settings service for runtime validation

---

*Wave 6 — Batch 1 (Patched): Seedable Config acceptance criteria with added edge cases for concurrent seeding race conditions and sync.Map concurrent write safety.*
