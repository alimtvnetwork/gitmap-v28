---
name: app-db-architecture
description: Design, implement, and maintain Split SQLite databases, Casbin RBAC schemas, PascalCase tables, and positive boolean fields.
---

# Split SQLite & Database Architecture Guidelines (`app-db-architecture`)

This skill guides the design, schema management, and querying of database layers in alignment with `02-spec/02-coding-guidelines/01-cross-language/07-database-naming.md`, `02-spec/04-database-conventions/02-naming-conventions.md`, and `02-spec/05-split-db-architecture/readme.md`.

## 1. Architecture Principles

1. **Split SQLite Architecture:**
   - **Root DB (`root.db`):** Global configurations, tenants, machine sessions, cross-repo leases, and system-wide metadata.
   - **App DB (`app.db`):** Core application entities, domain models, and business data.
   - **Session DB (`session.db`):** Transient user sessions, ephemeral tokens, and cache entries.
   - Integrated with Casbin RBAC model for fine-grained authorization checks.

2. **Schema & Table Naming Standards (PascalCase):**
   - Table names must strictly use **singular PascalCase** (e.g. `User`, `AgentSite`, `Transaction`).
   - Primary keys must follow the `{TableName}Id` format (e.g. `UserId`, `AgentSiteId`, `TransactionId`).
   - Column names must use **PascalCase** (e.g. `PluginSlug`, `CreatedAt`, `IsActive`, `DisplayName`) to align with `02-spec/02-coding-guidelines/01-cross-language/07-database-naming.md` and `02-spec/04-database-conventions/02-naming-conventions.md`.
   - Foreign key columns must use the exact same name as the referenced primary key (`AgentSiteId` references `AgentSite.AgentSiteId`).

3. **Boolean Column Discipline (Positive Only):**
   - Always prefix booleans with positive prefixes `Is` or `Has` (e.g. `IsActive`, `HasLicense`).
   - Total ban on negative column names (`IsNotActive`, `HasNoLicense` are forbidden).
   - Never evaluate boolean columns explicitly against `1` or `true` in application logic; use implicit checks (`if user.IsActive`).
   - For database columns representing negative states, provide computed/inversed properties in the application layer (`IsInactive() bool { return !u.IsActive }`) rather than storing negative column names.

4. **Foreign Keys & Explicit Joins:**
   - Enforce explicit foreign key constraints.
   - Avoid implicit ORM joins across split database boundaries.

5. **Repository Pattern (mandatory, language-agnostic):**
   - All SQL and data access lives behind a repository interface: one repository per aggregate or table group. Services, handlers, and use-cases NEVER embed SQL strings, query builders, or ORM calls directly.
   - Repository methods are intention-named (`FindByID`, `ListActive`, `Insert`, `Update`, `SoftDelete`) — never `RunQuery`, `ExecuteRaw`, or stringly-typed table/column parameters from callers.
   - Dialect and driver specifics stay inside the repository implementation; swapping the database means swapping the implementation, never the callers. This shape is identical in Go, TypeScript, Python, Rust, and PHP; only the idiom changes.
6. **Unit of Work:** Multi-repository transactions go through a unit-of-work abstraction; no ad-hoc transaction handling scattered in services.
7. **Repository Error Boundary:** Repository errors are wrapped with context (which repository, which operation) at the repository boundary before crossing into services (see `cg-error-management`).

---

## 2. Step 0 Preflight with Data Schema Verification

Before creating tables, running migrations, or claiming subtasks, inspect and verify the SQLite schema:

```bash
gitmap agent task status --ddl
gitmap agent task status --json
```

- **DDL & Schema Verification Protocol:** Inspect tables (`ParentTask`, `Subtask`, `AgentActionLog`), affirmative boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and foreign key linkages.
- **Forensic Autopsy Check:** Check `diagnostics.hasCrashesDetected`. If true, inspect `diagnostics.crashedAgents` to identify failed subtasks, targeting files, and last logged actions to facilitate clean resume.

---

## 3. 3-Tier Database Location Transparency & Diagnostic Suite

Whenever executing database operations or reporting telemetry, agents MUST fulfill the Database Location Transparency Contract:

| Database Tier | Physical File Location | Access Mode & Concurrency Protocol | Data Retained & Ownership |
| :--- | :--- | :--- | :--- |
| **Tier 1: Machine-Local Root Database** | `~/.gitmap/agent-root.db` | Machine sessions, cross-repo leases | Machine Sessions & Cross-Repo Leases |
| **Tier 2: Repo-Cache Companion Storage** | `repo-cache/.heatmap/<ai-slug>/<repo-slug>/heatmap.db` | Isolated Split-DB with periodic sync (`gitmap heatmap rc sync`) | Cross-repository multi-agent actions, tool usage traces, commit checkpoints |
| **Tier 3: In-Flight Local Run Database** | `.ai-memory/temp-agents/<slug>/agent-task.db` | WAL Mode (`PRAGMA journal_mode=WAL; busy_timeout=5000`) | Task status, worker claims, subtasks, in-flight action logs (<1ms Go engine) |
| **Tier 3: Local Heatmap Cache** | `.ai-memory/heatmap.db` | WAL Mode + Read-Only Queries (`mode=ro`) | Cognitive action telemetry, file exploration reads/edits, contention leases |

### Executable Diagnostic Commands:

```bash
# Verify schema DDL and table layout
gitmap agent task status --ddl

# Query companion telemetry with zero locks (mode=ro)
gitmap heatmap rc query "SELECT ActionType, FilePath, ReasoningText FROM AgentActionLog ORDER BY AgentActionLogId DESC LIMIT 5" --format table

# Inspect master tables catalog and row counts
gitmap heatmap rc tables

# View specific table records
gitmap heatmap rc view AgentActionLog --limit 10

# Synchronize local heatmap cache with repo-cache
gitmap heatmap rc sync
```

---

## 4. Cognitive Heatmap Telemetry & Blast Radius Verification

1. **Pre-Edit Contention & Blast Radius Check:**
   - BEFORE modifying any database entity or shared schema file, inspect active file leases:
     `gitmap heatmap discuss [--window 15m] [--threshold 2] [--json]` (or `gitmap agent heatmap`)
2. **Cognitive Decision Logging:**
   - Log schema reads and migration decisions with rationale:
     `gitmap heatmap log -t <task-slug> -a "READ" -f "<path>" -r "Reasoning: <why schema is inspected>" -s "EXPLORATION"`
     `gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <schema migration rationale>" -s "EXECUTION"`

---

## 5. Routine Execution & Commit Policy

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`) after schema edits.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run unit test suites (`go test ./...`, `pytest`) during routine turns.
- **GitMap Atomic Commits:** Stage and commit exclusively via GitMap hyphen format (`gitmap cpf "<module> - <summary>"` or `gitmap cpb "<module> - <summary>"`). TOTAL BAN on colons inside arguments and raw `git commit`.

---

## 6. Verification Checklist

- [ ] Schema adheres to PascalCase singular table names (`User`, `AgentSite`, `Transaction`).
- [ ] Columns use PascalCase (`PluginSlug`, `CreatedAt`, `IsActive`) and `{TableName}Id` primary keys.
- [ ] Booleans strictly use positive prefixes (`Is`, `Has`); zero negative columns.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
- [ ] All SQL/data access behind a repository interface (one per aggregate/table group); zero SQL, query-builder, or ORM calls in services/handlers/use-cases.
- [ ] Repository methods intention-named; no `RunQuery`/`ExecuteRaw`/stringly-typed callers; multi-repository transactions via unit-of-work only.
- [ ] Repository errors wrapped with repository + operation context at the boundary.
