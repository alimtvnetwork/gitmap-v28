---
name: gitmap-split-db-engine
description: Autonomously audit, extend, migrate, and modify GitMap's three-tier SQLite Split-DB architecture (gitmap.db, installation.db, repodb/pipeline.db).
---

# GitMap Split-DB Engine Skill (`gitmap-split-db-engine`)

> **Invocation:** `@[gitmap-split-db-engine](file;.agents/skills/gitmap-split-db-engine)`

## Mission & Purpose
This skill provides authoritative architectural guidance, code navigation, and execution rules for designing, migrating, and maintaining GitMap's multi-database SQLite Split-DB storage tier.

---

## 1. Key Architectural Components & Code Map

| Database | Primary Location | Scope & Stored Entities |
|---|---|---|
| **Master Store (`gitmap.db`)** | `cli/db/`, `cli/cmddb/` | Central entity store: `Project`, `Alias`, `ScanFolder`, `Bookmark`, `GitProfile`, `Transaction`, `Group`, `Task`, `SSHKey`, `VSCodeProject`, `Workdir`, `ZipGroup`. |
| **Tooling Store (`installation.db`)** | `cli/db/`, `cli/installer/` | Standalone installed packages, toolchains, system binaries, installation metadata, and uninstallation recipes. |
| **Telemetry Store (`repodb/pipeline.db`)** | `cli/repodb/`, `cli/pipelinedb/` | Repository-scoped execution logs, test execution durations, pipeline history, and compact error reports. |
| **Code Generators** | `03-ai-scripts/35-db-struct-enum-generator.py` | Generates Go structs, GORM mappings, and PascalCase enum constants from SQLite schemas. |

---

## 2. Essential Commands

```bash
# List all registered split-databases and system SQLite files
gitmap db ls
gitmap db-listing

# Inspect database file sizes, table counts, and disk footprints
gitmap db sizes list

# Purge pipeline database telemetry and GitHub Actions run artifacts (0.0 GB quota)
gitmap pipeline purge

# Read and explore database tables, views, and schemas (zero-lock WAL mode)
gitmap db query "SELECT conversation_id, title, status FROM conversation_summaries LIMIT 10" --db agy
gitmap db query "SELECT RepoId, AbsolutePath FROM TrackedRepository LIMIT 10" --db master
gitmap db tables --db agy
gitmap db schema conversation_summaries --db agy
gitmap db columns conversation_summaries --db agy
gitmap db-read master ls
gitmap db-read tasks table Subtask --columns
gitmap db-read ".ai-memory/heatmap.db" "SELECT COUNT(*) FROM HeatmapActionLog"

# Mutate and insert records into split-databases
gitmap db-add tasks Subtask csv "SUB-01","my-task","PENDING"
gitmap db-edit tasks -table Subtask -id 42 csv "SUB-01","my-task","COMPLETED"
gitmap db-record-del tasks -table Subtask -id 42 -y

# Search across split-database tables using full-text wildcard
gitmap db-find master -table Project,Alias --find "enrichment" --json

# Launch dynamic REST API server with interactive visualizer and Swagger UI
gitmap db-rest master,tasks 8080,8081
```

For comprehensive CRUD mutations, join views, and REST API gateway documentation, refer to `@[gitmap](file;.agents/skills/gitmap)` Section 11 and Section 15.

---

## 3. Core Invariants & Engineering Guardrails

1. **PascalCase Singular Tables:** All tables must be named in singular PascalCase (e.g. `Project`, `Task`, `ScanFolder`), never plural (`Projects`).
2. **Standardized Integer Primary Keys:** Primary keys must be named `{TableName}Id` with `INTEGER PRIMARY KEY AUTOINCREMENT`. UUIDs are strictly banned.
3. **Mandatory Documentation Columns:**
   - Reference/entity tables require a nullable `Description TEXT NULL`.
   - Transactional/operational tables require nullable `Notes TEXT NULL` and `Comments TEXT NULL`.
   - No hardcoded non-null defaults for these fields. Join tables are exempt.
4. **WAL Mode & Concurrency:** All SQLite connections must enable Write-Ahead Logging (`PRAGMA journal_mode=WAL;`) with busy timeouts to avoid database lock contention.
5. **Zero Raw SQL in Application Logic:** Database interactions must utilize GORM or registered typed query builders. Raw string concatenation in SQL queries is prohibited.
6. **Positive Boolean Column Prefixes:** All boolean columns must use affirmative prefixes (`IsActive`, `HasCompleted`, `IsValid`).
