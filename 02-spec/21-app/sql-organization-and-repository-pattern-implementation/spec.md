# Specification: SQL Organization, Repository Pattern & AGY Prompt Management

## 1. Overview & Strategic Intent
GitMap's persistence layer operates on a Three-Tier SQLite Split-DB engine (`gitmap.db`, `installation.db`, `repodb/*.db`, `heatmap.db`). To achieve enterprise-grade maintainability, predictability, and compile-time type safety without the empirical latency and deadlock penalties of heavyweight ORMs, this specification establishes:
1. **Logical SQL Organization:** Centralizes and modularizes fragmented SQL queries and schemas into structured domain files under `cli/store/sql/` and dedicated query structures.
2. **Typed Repository Pattern:** Wraps database handles in typed domain repository structs implementing standardized interfaces, returning monadic results (`EntityResult[T]`, `ListResult[T]`) and structured `*appfault.AppError` instances while eliminating raw unhandled errors.
3. **AGY Prompt Management & Multi-Format Serialization:** Extends GitMap's CLI and services with end-to-end prompt operations: backup, restore, copy, export (JSON, pure SQL DDL/DML, SQLite `.db`), and import (JSON, `.db`).

## 2. User Request (Verbatim)
```markdown
# SQL Organization and Repository Pattern Implementation: high priority instruction, non-negotiable task

Can we please take these SQL and organize to its files and folders so that the database are logically there? And from the business logic, we just call the database, wrapping around some classes, like the repository pattern. Can it be confirmed that you can do that at least? So that we do not get into so much trouble and try to have a Database (Database (Db)). We check new returns, the app error, things like that. I think we should proceed further with this. Also check the CI/CD until it is green. X issues if you found any issue. Make sure that using the `gitmap`, we can actually back up running prompts using AGY, restore prompts, copy prompts, export as JSON, import from JSON, export as SQL and Database (Database (Db)).

## task: SQL Organization and Repository Pattern Implementation
## slug: sql-organization-and-repository-pattern-implementation

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/sql-organization-and-repository-pattern-implementation/ and enqueue plan task in .ai-memory/plans/sql-organization-and-repository-pattern-implementation.md (subtasks in .ai-memory/plans/subtasks/sql-organization-and-repository-pattern-implementation/) first
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Organize SQL files into logical files and folders for better database management.
6. Implement the repository pattern by wrapping database calls around classes.
7. Confirm the implementation to avoid issues and ensure a functional Database (Database (Db)).
8. Check new returns and app errors, and ensure CI/CD pipelines are green.
9. Address any issues found during the process.
10. Utilize `gitmap` to back up running prompts using AGY, restore prompts, copy prompts, export as JSON, import from JSON, export as SQL and Database (Database (Db)).

## Must follow and spawn agent using

@[\1](file;.agents/skills/\1)

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[\1](file;.agents/skills/\1) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

## 3. Architecture Blueprint

### 3.1 Logical SQL Organization
- Directory: `cli/store/sql/`
- Structure:
  - `cli/store/sql/repos.sql`: Schema definitions and queries for repositories, groups, and scan metadata.
  - `cli/store/sql/prompts.sql`: Schema definitions and queries for prompt backup batches, items, and restore ledgers.
  - `cli/store/sql/cluster.sql`: Schema definitions and queries for cluster nodes, runs, and execution results.
- Embedded via Go `embed.FS` in `cli/store/sql/sql.go` or exposed as typed query constants to prevent runtime path resolution failures.

### 3.2 Repository Pattern & Database Wrapping
- Interface `IRepository[T, ID]`: Defines standardized CRUD queries (`FindByID`, `List`, `Create`, `Update`, `Delete`).
- Struct `PromptBackupRepository`: Encapsulates database queries on `backup-prompts/sql.db` (tables `PromptBackupBatch`, `PromptBackupItem`, `PromptRestoreLedger`).
- Return convention: Monadic result envelopes (`EntityResult[T]`, `ListResult[T]`) or typed `(*T, *appfault.AppError)`. Raw `error` is strictly prohibited.
- Single-Writer SQLite Governance: All repositories enforce `SetMaxOpenConns(1)`, WAL journal mode, and 5000ms busy timeout.

### 3.3 AGY Prompt Backup & Multi-Format Serialization
- CLI commands mounted under `gitmap prompt` and `gitmap agy running-prompts`:
  - `backup`: Captures active/queued prompts into a new backup batch.
  - `restore <batch-id>`: Reconstitutes prompts into active agent queues.
  - `copy <source-id> <dest-id>`: Clones prompt records between namespaces or sessions.
  - `export --format json|sql|db --out <path>`:
    - `json`: Serializes full batch hierarchy to JSON envelope.
    - `sql`: Generates idempotent SQL `INSERT INTO` script.
    - `db`: Clones or exports SQLite database file directly.
  - `import --format json|db --file <path>`: Ingests external JSON or SQLite database backup.

## 4. Acceptance Criteria & Quality Gates
- **AC-1:** All SQL files organized logically in `cli/store/sql/` with 100% parameterization (`?` placeholders).
- **AC-2:** Repository pattern implemented with typed domain models and `*appfault.AppError` returns.
- **AC-3:** Zero raw errors leaking across repository boundaries; all errors mapped to domain error codes.
- **AC-4:** AGY running prompt backup, restore, copy, and multi-format export/import functional via CLI.
- **AC-5:** All Go files $\le 300$ LoC, 100% positive boolean convention (`is*`, `has*`), zero nested `if` statements.
- **AC-6:** CI/CD telemetry verified green via GitMap Pipeline-AI.
