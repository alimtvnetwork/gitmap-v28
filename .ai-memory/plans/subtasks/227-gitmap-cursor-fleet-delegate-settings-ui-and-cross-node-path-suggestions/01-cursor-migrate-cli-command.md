# Subtask 227.1: Native GitMap Cursor Delegation CLI Command & Script Flags

- **Parent Spec:** [01-architecture-spec.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/01-architecture-spec.md)
- **Master Ledger:** [00-master-audit-ledger.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/00-master-audit-ledger.md)
- **Status:** Pending
- **Target Subsystems:** `cli/cmdcursor`, `repo-secrets/05-scripts`, `cli/store`

---

## 1. Objective

Implement the native first-class CLI command `gitmap cursor migrate` (with aliases `m`, `delegate`, `del`) in `cli/cmdcursor/`, wire it into `cursor_cmd.go`, update Python migration engine flags in `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py`, record dual Split-DB telemetry (`commands.db` and `ai-instruction/sql.db`), and provide comprehensive unit tests in `cursor_migrate_test.go`.

---

## 2. File Modification & Creation Inventory

| File | Action | Responsibilities | Max Lines |
|:---|:---|:---|:---|
| `cli/cmdcursor/cursor_migrate.go` | **Create** | `CursorMigrateOptions` struct, flag parsing, mutual exclusivity validation, help text rendering, `RunCursorMigrate` entry point | < 95 |
| `cli/cmdcursor/cursor_migrate_exec.go` | **Create** | Script location resolution, Python CLI argument assembly, streaming execution via `cmdpy`, and dual Split-DB telemetry dispatch | < 95 |
| `cli/cmdcursor/cursor_cmd.go` | **Modify** | Add `migrate`, `m`, `delegate`, `del` cases to `routeCursorSubcommand`; update `RenderCursorHelp()` | < 90 |
| `cli/cmdcursor/cursor_migrate_test.go` | **Create** | Unit tests for default options, custom flags, mutual exclusivity errors, argument builder, and help detection | < 120 |
| `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py` | **Modify** | Add `--exclude`, `--skip-settings`, `--settings-only` to `parse_args()`; enforce staging filters | — |

---

## 3. Step-by-Step Implementation Plan

### Step 3.1: Define Data Models & Option Parsing (`cli/cmdcursor/cursor_migrate.go`)
1. Create `CursorMigrateOptions` struct:
   - `Node string` (default: `"u1"`)
   - `Exclude string` (comma-separated exclusion patterns)
   - `SkipSettings bool`
   - `SettingsOnly bool`
   - `DryRun bool`
   - `NoBackup bool`
   - `Force bool`
   - `Output string`
2. Implement `parseCursorMigrateOptions(args []string) (CursorMigrateOptions, error)`:
   - Handle flags: `--node`, `-n`, `--exclude`, `--skip-settings`, `--settings-only`, `--dry-run`, `-d`, `--no-backup`, `--force`, `-f`, `--output`, `-o`.
   - Validate mutual exclusivity: if both `SkipSettings` and `SettingsOnly` are true, return `apperror.NewWithDetails("cmd.cursor.migrate", "E1035", "cannot specify both --skip-settings and --settings-only", "cmdcursor", apperror.ErrorTypeValidation, apperror.SeverityError, nil)`.
   - Handle help flags (`help`, `--help`, `-h`) by calling `renderCursorMigrateHelp()`.
3. Implement `renderCursorMigrateHelp()`:
   - Render clean ANSI reference documenting usage, aliases, flags, and examples.
4. Implement `RunCursorMigrate(args []string) error`:
   - Parse options, check for help requests, and invoke execution delegate `executeCursorMigration(opts)`.

### Step 3.2: Implement Execution & Split-DB Telemetry (`cli/cmdcursor/cursor_migrate_exec.go`)
1. Implement `resolveMigrationScriptPath() (string, error)`:
   - Check relative repo path `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py`.
   - Check working directory ancestors if run from subdirectories.
2. Implement `buildMigrationArgs(scriptPath string, opts CursorMigrateOptions) []string`:
   - Always prepend `scriptPath` and `--sync`.
   - Forward `--node <Node>`, `--output <Output>`.
   - Forward `--exclude <Exclude>` if non-empty.
   - Forward boolean flags: `--skip-settings`, `--settings-only`, `--dry-run`, `--no-backup`, `--force`.
3. Implement `executeCursorMigration(opts CursorMigrateOptions) error`:
   - Resolve script path and assemble arguments.
   - Print banner: `fmt.Printf("⚡ Initiating Cursor fleet migration to node '%s'...\n", opts.Node)`.
   - Invoke `cmdpy.ExecutePythonStreaming(binary, args, cwd)`.
   - Record dual telemetry:
     - `commands.db`: `store.OpenCommandHistorySplitDB("").InsertCommandRecord(...)`.
     - `ai-instruction/sql.db`: `store.RecordAiExecution("cursor_migrate", cmdLine, argsJson, cwd, "127.0.0.1", durationMs, exitCode, errMsg, isSuccess)`.
   - Return clean wrapped error if exit code != 0.

### Step 3.3: Wire Subcommands into `cli/cmdcursor/cursor_cmd.go`
1. In `routeCursorSubcommand(sub string, args []string) error`:
   - Add route cases:
     ```go
     case "migrate", "m", "delegate", "del":
         return RunCursorMigrate(args[1:])
     ```
2. In `RenderCursorHelp()`:
   - Add help entries:
     ```
     migrate, delegate, del        Migrate or delegate Cursor memories, chats & projects to fleet node
     ```

### Step 3.4: Upgrade Migration Script Flags (`repo-secrets/05-scripts/migrate-cursor-memories-conversations.py`)
1. In `parse_args()`:
   - Add `--exclude` argument: `nargs="?"`, comma-separated or multi-value patterns.
   - Add `--skip-settings`: `action="store_true"`, help text.
   - Add `--settings-only`: `action="store_true"`, help text.
   - Validate mutual exclusivity of `--skip-settings` and `--settings-only`.
2. In staging routines (`stage_user_config` and `stage_cursor_home`):
   - When `skip_settings` is enabled: skip copying `settings.json`, keybindings, snippet configs.
   - When `settings_only` is enabled: skip copying `state.vscdb`, `conversation-search.db`, and `workspaceStorage/`.
   - When `exclude` patterns are passed: skip any file or directory path matching the exclusion substring.

### Step 3.5: Author Unit Tests (`cli/cmdcursor/cursor_migrate_test.go`)
1. `TestParseCursorMigrateOptions_Defaults`: Validate default node is `"u1"` and booleans default to `false`.
2. `TestParseCursorMigrateOptions_CustomFlags`: Validate full flag parsing (`--node u2 --exclude "tests,tmp" --dry-run --force --no-backup`).
3. `TestParseCursorMigrateOptions_MutuallyExclusive`: Assert error when both `--skip-settings` and `--settings-only` are passed.
4. `TestBuildMigrationArgs`: Assert exact argument sequence passed to Python engine.
5. `TestCursorMigrateHelp`: Verify help flag returns nil without error.

---

## 4. Verification & Acceptance Criteria

1. **Test Suite:** Run `go test -v ./cli/cmdcursor/...` ensuring 100% pass rate.
2. **Flag Verification:** Test `gitmap cursor migrate --dry-run` and verify dry-run output is generated without altering disk.
3. **Telemetry Check:** Query SQLite Split-DB to confirm telemetry row insertion in `commands.db` and `ai-instruction/sql.db`.
4. **Linter Compliance:** Zero violations in `03-ai-scripts/09-check-nested-ifs.py` and `03-ai-scripts/10-check-enum-and-boolean.py`.
