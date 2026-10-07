# Master Execution Plan: Views AI Inspection, Cursor Settings Guard & CI/CD Remediation

> **Plan Version:** 1.0.0  
> **Status:** Active  
> **Parent Task:** Task-231  
> **Associated Spec:** `02-spec/21-app/231-views-ai-inspection-cursor-settings-guard-and-cicd-remediation/01-architecture-spec.md`  

---

## 1. Goal & Objectives

1. Review and validate codebase modifications related to Views AI, multi-instance Antigravity discovery, and Cursor settings.
2. Ensure `cli/cmdcursor/cursor_settings.go` and `cli/cmdcursor/cursor_settings_test.go` conform to all style and quality guidelines.
3. Fix the nested `if` statement in `src/pages/Settings.tsx` to satisfy the TypeScript/AST and nested-if linter.
4. Clean up vertical blank line spacing in `cli/cmdagy/agy_instance_discovery.go`.
5. Run all local quality gates: `check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-relative-paths.py`, and Go unit tests.
6. Commit all changes, execute minor version bump (`v6.505.0` -> `v6.506.0`), publish GitHub release, and confirm green CI/CD pipeline health via `gitmap pe -t`.

---

## 2. Subtask Breakdown

- **Subtask 01:** `01-cursor-settings-guard-and-style.md`
  - Implement and format `requireCursorSettingsFile` in `cli/cmdcursor/cursor_settings.go`.
  - Add unit tests in `cli/cmdcursor/cursor_settings_test.go`.
  - Enforce vertical blank line spacing rules.

- **Subtask 02:** `02-settings-tsx-nested-if-and-linter-remediation.md`
  - Flatten nested `if` block in `src/pages/Settings.tsx`.
  - Enforce blank line conventions in `cli/cmdagy/agy_instance_discovery.go`.
  - Run all repository linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`).

- **Subtask 03:** `03-cicd-verification-minor-bump-and-release.md`
  - Run unit test suites across `cli/`.
  - Execute minor version bump (`v6.506.0`).
  - Perform atomic commit and release ceremony via `37-bump-version.py` and `29-release-orchestrator.py`.
  - Monitor CI/CD telemetry via `gitmap pe -t` to confirm green status.

---

## 3. Worker Allocation (A = 2, H = 2)

- **Worker 01:** Owns Subtask 01 (`cli/cmdcursor/cursor_settings.go`, `cli/cmdcursor/cursor_settings_test.go`).
- **Worker 02:** Owns Subtask 02 (`src/pages/Settings.tsx`, `cli/cmdagy/agy_instance_discovery.go`).
- **Lead Orchestrator:** Owns Subtask 03 (testing, version bump, atomic commit, release ceremony, CI telemetry verification).

---

## 4. Execution Tracking Ledger

| Subtask ID | File | Owner | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| Subtask-01 | `.ai-memory/plans/subtasks/231-views-ai-inspection-cursor-settings-guard-and-cicd-remediation/01-cursor-settings-guard-and-style.md` | Worker 01 | COMPLETED | `go test ./cmdcursor/...` PASS (10/10 tests) |
| Subtask-02 | `.ai-memory/plans/subtasks/231-views-ai-inspection-cursor-settings-guard-and-cicd-remediation/02-settings-tsx-nested-if-and-linter-remediation.md` | Worker 02 | COMPLETED | 0 linter violations, `go test ./cmdagy/...` PASS |
| Subtask-03 | `.ai-memory/plans/subtasks/231-views-ai-inspection-cursor-settings-guard-and-cicd-remediation/03-cicd-verification-minor-bump-and-release.md` | Lead | IN_PROGRESS | Starting release ceremony |
