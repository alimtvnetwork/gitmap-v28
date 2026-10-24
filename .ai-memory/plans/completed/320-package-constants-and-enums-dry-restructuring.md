# Plan 320: Package Constants & Enums DRY Restructuring (Phase 2: 5 Core Packages)

## Executive Summary
This milestone executes the strict `01-prompts/15-cg-execute/06-constants-and-enums.md` coding guideline protocol across 5 central CLI packages in `cli/`:
1. `cli/cmdsync` (Repository Baseline Sync)
2. `cli/cmdplans` (Plan Inventory & Navigation)
3. `cli/cmdstatus` (Workspace Multi-Repository Status)
4. `cli/cmdlogin` (Authentication & GitHub Credential Management)
5. `cli/cmdcargo` (Cargo Runner & Rust Toolchain Automation)

Each package received:
- Dedicated leaf package `enums/` with strongly-typed enums ending in `Type`, deterministic methods (`IsValid()`, `String()`, `Parse*()`, `All*()`), positive booleans (`is*`, `has*`), and comprehensive unit tests (`enums_test.go`).
- Dedicated leaf package `constants/` with zero magic strings, centralized format tokens, default values, command identifiers, and unit tests (`constants_test.go`).
- Refactored caller files in each parent package eliminating magic string comparisons, raw literals, and repetitive token definitions.
- FastGate and `check-enum-guidelines.py` pre-commit validation passing with 100% green status.

---

## Technical Outcomes by Package

### 1. `cli/cmdsync` (Subtask 320-01)
- **Leaf Enums (`cli/cmdsync/enums/`):**
  - `SyncTargetType` (`all`, `dirty`, `clean`, `stashed`, `missing`)
  - `SyncRepoStatusType` (`synced`, `modified`, `conflict`, `skipped`, `error`)
  - `SyncActionType` (`fetch`, `pull`, `merge`, `rebase`, `checkout`)
  - `SyncModeType` (`fast-forward`, `rebase`, `dry-run`)
  - Unit tests in `enums_test.go` with 100% statement coverage.
- **Leaf Constants (`cli/cmdsync/constants/`):**
  - Centralized subcommands, CLI flags, target tokens, file/directory names, statuses, and error codes.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring & Modularization:**
  - Refactored callers across `sync.go`, `worker.go`, `sync_dispatch.go`, `mirror.go`, `boundary.go`, `defaults.go`, and `helpers_shim.go`.
  - Extracted `cli/cmdsync/sync_prettier.go` (140 lines) to keep `sync_dispatch.go` below the 500-line limit and flatten nested conditionals.
  - Enforced error return sovereignty; replaced `cliexit.HandleError` calls with direct error returns.

### 2. `cli/cmdplans` (Subtask 320-02)
- **Leaf Enums (`cli/cmdplans/enums/`):**
  - `PlanStatusType` (`pending`, `active`, `in_progress`, `completed`, `archived`, `failed`)
  - `PlansFilterModeType` (`all`, `active`, `pending`, `completed`)
  - `PlansOutputFormatType` (`table`, `json`, `markdown`, `csv`, `yaml`)
  - `PlanMetaFileType` (`readme.md`, `summary.md`, `checklist.md`, etc.)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmdplans/constants/`):**
  - CLI command identifiers, flag strings, directory path tokens, UI colors, and usage text.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring:**
  - Refactored `plans.go`, `plans_types.go`, `plans_parser.go`, `plans_parser_meta.go`, `plans_format.go`, and `plans_render.go`.
  - Decomposed functions exceeding 15 lines; flattened nested conditionals with early returns.

### 3. `cli/cmdstatus` (Subtask 320-03)
- **Leaf Enums (`cli/cmdstatus/enums/`):**
  - `StatusFilterModeType` (`all`, `dirty`, `ahead`, `behind`)
  - `StatusOutputModeType` (`table`, `json`, `compact`)
  - `StatusTableColumnType` (`Repository`, `Branch`, `Status`, `Sync`, `Stash`, `Files`)
  - `StatusRepoStateType` (`clean`, `dirty`, `ahead`, `behind`, `stashed`, `missing`)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmdstatus/constants/`):**
  - Command names (`CmdStatus`, `CmdStatusAlias`), table tokens, gap tokens (`StatusTableColGap`), banners, and suggestions.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring:**
  - Refactored `status.go`, `status_render.go`, `status_help_menu.go`, `helpers_shim.go`, and `flags_test.go`.

### 4. `cli/cmdlogin` (Subtask 320-04)
- **Leaf Enums (`cli/cmdlogin/enums/`):**
  - `LoginAuthMethodType` (`ssh`, `token`, `browser`, `status`)
  - `LoginInteractiveChoiceType` (`1`, `2`, `3`, `4`)
  - `GitAuthProviderType` (`github`, `gitlab`)
  - `ClipboardToolType` (`clip`, `pbcopy`, `wl-copy`, `xclip`, `xsel`)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmdlogin/constants/`):**
  - Subcommands, flags, Git config keys, endpoints (`GithubUserAPI`, etc.), and tool arguments.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring:**
  - Refactored `login_cmd.go`, `login_browser.go`, `login_ssh.go`, `login_status.go`, `login_token.go`, and created `login_test.go`.

### 5. `cli/cmdcargo` (Subtask 320-05)
- **Leaf Enums (`cli/cmdcargo/enums/`):**
  - `CargoSubcommandType` (`status`, `st`, `info`, `install`, `in`)
  - `CargoCommandType` (`build`, `check`, `test`, `clippy`, `run`, `clean`, `fmt`, `doc`, `update`, `init`, `new`, `bench`, `publish`)
  - `RustToolchainBinaryType` (`cargo`, `rustc`, `rustup`, `clang`)
  - `ToolchainEnvVarType` (`LIBCLANG_PATH`, `PATH`, `CARGO_HOME`, `RUSTUP_HOME`, etc.)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmdcargo/constants/`):**
  - CLI flags, binary names, path tokens, error code `ErrCodeMissingCargo = "E7100"`, Dracula UI colors, install suggestions.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring & Error Sovereignty:**
  - Refactored `cargo.go`, `cargo_flags.go`, `cargo_resolve.go`, `cargo_status.go`, `cargo_suggest.go`, `cargo_test.go`.
  - In `cargo_exec.go`, eliminated direct `cliexit.Exit()` call in `handleCargoExitError`, returning `error` directly to the outer dispatcher.
  - Decomposed `enrichCargoEnv` into functions $\le 11$ lines.

---

## Architectural Invariants Verified
- **Enum Suffix:** All 19 created enum types end in `Type`.
- **Zero Cyclic Dependencies:** All leaf packages (`enums/`, `constants/`) import only standard library packages.
- **Positive Booleans:** Zero negative boolean variables or parameters.
- **Function Limit:** All functions $\le 15$ lines ($\le 8$ preferred). Blank lines after `}` and before `return`.
- **FastGate Validation:** 0.054s execution, 100% PASS on relative paths, nested ifs, booleans, and file sizes.
- **Test Suite Pass:** 100% pass across all 15 subpackages (`go test -count=1`) and root `cli/cmd` test suite.
