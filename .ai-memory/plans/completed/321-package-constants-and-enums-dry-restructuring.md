# Plan 321: Package Constants & Enums DRY Restructuring (Phase 3: 5 Core Packages)

## Executive Summary
This milestone executes the strict `01-prompts/15-cg-execute/06-constants-and-enums.md` coding guideline protocol across 5 central CLI packages in `cli/`:
1. `cli/cmdconfig` (Configuration Management Engine — 12 files)
2. `cli/cmddiff` (Git Diff & Workspace Comparison — 6 files)
3. `cli/cmdsearch` (Fast Workspace & Repository Search — 5 files)
4. `cli/cmdbranch` (Git Branch Operations & Inspection — 3 files)
5. `cli/cmdhealth` (System & Diagnostic Health Checks — 2 files)

Each package received:
- Dedicated leaf package `enums/` with strongly-typed enums ending in `Type`, deterministic methods (`IsValid()`, `String()`, `Parse*()`, `All*()`), positive booleans (`is*`, `has*`), and comprehensive unit tests (`enums_test.go`).
- Dedicated leaf package `constants/` with zero magic strings, centralized format tokens, default values, command identifiers, and unit tests (`constants_test.go`).
- Refactored caller files in each parent package eliminating magic string comparisons, raw literals, and repetitive token definitions.
- FastGate and `check-enum-guidelines.py` pre-commit validation passing with 100% green status.

---

## Technical Outcomes by Package

### 1. `cli/cmdconfig` (Subtask 321-01)
- **Leaf Enums (`cli/cmdconfig/enums/`):**
  - `ConfigActionType` (`export`, `import`, `what-configs`, `which-format`, `get`, `set`, `list`, `reset`)
  - `ConfigOutputFormatType` (`table`, `json`, `yaml`, `compact`)
  - `ConfigScopeType` (`local`, `global`, `system`, `project`)
  - `ConfigKeyType` (`target_dir`, `default_tool`, `export_dir`, `schema_version`, `encoding`, `auto_import`)
  - `ConfigToolType` (`vscode`, `qtorrent`, `utorrent`, `all`)
  - `ConfigFileFormatType` (`json`, `yaml`, `db`, `sqlite`, `toml`, `env`, `ini`, `cfg`, `conf`)
  - `ConfigFileEncodingType` (`utf-8`, `base64`)
  - Unit tests in `enums_test.go` with 100% statement coverage.
- **Leaf Constants (`cli/cmdconfig/constants/`):**
  - Subcommands, flags, OS tokens (`OSWindows`, `OSDarwin`, `OSLinux`), tool subpaths, canonical file names, table column headers, starter templates, error codes (`E7201`-`E7204`).
  - Unit tests in `constants_test.go`.
- **Caller Refactoring:**
  - Refactored `config_model.go`, `config_tool_paths.go`, `config_export.go`, `config_export_writers.go`, `config_import.go`, `config_import_readers.go`, `what_configs_cmd.go`, `which_format_cmd.go`, and test fixtures.

### 2. `cli/cmddiff` (Subtask 321-02)
- **Leaf Enums (`cli/cmddiff/enums/`):**
  - `DiffViewModeType` (`unified`, `side-by-side`, `stat`, `summary`, `compact`)
  - `DiffOutputFormatType` (`text`, `json`, `raw`)
  - `DiffScopeType` (`working-tree`, `staged`, `branch`, `commit`)
  - `DiffFilterType` (`all`, `only-conflicts`, `only-missing`, `include-identical`)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmddiff/constants/`):**
  - Subcommands, flags (`--stat`, `--cached`, `--staged`, `--name-only`, `--json`), help text, layout tokens, error wrap tokens, positional limits (`ExpectedDiffPositionalArgs = 2`), exit code `ExitUsageError = 2`.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring & Exit Code Fix:**
  - In `diff.go:63`, replaced raw integer exit code `cliexit.HandleError(err, 2)` with `cliexit.HandleError(err, constants.ExitUsageError)`.
  - Refactored `diff_help_menu.go`, `dispatchcompare.go`, `dispatchdiff.go`, and created `cmddiff_test.go`.

### 3. `cli/cmdsearch` (Subtask 321-03)
- **Leaf Enums (`cli/cmdsearch/enums/`):**
  - `SearchScopeType` (`content`, `symbol`, `path`, `commit`, `repo`)
  - `SearchMatchModeType` (`keyword`, `regex`, `exact`, `fuzzy`)
  - `SearchOutputFormatType` (`text`, `json`, `table`, `compact`)
  - `SearchCategoryType` (`user`, `ai`)
  - `SearchActionType` (`query`, `history`, `top`, `stats`, `dh2d`, `clean`, `clear`, `purge`, `reset`)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmdsearch/constants/`):**
  - Commands, flags, default limits (`DefaultLocalSearchLimit = 20`), ignored directories (`node_modules`, `vendor`, `bin`), output path (`.gitmap/output/repo_search`), empty JSON array `"[]"`, format strings, error messages.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring:**
  - Refactored `search.go`, `search_entry.go`, `search_help_menu.go`, `helpers_shim.go`, and added `search_test.go`. Flattened conditional branches to satisfy zero nested ifs.

### 4. `cli/cmdbranch` (Subtask 321-04)
- **Leaf Enums (`cli/cmdbranch/enums/`):**
  - `BranchFilterType` (`all`, `local`, `remote`, `merged`, `unmerged`)
  - `BranchActionType` (`default`, `def`, `list`, `current`, `latest-branch`, `lb`, `help`)
  - `BranchSortOrderType` (`recent`, `committerdate`, `authordate`, `alphabetical`)
  - `BranchOutputFormatType` (`terminal`, `json`, `csv`)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmdbranch/constants/`):**
  - Command names & aliases (`CmdBranch`, `CmdBranchAlias`, `CmdSubDefault`, `CmdSubDefaultAlias`, `CmdSubLatest`), flags, usage lines, tips, error codes (`ErrCodeE9000`), current dir sentinel (`"."`).
  - Unit tests in `constants_test.go`.
- **Caller Refactoring:**
  - Refactored `branch.go`, `branch_help_menu.go`, `hook_shim.go`, and created `branch_test.go`.

### 5. `cli/cmdhealth` (Subtask 321-05)
- **Leaf Enums (`cli/cmdhealth/enums/`):**
  - `HealthCategoryType` (`git`, `pipeline`, `fleet`, `db`, `system`, `network`)
  - `HealthSeverityType` (`ok`, `warn`, `error`, `fatal`)
  - `HealthOutputFormatType` (`terminal`, `json`, `summary`, `markdown`)
  - `HealthCheckStatusType` (`green`, `yellow`, `red`, `skipped`)
  - Unit tests in `enums_test.go` with 100% coverage.
- **Leaf Constants (`cli/cmdhealth/constants/`):**
  - Command names (`CmdHealth`, `CmdHealthAlias`), environment variable `EnvGitmapWorkdir`, flags, limits/thresholds (`MaxDirtyReposDisplay = 5`, `MaxRedPipelinesDisplay = 5`, `MaxPipelineProbedRepos = 50`, `StaleDbDaysCutoff = 30`), headings, formatters, diagnostic notes, follow-up triage recommendations.
  - Unit tests in `constants_test.go`.
- **Caller Refactoring:**
  - Refactored `health.go`, `health_checks.go`, and created `health_test.go`.

---

## Architectural Invariants Verified
- **Enum Suffix:** All 23 created enum types end in `Type`.
- **Zero Cyclic Dependencies:** All leaf packages (`enums/`, `constants/`) import only standard library packages.
- **Positive Booleans:** Zero negative boolean variables or parameters.
- **Function Limit:** All functions $\le 15$ lines ($\le 8$ preferred). Blank lines after `}` and before `return`.
- **FastGate Validation:** 0.047s execution, 100% PASS on relative paths, nested ifs, booleans, and file sizes.
- **Test Suite Pass:** 100% pass across all 15 subpackages (`go test -count=1`).
