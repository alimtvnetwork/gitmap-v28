---
name: macro-multi-format-export-import
description: Autonomously export and import GitMap macros across JSON, YAML, and SQLite database formats with polymorphic parsing, schema validation, collision protection, and atomic persistence.
---

# Macro Multi-Format Export and Safe Import Skill

Use this skill when implementing, auditing, refactoring, or verifying macro export and import operations across GitMap.

## Capabilities & Architecture
1. **Multi-Format Serialization**:
   - **JSON**: Single macro object or all-macro catalog array with indentation.
   - **YAML**: Single macro or multi-macro definitions formatted via `gopkg.in/yaml.v3`.
   - **SQLite Database**: Portable `.db` / `.sqlite` file with normalized `macros` parent table and `macro_steps` child table via `modernc.org/sqlite`.
   - **ZIP Archives**: Compressed archive containing individual `.json` definitions.

2. **Polymorphic Dual-Shape Import**:
   - Transparently parses both single objects (`{"name": "...", "steps": [...]}`) and collections (`[{"name": "..."}, ...]`).
   - Automatically detects format via file extensions (`.json`, `.yaml`, `.yml`, `.db`, `.sqlite`, `.sqlite3`, `.zip`).
   - Supports explicit format flags (`--json`, `--yaml`, `--sqlite`, `--sqlitedb`, `--db`, `--zip`, `--format=<val>`).

3. **Safety & Integrity Contract**:
   - **Path Traversal Guards**: Rejects any macro names containing `..`, `/`, or `\`.
   - **Step Validation**: Ensures non-empty command lines and normalizes step numbering.
   - **Collision Guards**: Rejects overwriting existing macros by default; requires `--overwrite` or `--force`.
   - **Dry Run Simulation**: `--dry-run` performs full schema validation and conflict reporting without modifying disk.
   - **Atomic Writes**: Writes to `.tmp` file and atomically renames to prevent partial write corruption.
