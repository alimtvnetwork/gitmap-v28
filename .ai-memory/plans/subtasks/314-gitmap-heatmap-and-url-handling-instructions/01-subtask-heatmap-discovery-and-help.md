# Subtask 01: Heatmap Database Discovery Scanner, Repo Alias Enrichment & Flag Injection Routing Fix

## Status: READY
- **Task ID:** `314-subtask-01`
- **Parent Spec:** [02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/01-architecture-spec.md](../../../../02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/01-architecture-spec.md)
- **Module:** `cli/heatmap`, `cli/cmd`
- **Owned Files:**
  - `cli/cmd/special_repos_rc_heatmap.go`
  - `cli/heatmap/heatmap_discovery.go`
  - `cli/heatmap/heatmap_overview.go`
  - `cli/heatmap/heatmap_cmd.go`
  - `cli/heatmap/types.go`

---

## 1. Objectives

1. **Remediate Flag Injection Defect in `cli/cmd/special_repos_rc_heatmap.go`:**
   - Eliminate the routing failure where `enrichRCHeatmapFlags(args)` prepended or appended `--repo` and `--ai` into the subcommand position, resulting in `unknown rc subcommand: --repo`.
   - Implement `parseRcCommandTokens(args []string)` to cleanly isolate positional subcommands from CLI flag switches.
   - When no subcommand is specified (`gitmap rc heatmap`), route directly to the new `heatmap.RunHeatmapOverview`.
   - Guarantee that subcommands (`sync`, `status`, `commit`, `tree`, `search`, etc.) always occupy index 0 when forwarded to `heatmap.DispatchHeatmapSubsystem`.

2. **Implement 3-Tier Database Discovery Scanner in `cli/heatmap/heatmap_discovery.go`:**
   - Scan Tier 1 Local Root: `<effectiveDataDir>/heatmap.db` via `store.EffectiveDataDir()`.
   - Scan Tier 2 Local Split DBs: `<effectiveDataDir>/heatmap/<repoSlug>/sql.db` via directory traversal.
   - Scan Tier 3 Multi-AI Repo-Cache DBs: `<repo-cache>/.heatmap/<seq>-<ai>/` inspecting both `root-heatmap.db` and `<repoSlug>/heatmap.db`.
   - Query file metadata (size in bytes, formatted size, modification timestamp) and total record count (`SELECT COUNT(*) FROM HeatmapActionLog`) using read-only connections (`mode=ro`).

3. **Implement Overview & Repository Alias Enrichment in `cli/heatmap/heatmap_overview.go`:**
   - Query primary database via `db.ListAliasesWithRepo()` to match discovered repository slugs against registered aliases.
   - Attach aliases to `HeatmapDbEntry.Aliases` and populate `HeatmapDbEntry.PrimaryAlias`.
   - Provide dual output modes: formatted human-readable ASCII table (default) and atomic JSON payload (when `--json` flag is set).

4. **Wire Overview Subcommand in `cli/heatmap/heatmap_cmd.go`:**
   - Register `overview` as a recognized command in `dispatchHeatmapAdmin`.
   - When `gitmap heatmap` is run with no arguments or flags, display the combined overview and recent tasks summary rather than a bare help screen.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Disentangle Flag Routing in `cli/cmd/special_repos_rc_heatmap.go`
- Introduce helper `parseRcCommandTokens(args []string) (string, []string, bool)`:
  - Iterate through tokens; classify tokens starting with `-` as flags, and the first non-flag token as the subcommand.
- Refactor `RunRCHeatmapCommand(args []string) error`:
  - If `!hasSubcmd`: dispatch directly to `heatmap.RunHeatmapOverview(os.Stdout, flagMap)`.
  - If `hasSubcmd`: pass `subcmd` at index 0, followed by enriched flags.
- Verify that every refactored function stays strictly $\le 15$ lines.

### Step 2: Construct `cli/heatmap/heatmap_discovery.go`
- Implement `DiscoverAllHeatmapDatabases() ([]HeatmapDbEntry, error)`:
  - Call `scanTier1RootDb(dataDir)`.
  - Call `scanTier2SplitDbs(dataDir)`.
  - Call `scanTier3RepoCacheDbs(repoCacheRoot)`.
- Implement `inspectDatabaseStats(dbPath string, tier HeatmapTierType, slug, ai string) (HeatmapDbEntry, error)`:
  - Stat file size and mtime.
  - Open SQLite connection with `file:<path>?mode=ro&_busy_timeout=3000`.
  - Query `SELECT COUNT(*) FROM HeatmapActionLog`. If the table does not exist or fails, return 0 rows gracefully without terminating the scanner.
- Implement `formatFileSize(bytes int64) string` returning clean `B`, `KB`, or `MB` strings.

### Step 3: Author `cli/heatmap/heatmap_overview.go`
- Implement `RunHeatmapOverview(out io.Writer, flags map[string]string) error`:
  - Invoke `DiscoverAllHeatmapDatabases()`.
  - Invoke `enrichEntriesWithAliases(entries)`.
  - If `flags["json"] != ""` or `flags["format"] == "json"`, emit `renderOverviewJSON(out, report)`.
  - Otherwise, emit `renderOverviewText(out, report)`.
- Implement `enrichEntriesWithAliases(entries []HeatmapDbEntry) []HeatmapDbEntry`:
  - Open primary database via `store.OpenDefault()`.
  - Retrieve aliases via `db.ListAliasesWithRepo()`.
  - Match each entry's `RepoSlug` against alias definitions.

### Step 4: Hook Dispatcher & Update Help
- In `cli/heatmap/heatmap_cmd.go`:
  - Update `dispatchHeatmapAdmin` to handle `"overview"`.
  - Update `isHeatmapHelpRequested(args []string)` to allow bare `gitmap heatmap` to route to overview.
  - Update `printHeatmapHelp(out io.Writer)` to document the `overview` command and alias enrichment.

---

## 3. Strict Coding Guidelines & Constraints

- **Function Line Cap:** Every function MUST be $\le 15$ lines (preferred $\le 8$ lines).
- **Positive Booleans Only:** Use `isAccessible`, `hasSubcmd`, `isPrimary`, `isJsonFormat`. Never use negative booleans (`isNotValid`, `unaliased`).
- **Zero Nested Conditionals:** All conditionals must be flattened with early returns and guard clauses.
- **Error Handling:** Wrap all underlying SQLite or I/O errors using `*apperror.AppError` via `apperror.WrapSimple(err, "context")`.
- **Relative Path Hygiene:** Zero hardcoded drive letters (`C:`, `D:`) or URI schemes (`file:///`). Forward slashes only.
- **Zero Builds / Zero Tests:** Do not execute `go build` or `go test` during this authoring phase.

---

## 4. Verification Evidence & Acceptance Criteria

1. `gitmap rc heatmap` (bare command) successfully executes without error and displays the Heatmap Overview rather than `unknown rc subcommand: --repo`.
2. `gitmap rc heatmap --repo gitmap status` correctly dispatches `status` without argument displacement.
3. Discovered databases show accurate tier classification (`local_root`, `local_split_repo`, `repo_cache_root`, `repo_cache_repo`).
4. Repositories with registered aliases in `gitmap.db` display their primary alias badge.
5. All newly authored and modified Go files pass line limit checks ($\le 15$ lines) and negative boolean audit.
