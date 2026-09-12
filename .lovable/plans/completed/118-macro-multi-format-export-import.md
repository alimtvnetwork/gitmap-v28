# Plan 118: Macro Multi-Format Export and Safe Import Architecture

> **Execution Milestone**: Completed in 1 continuous loop session across 5 discrete subtasks (N=350 budget).
> **Originating Request**: "you don't have proper export, import option for macro and also need to have all export, all import and also export as a single, json, yaml, sqlitedb, and can also be imported safely propery, clear?? understood"
> **Completed At**: 2026-09-12T15:30:00Z
> **Author**: Antigravity Autonomous Continuous Loop Manager

---

## 1. Executive Summary & Architectural Overview
GitMap has been enhanced with a robust, production-grade Macro Export and Import subsystem allowing developers to:
- Export a single macro or all macros (`--all`).
- Multi-format serialization: **JSON** (indented), **YAML** (`gopkg.in/yaml.v3`), portable **SQLite Database** (`.db` / `.sqlite` via `modernc.org/sqlite`), and **ZIP Archives** (`.zip`).
- Safe polymorphic import with auto-inference, schema validation, path traversal guards (`..`, `/`, `\`), dry-run preview (`--dry-run`), collision management (`--force` / `--overwrite`), and exception exclusion filtering (`-except`).
- Direct CLI command routing (`gitmap macro export`, `gitmap macro import`, `gitmap macro-export`, `gitmap macro-import`).
- Full documentation parity across terminal help (`helptext/macro.md`), web reference (`docs/commands/macro.md`, `docs/commands/automation/readme.md`), root `readme.md`, and frontend commands catalog (`src/data/commands.ts`).

---

## 2. Implemented Subtasks (Consolidated)

### Subtask 01: Macro Core Multi-Format Export Engine
- **Target Files**: `gitmap/macro/export.go`, `gitmap/macro/export_sqlite.go`
- **Accomplishments**:
  - Implemented `ExportToJSON(data any) ([]byte, error)` with `constants.JSONIndent`.
  - Implemented `ExportToYAML(data any) ([]byte, error)` using `gopkg.in/yaml.v3`.
  - Implemented `ExportToZIP(macros []Macro, filePath string) error` packaging individual `.json` definitions into a zip archive.
  - Implemented `ExportMacrosToSQLite(macros []Macro, dbPath string) error` creating normalized `macros` and `macro_steps` relational tables with `modernc.org/sqlite`, single connection pooling (`SetMaxOpenConns(1)`), PRAGMAs (`WAL`, `busy_timeout=5000`), RFC3339 timestamps, and transactional commits.
  - Added `FilterMacrosForExport` and `isExcludedMacro` for exclusion filtering.

### Subtask 02: Macro Core Multi-Format Safe Import Engine
- **Target Files**: `gitmap/macro/import.go`, `gitmap/macro/import_sqlite.go`, `gitmap/macro/storage.go`
- **Accomplishments**:
  - Implemented `ValidateMacro(m *Macro)` with traversal rejection (`..`, `/`, `\`), empty name check, and zero-step check.
  - Implemented polymorphic dual-shape parsers for JSON (`ParseImportJSON`) and YAML (`ParseImportYAML`) transparently accepting both arrays (`[]Macro`) and single objects (`Macro`).
  - Implemented `ParseImportSQLite(dbPath string)` querying and assembling macros with steps.
  - Implemented `ParseImportZIP(filePath string)` extracting `.json` definitions.
  - Implemented `ImportMacros(macros []Macro, opts ImportOptions)` supporting collision detection, `--dry-run` simulation, `--force`/`--overwrite`, metrics collection (`ImportResult`), and atomic storage.
  - Added `MacroExists(name string) bool` to `storage.go`.

### Subtask 03: Macro CLI Handlers and Dispatch
- **Target Files**: `gitmap/cmd/macro_export.go`, `gitmap/cmd/macro_import.go`, `gitmap/cmd/macro_cmd.go`, `gitmap/cmd/rootdata.go`
- **Accomplishments**:
  - Implemented `runMacroExport` parsing `-f`, `--file`, `-o`, `--out`, `--format`, `--all`, `--json`, `--yaml`, `--sqlite`, `--zip`, `-except`.
  - Implemented `runMacroImport` parsing `-f`, `--file`, `--format`, `--force`/`--overwrite`, `--dry-run`, `-except`.
  - Extended `routeMacroSubcommand` in `macro_cmd.go` to handle `export`, `exp`, `dump`, `export-all`, `import`, `imp`, `load`, `restore`, `import-all`.
  - Registered root commands `macro-export` and `macro-import` in `rootdata.go`.
  - Updated `printMacroUsage()` in `macro_cmd.go`.

### Subtask 04: Macro E2E Unit Tests and Quality Gates
- **Target Files**: `gitmap/macro/export_import_test.go`, `gitmap/cmd/macro_export_import_test.go`
- **Accomplishments**:
  - Added 6 semantic behavior tests in `macro`: `TestExportMacros_SingleAndAllJSON`, `TestExportMacros_SingleAndAllYAML`, `TestExportMacros_SQLiteDatabaseRoundtrip`, `TestExportMacros_ZIPArchiveRoundtrip`, `TestValidateMacro_RejectsInvalidNameAndEmptySteps`, `TestImportMacros_HandlesDryRunAndExceptFilter`.
  - Added 3 semantic behavior tests in `cmd`: `TestParseMacroExportOpts_ParsesFlags`, `TestParseMacroImportOpts_ParsesFlags`, `TestInferMacroExportFormat_InfersFromExtension`.
  - Verified 100% green compilation via `go vet ./...` and `go build ./...`.

### Subtask 05: Documentation, Help Text, and Task Consolidation
- **Target Files**: `gitmap/helptext/macro.md`, `docs/commands/macro.md`, `docs/commands/automation/readme.md`, `readme.md`, `src/data/commands.ts`
- **Accomplishments**:
  - Updated terminal help with export/import usage, tables, flags, and copy-paste examples.
  - Updated web reference documentation and root `readme.md` under `### Interactive Macros & Automation`.
  - Synchronized `.lovable/test-inventory.json` indexing 3,522 tests across 123 packages.
  - Atomic tracking in `.lovable/temp/recent-file-changes.json`.
