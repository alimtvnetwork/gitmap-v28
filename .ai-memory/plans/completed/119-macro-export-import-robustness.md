# Plan 119: Macro Multi-Format Export and Safe Import Robustness Engine

> **Origin:** User verification and hardening request: "is it done properly? you don't have proper export, import option for macro and also need to have all export, all import and also export as a single, json, yaml, sqlitedb, and can also be imported safely properly, clear?? understood".
> **Execution Lifecycle:** Parent Task Continuous Loop completed in 12 steps across 2 planning subagents and master orchestration.
> **Status:** Completed & Consolidated.

---

## Executive Summary
This milestone delivers a production-grade, hardened, bidirectional Macro Export and Import subsystem for GitMap. It addresses all edge cases identified during deep-dive audits, ensuring complete parity across CLI commands, robust format handling, strict filesystem security, and atomic persistence.

### Key Capabilities Delivered
1. **Universal Export Commands**:
   - `gitmap macro export all [options]`: Exports all saved macros.
   - `gitmap macro export-all [options]` / `gitmap macro all export [options]`: Dedicated subcommands and modifier syntax.
   - `gitmap macro export single <name> [options]`: Exports a single specified macro.
   - `gitmap macro export-single <name> [options]` / `gitmap macro single export <name> [options]`: Dedicated subcommands and modifier syntax.
   - Root-level aliases in `gitmap/cmd/rootdata.go`: `macro-export-all`, `macro-export-single`.

2. **Universal Import Commands**:
   - `gitmap macro import <file>`: Imports macros with auto format detection.
   - `gitmap macro import all <file>` / `gitmap macro import-all <file>` / `gitmap macro all import <file>`: Explicit batch import.
   - `gitmap macro import single <file> [name]`: Single macro import with target isolation.
   - `gitmap macro import-single <file> [name]`: Dedicated single macro import command.
   - Renaming: `--as <newname>` / `--rename <newname>` allows importing a single macro under a distinct name without collisions.
   - Root-level aliases in `gitmap/cmd/rootdata.go`: `macro-import-all`, `macro-import-single`.

3. **Complete 4-Format Serialization Matrix**:
   - **JSON** (`.json`): Single macro object `{...}` or full catalog array `[...]` with indentation.
   - **YAML** (`.yaml`, `.yml`): Single macro mapping or multi-macro sequence.
   - **SQLite Database** (`.sqlitedb`, `.db`, `.sqlite`, `.sqlite3`): Relational schema with normalized `macros` parent table and `macro_steps` child table with foreign keys, index, and single-connection concurrency ceiling `db.SetMaxOpenConns(1)` with WAL mode.
   - **ZIP Bundle** (`.zip`): Compressed archive bundling individual `.json` macro definitions.

4. **Security & Integrity Hardening**:
   - **Win32 & Traversal Sanitization**: Macro names are checked against directory traversals (`..`, `/`, `\`), ASCII control characters ($< 32$), invalid Win32 characters (`:`, `*`, `?`, `"`, `<`, `>`, `|`, `\0`), reserved device names (`CON`, `PRN`, `AUX`, `NUL`, `COM1-9`, `LPT1-9`), and trailing dots or whitespace.
   - **Step Command Validation**: Replaces rune arithmetic with `strconv.Itoa(i + 1)` ensuring step indexes $\ge 10$ display clean error messages.
   - **Collision Protection**: Skips colliding macros by default; tracks and reports their names in `SkippedNames`. `--force` / `--overwrite` enables replacement.
   - **Missing Target Guard**: Explicit target name imports return a validation error if the requested macro is absent from the file.
   - **Atomic Persistence**: Writes via unique `.tmp` files (`fmt.Sprintf("%s.%d.tmp", target, time.Now().UnixNano())`), flushes via `file.Sync()`, and safely replaces destination to prevent file locking or zero-byte corruptions on Windows.

---

## Consolidated Subtasks

### Subtask 01: Macro Core Serialization & SQLite Export Hardening
- Target files: `gitmap/macro/export.go`, `gitmap/macro/export_sqlite.go`.
- Implemented atomic `.tmp` staging, sync, and safe destination replace in `WriteExportPayload`.
- Cleaned up residual `-wal` and `-shm` database files in `prepareSQLiteExportPath`.
- Added support for `.sqlitedb` in format inference and normalizers.

### Subtask 02: Macro Core Deserialization, Validation & Safety Guards
- Target files: `gitmap/macro/import.go`, `gitmap/macro/import_sqlite.go`, `gitmap/macro/storage.go`.
- Expanded `hasInvalidNameChars` with Win32 character, reserved device name, and trailing punctuation checks.
- Fixed step index reporting in `validateMacroSteps` via `strconv.Itoa(i + 1)`.
- Added `SkippedNames []string` and in-batch duplicate tracking to `ImportResult`.
- Added `RenameAs string`, `IsSingle bool`, and `IsAll bool` to `ImportOptions`.
- Implemented missing target validation in `filterAndValidateImportTargets`.
- Hardened `writeMacroFileAtomic` in `storage.go` with unique temp naming and destination replace.

### Subtask 03: CLI Commands, Flags & Dynamic Positional Reconcilers
- Target files: `gitmap/cmd/macro_export.go`, `gitmap/cmd/macro_import.go`, `gitmap/cmd/macro_cmd.go`, `gitmap/cmd/rootdata.go`.
- Fixed keyword modifier parsing (`all`, `single`) in `macro_export.go` so they do not claim `opts.TargetName`.
- Added support for second positional argument as output file (e.g., `gitmap macro export mymacro backup.json` or `gitmap macro export all backup.sqlitedb`).
- Added heuristic file path disambiguation in `reconcileExportTargets`.
- Added flags `--name`, `--target`, `--macro`, `--as`, `--rename`, `--sqlitedb`, `--zip`, `--json`, `--yaml`.
- Enforced single export validation requiring a macro name.
- Decoupled `macro-export-all`, `macro-export-single`, `macro-import-all`, `macro-import-single` in `rootdata.go`.
- Updated `printMacroUsage` in `macro_cmd.go`.

### Subtask 04: Test Suite Expansion and Verification
- Target files: `gitmap/macro/export_import_test.go`, `gitmap/cmd/macro_export_import_test.go`.
- Added tests for Win32 invalid chars, step index $\ge 10$ formatting, renaming, missing target error, and collision skipping.
- Added CLI tests for positional output file capture, format inference, single flag parsing, and import modifiers.
- Resynced test inventory to 3,532 tests across 123 packages.

### Subtask 05: Documentation, Help Text & Inventory Synchronization
- Target files: `gitmap/helptext/macro.md`, `docs/commands/macro.md`, `docs/commands/automation/readme.md`, `.agents/skills/macro-multi-format-export-import/SKILL.md`.
- Updated all reference docs with subcommand matrix, flags, and executable examples.

---

## Verification & Quality Gates
- `go vet ./...`: 100% clean (zero errors).
- `go build ./...`: 100% clean across all packages.
- Live CLI execution verified:
  - Single JSON/YAML/SQLiteDB exports.
  - All-macro JSON/YAML/SQLiteDB exports.
  - Safe import collision detection and `--dry-run` skip reporting.
  - Single macro import with `--as` renaming.
