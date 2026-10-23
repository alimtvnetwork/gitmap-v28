# Subtask Plan 03: gitmap Recent Commands Review

- **Subtask Slug:** `03-gitmap-recent-commands-review`
- **Parent Task:** `316-repository-management-and-testing-instructions`
- **Target Repository:** gitmap-v28

---

## 1. Overview & Objectives

Review the commands implemented, integrated, and created in the recent wave
(specs 271–315 in `02-spec/21-app/` and the live CLI surface) and identify
improvements worth making before the v6.600.0 release.

Known observation to verify and scope (2026-10-24): `gitmap cpf` in a clone with no
git identity configured surfaces a stack trace on first failure instead of a clean,
actionable error message; the real cause (missing `user.name` / `user.email`) only
becomes visible on a second run.

## 2. Method

- Search exclusively via GitMap (`gitmap aum search`, `gitmap find`, `gitmap cat`).
- Cross-check each recent spec's promised command surface against the implementation
  under `cli/`.

## 3. Acceptance Criteria

- Inventory of recent commands with a one-line purpose and implementation status.
- Improvement candidates ranked by user impact, each with relative-path evidence.
- The cpf identity/error-surfacing issue confirmed or refuted with code evidence.
