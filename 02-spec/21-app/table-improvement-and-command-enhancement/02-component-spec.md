# Component Specification: Dual SQLite Databases (Cache & Backup), New Commands Discovery Engine, Table Improvements, Help Documentation & LLM Skills Synchronization

**Module:** `cli/cmdpending`, `cli/constants`, `cli/helpdoc`, `.agents/skills/gitmap`  
**Parent Plan:** `.ai-memory/plans/table-improvement-and-command-enhancement.md`  
**Specification File:** `02-spec/21-app/table-improvement-and-command-enhancement/02-component-spec.md`  
**Version:** 1.2.0  
**Status:** Approved for Implementation  

---

## 1. Executive Summary & Component Scope

This component specification establishes the technical requirements, database schemas, API interfaces, CLI command catalogs, table rendering contracts, backup deserialization protocols, help documentation, and LLM agent skill integrations for five interconnected subsystems of the `table-improvement-and-command-enhancement` task:

1. **Table Rendering Contracts & Layout Engine (`cli/cmdpending/pending_commits_cmd.go`)**:
   - Replaces four fragmented file status columns (`DIRTY`, `UNTRACK`, `MODIF`, `STAGE`) with a single consolidated `UNCOMMITTED` column (`TotalUncommitted = Untracked + Modified + Staged`).
   - Introduces the `VER/BRANCH` column displaying tagged version milestones and active branch (`<version>/<branch>` or `<branch>`).
   - Renders hierarchical two-line tree view remediation hints under dirty repositories without border clipping or line wrapping (`│   ├── Option 1: ...`, `│   └── Option 2: ...`).
   - Displays a prominent fleet batch remediation command (`gitmap cpar "wip: save changes"`) in the table footer when dirty repositories exist.
   - Enforces strict 80-character terminal box alignment (78 inner characters bounded by `│ ... │`).

2. **SQLite Status Cache Engine (`cli/cmdpending/pending_commits_cache.go`)**:
   - Ultra-fast, local SQLite-backed status caching for `gitmap pending-commits` (`gitmap pc`).
   - Replaces expensive multi-repo filesystem git subprocess executions (`git status --porcelain`, `git rev-parse`, `git rev-list`) with sub-millisecond cached telemetry lookups.
   - Strictly enforces a **90-second Time-To-Live (TTL)** (within the required 1–2 minute window).
   - Stale records exceeding 90 seconds are immediately purged from SQLite and never trusted.
   - Dual-gate HEAD SHA verification guarantees cache invalidation whenever new commits occur.
   - First-class bypass flags: `--no-cache` (skip read and write) and `--refresh` (force immediate re-scan and cache update).

3. **SQLite Status Backup Database Engine & Deserialization (`cli/cmdpending/pending_commits_backup.go`)**:
   - Dedicated persistent historical backup database (`pending_commits_backup.db`).
   - Persists point-in-time snapshot records of repository statuses with full metadata (branch, commit SHA, version, uncommitted files count, unpushed count, dirty state, and complete JSON telemetry payload).
   - Unlike the cache DB, backup records are **not auto-purged on the 90-second TTL**, guaranteeing historical fleet states can be retrieved and served again.
   - First-class retrieval flag: `--backup` (`gitmap pc --backup`) allows viewing the latest backed-up snapshot when working offline or reviewing historical states.
   - **Backup Record Deserialization Contract**: Unpacks `payload_json` when serving backups to restore detailed `PendingFiles` and `UnpushedCommitSHAs` lists without data loss.

4. **New Commands Discovery Engine (`cli/cmdpending/new_commands_cmd.go`, `new_commands_types.go`)**:
   - Dedicated command `gitmap new-commands` (alias: `gitmap nc`) cataloging 100 new commands introduced across recent GitMap milestones.
   - 10 functional categories (10 commands each): `commits`, `diagnostics`, `scanner`, `fleet`, `ai`, `os`, `spec`, `storage`, `sync`, `tooling`.
   - Comprehensive filtering flags:
     - `--limit <n>` (alias `-n`, default: 100, clamped min 1)
     - `--filter <str>` (aliases `-f` and `-q`) supporting keyword search across name, alias, description, and copy-pasteable example.
     - `--category <cat>` (alias `-c`) supporting category isolation.
     - `--json` (alias `-j`) emitting structured machine-readable payload.
     - `-h` / `--help` displaying comprehensive contextual usage help.
   - Rich terminal table formatting with category badges and copy-pasteable runnable examples for every cataloged command.

5. **CLI Constants, Help Catalog & LLM Skills Synchronization**:
   - Centralized constants in `cli/constants/constants_cli.go`.
   - Comprehensive markdown documentation in `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md`.
   - Topic summary and alias routing in `cli/helpdoc/catalog.go` and `cli/helpdoc/print.go`.
   - LLM skill integration in `.agents/skills/gitmap/SKILL.md` enabling AI agents to discover dirty repos, parse consolidated `UNCOMMITTED` counts, apply hierarchical tree remedies, utilize cache flags safely, and explore new commands.

---

## 2. Component Interfaces & Go Contracts

### 2.1 Interface 1: `cli/cmdpending/pending_commits_types.go`

Defines core telemetry structures, aggregate counters, and affirmative options:

```go
package cmdpending

import "time"

// UncommittedChangesSummary aggregates metrics for uncommitted working tree changes.
type UncommittedChangesSummary struct {
	UntrackedFilesCount int `json:"untrackedFilesCount"`
	ModifiedFilesCount  int `json:"modifiedFilesCount"`
	StagedFilesCount    int `json:"stagedFilesCount"`
	TotalUncommitted    int `json:"totalUncommitted"`
}

// UnpushedCommitSummary aggregates metrics for local commits ahead of upstream branch.
type UnpushedCommitSummary struct {
	UnpushedCommitsCount int  `json:"unpushedCommitsCount"`
	HasUpstream          bool `json:"hasUpstream"`
}

// RemediationOption represents an actionable fix for a dirty repository.
type RemediationOption struct {
	OptionNumber int    `json:"optionNumber"`
	Label        string `json:"label"`
	Command      string `json:"command"`
}

// RepoPendingCommitRecord represents the pending changes status of an individual repository.
type RepoPendingCommitRecord struct {
	RepoName             string              `json:"repoName"`
	RelativePath         string              `json:"relativePath"`
	CurrentBranch        string              `json:"currentBranch"`
	Version              string              `json:"version,omitempty"`
	ShortVersionBranch   string              `json:"shortVersionBranch"`
	IsDirty              bool                `json:"isDirty"`
	IsClean              bool                `json:"isClean"`
	HasUncommitted       bool                `json:"hasUncommitted"`
	HasUnpushed          bool                `json:"hasUnpushed"`
	HasUpstream          bool                `json:"hasUpstream"`
	TotalUncommitted     int                 `json:"totalUncommitted"`
	UntrackedFilesCount  int                 `json:"untrackedFilesCount"`
	ModifiedFilesCount   int                 `json:"modifiedFilesCount"`
	StagedFilesCount     int                 `json:"stagedFilesCount"`
	UnpushedCommitsCount int                 `json:"unpushedCommitsCount"`
	PendingFiles         []string            `json:"pendingFiles,omitempty"`
	UnpushedCommitSHAs   []string            `json:"unpushedCommitShas,omitempty"`
	RemediationOptions   []RemediationOption `json:"remediationOptions,omitempty"`
}

// PendingCommitsPayload represents the top-level JSON telemetry for pending-commits.
type PendingCommitsPayload struct {
	Timestamp             time.Time                 `json:"timestamp"`
	TotalReposScanned     int                       `json:"totalReposScanned"`
	TotalDirtyRepos       int                       `json:"totalDirtyRepos"`
	TotalUncommittedFiles int                       `json:"totalUncommittedFiles"`
	TotalUnpushedCommits  int                       `json:"totalUnpushedCommits"`
	SortMode              string                    `json:"sortMode"`
	DetailMode            string                    `json:"detailMode"`
	Repositories          []RepoPendingCommitRecord `json:"repositories"`
}

// PendingCommitsOptions encapsulates CLI execution parameters for pending-commits.
type PendingCommitsOptions struct {
	SortMode    string
	DetailMode  string
	IsSSH       bool
	IsJSON      bool
	IsDirtyOnly bool
	IsAll       bool
	TargetRepo  string
	IsNoCache   bool
	IsRefresh   bool
	IsBackup    bool
}
```

### 2.2 Interface 2: `cli/cmdpending/pending_commits_cmd.go`

Implements CLI routing, options parsing, git inspection coordination, and terminal rendering:

```go
package cmdpending

import "database/sql"

// RunPendingCommits is the CLI dispatch entry point for gitmap pending-commits and gitmap pc.
func RunPendingCommits(args []string) error

// parsePendingCommitsOptions extracts execution flags, sorting, detail mode, and target repo.
func parsePendingCommitsOptions(args []string) PendingCommitsOptions

// runLocalPendingCommits orchestrates cache/backup DB connections and repository evaluation.
func runLocalPendingCommits(opts PendingCommitsOptions) error

// serveBackupPendingCommits queries the backup DB and renders historical status snapshots.
func serveBackupPendingCommits(backupConn *sql.DB, opts PendingCommitsOptions) error

// formatShortVersionBranch generates the combined VER/BRANCH string with truncation.
func formatShortVersionBranch(version, branch string, maxLen int) string

// buildRepoRemediationOptions computes Option 1 (commit & push) and Option 2 (stash) commands.
func buildRepoRemediationOptions(rec RepoPendingCommitRecord) []RemediationOption

// renderLocalPendingCommitsTerminal renders the rectangular 80-character box table.
func renderLocalPendingCommitsTerminal(payload PendingCommitsPayload, totalKnown int, inspected []RepoPendingCommitRecord, opts PendingCommitsOptions)

// renderRepoRemediationTree renders indented tree branches for a dirty repo row without clipping.
func renderRepoRemediationTree(rec RepoPendingCommitRecord)

// renderFleetBatchFooter renders the workspace-wide batch remediation box in the table footer.
func renderFleetBatchFooter(hasDirty bool)
```

### 2.3 Interface 3: `cli/cmdpending/pending_commits_cache.go`

Implements ephemeral status caching with 90-second TTL and dual-gate validation:

```go
package cmdpending

import "database/sql"

// PendingCommitCacheRecord stores cached git status metrics for a repository.
type PendingCommitCacheRecord struct {
	RepoPath         string
	HeadSHA          string
	UncommittedCount int
	UnpushedCount    int
	IsDirty          bool
	CachedAtUnix     int64
}

// DefaultCacheTTLSeconds defines the standard 90-second TTL (1.5 minutes).
const DefaultCacheTTLSeconds int64 = 90

// OpenPendingCommitsCache opens or initializes the cache SQLite database in WAL mode.
func OpenPendingCommitsCache(dbPath string) (*sql.DB, error)

// GetCachedPendingStatus retrieves a valid, non-expired cache record for a repo.
// Returns (nil, false, nil) on cache miss or expired record (expired records are purged).
func GetCachedPendingStatus(conn *sql.DB, repoPath, currentHeadSHA string, nowUnix int64) (*PendingCommitCacheRecord, bool, error)

// SaveCachedPendingStatus writes or updates a repository status record in SQLite.
func SaveCachedPendingStatus(conn *sql.DB, record PendingCommitCacheRecord) error

// PurgeExpiredPendingCache removes all records older than ttlSeconds.
func PurgeExpiredPendingCache(conn *sql.DB, nowUnix, ttlSeconds int64) (int64, error)

// InvalidatePendingCache deletes an individual repository record by path.
func InvalidatePendingCache(conn *sql.DB, repoPath string) error
```

### 2.4 Interface 4: `cli/cmdpending/pending_commits_backup.go`

Implements persistent snapshot storage, querying, and deserialization:

```go
package cmdpending

import "database/sql"

// PendingCommitBackupRecord represents a durable backup snapshot of a repository status.
type PendingCommitBackupRecord struct {
	RepoPath           string `json:"repoPath"`
	RepoName           string `json:"repoName"`
	CurrentBranch      string `json:"currentBranch"`
	Version            string `json:"version"`
	ShortVersionBranch string `json:"shortVersionBranch"`
	HeadSHA            string `json:"headSha"`
	UncommittedCount   int    `json:"uncommittedCount"`
	UnpushedCount      int    `json:"unpushedCount"`
	IsDirty            bool   `json:"isDirty"`
	BackedUpAtUnix     int64  `json:"backedUpAtUnix"`
	PayloadJSON        string `json:"payloadJson,omitempty"`
}

// OpenPendingCommitsBackup opens or initializes the backup SQLite database in WAL mode.
func OpenPendingCommitsBackup(dbPath string) (*sql.DB, error)

// SaveBackupPendingStatus writes or updates a durable repository backup status record.
func SaveBackupPendingStatus(conn *sql.DB, record RepoPendingCommitRecord, headSHA string, nowUnix int64) error

// GetBackupPendingStatus retrieves a backup status record by repository path.
func GetBackupPendingStatus(conn *sql.DB, repoPath string) (*PendingCommitBackupRecord, bool, error)

// ListAllBackupPendingStatuses returns all stored backup status records sorted by time descending.
func ListAllBackupPendingStatuses(conn *sql.DB) ([]PendingCommitBackupRecord, error)

// RestoreBackupToCache copies a backed up record into the active cache database.
func RestoreBackupToCache(backupConn, cacheConn *sql.DB, repoPath string, nowUnix int64) error

// PurgeBackupOlderThan removes backup records older than cutoffUnix.
func PurgeBackupOlderThan(conn *sql.DB, cutoffUnix int64) (int64, error)
```

### 2.5 Interface 5: `cli/cmdpending/new_commands_cmd.go` & `new_commands_types.go`

Implements discovery, filtering, and terminal/JSON rendering for the last 100 commands:

```go
package cmdpending

import (
	"time"
	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// NewCommandEntry represents a cataloged command entry.
type NewCommandEntry struct {
	Name        string   `json:"name"`
	Alias       string   `json:"alias,omitempty"`
	Version     string   `json:"version"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Example     string   `json:"example"`
	Flags       []string `json:"flags,omitempty"`
}

// NewCommandsPayload encapsulates structured output for new commands.
type NewCommandsPayload struct {
	Timestamp        time.Time         `json:"timestamp"`
	TotalCommands    int               `json:"totalCommands"`
	FilteredCommands int               `json:"filteredCommands"`
	Limit            int               `json:"limit"`
	Category         string            `json:"category,omitempty"`
	Filter           string            `json:"filter,omitempty"`
	Commands         []NewCommandEntry `json:"commands"`
}

// NewCommandsOptions encapsulates parsed CLI flags.
type NewCommandsOptions struct {
	Limit    int
	Category string
	Filter   string
	IsJSON   bool
	IsHelp   bool
}

// DefaultNewCommandsOptions returns options initialized with safe defaults.
func DefaultNewCommandsOptions() NewCommandsOptions

// RunNewCommands executes the new-commands discovery tool.
func RunNewCommands(args []string) error

// executeNewCommands coordinates catalog building, filtering, and rendering.
func executeNewCommands(opts NewCommandsOptions) error

// parseNewCommandsOptions extracts --limit, --category, --filter (-f, -q), and --json.
func parseNewCommandsOptions(args []string) (NewCommandsOptions, *appfault.AppError)

// buildNewCommandsCatalog returns the exhaustive 100-command catalog across 10 categories.
func buildNewCommandsCatalog() []NewCommandEntry

// filterNewCommands applies limit, category, and text filter to the catalog.
func filterNewCommands(catalog []NewCommandEntry, opts NewCommandsOptions) []NewCommandEntry

// renderNewCommandsTerminal renders the structured terminal view.
func renderNewCommandsTerminal(payload NewCommandsPayload) error

// renderNewCommandsJSON emits formatted JSON.
func renderNewCommandsJSON(payload NewCommandsPayload) error
```

---

## 3. Table Rendering Contracts, Tree View Formatting & Footer Layout

### 3.1 Terminal Box Dimensions & Width Arithmetic
The pending commits table enforces an exact 80-character terminal box layout:
- Left and right border characters: `│` (1 character each).
- Inner content width: 78 characters.
- Total line width: $1 + 78 + 1 = 80$ characters.

```text
┌──────────────────────────────────────────────────────────────────────────────┐ (80 chars)
│                               78 inner chars                                 │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Column Width Specifications

| Column Header | Format String | Max Width | Content Description |
| :--- | :--- | :--- | :--- |
| `REPOSITORY` | `%-22s` | 22 chars | Repository name or folder slug. Truncated with `…` if longer. |
| `VER/BRANCH` | `%-16s` | 16 chars | Tagged release and branch formatted via `formatShortVersionBranch`. |
| `UNCOMMITTED` | `%12s` | 12 chars | Consolidated uncommitted count (`TotalUncommitted`). |
| `UNPUSHED` | `%10s` | 10 chars | Count of commits ahead of remote tracking branch. |
| `STATUS` | `%-9s` | 9 chars | Status badge (`● PEND` or `○ CLEAN`). |

**Column Spacing Accounting:**
- `REPOSITORY` (22) + 1 space + `VER/BRANCH` (16) + 1 space + `UNCOMMITTED` (12) + 1 space + `UNPUSHED` (10) + 1 space + `STATUS` (9) + trailing padding (5) = 78 inner characters.

### 3.3 `VER/BRANCH` Formatting Contract
The `formatShortVersionBranch(version, branch string, maxLen int)` helper evaluates:
1. If `version` is empty: returns `branch`.
2. If `version` is present: returns `<version>/<branch>` (e.g. `v6.523.1/main`).
3. If combined length exceeds `maxLen` (16 characters):
   - Keeps `version` intact if possible.
   - Truncates `branch` with `…` (e.g. `v6.523/feat-a…`).
   - If `version` alone exceeds `maxLen`, truncates `version` with `…`.

### 3.4 Hierarchical Tree View Formatting Without Border Clipping

For every dirty repository (`IsDirty == true`), the rendering engine outputs two indented tree lines immediately following the repository row.

#### Exact Formatting Rules:
1. Each line begins with `│   ├── Option 1: ` or `│   └── Option 2: `.
2. Tree text format:
   - Line 1: `│   ├── Option 1: git -C "<repo>" add -A && git commit -m "wip: save changes" && git push`
   - Line 2: `│   └── Option 2: git -C "<repo>" stash -u`
3. **No-Clipping Invariant:**
   - Inner prefix length: `│   ├── Option 1: ` is 17 characters.
   - Inner command budget: $78 - 17 = 61$ characters.
   - If command length exceeds 61 characters, the command string is truncated to 60 characters with `…`.
   - The line is right-padded with spaces up to index 78.
   - The closing border `│` is appended at index 79.
   - Result: Every tree line is exactly 80 characters wide, preserving pristine rectangular box borders across all platforms without wrapping or clipping.
4. Clean repositories (`IsDirty == false`) **NEVER** render tree view lines.

### 3.5 Table Footer Layout Specification
The summary footer renders at the bottom of the table:

1. **Header Summary Box:**
   ```text
   ├──────────────────────────────────────────────────────────────────────────────┤
   │ Scanned: 42 repos    │ Dirty: 2 repos   │ Uncommitted: 15 files │ Unpushed: 1 │
   ├──────────────────────────────────────────────────────────────────────────────┤
   ```
2. **Clean Repositories Aggregation Row:**
   - If clean repositories exist and `--all` is not specified, clean repositories are collapsed into a single summary row:
   ```text
   │ 40 clean repos         -                            0          0   ○ CLEAN   │
   ```
3. **Fleet Batch Remediation Box (Conditional):**
   - If `TotalDirtyRepos > 0`, the footer includes the batch remediation command:
   ```text
   ├──────────────────────────────────────────────────────────────────────────────┤
   │ Fleet Remediation: gitmap cpar "wip: save changes"                           │
   └──────────────────────────────────────────────────────────────────────────────┘
   ```
   - If `TotalDirtyRepos == 0`, the fleet remediation line and separator are omitted, cleanly closing the table with `└─────────────...─┘`.

---

## 4. Backup Record Deserialization Contract

### 4.1 Problem & Motivation
When saving status snapshots via `SaveBackupPendingStatus`, the engine serializes the complete `RepoPendingCommitRecord` into `payload_json`. This payload contains rich internal state:
- `PendingFiles`: List of changed filenames.
- `UnpushedCommitSHAs`: List of commit hashes.
- `UntrackedFilesCount`, `ModifiedFilesCount`, `StagedFilesCount`.
- `RemediationOptions`.

Previously, when serving records via `--backup` (`gitmap pc --backup`), `convertBackupToPendingRecords` only copied top-level column values (`UncommittedCount`, `UnpushedCount`), discarding `payload_json`. As a result, detailed inspections (`--detail all`, `--detail 1`) or JSON outputs lost the actual file and commit lists.

### 4.2 Deserialization Contract Specification
When converting `PendingCommitBackupRecord` to `RepoPendingCommitRecord`:

```go
func convertBackupToPendingRecord(b PendingCommitBackupRecord) RepoPendingCommitRecord {
	// Step 1: Base initialization from database columns
	item := RepoPendingCommitRecord{
		RepoName:             b.RepoName,
		RelativePath:         b.RepoPath,
		CurrentBranch:        b.CurrentBranch,
		Version:              b.Version,
		ShortVersionBranch:   b.ShortVersionBranch,
		IsDirty:              b.IsDirty,
		IsClean:              !b.IsDirty && b.UnpushedCount == 0,
		HasUncommitted:       b.UncommittedCount > 0,
		HasUnpushed:          b.UnpushedCount > 0,
		TotalUncommitted:     b.UncommittedCount,
		UnpushedCommitsCount: b.UnpushedCount,
	}

	// Step 2: Deserialization Contract - Unpack payload_json if present
	if strings.TrimSpace(b.PayloadJSON) != "" && b.PayloadJSON != "{}" {
		var unpacked RepoPendingCommitRecord
		if err := json.Unmarshal([]byte(b.PayloadJSON), &unpacked); err == nil {
			item.PendingFiles = unpacked.PendingFiles
			item.UnpushedCommitSHAs = unpacked.UnpushedCommitSHAs
			item.UntrackedFilesCount = unpacked.UntrackedFilesCount
			item.ModifiedFilesCount = unpacked.ModifiedFilesCount
			item.StagedFilesCount = unpacked.StagedFilesCount
			item.HasUpstream = unpacked.HasUpstream
			if len(unpacked.RemediationOptions) > 0 {
				item.RemediationOptions = unpacked.RemediationOptions
			}
		}
	}

	// Step 3: Fallback remediation options generation if not restored from JSON
	if item.IsDirty && len(item.RemediationOptions) == 0 {
		item.RemediationOptions = buildRepoRemediationOptions(item)
	}

	return item
}
```

### 4.3 Invariants:
1. **Zero Data Loss:** Detailed file lists (`PendingFiles`) and commit SHAs (`UnpushedCommitSHAs`) are fully restored from backup snapshots.
2. **Graceful Fallback:** If `payload_json` is corrupted or empty, column counts are safely preserved and remediation options are synthesized on the fly.

---

## 5. New Commands Discovery Engine: Complete 100-Command Catalog

The discovery engine (`gitmap new-commands` / `gitmap nc`) indexes 100 commands introduced across recent GitMap history, organized across 10 functional categories (10 commands per category).

### 5.1 Category 1: `commits` (Commit & Workspace Operations)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 1 | `cpar` | `commit-push-all-repos` | `v6.520.0` | Commit and push all workspace repositories | `gitmap cpar "wip: save changes"` |
| 2 | `cpf` | `commit-push-feature` | `v6.500.0` | Atomic commit and push feature branch | `gitmap cpf "auth - add jwt validation"` |
| 3 | `cpb` | `commit-push-bug` | `v6.500.0` | Atomic commit and push bugfix branch | `gitmap cpb "cache - fix ttl expiration"` |
| 4 | `cpr` | `commit-push-release` | `v6.500.0` | Commit, push, and trigger release chore | `gitmap cpr "v6.524.0"` |
| 5 | `cin` | `commit-in` | `v6.510.0` | Commit in replay with JSON author rotation and AST heuristics | `gitmap cin --source repo-a --dest repo-b` |
| 6 | `pending-commits` | `pc` | `v6.523.0` | Discover uncommitted changes and unpushed commits with tree remediation | `gitmap pc` |
| 7 | `backup-branch` | `bb` | `v6.518.0` | Create automated backup branch for workspace repositories | `gitmap backup-branch "feat-login"` |
| 8 | `fix all` | `fix-all` | `v6.522.0` | Batch stash, commit, or discard pending changes across fleet | `gitmap fix all --action=wip` |
| 9 | `pcp` | `pull-commit-push` | `v6.505.0` | Pull rebase, stage all, commit with message, and push | `gitmap pcp "sync upstream"` |
| 10 | `commons` | `co` | `v6.515.0` | Apply standardized repo configurations and git hooks | `gitmap commons` |

### 5.2 Category 2: `diagnostics` (CI/CD Telemetry & Error Analysis)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 11 | `pe` | `pipeline-errors` | `v6.510.0` | Pipeline error analyzer with dynamic runner extraction | `gitmap pe -t` |
| 12 | `pe all` | `pipeline-errors-all` | `v6.515.0` | Extract failing step logs across all recent workflow runs | `gitmap pe all` |
| 13 | `te all` | `test-errors-all` | `v6.516.0` | Aggregate failing unit test output across all packages | `gitmap te all` |
| 14 | `sug` | `shutdown-until-green` | `v6.512.0` | Pause execution until CI/CD quality gates are green | `gitmap sug` |
| 15 | `rerun` | `rr` | `v6.511.0` | Rerun failed pipeline jobs or Antigravity prompts | `gitmap rerun --step=lint` |
| 16 | `regoldens` | `rg-gold` | `v6.517.0` | Regenerate golden test fixture outputs across test suites | `gitmap regoldens` |
| 17 | `doctor` | `doc` | `v6.508.0` | Comprehensive workstation and toolchain diagnostic health check | `gitmap doctor` |
| 18 | `audit-legacy` | `al` | `v6.514.0` | Scan codebase for deprecated APIs and legacy function signatures | `gitmap audit-legacy` |
| 19 | `fastgate` | `fg` | `v6.519.0` | Run fast preflight lint and formatting verification gates | `gitmap fastgate` |
| 20 | `pipeline-ai` | `pl-ai` | `v6.513.0` | Check CI/CD workflow state and wait dynamically on ETA | `gitmap pipeline-ai status -t 120` |

### 5.3 Category 3: `scanner` (Repository Discovery & Manifests)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 21 | `scan export` | `s-exp` | `v6.521.0` | Export scanned repository metadata to JSON or CSV manifest | `gitmap scan export --format=json` |
| 22 | `scan merge` | `s-mrg` | `v6.521.0` | Merge external scan databases into central repository index | `gitmap scan merge --src=other.db` |
| 23 | `rescan` | `rsc` | `v6.502.0` | Perform delta rescan of workspace filesystem changes | `gitmap rescan` |
| 24 | `rescan-subtree` | `rss` | `v6.504.0` | Narrowly re-run scan against an at-cap repository subtree | `gitmap rescan-subtree --max-depth=5` |
| 25 | `clone-only-missing` | `com` | `v6.506.0` | Clone only un-cloned repositories from workspace manifest | `gitmap com --manifest=repos.json` |
| 26 | `clone-sync` | `cs` | `v6.507.0` | Synchronize and clone all remote fleet repositories | `gitmap clone-sync` |
| 27 | `repo-create` | `repoc` | `v6.503.0` | Create a new local repository with standard folder scaffold | `gitmap repo-create my-service` |
| 28 | `dedupe` | `dd` | `v6.509.0` | Detect duplicate repository clones and directories across disks | `gitmap dedupe` |
| 29 | `size` | `sz` | `v6.501.0` | Inspect repository disk size and identify large bloated objects | `gitmap size --top=10` |
| 30 | `orphans` | `orph` | `v6.510.0` | Discover orphaned git directories not indexed in split-db | `gitmap orphans` |

### 5.4 Category 4: `fleet` (Multi-Node SSH Delegation & Cluster Topologies)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 31 | `ssh-bind` | `sb` | `v6.518.0` | Bind local port forwarding tunnels over SSH to fleet nodes | `gitmap ssh-bind --remote=node-01 --port=8080` |
| 32 | `ssh-deploy` | `sd` | `v6.516.0` | Deploy compiled binaries and configuration to remote SSH host | `gitmap ssh-deploy --node=prod-1` |
| 33 | `nodes` | `cluster-nodes` | `v6.512.0` | List registered fleet nodes with reachability and latency | `gitmap nodes` |
| 34 | `cluster` | `cl` | `v6.511.0` | Manage multi-node cluster topology and remote execution | `gitmap cluster status` |
| 35 | `sc` | `servers-clients` | `v6.514.0` | Distributed fan-out execution across server-client topologies | `gitmap sc exec "uptime"` |
| 36 | `deploy-keys` | `dk` | `v6.515.0` | Deploy authorized SSH public keys across all fleet machines | `gitmap deploy-keys --all` |
| 37 | `ssh-join` | `sj` | `v6.510.0` | Enroll target machine into SSH host registry with key auth | `gitmap ssh-join user@192.168.1.50` |
| 38 | `cluster-run-script` | `crs` | `v6.513.0` | Deploy and execute local script remotely across cluster nodes | `gitmap cluster-run-script deploy.sh` |
| 39 | `cluster-bootstrap` | `cb` | `v6.515.0` | Bootstrap node with RSA keys, passwordless sudo, and tools | `gitmap cluster-bootstrap 192.168.1.51` |
| 40 | `ssh-scan` | `ss` | `v6.512.0` | Scan local subnet for machines with open SSH port 22 | `gitmap ssh-scan 192.168.1.0/24` |

### 5.5 Category 5: `ai` (Autonomous Agents, Curriculum & Orchestration)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 41 | `agent task` | `at` | `v6.522.0` | Manage autonomous agent task lifecycle and action logs in SQLite | `gitmap task claim --db task.db --agent "Worker 01"` |
| 42 | `agy` | `antigravity` | `v6.510.0` | Google Antigravity workspace manager and prompt orchestrator | `gitmap agy deploy` |
| 43 | `agm` | `antigravity-mgr` | `v6.512.0` | Manage Antigravity Manager GUI application and fleet tools | `gitmap agm status` |
| 44 | `aum` | `auto` | `v6.511.0` | High-performance native Go automation and multi-core search | `gitmap aum search "pattern" cli` |
| 45 | `ai-analysis` | `aa` | `v6.514.0` | Perform automated codebase architecture and guideline analysis | `gitmap ai-analysis --rules=02-spec` |
| 46 | `macro` | `mc` | `v6.509.0` | Record, inspect, and replay interactive shell terminal macros | `gitmap macro replay build-all` |
| 47 | `llm-docs` | `ld` | `v6.515.0` | Generate consolidated markdown command matrix for LLMs | `gitmap llm-docs` |
| 48 | `llm train` | `llm-train` | `v6.516.0` | Run 4-stage LLM curriculum and generate developer skill | `gitmap llm train` |
| 49 | `prompt show --copy` | `psc` | `v6.520.0` | Render canonical prompt and copy directly to system clipboard | `gitmap prompt show --copy 01-prompts/v6.md` |
| 50 | `ai-fix` | `af` | `v6.513.0` | Run standardized repository autofix targets for coding guidelines | `gitmap ai-fix --target=naming` |

### 5.6 Category 6: `os` (Operating System & Desktop Configuration)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 51 | `os dock` | `dock` | `v6.517.0` | Configure taskbar and dock alignment across Windows and Linux | `gitmap os dock bottom` |
| 52 | `os panel` | `panel` | `v6.518.0` | Configure system display panels and multi-monitor geometry | `gitmap os panel --primary=1` |
| 53 | `apps` | `app` | `v6.516.0` | Installed application auditor and uninstaller framework | `gitmap apps list` |
| 54 | `fix-link` | `fl` | `v6.512.0` | Inspect and repair broken symlinks and VMware shared mounts | `gitmap fix-link --repair` |
| 55 | `vpm` | `vmware-power` | `v6.514.0` | VMware guest shared folder mounting and open-vm-tools manager | `gitmap vpm mount` |
| 56 | `which-format` | `wf` | `v6.519.0` | Inspect file line ending (CRLF/LF) and encoding formats | `gitmap which-format file.go` |
| 57 | `os-dns` | `dns` | `v6.513.0` | Inspect, benchmark, and switch system DNS servers | `gitmap os-dns switch 1.1.1.1` |
| 58 | `os-theme` | `theme` | `v6.511.0` | Switch desktop appearance between Dark Mode and Light Mode | `gitmap os-theme dark` |
| 59 | `power` | `pwr` | `v6.515.0` | Screen timeout and power plan management with state tracking | `gitmap power sleep 30` |
| 60 | `autologin` | `os-autologin` | `v6.510.0` | Configure OS auto-login credentials for Windows and Ubuntu | `gitmap autologin --user=admin` |

### 5.7 Category 7: `spec` (Specification Authoring, Auditing & Documentation)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 61 | `spec issue` | `si` | `v6.522.0` | Audit specification gap issues and generate remediation steps | `gitmap spec issue --spec=02-spec/21-app` |
| 62 | `spec-author` | `sa` | `v6.520.0` | Scaffold structured architecture and component specifications | `gitmap spec-author --slug=my-feature` |
| 63 | `spec-audit` | `spa` | `v6.519.0` | Conduct blind-AI readiness audits on specification documents | `gitmap spec-audit --target=02-spec` |
| 64 | `user` | `usr` | `v6.521.0` | Audit and configure git user name and email per repository | `gitmap user set --name="Dev" --email="dev@co.com"` |
| 65 | `changelog` | `clog` | `v6.508.0` | Generate SemVer changelog entries from recent atomic commits | `gitmap changelog generate` |
| 66 | `release-notes` | `rn` | `v6.509.0` | Extract markdown release notes for tagged versions | `gitmap release-notes v6.523.0` |
| 67 | `seowrite` | `seo` | `v6.510.0` | Format repository descriptions and keywords for SEO | `gitmap seowrite --keyword=git` |
| 68 | `help-json` | `hj` | `v6.515.0` | Emit machine-readable JSON schema for CLI help topics | `gitmap help --json --filter=task` |
| 69 | `prompt-template` | `pt` | `v6.512.0` | Manage reusable AI prompt prefix and verification templates | `gitmap prompt-template list` |
| 70 | `prompt list` | `plst` | `v6.518.0` | List registered canonical AI prompts in 01-prompts directory | `gitmap prompt list` |

### 5.8 Category 8: `storage` (Secrets Vault, Cache Stores & Disk Hygiene)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 71 | `rs file` | `rs-f` | `v6.515.0` | Copy secret file into encrypted repo-secrets and auto-push | `gitmap rs file .env --repo=backend` |
| 72 | `rs folder` | `rs-d` | `v6.515.0` | Copy secret folder into repo-secrets and auto-push | `gitmap rs folder certs/` |
| 73 | `rs text` | `rs-t` | `v6.515.0` | Store sensitive string directly into repo-secrets store | `gitmap rs text "key_123" --slug=api-key` |
| 74 | `rc file` | `rc-f` | `v6.516.0` | Offload reusable test script into repo-cache and auto-push | `gitmap rc file test_harness.py` |
| 75 | `rc folder` | `rc-d` | `v6.516.0` | Store shared fixture directory into repo-cache | `gitmap rc folder fixtures/` |
| 76 | `rc text` | `rc-t` | `v6.516.0` | Write reusable automation test script into repo-cache | `gitmap rc text "script" --slug=test --ext=.ps1` |
| 77 | `storage` | `stg` | `v6.511.0` | Inspect disk volumes, split-DB sqlite sizes, and snapshots | `gitmap storage` |
| 78 | `clean-dev` | `cld` | `v6.517.0` | Purge IDE temporary files, dev caches, and lock buffers | `gitmap clean-dev` |
| 79 | `clear-terminal` | `cls-term` | `v6.518.0` | Clear stuck terminal processes, handles, and buffer pipes | `gitmap clear-terminal` |
| 80 | `purge-cache` | `pgc` | `v6.520.0` | Invalidate and purge stale SQLite status cache entries | `gitmap purge-cache --ttl=90s` |

### 5.9 Category 9: `sync` (Multi-Repository Synchronization & Fleet Baselines)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 81 | `sync` | `sy` | `v6.522.0` | Synchronize prompts, skills, and specs across 43 repositories | `gitmap sync --workers=8` |
| 82 | `sync --repo` | `sy-r` | `v6.522.0` | Synchronize canonical assets to single target repository | `gitmap sync --repo=movie-cli-v8` |
| 83 | `sync --dry-run` | `sy-d` | `v6.522.0` | Preview synchronization file diffs without disk mutations | `gitmap sync --dry-run` |
| 84 | `sync --list` | `sy-l` | `v6.522.0` | List registered fleet repositories connected to sync engine | `gitmap sync --list` |
| 85 | `pull-all-efficient` | `pae` | `v6.508.0` | Pull all repositories concurrently with JSON telemetry | `gitmap pae --json` |
| 86 | `reconcile` | `recon` | `v6.505.0` | Reconcile divergent branch heads and resolve conflicts | `gitmap reconcile` |
| 87 | `has-any-updates` | `hau` | `v6.509.0` | Check remote tracking branches for incoming commits | `gitmap hau` |
| 88 | `latest-branch` | `lb` | `v6.510.0` | Discover the most recently updated remote branch across fleet | `gitmap lb` |
| 89 | `desktop-sync` | `ds` | `v6.503.0` | Synchronize local repositories with GitHub Desktop state | `gitmap desktop-sync` |
| 90 | `watch` | `w` | `v6.504.0` | Live-refresh terminal dashboard monitoring repo changes | `gitmap watch` |

### 5.10 Category 10: `tooling` (Developer Utilities, Search & Replacement)

| # | Name | Alias | Version | Description | Example |
|---|---|---|---|---|---|
| 91 | `new-commands` | `nc` | `v6.524.0` | Inspect recent commands added across releases with examples | `gitmap nc --limit=20` |
| 92 | `which-format` | `whichfmt` | `v6.519.0` | Detect line endings, encoding BOM, and file permissions | `gitmap which-format ./cli` |
| 93 | `find-files` | `ff` | `v6.512.0` | Find exact filename across workspace within milliseconds | `gitmap ff "types.go"` |
| 94 | `find-files-any` | `ffa` | `v6.512.0` | Find files matching substring across 10,000+ files | `gitmap ffa "cache"` |
| 95 | `find-files-startswith` | `ffs` | `v6.512.0` | Find files matching filename prefix | `gitmap ffs "test_"` |
| 96 | `find-files-endswith` | `ffe` | `v6.512.0` | Find files matching filename extension or suffix | `gitmap ffe "_test.go"` |
| 97 | `list-files` | `lf` | `v6.511.0` | Stream relative file paths matching pattern or folder | `gitmap lf cli/cmdpending` |
| 98 | `replace` | `rep` | `v6.513.0` | Literal multi-file string replacement with audit trail | `gitmap replace "oldString" "newString"` |
| 99 | `replace-regex` | `repr` | `v6.513.0` | Regex pattern replacement across repository code | `gitmap replace-regex "v[0-9]+" "v6.524.0"` |
| 100 | `cat` | `ct` | `v6.514.0` | Zero-disk stdout stream of file content in CLI sessions | `gitmap cat cli/constants/constants_cli.go` |

---

## 6. Documentation & LLM Integration Requirements

### 6.1 `cli/constants/constants_cli.go`
The following constants are registered to support command routing and help text:

```go
package constants

const (
	// CmdPendingCommits discovers uncommitted changes and unpushed commits across repos.
	CmdPendingCommits      = "pending-commits"
	CmdPendingCommitsAlias = "pc"

	// CmdNewCommands catalogs and filters commands added in recent git history.
	CmdNewCommands         = "new-commands"
	CmdNewCommandsAlias    = "nc"

	// HelpPendingCommits describes pending-commits in root help.
	HelpPendingCommits = "  pending-commits (pc) Discover uncommitted changes & unpushed commits with tree remediation"

	// HelpNewCommands describes new-commands in root help.
	HelpNewCommands    = "  new-commands (nc)    Inspect last 100 commands added in recent commits with examples"

	// Flags for cache and backup control
	FlagNoCache     = "no-cache"
	FlagDescNoCache = "Bypass status cache and query git status directly from filesystem"

	FlagRefresh     = "refresh"
	FlagDescRefresh = "Force refresh status cache by re-evaluating git status and updating SQLite"

	FlagBackup      = "backup"
	FlagDescBackup  = "Serve repository statuses from persistent backup database without live git scan"

	FlagLimit       = "limit"
	FlagDescLimit   = "Maximum number of items to display"

	FlagFilter      = "filter"
	FlagDescFilter  = "Keyword filter matching command name, alias, description, or example"

	FlagCategory      = "category"
	FlagDescCategory  = "Category filter (e.g. commits, diagnostics, scanner, fleet, ai, os, spec, storage, sync, tooling)"
)
```

### 6.2 `cli/helpdoc/pending-commits.md`
Must document:
1. Command aliases: `pc` and `pending-commits`.
2. All flags: `--sort`, `--detail`, `--no-cache`, `--refresh`, `--backup`, `--ssh`, `--json`, `--dirty-only`, `--all`.
3. Explanation of dual SQLite architecture (90s status cache vs persistent backup DB).
4. Description of consolidated `UNCOMMITTED` column and `VER/BRANCH` display.
5. Indented tree view remediation options (`Option 1` vs `Option 2`).
6. Global fleet batch command in footer (`gitmap cpar "wip: save changes"`).
7. Concrete runnable examples covering cache refresh, cache bypass, backup serving, and JSON output.

### 6.3 `cli/helpdoc/new-commands.md`
Must document:
1. Command aliases: `nc` and `new-commands`.
2. All flags: `--limit` / `-n`, `--filter` / `-f` / `-q`, `--category` / `-c`, `--json` / `-j`, `-h` / `--help`.
3. Description of the 10 functional categories.
4. Concrete runnable examples for filtering by keyword, filtering by category, limiting count, and emitting JSON.

### 6.4 `cli/helpdoc/catalog.go` & `cli/helpdoc/print.go`
- `topicSummaries` in `cli/helpdoc/catalog.go` registers concise 1-2 sentence summaries for `pending-commits`, `pc`, `new-commands`, `nc`.
- `helpAliases` in `cli/helpdoc/print.go` maps `pc` -> `pending-commits` and `nc` -> `new-commands`.

### 6.5 `.agents/skills/gitmap/SKILL.md`
- **Replacement Matrix (Section 2):**
  - Maps manual `git status` loops -> `gitmap pc` (`gitmap pending-commits`).
  - Maps guessing new commands from git log -> `gitmap nc` (`gitmap new-commands`).
- **Cheat Sheet (Section 2 & 8):**
  - Adds `gitmap pc` with flags `--refresh`, `--no-cache`, `--backup`, `--json`.
  - Adds `gitmap nc` with flags `--filter`, `-f`, `-q`, `--category`, `--limit`, `--json`.
- **Operational Guardrails (Rule 9):**
  - Mandates using `gitmap pc` for multi-repo status.
  - Explains respecting 90s cache TTL and using `--refresh` or `--no-cache` after mutations.
  - Explains hierarchical tree remediation (`Option 1` vs `Option 2`) and batch footer remediation (`gitmap cpar`).
  - Explains using `--backup` for offline reviews.
  - Explains using `gitmap nc` for tool discovery.

---

## 7. Acceptance Criteria (AC1 Through AC8)

### AC1: Unified `UNCOMMITTED` Column
- The pending commits table replaces separate `DIRTY`, `UNTRACK`, `MODIF`, and `STAGE` columns with a single `UNCOMMITTED` column.
- The metric equals `UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`.
- Renders clean `-` or `0` for clean repos, and highlighted integer count for dirty repos.

### AC2: Repository Version and Branch (`VER/BRANCH`)
- The table displays `<version>/<branch>` if a SemVer tag exists (e.g. `v6.523.1/main`), or `<branch>` if unversioned.
- Column width is standardized at 16 characters and safely truncated with `…` to preserve 80-character box alignment.

### AC3: Batch Remediation Footer Command
- When dirty repositories exist, the table footer renders:
  `Fleet Remediation: gitmap cpar "wip: save changes"`
- Clean status output omits the fleet remediation line.

### AC4: Hierarchical Tree-View Remediation Without Clipping
- Every dirty repository row renders two indented tree branches:
  - `├── Option 1: git -C "<repo>" add -A && git commit -m "wip: save changes" && git push`
  - `└── Option 2: git -C "<repo>" stash -u`
- Tree lines are right-padded to 78 inner characters and enclosed in `│ ... │` borders without wrapping or clipping.
- Clean repositories do not render tree lines.

### AC5: SQLite Status Cache with 90-Second TTL
- Primary cache stored in `store.BinaryDataDir()` / `pending_commits_cache.db`.
- Records younger than 90 seconds with matching `HEAD` SHA return cached counts instantly (<5ms).
- Records older than 90 seconds are deleted immediately on access and never trusted.
- Bypass flags `--no-cache` and `--refresh` work as specified.

### AC6: SQLite Status Backup Database & Deserialization
- Persistent backup DB stored in `store.BinaryDataDir()` / `pending_commits_backup.db`.
- Holds point-in-time snapshots of repository status and JSON payloads.
- Records are NOT auto-purged by the 90-second TTL.
- `--backup` flag allows serving backed-up status snapshots again.
- Deserialization contract unpacks `payload_json` to restore `PendingFiles` and `UnpushedCommitSHAs`.

### AC7: New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)
- Both `gitmap new-commands` and `gitmap nc` invoke the catalog engine.
- Indexes 100 new commands across 10 functional categories (10 commands each).
- Supports `--limit` / `-n`, `--filter` (both `-f` and `-q`), `--category` / `-c`, and `--json` / `-j`.
- Displays syntax, version, category, description, and concrete copy-pasteable examples for every command.

### AC8: CLI Help Documentation & LLM Skills Sync
- `cli/constants/constants_cli.go` defines all command, alias, help, and flag constants.
- `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md` provide complete user documentation.
- `cli/helpdoc/catalog.go` and `cli/helpdoc/print.go` register help summaries and aliases.
- `.agents/skills/gitmap/SKILL.md` includes `gitmap pc` and `gitmap nc` in replacement matrix, cheat sheet, and operational guardrail rule 9.

---

## 8. Implementation Standards & Invariants

1. **Affirmative Booleans:** All boolean fields, variables, parameters, and return values MUST use affirmative naming (`isDirty`, `isClean`, `hasUncommitted`, `hasUpstream`, `isJSON`, `isHelpRequested`). Negative or inverted booleans (`notClean`, `noDirty`, `skipCache`) are strictly forbidden.
2. **Function Sizing Guard (8–15 LOC Cap):** Decompose all logic into small, single-purpose functions.
3. **Structured Error Handling:** All errors must be wrapped with `appfault.AppError` and standard error codes (`E9001`, `E9002`, `E9003`).
4. **Strict Relative Git Paths:** All file paths cited in documentation, error messages, and tests must be relative to repository root (`02-spec/...`, `cli/...`, `.ai-memory/...`). Absolute paths and `file:///` URIs are forbidden.
5. **No Git Commands:** Never execute git commands directly during multi-agent subtask execution.
