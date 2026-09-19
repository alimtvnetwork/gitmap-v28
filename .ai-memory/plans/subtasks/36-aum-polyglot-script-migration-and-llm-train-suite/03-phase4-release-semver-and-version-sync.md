# Subtask 03: Phase 4 Release, SemVer & Version Synchronization

> **Assigned Worker:** Sub-Agent 2 (Quality & Release Engine)
> **Parent Plan:** `36-aum-polyglot-script-migration-and-llm-train-suite.md`
> **Target Scope:** Phase 4 of Spec 128 (`02-spec/21-app/128-aum-automation-suite-and-roadmap.md`)

---

## 1. Objectives & Scripts to Migrate

Migrate the following 3 Python scripts from `03-ai-scripts/` into compiled Go subcommands under `gitmap aum`:

1. **`14-version-sync-checker.py` ➔ `gitmap aum version-sync [flags]`**:
   - Audits version string alignment across:
     - `cli/constants/constants_version.go`
     - `package.json`
     - `Cargo.toml`
     - `02-spec/00-overview.md`
   - Flags: `--fix`, `--target-version <v>`, `--json`.

2. **`29-release-bumper.py` & `29-release-orchestrator.py` ➔ `gitmap aum release-bump [major|minor|patch]`**:
   - Computes next semantic version number based on commit types (fix ➔ patch, feat ➔ minor, breaking ➔ major).
   - Updates changelog, synchronizes package files, and tags git repository.
   - Flags: `--dry-run`, `--tag`, `--push`.

3. **`38-milestone-consolidator.py` ➔ `gitmap aum milestones [milestone-id]`**:
   - Queries GitHub milestone issues and PRs.
   - Formats closed items into structured release notes for changelog inclusion.
   - Flags: `--json`, `--out <path>`.

---

## 2. Implementation Files & Architecture

- `cli/cmdautomation/version_sync.go` & `cli/cmdautomation/version_sync_cmd.go`
- `cli/cmdautomation/release_bump.go` & `cli/cmdautomation/release_bump_cmd.go`
- `cli/cmdautomation/milestones.go` & `cli/cmdautomation/milestones_cmd.go`
- Unit tests: `cli/cmdautomation/phase4_test.go`

---

## 3. Strict Guidelines for Sub-Agent 2
- Functions <= 8–15 lines; files <= 100–200 lines.
- Affirmative booleans (`is*`, `has*`, `can*`).
- Single return types with `result.Result[T]` or `*apperror.AppError`.
- Zero nested ifs (nesting depth > 1 is forbidden).
- Do NOT run `go build` or `go test` in the loop.
