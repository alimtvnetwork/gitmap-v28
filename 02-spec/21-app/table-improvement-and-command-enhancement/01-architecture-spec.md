# Architecture Specification: Pending Commits Table Improvement & New Commands Discovery Engine

> **Spec Version:** 1.0.0  
> **Status:** Approved / Grounded  
> **Task Slug:** `table-improvement-and-command-enhancement`  
> **Target Path:** `02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md`  

---

## 1. System Architecture Context & Problem Statement

### 1.1 Existing Pending Commits UX Shortcomings
GitMap provides the `gitmap pending-commits` (alias: `pc`) command implemented in `cli/cmdpending/pending_commits_cmd.go`. While functional, the current terminal table rendering exhibits several critical architectural and usability deficiencies:

1. **Fragmented Column Allocation & Screen Real Estate Waste:**
   - The current table allocates four separate columns for file status: `DIRTY` (5 chars), `UNTRACK` (7 chars), `MODIF` (6 chars), and `STAGE` (6 chars), totaling 24+ characters of terminal width.
   - For repository health assessment, developers and automated agents primarily need to know whether uncommitted files exist and the aggregate change count. Fragmenting these numbers forces aggressive truncation on repository names (capped at 20 chars) and branch names (capped at 8 chars).

2. **Absence of Repository Version Identification:**
   - The table displays only `BRANCH` (truncated to 8 characters, e.g. `main`, `feat-x…`), giving zero visibility into which release milestone or SemVer tag the local repository is anchored to (e.g., `v6.523.1` or `v6.524.0`).
   - When managing multi-repo fleets or versioned submodules, developers must run secondary commands to identify whether a dirty repository is tracking a released tag or an unversioned feature branch.

3. **Missing In-Place Tree-View Remediation Hints:**
   - When a repository is identified as dirty or pending, developers must manually formulate remediation commands.
   - There are no inline contextual suggestions for quick resolution (e.g., fast commit and push vs. stash or discard).

4. **Lack of Workspace Fleet Batch Guidance:**
   - When 5, 10, or 20 repositories have pending uncommitted changes across a workspace, fixing them repo-by-repo is slow and error-prone.
   - GitMap already possesses the enterprise fleet commit-and-push engine `gitmap cpar` (`commit-push-all-repos`), but `gitmap pending-commits` never surfaces this batch capability in its table footer.

### 1.2 The Need for `gitmap new-commands` (`gitmap nc`)
GitMap has grown through dozens of minor and major releases, accumulating over 100 enterprise commands spanning fleet SSH delegation, split-DB indexing, CI/CD diagnostics, automation units, and workspace sync. However:
- Discoverability of recently added commands is low; developers rely on memory or fragmented help text.
- Standard help lists all commands alphabetically or in monolithic clusters without highlighting recent additions.
- Developers require a focused command: `gitmap new-commands` (and short alias `gitmap nc`) that filters the last 100 new commands, provides concise descriptions, version metadata, and practical copy-pasteable CLI examples.

---

## 2. Pending Commits Table Improvement Architecture

### 2.1 Unified `UNCOMMITTED` Column
Instead of splitting file states into `DIRTY`, `UNTRACK`, `MODIF`, and `STAGE`, the table replaces these four columns with a single unified `UNCOMMITTED` column.
- **Metric Calculation:** `TotalUncommitted = UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`.
- **Display Formatting:**
  - When `TotalUncommitted == 0`: Formatted as `-` or `0` in dimmed or green text.
  - When `TotalUncommitted > 0`: Formatted as clean numeric count `N` (colored yellow).
- **Space Savings:** Reclaims 15+ columns of horizontal space to expand repository names and version strings.

### 2.2 Repository Short Version and Branch (`VER/BRANCH`)
The improved table introduces a combined `VER/BRANCH` column:
- **Version Discovery Logic:**
  1. Inspect local repository Git tags via `git -C <dir> describe --tags --abbrev=0`.
  2. Fall back to reading `.gitmap/version.json` or `version.json` if present.
  3. If no tag or version file is found, version is empty string `""`.
- **Formatting Algorithm (`formatShortVersionBranch`):**
  - If version is present: `<version>/<branch>` (e.g. `v6.523/main`, `v1.4.0/develop`).
  - If version is absent: `<branch>` (e.g. `main`, `master`, `feat-login`).
  - If combined string exceeds column width (16–18 characters), apply smart abbreviation keeping the version intact and truncating the branch with `…` (e.g. `v6.523/feat-auth…`).

### 2.3 Hierarchical Tree-View Remediation Hints
Beneath each dirty repository row in the terminal table, GitMap renders a two-line hierarchical tree view showing concrete, actionable remediation options:
- **Branch Symbols:**
  - Line 1: `├── Option 1: <command>`
  - Line 2: `└── Option 2: <command>`
- **Option 1 (Commit & Push Local WIP):**
  - Targets committing and pushing the uncommitted working tree.
  - Command: `git -C "<repoPath>" add -A && git commit -m "wip: save changes" && git push` or `gitmap cpf "wip: save changes"`.
- **Option 2 (Stash / Isolate Local Changes):**
  - Targets stashing uncommitted changes cleanly.
  - Command: `git -C "<repoPath>" stash -u`.
- **Indentation & Alignment:** Tree lines are indented beneath the table border using the table border style (`│   ├── Option 1: ...`), ensuring the table bounding box remains rectangular, visually aligned, and clean.

### 2.4 Overall Fleet Batch Command in Table Footer
In the summary footer box of the pending commits table, GitMap prominently features the workspace fleet remediation command:
```
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ Fleet Remediation: gitmap cpar "wip: save changes"                           │
  └──────────────────────────────────────────────────────────────────────────────┘
```
This informs the developer immediately that they can batch-commit and push all pending changes across all dirty repositories in a single atomic operation without manual per-repo intervention.

### 2.5 Table Layout & Column Alignment Specifications
The bounding box width is standardized at 78 inner characters (80 characters total including border edges `│`), preserving full cross-platform terminal compatibility:

```
  ┌──────────────────────────────────────────────────────────────────────────────┐
  │                    GITMAP PENDING COMMITS SUMMARY                            │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ Scanned: 42 repos    │ Dirty: 2 repos   │ Uncommitted: 15 files │ Unpushed: 1 │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ REPOSITORY             VER/BRANCH         UNCOMMITTED   UNPUSHED   STATUS    │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ gitmap                 v6.523.1/main               12          1   ● PEND    │
  │   ├── Option 1: git -C "gitmap" add -A && git commit -m "wip"                │
  │   └── Option 2: git -C "gitmap" stash -u                                     │
  │ movie-cli-v8           v8.0.0/main                  3          0   ● PEND    │
  │   ├── Option 1: git -C "movie-cli-v8" add -A && git commit -m "wip"          │
  │   └── Option 2: git -C "movie-cli-v8" stash -u                               │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ 40 clean repos         -                            0          0   ○ CLEAN   │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ Fleet Remediation: gitmap cpar "wip: save changes"                           │
  └──────────────────────────────────────────────────────────────────────────────┘
```

Width specifications:
- `REPOSITORY`: 22 characters (`%-22s`)
- `VER/BRANCH`: 16 characters (`%-16s`)
- `UNCOMMITTED`: 12 characters (`%12d`)
- `UNPUSHED`: 10 characters (`%10d`)
- `STATUS`: 9 characters (`%-9s`)

---

## 3. `gitmap new-commands` (`gitmap nc`) Architectural Design

### 3.1 Overview & Command Routing
The `new-commands` command provides a curated discovery engine for the last 100 newly introduced commands across recent GitMap releases.
- **Primary Command:** `gitmap new-commands`
- **Short Alias:** `gitmap nc`
- **CLI Registration:** Registered in `cli/cmd/roottooling.go` inside `toolingDevEntries()` and constants defined in `cli/constants/constants_cli.go`.

### 3.2 Data Models (`cli/cmdpending/new_commands_types.go`)
```go
package cmdpending

import "time"

// NewCommandEntry represents a single cataloged command entry.
type NewCommandEntry struct {
	Name        string   `json:"name"`
	Alias       string   `json:"alias,omitempty"`
	Version     string   `json:"version"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Example     string   `json:"example"`
	Flags       []string `json:"flags,omitempty"`
}

// NewCommandsPayload defines the top-level structured telemetry for JSON output.
type NewCommandsPayload struct {
	Timestamp      time.Time         `json:"timestamp"`
	TotalCataloged int               `json:"totalCataloged"`
	ReturnedCount  int               `json:"returnedCount"`
	CategoryFilter string            `json:"categoryFilter,omitempty"`
	Commands       []NewCommandEntry `json:"commands"`
}

// NewCommandsOptions encapsulates CLI flags for the new-commands command.
type NewCommandsOptions struct {
	Limit          int
	FilterQuery    string
	CategoryFilter string
	IsJSON         bool
}
```

### 3.3 Command Catalog Engine (`cli/cmdpending/new_commands_cmd.go`)
The catalog maintains a structured, chronologically sorted registry of the last 100 new commands introduced across GitMap milestones:
1. **Fleet & Workspace Batch Operations:**
   - `cpar` (Commit and push all repositories in workspace): `gitmap cpar "wip: save changes"`
   - `commons` / `co` (Apply standardized repo configurations): `gitmap commons`
   - `space` (Workspace backup branches and baselines): `gitmap space backup-branch "feat-login"`
2. **Pipeline Telemetry & Diagnostics:**
   - `pe` (Pipeline error analyzer with dynamic runner extraction): `gitmap pe -t`
   - `sug` (Shutdown execution until quality gates are green): `gitmap sug`
   - `rerun` / `rr` (Rerun failed pipeline jobs): `gitmap rerun --step=lint`
3. **Commit & Pull Modernization:**
   - `pending-commits` / `pc` (Inspect dirty working trees and unpushed commits): `gitmap pc`
   - `cpf` (Atomic commit, push, feature): `gitmap cpf "feat: user authentication"`
   - `cpb` (Atomic commit, push, bugfix): `gitmap cpb "fix: buffer overflow"`
   - `cpr` (Commit, push, release workflow): `gitmap cpr "v6.524.0"`
   - `cin` (Commit in replay): `gitmap cin --source repo-a --dest repo-b`
4. **AI & Automation Systems:**
   - `agy` (Antigravity developer tools and queue): `gitmap agy deploy`
   - `agm` (Antigravity fleet agent manager): `gitmap agm status`
   - `aum` (Automation Unit Manager): `gitmap aum run 06-cicd-local-runner.py`
5. **System & OS Integration:**
   - `os dock` (Cross-platform taskbar/dock positioner): `gitmap os dock bottom`
   - `apps` (Installed application auditor and uninstaller): `gitmap apps list`
   - `fix-link` (Symlink and junction repairs): `gitmap fix-link --repair`

### 3.4 CLI Flags & Filtering Capabilities
- `--limit N` (default: 100): Cap output to the most recent `N` commands.
- `--category <cat>`: Filter by functional category (e.g., `commits`, `pipeline`, `fleet`, `ai`, `os`).
- `--filter <query>` / `-q <query>`: Text substring search matching name, alias, description, or example.
- `--json`: Emit complete structured machine-readable payload.
- `-h, --help`: Display contextual help card with examples and flag guide.

---

## 4. Acceptance Criteria

### AC1: Pending Commits Unified Column
- The table output replaces the four disparate columns (`DIRTY`, `UNTRACK`, `MODIF`, `STAGE`) with a single consolidated `UNCOMMITTED` column.
- The value represents `UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`.
- The summary row accurately aggregates total uncommitted files across all inspected repositories.

### AC2: Version/Branch Identification
- The table displays the combined `VER/BRANCH` column for every repository row.
- Repositories with valid tags display `<tag>/<branch>` (e.g. `v6.523.1/main`).
- Repositories without tags display `<branch>`.
- The column gracefully truncates long names with `…` without breaking box border alignment.

### AC4: Hierarchical Tree-View Remediation & Fleet Footer
- Every dirty repository row displays two hierarchical tree branches immediately below:
  - `├── Option 1: <commit-and-push-command>`
  - `└── Option 2: <stash-or-discard-command>`
- Clean repositories do not render tree-view hints.
- The table footer displays the fleet-wide batch command:
  `Fleet Remediation: gitmap cpar "wip: save changes"`

### AC8: New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)
- Both `gitmap new-commands` and `gitmap nc` invoke the command catalog engine.
- Supports filtering the last 100 commands by default.
- Every entry displays command name, alias, version introduced, description, and practical CLI example.
- Supports `--json` flag producing valid `NewCommandsPayload` JSON.
- Supports `--limit`, `--category`, and `--filter` flags.

---

## 5. Coding Guidelines & Implementation Invariants

### 5.1 Affirmative Boolean Naming
All boolean fields, variables, parameters, and return values MUST strictly use affirmative naming with `is*` or `has*` prefixes:
- Permitted: `isDirty`, `isClean`, `hasUncommitted`, `hasUpstream`, `isJSON`, `isHelpRequested`, `hasRemediation`.
- Strictly Forbidden: `notClean`, `noDirty`, `uncommitted`, `disableTree`, `skipFooter`.

### 5.2 Structured Error Handling
All operations returning errors MUST wrap them with `apperror.AppError` and standard error codes:
- `E9001`: Directory inspection failure.
- `E9002`: Repository target resolution error.
- `E9003`: Invalid command filter or limit argument.

### 5.3 Function Sizing Guard (8–15 LOC Cap)
Every function MUST strictly adhere to the 8–15 line length cap:
- Decompose rendering logic into discrete helpers (`renderTableHeader`, `renderRepoRow`, `renderTreeHints`, `renderTableFooter`).
- Separate CLI argument parsing, model transformation, and output rendering into dedicated micro-functions.

### 5.4 Strict Relative Paths
All file references and error messages MUST use relative paths relative to repository root. Absolute filesystem paths (e.g. `C:\...`) and `file:///` URIs are strictly prohibited.
