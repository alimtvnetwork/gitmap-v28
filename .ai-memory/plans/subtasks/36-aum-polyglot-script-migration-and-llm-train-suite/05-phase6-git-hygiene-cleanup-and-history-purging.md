# Subtask 05: Phase 6 Git Hygiene, Cleanup & History Purging

> **Assigned Worker:** Sub-Agent 2 (Quality & Release Engine)
> **Parent Plan:** `36-aum-polyglot-script-migration-and-llm-train-suite.md`
> **Target Scope:** Phase 6 of Spec 128 (`02-spec/21-app/128-aum-automation-suite-and-roadmap.md`)

---

## 1. Objectives & Scripts to Migrate

Migrate the following 4 Python scripts from `03-ai-scripts/` into compiled Go subcommands under `gitmap aum`:

1. **`19-artifact-remover.py` & `03-file-manipulator.py` ➔ `gitmap aum clean-artifacts [flags]`**:
   - Safely deletes build binaries (`gitmap.exe`, `*.syso`), test dumps, `.pytest_cache`, and temporary directories.
   - Preserves all tracked source files and configuration.
   - Flags: `--dry-run`, `--all`, `--verbose`.

2. **`27-git-changed-files.py` ➔ `gitmap aum changed-files [flags]`**:
   - Discovers all modified, staged, and untracked files relative to `origin/main` or specified commit.
   - Outputs machine-readable JSON for targeted linter feeds.
   - Flags: `--staged-only`, `--json`, `--base <ref>`.

3. **`30-purge-history.py` & `33-git-history-tracer-and-purge.py` ➔ `gitmap aum purge-history [flags]`**:
   - Safely traces and identifies large historical git blobs without rewriting recent commit lineage.
   - Flags: `--dry-run`, `--min-size-mb <N>`, `--json`.

4. **`26-go-code-formatter.py` & `10-encoding-normalizer.py` ➔ `gitmap aum format-go [dir]`**:
   - AST-aware Go formatter that organizes imports into standard, 3rd-party, and repository groups.
   - Enforces UTF-8 BOM-free encoding.
   - Flags: `--write`, `--check`.

---

## 2. Implementation Files & Architecture

- `cli/cmdautomation/clean_artifacts.go` & `cli/cmdautomation/clean_artifacts_cmd.go`
- `cli/cmdautomation/changed_files.go` & `cli/cmdautomation/changed_files_cmd.go`
- `cli/cmdautomation/purge_history.go` & `cli/cmdautomation/purge_history_cmd.go`
- `cli/cmdautomation/format_go.go` & `cli/cmdautomation/format_go_cmd.go`
- Unit tests: `cli/cmdautomation/phase6_test.go`

---

## 3. Strict Guidelines for Sub-Agent 2
- Functions <= 8–15 lines; files <= 100–200 lines.
- Affirmative booleans (`is*`, `has*`, `can*`).
- Single return types with `result.Result[T]` or `*apperror.AppError`.
- Zero nested ifs (nesting depth > 1 is forbidden).
- Do NOT run `go build` or `go test` in the loop.
