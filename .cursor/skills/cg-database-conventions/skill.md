---
name: cg-database-conventions
description: Autonomously audits, refactors, and validates repository-wide database conventions using SQLite split-db patterns, PascalCase rules, and ORM standards.
---

# Database Conventions Coding Guidelines (`cg-database-conventions`)

This skill provides autonomous audit, refactoring, and validation of repository-wide database management based on `02-spec/02-coding-guidelines/01-cross-language/07-database-naming.md`, `02-spec/04-database-conventions/`, and `02-spec/05-split-db-architecture/readme.md`.

## 1. Concurrency Invariants & Header Parameters

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Fleet Sync, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

## 2. Core Architectural Invariants

1. **PascalCase Singular Tables**: Table names must strictly use **singular PascalCase** (e.g. `User`, `AgentSite`, `Transaction`). Never pluralize table names.
2. **PascalCase Columns**: Column names must strictly use **PascalCase** (`PluginSlug`, `CreatedAt`, `IsActive`, `DisplayName`) across all database entities.
3. **Primary Keys**: Primary keys must follow the `{TableName}Id` format (e.g. `UserId`, `AgentSiteId`, `TransactionId`). Never use bare `Id`.
4. **Foreign Keys**: Foreign key columns must use the exact same name as the referenced primary key (`AgentSiteId` references `AgentSite.AgentSiteId`).
5. **Split SQLite Architecture**: Split database layout across `root.db` (tenants, machine sessions, leases), `app.db` (domain entities), and `session.db` (ephemeral sessions, cache).
6. **Positive Boolean Columns**: Always use positive `Is` or `Has` prefixes (`IsActive`, `HasLicense`). Total ban on negative column names (`IsNotActive`, `HasNoLicense`).
7. **Timestamps as Epoch Integers**: Store all timestamps as `INTEGER` representing Unix epoch seconds in UTC (`CreatedAt INTEGER NOT NULL DEFAULT (unixepoch())`).
8. **Entity & Transactional Descriptions**: Entity/reference tables MUST include `Description TEXT NULL`. Transactional tables MUST include `Notes TEXT NULL` and `Comments TEXT NULL`.

---

## 3. Step 0 Preflight with Data Schema Verification

Before executing database migrations, generating models, or modifying schemas, inspect the SQLite schema:

```bash
gitmap agent task status --ddl
gitmap agent task status --json
```

- **DDL & Schema Verification Protocol:** Inspect tables (`ParentTask`, `Subtask`, `AgentActionLog`), affirmative boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and foreign key linkages.
- **Forensic Autopsy Check:** Check `diagnostics.hasCrashesDetected`. If true, inspect `diagnostics.crashedAgents` to identify failed subtasks, targeting files, and last logged actions to facilitate clean resume.

---

## 4. 3-Tier Database Location Transparency & Diagnostic Suite

Whenever executing database operations or reporting telemetry, agents MUST fulfill the Database Location Transparency Contract:

| Database Tier | Physical File Location | Access Mode & Concurrency Protocol | Data Retained & Ownership |
| :--- | :--- | :--- | :--- |
| **Tier 1: Machine-Local Root Database** | `~/.gitmap/agent-root.db` | Machine sessions, cross-repo leases | Machine Sessions & Cross-Repo Leases |
| **Tier 2: Repo-Cache Companion Storage** | `repo-cache/.heatmap/<ai-slug>/<repo-slug>/heatmap.db` | Isolated Split-DB with periodic sync (`gitmap heatmap rc sync`) | Cross-repository multi-agent actions, tool usage traces, commit checkpoints |
| **Tier 3: In-Flight Local Run Database** | `.ai-memory/temp-agents/<slug>/agent-task.db` | WAL Mode (`PRAGMA journal_mode=WAL; busy_timeout=5000`) | Task status, worker claims, subtasks, in-flight action logs (<1ms Go engine) |
| **Tier 3: Local Heatmap Cache** | `.ai-memory/heatmap.db` | WAL Mode + Read-Only Queries (`mode=ro`) | Cognitive action telemetry, file exploration reads/edits, contention leases |

### Executable Diagnostic Commands:

```bash
# Verify schema DDL and table definitions
gitmap agent task status --ddl

# Pre-edit contention and lease check
gitmap heatmap discuss [--window 15m] [--threshold 2] [--json]

# Cognitive file exploration and schema modification telemetry
gitmap heatmap log -t <task-slug> -a "READ" -f "<path>" -r "Reasoning: <explanation>" -s "EXPLORATION"
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <schema-refactor-rationale>" -s "EXECUTION"

# Real-time continuous activity stream & tree
gitmap heatmap tree [-t <task-slug>]
gitmap heatmap live

# Read-only companion query and schema inspection
gitmap heatmap rc query "SELECT ActionType, FilePath, ReasoningText FROM AgentActionLog ORDER BY AgentActionLogId DESC LIMIT 5" --format table
gitmap heatmap rc tables
gitmap heatmap rc view AgentActionLog --limit 10
gitmap heatmap rc sync
```

---

## 5. Routine Execution & GitMap Hyphen Commit Policy

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`) after schema edits.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run unit test suites (`go test ./...`, `pytest`) during routine turns.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 6. Verification Checklist

- [ ] Table names are singular PascalCase (`User`, `AgentSite`, `Transaction`).
- [ ] Column names use PascalCase (`PluginSlug`, `CreatedAt`, `IsActive`) and `{TableName}Id` keys.
- [ ] Boolean columns strictly use positive prefixes (`Is`, `Has`); zero negative columns.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
