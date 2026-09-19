# Subtask 01: Phase 2 Database Topology, Code Generator & Migration Engine

> **Assigned Worker:** Sub-Agent 1 (Systems & Database Engine)
> **Parent Plan:** `36-aum-polyglot-script-migration-and-llm-train-suite.md`
> **Target Scope:** Phase 2 of Spec 128 (`02-spec/21-app/128-aum-automation-suite-and-roadmap.md`)

---

## 1. Objectives & Scripts to Migrate

Migrate the following 4 Python scripts from `03-ai-scripts/` into compiled Go subcommands under `gitmap aum`:

1. **`18-codebase-topology-discoverer.py` ➔ `gitmap aum topology [dir]`**:
   - Analyzes repository entry points (`main.go`, `index.ts`, `cli/`), database schemas, CI workflows, and language distributions.
   - Caches discovery results in `.gitmap/data/automation/sql.db` (or binary data dir).
   - Flags: `--json`, `--refresh`.

2. **`30-db-struct-enum-generator.py` ➔ `gitmap aum db-generate [db-path]`**:
   - Inspects SQLite database tables and column types.
   - Generates strongly typed Go structs, TypeScript interfaces, and enum definitions adhering to PascalCase rules.
   - Flags: `--lang <go|ts|all>`, `--out <dir>`.

3. **`31-db-migration-runner.py` ➔ `gitmap aum db-migrate [db-path] [migrations-dir]`**:
   - Executes ordered SQL migration files with transaction rollback protection.
   - Records applied migrations in `_migrations` history table.
   - Flags: `--dry-run`, `--rollback <step>`.

4. **`34-schema-scanner.py` ➔ `gitmap aum schema-audit [db-path]`**:
   - Validates SQLite split-db schemas against repository naming conventions (PascalCase table names, positive booleans, primary keys).
   - Flags: `--json`, `--strict`.

---

## 2. Implementation Files & Architecture

- `cli/cmdautomation/topology.go` & `cli/cmdautomation/topology_cmd.go`
- `cli/cmdautomation/db_generate.go` & `cli/cmdautomation/db_generate_cmd.go`
- `cli/cmdautomation/db_migrate.go` & `cli/cmdautomation/db_migrate_cmd.go`
- `cli/cmdautomation/schema_audit.go` & `cli/cmdautomation/schema_audit_cmd.go`
- Unit tests: `cli/cmdautomation/phase2_test.go`

---

## 3. Strict Guidelines for Sub-Agent 1
- Functions <= 8–15 lines; files <= 100–200 lines.
- Affirmative booleans (`is*`, `has*`, `can*`).
- Single return types with `result.Result[T]` or `*apperror.AppError`.
- Zero nested ifs (nesting depth > 1 is forbidden).
- Do NOT run `go build` or `go test` in the loop.
