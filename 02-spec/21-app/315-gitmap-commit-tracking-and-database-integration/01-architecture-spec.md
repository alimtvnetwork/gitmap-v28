# Architecture Specification: Commit Info, Universal SHA Resolution, and Cross-Database AI Provenance

> **Module:** `02-spec/21-app/315-gitmap-commit-tracking-and-database-integration`  
> **Sequence ID:** `315`  
> **Status:** `Draft / Implementation-Ready`  
> **Subtask 01:** [.ai-memory/plans/subtasks/315-gitmap-commit-tracking-and-database-integration/01-subtask-commit-info-command.md](../../../.ai-memory/plans/subtasks/315-gitmap-commit-tracking-and-database-integration/01-subtask-commit-info-command.md)  
> **Subtask 02:** [.ai-memory/plans/subtasks/315-gitmap-commit-tracking-and-database-integration/02-subtask-sha-resolver-and-tree-view.md](../../../.ai-memory/plans/subtasks/315-gitmap-commit-tracking-and-database-integration/02-subtask-sha-resolver-and-tree-view.md)  
> **Version:** `1.0.0`  
> **Updated:** `2026-10-23`  

---

## 1. Executive Summary

This specification establishes the architectural foundation for GitMap's deep commit inspection subsystem (`gitmap commit info`, aliased to `gitmap commit-info` and `gitmap ci`). The subsystem bridges raw Git object graphs with GitMap's multi-tier SQLite state, delivering universal commit identification, cross-database cognitive provenance, and high-clarity terminal tree visualization.

Key capabilities delivered by this architecture:
1. **4-Data-Layer Integration Topology:** Unifies data across raw Git object storage, the dedicated Split-DB (`repodb/commit_history.db`), Heatmap telemetry databases (`HeatmapCommitLog`, `HeatmapAiSession`, `HeatmapActionLog`), and Multi-Agent Orchestration databases (`ai_agents.db` and `agent-task.db`).
2. **Universal 4-Stage SHA Resolution Pipeline:** Resolves full 40-character SHAs, standard 7-character short SHAs, arbitrary prefixes, suffixes, and infix substrings across all branches and historical records without requiring full checkout.
3. **Machine and Agent Origin Attribution:** Formally distinguishes local AI agent commits, foreign worker/fleet commits, and human developer commits by matching hostnames, session UUIDs, and database audit logs.
4. **Rich Terminal Tree Visualization:** Renders structured commit lineage, database cross-references, file mutation metrics, and agent reasoning traces using GitMap's standardized terminal drawing glyphs (`constants.TreeBranch`, `constants.TreeCorner`, `constants.TreePipe`, and `constants.TreeSpace`).

---

## 2. 4-Data-Layer Integration Topology

GitMap reconciles Git commit hashes with cognitive traces across four distinct persistence layers:

```
+---------------------------------------------------------------------------------------+
|                       4-Data-Layer Commit Integration Topology                        |
+---------------------------------------------------------------------------------------+
|  Layer 1: Native Git Object Repository                                                |
|  - Raw objects (.git/objects), refs, log history, packfiles                           |
|  - Low-level Git operations: rev-parse, show, log, cat-file, diff-tree                |
+---------------------------------------------------------------------------------------+
|  Layer 2: Dedicated Commit Split-DB                                                   |
|  - Location: <effectiveDataDir>/repodb/commit_history.db                              |
|  - Primary Record: CommitHistoryRecord                                                |
|  - Indexed SHA, ShortSha, RepoSlug, Author, AI markers, CommitMessage                 |
+---------------------------------------------------------------------------------------+
|  Layer 3: Heatmap Cognitive Telemetry DBs                                             |
|  - Locations: heatmap.db, split heatmap/<repoSlug>/sql.db, repo-cache .heatmap/       |
|  - Tables: HeatmapCommitLog, HeatmapAiSession, HeatmapActionLog                       |
|  - Traces cognitive sessions, reasoning states, tool actions, and step budgets        |
+---------------------------------------------------------------------------------------+
|  Layer 4: Multi-Agent Orchestration DBs                                               |
|  - Tier 1 Master: .ai-memory/temp-agents/ai_agents.db (ParentTaskRegistry, Agent)     |
|  - Tier 2 Run-Scoped: <runDir>/agent-task.db (ParentTask, Subtask, AgentActionLog)     |
|  - Traces parent task plans, agent role claims, subtask codes, and execution evidence |
+---------------------------------------------------------------------------------------+
```

### 2.1 Layer 1: Native Git Object Repository
- **Scope:** Local working tree repository.
- **Access Pattern:** Read-only invocation of low-level Git commands using process execution wrappers.
- **Operations:**
  - `git rev-parse --verify <ref>^{commit}`: Fast resolution of valid commit references.
  - `git log -1 --format=... <sha>`: Extraction of raw author name, email, Unix timestamp, parent SHAs, and subject/body.
  - `git show --numstat --format= <sha>`: Detailed enumeration of changed file paths, added lines, and deleted lines.

### 2.2 Layer 2: Dedicated Commit Split-DB (`repodb/commit_history.db`)
- **Scope:** Host-level Split-DB under `<effectiveDataDir>/repodb/commit_history.db`.
- **Purpose:** High-speed lookup index preventing repetitive, heavy Git subprocess spawning during batch queries and cross-repository audits.
- **Schema Specification:**
  ```sql
  CREATE TABLE IF NOT EXISTS CommitHistoryRecord (
      CommitHistoryRecordId INTEGER PRIMARY KEY AUTOINCREMENT,
      RepoSlug TEXT NOT NULL,
      CommitSha TEXT NOT NULL UNIQUE,
      ShortSha TEXT NOT NULL,
      AuthorName TEXT NOT NULL DEFAULT '',
      AuthorEmail TEXT NOT NULL DEFAULT '',
      CommitMessage TEXT NOT NULL DEFAULT '',
      CommittedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
      IsAiCommitted INTEGER NOT NULL DEFAULT 0,
      AiSessionId INTEGER NULL,
      AiName TEXT NOT NULL DEFAULT '',
      AgentRole TEXT NOT NULL DEFAULT '',
      TaskSlug TEXT NOT NULL DEFAULT '',
      MachineHostname TEXT NOT NULL DEFAULT '',
      FilesChangedCount INTEGER NOT NULL DEFAULT 0,
      InsertionsCount INTEGER NOT NULL DEFAULT 0,
      DeletionsCount INTEGER NOT NULL DEFAULT 0,
      ParentSha TEXT NOT NULL DEFAULT '',
      Notes TEXT NULL,
      CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  CREATE UNIQUE INDEX IF NOT EXISTS IdxCommitHistory_Sha ON CommitHistoryRecord(CommitSha);
  CREATE INDEX IF NOT EXISTS IdxCommitHistory_ShortSha ON CommitHistoryRecord(ShortSha);
  CREATE INDEX IF NOT EXISTS IdxCommitHistory_Repo ON CommitHistoryRecord(RepoSlug);
  CREATE INDEX IF NOT EXISTS IdxCommitHistory_Task ON CommitHistoryRecord(TaskSlug);
  ```

### 2.3 Layer 3: Heatmap Cognitive Telemetry DBs
- **Scope:** Workstation root (`heatmap.db`), repo-specific split (`heatmap/<slug>/sql.db`), and shared repo-cache (`.heatmap/<seq>-<ai>/`).
- **Target Tables:**
  - `HeatmapCommitLog`: Maps `CommitSha` to `TaskSlug`, `RepoSlug`, `CommittedAt`, `FilesCommittedCount`, and `PushedSuccessfully`.
  - `HeatmapAiSession`: Maps AI sessions by `SessionUuid` or `AiName`, identifying `MachineHostname`, `OsPlatform`, and `AgentRole`.
  - `HeatmapActionLog`: Reconstructs the exact cognitive action log entries preceding and producing the commit.

### 2.4 Layer 4: Multi-Agent Orchestration DBs
- **Scope:** Global agent coordination (`ai_agents.db`) and task-scoped runs (`agent-task.db`).
- **Target Tables:**
  - `ParentTaskRegistry` / `ParentTask`: Contains task plan titles, total step budgets, completed step counts, and active status.
  - `Subtask`: Associates the commit with individual task codes (e.g., `315-subtask-01`), assigned agent roles, and evidence links.
  - `AgentActionLog`: Captures fine-grained file modification line bounds, executed commands, and tool duration metrics.

---

## 3. Universal SHA Resolution Engine

The universal SHA resolver (`cli/cmdcommit/commit_info_resolve.go`) accepts any commit identifier: full 40-char SHA, 7-char short SHA, prefix, suffix, branch name, tag, or arbitrary substring. It processes inputs through a bounded 4-stage pipeline.

```
User Input Query (e.g., "7a3b4c", "c12d9", "fix-typo-in-parser")
       |
       v
+-------------------------------------------------------------+
| Stage 1: Native Git rev-parse Fast-Path                     |
| Command: git rev-parse --verify <input>^{commit}            |
| Match? -----> YES -----> Extract 40-char Full SHA           |
+-------------------------------------------------------------+
       | NO
       v
+-------------------------------------------------------------+
| Stage 2: Commit Split-DB Indexed Lookup                     |
| Database: repodb/commit_history.db                          |
| Query: WHERE CommitSha = ? OR ShortSha = ? OR               |
|        CommitSha LIKE ? || '%' OR CommitSha LIKE '%' || ?   |
| Match? -----> YES -----> Return Full Record                 |
+-------------------------------------------------------------+
       | NO
       v
+-------------------------------------------------------------+
| Stage 3: Heatmap Telemetry Database Lookup                  |
| Database: heatmap.db / .heatmap/                            |
| Query: SELECT CommitSha FROM HeatmapCommitLog               |
|        WHERE CommitSha LIKE '%' || ? || '%'                 |
| Match? -----> YES -----> Resolve Full SHA                   |
+-------------------------------------------------------------+
       | NO
       v
+-------------------------------------------------------------+
| Stage 4: Streaming Full Git Object Substring Search         |
| Command: git rev-list --all                                 |
| Stream: Scan commit SHAs for substring match                |
| Match? -----> YES -----> Return First Unique SHA            |
+-------------------------------------------------------------+
       | NO
       v
   Error: CommitNotFound ("no commit matching query found")
```

### 3.1 Stage 1: Native Git `rev-parse` Fast-Path
1. Execute `git rev-parse --verify --quiet <query>^{commit}` in the target repository directory.
2. If the exit code is 0 and output is a valid 40-hex-character string, resolution terminates immediately with high confidence.
3. Handles standard Git references: exact full SHAs, prefixes, HEAD references, branch heads, and tags.

### 3.2 Stage 2: Commit Split-DB Indexed Lookup
1. Open `<effectiveDataDir>/repodb/commit_history.db` in read-only mode (`mode=ro`).
2. Execute parameterized query matching exact SHA, prefix, or suffix:
   ```sql
   SELECT CommitSha, RepoSlug, CommitMessage FROM CommitHistoryRecord
   WHERE RepoSlug = ? AND (
       CommitSha = ? OR
       ShortSha = ? OR
       CommitSha LIKE ? || '%' OR
       CommitSha LIKE '%' || ?
   ) LIMIT 5;
   ```
3. If exactly one record matches, return the resolved SHA. If multiple records match, return an ambiguous match error listing candidate short SHAs.

### 3.3 Stage 3: Heatmap Telemetry Database Lookup
1. Scan active heatmap databases for matching `HeatmapCommitLog` entries.
2. Search query:
   ```sql
   SELECT CommitSha, TaskSlug, RepoSlug FROM HeatmapCommitLog
   WHERE CommitSha LIKE '%' || ? || '%'
   ORDER BY CommittedAt DESC LIMIT 5;
   ```
3. If matches are found, resolve the full SHA and seed the metadata query with known `TaskSlug`.

### 3.4 Stage 4: Streaming Full Git Object Substring Search
1. When queries are arbitrary substrings (e.g., middle characters of a SHA) not resolvable via Git's native prefix parser:
2. Stream stdout from `git rev-list --all` line by line.
3. Compare each 40-character line using `strings.Contains(sha, query)`.
4. Bounded scan: limit scan to the first 25,000 commits or terminate on the first unambiguous match to protect execution latency.

---

## 4. Cross-Database Provenance & Origin Attribution

A critical requirement of GitMap is answering: **Who made this commit, and was it generated by an AI agent on this machine, an AI agent on another machine, or a human?**

```
+---------------------------------------------------------------------------------------+
|                               Origin Attribution Logic                                |
+---------------------------------------------------------------------------------------+
|  Condition 1: IsAiCommitted == true AND MachineHostname == LocalHost                  |
|  => Result: LOCAL AI COMMIT                                                           |
|     Attribution: "Local Machine AI Agent (<AiName> / <AgentRole>)"                    |
|     Provenance: Attached to local session UUID, task slug, and action logs            |
+---------------------------------------------------------------------------------------+
|  Condition 2: IsAiCommitted == true AND MachineHostname != LocalHost                  |
|  => Result: FOREIGN AI WORKER COMMIT                                                  |
|     Attribution: "Remote Cluster Worker AI (<AiName> on <MachineHostname>)"           |
|     Provenance: Recorded by remote fleet node or synchronized via repo-cache          |
+---------------------------------------------------------------------------------------+
|  Condition 3: IsAiCommitted == false AND (No AI Session OR No Task Log)               |
|  => Result: EXTERNAL / HUMAN COMMIT                                                   |
|     Attribution: "The current machine's AI has not committed this                     |
|                   (likely another worker/node/human)"                                 |
|     Provenance: Native Git author credentials (<AuthorName> <AuthorEmail>)            |
+---------------------------------------------------------------------------------------+
```

### 4.1 Attribution Classification Algorithm
1. Retrieve commit metadata from Git and match against `CommitHistoryRecord` and `HeatmapCommitLog`.
2. Inspect `IsAiCommitted` flag:
   - If `!isAiCommitted`, classify as `OriginHumanOrExternal`.
   - If `isAiCommitted`, compare `record.MachineHostname` against the host machine's current hostname (`os.Hostname()`).
3. If `record.MachineHostname` matches the local host, verify existence of local task artifacts in `.ai-memory/temp-agents/` or `agent-task.db`:
   - If local task found: label as `OriginLocalAI` with full task, subtask, and step references.
   - If not found locally: label as `OriginLocalAIHistoric`.
4. If `record.MachineHostname` is non-empty and differs from local host, label as `OriginForeignAI` and format the remote hostname badge.
5. If no AI metadata exists anywhere across databases, render the authoritative status:
   `The current machine's AI has not committed this (likely another worker/node/human)`.

---

## 5. Rich Terminal Tree Visualization

The commit tree formatter (`cli/cmdcommit/commit_info_tree.go`) renders an ASCII tree structure utilizing GitMap's terminal glyph constants.

### 5.1 Terminal Tree Layout Design

```text
commit 7a3b4c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b (short: 7a3b4c8)
├── Author:      Alex Developer <alex@example.com>
├── Date:        2026-10-23 21:15:30 +0800 (15 minutes ago)
├── Message:     feat(commit): implement universal sha resolution engine
├── Provenance:  Local Machine AI Agent (antigravity / ORCHESTRATOR)
│   ├── Hostname:     WORKSTATION-01
│   ├── Session UUID: sess-9a8b-4c2d-11ef-8412
│   └── Task Slug:    315-gitmap-commit-tracking-and-database-integration
├── Database Cross-References:
│   ├── Split-DB:     repodb/commit_history.db (Record #1042)
│   ├── Heatmap:      heatmap.db (HeatmapCommitLog #812, 14 Actions logged)
│   └── Agent Tasks:  agent-task.db (Subtask 315-subtask-01: COMPLETED)
├── File Mutations:
│   ├── cli/cmdcommit/commit_info.go         (+142, -0)
│   ├── cli/cmdcommit/commit_info_resolve.go (+89, -0)
│   └── cli/cmd/rootcore.go                  (+4, -1)
└── Summary:     3 files changed, 235 insertions(+), 1 deletion(-)
```

### 5.2 Terminal Glyph Usage
- `constants.TreeBranch`: `"├── "` connects intermediate tree nodes.
- `constants.TreeCorner`: `"└── "` connects the final node of a subtree.
- `constants.TreePipe`: `"│   "` provides vertical indentation guides for child hierarchies.
- `constants.TreeSpace`: `"    "` provides blank indentation padding for terminal leaves.

### 5.3 Color Tokens & Styling
- Full & Short SHA: `constants.ColorNeonGreen` or `constants.ColorNeonGold`.
- Field Labels (`Author:`, `Date:`, `Message:`): `constants.ColorCyan`.
- AI Provenance Badges: `constants.ColorMagenta` for AI agents, `constants.ColorYellow` for foreign nodes, `constants.ColorDim` for external/human.
- Database Cross-References: `constants.ColorWhite` with `constants.ColorDim` file paths.
- File Mutations: `constants.ColorGreen` for insertions (`+`), `constants.ColorRed` for deletions (`-`).

### 5.4 Dual Output Modes
- **Human Terminal View (Default):** ANSI colorized ASCII tree rendered to stdout.
- **Machine JSON View (`--json`):** Emits structured `CommitInfoReport` with zero terminal formatting or ANSI escapes for downstream automation.

---

## 6. Data Models and Contracts

### 6.1 Core Data Structures (`cli/cmdcommit/commit_info_types.go`)

```go
package cmdcommit

import "time"

// OriginType designates the origin classification of a commit.
type OriginType string

const (
	OriginLocalAI         OriginType = "LOCAL_AI"
	OriginForeignAI       OriginType = "FOREIGN_AI"
	OriginHumanOrExternal OriginType = "HUMAN_OR_EXTERNAL"
)

// CommitFileDelta captures changes to a single file in a commit.
type CommitFileDelta struct {
	FilePath    string `json:"filePath"`
	Insertions  int    `json:"insertions"`
	Deletions   int    `json:"deletions"`
	IsBinary    bool   `json:"isBinary"`
}

// CommitProvenance holds origin and attribution metadata.
type CommitProvenance struct {
	OriginType      OriginType `json:"originType"`
	AttributionText string     `json:"attributionText"`
	AiName          string     `json:"aiName,omitempty"`
	AgentRole       string     `json:"agentRole,omitempty"`
	SessionUuid     string     `json:"sessionUuid,omitempty"`
	TaskSlug        string     `json:"taskSlug,omitempty"`
	MachineHostname string     `json:"machineHostname,omitempty"`
	IsLocalMachine  bool       `json:"isLocalMachine"`
}

// CommitDbCrossReferences captures linkages across SQLite databases.
type CommitDbCrossReferences struct {
	HasSplitDbRecord  bool   `json:"hasSplitDbRecord"`
	SplitDbPath       string `json:"splitDbPath,omitempty"`
	SplitDbRecordId   int64  `json:"splitDbRecordId,omitempty"`
	HasHeatmapRecord  bool   `json:"hasHeatmapRecord"`
	HeatmapDbPath     string `json:"heatmapDbPath,omitempty"`
	HeatmapActionCount int   `json:"heatmapActionCount"`
	HasAgentTask      bool   `json:"hasAgentTask"`
	AgentTaskDbPath   string `json:"agentTaskDbPath,omitempty"`
	SubtaskCode       string `json:"subtaskCode,omitempty"`
	SubtaskStatus     string `json:"subtaskStatus,omitempty"`
}

// CommitInfoReport represents the unified commit report payload.
type CommitInfoReport struct {
	FullSha          string                  `json:"fullSha"`
	ShortSha         string                  `json:"shortSha"`
	AuthorName       string                  `json:"authorName"`
	AuthorEmail      string                  `json:"authorEmail"`
	CommittedAt      time.Time               `json:"committedAt"`
	FormattedDate    string                  `json:"formattedDate"`
	RelativeDate     string                  `json:"relativeDate"`
	CommitMessage    string                  `json:"commitMessage"`
	ParentShas       []string                `json:"parentShas"`
	RepoSlug         string                  `json:"repoSlug"`
	Provenance       CommitProvenance        `json:"provenance"`
	CrossReferences  CommitDbCrossReferences `json:"crossReferences"`
	FileDeltas       []CommitFileDelta       `json:"fileDeltas"`
	TotalFiles       int                     `json:"totalFiles"`
	TotalInsertions  int                     `json:"totalInsertions"`
	TotalDeletions   int                     `json:"totalDeletions"`
}
```

---

## 7. Coding Guidelines & Invariant Adherence

All files implemented under this architecture specification must adhere strictly to repository coding standards:

1. **Relative Path Hygiene:** Zero hardcoded drive letters (`C:`, `D:`) or protocol schemes (`file:///`). All paths must be relative or resolved via standard workspace resolvers.
2. **Positive Booleans:** Boolean fields and variables must strictly use `is*` or `has*` prefixes (e.g., `isLocalMachine`, `hasSplitDbRecord`, `isResolved`, `hasChanges`). Negative booleans are banned.
3. **15-Line Function Limit:** Functions must be $\le 15$ lines, with $\le 8$ lines preferred. Long logic blocks must be decomposed into dedicated helper functions.
4. **Zero Nested Conditionals:** Nested `if` statements are forbidden. Flatten all conditionals using guard clauses, early returns, and map lookups.
5. **AppError Wrapping:** All error returns must wrap underlying errors using `*apperror.AppError` via `apperror.WrapSimple(err, "context")`.
6. **Zero Builds and Git Commands During Authoring:** No build, test, or Git commands may be executed during specification authoring.

---

## 8. Acceptance Criteria and Verification Protocol

1. **Universal SHA Resolution:**
   - 40-character full SHA resolves directly via Stage 1.
   - 7-character prefix resolves correctly via Stage 1 or Stage 2.
   - Arbitrary middle substring (e.g. 5 hex chars) resolves via Stage 3 or Stage 4.
   - Non-existent query returns descriptive `CommitNotFound` error without crashing.

2. **Attribution Correctness:**
   - Commits authored by local AI sessions show `LOCAL_AI` with session and task badges.
   - Commits matching foreign hostnames show `FOREIGN_AI` with worker hostname.
   - Commits without AI traces display: `"The current machine's AI has not committed this (likely another worker/node/human)"`.

3. **Rendering Fidelity:**
   - ASCII tree displays standard box-drawing glyphs (`├──`, `└──`, `│   `, `    `).
   - Terminal colors adhere to neon/cyan/dim palette.
   - `--json` flag emits valid, unescaped JSON payload matching `CommitInfoReport`.
