---
name: gitmap-split-db-engine
description: Autonomously audit, extend, migrate, and modify GitMap's three-tier SQLite Split-DB architecture (gitmap.db, installation.db, repodb/pipeline.db).
---

# GitMap Split-DB Engine Skill (`gitmap-split-db-engine`)

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
# Check database health, connection status, and schema integrity
gitmap db status

# Run forward schema migrations across split databases
gitmap db migrate

# Reclaim space and auto-prune pipeline database telemetry (10MB ceiling)
gitmap db prune
gitmap pipeline purge

# Reset local database to clean state (with safety prompt)
gitmap db reset
```

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
