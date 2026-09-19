# Subtask 04: Phase 5 Documentation, Spec Migration & Memory Consolidation

> **Assigned Worker:** Sub-Agent 2 (Quality & Release Engine)
> **Parent Plan:** `36-aum-polyglot-script-migration-and-llm-train-suite.md`
> **Target Scope:** Phase 5 of Spec 128 (`02-spec/21-app/128-aum-automation-suite-and-roadmap.md`)

---

## 1. Objectives & Scripts to Migrate

Migrate the following 4 Python scripts from `03-ai-scripts/` into compiled Go subcommands under `gitmap aum`:

1. **`09-cli-help-auditor.py` ➔ `gitmap aum help-audit [flags]`**:
   - Simulates CLI AST command definitions against markdown documentation in `cli/helptext/`.
   - Flags missing documentation files or undocumented subcommands.
   - Flags: `--strict`, `--json`.

2. **`20-plan-consolidator.py` & `32-deep-consolidator.py` ➔ `gitmap aum plan-consolidate [flags]`**:
   - Clusters completed tasks and subtasks in `.ai-memory/plans/completed/` into single milestone summary documents.
   - Drastically reduces file counts while preserving all verified outcomes.
   - Flags: `--dry-run`, `--threshold <N>`.

3. **`22-doc-path-linter.py` & `23-coding-guideline-path-consolidator.py` ➔ `gitmap aum doc-links [dir]`**:
   - Validates all markdown relative links across `02-spec/`, `.ai-memory/`, and `README.md`.
   - Checks that referenced targets exist on disk.
   - Flags: `--fix`, `--json`.

4. **`24-spec-path-migrator.py` & `25-repo-migrator.py` ➔ `gitmap aum spec-migrate [flags]`**:
   - Re-sequences specification file prefixes (e.g. `127-...`, `128-...`) and updates all cross-links across the codebase.
   - Flags: `--dry-run`, `--from <num>`, `--to <num>`.

---

## 2. Implementation Files & Architecture

- `cli/cmdautomation/help_audit.go` & `cli/cmdautomation/help_audit_cmd.go`
- `cli/cmdautomation/plan_consolidate.go` & `cli/cmdautomation/plan_consolidate_cmd.go`
- `cli/cmdautomation/doc_links.go` & `cli/cmdautomation/doc_links_cmd.go`
- `cli/cmdautomation/spec_migrate.go` & `cli/cmdautomation/spec_migrate_cmd.go`
- Unit tests: `cli/cmdautomation/phase5_test.go`

---

## 3. Strict Guidelines for Sub-Agent 2
- Functions <= 8–15 lines; files <= 100–200 lines.
- Affirmative booleans (`is*`, `has*`, `can*`).
- Single return types with `result.Result[T]` or `*apperror.AppError`.
- Zero nested ifs (nesting depth > 1 is forbidden).
- Do NOT run `go build` or `go test` in the loop.
