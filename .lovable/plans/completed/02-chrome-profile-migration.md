---
slug: "chrome-profile-migration"
status: "pending"
total_tasks: 30
---

# Plan: Chrome Profile Migration (Cross-OS)

## Overview
Enhance gitmap chrome-profile-export and import to support .zip and .sqlite formats for cross-OS migration.

## Tasks
- [ ] 001-task.md: Analyze Chrome SQLite Schema for Export (History, Web Data)
- [ ] 002-task.md: Define Cross-OS File Mapping (Windows vs Linux)
- [ ] 003-task.md: Update constants_chromeprofile.go with Format Flags
- [ ] 004-task.md: Update chromeprofile.go to Parse --format=zip|json|sqlite
- [ ] 005-task.md: Define zipExportRecord struct in chromeprofile_db.go
- [ ] 006-task.md: Add SQLite inclusion flag logic to runChromeProfileExport
- [ ] 007-task.md: Create chromeprofile_sqlite_export.go scaffolding
- [ ] 008-task.md: Implement Chrome SQLite file copy routine (Windows)
- [ ] 009-task.md: Implement Chrome SQLite file copy routine (Linux)
- [ ] 010-task.md: Implement Chrome SQLite file copy routine (macOS)
- [ ] 011-task.md: Filter out encrypted DBs (Login Data, Cookies) from Export
- [ ] 012-task.md: Write SQLite payloads to ZIP archive
- [ ] 013-task.md: Write JSON payload to ZIP archive
- [ ] 014-task.md: Update writeChromeExport to conditionally output sqlite only
- [ ] 015-task.md: Update writeChromeExport to conditionally output zip format
- [ ] 016-task.md: Update defaultChromeExportPath to handle .zip and .sqlite
- [ ] 017-task.md: Update emitChromeSnapshots for multi-format output
- [ ] 018-task.md: Update printChromeArtifacts to handle zip/sqlite files
- [ ] 019-task.md: Add --format flag support to runChromeProfileImport
- [ ] 020-task.md: Create chromeprofile_sqlite_import.go scaffolding
- [ ] 021-task.md: Implement ZIP extraction for profile import
- [ ] 022-task.md: Implement SQLite restoration for profile import
- [ ] 023-task.md: Handle Profile merging when SQLite databases exist
- [ ] 024-task.md: Update applyChromeExport for SQLite blobs
- [ ] 025-task.md: Migrate gitmap.db schema to support Format='zip' and 'sqlite'
- [ ] 026-task.md: Update persistChromeProfile to track new format extensions
- [ ] 027-task.md: Update listChromeProfilesFromDB to summarize zip exports
- [ ] 028-task.md: Add Integration Test: Windows to Linux Zip Export/Import
- [ ] 029-task.md: Add Unit Tests for sqlite filtering and zip bundling
- [ ] 030-task.md: Update README and CLI Help Text for new format flags
