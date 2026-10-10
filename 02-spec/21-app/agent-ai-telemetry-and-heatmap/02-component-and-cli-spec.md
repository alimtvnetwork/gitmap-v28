# Component & CLI Specification: Antigravity & GitMap AI Agent Telemetry & Heatmap Suite

- **Feature Slug:** `agent-ai-telemetry-and-heatmap`
- **Module:** `cli/cmdagent`, `cli/cmdagentai`, `03-ai-scripts/46-agent-sqlite-task-manager.py`
- **Specification Version:** 2.0.0
- **Status:** APPROVED
- **Target Audience:** Autonomous AI Agents, Multi-Agent Workers, Lead Orchestrators, GitMap Engineers
- **Target Release:** Minor Bump

---

## 1. CLI Routing Architecture & Command Taxonomy

The telemetry suite exposes two primary command namespaces designed for zero-friction invocation by autonomous agents and human developers:

1. **`gitmap agent-ai` (alias: `gitmap aai`):**
   High-level workflow commands for task lifecycle, pre-touch editing declarations, interactive heatmaps, and continuous learning.
2. **`gitmap agent` (aliases: `gitmap ai-agents`, `gitmap agents`):**
   Low-level atomic engine commands for subtask queuing, file box look-ahead claims, and crash forensic autopsies.

```
gitmap agent-ai
├── create-task "<title or slug>" [--budget N]
├── subtask
│   └── create <parent-slug/id> "<subtask-title>" [--code <Code>] [--files <paths>] [--role <role>]
├── editing <subtask-or-task-id> "<file>" --action <action> --lines <x,y> --reasoning "<why>" [--agent <role>]
├── heatmap [<task-id-or-slug>]
│   ├── ls [N]
│   ├── search "<query>"
│   └── grep "<regex>"
├── subtask [<task-id-or-slug>]
└── learn "<text>" [--tag <tag>]

gitmap agent
└── subtask
    └── claim-files --task-id <id> --subtask <code/id> --files "<paths>" --agent "<role>" [--json]
```

---

## 2. Command Specifications & Syntax Contracts

### 2.1 Parent Task Creation (`gitmap agent-ai create-task`)

Initializes a parent task in the Tier 1 registry (`.ai-memory/temp-agents/ai_agents.db`) and automatically scaffolds the Tier 2 database (`.ai-memory/temp-agents/<slug>/agent-task.db`).

#### Syntax
```bash
gitmap agent-ai create-task "<title or slug>" [--budget <steps>] [--dir <custom-dir>] [--json]
```

#### Arguments & Flags
| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `<title or slug>` | String | Yes | - | Task name or slug (e.g. `"Repository Fix and Solution Implementation"`) |
| `--budget`, `-b` | Integer | No | `300` | Monotonic step budget for autonomous loop |
| `--dir`, `-d` | String | No | auto | Custom run directory relative to repo root |
| `--json` | Boolean | No | `false` | Emits structured JSON response |

#### Behavioral Logic
1. Sanitizes input string into a deterministic lowercase kebab-case slug:
   - Replaces non-alphanumeric characters with hyphens.
   - Truncates to max 6 primary words if overly verbose.
2. Inserts or updates row in Tier 1 `ParentTaskRegistry` table.
3. Creates directory `.ai-memory/temp-agents/<slug>/`.
4. Initializes Tier 2 SQLite DB `agent-task.db` with WAL mode and `ParentTask` record.
5. Emits `ParentTaskId`, `TaskSlug`, and resolved `DbPath`.

#### JSON Output Schema
```json
{
  "isSuccess": true,
  "parentTaskId": 1,
  "taskSlug": "repository-fix-and-solution-implementation",
  "taskName": "Repository Fix and Solution Implementation",
  "runDirectory": ".ai-memory/temp-agents/repository-fix-and-solution-implementation",
  "dbPath": ".ai-memory/temp-agents/repository-fix-and-solution-implementation/agent-task.db",
  "totalStepsBudget": 300,
  "status": "ACTIVE"
}
```

---

### 2.2 Subtask Creation (`gitmap agent-ai subtask create`)

Registers a discrete subtask under the designated parent task across both database tiers.

#### Syntax
```bash
gitmap agent-ai subtask create <parent-slug/id> "<subtask-title>" [--code <Code>] [--files <paths>] [--role <agent-role>] [--json]
```

#### Arguments & Flags
| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `<parent-slug/id>` | String | Yes | - | Parent task ID or slug (fuzzy-matched if partial) |
| `"<subtask-title>"` | String | Yes | - | Subtask title description |
| `--code`, `-c` | String | No | auto | Unique code (e.g. `Task-01`). Auto-increments if omitted |
| `--files`, `-f` | String | No | `""` | Comma-separated relative file paths initially allocated |
| `--role`, `-r` | String | No | `"Worker"`| Assigned agent role (e.g. `"Worker 01"`) |
| `--json` | Boolean | No | `false` | Emits structured JSON output |

#### Behavioral Logic
1. Resolves parent task from Tier 1 registry by exact ID, exact slug, or closest fuzzy match.
2. If `--code` is omitted, scans existing subtasks in Tier 2 DB and assigns `Task-01`, `Task-02`, etc.
3. Inserts row into Tier 2 `Subtask` table.
4. If `--files` are provided, registers them in `OwnedFilesJson` and mirrors into Tier 1 `FileClaim`.
5. Updates `TotalSubtasks` count in Tier 1 `ParentTaskRegistry`.

---

### 2.3 Subtask File Box Claiming (`gitmap agent subtask claim-files`)

Declares an agent's look-ahead file box. Enables early collision detection before files are touched.

#### Syntax
```bash
gitmap agent subtask claim-files \
  --task-id "<task-id-or-slug>" \
  --subtask "<subtask-id-or-code>" \
  --files "<comma-separated-paths>" \
  --agent "<agent-role>" \
  [--json]
```

#### Example Invocation
```bash
gitmap agent subtask claim-files \
  --task-id "command-development-for-release-tag-management" \
  --subtask "Task-01" \
  --files "cli/cmdfixreleasetags/audit_engine.go,cli/cmdfixreleasetags/delete_executor.go" \
  --agent "Worker 01"
```

#### Behavioral Logic
1. Normalizes all file paths to relative forward-slash paths; validates against repository tree.
2. Updates `OwnedFilesJson` in Tier 2 `Subtask` table.
3. Inserts active claims into Tier 1 `FileClaim` table.
4. **Collision Detection:** Queries Tier 1 `FileClaim` for other active subtasks holding identical files:
   - If collision detected: Inserts row into `CollisionEvent` table and prints high-contrast warning banner:
     ```
     ⚠️ FILE CLAIM COLLISION DETECTED:
     File: cli/cmdfixreleasetags/audit_engine.go
     Claimed by: Worker 01 (Task-01) AND Worker 02 (Task-02)
     Resolution: Disjoint file boxes mandated. Worker 02 must yield or re-partition!
     ```

---

### 2.4 Pre-Touch In-Flight Editing Declaration (`gitmap agent-ai editing`)

Mandatory pre-flight command that MUST be executed immediately before invoking any file mutation tool (`replace_file_content`, `write_to_file`, etc.).

#### Syntax
```bash
gitmap agent-ai editing <subtask-or-task-id> "<relative-file>" \
  --action <action> \
  --lines <start,end> \
  --reasoning "<rationale>" \
  [--agent "<role>"] \
  [--json]
```

#### Arguments & Flags
| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `<subtask-or-task-id>` | String/Int | Yes | - | Subtask ID, Task Code (`Task-01`), or Task Slug |
| `"<relative-file>"` | String | Yes | - | Strict relative path of the file to modify or read |
| `--action`, `-a` | String | Yes | - | Tool action: `replace_file_content`, `write_to_file`, `view_file` |
| `--lines`, `-l` | String | Yes | `0,0` | Target line range: `<start>,<end>` (e.g. `120,145`) |
| `--reasoning`, `-r` | String | Yes | - | Clear technical rationale explaining the modification |
| `--agent` | String | No | auto | Agent role (e.g. `"Worker 01"`). Auto-resolved from subtask if omitted |
| `--json` | Boolean | No | `false` | Emits structured JSON response |

#### Example Invocations
```bash
# Surgical file edit declaration:
gitmap agent-ai editing Task-01 "cli/cmdfixreleasetags/audit_engine.go" \
  --action "replace_file_content" \
  --lines 120,145 \
  --reasoning "Implement 5-point release tag audit classification logic"

# Legacy Python fallback equivalent:
python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action \
  --db ".ai-memory/temp-agents/command-development-for-release-tag-management/agent-task.db" \
  --subtask-id 1 \
  --agent "Worker 01" \
  --action "replace_file_content" \
  --file "cli/cmdfixreleasetags/audit_engine.go" \
  --details "Implement 5-point release tag audit classification logic"
```

#### Behavioral Logic
1. Disambiguates `<subtask-or-task-id>`:
   - If numeric integer: Matches `SubtaskId`.
   - If format `Task-NN`: Matches `TaskCode` under the active task.
   - If slug: Resolves parent task and matches currently claimed subtask for `--agent`.
2. Inserts row into Tier 2 `AgentActionLog` with status `'IN_PROGRESS'`.
3. Mirrors entry into Tier 1 `GlobalActionIndex` for cross-task searching and instant heatmap reflection.
4. Returns sub-millisecond confirmation (`[0.8ms] Action recorded: ActionLogId 42`).

---

### 2.5 Heatmap Inspection (`gitmap agent-ai heatmap`)

Generates real-time file churn rankings and heatmaps for a specific task or the entire workspace.

#### Syntax
```bash
gitmap agent-ai heatmap [<task-id-or-slug>] [--json]
```

#### Behavioral Logic
1. If `<task-id-or-slug>` is provided:
   - Resolves target task using closest fuzzy matching across Tier 1 `ParentTaskRegistry`.
   - Aggregates write claims from `FileClaim` and touches from `AgentActionLog`.
2. If omitted:
   - Evaluates active parent tasks from the last 48 hours and produces workspace-wide rollup.
3. Renders high-contrast ANSI table sorted by total heat score descending:

```
[1.2ms] agent-ai heatmap: command-development-for-release-tag-management

FILE PATH                                    CLAIMS  TOUCHES  HEAT SCORE  TIER
cli/cmdfixreleasetags/audit_engine.go             2        8          12  🔥🔥 CRITICAL
cli/cmdfixreleasetags/delete_executor.go          2        5           9  🔥 ACTIVE
cli/cmdfixreleasetags/types.go                    1        3           5  🔥 ACTIVE
02-spec/21-app/command-development/...            1        2           4  • TOUCHED
```

---

### 2.6 Heatmap Subcommand Family (`ls`, `search`, `grep`)

#### 1. List Recent Heatmaps (`gitmap agent-ai heatmap ls [N]`)
Lists the last $N$ task heatmaps (default $N=30$) with activity metrics.
```bash
gitmap agent-ai heatmap ls 30
```
- **Output:** Aligned table displaying `SLUG`, `STATUS`, `SUBTASKS`, `TOUCHES`, `ACTIVE AGENTS`, `LAST ACTIVITY`.

#### 2. Search Tasks & Heatmaps (`gitmap agent-ai heatmap search "<query>"`)
Fuzzy-searches task titles, slugs, and subtask descriptions in Tier 1 registry.
```bash
gitmap agent-ai heatmap search "release tags"
```
- Locates matching tasks and displays file churn rankings for the top match.

#### 3. High-Speed Regex Grep (`gitmap agent-ai heatmap grep "<regex>"`)
Executes an ultra-fast regex query over `GlobalActionIndex` in the Tier 1 master DB.
```bash
gitmap agent-ai heatmap grep "audit.*engine"
```
- Identifies every file, subtask, agent, and line range matching the regex across all historical sessions in $<10\text{ms}$.

---

### 2.7 Subtask View (`gitmap agent-ai subtask [<task-id-or-slug>]`)

Displays an interactive tree of subtasks, worker assignments, claimed file boxes, and progress state.
```bash
gitmap agent-ai subtask command-development-for-release-tag-management
```

```
Parent Task: command-development-for-release-tag-management [ACTIVE]
├── [DONE] Task-01: Audit Engine & Safety Guard (Worker 01)
│   ├── Files: cli/cmdfixreleasetags/audit_engine.go, cli/cmdfixreleasetags/safety_guard.go
│   └── Evidence: PASS exit 0 (go test ./cmdfixreleasetags)
└── [IN_PROGRESS] Task-02: CLI Routing & UI Integration (Worker 02)
    ├── Files: cli/cmdfix/fix_cmd.go, cli/cmdfixreleasetags/fix_release_tags.go
    └── Last Touch: cli/cmdfix/fix_cmd.go (lines 45-68) [0.4m ago]
```

---

### 2.8 Continuous Learning (`gitmap agent-ai learn "<text>"`)

Persists architectural conventions, discovered pitfalls, and domain rules into Tier 1 `AgentLearning`.

#### Syntax
```bash
gitmap agent-ai learn "<content>" [--tag <tag>] [--task <slug>] [--json]
```

#### Example
```bash
gitmap agent-ai learn "Always use captureStdout asynchronous reader on Windows to prevent pipe deadlock in pipeline tests" --tag "windows-testing"
```

#### Behavioral Logic
1. Stores record in `AgentLearning` with ISO-8601 timestamp.
2. Automatically available for semantic and tag queries by subsequent AI agents during Phase 1 preflight.

---

## 3. Implementation Blueprint for Go & Python Engines

### 3.1 Go Package Structure
- `cli/cmdagent/agent_claim.go`: Implements `claim-files` with collision detection.
- `cli/cmdagent/agent_heatmap.go`: Implements high-speed table rendering and heat score calculations.
- `cli/cmdagentai/agent_ai_cmd.go`: Top-level router for `gitmap agent-ai` / `gitmap aai`.
- `cli/cmdagentai/agent_ai_editing.go`: Pre-touch editing declaration handler.
- `cli/cmdagentai/agent_ai_learn.go`: Heuristic memory storage and retrieval.
- `cli/cmdagentai/agent_ai_heatmap_query.go`: `ls`, `search`, and `grep` implementations.
- `cli/constants/constants_cli.go`: Command names (`CmdAgentAI = "agent-ai"`, `CmdAgentAIAlias = "aai"`).

### 3.2 Python Fallback Script (`03-ai-scripts/46-agent-sqlite-task-manager.py`)
- Maintains 100% schema parity and command flag parity.
- Supports `init`, `add-subtasks`, `claim`, `claim-files`, `log-action`, `complete`, `diagnose`, `heatmap`.

---

## 4. Acceptance Criteria & Verification Protocol

| ID | Criterion | Verification Command | Expected Outcome |
|---|---|---|---|
| **AC-01** | Create Parent Task | `gitmap agent-ai create-task "Telemetry Test"` | Returns slug `telemetry-test`, scaffolds Tier 2 DB, exits 0 |
| **AC-02** | Create Subtask | `gitmap agent-ai subtask create telemetry-test "Subtask 01"` | Inserts into Tier 2 `Subtask`, assigns `Task-01`, exits 0 |
| **AC-03** | Claim File Box | `gitmap agent subtask claim-files -t telemetry-test --subtask Task-01 --files "a.go,b.go" -a "Worker 01"` | Updates `OwnedFilesJson`, mirrors to Tier 1 `FileClaim`, exits 0 |
| **AC-04** | Collision Detection | Second claim on `a.go` by `Worker 02` | High-contrast collision warning banner rendered, exits 0 |
| **AC-05** | Pre-touch Declaration | `gitmap agent-ai editing Task-01 "a.go" -a replace_file_content -l 10,20 -r "Refactor"` | Commits to `AgentActionLog` in $<2\text{ms}$, exits 0 |
| **AC-06** | Heatmap Rendering | `gitmap agent-ai heatmap telemetry-test` | Renders ANSI table with file ranks and heat scores, exits 0 |
| **AC-07** | Heatmap List & Grep | `gitmap agent-ai heatmap ls 10` & `gitmap agent-ai heatmap grep "a\.go"` | Returns formatted list and matched regex rows, exits 0 |
| **AC-08** | Continuous Learning | `gitmap agent-ai learn "Heuristic rule" --tag test` | Persists into `AgentLearning`, exits 0 |
| **AC-09** | Relative Paths Hygiene | `python linter-scripts/check-relative-paths.py -c` | 0 absolute path violations |
| **AC-10** | Positive Booleans | Schema inspection | Zero negative boolean flags (`is_not_*`, `disable_*` banned) |
