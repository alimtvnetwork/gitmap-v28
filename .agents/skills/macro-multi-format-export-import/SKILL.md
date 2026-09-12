---
name: macro-multi-format-export-import
description: Autonomously export and import GitMap macros across JSON, YAML, and SQLite database formats with polymorphic parsing, schema validation, collision protection, and atomic persistence.
---

# Macro Multi-Format Export and Safe Import Skill

Use this skill when implementing, auditing, refactoring, or verifying macro export and import operations across GitMap.

## Capabilities & Architecture
1. **Multi-Format Serialization**:
   - **JSON**: Single macro object `{...}` or all-macro catalog array `[...]` with indentation.
   - **YAML**: Single macro or multi-macro definitions formatted via `gopkg.in/yaml.v3`.
   - **SQLite Database**: Portable `.db`, `.sqlite`, `.sqlite3`, `.sqlitedb` file with normalized `macros` parent table and `macro_steps` child table via `modernc.org/sqlite`. Concurrency ceiling `db.SetMaxOpenConns(1)` with `WAL` journal mode and busy timeout.
   - **ZIP Archives**: Compressed archive containing individual `.json` definitions.

2. **Polymorphic Dual-Shape Import**:
   - Transparently parses both single objects (`{"name": "...", "steps": [...]}`) and collections (`[{"name": "..."}, ...]`).
   - Automatically detects format via file extensions (`.json`, `.yaml`, `.yml`, `.db`, `.sqlite`, `.sqlite3`, `.sqlitedb`, `.zip`).
   - Supports explicit format flags (`--json`, `--yaml`, `--sqlite`, `--sqlitedb`, `--db`, `--zip`, `--format=<val>`).

3. **Subcommand & Alias Matrix**:
   - `gitmap macro export [all|single] [name] [-o <file>]`
   - `gitmap macro export-all` / `gitmap macro all export`
   - `gitmap macro export-single <name>` / `gitmap macro single export <name>`
   - `gitmap macro import [all|single] <file> [name] [--as <name>] [--force] [--dry-run]`
   - `gitmap macro import-all <file>` / `gitmap macro all import <file>`
   - `gitmap macro import-single <file> [name] [--as <name>]` / `gitmap macro single import <file> [name]`

4. **Safety & Integrity Contract**:
   - **Path Traversal & Win32 Guards**: Rejects macro names containing `..`, `/`, `\`, control chars $< 32$, invalid Win32 characters (`:`, `*`, `?`, `"`, `<`, `>`, `|`, `\0`), reserved device names (`CON`, `PRN`, `AUX`, `NUL`, `COM1-9`, `LPT1-9`), and trailing dots or spaces.
   - **Step Validation**: Ensures non-empty command lines, normalizes step numbering, and formats step index errors cleanly using `strconv.Itoa(i + 1)`.
   - **Collision Guards**: Skips existing macros by default and records their names in `SkippedNames`; requires `--force` or `--overwrite` to replace.
   - **Renaming Capability**: `--as <name>` / `--rename <name>` allows renaming a single imported macro to prevent collisions.
   - **Dry Run Simulation**: `--dry-run` performs full schema validation and conflict reporting without modifying disk.
   - **Atomic Writes**: Writes to `.tmp` file, flushes via `file.Sync()`, and safely replaces destination to prevent partial write corruption or Windows file locks.
