# Subtask 02: Phase 3 CI/CD Local Preflight & Multi-Core Checkers

> **Assigned Worker:** Sub-Agent 1 (Systems & Database Engine)
> **Parent Plan:** `36-aum-polyglot-script-migration-and-llm-train-suite.md`
> **Target Scope:** Phase 3 of Spec 128 (`02-spec/21-app/128-aum-automation-suite-and-roadmap.md`)

---

## 1. Objectives & Scripts to Migrate

Migrate the following 4 Python scripts from `03-ai-scripts/` into compiled Go subcommands under `gitmap aum`:

1. **`28-go-preflight-ci.py` & `06-cicd-local-runner.py` ➔ `gitmap aum preflight [flags]`**:
   - Executes parallel multi-core checks across CPU cores (gofmt check, relative paths, nested ifs, boolean naming, tests).
   - Real-time progress percentages and elapsed timing.
   - Flags: `--fail-fast`, `--json`, `--workers <N>`.

2. **`33-test-inventory-generator.py` ➔ `gitmap aum test-inventory [flags]`**:
   - Discovers all unit and integration tests across Go packages and script suites.
   - Generates `.ai-memory/test-inventory.json` with duration baselines and test counts.
   - Flags: `--out <path>`, `--refresh`.

3. **`34-purge-github-actions-artifacts.py` ➔ `gitmap aum purge-actions [flags]`**:
   - Queries GitHub Actions API and safely purges obsolete workflow run logs and artifacts.
   - Enforces 0.0 GB storage footprint (Rule R18).
   - Flags: `--dry-run`, `--older-than <days>`, `--repo <slug>`.

4. **`16-installer-smoke-tester.py` ➔ `gitmap aum smoke-test [flags]`**:
   - Cross-platform dry-run validator for installed tools, verifying `--version` and PATH availability.
   - Flags: `--tools <list>`, `--json`.

---

## 2. Implementation Files & Architecture

- `cli/cmdautomation/preflight.go` & `cli/cmdautomation/preflight_cmd.go`
- `cli/cmdautomation/test_inventory.go` & `cli/cmdautomation/test_inventory_cmd.go`
- `cli/cmdautomation/purge_actions.go` & `cli/cmdautomation/purge_actions_cmd.go`
- `cli/cmdautomation/smoke_test.go` & `cli/cmdautomation/smoke_test_cmd.go`
- Unit tests: `cli/cmdautomation/phase3_test.go`

---

## 3. Strict Guidelines for Sub-Agent 1
- Functions <= 8–15 lines; files <= 100–200 lines.
- Affirmative booleans (`is*`, `has*`, `can*`).
- Single return types with `result.Result[T]` or `*apperror.AppError`.
- Zero nested ifs (nesting depth > 1 is forbidden).
- Do NOT run `go build` or `go test` in the loop.
