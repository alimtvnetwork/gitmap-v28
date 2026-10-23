# Architecture Specification: Heatmap Database Discovery, Terminal Suggestions, and AI-Optimized Help

> **Module:** `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions`  
> **Sequence ID:** `314`  
> **Status:** `Draft / Implementation-Ready`  
> **Subtask 01:** [.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/01-subtask-heatmap-discovery-and-help.md](../../../.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/01-subtask-heatmap-discovery-and-help.md)  
> **Subtask 02:** [.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/02-subtask-recent-tasks-and-suggestions.md](../../../.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/02-subtask-recent-tasks-and-suggestions.md)  
> **Version:** `1.0.0`  
> **Updated:** `2026-10-23`  

---

## 1. Executive Summary

This specification establishes the system architecture for GitMap's Heatmap Database Discovery, Cross-Repository Recent Tasks Aggregation, Dynamic AI Terminal Suggestions, and the resolution of the command flag injection routing defect:

1. **3-Tier Heatmap Database Topology:** Formalizes the tripartite SQLite storage layout across Local Workstation Root (`<effectiveDataDir>/heatmap.db`), Local Per-Repository Split DBs (`<effectiveDataDir>/heatmap/<repoSlug>/sql.db`), and Multi-AI Repo-Cache Telemetry (`<repo-cache>/.heatmap/<seq>-<ai>/`).
2. **Comprehensive Database Discovery Scanner (`cli/heatmap/heatmap_discovery.go`):** Implements automated filesystem introspection identifying all active heatmap databases, file sizes, last modified timestamps, and total action counts across all three tiers.
3. **Repository Alias Enrichment Engine:** Joins discovered database paths against `gitmap.db`'s `Alias` registry via `db.ListAliasesWithRepo()`, presenting human-readable aliases alongside raw repository slugs.
4. **Flag Injection Routing Bug Remediation:** Eliminates the defect in `cli/cmd/special_repos_rc_heatmap.go:RunRCHeatmapCommand` where `--repo` and `--ai` flags were mistakenly prepended as subcommands, breaking subcommand dispatch.
5. **Multi-Repository Recent Tasks Aggregator (`cli/heatmap/heatmap_recent_tasks.go`):** Queries top 5 to 10 recent tasks across 2 to 5 repositories from `HeatmapActionLog`, calculating completion rates, action volumes, and reasoning state distributions.
6. **Dynamic Terminal Suggestions Engine (`cli/heatmap/heatmap_suggestions.go`):** Synthesizes live database state to emit high-leverage terminal command recommendations for AI agents and human developers.

---

## 2. 3-Tier Heatmap Database Architecture

GitMap organizes cognitive telemetry and action traces into three isolated, purpose-built SQLite tiers to balance low-latency local logging, per-project encapsulation, and multi-agent coordination.

```
+---------------------------------------------------------------------------------------+
|                                3-Tier Heatmap Storage                                 |
+---------------------------------------------------------------------------------------+
|  Tier 1: Local Workstation Root DB                                                    |
|  Location: <effectiveDataDir>/heatmap.db                                              |
|  Scope: Global workstation coordination, uncommitted action delta queues              |
+---------------------------------------------------------------------------------------+
|  Tier 2: Local Per-Repository Split DBs                                               |
|  Location: <effectiveDataDir>/heatmap/<repoSlug>/sql.db                               |
|  Scope: Local repository-specific action trees, file touches, and task traces         |
+---------------------------------------------------------------------------------------+
|  Tier 3: Multi-AI Repo-Cache Telemetry DBs                                            |
|  Location: <repo-cache>/.heatmap/<seq>-<ai>/                                          |
|    - root-heatmap.db                                 (AI global registry)             |
|    - <repoSlug>/heatmap.db                           (AI per-repo telemetry)          |
|  Scope: Multi-agent shared sync, Git-committed audit trails, team collaboration       |
+---------------------------------------------------------------------------------------+
```

### 2.1 Tier 1: Local Workstation Root Database
- **Resolved Path:** `<effectiveDataDir>/heatmap.db` via `store.EffectiveDataDir()`.
- **Primary Responsibility:** Coordinates local workstation actions, tracks overall execution state, manages uncommitted mutation counters, and maintains local sync state with remote repositories.
- **SQLite Configuration Pragmas:**
  - `PRAGMA journal_mode = WAL;`
  - `PRAGMA busy_timeout = 5000;`
  - `PRAGMA synchronous = NORMAL;`
  - `PRAGMA foreign_keys = ON;`

### 2.2 Tier 2: Local Per-Repository Split Databases
- **Resolved Path:** `<effectiveDataDir>/heatmap/<repoSlug>/sql.db` via `store.ResolveHeatmapRepoSplitDbPath(repoSlug, repoRoot)`.
- **Primary Responsibility:** Isolates large action log volumes per repository. Prevents a single high-churn project from bloating global data or causing lock contention. Stores local execution trees, file touch counts, and task summaries.
- **Directory Layout:**
  ```text
  <effectiveDataDir>/heatmap/
    +-- gitmap/
    |   +-- sql.db
    +-- repo-secrets/
    |   +-- sql.db
    +-- riseup-core/
        +-- sql.db
  ```

### 2.3 Tier 3: Multi-AI Repo-Cache Telemetry Databases
- **Resolved Path:** `<repo-cache>/.heatmap/<seq>-<ai>/` via `store.ResolveRcHeatmapAiDir(aiName)`.
- **Primary Responsibility:** Serves as the Git-tracked repository coordination medium. Each AI agent family operates inside a partitioned sequential namespace:
  - `01-antigravity`: Google Antigravity IDE and subagents
  - `02-muse`: Muse autonomous orchestration agents
  - `03-cloud-code`: Cloud Code developer sessions
  - `04-codex`: OpenAI / GitHub Copilot agents
  - `05-grok`: Grok reasoning agents
  - `06-cursor`: Cursor IDE sessions
- **Storage Artifacts per AI Namespace:**
  - `root-heatmap.db`: High-level session logs and roadmap milestones.
  - `<repoSlug>/heatmap.db`: Project-specific cognitive traces, reasoning states, and verification results.

---

## 3. Comprehensive Database Discovery Scanner

The discovery scanner in `cli/heatmap/heatmap_discovery.go` systematically surveys all three tiers on the active host machine.

### 3.1 Scanner Data Contracts

```go
package heatmap

import (
	"time"
)

// HeatmapTierType designates the architectural tier of a discovered database.
type HeatmapTierType string

const (
	TierLocalRoot      HeatmapTierType = "local_root"
	TierLocalSplitRepo HeatmapTierType = "local_split_repo"
	TierRepoCacheRoot  HeatmapTierType = "repo_cache_root"
	TierRepoCacheRepo  HeatmapTierType = "repo_cache_repo"
)

// HeatmapDbEntry contains verified metadata for an individual heatmap database.
type HeatmapDbEntry struct {
	TierType       HeatmapTierType `json:"tierType"`
	RepoSlug       string          `json:"repoSlug"`
	AiNamespace    string          `json:"aiNamespace,omitempty"`
	Aliases        []string        `json:"aliases"`
	PrimaryAlias   string          `json:"primaryAlias,omitempty"`
	RelativePath   string          `json:"relativePath"`
	FileSizeBytes  int64           `json:"fileSizeBytes"`
	FormattedSize  string          `json:"formattedSize"`
	TotalActions   int64           `json:"totalActions"`
	LastModifiedAt time.Time       `json:"lastModifiedAt"`
	IsAccessible   bool            `json:"isAccessible"`
}

// HeatmapDiscoveryReport bundles discovered databases and aggregation metrics.
type HeatmapDiscoveryReport struct {
	ScannedAt      time.Time        `json:"scannedAt"`
	TotalDatabases int              `json:"totalDatabases"`
	TotalActions   int64            `json:"totalActions"`
	Databases      []HeatmapDbEntry `json:"databases"`
}
```

### 3.2 Discovery Scanning Flow
1. **Tier 1 Probe:** Check existence of `<effectiveDataDir>/heatmap.db`. If present, record size, modification timestamp, and run `SELECT COUNT(*) FROM HeatmapActionLog`.
2. **Tier 2 Discovery:** Enumerate subdirectories of `<effectiveDataDir>/heatmap/`. For each subfolder containing `sql.db`, record `repoSlug`, size, mtime, and row count.
3. **Tier 3 Discovery:** Resolve `<repo-cache>/.heatmap/`. Enumerate all sequence-prefixed folders (e.g. `01-antigravity/`). Inspect `root-heatmap.db` and any nested project directories containing `heatmap.db`.
4. **Safety & Zero-Contention:** Open discovered files using read-only DSN:
   `file:<path>?mode=ro&_busy_timeout=3000` to prevent blocking active writers.

---

## 4. Repository Alias Enrichment Engine

GitMap maintains repository aliases in `gitmap.db`'s `Alias` table. To maximize terminal readability for both developers and AI models, discovered repositories are enriched with registered aliases.

### 4.1 Lookup Mechanism
- Call `db.ListAliasesWithRepo()` on the primary store connection.
- Construct a fast bidirectional lookup map:
  ```go
  type RepoAliasMap struct {
  	slugToAliases  map[string][]string
  	slugToPrimary  map[string]string
  }
  ```
- Normalization handles hyphen and underscore variations using `store.SanitizeSlug(slug)`.
- If a discovered repository matches registered aliases (e.g., `gitmap` has aliases `gm`, `gmap`), the primary alias is displayed as a badge in terminal tables: `gitmap [alias: gm]`.

---

## 5. 4-Part RCA: Flag Injection Routing Defect

### 5.1 Symptoms & Observed Failure
When invoking `gitmap rc heatmap` with no arguments, or when passing flags before subcommands (such as `gitmap rc heatmap --repo gitmap status`), the command consistently fails with validation errors or prints generic help:
```text
unknown rc subcommand: --repo
```
The fallback execution path in `executeFallbackRCHeatmap` then receives the invalid argument list, fails its switch block, and displays the command help banner instead of executing the requested action.

### 5.2 Deep Technical Root Cause
The defect exists in `cli/cmd/special_repos_rc_heatmap.go:RunRCHeatmapCommand`:

```go
// Defective implementation
func RunRCHeatmapCommand(args []string) error {
	enrichedArgs := enrichRCHeatmapFlags(args)
	err := heatmap.DispatchHeatmapSubsystem(context.Background(), append([]string{"rc"}, enrichedArgs...), os.Stdout)
	if err == nil {
		return nil
	}
	return executeFallbackRCHeatmap(enrichedArgs)
}
```

1. **Flag Appending Without Subcommand Awareness:** `enrichRCHeatmapFlags(args)` scans `args` for `--repo` and `--ai`. If missing, it blindly appends `--repo <detected>` and `--ai <detected>`.
2. **Subcommand Position Displacement:** When `args` is empty (`gitmap rc heatmap`), `enrichedArgs` becomes `["--repo", "gitmap", "--ai", "antigravity"]`.
3. **Dispatcher Collision:** `DispatchHeatmapSubsystem` strips `"rc"` and passes `enrichedArgs` to `DispatchRcHeatmap`.
4. **Invalid Subcommand Parsing:** `DispatchRcHeatmap` executes `subcmd := strings.ToLower(args[0])`. Because `args[0]` is `"--repo"`, the switch block routes to `default:` and returns error `E1000: unknown rc subcommand: --repo`.
5. **Fallback Loop Failure:** When fallback `executeFallbackRCHeatmap(enrichedArgs)` executes, `args[0]` remains `"--repo"`, failing the second switch statement.

### 5.3 Architecture Fix & Decoupled Dispatch
The dispatch architecture must separate subcommand tokens from flags prior to flag enrichment:

```go
// Extracted clean tokens
type RcCommandTokens struct {
	Subcommand string
	Flags      []string
	HasSubcmd  bool
}

func parseRcCommandTokens(args []string) RcCommandTokens {
	var subcmd string
	var flags []string
	var hasSubcmd bool

	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			flags = append(flags, a)
			continue
		}
		if !hasSubcmd {
			subcmd = strings.ToLower(strings.TrimSpace(a))
			hasSubcmd = true
			continue
		}
		flags = append(flags, a)
	}

	return RcCommandTokens{Subcommand: subcmd, Flags: flags, HasSubcmd: hasSubcmd}
}
```

1. **Explicit Overview Routing:** If `tokens.HasSubcmd` is `false`, route immediately to `heatmap.RunHeatmapOverview(out, flags)`.
2. **Normalized Dispatch Slice:** Reconstruct the dispatched slice with the subcommand strictly at index 0: `append([]string{"rc", tokens.Subcommand}, enrichedFlags...)`.
3. **No Flag Prepending:** Flags are guaranteed to follow the subcommand, preventing argument displacement.

---

## 6. Multi-Repository Recent Tasks Aggregator

To provide immediate situational awareness across polyglot workspaces, GitMap inspects the most active repositories and aggregates recent task activity.

### 6.1 Aggregation Models

```go
package heatmap

import (
	"time"
)

// RecentTaskSummary captures cognitive performance metrics for a recent task.
type RecentTaskSummary struct {
	TaskSlug            string          `json:"taskSlug"`
	RepoSlug            string          `json:"repoSlug"`
	RepoAlias           string          `json:"repoAlias,omitempty"`
	TierType            HeatmapTierType `json:"tierType"`
	TotalActions        int64           `json:"totalActions"`
	ExpectationsMet     int64           `json:"expectationsMet"`
	ExpectationsMetPct  float64         `json:"expectationsMetPct"`
	ReasoningStates     map[string]int  `json:"reasoningStates"`
	LastReasoningState  string          `json:"lastReasoningState"`
	LastActionType      string          `json:"lastActionType"`
	LastActiveAt        time.Time       `json:"lastActiveAt"`
	HasFailedExpectations bool          `json:"hasFailedExpectations"`
}

// MultiRepoRecentSummary holds cross-project task summaries.
type MultiRepoRecentSummary struct {
	GeneratedAt time.Time           `json:"generatedAt"`
	Tasks       []RecentTaskSummary `json:"tasks"`
}
```

### 6.2 Aggregation Logic
1. Scan up to 5 most recently updated heatmap database files.
2. Execute bounded aggregate queries:
   ```sql
   SELECT TaskSlug,
          COUNT(*) as TotalActions,
          COALESCE(SUM(CASE WHEN IsExpectationMet = 1 THEN 1 ELSE 0 END), 0) as MetCount,
          MAX(CreatedAt) as LastActive,
          COALESCE(MAX(ReasoningState), 'UNKNOWN') as LastState
   FROM HeatmapActionLog
   GROUP BY TaskSlug
   ORDER BY LastActive DESC
   LIMIT 10;
   ```
3. Calculate expectation adherence percentage:
   $$\text{ExpectationsMetPct} = \frac{\text{ExpectationsMet}}{\text{TotalActions}} \times 100.0$$
4. Compute boolean flag `hasFailedExpectations`:
   $$\text{hasFailedExpectations} = (\text{TotalActions} - \text{ExpectationsMet}) > 0$$

---

## 7. Dynamic Terminal Suggestions Engine

The suggestions engine in `cli/heatmap/heatmap_suggestions.go` translates discovered database states and recent task health into actionable terminal commands.

### 7.1 Suggestion Matrix

| Task / DB State | Detected Condition | Recommended Terminal Suggestion | Intent & Action |
|---|---|---|---|
| **Failed Expectations** | `hasFailedExpectations == true` | `gitmap heatmap tree --task <taskSlug>` | Review hierarchical causal tree to identify failing steps |
| **Error Search** | `ReasoningStates["FAIL"] > 0` | `gitmap heatmap search --task <taskSlug> --state FAIL` | Filter exact failure logs and error messages |
| **Uncommitted Deltas** | `TierType == TierLocalRoot` and dirty | `gitmap rc heatmap sync --task <taskSlug>` | Synchronize local traces into shared repo-cache |
| **Healthy Task** | `ExpectationsMetPct == 100.0` | `gitmap heatmap status --task <taskSlug>` | Confirm overall task health and file touch counts |
| **Empty Databases** | `TotalDatabases == 0` | `gitmap heatmap log --task <task> --action READ --file <path>` | Start recording cognitive decisions |

### 7.2 Terminal Formatting Example
```text
GITMAP HEATMAP OVERVIEW & TELEMETRY DISCOVERY
Scanned: 3 databases | Total Actions: 142 | Repositories: 2

Discovered Databases:
  REPO       ALIAS   TIER              ACTIONS   SIZE       LAST UPDATED
  gitmap     gm      local_split_repo  128       48.0 KB    2026-10-23 20:15:10
  gitmap     gm      local_root        14        16.0 KB    2026-10-23 20:18:22
  repo-sec   rs      repo_cache_repo   0         12.0 KB    2026-10-23 18:04:12

Recent Multi-Repository Tasks:
  TASK SLUG                             REPO     ACTIONS   MET %    LAST STATE   UPDATED
  313-gitmap-execution-optimization     gitmap   94        97.8%    PASS         10m ago
  314-heatmap-and-url-handling          gitmap   34        100.0%   PASS          2m ago

Dynamic Suggestions for AI & Developers:
  ▶ Inspect failing action branches:
    gitmap heatmap tree --task 313-gitmap-execution-optimization

  ▶ Synchronize local telemetry to shared repo-cache:
    gitmap rc heatmap sync --repo gitmap --task 314-heatmap-and-url-handling

  ▶ Search critical decision logs:
    gitmap heatmap search "routing fix" --json
```

---

## 8. Implementation Contracts & Function Signatures

All implementation functions strictly adhere to GitMap coding guidelines:
- Function line count $\le 15$ lines (preferred $\le 8$ lines).
- Positive booleans (`is*`, `has*` only).
- Zero nested `if` statements (guard clauses only).
- Explicit error propagation using `*apperror.AppError`.

### 8.1 Discovery Signatures (`cli/heatmap/heatmap_discovery.go`)
```go
func DiscoverAllHeatmapDatabases() ([]HeatmapDbEntry, error)
func scanTier1RootDb(dataDir string) (*HeatmapDbEntry, error)
func scanTier2SplitDbs(dataDir string) ([]HeatmapDbEntry, error)
func scanTier3RepoCacheDbs(repoCacheRoot string) ([]HeatmapDbEntry, error)
func inspectDatabaseStats(dbPath string, tier HeatmapTierType, slug string) (HeatmapDbEntry, error)
func formatFileSize(bytes int64) string
```

### 8.2 Overview Signatures (`cli/heatmap/heatmap_overview.go`)
```go
func RunHeatmapOverview(out io.Writer, flags map[string]string) error
func enrichEntriesWithAliases(entries []HeatmapDbEntry) []HeatmapDbEntry
func renderOverviewText(out io.Writer, report HeatmapDiscoveryReport, recent MultiRepoRecentSummary)
func renderOverviewJSON(out io.Writer, report HeatmapDiscoveryReport, recent MultiRepoRecentSummary) error
```

### 8.3 Recent Tasks Signatures (`cli/heatmap/heatmap_recent_tasks.go`)
```go
func QueryMultiRepoRecentTasks(entries []HeatmapDbEntry, limit int) (MultiRepoRecentSummary, error)
func queryTasksFromDatabase(entry HeatmapDbEntry, limit int) ([]RecentTaskSummary, error)
func computeTaskAdherence(met, total int64) (float64, bool)
```

### 8.4 Dynamic Suggestions Signatures (`cli/heatmap/heatmap_suggestions.go`)
```go
func GenerateDynamicSuggestions(recent MultiRepoRecentSummary, entries []HeatmapDbEntry) []string
func buildFailureSuggestions(task RecentTaskSummary) []string
func buildSyncSuggestions(task RecentTaskSummary) []string
func printSuggestionsBlock(out io.Writer, suggestions []string)
```

---

## 9. Verification & Safety Invariants

1. **Read-Only SQLite Access:** All discovery scans and recent task queries MUST open databases with `mode=ro` and `_busy_timeout=3000`. Scans must NEVER acquire write locks or mutate telemetry data.
2. **Relative Path Hygiene:** Zero hardcoded drive letters (`C:`, `D:`) or URI schemes (`file:///`). All paths formatted with forward slashes via `filepath.ToSlash`.
3. **Graceful Degradation:** If `gitmap.db` is locked or unreachable, alias enrichment skips silently without blocking overview rendering.
4. **Zero Build & Zero Test Invariant:** Subagent execution must NOT invoke `go build` or `go test`.
